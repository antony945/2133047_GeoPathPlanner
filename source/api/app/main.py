from fastapi import FastAPI, Header, HTTPException, Depends, Query
from fastapi.responses import JSONResponse
from app.models import RoutingRequest, RoutingResponse
from app.token import verify_jwt_token
from app.kafka import KafkaService
from app.logger import logger
from app.config import RESPONSE_TIMEOUT_SECONDS, APP_NAME, APP_VERSION
from app.db import init_db, insert_routing_response, get_responses_by_user, db_healthcheck
from contextlib import asynccontextmanager

# Create kafka service
kafka = KafkaService()

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    logger.info("🚀 Starting up API service...")
    await kafka.start()
    logger.info("✅ Kafka client started successfully.")
    await init_db()
    logger.info("🗄️ Database initialized.")

    yield

    # Shutdown
    logger.info("🛑 Shutting down API service...")
    await kafka.stop()
    logger.info("✅ Kafka client stopped successfully.")

# Create api
app = FastAPI(title=APP_NAME, version=APP_VERSION, lifespan=lifespan)

# -------------------------------
# 🌍 Base endpoint
# -------------------------------
@app.get("/")
async def root():
    logger.info("📡 Root endpoint called.")
    return {
        "service": app.title,
        "status": "running",
        "message": "Welcome to the GeoPathPlanner Routing API 🚀",
        "version": app.version
    }

# -------------------------------
# 💓 Healthcheck endpoint
# -------------------------------
@app.get("/health")
async def health_check():
    logger.debug("🩺 Healthcheck endpoint called.")
    
    # Check Kafka status
    kafka_status = "unknown"
    try:
        # This assumes your Kafka client has a method to check if it's ready
        kafka_status = "ok" if kafka.consumer is not None and kafka.producer is not None else "not ready"
    except Exception:
        kafka_status = "error"

    # Check PostgreSQL / PostGIS status
    db_status = "ok" if await db_healthcheck() else "error"

    status = {
        "api": "ok",
        "kafka": kafka_status,
        "database": db_status
    }

    # If any service is not ok, return 503
    if kafka_status != "ok" or db_status != "ok":
        return JSONResponse(content=status, status_code=503)

    return status

# -------------------------------
# 🧭 Compute route endpoint
# -------------------------------
@app.post("/compute", response_model=RoutingResponse)
async def compute_route(
    request: RoutingRequest,
    user_id: str | None = Query(None),
    token_payload: dict | None = Depends(verify_jwt_token)  # <- inject verification
):
    logger.info(f"📨 Received routing request (request_id={request.request_id})")

    # Handle authenticated users (with token)
    # TODO: To check
    # if user_id and token_payload:
    #     jwt_sub = token_payload.get("sub")
    #     if user_id and str(jwt_sub) != str(user_id):
    #         logger.warning(f"🔒 User ID mismatch: token.sub={jwt_sub}, query.user_id={user_id}")
    #         raise HTTPException(status_code=403, detail="User ID mismatch")
    #     else:
    #         logger.debug(f"✅ JWT validated for user_id={jwt_sub or user_id}")
    # else:
    #     logger.info("⚠️ No JWT token provided — treating as anonymous request.")

    # Produce to Kafka
    logger.info(f"📤 Producing routing request {request.request_id} to Kafka...")
    await kafka.produce_request(request)

    # Wait for response with a timeout
    timeout = RESPONSE_TIMEOUT_SECONDS
    logger.info(f"⏳ Waiting for response for request_id={request.request_id} (timeout={timeout}s)...")
    response = await kafka.wait_for_response(request.request_id, timeout=timeout)

    if not response:
        logger.error(f"⛔ No response for request_id={request.request_id} within {timeout}s.")
        raise HTTPException(status_code=504, detail="No response received in time")

    logger.info(f"✅ Received routing response for request_id={request.request_id} (route_found={response.route_found})")

    # Persist in DB if user_id exists
    if user_id:
        try:
            await insert_routing_response(response, user_id=user_id)
            logger.info(f"🗄️ Routing response saved for user_id={user_id}")
        except Exception as e:
            logger.error(f"❌ Failed to save routing response to DB: {e}")
    
    return response

# -------------------------------
# 🕘 Retrieve past routes/history
# -------------------------------
@app.get("/history", response_model=list[RoutingResponse])
async def get_user_history(user_id: str = Query(..., description="User ID to retrieve route history for")):
    """
    Retrieve all past routing responses associated with a specific user_id.
    """
    logger.info(f"📥 Retrieving routing history for user_id={user_id}")

    try:
        db_entries = await get_responses_by_user(user_id)
        
        # TODO: Think what to do when no entries
        # if not db_entries:
        #     raise HTTPException(status_code=404, detail="No routing history found for this user")
        
        logger.info(f"✅ Found {len(db_entries)} routes for user_id={user_id}")
    except Exception as e:
        logger.error(f"❌ Failed to retrieve routing history from DB: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve routing history")

    # return [RoutingResponse(**entry.response.__dict__) for entry in db_entries]
    # return db_entries

    routes = []
    for entry in db_entries:
        try:
            # Parse the JSON stored in the "response" column
            routes.append(RoutingResponse.model_validate(entry.response))
        except Exception as e:
            logger.warning(f"⚠️ Failed to parse routing response for request_id={entry.request_id}: {e}")

    return routes
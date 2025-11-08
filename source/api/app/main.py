from fastapi import FastAPI, Header, HTTPException, Depends, Query, Path
from fastapi.responses import JSONResponse
from fastapi.encoders import jsonable_encoder
from fastapi.middleware.cors import CORSMiddleware
from app.models import RoutingRequest, RoutingResponse
from app.token import verify_jwt_token
from app.kafka import KafkaService
from app.logger import logger
from app.config import RESPONSE_TIMEOUT_SECONDS, APP_NAME, APP_VERSION
from app.db import init_db, insert_routing_response, get_routing_response, get_routing_responses_by_user, db_healthcheck, delete_routing_response
from contextlib import asynccontextmanager

# -------------------------------
# Standard response helper
# -------------------------------
def standard_response(data=None, status="success", message=None):
    return {"status": status, "message": message, "data": data}

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

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allows all origins
    allow_credentials=True,
    allow_methods=["*"],  # Allows all methods
    allow_headers=["*"],  # Allows all headers
)

# -------------------------------
# 🌍 Base endpoint
# -------------------------------
@app.get("/")
async def root():
    logger.info("📡 Root endpoint called.")
    return standard_response(
        data = {
            "service": app.title,
            "status": "running",
            "message": "Welcome to the GeoPathPlanner Routing API 🚀",
            "version": app.version
        }
    )

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

    # If any service is not ok, return 503 otherwise 200
    code = 503 if kafka_status != "ok" or db_status != "ok" else 200
    return JSONResponse(content=standard_response(status), status_code=code)

# -------------------------------
# 🧭 Compute route endpoint
# -------------------------------
@app.post("/routes/compute", response_model=RoutingResponse)
async def compute_route(
    request: RoutingRequest,
    user_id: str | None = Query(None),
    token_payload: dict = Depends(verify_jwt_token)
):
    logger.info(f"📨 Received routing request (request_id={request.request_id})")

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
    
    return JSONResponse(content=standard_response(jsonable_encoder(response)), status_code=200)

# -------------------------------
# 🕘 Retrieve past routes/history
# -------------------------------
@app.get("/routes", response_model=list[RoutingResponse])
async def get_user_history(
    user_id: str = Query(..., description="User ID to retrieve route history for"),
    token_payload: dict = Depends(verify_jwt_token)
):
    """
    Retrieve all past routing responses associated with a specific user_id.
    """
    logger.info(f"📥 Retrieving routing history for user_id={user_id}")

    try:
        db_entries = await get_routing_responses_by_user(user_id)
        
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

    return JSONResponse(content=standard_response(jsonable_encoder(routes)), status_code=200)

# -------------------------------
# 🗑️ Delete a past route
# -------------------------------
@app.delete("/routes/{request_id}")
async def delete_route(
    request_id: str = Path(..., description="Request ID of the route to delete"),
    user_id: str = Query(..., description="User ID associated with the route"),
    token_payload: dict = Depends(verify_jwt_token)
):
    logger.info(f"🗑️ Deleting routing response request_id={request_id} for user_id={user_id}")

    try:
        deleted = await delete_routing_response(request_id)
        if not deleted:
            raise HTTPException(status_code=404, detail="Route not found")
    except Exception as e:
        logger.error(f"❌ Failed to delete routing response: {e}")
        raise HTTPException(status_code=500, detail="Failed to delete routing response")

    return JSONResponse(
        content=standard_response(
            data=jsonable_encoder(deleted.response),
            message="Route successfully removed"
        ),
        status_code=200
    )

# -------------------------------
# 📦 Retrieve a single route
# -------------------------------
@app.get("/routes/{request_id}")
async def get_route(
    request_id: str = Path(..., description="Request ID of the route to retrieve"),
    user_id: str = Query(..., description="User ID associated with the route"),
    token_payload: dict = Depends(verify_jwt_token)
):
    logger.info(f"📦 Retrieving routing response request_id={request_id} for user_id={user_id}")

    try:
        route = await get_routing_response(request_id)
        if not route:
            raise HTTPException(status_code=404, detail="Route not found")

    except Exception as e:
        logger.error(f"❌ Failed to retrieve routing response: {e}")
        raise HTTPException(status_code=500, detail="Failed to retrieve routing response")
    
    return JSONResponse(
        content=standard_response(
            data=jsonable_encoder(route.response),
            message="Route successfully retrieved"
        ),
        status_code=200
    )

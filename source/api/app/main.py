from fastapi import FastAPI, Header, HTTPException, Depends, Query
from fastapi.responses import JSONResponse
from app.models import RoutingRequest, RoutingResponse
from app.token import verify_jwt_token
from app.kafka import KafkaService
from app.logger import logger
from app.config import RESPONSE_TIMEOUT_SECONDS

# Create kafka service
kafka = KafkaService()

# Create api
app = FastAPI(title="GeoPathPlanner API", version="1.0.0")

@app.on_event("startup")
async def startup_event():
    logger.info("🚀 Starting up API service...")
    await kafka.start()
    logger.info("✅ Kafka client started successfully.")

@app.on_event("shutdown")
async def shutdown_event():
    logger.info("🛑 Shutting down API service...")
    await kafka.stop()
    logger.info("✅ Kafka client stopped successfully.")

# -------------------------------
# 🌍 Base endpoint
# -------------------------------
@app.get("/")
async def root():
    logger.info("📡 Root endpoint called.")
    return {
        "service": "GeoPathPlanner API",
        "status": "running",
        "message": "Welcome to the GeoPathPlanner Routing API 🚀",
        "version": "1.0.0"
    }

# -------------------------------
# 💓 Healthcheck endpoint
# -------------------------------
@app.get("/health")
async def health_check():
    logger.debug("🩺 Healthcheck endpoint called.")
    kafka_status = "unknown"
    try:
        # This assumes your Kafka client has a method to check if it's ready
        kafka_status = "ok" if kafka.consumer is not None and kafka.producer is not None else "not ready"
    except Exception:
        kafka_status = "error"

    status = {
        "api": "ok",
        "kafka": kafka_status
    }

    # If any service is not ok, return 503
    if kafka_status != "ok":
        return JSONResponse(content=status, status_code=503)

# -------------------------------
# 🧭 Compute route endpoint
# -------------------------------
@app.post("/compute", response_model=RoutingResponse)
async def compute_route(
    request: RoutingRequest,
    user_id: str = Query(...),
    token_payload: dict = Depends(verify_jwt_token)
):
    logger.info(f"📨 Received routing request (request_id={request.request_id})")

    # Handle authenticated users (with token)
    # TODO: To check
    # if token_payload:
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

    # TODO: persist in DB only if user_id exists
    if user_id:
        logger.info(f"🗄️ Saving route for user_id={user_id} to DB (future implementation).")
    
    return response

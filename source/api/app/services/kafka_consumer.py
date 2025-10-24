import os
import json
import asyncio
import logging
from datetime import datetime

from aiokafka import AIOKafkaConsumer
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

# ── Your project imports (keep these paths exactly as in your API)
# If your paths differ, adjust the import paths accordingly.
from app.models.routing_request import RoutingRequest
from app.core.routing import compute_path  # your existing routing logic

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("GeoPathPlannerConsumer")

DATABASE_URL = os.getenv("DATABASE_URL")  # e.g. postgresql+asyncpg://postgres:postgres@db:5432/geopathplanner
KAFKA_BROKER_URL = os.getenv("KAFKA_BROKER_URL", "kafka:9092")
KAFKA_TOPIC = os.getenv("KAFKA_TOPIC", "routing_requests")

# Async SQLAlchemy
engine = create_async_engine(DATABASE_URL, echo=False, future=True)
SessionLocal = sessionmaker(bind=engine, class_=AsyncSession, expire_on_commit=False)


async def process_message(payload: dict):
    """
    Expected payload structure produced by /routing/compute:
      {
        "user_id": <int>,
        "waypoints": [[lon, lat], [lon, lat], ...],
        "obstacles": [<geojson-like polygons/segments>],  # optional
        ... (anything else your compute_path expects)
      }
    """
    logger.info(f"Processing payload: {payload}")
    async with SessionLocal() as session:
        try:
            user_id = payload["user_id"]
            waypoints = payload["waypoints"]
            obstacles = payload.get("obstacles", [])

            # Run your routing logic (sync or async). If it's sync, this is fine.
            result = compute_path(waypoints, obstacles)

            # Persist request row
            req = RoutingRequest(
                user_id=user_id,
                waypoints=waypoints,
                obstacles=obstacles,
                result=result,
                status="completed",
                timestamp=datetime.utcnow(),
            )
            session.add(req)
            await session.commit()
            logger.info(f"✅ Stored request for user_id={user_id}")

        except Exception as e:
            logger.exception(f"❌ Failed to process message: {e}")


async def consume():
    consumer = AIOKafkaConsumer(
        KAFKA_TOPIC,
        bootstrap_servers=KAFKA_BROKER_URL,
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
        group_id="geo_routing_consumers",
        enable_auto_commit=True,
        auto_offset_reset="earliest",
    )
    await consumer.start()
    logger.info(f"Kafka consumer listening on topic '{KAFKA_TOPIC}'...")
    try:
        async for msg in consumer:
            await process_message(msg.value)
    finally:
        await consumer.stop()


if __name__ == "__main__":
    asyncio.run(consume())

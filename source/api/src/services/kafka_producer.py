import asyncio
import json
import logging
from datetime import datetime
from aiokafka import AIOKafkaProducer

producer = None

async def get_producer():
    """Create or return a global Kafka producer"""
    global producer
    if producer is None:
        producer = AIOKafkaProducer(
            bootstrap_servers="kafka:9092",
            value_serializer=lambda v: json.dumps(v).encode("utf-8"),
        )
        await producer.start()
    return producer

async def send_request_message(request_id: str, user_id: str | None, payload: dict):
    """Send the 'request_to_send' event to Kafka"""
    try:
        p = await get_producer()
        message = {
            "type": "request_to_send",
            "request_id": request_id,
            "user_id": user_id,
            "timestamp": datetime.utcnow().isoformat(),
            "payload": payload,
        }
        await p.send_and_wait("routing_requests", message)
        logging.info(f"[Kafka] sent request {request_id}")
    except Exception as e:
        logging.error(f"[Kafka] send failed: {e}")

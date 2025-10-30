# app/kafka_service.py
import asyncio
import json
import logging
from aiokafka import AIOKafkaProducer, AIOKafkaConsumer
from app.models import RoutingRequest, RoutingResponse
from app.config import KAFKA_REQUEST_TOPIC, KAFKA_RESPONSE_TOPIC, KAFKA_BROKERS, KAFKA_CONSUMER_GROUP
from app.logger import logger

class KafkaService:
    def __init__(self):
        self.producer: AIOKafkaProducer | None = None
        self.consumer: AIOKafkaConsumer | None = None
        self._response_futures: dict[str, asyncio.Future] = {}

    async def start(self):
        """Initialize and start producer and consumer."""
        self.producer = AIOKafkaProducer(
            bootstrap_servers=KAFKA_BROKERS,
            value_serializer=lambda v: json.dumps(v).encode("utf-8"),
        )

        self.consumer = AIOKafkaConsumer(
            KAFKA_RESPONSE_TOPIC,
            bootstrap_servers=KAFKA_BROKERS,
            value_deserializer=None,
            group_id=KAFKA_CONSUMER_GROUP,
            enable_auto_commit=True,
            auto_offset_reset="earliest",
        )

        await self.producer.start()
        await self.consumer.start()

        # Background task for consuming responses
        asyncio.create_task(self._consume_responses())
        logger.info("Kafka producer and consumer started.")

    async def stop(self):
        """Stop producer and consumer."""
        if self.producer:
            await self.producer.stop()
        if self.consumer:
            await self.consumer.stop()
        logger.info("Kafka service stopped.")

    async def produce_request(self, request: RoutingRequest):
        """Send a RoutingRequest message to Kafka."""
        message = request.model_dump_json()
        await self.producer.send_and_wait(KAFKA_REQUEST_TOPIC, message)
        logger.info(f"Produced routing request: {request.request_id}")

    async def wait_for_response(self, request_id: str, timeout: float = 10.0) -> RoutingResponse | None:
        """Wait for the response that matches a given request_id."""
        loop = asyncio.get_running_loop()
        fut = loop.create_future()
        self._response_futures[request_id] = fut

        try:
            return await asyncio.wait_for(fut, timeout)
        except asyncio.TimeoutError:
            logger.warning(f"Timeout waiting for response to request_id={request_id}")
            return None
        finally:
            self._response_futures.pop(request_id, None)

    async def _consume_responses(self):
        """Background task to continuously consume routing responses."""
        logger.info(f"Kafka consumer listening on '{KAFKA_RESPONSE_TOPIC}'...")
        async for msg in self.consumer:
            try:
                # msg.value is bytes → decode to str
                payload_str = msg.value.decode("utf-8")
                # Pydantic parses JSON (including datetime & nested models)
                response = RoutingResponse.model_validate_json(payload_str)
            except Exception as e:
                logger.warning(f"Invalid message format: {e}")
                continue

            fut = self._response_futures.get(response.request_id)
            if fut and not fut.done():
                try:
                    fut.set_result(response)
                    logger.info(f"Delivered response for request_id={response.request_id}")
                except asyncio.InvalidStateError:
                    logger.warning(f"Future for request_id={response.request_id} was already done or cancelled")
            else:
                logger.info(f"No waiting future for request_id={response.request_id}, ignoring message")
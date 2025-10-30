import os
from dotenv import load_dotenv

# Load variables from a .env file if present (useful in local dev)
load_dotenv()

# ===============================
# ✅ Kafka Configuration
# ===============================
KAFKA_BROKERS = os.getenv("KAFKA_BROKERS", "kafka:9093")
# The topic to which this microservice sends routing requests
KAFKA_REQUEST_TOPIC = os.getenv("KAFKA_REQUEST_TOPIC", "routing-requests")
# The topic this microservice listens on for routing responses
KAFKA_RESPONSE_TOPIC = os.getenv("KAFKA_RESPONSE_TOPIC", "routing-responses")
# Optional: Kafka consumer group (you can override per microservice)
KAFKA_CONSUMER_GROUP = os.getenv("KAFKA_API_CONSUMER_GROUP_ID", "geo_routing_api_group")

# ===============================
# ✅ JWT / Authentication
# ===============================
JWT_SECRET_KEY = os.getenv("JWT_SECRET_KEY", "default_secret_key")  # ⚠️ override in prod
JWT_ALGORITHM = os.getenv("JWT_ALGORITHM", "HS256")
JWT_ISSUER = os.getenv("JWT_ISSUER", "geopathplanner-auth-service")

# ===============================
# ✅ Application Settings
# ===============================
APP_NAME = os.getenv("APP_NAME", "GeoPathPlannerAPI")
APP_ENV = os.getenv("APP_ENV", "dev")
APP_PORT = int(os.getenv("APP_PORT", "8000"))
LOG_LEVEL = os.getenv("LOG_LEVEL", "INFO")
RESPONSE_TIMEOUT_SECONDS = float(os.getenv("RESPONSE_TIMEOUT_SECONDS", "20.0"))
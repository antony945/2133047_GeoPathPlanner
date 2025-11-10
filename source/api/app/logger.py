import logging
import os
from logging.handlers import RotatingFileHandler
from app.config import LOG_LEVEL, APP_ENV

# Create logs directory if missing
LOG_DIR = "logs"
os.makedirs(LOG_DIR, exist_ok=True)

# Parse log level from env
numeric_level = getattr(logging, LOG_LEVEL.upper(), logging.INFO)

# Define log file paths
INFO_LOG_FILE = os.path.join(LOG_DIR, "info.log")
ERROR_LOG_FILE = os.path.join(LOG_DIR, "error.log")
DEBUG_LOG_FILE = os.path.join(LOG_DIR, "debug.log")

# Define format
LOG_FORMAT = "[%(asctime)s] [%(levelname)s] [%(name)s] %(message)s"

# Configure root logger
logging.basicConfig(
    level=numeric_level,
    format=LOG_FORMAT,
    handlers=[
        RotatingFileHandler(INFO_LOG_FILE, maxBytes=5_000_000, backupCount=5),
        RotatingFileHandler(ERROR_LOG_FILE, maxBytes=5_000_000, backupCount=5),
        RotatingFileHandler(DEBUG_LOG_FILE, maxBytes=5_000_000, backupCount=5),
        logging.StreamHandler(),  # Print to console too
    ],
)

# Create logger instance for use in app
logger = logging.getLogger("geo_api")

# Adjust per-handler filtering
for handler in logger.handlers:
    if isinstance(handler, RotatingFileHandler):
        if "error" in handler.baseFilename:
            handler.setLevel(logging.ERROR)
        elif "debug" in handler.baseFilename:
            handler.setLevel(logging.DEBUG)
        else:
            handler.setLevel(logging.INFO)

logger.info(f"Logger initialized with level={LOG_LEVEL} in env={APP_ENV}")

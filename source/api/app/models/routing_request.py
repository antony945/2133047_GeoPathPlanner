from sqlalchemy import Column, Integer, String, DateTime, ForeignKey
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.orm import declarative_base, relationship
from datetime import datetime

# If you already have a Base in your project (e.g., app.db.base),
# import that Base instead of creating a new one.
try:
    from app.db.base import Base
except Exception:
    Base = declarative_base()

class RoutingRequest(Base):
    __tablename__ = "routing_requests"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False, index=True)

    # Store arrays/objects directly (Postgres JSONB)
    waypoints = Column(JSONB, nullable=False)
    obstacles = Column(JSONB, nullable=True)
    result = Column(JSONB, nullable=True)

    status = Column(String, default="pending", nullable=False)
    timestamp = Column(DateTime, default=datetime.utcnow, nullable=False)

    # Optional: if you have a User model
    # user = relationship("User", back_populates="routing_requests")

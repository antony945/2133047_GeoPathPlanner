from pydantic import BaseModel
from datetime import datetime
from typing import Any, Optional

class RoutingRequestOut(BaseModel):
    id: int
    user_id: int
    waypoints: Any
    obstacles: Optional[Any]
    result: Optional[Any]
    status: str
    timestamp: datetime

    class Config:
        from_attributes = True  # Pydantic v2 (ORM mode)

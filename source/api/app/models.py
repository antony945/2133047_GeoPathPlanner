from pydantic import BaseModel
from geojson_pydantic import Feature
from datetime import datetime
from typing import Optional, List

class RoutingRequest(BaseModel):
    request_id: str
    waypoints: List[Feature]
    constraints: Optional[List[Feature]] = None
    search_volume: Optional[Feature] = None
    parameters: Optional[dict] = None
    received_at: datetime

class RoutingResponse(RoutingRequest):
    route_found: bool
    route: Optional[List[Feature]] = None
    cost_km: Optional[float] = None
    message: Optional[str] = None
    completed_at: datetime
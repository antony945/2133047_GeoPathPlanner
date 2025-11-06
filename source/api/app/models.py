from pydantic import BaseModel, ConfigDict
from geojson_pydantic import Feature
from datetime import datetime
from typing import Optional, List

class RoutingRequest(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    request_id: str
    waypoints: List[Feature]
    constraints: Optional[List[Feature]] = None
    search_volume: Feature
    parameters: Optional[dict] = None
    received_at: datetime

class RoutingResponse(RoutingRequest):
    route_found: bool
    route: Optional[List[Feature]] = None
    cost_km: Optional[float] = None
    message: Optional[str] = None
    completed_at: datetime
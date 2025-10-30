from typing import Optional, List, Literal

from pydantic import BaseModel, Field


# ----------------- Schemas -----------------
Unit = Literal["m", "ft"]


class LatLonAlt(BaseModel):
    lat: float
    lon: float
    alt: float
    alt_unit: Unit = "m"


class Constraint(BaseModel):
    polygon: dict
    min_alt: Optional[float] = None
    max_alt: Optional[float] = None
    alt_unit: Optional[Unit] = None
    name: Optional[str] = None


class SearchVolume(BaseModel):
    bbox: tuple[float, float, float, float]
    margin_m: int


class AlgorithmParams(BaseModel):
    algorithm: str = Field(default="RRT*")
    iterations: Optional[int] = 10000
    speed_mps: Optional[float] = 12


class RoutingInput(BaseModel):
    waypoints: List[LatLonAlt]
    constraints: List[Constraint] = []
    search_volume: Optional[SearchVolume] = None
    params: AlgorithmParams
    allow_contains: bool = False


class RoutingResponse(BaseModel):
    request_id: str
    route_found: bool
    route_output_waypoints: Optional[List[LatLonAlt]] = None
    route_output_message: Optional[str] = None
    total_distance_km: Optional[float] = None
    compute_ms: Optional[int] = None

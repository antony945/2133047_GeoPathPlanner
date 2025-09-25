import os
import uuid
from datetime import datetime, timedelta, timezone
from typing import Optional, List, Literal

from fastapi import FastAPI, Depends, HTTPException, status, Header
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import OAuth2PasswordBearer
from jose import jwt, JWTError
from passlib.context import CryptContext
from pydantic import BaseModel, Field

# --- SQLAlchemy / Postgres (async) ---
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column
from sqlalchemy.dialects.postgresql import UUID as PG_UUID, JSONB
from sqlalchemy import String, Boolean, Integer, Text, TIMESTAMP, func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

# ----------------- Config -----------------
SECRET_KEY = os.getenv("SECRET_KEY", "dev-change-me")
ALGORITHM = os.getenv("ALGORITHM", "HS256")
ACCESS_TOKEN_EXPIRE_MINUTES = int(os.getenv("ACCESS_TOKEN_EXPIRE_MINUTES", "60"))
DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgresql+asyncpg://postgres:postgres@localhost:5432/geopathplanner",
)

# ----------------- Auth helpers -----------------
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="/auth/login")

def create_access_token(user_id: str) -> str:
    expire = datetime.now(timezone.utc) + timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    return jwt.encode({"sub": user_id, "exp": expire}, SECRET_KEY, algorithm=ALGORITHM)

def decode_token(token: str) -> Optional[str]:
    try:
        return jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM]).get("sub")
    except JWTError:
        return None

# ----------------- DB setup -----------------
class Base(DeclarativeBase):
    pass

engine = create_async_engine(DATABASE_URL, echo=False, future=True)
AsyncSessionLocal = async_sessionmaker(bind=engine, class_=AsyncSession, expire_on_commit=False)

async def get_db() -> AsyncSession:
    async with AsyncSessionLocal() as session:
        yield session

# ----------------- Models -----------------
class User(Base):
    __tablename__ = "users"
    id: Mapped[uuid.UUID] = mapped_column(PG_UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    username: Mapped[str] = mapped_column(String(255), unique=True, nullable=False)
    password_hash: Mapped[str] = mapped_column(String(255), nullable=False)

class RoutingRequest(Base):
    __tablename__ = "routing_requests"
    request_id: Mapped[uuid.UUID] = mapped_column(PG_UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    request_datetime: Mapped[datetime] = mapped_column(TIMESTAMP(timezone=True), server_default=func.now())
    user_id: Mapped[Optional[uuid.UUID]] = mapped_column(PG_UUID(as_uuid=True), nullable=True)
    input_json: Mapped[dict] = mapped_column(JSONB, nullable=False)
    route_found: Mapped[bool] = mapped_column(Boolean, nullable=False)
    route_output: Mapped[Optional[dict]] = mapped_column(JSONB)
    route_message: Mapped[Optional[str]] = mapped_column(Text)
    compute_ms: Mapped[Optional[int]] = mapped_column(Integer)

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

class UserCreate(BaseModel):
    username: str
    password: str

class UserLogin(BaseModel):
    username: str
    password: str

class UserOut(BaseModel):
    id: str
    username: str

# ----------------- Helpers -----------------
def normalize_alt(value: float, unit: str) -> float:
    return value if unit == "m" else value * 0.3048

def haversine_km(lat1, lon1, lat2, lon2) -> float:
    from math import radians, sin, cos, sqrt, atan2
    R = 6371.0088
    dlat = radians(lat2 - lat1); dlon = radians(lon2 - lon1)
    a = sin(dlat/2)**2 + cos(radians(lat1))*cos(radians(lat2))*sin(dlon/2)**2
    return 2*R*atan2(sqrt(a), sqrt(1-a))

def path_length_km(points: List[LatLonAlt]) -> float:
    total = 0.0
    for i in range(len(points)-1):
        a, b = points[i], points[i+1]
        total += haversine_km(a.lat, a.lon, b.lat, b.lon)
    return round(total, 3)

async def get_current_user(db: AsyncSession = Depends(get_db), token: str = Depends(oauth2_scheme)) -> User:
    sub = decode_token(token)
    if not sub:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid token")
    uid = uuid.UUID(sub)
    row = await db.execute(select(User).where(User.id == uid))
    user = row.scalar_one_or_none()
    if not user:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="User not found")
    return user

# ----------------- FastAPI app -----------------
app = FastAPI(title="GeoPathPlanner API")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.on_event("startup")
async def on_startup():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)

# -------- Root (health) --------
@app.get("/", tags=["health"])
def root():
    return {"ok": True, "msg": "API is running"}

# -------- Auth --------
@app.post("/auth/register", status_code=201, tags=["auth"])
async def register(payload: UserCreate, db: AsyncSession = Depends(get_db)):
    exists = await db.execute(select(User).where(User.username == payload.username))
    if exists.scalar_one_or_none():
        raise HTTPException(status_code=400, detail="Username already exists")
    u = User(username=payload.username, password_hash=pwd_context.hash(payload.password))
    db.add(u); await db.commit()
    return {"message": "registered"}

@app.post("/auth/login", tags=["auth"])
async def login(payload: UserLogin, db: AsyncSession = Depends(get_db)):
    q = await db.execute(select(User).where(User.username == payload.username))
    u = q.scalar_one_or_none()
    if not u or not pwd_context.verify(payload.password, u.password_hash):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid credentials")
    return {"access_token": create_access_token(str(u.id)), "token_type": "bearer"}

@app.get("/auth/me", response_model=UserOut, tags=["auth"])
async def me(current_user: User = Depends(get_current_user)):
    return {"id": str(current_user.id), "username": current_user.username}

# -------- Routing (anonymous allowed) --------
@app.post("/routing/compute", response_model=RoutingResponse, tags=["routing"])
async def compute_route(
    payload: RoutingInput,
    db: AsyncSession = Depends(get_db),
    authorization: Optional[str] = Header(default=None),  # optional Bearer token
):
    # normalize waypoints (ft -> m)
    norm = [LatLonAlt(lat=w.lat, lon=w.lon, alt=normalize_alt(w.alt, w.alt_unit), alt_unit="m")
            for w in payload.waypoints]
    if len(norm) < 2:
        return {"request_id": str(uuid.uuid4()), "route_found": False,
                "route_output_message": "Need at least 2 waypoints"}

    # trivial route (placeholder): through the waypoints
    dist_km = path_length_km(norm)
    resp = {"request_id": str(uuid.uuid4()), "route_found": True,
            "route_output_waypoints": norm, "total_distance_km": dist_km, "compute_ms": 5}

    # if Authorization: Bearer <token> present, save to history
    user_id: Optional[uuid.UUID] = None
    if authorization and authorization.lower().startswith("bearer "):
        sub = decode_token(authorization.split(" ", 1)[1])
        if sub:
            try:
                user_id = uuid.UUID(sub)
            except Exception:
                user_id = None

    if user_id:
        rr = RoutingRequest(
            user_id=user_id, input_json=payload.model_dump(),
            route_found=True, route_output=resp, route_message=None, compute_ms=5
        )
        db.add(rr); await db.commit()

    return resp

# -------- History (auth only) --------
@app.get("/routes", tags=["history"])
async def list_routes(db: AsyncSession = Depends(get_db), current_user: User = Depends(get_current_user)):
    q = await db.execute(
        select(RoutingRequest).where(RoutingRequest.user_id == current_user.id)
        .order_by(RoutingRequest.request_datetime.desc()).limit(50)
    )
    items = [{
        "request_id": str(r.request_id),
        "when": r.request_datetime.isoformat() if r.request_datetime else None,
        "route_found": r.route_found,
        "compute_ms": r.compute_ms
    } for r in q.scalars().all()]
    return {"items": items}

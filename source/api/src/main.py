import os
import uuid
from typing import Optional
from services.kafka_producer import send_request_message
from fastapi import FastAPI, Depends, HTTPException, status, Header
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import OAuth2PasswordBearer
from jose import jwt, JWTError

# --- SQLAlchemy / Postgres (async) ---
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column
from sqlalchemy.dialects.postgresql import UUID as PG_UUID, JSONB
from sqlalchemy import String, Boolean, Integer, Text, TIMESTAMP, func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from schemas.schemas import RoutingResponse

# ----------------- Config -----------------
SECRET_KEY = os.getenv("SECRET_KEY", "dev-change-me")
ALGORITHM = os.getenv("ALGORITHM", "HS256")
ACCESS_TOKEN_EXPIRE_MINUTES = int(os.getenv("ACCESS_TOKEN_EXPIRE_MINUTES", "60"))
DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgresql+asyncpg://postgres:postgres@localhost:5432/geopathplanner",
)

# # ----------------- Auth helpers -----------------
# pwd_context = CryptContext(
#     schemes=["bcrypt"],
#     bcrypt__rounds=12,
#     bcrypt__ident="2b",
#     deprecated="auto",
# )

# oauth2_scheme = OAuth2PasswordBearer(tokenUrl="/auth/login")


# def create_access_token(user_id: str) -> str:
#     expire = datetime.now(timezone.utc) + timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
#     return jwt.encode({"sub": user_id, "exp": expire}, SECRET_KEY, algorithm=ALGORITHM)


# def decode_token(token: str) -> Optional[str]:
#     try:
#         return jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM]).get("sub")
#     except JWTError:
#         return None


# ----------------- DB setup -----------------
class Base(DeclarativeBase):
    pass


# engine = create_async_engine(DATABASE_URL, echo=False, future=True)
# AsyncSessionLocal = async_sessionmaker(bind=engine, class_=AsyncSession, expire_on_commit=False)


# async def get_db() -> AsyncSession:
#     async with AsyncSessionLocal() as session:
#         yield session


# # ----------------- Models -----------------

# class RoutingRequest(Base):
#     __tablename__ = "routing_requests"
#     request_id: Mapped[uuid.UUID] = mapped_column(PG_UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
#     request_datetime: Mapped[datetime] = mapped_column(TIMESTAMP(timezone=True), server_default=func.now())
#     user_id: Mapped[Optional[uuid.UUID]] = mapped_column(PG_UUID(as_uuid=True), nullable=True)
#     input_json: Mapped[dict] = mapped_column(JSONB, nullable=False)
#     route_found: Mapped[bool] = mapped_column(Boolean, nullable=False)
#     route_output: Mapped[Optional[dict]] = mapped_column(JSONB)
#     route_message: Mapped[Optional[str]] = mapped_column(Text)
#     compute_ms: Mapped[Optional[int]] = mapped_column(Integer)

# ----------------- FastAPI app -----------------
app = FastAPI(title="GeoPathPlanner API")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# @app.on_event("startup")
# async def on_startup():
#     async with engine.begin() as conn:
#         await conn.run_sync(Base.metadata.create_all)


# -------- Root (health) --------
@app.get("/", tags=["health"])
def root():
    return {"ok": True, "msg": "API is running"}

# # -------- Routing (anonymous allowed) --------
@app.post("/routing/compute", response_model=RoutingResponse, tags=["routing"])
async def compute_route(
    payload: RoutingInput,
    db: AsyncSession = Depends(get_db),
    authorization: Optional[str] = Header(default=None),
):
    request_id = str(uuid.uuid4())

    norm = [
        LatLonAlt(lat=w.lat, lon=w.lon, alt=normalize_alt(w.alt, w.alt_unit), alt_unit="m")
        for w in payload.waypoints
    ]

    if len(norm) < 2:
        return {
            "request_id": request_id,
            "route_found": False,
            "route_output_message": "Need at least 2 waypoints",
        }

    dist_km = path_length_km(norm)
    resp = {
        "request_id": request_id,
        "route_found": True,
        "route_output_waypoints": norm,
        "total_distance_km": dist_km,
        "compute_ms": 5,
    }

    # Extract user_id first
    user_id: Optional[uuid.UUID] = None
    # Extract user_id from Bearer token (if provided)
    user_id: Optional[uuid.UUID] = None
    if authorization and authorization.lower().startswith("bearer "):
        token = authorization.split(" ", 1)[1]
        sub = decode_token(token)
        if sub:
            try:
                user_id = uuid.UUID(sub)
            except Exception:
                user_id = None

    # Always send the Kafka message
    await send_request_message(request_id, str(user_id) if user_id else None, payload.model_dump())

    # Always store the routing request (even if user_id is None)
    serializable_resp = {
        **resp,
        "route_output_waypoints": [w.model_dump() for w in resp["route_output_waypoints"]],
    }

    rr = RoutingRequest(
        request_id=uuid.UUID(request_id),
        user_id=user_id,
        input_json=payload.model_dump(),
        route_found=True,
        route_output=serializable_resp,
        route_message=None,
        compute_ms=5,
    )
    db.add(rr)
    await db.commit()

    return resp

# # -------- History (auth only) --------
# @app.get("/routes", tags=["history"])
# async def list_routes(db: AsyncSession = Depends(get_db), current_user: User = Depends(get_current_user)):
#     q = await db.execute(
#         select(RoutingRequest)
#         .where(RoutingRequest.user_id == current_user.id)
#         .order_by(RoutingRequest.request_datetime.desc())
#         .limit(50)
#     )
#     items = [
#         {
#             "request_id": str(r.request_id),
#             "when": r.request_datetime.isoformat() if r.request_datetime else None,
#             "route_found": r.route_found,
#             "compute_ms": r.compute_ms,
#         }
#         for r in q.scalars().all()
#     ]
#     return {"items": items}

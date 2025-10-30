import os
from typing import List, Optional
from datetime import datetime, timezone
from app.logger import logger
from sqlalchemy import select, delete
from fastapi.encoders import jsonable_encoder

from sqlalchemy import Column, String, JSON, DateTime
from sqlalchemy.ext.asyncio import AsyncEngine, create_async_engine, AsyncSession
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from sqlalchemy.sql import func
# from geoalchemy2 import Geometry

from app.models import RoutingResponse

from app.config import DATABASE_URL

# --- DATABASE SETUP ---
engine: AsyncEngine = create_async_engine(DATABASE_URL, echo=True)
async_session = sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)

Base = declarative_base()

# --- TABLE DEFINITION ---
class RoutingResponseDB(Base):
    __tablename__ = "routing_responses"

    request_id = Column(String, primary_key=True)
    user_id = Column(String, index=True, nullable=False)
    # waypoints = Column(JSON, nullable=False)
    # constraints = Column(JSON, nullable=True)
    # search_volume = Column(JSON, nullable=False)
    # parameters = Column(JSON, nullable=True)
    # received_at = Column(DateTime, nullable=True)

    # Store the full response as JSON
    response = Column(JSON, nullable=False)

    # TODO: Think about storing additional parameters or not

    # # Optionally store geometries for PostGIS queries
    # # Waypoints as MULTIPOINT (from list of Point features)
    # waypoints_geom = Column(Geometry(geometry_type='MULTIPOINT', srid=4326), nullable=True)
    
    # # Constraints as MULTIPOLYGON (from list of Polygon features)
    # constraints_geom = Column(Geometry(geometry_type='MULTIPOLYGON', srid=4326), nullable=True)

    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())

# --- DB INIT FUNCTION ---
async def init_db():
    """Create tables in database."""
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)

# --- DB HEALTHCHECK ---
async def db_healthcheck() -> bool:
    """
    Check if the database connection works by executing a simple SELECT 1.
    Returns True if OK, False otherwise.
    """
    try:
        async with async_session() as session:
            # Modern, recommended syntax using select()
            result = await session.execute(select(1))
            _ = result.scalar()  # Consume the result to ensure execution
        return True
    except Exception as e:
        logger.error(f"DB healthcheck failed: {e}")
        return False

# --- HELPER FUNCTIONS ---
async def insert_routing_response(response: RoutingResponse, user_id: str):
    """Insert a routing response into the database."""
    async with async_session() as session:
        db_entry = RoutingResponseDB(
            request_id=response.request_id,
            user_id=user_id,
            response=jsonable_encoder(response),
            # waypoints_geom=convert_waypoints_to_multipoint(response.waypoints),
            # constraints_geom=convert_constraints_to_multipolygon(response.constraints)
        )
        session.add(db_entry)
        await session.commit()

async def get_responses_by_user(user_id: str) -> List[RoutingResponseDB]:
    """Return all routing responses associated with a specific user_id."""
    async with async_session() as session:
        result = await session.execute(
            select(RoutingResponseDB)
                .where(RoutingResponseDB.user_id == user_id)
                .order_by(RoutingResponseDB.created_at.desc())
        )
        return result.scalars().all()

async def delete_routing_response(request_id: str) -> bool:
    """
    Delete a routing response by request_id and user_id.
    Returns True if a row was deleted, False if no matching row was found.
    """
    async with async_session() as session:
        try:
            stmt = delete(RoutingResponseDB).where(
                RoutingResponseDB.request_id == request_id
                # RoutingResponseDB.user_id == user_id
            ).returning(RoutingResponseDB.request_id)

            result = await session.execute(stmt)
            await session.commit()

            deleted_request_id = result.scalar()
            return deleted_request_id is not None

        except Exception as e:
            logger.error(f"Failed to delete routing response request_id={request_id}: {e}")
            return False

# # --- GEOMETRY CONVERSION HELPERS ---
# from shapely.geometry import MultiPoint, MultiPolygon, shape
# from geoalchemy2.shape import from_shape

# def convert_waypoints_to_multipoint(waypoints) -> Optional[str]:
#     """Convert list of GeoJSON Point features into a PostGIS MULTIPOINT."""
#     if not waypoints:
#         return None
#     points = [shape(wp.geometry.dict()) for wp in waypoints]
#     multipoint = MultiPoint(points)
#     return from_shape(multipoint, srid=4326)

# def convert_constraints_to_multipolygon(constraints) -> Optional[str]:
#     """Convert list of GeoJSON Polygon features into a PostGIS MULTIPOLYGON."""
#     if not constraints:
#         return None
#     polygons = [shape(c.geometry.dict()) for c in constraints]
#     multipolygon = MultiPolygon(polygons)
#     return from_shape(multipolygon, srid=4326)

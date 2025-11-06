from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from typing import List

from app.db.session import get_db
from app.models.routing_request import RoutingRequest
from app.schemas.routing_history import RoutingRequestOut
from app.auth.dependencies import get_current_user  # adjust import if your auth path differs

router = APIRouter()

@router.get("/history", response_model=List[RoutingRequestOut])
async def get_routing_history(
    db: AsyncSession = Depends(get_db),
    current_user=Depends(get_current_user)
):
    stmt = (
        select(RoutingRequest)
        .where(RoutingRequest.user_id == current_user.id)
        .order_by(RoutingRequest.timestamp.desc())
    )
    result = await db.execute(stmt)
    return result.scalars().all()

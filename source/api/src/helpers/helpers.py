from typing import List

from schemas.schemas import LatLonAlt

# ----------------- Helpers -----------------

def normalize_alt(value: float, unit: str) -> float:
    return value if unit == "m" else value * 0.3048


def haversine_km(lat1, lon1, lat2, lon2) -> float:
    from math import radians, sin, cos, sqrt, atan2

    R = 6371.0088
    dlat = radians(lat2 - lat1)
    dlon = radians(lon2 - lon1)
    a = sin(dlat / 2) ** 2 + cos(radians(lat1)) * cos(radians(lat2)) * sin(dlon / 2) ** 2
    return 2 * R * atan2(sqrt(a), sqrt(1 - a))


def path_length_km(points: List[LatLonAlt]) -> float:
    total = 0.0
    for i in range(len(points) - 1):
        a, b = points[i], points[i + 1]
        total += haversine_km(a.lat, a.lon, b.lat, b.lon)
    return round(total, 3)


# async def get_current_user(db: AsyncSession = Depends(get_db), token: str = Depends(oauth2_scheme)) -> User:
#     sub = decode_token(token)
#     if not sub:
#         raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid token")
#     uid = uuid.UUID(sub)
#     row = await db.execute(select(User).where(User.id == uid))
#     user = row.scalar_one_or_none()
#     if not user:
#         raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="User not found")
#     return user

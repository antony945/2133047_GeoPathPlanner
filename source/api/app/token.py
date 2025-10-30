from fastapi import Query, Header, HTTPException
from jwt import exceptions as jwt_exceptions
import jwt
import os
from app.config import JWT_SECRET_KEY, JWT_ALGORITHM
from typing import Optional

# TODO: Just for now, maybe it's better to always have a token, also for unregistered user, boh
async def verify_jwt_token(
    authorization: Optional[str] = Header(None),
    user_id: str | None = Query(None),
) -> dict | None:
    """
    Verifies JWT token from the Authorization header for a given user_id.

    - If `user_id` is not provided, returns None (no verification needed).
    - If `user_id` is provided but the Authorization header is missing or invalid, raises HTTPException.
    - Returns the decoded JWT payload as a dictionary if valid.
    """
    # TODO: Just for testing, now return always None
    return None

    if not user_id:
        return None

    if not authorization:
        raise HTTPException(status_code=401, detail="Missing Authorization header")
    
    try:
        # Remove "Bearer " prefix if present
        token = authorization.replace("Bearer ", "")
        payload = jwt.decode(token, JWT_SECRET_KEY, algorithms=[JWT_ALGORITHM])
        return payload
    except jwt_exceptions.ExpiredSignatureError as e:
        raise HTTPException(status_code=401, detail=f"Expired JWT token: {e}")
    except jwt_exceptions.InvalidTokenError as e:
        raise HTTPException(status_code=401, detail=f"Invalid JWT token: {e}")
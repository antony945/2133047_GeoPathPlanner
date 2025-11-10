from fastapi import Query, Header, HTTPException
from jwt import exceptions as jwt_exceptions
import jwt
import os
from app.config import JWT_SECRET_KEY, JWT_ALGORITHM
from typing import Optional
from app.logger import logger

# TODO: Just for now, maybe it's better to always have a token, also for unregistered user, boh
async def verify_jwt_token(    
    authorization: Optional[str] = Header(None),
    user_id: str | None = Query(None),
) -> dict:
    """
    Verify a JWT token from the Authorization header.

    Raises an HTTPException if:
        - The Authorization header is missing (401).
        - The JWT token is expired (401).
        - The JWT token is invalid (401).
        - The user_id is provided and does not match the token's 'sub' claim (403).

    If no exception is raised, the token is valid, and execution can continue.

    Args:
        authorization (Optional[str]): JWT token from the "Authorization" header.
        user_id (Optional[str]): User ID to verify against the token's "sub" claim.

    Returns:
        dict: The decoded JWT payload.
    """
        
    # TODO: For now return here so we can test without the actual jwt token
    return None

    if not authorization:
        raise HTTPException(status_code=401, detail="Missing Authorization header")
    
    try:
        # Remove "Bearer " prefix if present
        token = authorization.replace("Bearer ", "")
        payload = jwt.decode(token, JWT_SECRET_KEY, algorithms=[JWT_ALGORITHM])
        
        # TODO: check if decode could raise other expeections 

        if user_id and payload:
            jwt_sub = payload.get("sub")
            if user_id and str(jwt_sub) != str(user_id):
                logger.warning(f"🔒 User ID mismatch: token.sub={jwt_sub}, query.user_id={user_id}")
                raise HTTPException(status_code=403, detail="User ID mismatch")
            else:
                logger.debug(f"✅ JWT validated for user_id={jwt_sub or user_id}")
        else:
            logger.info("⚠️ No User ID provided — treating as anonymous request.")

    except jwt_exceptions.ExpiredSignatureError as e:
        raise HTTPException(status_code=401, detail=f"Expired JWT token: {e}")
    except jwt_exceptions.InvalidTokenError as e:
        raise HTTPException(status_code=401, detail=f"Invalid JWT token: {e}")
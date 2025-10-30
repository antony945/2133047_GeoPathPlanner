from fastapi import Header, HTTPException, Depends, Query
from jwt import exceptions as jwt_exceptions
import jwt
import os
from app.config import JWT_SECRET_KEY, JWT_ALGORITHM

JWT_SECRET = os.getenv("JWT_SECRET", "default_secret")
JWT_ALGORITHM = os.getenv("JWT_ALGORITHM", "HS256")

def verify_jwt_token(auth_header: str = Header(...)):
    """
    Extract and verify JWT from Authorization header.
    """
    # TODO: to check this
    return
    if not auth_header.startswith("Bearer "):
        raise HTTPException(status_code=401, detail="Invalid Authorization header")

    token = auth_header.split(" ")[1]

    try:
        payload = jwt.decode(token, JWT_ALGORITHM, algorithms=[JWT_SECRET_KEY])
    except jwt_exceptions.ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="Token expired")
    except jwt_exceptions.InvalidTokenError:
        raise HTTPException(status_code=401, detail="Invalid token")
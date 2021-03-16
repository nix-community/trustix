from trustix_nix_reprod.cache import connection
from pydantic import BaseModel, ValidationError
from functools import wraps
from typing import (
    Optional,
    Type,
)
import hashlib
import asyncio
import orjson


def _cache_key(fn, args, kwargs) -> str:
    return hashlib.sha256(orjson.dumps([
        fn.__module__,
        fn.__name__,
        args,
        kwargs,
    ])).hexdigest()


def cached(model: Type[BaseModel], ttl: int):
    """Cache decorator for pydantic models"""

    def wrapper(fn):
        @wraps(fn)
        async def inner(*args, **kwargs):
            conn = await connection.get()
            key = _cache_key(fn, args, kwargs)

            cached_result: Optional[str] = await conn.get(key)

            if cached_result:
                try:
                    resp = model.parse_raw(cached_result)
                except ValidationError as e:
                    print("Incompatible schema detected:")
                    print(e)
                else:
                    asyncio.create_task(conn.expire(key, ttl))
                    return resp

            result: BaseModel = await fn(*args, **kwargs)
            asyncio.create_task(conn.set(key, result.json(), expire=ttl))

            return result

        return inner
    return wrapper

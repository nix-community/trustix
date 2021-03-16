from trustix_nix_reprod.conf import settings
from typing import Optional
import aioredis  # type: ignore
import asyncio


__all__ = (
    "close",
    "get",
)


_lock = asyncio.Lock()
_connection: Optional[aioredis.ConnectionsPool] = None


async def get():
    global _connection
    if _connection:
        return _connection

    await _lock.acquire()
    try:
        _connection = await aioredis.create_redis_pool(settings.redis_uri)
        return _connection
    finally:
        _lock.release()


async def close():
    global _connection
    if _connection:
        await _connection.close()

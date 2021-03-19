# Trustix
# Copyright (C) 2021 Tweag IO

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

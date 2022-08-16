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

from trustix_nix_reprod.cache import connection
from pydantic import BaseModel, ValidationError
from functools import wraps
from typing import (
    Optional,
    Type,
)
import hashlib
import asyncio

import json


def _cache_key(fn, args, kwargs) -> str:
    return hashlib.sha256(
        json.dumps(
            [
                fn.__module__,
                fn.__name__,
                args,
                kwargs,
            ]
        ).encode()
    ).hexdigest()


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

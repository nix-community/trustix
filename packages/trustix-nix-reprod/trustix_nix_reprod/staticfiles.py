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

from starlette.staticfiles import StaticFiles as _StaticFiles  # type: ignore
from starlette.staticfiles import NotModifiedResponse  # type: ignore
from starlette.datastructures import Headers  # type: ignore
from starlette.responses import (  # type: ignore
    FileResponse,
    Response,
)
from starlette.types import Scope  # type: ignore
from async_lru import alru_cache  # type: ignore
import aiofiles
import hashlib
import os


# Note: The entire app is reloaded on JS changes in development
# and in production they are immutable so we can enable the lru
# cache regardless of the environment
@alru_cache(maxsize=None)
async def _hash_file(path: str) -> str:
    h = hashlib.sha256()
    async with aiofiles.open(path, mode="rb") as file:
        while True:
            chunk = await file.read(FileResponse.chunk_size)
            h.update(chunk)
            if len(chunk) < FileResponse.chunk_size:
                break
    return h.hexdigest()


class StaticFiles(_StaticFiles):
    """
    A version of StaticFiles that sets etag based on file contents
    """

    def file_response(
        self,
        full_path: str,
        stat_result: os.stat_result,
        scope: Scope,
        status_code: int = 200,
    ) -> Response:
        method = scope["method"]
        return FileResponse(
            full_path, status_code=status_code, stat_result=stat_result, method=method
        )

    async def get_response(self, path: str, scope: Scope) -> Response:
        response = await super().get_response(path, scope)
        if isinstance(response, FileResponse):
            etag = await _hash_file(response.path)
            response.headers["etag"] = etag
            if self.is_not_modified(response.headers, Headers(scope=scope)):
                return NotModifiedResponse(response.headers)
        return response

#!/usr/bin/env python
# The MIT License
# Copyright (c) 2018 Adam Hose (adisbladis)

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

import contextlib
import functools
import asyncio


class DeferStack:
    """Use go-style defer from python"""

    def __init__(self):
        self._stack = contextlib.ExitStack()
        self._astack = contextlib.AsyncExitStack()

    def __enter__(self):
        def defer(fn, *args, **kwargs):
            partial = functools.partial(fn, *args, **kwargs)
            self._stack.callback(partial)

        return defer

    def __exit__(self, *exc_details):
        self._stack.__exit__(*exc_details)

    async def __aenter__(self):
        def defer(fn, *args, **kwargs):
            partial = functools.partial(fn, *args, **kwargs)
            if asyncio.iscoroutinefunction(fn):
                self._astack.push_async_callback(partial)
            else:
                self._astack.callback(partial)

        return defer

    async def __aexit__(self, *exc_details):
        await self._astack.__aexit__(*exc_details)


# # Usage examples

# # Async
# async def amain():
#     async def adone():
#         print('adone')
#     def done():
#         print('done')
#     async with DeferStack() as defer:
#         defer(adone)
#         defer(done)

# import asyncio
# loop = asyncio.get_event_loop()
# loop.run_until_complete(amain())

# # Sync
# def main():
#     with DeferStack() as defer:
#         defer(print, 'done')

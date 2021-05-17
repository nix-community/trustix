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

from functools import lru_cache
import grpc  # type: ignore

from trustix_python.rpc import rpc_pb2_grpc  # type: ignore
from trustix_nix_reprod.conf import settings


@lru_cache(maxsize=None)
def get_channel() -> grpc.aio.Channel:
    return grpc.aio.insecure_channel(settings.trustix_rpc)


@lru_cache(maxsize=None)
def get_rpcapi() -> rpc_pb2_grpc.RPCApiStub:  # type: ignore
    return rpc_pb2_grpc.RPCApiStub(get_channel())  # type: ignore


@lru_cache(maxsize=None)
def get_logrpc() -> rpc_pb2_grpc.LogRPCStub:  # type: ignore
    return rpc_pb2_grpc.LogRPCStub(get_channel())  # type: ignore

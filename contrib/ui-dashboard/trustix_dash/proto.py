from functools import lru_cache
import grpc  # type: ignore

from trustix_proto import trustix_pb2_grpc  # type: ignore
from trustix_dash.conf import settings


@lru_cache(maxsize=None)
def get_channel() -> grpc.aio.Channel:
    return grpc.aio.insecure_channel(settings.trustix_rpc)


@lru_cache(maxsize=None)
def get_combined_rpc() -> trustix_pb2_grpc.TrustixCombinedRPCStub:  # type: ignore
    return trustix_pb2_grpc.TrustixCombinedRPCStub(get_channel())  # type: ignore

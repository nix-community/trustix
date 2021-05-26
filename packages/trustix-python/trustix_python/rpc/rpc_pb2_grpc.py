# type: ignore
# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from trustix_python.api import api_pb2 as api_dot_api__pb2
from trustix_python.rpc import rpc_pb2 as rpc_dot_rpc__pb2
from trustix_python.schema import loghead_pb2 as schema_dot_loghead__pb2


class RPCApiStub(object):
    """RPCApi are "private" rpc methods for an instance.
    This should only be available to trusted parties.
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Logs = channel.unary_unary(
            "/trustix.RPCApi/Logs",
            request_serializer=api_dot_api__pb2.LogsRequest.SerializeToString,
            response_deserializer=api_dot_api__pb2.LogsResponse.FromString,
        )
        self.Decide = channel.unary_unary(
            "/trustix.RPCApi/Decide",
            request_serializer=rpc_dot_rpc__pb2.DecideRequest.SerializeToString,
            response_deserializer=rpc_dot_rpc__pb2.DecisionResponse.FromString,
        )
        self.GetValue = channel.unary_unary(
            "/trustix.RPCApi/GetValue",
            request_serializer=api_dot_api__pb2.ValueRequest.SerializeToString,
            response_deserializer=api_dot_api__pb2.ValueResponse.FromString,
        )


class RPCApiServicer(object):
    """RPCApi are "private" rpc methods for an instance.
    This should only be available to trusted parties.
    """

    def Logs(self, request, context):
        """Get a list of all logs published/subscribed by this node"""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def Decide(self, request, context):
        """Decide on an output for key based on the configured decision method"""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def GetValue(self, request, context):
        """Get values by their content-address"""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")


def add_RPCApiServicer_to_server(servicer, server):
    rpc_method_handlers = {
        "Logs": grpc.unary_unary_rpc_method_handler(
            servicer.Logs,
            request_deserializer=api_dot_api__pb2.LogsRequest.FromString,
            response_serializer=api_dot_api__pb2.LogsResponse.SerializeToString,
        ),
        "Decide": grpc.unary_unary_rpc_method_handler(
            servicer.Decide,
            request_deserializer=rpc_dot_rpc__pb2.DecideRequest.FromString,
            response_serializer=rpc_dot_rpc__pb2.DecisionResponse.SerializeToString,
        ),
        "GetValue": grpc.unary_unary_rpc_method_handler(
            servicer.GetValue,
            request_deserializer=api_dot_api__pb2.ValueRequest.FromString,
            response_serializer=api_dot_api__pb2.ValueResponse.SerializeToString,
        ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
        "trustix.RPCApi", rpc_method_handlers
    )
    server.add_generic_rpc_handlers((generic_handler,))


# This class is part of an EXPERIMENTAL API.
class RPCApi(object):
    """RPCApi are "private" rpc methods for an instance.
    This should only be available to trusted parties.
    """

    @staticmethod
    def Logs(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/trustix.RPCApi/Logs",
            api_dot_api__pb2.LogsRequest.SerializeToString,
            api_dot_api__pb2.LogsResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
        )

    @staticmethod
    def Decide(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/trustix.RPCApi/Decide",
            rpc_dot_rpc__pb2.DecideRequest.SerializeToString,
            rpc_dot_rpc__pb2.DecisionResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
        )

    @staticmethod
    def GetValue(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/trustix.RPCApi/GetValue",
            api_dot_api__pb2.ValueRequest.SerializeToString,
            api_dot_api__pb2.ValueResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
        )


class LogRPCStub(object):
    """RPCApi are "private" rpc methods for an instance related to a specific log.
    This should only be available to trusted parties.
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetHead = channel.unary_unary(
            "/trustix.LogRPC/GetHead",
            request_serializer=api_dot_api__pb2.LogHeadRequest.SerializeToString,
            response_deserializer=schema_dot_loghead__pb2.LogHead.FromString,
        )
        self.GetLogEntries = channel.unary_unary(
            "/trustix.LogRPC/GetLogEntries",
            request_serializer=api_dot_api__pb2.GetLogEntriesRequest.SerializeToString,
            response_deserializer=api_dot_api__pb2.LogEntriesResponse.FromString,
        )
        self.Submit = channel.unary_unary(
            "/trustix.LogRPC/Submit",
            request_serializer=rpc_dot_rpc__pb2.SubmitRequest.SerializeToString,
            response_deserializer=rpc_dot_rpc__pb2.SubmitResponse.FromString,
        )
        self.Flush = channel.unary_unary(
            "/trustix.LogRPC/Flush",
            request_serializer=rpc_dot_rpc__pb2.FlushRequest.SerializeToString,
            response_deserializer=rpc_dot_rpc__pb2.FlushResponse.FromString,
        )


class LogRPCServicer(object):
    """RPCApi are "private" rpc methods for an instance related to a specific log.
    This should only be available to trusted parties.
    """

    def GetHead(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def GetLogEntries(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def Submit(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def Flush(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")


def add_LogRPCServicer_to_server(servicer, server):
    rpc_method_handlers = {
        "GetHead": grpc.unary_unary_rpc_method_handler(
            servicer.GetHead,
            request_deserializer=api_dot_api__pb2.LogHeadRequest.FromString,
            response_serializer=schema_dot_loghead__pb2.LogHead.SerializeToString,
        ),
        "GetLogEntries": grpc.unary_unary_rpc_method_handler(
            servicer.GetLogEntries,
            request_deserializer=api_dot_api__pb2.GetLogEntriesRequest.FromString,
            response_serializer=api_dot_api__pb2.LogEntriesResponse.SerializeToString,
        ),
        "Submit": grpc.unary_unary_rpc_method_handler(
            servicer.Submit,
            request_deserializer=rpc_dot_rpc__pb2.SubmitRequest.FromString,
            response_serializer=rpc_dot_rpc__pb2.SubmitResponse.SerializeToString,
        ),
        "Flush": grpc.unary_unary_rpc_method_handler(
            servicer.Flush,
            request_deserializer=rpc_dot_rpc__pb2.FlushRequest.FromString,
            response_serializer=rpc_dot_rpc__pb2.FlushResponse.SerializeToString,
        ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
        "trustix.LogRPC", rpc_method_handlers
    )
    server.add_generic_rpc_handlers((generic_handler,))


# This class is part of an EXPERIMENTAL API.
class LogRPC(object):
    """RPCApi are "private" rpc methods for an instance related to a specific log.
    This should only be available to trusted parties.
    """

    @staticmethod
    def GetHead(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/trustix.LogRPC/GetHead",
            api_dot_api__pb2.LogHeadRequest.SerializeToString,
            schema_dot_loghead__pb2.LogHead.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
        )

    @staticmethod
    def GetLogEntries(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/trustix.LogRPC/GetLogEntries",
            api_dot_api__pb2.GetLogEntriesRequest.SerializeToString,
            api_dot_api__pb2.LogEntriesResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
        )

    @staticmethod
    def Submit(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/trustix.LogRPC/Submit",
            rpc_dot_rpc__pb2.SubmitRequest.SerializeToString,
            rpc_dot_rpc__pb2.SubmitResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
        )

    @staticmethod
    def Flush(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/trustix.LogRPC/Flush",
            rpc_dot_rpc__pb2.FlushRequest.SerializeToString,
            rpc_dot_rpc__pb2.FlushResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
        )
# python3 -m grpc_tools.protoc --proto_path=. ./proto/transistor/v1/transistor.proto --python_out=. --grpc_python_out=.

import grpc
from gen.proto.transistor.v1 import transistor_pb2
from gen.proto.transistor.v1 import transistor_pb2_grpc

class TransistorClient:
    def __init__(self, addr="localhost:50051"):
        self.addr = addr
        self.channel = grpc.insecure_channel(self.addr)
        self.stub = self._connect()

    def _connect(self):
        return transistor_pb2_grpc.TransistorServiceStub(self.channel)

    def close(self):
        self.channel.close()
    
    def publish(self):
        self.stub.Publish()

# client = TransistorClient(addr="localhost:50051")
# client.close()
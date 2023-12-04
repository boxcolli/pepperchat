import grpc
import transistor_pb2
import transistor_pb2_grpc

class TransistorClient:
    def __init__(self, addr="localhost:50051"):
        self.addr = addr
        self.channel = grpc.insecure_channel(self.addr)
        self.stub = self._connect()

    def _connect(self):
        return transistor_pb2_grpc.TransistorServiceStub(self.channel)

    def close(self):
        self.channel.close()

# client = TransistorClient(addr="localhost:50051")
# client.close()
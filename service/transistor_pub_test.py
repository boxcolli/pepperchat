import asyncio
import logging

import grpc
from gen.proto.transistor.v1 import transistor_pb2
from gen.proto.transistor.v1 import transistor_pb2_grpc

async def run() -> None:
    async with grpc.aio.insecure_channel("localhost:50050") as channel:
        stub = transistor_pb2_grpc.TransistorServiceStub(channel)

        # Command Ping
        async for response in stub.Command(
            transistor_pb2.CommandRequest(args=["ping"])
        ):
            print(
                "Transistor client received from async generator: "
                + response.message
            )

async def pub(stub: transistor_pb2_grpc.TransistorServiceStub) -> None:
    stub.Publish()

if __name__ == "__main__":
    logging.basicConfig()
    asyncio.run(run())
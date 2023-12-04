from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf import any_pb2 as _any_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Method(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    METHOD_UNSPECIFIED: _ClassVar[Method]
    METHOD_EMPTY: _ClassVar[Method]
    METHOD_CREATE: _ClassVar[Method]
    METHOD_UPDATE: _ClassVar[Method]
    METHOD_DELETE: _ClassVar[Method]

class Mode(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    MODE_UNSPECIFIED: _ClassVar[Mode]
    MODE_ANY: _ClassVar[Mode]
    MODE_ROUTE: _ClassVar[Mode]
    MODE_ROOT: _ClassVar[Mode]

class Operation(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    OPERATION_UNSPECIFIED: _ClassVar[Operation]
    OPERATION_ADD: _ClassVar[Operation]
    OPERATION_DEL: _ClassVar[Operation]
METHOD_UNSPECIFIED: Method
METHOD_EMPTY: Method
METHOD_CREATE: Method
METHOD_UPDATE: Method
METHOD_DELETE: Method
MODE_UNSPECIFIED: Mode
MODE_ANY: Mode
MODE_ROUTE: Mode
MODE_ROOT: Mode
OPERATION_UNSPECIFIED: Operation
OPERATION_ADD: Operation
OPERATION_DEL: Operation

class Topic(_message.Message):
    __slots__ = ("tokens",)
    TOKENS_FIELD_NUMBER: _ClassVar[int]
    tokens: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, tokens: _Optional[_Iterable[str]] = ...) -> None: ...

class Message(_message.Message):
    __slots__ = ("mode", "topic", "method", "data", "timestamp")
    MODE_FIELD_NUMBER: _ClassVar[int]
    TOPIC_FIELD_NUMBER: _ClassVar[int]
    METHOD_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    mode: Mode
    topic: Topic
    method: Method
    data: _any_pb2.Any
    timestamp: _timestamp_pb2.Timestamp
    def __init__(self, mode: _Optional[_Union[Mode, str]] = ..., topic: _Optional[_Union[Topic, _Mapping]] = ..., method: _Optional[_Union[Method, str]] = ..., data: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., timestamp: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class Change(_message.Message):
    __slots__ = ("mode", "op", "topic")
    MODE_FIELD_NUMBER: _ClassVar[int]
    OP_FIELD_NUMBER: _ClassVar[int]
    TOPIC_FIELD_NUMBER: _ClassVar[int]
    mode: Mode
    op: Operation
    topic: Topic
    def __init__(self, mode: _Optional[_Union[Mode, str]] = ..., op: _Optional[_Union[Operation, str]] = ..., topic: _Optional[_Union[Topic, _Mapping]] = ...) -> None: ...

class PublishRequest(_message.Message):
    __slots__ = ("msg",)
    MSG_FIELD_NUMBER: _ClassVar[int]
    msg: Message
    def __init__(self, msg: _Optional[_Union[Message, _Mapping]] = ...) -> None: ...

class PublishResponse(_message.Message):
    __slots__ = ("change",)
    CHANGE_FIELD_NUMBER: _ClassVar[int]
    change: _containers.RepeatedCompositeFieldContainer[Change]
    def __init__(self, change: _Optional[_Iterable[_Union[Change, _Mapping]]] = ...) -> None: ...

class SubscribeRequest(_message.Message):
    __slots__ = ("change",)
    CHANGE_FIELD_NUMBER: _ClassVar[int]
    change: Change
    def __init__(self, change: _Optional[_Union[Change, _Mapping]] = ...) -> None: ...

class SubscribeResponse(_message.Message):
    __slots__ = ("msg",)
    MSG_FIELD_NUMBER: _ClassVar[int]
    msg: Message
    def __init__(self, msg: _Optional[_Union[Message, _Mapping]] = ...) -> None: ...

class CommandRequest(_message.Message):
    __slots__ = ("args",)
    ARGS_FIELD_NUMBER: _ClassVar[int]
    args: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, args: _Optional[_Iterable[str]] = ...) -> None: ...

class CommandResponse(_message.Message):
    __slots__ = ("line",)
    LINE_FIELD_NUMBER: _ClassVar[int]
    line: str
    def __init__(self, line: _Optional[str] = ...) -> None: ...

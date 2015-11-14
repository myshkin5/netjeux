package factory

import (
	"time"

	"github.com/myshkin5/netspel/jsonstruct"
)

type Scheme interface {
	Init(config jsonstruct.JSONStruct) error

	BytesPerMessage() int
	MessagesPerRun() int

	RunWriter(writer Writer)
	RunReader(reader Reader)

	ByteCount() uint64
	ErrorCount() uint32
	FirstError() error
	RunTime() time.Duration
}

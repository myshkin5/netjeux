package factory

import (
	"time"
)

type Scheme interface {
	Init(config map[string]interface{}) error

	BytesPerMessage() int
	MessagesPerRun() int

	RunWriter(writer Writer)
	RunReader(reader Reader)

	ByteCount() uint64
	ErrorCount() uint32
	FirstError() error
	RunTime() time.Duration
}

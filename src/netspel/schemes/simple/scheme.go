package simple

import (
	"sync/atomic"

	"netspel/factory"
)

const (
	MessagesPerRun  = 10000
	BytesPerMessage = 1000
)

type Scheme struct {
	buffer     []byte
	byteCount  uint64
	errorCount uint32

	DefaultReport  string
	LessThanReport string
	ErrorReport    string
}

func (s *Scheme) Init(config factory.Config) error {
	s.buffer = make([]byte, BytesPerMessage)

	s.DefaultReport = "."
	s.LessThanReport = "<"
	s.ErrorReport = "#"

	return nil
}

func (s *Scheme) ByteCount() uint64 {
	return atomic.LoadUint64(&s.byteCount)
}

func (s *Scheme) ErrorCount() uint32 {
	return atomic.LoadUint32(&s.errorCount)
}

func (s *Scheme) RunWriter(writer factory.Writer) {
	for i := 0; i < MessagesPerRun; i++ {
		s.countMessage(writer.Write(s.buffer))
	}
}

func (s *Scheme) RunReader(reader factory.Reader) {
	buffer := make([]byte, BytesPerMessage)
	for {
		s.countMessage(reader.Read(buffer))
	}
}

func (s *Scheme) countMessage(bytes int, err error) {
	atomic.AddUint64(&s.byteCount, uint64(bytes))
	if err != nil {
		atomic.AddUint32(&s.errorCount, 1)
	}
	switch {
	case err != nil:
		print(s.ErrorReport)
	case bytes < BytesPerMessage:
		print(s.LessThanReport)
	default:
		print(s.DefaultReport)
	}
}

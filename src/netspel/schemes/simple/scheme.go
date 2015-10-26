package simple

import (
	"io"
	"sync/atomic"
)

const (
	MessagesPerRun  = 10000
	BytesPerMessage = 32768
)

type Scheme struct {
	writer     io.Writer
	reader     io.Reader
	buffer     []byte
	byteCount  uint64
	errorCount uint32

	DefaultReport  string
	LessThanReport string
	ErrorReport    string
}

func New(writer io.Writer, reader io.Reader) *Scheme {
	return &Scheme{
		writer: writer,
		reader: reader,
		buffer: make([]byte, BytesPerMessage),

		DefaultReport:  ".",
		LessThanReport: "<",
		ErrorReport:    "#",
	}
}

func (s *Scheme) ByteCount() uint64 {
	return atomic.LoadUint64(&s.byteCount)
}

func (s *Scheme) ErrorCount() uint32 {
	return atomic.LoadUint32(&s.errorCount)
}

func (s *Scheme) RunWriter() {
	for i := 0; i < MessagesPerRun; i++ {
		s.countMessage(s.writer.Write(s.buffer))
	}
}

func (s *Scheme) RunReader() {
	buffer := make([]byte, BytesPerMessage)
	for {
		s.countMessage(s.reader.Read(buffer))
	}
}

func (s *Scheme) countMessage(bytes int, err error) {
	atomic.AddUint64(&s.byteCount, uint64(bytes))
	if err != nil {
		atomic.AddUint32(&s.errorCount, 1)
	}
	switch {
	case bytes < BytesPerMessage:
		print(s.LessThanReport)
	case err != nil:
		print(s.ErrorReport)
	default:
		print(s.DefaultReport)
	}
}

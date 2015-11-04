package simple

import (
	"fmt"
	"sync/atomic"

	"netspel/factory"
)

const (
	prefix = "simple."

	MessagesPerRun  = prefix + "messages-per-run"
	BytesPerMessage = prefix + "bytes-per-message"

	reportPrefix = prefix + "report-flags."

	DefaultReport  = reportPrefix + "default"
	LessThanReport = reportPrefix + "less-than"
	ErrorReport    = reportPrefix + "error"
)

type Scheme struct {
	buffer     []byte
	byteCount  uint64
	errorCount uint32

	bytesPerMessage int
	messagesPerRun  int

	defaultReport  string
	lessThanReport string
	errorReport    string
}

func (s *Scheme) Init(config factory.Config) error {
	var ok bool
	s.bytesPerMessage, ok = config.AdditionalInt(BytesPerMessage)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", BytesPerMessage)
	}
	s.buffer = make([]byte, s.bytesPerMessage)

	s.messagesPerRun, ok = config.AdditionalInt(MessagesPerRun)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", MessagesPerRun)
	}

	s.defaultReport, ok = config.AdditionalString(DefaultReport)
	if !ok {
		s.defaultReport = "."
	}
	s.lessThanReport, ok = config.AdditionalString(LessThanReport)
	if !ok {
		s.lessThanReport = "<"
	}
	s.errorReport, ok = config.AdditionalString(ErrorReport)
	if !ok {
		s.errorReport = "#"
	}

	return nil
}

func (s *Scheme) ByteCount() uint64 {
	return atomic.LoadUint64(&s.byteCount)
}

func (s *Scheme) ErrorCount() uint32 {
	return atomic.LoadUint32(&s.errorCount)
}

func (s *Scheme) RunWriter(writer factory.Writer) {
	for i := 0; i < s.messagesPerRun; i++ {
		s.countMessage(writer.Write(s.buffer))
	}
}

func (s *Scheme) RunReader(reader factory.Reader) {
	buffer := make([]byte, s.bytesPerMessage)
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
		print(s.errorReport)
	case bytes < s.bytesPerMessage:
		print(s.lessThanReport)
	default:
		print(s.defaultReport)
	}
}

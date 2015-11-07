package simple

import (
	"fmt"
	"sync/atomic"

	"netspel/factory"
	"sync"
	"time"
)

const (
	prefix = "simple."

	MessagesPerRun     = prefix + "messages-per-run"
	BytesPerMessage    = prefix + "bytes-per-message"
	WaitForLastMessage = prefix + "wait-for-last-message"

	reportPrefix = prefix + "report-flags."

	DefaultReport     = reportPrefix + "default"
	LessThanReport    = reportPrefix + "less-than"
	ErrorReport       = reportPrefix + "error"
	GreaterThanReport = reportPrefix + "greater-than"
)

type Scheme struct {
	buffer     []byte
	byteCount  uint64
	errorCount uint32
	runtime    time.Duration

	bytesPerMessage    int
	messagesPerRun     int
	waitForLastMessage time.Duration

	defaultReport     string
	lessThanReport    string
	errorReport       string
	greaterThanReport string
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

	wait, ok := config.AdditionalString(WaitForLastMessage)
	if !ok {
		wait = "5s"
	}
	var err error
	s.waitForLastMessage, err = time.ParseDuration(wait)
	if err != nil {
		return err
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
	s.greaterThanReport, ok = config.AdditionalString(GreaterThanReport)
	if !ok {
		s.greaterThanReport = ">"
	}

	return nil
}

func (s *Scheme) ByteCount() uint64 {
	return atomic.LoadUint64(&s.byteCount)
}

func (s *Scheme) ErrorCount() uint32 {
	return atomic.LoadUint32(&s.errorCount)
}

func (s *Scheme) Runtime() time.Duration {
	return s.runtime
}

func (s *Scheme) RunWriter(writer factory.Writer) {
	startTime := time.Now()
	for i := 0; i < s.messagesPerRun; i++ {
		s.countMessage(writer.Write(s.buffer))
	}
	s.runtime = time.Now().Sub(startTime)
}

func (s *Scheme) RunReader(reader factory.Reader) {
	timer := time.NewTimer(time.Duration(1<<63 - 1))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.runReader(reader, timer)
	}()

	<-timer.C

	reader.Stop()

	wg.Wait()
}

func (s *Scheme) runReader(reader factory.Reader, timer *time.Timer) {
	var startTime time.Time
	buffer := make([]byte, s.bytesPerMessage*2)
	for {
		count, err := reader.Read(buffer)

		if startTime.IsZero() {
			startTime = time.Now()
		}

		timer.Reset(s.waitForLastMessage)

		if err == factory.ErrReaderClosed {
			break
		}

		s.countMessage(count, err)
	}

	s.runtime = time.Now().Sub(startTime)
}

func (s *Scheme) countMessage(count int, err error) {
	atomic.AddUint64(&s.byteCount, uint64(count))
	if err != nil {
		atomic.AddUint32(&s.errorCount, 1)
	}
	switch {
	case err != nil:
		print(s.errorReport)
	case count < s.bytesPerMessage:
		print(s.lessThanReport)
	case count > s.bytesPerMessage:
		print(s.greaterThanReport)
	default:
		print(s.defaultReport)
	}
}

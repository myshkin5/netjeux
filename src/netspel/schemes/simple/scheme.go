package simple

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"netspel/factory"
	"netspel/json"
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
	firstError error
	runTime    time.Duration

	bytesPerMessage    int
	messagesPerRun     int
	waitForLastMessage time.Duration

	defaultReport     string
	lessThanReport    string
	errorReport       string
	greaterThanReport string
}

func (s *Scheme) Init(config map[string]interface{}) error {
	var ok bool
	s.bytesPerMessage, ok = json.Int(BytesPerMessage, config)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", BytesPerMessage)
	}
	s.buffer = make([]byte, s.bytesPerMessage)

	s.messagesPerRun, ok = json.Int(MessagesPerRun, config)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", MessagesPerRun)
	}

	wait, ok := json.String(WaitForLastMessage, config)
	if !ok {
		wait = "5s"
	}
	var err error
	s.waitForLastMessage, err = time.ParseDuration(wait)
	if err != nil {
		return err
	}

	s.defaultReport, ok = json.String(DefaultReport, config)
	if !ok {
		s.defaultReport = "."
	}
	s.lessThanReport, ok = json.String(LessThanReport, config)
	if !ok {
		s.lessThanReport = "<"
	}
	s.errorReport, ok = json.String(ErrorReport, config)
	if !ok {
		s.errorReport = "#"
	}
	s.greaterThanReport, ok = json.String(GreaterThanReport, config)
	if !ok {
		s.greaterThanReport = ">"
	}

	return nil
}

func (s *Scheme) BytesPerMessage() int {
	return s.bytesPerMessage
}

func (s *Scheme) MessagesPerRun() int {
	return s.messagesPerRun
}

func (s *Scheme) ByteCount() uint64 {
	return atomic.LoadUint64(&s.byteCount)
}

func (s *Scheme) ErrorCount() uint32 {
	return atomic.LoadUint32(&s.errorCount)
}

func (s *Scheme) FirstError() error {
	return s.firstError
}

func (s *Scheme) RunTime() time.Duration {
	return s.runTime
}

func (s *Scheme) RunWriter(writer factory.Writer) {
	startTime := time.Now()
	for i := 0; i < s.messagesPerRun; i++ {
		s.countMessage(writer.Write(s.buffer))
	}
	s.runTime = time.Now().Sub(startTime)
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

	s.runTime = time.Now().Sub(startTime)
}

func (s *Scheme) countMessage(count int, err error) {
	atomic.AddUint64(&s.byteCount, uint64(count))
	if err != nil {
		atomic.AddUint32(&s.errorCount, 1)
		if s.firstError == nil {
			s.firstError = err
		}
	}
	switch {
	case err != nil:
		fmt.Print(s.errorReport)
	case count < s.bytesPerMessage:
		fmt.Print(s.lessThanReport)
	case count > s.bytesPerMessage:
		fmt.Print(s.greaterThanReport)
	default:
		fmt.Print(s.defaultReport)
	}
}

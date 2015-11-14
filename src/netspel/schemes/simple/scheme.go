package simple

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"netspel/factory"
	"netspel/jsonstruct"
	"netspel/logs"
)

const (
	prefix = "simple."

	MessagesPerRun     = prefix + "messages-per-run"
	BytesPerMessage    = prefix + "bytes-per-message"
	WaitForLastMessage = prefix + "wait-for-last-message"

	WarmupMessagesPerRun = prefix + "warmup-messages-per-run"
	WarmupWait           = prefix + "warmup-wait"

	reportPrefix      = prefix + "report-flags."
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

	warmupMessagesPerRun int
	warmupWait           time.Duration

	defaultReport     string
	lessThanReport    string
	errorReport       string
	greaterThanReport string
}

func (s *Scheme) Init(config jsonstruct.JSONStruct) error {
	var ok bool
	s.bytesPerMessage, ok = config.Int(BytesPerMessage)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", BytesPerMessage)
	}
	s.buffer = make([]byte, s.bytesPerMessage)

	s.messagesPerRun, ok = config.Int(MessagesPerRun)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", MessagesPerRun)
	}

	var err error
	s.waitForLastMessage, err = config.DurationWithDefault(WaitForLastMessage, 5*time.Second)
	if err != nil {
		return err
	}

	s.warmupMessagesPerRun = config.IntWithDefault(WarmupMessagesPerRun, 0)
	s.warmupWait, err = config.DurationWithDefault(WarmupWait, 5*time.Second)
	if err != nil {
		return err
	}

	s.defaultReport = config.StringWithDefault(DefaultReport, ".")
	s.lessThanReport = config.StringWithDefault(LessThanReport, "<")
	s.errorReport = config.StringWithDefault(ErrorReport, "#")
	s.greaterThanReport = config.StringWithDefault(GreaterThanReport, ">")

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
	for i := 0; i < s.warmupMessagesPerRun; i++ {
		writer.Write(s.buffer)
	}

	if s.warmupMessagesPerRun > 0 {
		time.Sleep(s.warmupWait)
	}

	startTime := time.Now()
	for i := 0; i < s.messagesPerRun; i++ {
		s.countMessage(writer.Write(s.buffer))
	}
	s.runTime = time.Now().Sub(startTime)

	err := writer.Close()
	if err != nil {
		logs.Logger.Warning("Error closing writer, %s", err.Error())
	}
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

	err := reader.Close()
	if err != nil {
		logs.Logger.Warning("Error closing reader, %s", err.Error())
	}

	wg.Wait()
}

func (s *Scheme) runReader(reader factory.Reader, timer *time.Timer) {
	var startTime, lastMessageTime time.Time
	buffer := make([]byte, s.bytesPerMessage*2)
	for i := 0; i < s.warmupMessagesPerRun; i++ {
		reader.Read(buffer)
	}

	for {
		count, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		lastMessageTime = time.Now()
		if startTime.IsZero() {
			startTime = lastMessageTime
		}

		timer.Reset(s.waitForLastMessage)

		s.countMessage(count, err)
	}

	s.runTime = lastMessageTime.Sub(startTime)
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

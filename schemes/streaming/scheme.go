package streaming

import (
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/myshkin5/jsonstruct"
	"github.com/myshkin5/netspel/factory"
	"github.com/myshkin5/netspel/logs"
)

const (
	prefix = ".streaming."

	MessagesPerSecond         = prefix + "messages-per-second"
	ExpectedMessagesPerSecond = prefix + "expected-messages-per-second"
	BytesPerMessage           = prefix + "bytes-per-message"
	ReportCycle               = prefix + "report-cycle"

	DefaultMessagesPerSecond         = 1000
	DefaultExpectedMessagesPerSecond = 0
	DefaultBytesPerMessage           = 1024
	DefaultReportCycle               = time.Second
)

type Scheme struct {
	buffer       []byte
	messageCount uint32
	byteCount    uint64
	errorCount   uint32

	messagesPerSecond int
	bytesPerMessage   int
	reportCycle       time.Duration

	tickerTime time.Duration
	closed     int32
	done       sync.WaitGroup
	closer     io.Closer
	reporter   Reporter
}

type Reporter interface {
	Init(expectedMessagesPerSecond int, reportCycle time.Duration)
	Report(report Report)
}

func (s *Scheme) Init(config jsonstruct.JSONStruct) error {
	s.messagesPerSecond = config.IntWithDefault(MessagesPerSecond, DefaultMessagesPerSecond)
	expectedMessagesPerSecond := config.IntWithDefault(ExpectedMessagesPerSecond, DefaultExpectedMessagesPerSecond)
	s.bytesPerMessage = config.IntWithDefault(BytesPerMessage, DefaultBytesPerMessage)

	var err error
	s.reportCycle, err = config.DurationWithDefault(ReportCycle, DefaultReportCycle)
	if err != nil {
		return err
	}

	if expectedMessagesPerSecond == 0 {
		expectedMessagesPerSecond = s.messagesPerSecond
	}

	s.buffer = make([]byte, s.bytesPerMessage)
	if s.messagesPerSecond > 0 {
		s.tickerTime = time.Second / time.Duration(s.messagesPerSecond)
	}

	s.done.Add(1)

	if s.reporter == nil {
		s.reporter = &ReporterImpl{}
	}
	s.reporter.Init(expectedMessagesPerSecond, s.reportCycle)

	return nil
}

func (s *Scheme) SetReporter(reporter Reporter) {
	s.reporter = reporter
}

func (s *Scheme) RunWriter(writer factory.Writer) {
	s.closer = writer
	defer s.done.Done()
	s.startReporter()

	var ticker *time.Ticker
	if s.tickerTime > 0 {
		ticker = time.NewTicker(s.tickerTime)
		defer ticker.Stop()
	}

	for {
		if ticker != nil {
			<-ticker.C
		}

		if s.isClosed() {
			break
		}

		count, err := writer.Write(s.buffer)
		s.countMessage(count, err)
	}
}

func (s *Scheme) RunReader(reader factory.Reader) {
	s.closer = reader
	defer s.done.Done()
	s.startReporter()

	buffer := make([]byte, s.bytesPerMessage*2)

	var ticker *time.Ticker
	if s.tickerTime > 0 {
		ticker = time.NewTicker(s.tickerTime)
		defer ticker.Stop()
	}

	for {
		if ticker != nil {
			<-ticker.C
		}

		if s.isClosed() {
			break
		}

		count, err := reader.Read(buffer)
		s.countMessage(count, err)
	}
}

func (s *Scheme) Close() error {
	atomic.StoreInt32(&s.closed, 1)
	if s.closer != nil {
		err := s.closer.Close()
		if err != nil {
			return err
		}
	}
	s.done.Wait()
	return nil
}

func (s *Scheme) isClosed() bool {
	return atomic.LoadInt32(&s.closed) == 1
}

func (s *Scheme) countMessage(count int, err error) {
	if err != nil {
		logs.Logger.Debug("Adapter error, %v", err)
	}
	if count > 0 {
		atomic.AddUint32(&s.messageCount, 1)
		atomic.AddUint64(&s.byteCount, uint64(count))
	}
	if err != nil {
		atomic.AddUint32(&s.errorCount, 1)
	}
}

func (s *Scheme) startReporter() {
	s.done.Add(1)

	go func() {
		defer s.done.Done()

		ticker := time.NewTicker(s.reportCycle)
		for {
			<-ticker.C
			if s.isClosed() {
				break
			}

			report := Report{}
			report.MessageCount = atomic.SwapUint32(&s.messageCount, 0)
			report.ByteCount = atomic.SwapUint64(&s.byteCount, 0)
			report.ErrorCount = atomic.SwapUint32(&s.errorCount, 0)
			s.reporter.Report(report)
		}
	}()
}

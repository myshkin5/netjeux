package simple

import (
	"io"
	"sync"
	"time"

	"github.com/myshkin5/jsonstruct"
	"github.com/myshkin5/netspel/factory"
	"github.com/myshkin5/netspel/logs"
	"github.com/myshkin5/netspel/utils"
)

const (
	prefix = ".simple."

	MessagesPerRun     = prefix + "messages-per-run"
	BytesPerMessage    = prefix + "bytes-per-message"
	WaitForLastMessage = prefix + "wait-for-last-message"

	WarmupMessagesPerRun = prefix + "warmup-messages-per-run"
	WarmupWait           = prefix + "warmup-wait"

	DefaultMessagesPerRun     = 10000
	DefaultBytesPerMessage    = 1024
	DefaultWaitForLastMessage = 5 * time.Second

	DefaultWarmupMessagesPerRun = 0
	DefaultWarmupWait           = 5 * time.Second
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
}

func (s *Scheme) Init(config jsonstruct.JSONStruct) error {
	s.messagesPerRun = config.IntWithDefault(MessagesPerRun, DefaultMessagesPerRun)
	s.bytesPerMessage = config.IntWithDefault(BytesPerMessage, DefaultBytesPerMessage)
	s.buffer = make([]byte, s.bytesPerMessage)

	var err error
	s.waitForLastMessage, err = config.DurationWithDefault(WaitForLastMessage, DefaultWaitForLastMessage)
	if err != nil {
		return err
	}

	s.warmupMessagesPerRun = config.IntWithDefault(WarmupMessagesPerRun, DefaultWarmupMessagesPerRun)
	s.warmupWait, err = config.DurationWithDefault(WarmupWait, DefaultWarmupWait)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scheme) ByteCount() uint64 {
	return s.byteCount
}

func (s *Scheme) ErrorCount() uint32 {
	return s.errorCount
}

func (s *Scheme) FirstError() error {
	return s.firstError
}

func (s *Scheme) RunTime() time.Duration {
	return s.runTime
}

func (s *Scheme) RunWriter(writer factory.Writer) {
	if s.warmupMessagesPerRun > 0 {
		logs.Logger.Info("Writing %d warmup messages", s.warmupMessagesPerRun)
	}

	for i := 0; i < s.warmupMessagesPerRun; i++ {
		writer.Write(s.buffer)
	}

	if s.warmupMessagesPerRun > 0 {
		time.Sleep(s.warmupWait)
	}

	logs.Logger.Info("Starting writing %d messages...", s.messagesPerRun)
	startTime := time.Now()
	for i := 0; i < s.messagesPerRun; i++ {
		s.countMessage(writer.Write(s.buffer))
	}
	s.runTime = time.Now().Sub(startTime)
	logs.Logger.Info("Finished.")

	err := writer.Close()
	if err != nil {
		logs.Logger.Warning("Error closing writer, %s", err.Error())
	}

	s.outputReport()
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

	s.outputReport()
}

func (s *Scheme) runReader(reader factory.Reader, timer *time.Timer) {
	if s.warmupMessagesPerRun > 0 {
		logs.Logger.Info("Reading %d warmup messages", s.warmupMessagesPerRun)
	}

	var startTime, lastMessageTime time.Time
	buffer := make([]byte, s.bytesPerMessage*2)
	for i := 0; i < s.warmupMessagesPerRun; i++ {
		reader.Read(buffer)
	}

	logs.Logger.Info("Starting reading %d messages...", s.messagesPerRun)
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
	logs.Logger.Info("Finished.")

	s.runTime = lastMessageTime.Sub(startTime)
}

func (s *Scheme) countMessage(count int, err error) {
	s.byteCount += uint64(count)
	if err != nil {
		s.errorCount++
		if s.firstError == nil {
			s.firstError = err
		}
	}
}

func (s *Scheme) outputReport() {
	bytesPerSec := utils.ByteSize(s.ByteCount()) * utils.ByteSize(time.Second) / utils.ByteSize(s.RunTime().Nanoseconds())
	messagesPerSec := float64(s.messagesPerRun) * float64(time.Second) / float64(s.RunTime().Nanoseconds())

	logs.Logger.Info("Byte count: %d", s.ByteCount())
	logs.Logger.Info("Rates: %s/s %.1f messages/s", bytesPerSec.String(), messagesPerSec)
	logs.Logger.Info("Error count: %d", s.ErrorCount())
	logs.Logger.Info("Run time: %s", s.RunTime().String())
	if s.FirstError() != nil {
		logs.Logger.Info("First error: %s", s.FirstError().Error())
	}
}

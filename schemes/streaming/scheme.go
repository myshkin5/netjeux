package streaming

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/myshkin5/netspel/factory"
	"github.com/myshkin5/netspel/jsonstruct"
	"github.com/myshkin5/netspel/logs"
	"github.com/myshkin5/netspel/utils"
)

const (
	prefix = "streaming."

	MessagesPerSecond = prefix + "messages-per-second"
	BytesPerMessage   = prefix + "bytes-per-message"
)

type Scheme struct {
	buffer       []byte
	messageCount uint32
	byteCount    uint64
	errorCount   uint32

	messagesPerSecond int
	bytesPerMessage   int

	tickerTime time.Duration
	closed     int32
	done       sync.WaitGroup
	closer     io.Closer
}

func (s *Scheme) Init(config jsonstruct.JSONStruct) error {
	var ok bool
	s.bytesPerMessage, ok = config.Int(BytesPerMessage)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", BytesPerMessage)
	}
	s.buffer = make([]byte, s.bytesPerMessage)

	s.messagesPerSecond, ok = config.Int(MessagesPerSecond)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", MessagesPerSecond)
	}

	s.tickerTime = time.Second / time.Duration(s.messagesPerSecond)

	s.done.Add(1)

	return nil
}

func (s *Scheme) BytesPerMessage() int {
	return s.bytesPerMessage
}

func (s *Scheme) MessagesPerSecond() int {
	return s.messagesPerSecond
}

func (s *Scheme) MessageCount() uint32 {
	return atomic.LoadUint32(&s.messageCount)
}

func (s *Scheme) ByteCount() uint64 {
	return atomic.LoadUint64(&s.byteCount)
}

func (s *Scheme) ErrorCount() uint32 {
	return atomic.LoadUint32(&s.errorCount)
}

func (s *Scheme) RunWriter(writer factory.Writer) {
	s.closer = writer
	defer s.done.Done()
	s.startReporter()

	ticker := time.NewTicker(s.tickerTime)
	defer ticker.Stop()

	for {
		<-ticker.C
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

	for {
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

		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			if s.isClosed() {
				break
			}

			messageCount := atomic.SwapUint32(&s.messageCount, 0)
			byteCount := atomic.SwapUint64(&s.byteCount, 0)
			errorCount := atomic.SwapUint32(&s.errorCount, 0)
			logs.Logger.Info("Message count: %7d, Error count: %7d, Byte count: %s", messageCount, errorCount, utils.ByteSize(byteCount).String())
		}
	}()
}

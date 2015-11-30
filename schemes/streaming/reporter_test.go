package streaming_test

import (
	"github.com/myshkin5/netspel/schemes/streaming"

	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporter", func() {
	var (
		logger mockLogger
	)

	BeforeEach(func() {
		logger = mockLogger{
			logs: make(chan string, 100),
		}
		streaming.ReporterLogger = &logger
	})

	expectLog := func(log string, messageCount uint32, byteCount uint64, errorCount uint32, expectedMessagesPerSecond int, reportCycle time.Duration) {
		reporter := streaming.ReporterImpl{}
		reporter.Init(expectedMessagesPerSecond, reportCycle)
		reporter.Report(streaming.Report{
			MessageCount: messageCount,
			ByteCount:    byteCount,
			ErrorCount:   errorCount,
		})
		ExpectWithOffset(1, logger.logs).To(Receive(Equal(log)))
	}

	It("reports to the logger", func() {
		expectLog("       0 messages/s (   NaN%),        0 errors/s, 0.00 B/s", 0, 0, 0, 0, time.Second)
		expectLog("     100 messages/s (  +Inf%),       10 errors/s, 1.00 KB/s", 100, 1024, 10, 0, time.Second)
		expectLog("     100 messages/s (100.00%),       10 errors/s, 1.00 KB/s", 100, 1024, 10, 100, time.Second)
		expectLog("      50 messages/s ( 50.00%),       10 errors/s, 1.00 KB/s", 50, 1024, 10, 100, time.Second)
		expectLog("     500 messages/s ( 50.00%),      100 errors/s, 10.00 KB/s", 50, 1024, 10, 1000, 100*time.Millisecond)
		expectLog("       5 messages/s (  0.50%),        1 errors/s, 1.00 KB/s", 50, 10240, 10, 1000, 10*time.Second)
	})
})

type mockLogger struct {
	logs chan string
}

func (m *mockLogger) Info(format string, args ...interface{}) {
	m.logs <- fmt.Sprintf(format, args...)
}

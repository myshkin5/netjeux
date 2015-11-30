package streaming_test

import (
	"github.com/myshkin5/jsonstruct"
	"github.com/myshkin5/netspel/schemes/internal/mocks"
	"github.com/myshkin5/netspel/schemes/streaming"

	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scheme", func() {
	var (
		writer   *mocks.MockWriter
		reader   *mocks.MockReader
		scheme   *streaming.Scheme
		config   jsonstruct.JSONStruct
		reporter *mockReporter
	)

	BeforeEach(func() {
		writer = mocks.NewMockWriter()
		reader = mocks.NewMockReader()
		scheme = &streaming.Scheme{}
		config = jsonstruct.New()
		reporter = &mockReporter{
			reports: make(chan streaming.Report, 100),
		}
		scheme.SetReporter(reporter)
	})

	emptyReport := streaming.Report{
		MessageCount: 0,
		ByteCount:    0,
		ErrorCount:   0,
	}
	singleMessageReport := streaming.Report{
		MessageCount: 1,
		ByteCount:    1024,
		ErrorCount:   0,
	}

	Context("with a 5-messages-per-second configuration", func() {
		BeforeEach(func() {
			// one message per 100ms
			config.SetInt(streaming.MessagesPerSecond, 10)
			// will alternate reports of 0 and 1 messages per report
			config.SetDuration(streaming.ReportCycle, 80*time.Millisecond)

			err := scheme.Init(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets up the reporter properly", func() {
			Expect(reporter.expectedMessagesPerSecond).To(Equal(10))
			Expect(reporter.reportCycle).To(Equal(80 * time.Millisecond))
		})

		It("writes messages at the rate specified", func() {
			go scheme.RunWriter(writer)

			Eventually(reporter.reports, 100*time.Millisecond).Should(Receive(Equal(emptyReport)))
			Eventually(reporter.reports, 100*time.Millisecond).Should(Receive(Equal(singleMessageReport)))

			scheme.Close()

			Expect(len(writer.Messages)).To(Equal(1))
		})

		It("reads messages at the rate specified", func() {
			for i := 0; i < 100; i++ {
				reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 1024), Error: nil}
			}

			go scheme.RunReader(reader)

			Eventually(reporter.reports, 100*time.Millisecond).Should(Receive(Equal(emptyReport)))
			Eventually(reporter.reports, 100*time.Millisecond).Should(Receive(Equal(singleMessageReport)))

			scheme.Close()
		})
	})

	Context("with a zero-messages-per-second (infinite) configuration", func() {
		BeforeEach(func() {
			config.SetInt(streaming.MessagesPerSecond, 0)
			config.SetDuration(streaming.ReportCycle, 50*time.Millisecond)

			err := scheme.Init(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets up the reporter properly", func() {
			Expect(reporter.expectedMessagesPerSecond).To(Equal(0))
		})

		It("reads messages as quickly as possible", func() {
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 10), Error: nil}
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 1000), Error: nil}

			go scheme.RunReader(reader)

			Eventually(reporter.reports, 100*time.Millisecond).Should(Receive(Equal(streaming.Report{
				MessageCount: 2,
				ByteCount:    1010,
				ErrorCount:   0,
			})))
			Eventually(reporter.reports, 60*time.Millisecond).Should(Receive(Equal(emptyReport)))

			scheme.Close()

			Expect(reader.ReadMessages).To(BeEmpty())
		})

		It("writes messages as quickly as possible", func() {
			go scheme.RunWriter(writer)

			Eventually(reporter.reports, 100*time.Millisecond).Should(Receive(Equal(streaming.Report{
				MessageCount: 10000,
				ByteCount:    10240000,
				ErrorCount:   0,
			})))
			Eventually(reporter.reports, 60*time.Millisecond).Should(Receive(Equal(emptyReport)))

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				scheme.Close()
			}()

			time.Sleep(100 * time.Millisecond)
			<-writer.Messages

			wg.Wait()
		})
	})
})

type mockReporter struct {
	expectedMessagesPerSecond int
	reportCycle               time.Duration
	reports                   chan streaming.Report
}

func (m *mockReporter) Init(expectedMessagesPerSecond int, reportCycle time.Duration) {
	m.expectedMessagesPerSecond = expectedMessagesPerSecond
	m.reportCycle = reportCycle
}

func (m *mockReporter) Report(report streaming.Report) {
	m.reports <- report
}

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
		writer *mocks.MockWriter
		reader *mocks.MockReader
		scheme *streaming.Scheme
		config jsonstruct.JSONStruct
	)

	BeforeEach(func() {
		writer = mocks.NewMockWriter()
		reader = mocks.NewMockReader()
		scheme = &streaming.Scheme{}
		config = jsonstruct.New()
		config.SetInt(streaming.BytesPerMessage, 1000)
	})

	Context("with a 5-messages-per-second configuration", func() {
		BeforeEach(func() {
			config.SetInt(streaming.MessagesPerSecond, 5)

			err := scheme.Init(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("writes messages at the rate specified", func() {
			go scheme.RunWriter(writer)

			time.Sleep(500 * time.Millisecond)
			Expect(scheme.ByteCount()).To(BeNumerically(">", 1000))

			time.Sleep(600 * time.Millisecond)
			Expect(scheme.ByteCount()).To(BeNumerically("<=", 1000))
			scheme.Close()

			Expect(len(writer.Messages)).To(BeNumerically(">=", 4))
			Expect(len(writer.Messages)).To(BeNumerically("<=", 6))
		})

		It("reads messages at the rate specified", func() {
			for i := 0; i < 10; i++ {
				reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 1000), Error: nil}
			}

			go scheme.RunReader(reader)

			time.Sleep(500 * time.Millisecond)
			Expect(scheme.ByteCount()).To(BeNumerically(">", 1000))

			time.Sleep(600 * time.Millisecond)
			Expect(scheme.ByteCount()).To(BeNumerically("<=", 1000))
			scheme.Close()

			Expect(len(reader.ReadMessages)).To(BeNumerically(">=", 4))
			Expect(len(reader.ReadMessages)).To(BeNumerically("<=", 6))
		})
	})

	Context("with a zero-messages-per-second (infinite) configuration", func() {
		BeforeEach(func() {
			config.SetInt(streaming.MessagesPerSecond, 0)

			err := scheme.Init(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("reads messages as quickly as possible", func() {
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 10), Error: nil}
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 1000), Error: nil}

			go scheme.RunReader(reader)

			time.Sleep(100 * time.Millisecond)
			Expect(scheme.ByteCount()).To(BeEquivalentTo(1010))

			time.Sleep(time.Second)
			Expect(scheme.ByteCount()).To(BeEquivalentTo(0))

			scheme.Close()

			Expect(reader.ReadMessages).To(BeEmpty())
		})

		It("writes messages as quickly as possible", func() {
			go scheme.RunWriter(writer)

			time.Sleep(100 * time.Millisecond)
			Expect(scheme.ByteCount()).To(BeEquivalentTo(10000000))

			time.Sleep(time.Second)
			Expect(scheme.ByteCount()).To(BeEquivalentTo(0))

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

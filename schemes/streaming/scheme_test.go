package streaming_test

import (
	"github.com/myshkin5/netspel/jsonstruct"
	"github.com/myshkin5/netspel/schemes/internal/mocks"
	"github.com/myshkin5/netspel/schemes/streaming"

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
		config.SetInt(streaming.MessagesPerSecond, 5)
		config.SetInt(streaming.BytesPerMessage, 1000)

		err := scheme.Init(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns the values it is configured with", func() {
		Expect(scheme.MessagesPerSecond()).To(Equal(5))
		Expect(scheme.BytesPerMessage()).To(Equal(1000))
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
})

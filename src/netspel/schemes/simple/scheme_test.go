package simple_test

import (
	"errors"

	"netspel/factory"
	"netspel/schemes/internal/mocks"
	"netspel/schemes/simple"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scheme", func() {
	var (
		writer *mocks.MockWriter
		reader *mocks.MockReader
		scheme *simple.Scheme
	)

	BeforeEach(func() {
		writer = mocks.NewMockWriter()
		reader = mocks.NewMockReader()
		scheme = &simple.Scheme{}
		err := scheme.Init(factory.Config{})
		Expect(err).NotTo(HaveOccurred())
		scheme.DefaultReport = ""
		scheme.LessThanReport = ""
		scheme.ErrorReport = ""
	})

	It("writes messages to a writer", func() {
		go scheme.RunWriter(writer)

		var sentMessage []byte
		for i := 0; i < simple.MessagesPerRun; i++ {
			Eventually(writer.Messages).Should(Receive(&sentMessage))
			Expect(sentMessage).To(HaveLen(simple.BytesPerMessage))
		}

		Expect(scheme.ByteCount()).To(BeEquivalentTo(simple.MessagesPerRun * simple.BytesPerMessage))
		Expect(scheme.ErrorCount()).To(BeZero())
	})

	It("reads messages from a reader", func() {
		reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 10), Error: nil}
		reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 1000), Error: nil}
		reader.ReadMessages <- mocks.ReadMessage{Buffer: []byte{}, Error: errors.New("Bad stuff")}

		go scheme.RunReader(reader)

		Eventually(scheme.ByteCount).Should(BeEquivalentTo(1010))
		Eventually(scheme.ErrorCount).Should(BeEquivalentTo(1))
	})
})

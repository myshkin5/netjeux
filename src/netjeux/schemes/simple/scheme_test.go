package simple_test

import (
	"netjeux/schemes/simple"

	"errors"
	"netjeux/schemes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scheme", func() {
	var (
		writer *schemes.MockWriter
		reader *schemes.MockReader
		scheme *simple.Scheme
	)

	BeforeEach(func() {
		writer = schemes.NewMockWriter()
		reader = schemes.NewMockReader()
		scheme = simple.New(writer, reader)
		scheme.DefaultReport = ""
		scheme.LessThanReport = ""
		scheme.ErrorReport = ""
	})

	It("writes messages to a writer", func() {
		go scheme.RunWriter()

		var sentMessage []byte
		for i := 0; i < simple.MessagesPerRun; i++ {
			Eventually(writer.Messages).Should(Receive(&sentMessage))
			Expect(sentMessage).To(HaveLen(simple.BytesPerMessage))
		}

		Expect(scheme.ByteCount()).To(BeEquivalentTo(simple.MessagesPerRun * simple.BytesPerMessage))
		Expect(scheme.ErrorCount()).To(BeZero())
	})

	It("reads messages from a reader", func() {
		reader.ReadMessages <- schemes.ReadMessage{Buffer: make([]byte, 100), Error: nil}
		reader.ReadMessages <- schemes.ReadMessage{Buffer: make([]byte, 10000), Error: nil}
		reader.ReadMessages <- schemes.ReadMessage{Buffer: []byte{}, Error: errors.New("Bad stuff")}

		go scheme.RunReader()

		Eventually(scheme.ByteCount).Should(BeEquivalentTo(10100))
		Eventually(scheme.ErrorCount).Should(BeEquivalentTo(1))
	})
})

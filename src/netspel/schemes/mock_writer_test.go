package schemes_test

import (
	"netspel/schemes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockWriter", func() {
	It("stores all writes in order", func() {
		writer := schemes.NewMockWriter()

		message1 := []byte("message 1")
		bytesWritten, err := writer.Write(message1)
		Expect(err).NotTo(HaveOccurred())
		Expect(bytesWritten).To(Equal(len(message1)))

		message2 := []byte("message 2 - with more stuff")
		bytesWritten, err = writer.Write(message2)
		Expect(err).NotTo(HaveOccurred())
		Expect(bytesWritten).To(Equal(len(message2)))

		var sentMessage []byte
		Expect(writer.Messages).To(Receive(&sentMessage))
		Expect(sentMessage).To(Equal(message1))

		Expect(writer.Messages).To(Receive(&sentMessage))
		Expect(sentMessage).To(Equal(message2))
	})
})

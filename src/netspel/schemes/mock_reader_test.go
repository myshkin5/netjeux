package schemes_test

import (
	"netspel/schemes"

	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockReader", func() {
	It("returns all messages", func() {
		reader := schemes.NewMockReader()
		message1 := []byte("message 1")
		reader.ReadMessages <- schemes.ReadMessage{Buffer: message1, Error: nil}
		message2 := []byte("message 2 - something else")
		reader.ReadMessages <- schemes.ReadMessage{Buffer: message2, Error: nil}
		reader.ReadMessages <- schemes.ReadMessage{Buffer: []byte{}, Error: errors.New("Bad stuff")}

		buffer := make([]byte, 30)
		bytesRead, err := reader.Read(buffer)
		Expect(bytesRead).To(Equal(len(message1)))
		Expect(err).NotTo(HaveOccurred())
		Expect(buffer[0:bytesRead]).To(Equal(message1))

		bytesRead, err = reader.Read(buffer)
		Expect(bytesRead).To(Equal(len(message2)))
		Expect(err).NotTo(HaveOccurred())
		Expect(buffer[0:bytesRead]).To(Equal(message2))

		bytesRead, err = reader.Read(buffer)
		Expect(err).To(HaveOccurred())
	})
})

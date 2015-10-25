package udp_test

import (
	"netjeux/adapters/udp"

	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	It("reads from a UDP port", func() {
		reader, err := udp.NewReader(51040)
		Expect(err).NotTo(HaveOccurred())

		messageSent := []byte("hello")
		messages := make(chan []byte, 100)
		go func() {
			messageRead := make([]byte, 1024)
			for {
				bytesRead, err := reader.Read(messageRead)
				Expect(err).NotTo(HaveOccurred())
				send := make([]byte, bytesRead)
				copy(send, messageRead[0:bytesRead])
				messages <- send
			}
		}()

		raddr, err := net.ResolveUDPAddr("udp4", "localhost:51040")
		Expect(err).NotTo(HaveOccurred())

		connection, err := net.DialUDP("udp4", nil, raddr)
		Expect(err).NotTo(HaveOccurred())

		bytesWritten, err := connection.Write(messageSent)
		Expect(err).NotTo(HaveOccurred())
		Expect(bytesWritten).To(Equal(len(messageSent)))

		var messageRead []byte
		Eventually(messages).Should(Receive(&messageRead))
		Expect(messageRead).To(Equal(messageSent))
	})
})

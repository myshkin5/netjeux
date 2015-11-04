package udp_test

import (
	"netspel/adapters/udp"
	"netspel/factory"

	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Writer", func() {
	It("writes to a UDP port", func() {
		laddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", 51041))
		Expect(err).NotTo(HaveOccurred())

		connection, err := net.ListenUDP("udp4", laddr)
		Expect(err).NotTo(HaveOccurred())

		messages := make(chan []byte, 100)
		go func() {
			messageRead := make([]byte, 1024)
			for {
				bytesRead, err := connection.Read(messageRead)
				Expect(err).NotTo(HaveOccurred())
				message := make([]byte, bytesRead)
				copy(message, messageRead[0:bytesRead])
				messages <- message
			}
		}()

		config := factory.NewConfig()
		config.ParseAndSetAdditionalInt(udp.Port + "=51041")
		config.ParseAndSetAdditionalString(udp.RemoteAddr + "=localhost")

		writer := udp.Writer{}
		err = writer.Init(*config)
		Expect(err).NotTo(HaveOccurred())

		messageSent := []byte("hello")
		bytesWritten, err := writer.Write(messageSent)
		Expect(err).NotTo(HaveOccurred())
		Expect(bytesWritten).To(Equal(len(messageSent)))

		var messageRead []byte
		Eventually(messages).Should(Receive(&messageRead))
		Expect(messageRead).To(Equal(messageSent))
	})
})

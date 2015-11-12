package udp_test

import (
	"fmt"
	"io"
	"net"
	"time"

	"netspel/adapters/udp"
	"netspel/jsonstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	var (
		port   int
		reader udp.Reader
	)

	BeforeEach(func() {
		if port == 0 {
			port = 51010
		}
		port++
		config := jsonstruct.JSONStruct(make(map[string]interface{}))
		config.SetInt(udp.Port, port)

		reader = udp.Reader{}
		err := reader.Init(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("reads from a UDP port", func() {
		done := make(chan struct{})
		messages := make(chan []byte, 100)
		go func() {
			defer GinkgoRecover()
			messageRead := make([]byte, 1024)
			for {
				bytesRead, err := reader.Read(messageRead)
				if err == io.EOF {
					break
				}
				Expect(err).NotTo(HaveOccurred())
				message := make([]byte, bytesRead)
				copy(message, messageRead[0:bytesRead])
				messages <- message
			}
			close(done)
		}()

		raddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("localhost:%d", port))
		Expect(err).NotTo(HaveOccurred())

		connection, err := net.DialUDP("udp4", nil, raddr)
		Expect(err).NotTo(HaveOccurred())

		messageSent := []byte("hello")
		bytesWritten, err := connection.Write(messageSent)
		Expect(err).NotTo(HaveOccurred())
		Expect(bytesWritten).To(Equal(len(messageSent)))

		var messageRead []byte
		Eventually(messages).Should(Receive(&messageRead))
		Expect(messageRead).To(Equal(messageSent))

		err = reader.Close()
		Expect(err).NotTo(HaveOccurred())
		Eventually(done).Should(BeClosed())
	})

	It("cancels a read when told to stop", func() {
		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			messageRead := make([]byte, 1024)
			bytesRead, err := reader.Read(messageRead)
			Expect(err).To(Equal(io.EOF))
			Expect(bytesRead).To(Equal(0))
			close(done)
		}()

		time.Sleep(50 * time.Millisecond)

		err := reader.Close()
		Expect(err).NotTo(HaveOccurred())
		Eventually(done).Should(BeClosed())
	})
})

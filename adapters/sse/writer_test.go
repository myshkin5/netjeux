package sse_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/myshkin5/netspel/adapters/sse"
	"github.com/myshkin5/netspel/jsonstruct"
	vitosse "github.com/vito/go-sse/sse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Writer", func() {
	var (
		port   int
		writer *sse.Writer
	)

	writerReady := func() bool {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ready", port))
		if err != nil {
			return false
		}

		if resp.StatusCode != http.StatusOK {
			return false
		}

		return true
	}

	BeforeEach(func() {
		if port == 0 {
			port = 29687
		}
		port++

		config := jsonstruct.New()
		config.SetInt(sse.Port, port)

		writer = &sse.Writer{}
		err := writer.Init(config)
		Expect(err).NotTo(HaveOccurred())

		Eventually(writerReady).Should(BeTrue())
	})

	It("returns the proper SSE headers", func() {
		defer writer.Close()

		resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(resp.Header.Get("Content-Type")).To(Equal("text/event-stream; charset=utf-8"))
		Expect(resp.Header.Get("Cache-Control")).To(Equal("no-cache"))
		Expect(resp.Header.Get("Connection")).To(Equal("keep-alive"))
	})

	It("blocks a writer until a reader starts reading", func() {
		defer writer.Close()

		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			writer.Write(make([]byte, 1024))
			close(done)
		}()

		Consistently(done).ShouldNot(BeClosed())
	})

	It("sends data to a reader", func(ginkgoDone Done) {
		defer close(ginkgoDone)

		message1 := make([]byte, 10)
		message2 := make([]byte, 1024)
		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()

			count, err := writer.Write(message1)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(10))

			count, err = writer.Write(message2)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1024))

			err = writer.Close()
			Expect(err).NotTo(HaveOccurred())

			close(done)
		}()

		resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
		Expect(err).NotTo(HaveOccurred())

		reader := vitosse.NewReadCloser(resp.Body)

		event, err := reader.Next()
		Expect(err).NotTo(HaveOccurred())
		Expect(event.Data).To(Equal(message1))

		event, err = reader.Next()
		Expect(err).NotTo(HaveOccurred())
		Expect(event.Data).To(Equal(message2))

		Eventually(done).Should(BeClosed())
	}, 10)

	It("stops attempting to write when closed", func() {
		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()

			count, err := writer.Write(make([]byte, 10))
			Expect(err).To(HaveOccurred())
			Expect(count).To(Equal(0))

			close(done)
		}()

		// Test needs to wait for Write() to be called
		time.Sleep(10 * time.Millisecond)

		err := writer.Close()
		Expect(err).NotTo(HaveOccurred())

		Eventually(done).Should(BeClosed())
	})

	It("closes when there are no writes pending", func() {
		err := writer.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can be reconnected after a reader completes and closes the connection", func() {
		mCount := 0

		readWrite := func() {
			message := make([]byte, 1024)
			done := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				message[0] = byte(mCount)
				mCount++
				count, err := writer.Write(message)
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(1024))

				close(done)
			}()

			resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
			Expect(err).NotTo(HaveOccurred())

			reader := vitosse.NewReadCloser(resp.Body)

			event, err := reader.Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(event.Data).To(Equal(message))

			err = reader.Close()
			Expect(err).NotTo(HaveOccurred())

			Eventually(done).Should(BeClosed())
		}

		readWrite()

		// Need the first request to close before we write another message
		time.Sleep(100 * time.Millisecond)

		readWrite()

		err := writer.Close()
		Expect(err).NotTo(HaveOccurred())
	})
})

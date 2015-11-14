package sse_test

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"netspel/adapters/sse"
	"netspel/jsonstruct"

	vitosse "github.com/vito/go-sse/sse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	var (
		port            int
		sleepBeforeSend time.Duration
		events          []*vitosse.Event
		config          jsonstruct.JSONStruct
		reader          sse.Reader
	)

	handle := func(w http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()

		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("\n"))
		flusher := w.(http.Flusher)
		flusher.Flush()

		time.Sleep(sleepBeforeSend)

		for _, event := range events {
			err := event.Write(w)
			Expect(err).NotTo(HaveOccurred())
			flusher.Flush()
		}
	}

	BeforeEach(func() {
		if port == 0 {
			port = 49282
		}
		port++

		go func() {
			defer GinkgoRecover()
			err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), http.HandlerFunc(handle))
			Expect(err).NotTo(HaveOccurred())
		}()

		events = []*vitosse.Event{}

		time.Sleep(50 * time.Millisecond)

		config = jsonstruct.New()
		config.SetString(sse.RemoteAddr, "localhost")
		config.SetInt(sse.Port, port)

		reader = sse.Reader{}
	})

	It("connects to a writer via HTTP and reads SSE messages", func() {
		event := &vitosse.Event{
			ID:   "event-id",
			Name: "event-name",
			Data: make([]byte, 100),
		}
		events = append(events, event)
		events = append(events, event)

		err := reader.Init(config)
		Expect(err).NotTo(HaveOccurred())

		message := make([]byte, 200)
		count, err := reader.Read(message)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(100))

		count, err = reader.Read(message)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(Equal(100))

		count, err = reader.Read(message)
		Expect(err).To(Equal(io.EOF))

		err = reader.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns from a call to Read() when Close() is called", func() {
		err := reader.Init(config)
		Expect(err).NotTo(HaveOccurred())

		sleepBeforeSend = time.Second
		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			messageRead := make([]byte, 1024)
			bytesRead, err := reader.Read(messageRead)
			Expect(err).To(Equal(io.EOF))
			Expect(bytesRead).To(Equal(0))
			close(done)
		}()

		time.Sleep(10 * time.Millisecond)

		err = reader.Close()
		Expect(err).NotTo(HaveOccurred())
		Eventually(done).Should(BeClosed())
	})
})

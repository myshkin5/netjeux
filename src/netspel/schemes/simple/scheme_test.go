package simple_test

import (
	"netspel/schemes/internal/mocks"
	"netspel/schemes/simple"

	"errors"
	"sync"
	"time"

	"netspel/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scheme", func() {
	var (
		writer *mocks.MockWriter
		reader *mocks.MockReader
		scheme *simple.Scheme
		config map[string]interface{}
	)

	BeforeEach(func() {
		writer = mocks.NewMockWriter()
		reader = mocks.NewMockReader()
		scheme = &simple.Scheme{}
		config = make(map[string]interface{})
		json.SetInt(simple.BytesPerMessage, 1000, config)
		json.SetInt(simple.MessagesPerRun, 100, config)
		json.SetString(simple.WaitForLastMessage, "100ms", config)
		json.SetString(simple.DefaultReport, "", config)
		json.SetString(simple.LessThanReport, "", config)
		json.SetString(simple.ErrorReport, "", config)
		json.SetString(simple.GreaterThanReport, "", config)
	})

	Context("with a short wait time", func() {
		JustBeforeEach(func() {
			json.SetString(simple.WaitForLastMessage, "100ms", config)
			err := scheme.Init(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the values it is configured with", func() {
			Expect(scheme.BytesPerMessage()).To(Equal(1000))
			Expect(scheme.MessagesPerRun()).To(Equal(100))
		})

		It("writes messages to a writer", func() {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				scheme.RunWriter(writer)
			}()

			var sentMessage []byte
			for i := 0; i < 100; i++ {
				Eventually(writer.Messages).Should(Receive(&sentMessage))
				Expect(sentMessage).To(HaveLen(1000))
			}

			wg.Wait()

			Expect(scheme.ByteCount()).To(BeEquivalentTo(100 * 1000))
			Expect(scheme.ErrorCount()).To(BeZero())
			Expect(scheme.RunTime()).To(BeNumerically(">", 0))
		})

		It("reads messages from a reader", func() {
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 10), Error: nil}
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 1000), Error: nil}
			firstError := errors.New("Bad stuff")
			reader.ReadMessages <- mocks.ReadMessage{Buffer: []byte{}, Error: firstError}

			scheme.RunReader(reader)

			Eventually(scheme.ByteCount).Should(BeEquivalentTo(1010))
			Eventually(scheme.ErrorCount).Should(BeEquivalentTo(1))
			Expect(scheme.RunTime()).To(BeNumerically(">", time.Duration(0)))
			Expect(scheme.FirstError()).To(Equal(firstError))
		})

		It("can read upto twice the size message as it is expected to read", func() {
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 2000), Error: nil}

			scheme.RunReader(reader)

			Eventually(scheme.ByteCount).Should(BeEquivalentTo(2000))
			Eventually(scheme.ErrorCount).Should(BeEquivalentTo(0))
		})
	})

	Context("with a longer wait time", func() {
		JustBeforeEach(func() {
			json.SetString(simple.WaitForLastMessage, "1s", config)
			err := scheme.Init(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("waits the proper amount time for the last message", func() {
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 100), Error: nil}

			done := make(chan struct{})
			go func() {
				defer close(done)
				scheme.RunReader(reader)
			}()

			Consistently(done, 800*time.Millisecond).ShouldNot(BeClosed())

			Eventually(done).Should(BeClosed())
		})
	})
})

package simple_test

import (
	"netspel/factory"
	"netspel/schemes/internal/mocks"
	"netspel/schemes/simple"

	"errors"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scheme", func() {
	var (
		writer *mocks.MockWriter
		reader *mocks.MockReader
		scheme *simple.Scheme
		config *factory.Config
	)

	BeforeEach(func() {
		writer = mocks.NewMockWriter()
		reader = mocks.NewMockReader()
		scheme = &simple.Scheme{}
		config = factory.NewConfig()
		err := config.ParseAndSetAdditionalInt(simple.MessagesPerRun + "=100")
		Expect(err).NotTo(HaveOccurred())
		err = config.ParseAndSetAdditionalInt(simple.BytesPerMessage + "=1000")
		Expect(err).NotTo(HaveOccurred())
		err = config.ParseAndSetAdditionalString(simple.WaitForLastMessage + "=100ms")
		Expect(err).NotTo(HaveOccurred())
		err = config.ParseAndSetAdditionalString(simple.DefaultReport + "=")
		Expect(err).NotTo(HaveOccurred())
		err = config.ParseAndSetAdditionalString(simple.LessThanReport + "=")
		Expect(err).NotTo(HaveOccurred())
		err = config.ParseAndSetAdditionalString(simple.ErrorReport + "=")
		Expect(err).NotTo(HaveOccurred())
		err = config.ParseAndSetAdditionalString(simple.GreaterThanReport + "=")
		Expect(err).NotTo(HaveOccurred())
	})

	Context("with a short wait time", func() {
		JustBeforeEach(func() {
			err := config.ParseAndSetAdditionalString(simple.WaitForLastMessage + "=100ms")
			Expect(err).NotTo(HaveOccurred())
			err = scheme.Init(*config)
			Expect(err).NotTo(HaveOccurred())
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
			Expect(scheme.Runtime()).To(BeNumerically(">", 0))
		})

		It("reads messages from a reader", func() {
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 10), Error: nil}
			reader.ReadMessages <- mocks.ReadMessage{Buffer: make([]byte, 1000), Error: nil}
			reader.ReadMessages <- mocks.ReadMessage{Buffer: []byte{}, Error: errors.New("Bad stuff")}

			scheme.RunReader(reader)

			Eventually(scheme.ByteCount).Should(BeEquivalentTo(1010))
			Eventually(scheme.ErrorCount).Should(BeEquivalentTo(1))
			Expect(scheme.Runtime()).To(BeNumerically(">", time.Duration(0)))
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
			err := config.ParseAndSetAdditionalString(simple.WaitForLastMessage + "=1s")
			Expect(err).NotTo(HaveOccurred())
			err = scheme.Init(*config)
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

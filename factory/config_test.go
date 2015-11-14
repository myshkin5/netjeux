package factory_test

import (
	"github.com/myshkin5/netspel/factory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	It("loads config from a file", func() {
		config, err := factory.LoadFromFile("./simple.json")

		Expect(err).NotTo(HaveOccurred())
		Expect(config.WriterType).To(Equal("netspel.adapters.udp.Writer"))
		Expect(config.ReaderType).To(Equal("netspel.adapters.udp.Reader"))
		Expect(config.SchemeType).To(Equal("netspel.schemes.simple.Scheme"))
	})

	It("returns an error when attempting to load a non-existent file", func() {
		_, err := factory.LoadFromFile("./not-there.json")

		Expect(err).To(HaveOccurred())
	})

	It("parses a JSON object and stores the results", func() {
		config, err := factory.Parse([]byte(`{
			"writer-type": "SomeNeatWriter",
			"reader-type": "SomeNeatReader",
			"scheme-type": "SomeNeatScheme",
			"additional": {
				"this": "that"
			}
		}`))

		Expect(err).NotTo(HaveOccurred())
		Expect(config.WriterType).To(Equal("SomeNeatWriter"))
		Expect(config.ReaderType).To(Equal("SomeNeatReader"))
		Expect(config.SchemeType).To(Equal("SomeNeatScheme"))

		this, ok := config.Additional["this"]
		Expect(ok).To(BeTrue())
		value := this.(string)
		Expect(value).To(Equal("that"))
	})

	It("returns an error when the writer type isn't specified", func() {
		_, err := factory.Parse([]byte(`{
			"scheme-type": "SomeNeatScheme",
			"reader-type": "SomeNeatReader"
		}`))

		Expect(err).To(HaveOccurred())
	})

	It("returns an error when the reader type isn't specified", func() {
		_, err := factory.Parse([]byte(`{
			"scheme-type": "SomeNeatScheme",
			"writer-type": "SomeNeatWriter"
		}`))

		Expect(err).To(HaveOccurred())
	})

	It("returns an error when the scheme type isn't specified", func() {
		_, err := factory.Parse([]byte(`{
			"writer-type": "SomeNeatWriter",
			"reader-type": "SomeNeatReader"
		}`))

		Expect(err).To(HaveOccurred())
	})
})

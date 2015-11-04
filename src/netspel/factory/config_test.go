package factory_test

import (
	"netspel/factory"

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

	Describe("AdditionalString()", func() {
		var (
			config *factory.Config
		)

		BeforeEach(func() {
			config = factory.NewConfig()
		})

		It("returns not ok when requesting a non-existent additional string", func() {
			_, ok := config.AdditionalString("not there")
			Expect(ok).To(BeFalse())
		})

		It("returns not ok when requesting an additional string value that can't be coerced into a string", func() {
			config.Additional["something"] = 1.2
			_, ok := config.AdditionalString("something")
			Expect(ok).To(BeFalse())
		})

		It("returns ok and the value when requesting an additional string value", func() {
			config.Additional["something"] = "that's really there"
			value, ok := config.AdditionalString("something")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("that's really there"))
		})

		It("returns a child value", func() {
			config, err := factory.Parse([]byte(`{
				"writer-type": "SomeNeatWriter",
				"reader-type": "SomeNeatReader",
				"scheme-type": "SomeNeatScheme",
				"additional": {
					"parent": {
						"child": "value"
					}
				}
			}`))

			Expect(err).NotTo(HaveOccurred())

			value, ok := config.AdditionalString("parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("value"))
		})
	})

	Describe("AdditionalInt()", func() {
		var (
			config *factory.Config
		)

		BeforeEach(func() {
			config = factory.NewConfig()
		})

		It("returns not ok when requesting a non-existent additional int", func() {
			_, ok := config.AdditionalInt("not there")
			Expect(ok).To(BeFalse())
		})

		It("returns not ok when requesting an additional int value that can't be coerced into a int", func() {
			config.Additional["something"] = "not an int"
			_, ok := config.AdditionalInt("something")
			Expect(ok).To(BeFalse())
		})

		It("returns ok and the value when requesting an additional string value", func() {
			config.Additional["something"] = 1234
			value, ok := config.AdditionalInt("something")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(1234))
		})

		It("returns a child value", func() {
			config, err := factory.Parse([]byte(`{
				"writer-type": "SomeNeatWriter",
				"reader-type": "SomeNeatReader",
				"scheme-type": "SomeNeatScheme",
				"additional": {
					"parent": {
						"child": 98765
					}
				}
			}`))

			Expect(err).NotTo(HaveOccurred())

			value, ok := config.AdditionalInt("parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(98765))
		})
	})

	Describe("ParseAndSetAdditionalString()", func() {
		It("parses and overrides CLI values", func() {
			config, err := factory.Parse([]byte(`{
				"writer-type": "SomeNeatWriter",
				"reader-type": "SomeNeatReader",
				"scheme-type": "SomeNeatScheme",
				"additional": {
					"this": "that",
					"parent": {
						"child": "value"
					}
				}
			}`))

			Expect(err).NotTo(HaveOccurred())

			err = config.ParseAndSetAdditionalString("this=something else")
			Expect(err).NotTo(HaveOccurred())

			value, ok := config.AdditionalString("this")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("something else"))

			err = config.ParseAndSetAdditionalString("parent.child=new value")
			Expect(err).NotTo(HaveOccurred())

			value, ok = config.AdditionalString("parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("new value"))
		})

		It("returns an error if there isn't exactly one equals sign", func() {
			config := factory.NewConfig()
			err := config.ParseAndSetAdditionalString("neat")
			Expect(err).To(HaveOccurred())
			err = config.ParseAndSetAdditionalString("this_is=too=neat")
			Expect(err).To(HaveOccurred())
		})

		It("can add string values even when the config was parsed but there was no additional section", func() {
			config, err := factory.Parse([]byte(`{
				"writer-type": "SomeNeatWriter",
				"reader-type": "SomeNeatReader",
				"scheme-type": "SomeNeatScheme"
			}`))

			Expect(err).NotTo(HaveOccurred())

			err = config.ParseAndSetAdditionalString("this=something else")
			Expect(err).NotTo(HaveOccurred())
		})

		It("can set a string to empty string", func() {
			config := factory.NewConfig()
			err := config.ParseAndSetAdditionalString("value=something")
			Expect(err).NotTo(HaveOccurred())
			err = config.ParseAndSetAdditionalString("value=")
			Expect(err).NotTo(HaveOccurred())

			value, ok := config.AdditionalString("value")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(""))

			err = config.ParseAndSetAdditionalString("another=")
			Expect(err).NotTo(HaveOccurred())

			another, ok := config.AdditionalString("another")
			Expect(ok).To(BeTrue())
			Expect(another).To(Equal(""))
		})

		It("can set values multiple levels deep", func() {
			config := factory.NewConfig()
			err := config.ParseAndSetAdditionalString("one.two.three=hi")
			Expect(err).NotTo(HaveOccurred())
			value, ok := config.AdditionalString("one.two.three")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("hi"))
		})
	})

	Describe("ParseAndSetAdditionalInt()", func() {
		It("parses and overrides CLI values", func() {
			config, err := factory.Parse([]byte(`{
				"writer-type": "SomeNeatWriter",
				"reader-type": "SomeNeatReader",
				"scheme-type": "SomeNeatScheme",
				"additional": {
					"this": "that",
					"parent": {
						"child": 98765
					}
				}
			}`))

			Expect(err).NotTo(HaveOccurred())

			err = config.ParseAndSetAdditionalInt("this=1000000")
			Expect(err).NotTo(HaveOccurred())

			value, ok := config.AdditionalInt("this")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(1000000))

			err = config.ParseAndSetAdditionalInt("parent.child=12345")
			Expect(err).NotTo(HaveOccurred())

			value, ok = config.AdditionalInt("parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(12345))
		})

		It("returns an error if there isn't exactly one equals sign or if the value isn't an int", func() {
			config := factory.NewConfig()
			err := config.ParseAndSetAdditionalInt("neat")
			Expect(err).To(HaveOccurred())
			err = config.ParseAndSetAdditionalInt("this_is=too=neat")
			Expect(err).To(HaveOccurred())
			err = config.ParseAndSetAdditionalInt("this_is=not_an_int")
			Expect(err).To(HaveOccurred())
		})
	})
})

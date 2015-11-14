package jsonstruct_test

import (
	"encoding/json"
	"netspel/jsonstruct"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON", func() {
	var (
		values jsonstruct.JSONStruct
	)

	BeforeEach(func() {
		values = jsonstruct.New()
	})

	Describe("String()", func() {
		It("returns not ok when requesting a non-existent additional string", func() {
			_, ok := values.String("not there")
			Expect(ok).To(BeFalse())
		})

		It("returns not ok when requesting an additional string value that can't be coerced into a string", func() {
			values["something"] = 1.2
			_, ok := values.String("something")
			Expect(ok).To(BeFalse())
		})

		It("returns ok and the value when requesting an additional string value", func() {
			values["something"] = "that's really there"
			value, ok := values.String("something")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("that's really there"))
		})

		It("returns a child value", func() {
			err := json.Unmarshal([]byte(`{
				"parent": {
					"child": "value"
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			value, ok := values.String("parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("value"))
		})
	})

	Describe("StringWithDefault()", func() {
		It("returns the default value when a value isn't found", func() {
			Expect(values.StringWithDefault("not-present-path", "default-value")).To(Equal("default-value"))
		})

		It("returns the non-default value when a value is found", func() {
			values["present-path"] = "non-default-value"

			Expect(values.StringWithDefault("present-path", "default-value")).To(Equal("non-default-value"))
		})
	})

	Describe("Int()", func() {
		It("returns not ok when requesting a non-existent additional int", func() {
			_, ok := values.Int("not there")
			Expect(ok).To(BeFalse())
		})

		It("returns not ok when requesting an additional int value that can't be coerced into a int", func() {
			values["something"] = "not an int"
			_, ok := values.Int("something")
			Expect(ok).To(BeFalse())
		})

		It("returns ok and the value when requesting an additional string value", func() {
			values["something"] = 1234
			value, ok := values.Int("something")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(1234))
		})

		It("returns a child value", func() {
			err := json.Unmarshal([]byte(`{
				"parent": {
					"child": 98765
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			value, ok := values.Int("parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(98765))
		})
	})

	Describe("IntWithDefault()", func() {
		It("returns the default value when a value isn't found", func() {
			Expect(values.IntWithDefault("not-present-path", 42)).To(Equal(42))
		})

		It("returns the non-default value when a value is found", func() {
			values["present-path"] = 84

			Expect(values.IntWithDefault("present-path", 42)).To(Equal(84))
		})
	})

	Describe("Duration()", func() {
		It("returns a not found error when the value doesn't exist", func() {
			_, err := values.Duration("not there")
			Expect(err).To(Equal(jsonstruct.ErrValueNotFound))
		})

		It("returns an error when there is an error parsing the duration", func() {
			values["not-a-duration-path"] = "not-a-duration"
			_, err := values.Duration("not-a-duration-path")
			Expect(err).To(HaveOccurred())
		})

		It("returns valid durations", func() {
			values["valid-duration"] = "20s"
			duration, err := values.Duration("valid-duration")
			Expect(err).NotTo(HaveOccurred())
			Expect(duration).To(Equal(20 * time.Second))
		})
	})

	Describe("DurationWithDefault()", func() {
		It("returns the default value when a value isn't found", func() {
			Expect(values.DurationWithDefault("not-present-path", 15*time.Millisecond)).To(Equal(15 * time.Millisecond))
		})

		It("returns the non-default value when a value is found", func() {
			values["present-path"] = "84ms"

			Expect(values.DurationWithDefault("present-path", 42*time.Millisecond)).To(Equal(84 * time.Millisecond))
		})
	})
})

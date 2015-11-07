package json_test

import (
	encoding_json "encoding/json"
	"netspel/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON", func() {
	var (
		values map[string]interface{}
	)

	Describe("String()", func() {
		BeforeEach(func() {
			values = make(map[string]interface{})
		})

		It("returns not ok when requesting a non-existent additional string", func() {
			_, ok := json.String("not there", values)
			Expect(ok).To(BeFalse())
		})

		It("returns not ok when requesting an additional string value that can't be coerced into a string", func() {
			values["something"] = 1.2
			_, ok := json.String("something", values)
			Expect(ok).To(BeFalse())
		})

		It("returns ok and the value when requesting an additional string value", func() {
			values["something"] = "that's really there"
			value, ok := json.String("something", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("that's really there"))
		})

		It("returns a child value", func() {
			err := encoding_json.Unmarshal([]byte(`{
				"parent": {
					"child": "value"
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			value, ok := json.String("parent.child", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("value"))
		})
	})

	Describe("Int()", func() {
		BeforeEach(func() {
			values = make(map[string]interface{})
		})

		It("returns not ok when requesting a non-existent additional int", func() {
			_, ok := json.Int("not there", values)
			Expect(ok).To(BeFalse())
		})

		It("returns not ok when requesting an additional int value that can't be coerced into a int", func() {
			values["something"] = "not an int"
			_, ok := json.Int("something", values)
			Expect(ok).To(BeFalse())
		})

		It("returns ok and the value when requesting an additional string value", func() {
			values["something"] = 1234
			value, ok := json.Int("something", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(1234))
		})

		It("returns a child value", func() {
			err := encoding_json.Unmarshal([]byte(`{
				"parent": {
					"child": 98765
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			value, ok := json.Int("parent.child", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(98765))
		})
	})
})

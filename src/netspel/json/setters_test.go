package json_test

import (
	encoding_json "encoding/json"
	"netspel/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Setters", func() {
	var (
		values map[string]interface{}
	)

	Describe("SetString()", func() {
		It("overrides existing values", func() {
			err := encoding_json.Unmarshal([]byte(`{
				"this": "that",
				"parent": {
					"child": "value"
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			json.SetString("this", "something else", values)

			value, ok := json.String("this", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("something else"))

			json.SetString("parent.child", "new value", values)

			value, ok = json.String("parent.child", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("new value"))
		})

		It("can set a string to empty string", func() {
			values = make(map[string]interface{})
			json.SetString("value", "something", values)
			json.SetString("value", "", values)

			value, ok := json.String("value", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(""))

			json.SetString("another", "", values)

			another, ok := json.String("another", values)
			Expect(ok).To(BeTrue())
			Expect(another).To(Equal(""))
		})

		It("can set values multiple levels deep", func() {
			values = make(map[string]interface{})
			json.SetString("one.two.three", "hi", values)
			value, ok := json.String("one.two.three", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("hi"))
		})
	})

	Describe("SetInt()", func() {
		It("overrides existing values", func() {
			err := encoding_json.Unmarshal([]byte(`{
				"this": "that",
				"parent": {
					"child": 98765
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			json.SetInt("this", 1000000, values)

			value, ok := json.Int("this", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(1000000))

			json.SetInt("parent.child", 12345, values)

			value, ok = json.Int("parent.child", values)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(12345))
		})
	})
})

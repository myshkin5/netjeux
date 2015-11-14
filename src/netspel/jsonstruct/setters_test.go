package jsonstruct_test

import (
	"encoding/json"
	"netspel/jsonstruct"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Setters", func() {
	var (
		values jsonstruct.JSONStruct
	)

	Describe("SetString()", func() {
		It("overrides existing values", func() {
			err := json.Unmarshal([]byte(`{
				"this": "that",
				"parent": {
					"child": "value"
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			values.SetString("this", "something else")

			value, ok := values.String("this")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("something else"))

			values.SetString("parent.child", "new value")

			value, ok = values.String("parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("new value"))
		})

		It("can set a string to empty string", func() {
			values = jsonstruct.New()
			values.SetString("value", "something")
			values.SetString("value", "")

			value, ok := values.String("value")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(""))

			values.SetString("another", "")

			another, ok := values.String("another")
			Expect(ok).To(BeTrue())
			Expect(another).To(Equal(""))
		})

		It("can set values multiple levels deep", func() {
			values = jsonstruct.New()
			values.SetString("one.two.three", "hi")
			value, ok := values.String("one.two.three")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("hi"))
		})
	})

	Describe("SetInt()", func() {
		It("overrides existing values", func() {
			err := json.Unmarshal([]byte(`{
				"this": "that",
				"parent": {
					"child": 98765
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			values.SetInt("this", 1000000)

			value, ok := values.Int("this")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(1000000))

			values.SetInt("parent.child", 12345)

			value, ok = values.Int("parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(12345))
		})
	})

	Describe("SetDuration()", func() {
		It("sets a duration value as a string", func() {
			values = jsonstruct.New()

			values.SetDuration("duration-path", 32*time.Second)

			Expect(values["duration-path"]).To(Equal("32s"))
		})
	})
})

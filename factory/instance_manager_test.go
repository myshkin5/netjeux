package factory_test

import (
	"reflect"

	"github.com/myshkin5/netspel/factory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type coolInterface interface {
	coolMethod() int
}

type coolType struct {
	coolValue int
}

func (t *coolType) coolMethod() int {
	return t.coolValue
}

var _ = Describe("Manager", func() {
	var (
		manager *factory.InstanceManager
	)

	BeforeEach(func() {
		manager = factory.NewInstanceManager()
	})

	It("returns an error when attempting to create an unregistered type", func() {
		_, err := manager.CreateInstance("don't matter")
		Expect(err).To(HaveOccurred())
	})

	It("returns an new instance of a registered type", func() {
		manager.RegisterType("cool type", reflect.TypeOf(coolType{}))
		coolInstanceValue, err := manager.CreateInstance("cool type")
		Expect(err).NotTo(HaveOccurred())
		coolInstance, ok := coolInstanceValue.Interface().(coolInterface)
		Expect(ok).To(BeTrue())
		Expect(coolInstance.coolMethod()).To(BeZero())
	})
})

package schemes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSchemes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Schemes Suite")
}

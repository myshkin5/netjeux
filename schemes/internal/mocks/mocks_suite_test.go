package mocks_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSchemes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Schemes - Internal - Mocks Suite")
}

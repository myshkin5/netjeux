package sse_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adapters - SSE Suite")
}

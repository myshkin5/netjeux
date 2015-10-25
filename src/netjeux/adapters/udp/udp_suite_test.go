package udp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUDP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UDP Suite")
}

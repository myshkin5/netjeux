package simple_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/myshkin5/netspel/logs"
	"github.com/op/go-logging"
)

func TestSimple(t *testing.T) {
	RegisterFailHandler(Fail)
	logs.LogLevel.SetLevel(logging.CRITICAL, "netspel")
	RunSpecs(t, "Schemes - Simple Suite")
}

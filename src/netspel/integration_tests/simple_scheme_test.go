package integration_tests_test

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("simple.Scheme", func() {
	var (
		readerSession *gexec.Session
		writerSession *gexec.Session
	)

	AfterEach(func() {
		if readerSession != nil {
			readerSession.Terminate()
		}
		if writerSession != nil {
			writerSession.Terminate()
		}

		gexec.CleanupBuildArtifacts()
	})

	It("compiles, sends and receives messages via UDP", func() {
		executablePath, err := gexec.Build("netspel")
		Expect(err).NotTo(HaveOccurred())

		readerCommand := exec.Command(executablePath, "--config", "./simple.json", "read")
		readerSession, err = gexec.Start(readerCommand,
			gexec.NewPrefixedWriter("\x1b[37m[o]\x1b[32m[reader]\x1b[0m ", GinkgoWriter),
			gexec.NewPrefixedWriter("\x1b[31m[e]\x1b[32m[reader]\x1b[0m ", GinkgoWriter))
		Expect(err).NotTo(HaveOccurred())
		writerCommand := exec.Command(executablePath, "--config", "./simple.json", "write")
		writerSession, err = gexec.Start(writerCommand,
			gexec.NewPrefixedWriter("\x1b[37m[o]\x1b[31m[writer]\x1b[0m ", GinkgoWriter),
			gexec.NewPrefixedWriter("\x1b[31m[e]\x1b[31m[writer]\x1b[0m ", GinkgoWriter))
		Expect(err).NotTo(HaveOccurred())

		Eventually(readerSession, 10*time.Second).Should(gexec.Exit(0))
		Eventually(writerSession, 10*time.Second).Should(gexec.Exit(0))
	})
})

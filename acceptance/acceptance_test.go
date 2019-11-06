package acceptance_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Acceptance", func() {
	It("should create the process in a new uts namespace", func() {
		cmd := exec.Command(barrelPath, "roll", "--rootfs", rootfsPath, "/bin/sh", "--", "-c", "hostname foo; hostname")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("foo"))
	})

	It("should exit with the exit code of the container", func() {
		cmd := exec.Command(barrelPath, "roll", "-r", rootfsPath, "/bin/sh", "--", "-c", "exit 12")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.ExitCode).Should(Equal(12))
	})

	It("should fail when the rootfs is not set", func() {
		cmd := exec.Command(barrelPath, "roll", "/bin/sh", "--", "-c", "echo", "this will not run")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.ExitCode).Should(Equal(1))
	})
})

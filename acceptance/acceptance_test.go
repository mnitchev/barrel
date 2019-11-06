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
		cmdPath, err := gexec.Build("github.com/mnitchev/barrel")
		Expect(err).ToNot(HaveOccurred())

		cmd := exec.Command(cmdPath, "roll", "bash", "--", "-c", "hostname foo; hostname")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("foo"))
	})

	It("should exit with the exit code of the container", func() {
		cmdPath, err := gexec.Build("github.com/mnitchev/barrel")
		Expect(err).ToNot(HaveOccurred())

		cmd := exec.Command(cmdPath, "roll", "bash", "--", "-c", "exit 12")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.ExitCode).Should(Equal(12))
	})
})

package acceptance_test

import (
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Acceptance", func() {

	It("should create the process in a new uts namespace", func() {
		cmd := exec.Command(barrelPath, "roll", "--rootfs", rootfsPath, "--cgroup-name", "test", "/bin/sh", "--", "-c", "hostname foo; hostname")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.ExitCode).Should(Equal(0))
		Expect(session.Out).To(gbytes.Say("foo"))
	})

	It("should exit with the exit code of the container", func() {
		cmd := exec.Command(barrelPath, "roll", "-r", rootfsPath, "--cgroup-name", "test", "/bin/sh", "--", "-c", "exit 12")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.ExitCode).Should(Equal(12))
	})

	It("should fail when the rootfs is not set", func() {
		cmd := exec.Command(barrelPath, "roll", "/bin/sh", "--cgroup-name", "test", "--", "-c", "echo", "this will not run")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.ExitCode).Should(Equal(1))
	})

	It("should fail when the cgroup name is not set", func() {
		cmd := exec.Command(barrelPath, "roll", "/bin/sh", "--cgroup-name", "test", "--", "-c", "echo", "this will not run")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session.ExitCode).Should(Equal(1))
	})

	It("should not be able to list the host's processes", func() {
		cmd := exec.Command(barrelPath, "roll", "-r", rootfsPath, "--cgroup-name", "test", "/bin/sh", "--", "-c", "ps", "aux")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session.ExitCode).Should(Equal(0))
		output := strings.Trim(string(session.Out.Contents()), "\n")
		lines := strings.Split(output, "\n")
		Expect(lines).To(HaveLen(3))
		Expect(lines).To(ContainElement(MatchRegexp("1.*/bin/sh -c ps")))
		Expect(lines).To(ContainElement(MatchRegexp("[0-9]+.*ps")))
	})
})

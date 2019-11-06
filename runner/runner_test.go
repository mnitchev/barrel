package runner_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/mnitchev/barrel/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	When("running a command", func() {
		verifyNamespaceIsCreated := func(ns string) {
			output := bytes.Buffer{}
			procPath := fmt.Sprintf("/proc/self/ns/%s", ns)
			container := runner.Container{
				Command: "readlink",
				Args:    []string{"-n", procPath},
				Stdin:   os.Stdin,
				Stdout:  &output,
				Stderr:  os.Stderr,
			}

			parentNs, err := os.Readlink(procPath)
			Expect(err).NotTo(HaveOccurred())

			exitCode, err := runner.Run(container)
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(Equal(0))

			containerNs := output.String()
			symlincRegex := fmt.Sprintf("%s:[[0-9]+]", ns)
			Expect(containerNs).To(MatchRegexp(symlincRegex))
			Expect(containerNs).NotTo(Equal(parentNs))

		}
		It("should run it in a new UTS namespace", func() {
			verifyNamespaceIsCreated("uts")
		})

		It("should run it in a new mount namespace", func() {
			output := bytes.Buffer{}
			container := runner.Container{
				Command: "readlink",
				Args:    []string{"-n", "/proc/self/ns/mnt"},
				Stdin:   os.Stdin,
				Stdout:  &output,
				Stderr:  os.Stderr,
			}

			parentMnt, err := os.Readlink("/proc/self/ns/mnt")
			Expect(err).NotTo(HaveOccurred())

			exitCode, err := runner.Run(container)
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(Equal(0))

			containerMnt := output.String()
			Expect(containerMnt).To(MatchRegexp("mnt:[[0-9]+]"))
			Expect(containerMnt).NotTo(Equal(parentMnt))
		})

		It("should set the PS1 env variable", func() {
			output := bytes.Buffer{}
			container := runner.Container{
				Command: "env",
				Args:    []string{},
				Stdin:   os.Stdin,
				Stdout:  &output,
				Stderr:  os.Stderr,
			}

			exitCode, err := runner.Run(container)
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(Equal(0))
			Expect(output.String()).To(ContainSubstring("PS1=λ [contained-process] → "))
		})
	})

	When("the command exits with a non-zero exit code", func() {
		It("should return the exit code", func() {
			errOutput := bytes.Buffer{}
			container := runner.Container{
				Command: "/bin/sh",
				Args:    []string{"-c", "nonexistent-command"},
				Stdin:   os.Stdin,
				Stdout:  os.Stdout,
				Stderr:  &errOutput,
			}

			exitCode, err := runner.Run(container)
			Expect(err).To(HaveOccurred())
			Expect(exitCode).To(Equal(127))
			Expect(errOutput.String()).To(ContainSubstring("nonexistent-command: not found"))
		})
	})

	When("the command cannot be started", func() {
		It("should exit with exit code 1", func() {
			container := runner.Container{
				Command: "non-existent-command",
				Args:    []string{"-c", "echo", "hello"},
				Stdin:   os.Stdin,
				Stdout:  os.Stdout,
				Stderr:  os.Stderr,
			}

			exitCode, err := runner.Run(container)
			Expect(err).To(HaveOccurred())
			Expect(exitCode).To(Equal(1))
		})
	})
})

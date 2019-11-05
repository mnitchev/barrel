package runner_test

import (
	"bytes"
	"os"

	"github.com/mnitchev/barrel/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	When("running a command", func() {
		It("should run it in a new UTS namespace", func() {
			output := bytes.Buffer{}
			container := runner.Container{
				Command: "readlink",
				Args:    []string{"-n", "/proc/self/ns/uts"},
				Stdin:   os.Stdin,
				Stdout:  &output,
				Stderr:  os.Stderr,
			}

			parentUts, err := os.Readlink("/proc/self/ns/uts")
			Expect(err).NotTo(HaveOccurred())
			containerUts := output.String()

			err = runner.Run(container)
			Expect(err).NotTo(HaveOccurred())
			Expect(containerUts).NotTo(Equal(parentUts))
		})
	})
})

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
			buf := bytes.Buffer{}
			container := runner.Container{
				Command: "readlink",
				Args:    []string{"-n", "/proc/self/ns/uts"},
				Stdin:   os.Stdin,
				Stdout:  &buf,
				Stderr:  os.Stderr,
			}

			parentUts, err := os.Readlink("/proc/self/ns/uts")
			Expect(err).NotTo(HaveOccurred())

			err = runner.Run(container)
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).NotTo(Equal(parentUts))
		})

	})
})

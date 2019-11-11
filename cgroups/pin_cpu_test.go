package cgroups_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mnitchev/barrel/cgroups"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PinCpu", func() {
	const cgroupName = "test"

	BeforeEach(func() {
		cgroupPath := filepath.Join("/sys/fs/cgroup/cpuset/", cgroupName)
		Expect(os.MkdirAll(cgroupPath, 0755)).To(Succeed())
	})

	AfterEach(func() {
		cgroupPath := filepath.Join("/sys/fs/cgroup/cpuset/", cgroupName)
		Expect(os.RemoveAll(cgroupPath)).To(Succeed())
	})

	When("pinning a cpu to a cgrpup", func() {
		It("should write the cpu indexes in the cpuset.cpus file of the cgroup", func() {
			Expect(cgroups.PinCPU(cgroupName, "0-1,3")).To(Succeed())
			cpusFilePath := filepath.Join("/sys/fs/cgroup/cpuset/", cgroupName, "cpuset.cpus")
			cpus, err := ioutil.ReadFile(cpusFilePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(cpus)).To(Equal("0-1,3\n"))
		})
	})
})

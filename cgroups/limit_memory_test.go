package cgroups_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mnitchev/barrel/cgroups"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LimitMemory", func() {
	const cgroupName = "test"

	BeforeEach(func() {
		cgroupPath := filepath.Join("/sys/fs/cgroup/memory/", cgroupName)
		Expect(os.MkdirAll(cgroupPath, 0755)).To(Succeed())
	})

	AfterEach(func() {
		cgroupPath := filepath.Join("/sys/fs/cgroup/memory/", cgroupName)
		Expect(os.RemoveAll(cgroupPath)).To(Succeed())
	})

	When("limiting the memory of a cgroup", func() {
		It("should write the memory limit in the memory.limit_in_bytes file", func() {
			Expect(cgroups.LimitMemory(cgroupName, "1024M")).To(Succeed())
			memoryLimitFile := filepath.Join("/sys/fs/cgroup/memory/", cgroupName, "memory.limit_in_bytes")
			limit, err := ioutil.ReadFile(memoryLimitFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(limit)).To(Equal("1073741824\n"))
		})
	})

	It("should write the memory limit + 2MB in the memory.memsw.limit_in_bytes file", func() {
		Expect(cgroups.LimitMemory(cgroupName, "1024M")).To(Succeed())
		memoryLimitFile := filepath.Join("/sys/fs/cgroup/memory/", cgroupName, "memory.memsw.limit_in_bytes")
		limit, err := ioutil.ReadFile(memoryLimitFile)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(limit)).To(Equal("1075838976\n"))
	})

})

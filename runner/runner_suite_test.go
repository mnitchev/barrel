package runner_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var rootfsPath string

func TestRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Runner Suite")
}

var _ = BeforeSuite(func() {
	extractRootfs()
})

var _ = AfterSuite(func() {
	Expect(os.RemoveAll(rootfsPath)).To(Succeed())
})

func extractRootfs() {
	var err error
	rootfsPath, err = ioutil.TempDir("", "rootfs")
	Expect(os.Chmod(rootfsPath, 0777)).To(Succeed())
	Expect(err).ToNot(HaveOccurred())
	untarCmd := exec.Command("tar", "xfz", "../acceptance/assets/busybox/busybox1.25.0-uclibc.tar.gz", "--directory", rootfsPath)
	Expect(untarCmd.Run()).To(Succeed())
}

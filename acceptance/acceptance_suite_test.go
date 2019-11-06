package acceptance_test

import (
	"io/ioutil"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	rootfsPath string
	barrelPath string
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var _ = BeforeSuite(func() {
	extractRootfs()
	buildBarrel()
})

func buildBarrel() {
	var err error
	barrelPath, err = gexec.Build("github.com/mnitchev/barrel")
	Expect(err).ToNot(HaveOccurred())
}

func extractRootfs() {
	var err error
	rootfsPath, err = ioutil.TempDir("", "rootfs")
	Expect(err).ToNot(HaveOccurred())
	untarCmd := exec.Command("tar", "xfz", "assets/busybox/busybox1.25.0-uclibc.tar.gz", "--directory", rootfsPath)
	Expect(untarCmd.Run()).To(Succeed())
}

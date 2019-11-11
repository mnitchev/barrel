package cgroups_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCgroups(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cgroups Suite")
}

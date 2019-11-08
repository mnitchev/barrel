package cgroups

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func PinCPU(cgroupName, cpuIndexes string) error {
	if err := ioutil.WriteFile(filepath.Join("/sys/fs/cgroup/cpuset/", cgroupName, "cpuset.cpus"), []byte(cpuIndexes), 0644); err != nil {
		fmt.Printf("Failed to pin cpu %s: %s", cpuIndexes, err)
		return err
	}
	return nil
}

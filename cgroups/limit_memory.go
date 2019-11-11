package cgroups

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
)

func LimitMemory(cgroupName string, maxMemory string) error {
	swapMemory, err := calculateSwapMemory(maxMemory)
	if err != nil {
		fmt.Printf("Failed to parse memory %s: %s\n", maxMemory, err)
		return err
	}

	memoryLimitFile := filepath.Join("/sys/fs/cgroup/memory/", cgroupName, "memory.limit_in_bytes")
	if err := ioutil.WriteFile(memoryLimitFile, []byte(maxMemory), 0644); err != nil {
		fmt.Printf("Failed to set memory limit %s: %s\n", maxMemory, err)
		return err
	}

	swapMemoryLimitFile := filepath.Join("/sys/fs/cgroup/memory/", cgroupName, "memory.memsw.limit_in_bytes")
	if err := ioutil.WriteFile(swapMemoryLimitFile, []byte(swapMemory), 0644); err != nil {
		fmt.Printf("Failed to set swap memory limit  %s: %s\n", maxMemory, err)
		return err
	}
	return nil
}

func calculateSwapMemory(memory string) (string, error) {
	memoryFormat := regexp.MustCompile(`[0-9]+(M|G|K|m|g|k)?`)
	if !memoryFormat.MatchString(memory) {
		return "", errors.New("memory limit not in correct format")
	}
	amountRegex := regexp.MustCompile(`[0-9]+`)
	amount, err := strconv.Atoi(amountRegex.FindString(memory))
	if err != nil {
		return "", err
	}

	unitRegex := regexp.MustCompile(`(M|G|K|m|g|k)`)
	unit := unitRegex.FindString(memory)

	return strconv.Itoa(amount+2) + unit, nil
}

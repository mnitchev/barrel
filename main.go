package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	command := os.Args[1]
	arguments := os.Args[2:]
	cmd := exec.Command(command, arguments...)
	cmd.Env = []string{"PS1=-[contained-process]- #", "FOO=bar"}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to run the command %s \n", command)
		panic(err)
	}
}

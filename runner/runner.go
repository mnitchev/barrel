package runner

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

const promt = "λ [contained-process] → "

type Container struct {
	Command string
	Args    []string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
}

func Run(container Container) error {
	cmd := exec.Command(container.Command, container.Args...)
	promtEnv := fmt.Sprintf("PS1=%s", promt)
	cmd.Env = []string{promtEnv}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Stdin = container.Stdin
	cmd.Stdout = container.Stdout
	cmd.Stderr = container.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to run the command %s with args: %s \n", container.Command, container.Args)
		return err
	}
	return nil
}

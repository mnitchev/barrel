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

func Run(container Container) (int, error) {
	cmd := exec.Command(container.Command, container.Args...)
	promtEnv := fmt.Sprintf("PS1=%s", promt)
	cmd.Env = []string{promtEnv}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS,
	}
	cmd.Stdin = container.Stdin
	cmd.Stdout = container.Stdout
	cmd.Stderr = container.Stderr

	err := cmd.Run()
	return parseExitCode(err), err
}

func parseExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			fmt.Printf("Command exited with non-zero exit code: %d\n", status.ExitStatus())
			return status.ExitStatus()
		}
	}
	fmt.Printf("Command failed to start, %s", err)
	return 1
}

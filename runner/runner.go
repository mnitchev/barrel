package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

const promt = "λ [contained-process] → "

type Container struct {
	Command    string
	Args       []string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	RootfsPath string
}

func init() {
	reexec.Register("installNamespaces", installNamespaces)
	if reexec.Init() {
		os.Exit(0)
	}
}

func installNamespaces() {
	runContainer()
}

func runContainer() {
	command := os.Args[1]
	args := os.Args[2:]
	cmd := exec.Command(command, args...)
	promtEnv := fmt.Sprintf("PS1=%s", promt)
	cmd.Env = []string{promtEnv}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running process: %s\n", err)
		os.Exit(parseExitCode(err))
	}
}

func Run(container Container) (int, error) {
	args := append([]string{"installNamespaces", container.Command}, container.Args...)
	cmd := reexec.Command(args...)

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

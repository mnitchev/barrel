package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

const prompt = "λ [contained-process] → "

type Container struct {
	Command    string
	Args       []string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	RootfsPath string
	CgroupName string
}

func init() {
	reexec.Register("installNamespaces", installNamespaces)
	if reexec.Init() {
		os.Exit(0)
	}
}

func installNamespaces() {
	rootfsPath := os.Args[1]
	procfs := filepath.Join(rootfsPath, "/proc")
	if err := syscall.Mount("proc", procfs, "proc", 0, ""); err != nil {
		fmt.Errorf("Failed to chdir: %s", err)
		panic(err)
	}

	if err := syscall.Mount(rootfsPath, rootfsPath, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		fmt.Printf("Failed to mount new rootfs at %s: %s", rootfsPath, err)
		panic(err)
	}

	oldPath := filepath.Join(rootfsPath, ".oldroot")
	if err := os.MkdirAll(oldPath, 0700); err != nil {
		fmt.Errorf("Failed to create directory for old rootfs: %s", err)
		panic(err)
	}

	if err := syscall.PivotRoot(rootfsPath, oldPath); err != nil {
		fmt.Errorf("Failed to pivot root: %s", err)
		panic(err)
	}

	if err := os.Chdir("/"); err != nil {
		fmt.Errorf("Failed to chdir: %s", err)
		panic(err)
	}

	if err := syscall.Unmount("/.oldroot", syscall.MNT_DETACH); err != nil {
		fmt.Errorf("Failed to unmount old rootfs")
		panic(err)
	}

	if err := os.RemoveAll("/.oldroot"); err != nil {
		fmt.Errorf("Failed to remove old rootfs")
		panic(err)
	}

	runContainer()
}

func runContainer() {
	args := os.Args[2:]
	commandPath, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Printf("Command %s not found in PATH\n", args[0])
		os.Exit(127)
	}

	env := append(os.Environ(), fmt.Sprintf("PS1=%s", prompt))
	if err := syscall.Exec(commandPath, args, env); err != nil {
		fmt.Printf("Error running process: %s\n", err)
		os.Exit(parseExitCode(err))
	}
}

func Run(container Container) (int, error) {
	args := append([]string{"installNamespaces", container.RootfsPath, container.Command}, container.Args...)
	cmd := reexec.Command(args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID,
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

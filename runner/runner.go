package runner

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	cgroupName := os.Args[2]

	setUpCgroups(cgroupName)
	mountProcfs(rootfsPath)
	mountRootfs(rootfsPath)
	runContainer()
}

func runContainer() {
	args := os.Args[3:]
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
	args := append([]string{"installNamespaces", container.RootfsPath, container.CgroupName, container.Command}, container.Args...)
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

func mountProcfs(rootfsPath string) {
	procfs := filepath.Join(rootfsPath, "/proc")
	if err := syscall.Mount("proc", procfs, "proc", 0, ""); err != nil {
		fmt.Printf("Failed to chdir: %s", err)
		panic(err)
	}
}

func setUpCgroups(cgroupName string) {
	setUpCpusetCgroup(cgroupName)
	setUpMemoryCgroup(cgroupName)
}

func setUpMemoryCgroup(cgroupName string) {
	memoryCgroup := "/sys/fs/cgroup/memory/"
	cgroupPath := filepath.Join(memoryCgroup, cgroupName)
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		fmt.Printf("Failed to create new memory cgroup %s: %s", cgroupName, err)
		panic(err)
	}
	pid := os.Getpid()
	if err := ioutil.WriteFile(filepath.Join(cgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		fmt.Printf("Failed to set pid %d in cgroup tasks file: %s", pid, err)
		panic(err)
	}
}

func setUpCpusetCgroup(cgroupName string) {
	cpusetCgroup := "/sys/fs/cgroup/cpuset/"
	cgroupPath := filepath.Join(cpusetCgroup, cgroupName)
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		fmt.Printf("Failed to create new cpuset cgroup %s: %s", cgroupName, err)
		panic(err)
	}

	copyFile(filepath.Join(cpusetCgroup, "cpuset.mems"), filepath.Join(cgroupPath, "cpuset.mems"))
	copyFile(filepath.Join(cpusetCgroup, "cpuset.cpus"), filepath.Join(cgroupPath, "cpuset.cpus"))

	pid := os.Getpid()
	if err := ioutil.WriteFile(filepath.Join(cgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		fmt.Printf("Failed to set pid %d in cgroup tasks file: %s", pid, err)
		panic(err)
	}
}

func mountRootfs(rootfsPath string) {
	ns, _ := os.Readlink("/proc/self/ns/mnt")
	fmt.Printf("Mount namespace: %s\n", ns)
	upperDir := "/barrel/upper"
	if err := os.MkdirAll(upperDir, 0755); err != nil {
		fmt.Printf("Failed to create temp dir at %s: %s", upperDir, err)
		panic(err)
	}
	workDir := "/barrel/work"
	if err := os.MkdirAll(workDir, 0755); err != nil {
		fmt.Printf("Failed to create temp dir at %s: %s", workDir, err)
		panic(err)
	}
	mergedDir := "/barrel/merged"
	if err := os.MkdirAll(mergedDir, 0755); err != nil {
		fmt.Printf("Failed to create temp dir at %s: %s", mergedDir, err)
		panic(err)
	}
	data := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", rootfsPath, upperDir, workDir)
	if err := syscall.Mount("overlay", mergedDir, "overlay", 0, data); err != nil {
		fmt.Printf("Failed to mount new rootfs at %s: %s\n", mergedDir, err)
		panic(err)
	}

	oldPath := filepath.Join(mergedDir, ".oldroot")
	if err := os.MkdirAll(oldPath, 0755); err != nil {
		fmt.Printf("Failed to create directory for old rootfs: %s", err)
		panic(err)
	}

	if err := syscall.PivotRoot(mergedDir, oldPath); err != nil {
		fmt.Printf("Failed to pivot root at %s: %s\n", mergedDir, err)
		panic(err)
	}

	if err := os.Chdir("/"); err != nil {
		fmt.Printf("Failed to chdir: %s", err)
		panic(err)
	}

	if err := syscall.Unmount("/.oldroot", syscall.MNT_DETACH); err != nil {
		fmt.Printf("Failed to unmount old rootfs")
		panic(err)
	}

	if err := os.RemoveAll("/.oldroot"); err != nil {
		fmt.Printf("Failed to remove old rootfs")
		panic(err)
	}
}

func copyFile(source, target string) {
	contents, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Printf("Failed to open cpuset.mems file of parent cgroup: %s", err)
		panic(err)
	}

	if err := ioutil.WriteFile(target, contents, 0644); err != nil {
		fmt.Printf("Failed to write file %s : %s", target, err)
		panic(err)
	}
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

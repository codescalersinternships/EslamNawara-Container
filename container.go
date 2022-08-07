package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("invalid command")
	}
}

func run() {
	fmt.Println("Running as", os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	cmd.Run()
}
func child() {
	fmt.Println("Running as", os.Getpid())

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("./alpine-fs")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	syscall.Unmount("/proc", 0)
}

func cgroup() {
	cgroups := "/sys/fs/cgroup/containers"
	pids := filepath.Join(cgroups, "pids")
	err := os.Mkdir(filepath.Join(pids, "eslam"), 0777)
	if err != nil && os.IsExist(err) {
		panic(err)
	}
	check(os.WriteFile(path.Join(pids, "pids.max"), []byte("30"), 0700))
	check(os.WriteFile(path.Join(pids, "notify_on_release"), []byte("1"), 0700))
	check(os.WriteFile(path.Join(pids, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

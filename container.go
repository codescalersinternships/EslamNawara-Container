package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		panic("Few arguments")
	}

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
		Credential: &syscall.Credential{
			Uid: uint32(syscall.Getuid()),
			Gid: uint32(syscall.Getgid()),
		},
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		UidMappings:  []syscall.SysProcIDMap{{ContainerID: 0, HostID: syscall.Getuid(), Size: 1}},
		GidMappings:  []syscall.SysProcIDMap{{ContainerID: 0, HostID: syscall.Getgid(), Size: 1}},
	}

	check(cmd.Run())
}

func child() {
	fmt.Println("Running as", os.Getpid())

	cgroup()

	check(syscall.Sethostname([]byte(NameGenerator())))
	check(syscall.Chroot("./rootfs"))
	check(syscall.Chdir("/"))
	check(syscall.Mount("proc", "proc", "proc", 0, "rw"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	check(cmd.Run())
	check(syscall.Unmount("/proc", 0))
}

func cgroup() {
	pids := "/sys/fs/cgroup/pids"
	os.MkdirAll(pids, 0755)

	check(os.WriteFile(path.Join(pids, "pids.max"), []byte("10"), 0700))
	check(os.WriteFile(path.Join(pids, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

package internal

import "syscall"

const (
	sysSigChild syscall.Signal = syscall.SIGCHLD
)

var sysProcAttrSetPgid = &syscall.SysProcAttr{Setpgid: true}

func sysGetPgid(pid int) (pgid int, err error) {
	return syscall.Getpgid(pid)
}

func sysKill(pid int, sig syscall.Signal) (err error) {
	return syscall.Kill(pid, sig)
}

func sysWaitFor(pid int, wstatus *syscall.WaitStatus, options int, rusage *syscall.Rusage) (wpid int, err error) {
	return syscall.Wait4(pid, wstatus, options, rusage)
}

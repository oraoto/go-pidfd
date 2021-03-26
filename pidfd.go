// +build linux

// Package pidfd provides pidfd_open, pidfd_getfd, pidfd_send_signal support on linux 5.6+.
package pidfd

import "syscall"

const (
	sys_PIDFD_SEND_SIGNAL = 424
	sys_PIDFD_OPEN        = 434
	sys_PIDFD_GETFD       = 438
)

// pidfd, a file descriptor that refers to a process.
type PidFd int

// Obtain a file descriptor that refers to a process.
//
// The flags argument is reserved for future use; currently, this argument must be specified as 0.
func Open(pid int, flags uint) (PidFd, error) {
	fd, _, errno := syscall.Syscall(sys_PIDFD_OPEN, uintptr(pid), uintptr(flags), 0)
	if errno != 0 {
		return 0, errno
	}
	return PidFd(fd), nil
}

// Obtain a duplicate of another process's file descriptor.
//
// The flags argument is reserved for future use; currently, this argument must be specified as 0.
//
// PTRACE_MODE_ATTACH_REALCREDS permission is required.
func (fd PidFd) GetFd(targetfd int, flags uint) (int, error) {
	newfd, _, errno := syscall.Syscall(sys_PIDFD_GETFD, uintptr(fd), uintptr(targetfd), uintptr(flags))

	if errno != 0 {
		return 0, errno
	}
	return int(newfd), nil
}

// Send a signal to a process specified by a pidfd.
//
// The flags argument is reserved for future use; currently, this argument must be specified as 0.
func (fd PidFd) SendSignal(signal syscall.Signal, flags uint) error {
	_, _, errno := syscall.Syscall6(sys_PIDFD_SEND_SIGNAL, uintptr(fd), uintptr(signal), 0, uintptr(flags), 0, 0)

	if errno != 0 {
		return errno
	}
	return nil
}

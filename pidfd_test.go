package pidfd_test

import (
	"os/exec"
	"syscall"
	"testing"

	"github.com/oraoto/go-pidfd"
)

func TestOpen(t *testing.T) {
	cmd := exec.Command("sleep", "2")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	p, err := pidfd.Open(cmd.Process.Pid, 0)

	if err != nil {
		t.Fatal(p)
	}
}

func TestGetFD(t *testing.T) {
	cmd := exec.Command("cat")
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	p, err := pidfd.Open(cmd.Process.Pid, 0)

	// get cat's stdin
	stdin, err := p.GetFd(1, 0)
	if err != nil {
		t.Fatal(err)
	}
	// write something to cat's stdin
	syscall.Write(stdin, []byte("Hello\n"))

	// read it back from cat's stdout
	buf := make([]byte, 6)
	stdout.Read(buf)

	if string(buf) != "Hello\n" {
		t.Error("expect Hello\\n, got " + string(buf))
	}
}

func TestSendSignal(t *testing.T) {
	cmd := exec.Command("cat")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	p, err := pidfd.Open(cmd.Process.Pid, 0)

	err = p.SendSignal(syscall.SIGKILL, 0)

	if err != nil {
		t.Fatal(err)
	}

	cmd.Wait()
	if cmd.ProcessState.ExitCode() != -1 {
		t.Fatalf("expect exit code -1, got %d", cmd.ProcessState.ExitCode())
	}
}

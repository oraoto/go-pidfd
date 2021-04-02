package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/oraoto/go-pidfd"
)

var pid = flag.Int("pid", 0, "http server pid")
var fd = flag.Int("fd", 0, "http server listen fd")

func main() {
	flag.Parse()

	if *pid != 0 && *fd != 0 {
		retrieveAndServer(*pid, *fd)
	} else {
		startHttpServer()
	}
}

// start a http server on existing listen fd
func retrieveAndServer(pid int, fd int) {
	p, err := pidfd.Open(pid, 0)
	handleError(err)

	listenFd, err := p.GetFd(fd, 0)
	handleError(err)

	file := os.NewFile(uintptr(listenFd), "")

	ln, err := net.FileListener(file)
	handleError(err)

	err = http.Serve(ln, http.HandlerFunc(handler))
	handleError(err)
}

// start a normal server
func startHttpServer() {
	// print listen fd
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			c.Control(func(fd uintptr) {
				fmt.Printf("Listen fd = %d , pid = %d\n", fd, os.Getpid())
			})
			return nil
		},
	}

	ln, err := lc.Listen(context.Background(), "tcp", ":8080")
	handleError(err)

	err = http.Serve(ln, http.HandlerFunc(handler))
	handleError(err)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from %d \n", os.Getpid())
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

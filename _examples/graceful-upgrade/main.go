package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"log"

	"github.com/oraoto/go-pidfd"
)

func main() {
	server := &Server{}

	server.ListenAndServe(":8001", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %d \n", os.Getpid())
	}))
}

type Server struct {
	listenFd int
	listener net.Listener
	http.Server
}

func (s *Server) ListenAndServe(addr string, handler http.Handler) {
	pid, _ := strconv.Atoi(os.Getenv("_PID"))
	fd, _ := strconv.Atoi(os.Getenv("_FD"))

	// SIGHUP
	go s.handleUpgradeSignal()
	// SIGTERM, SIGINT
	go s.handleShutdownSignal()

	if pid > 0 && fd > 0 {
		s.retrieveAndServer(pid, fd, handler)
	} else {
		s.startHttpServer(addr, handler)
	}
}

func (s *Server) handleUpgradeSignal() {
	upgrade := make(chan os.Signal, 1)
	signal.Notify(upgrade, syscall.SIGHUP)

	for {
		<-upgrade

		cmd := exec.Command(os.Args[0], os.Args[1:]...)

		// pass pid and listenfd to child process
		env := os.Environ()
		env = append(env, fmt.Sprintf("_PID=%d", os.Getpid()))
		env = append(env, fmt.Sprintf("_FD=%d", s.listenFd))
		cmd.Env = env

		if err := cmd.Start(); err != nil {
			log.Print(err)
			cmd.Wait()
		}
	}
}

func (s *Server) handleShutdownSignal() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Println("fail to shoudown:", err)
	}
}

func (s *Server) retrieveAndServer(pid int, fd int, handler http.Handler) {
	p, err := pidfd.Open(pid, 0)
	handleError(err)

	listenFd, err := p.GetFd(fd, 0)
	handleError(err)
	s.listenFd = listenFd

	file := os.NewFile(uintptr(listenFd), "")

	ln, err := net.FileListener(file)
	handleError(err)

	var errorChan = make(chan error)
	go func() {
		errorChan <- http.Serve(ln, handler)
	}()

	select {
	case err := <-errorChan:
		handleError(err)
	case <-time.After(time.Second * 5):
		log.Println("Stop old server")
		p.SendSignal(syscall.SIGTERM, 0)
	}

	err = <-errorChan
	handleError(err)
}

func (s *Server) startHttpServer(addr string, handler http.Handler) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			c.Control(func(fd uintptr) {
				s.listenFd = int(fd)
			})
			return nil
		},
	}

	ln, err := lc.Listen(context.Background(), "tcp", addr)
	handleError(err)
	s.listener = ln

	s.Addr = addr
	s.Handler = handler

	s.Serve(s.listener)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

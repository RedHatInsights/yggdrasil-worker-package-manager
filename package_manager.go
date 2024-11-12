package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/subpop/go-log"
)

type ExitError struct {
	err *exec.ExitError
}

func (e ExitError) Error() string {
	return "program exited: " + e.err.Error()
}

type Installer interface {
	Install(name string) (stdout, stderr []byte, code int, err error)
}

type Uninstaller interface {
	Uninstall(name string) (stdout, stderr []byte, code int, err error)
}

type RepositoryManager interface {
	AddRepo(name string, content []byte) (stdout, stderr []byte, code int, err error)
	RemoveRepo(name string) (stdout, stderr []byte, code int, err error)
	EnableRepo(name string) (stdout, stderr []byte, code int, err error)
	DisableRepo(name string) (stdout, stderr []byte, code int, err error)
}

type PackageManager interface {
	Installer
	Uninstaller
	RepositoryManager
	Stdout() chan []byte
	Stderr() chan []byte
}

// run runs cmd and returns the command's stdout, stderr and exit code. If
// stdoutFunc is not nil, a pipe is attached to the command's stdout. Bytes are
// read from stdout into a buffer and stdoutFunc is called. The same is true for
// stderrFunc.
func run(
	cmd *exec.Cmd,
	stdoutFunc func(buf []byte),
	stderrFunc func(buf []byte),
) (stdout, stderr []byte, code int, err error) {
	outb := new(bytes.Buffer)
	errb := new(bytes.Buffer)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, -1, fmt.Errorf("cannot connect stdout to pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, -1, fmt.Errorf("cannot connect stderr to pipe: %v", err)
	}

	outt := io.TeeReader(stdoutPipe, outb)
	errt := io.TeeReader(stderrPipe, errb)

	go read(outt, stdoutFunc)
	go read(errt, stderrFunc)

	if err := cmd.Run(); err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			return outb.Bytes(), errb.Bytes(), e.ExitCode(), ExitError{e}
		default:
			return nil, nil, -1, fmt.Errorf("failed to start program: %w", err)
		}
	}

	return outb.Bytes(), errb.Bytes(), cmd.ProcessState.ExitCode(), nil
}

// read reads from r into a buffer and calls f with the resulting buffer.
func read(r io.Reader, f func(buf []byte)) {
	for {
		buf := make([]byte, 4096)
		n, err := r.Read(buf)
		if err != nil {
			switch err {
			case io.EOF:
				log.Debugf("reached EOF: %v", err)
				return
			default:
				log.Errorf("cannot read from reader: %v", err)
				continue
			}
		}
		if n > 0 {
			f(buf)
		}
	}
}

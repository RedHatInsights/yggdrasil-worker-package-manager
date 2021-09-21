package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Installer interface {
	Install(name string) (stdout, stderr []byte, code int, err error)
}

type Uninstaller interface {
	Uninstall(name string) (stdout, stderr []byte, code int, err error)
}

type PackageManager interface {
	Installer
	Uninstaller
}

func run(cmd *exec.Cmd) (stdout, stderr []byte, code int, err error) {
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			return outb.Bytes(), errb.Bytes(), e.ExitCode(), fmt.Errorf("failed running program: %w", e)
		default:
			return nil, nil, -1, fmt.Errorf("failed to start program: %w", err)
		}
	}

	return outb.Bytes(), errb.Bytes(), 0, nil
}

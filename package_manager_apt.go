package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const sourcesFile = "/etc/apt/sources.list.d/yggdrasil-package-manager.list"

type PackageManagerApt struct{}

func (p *PackageManagerApt) Install(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("install", name)
}

func (p *PackageManagerApt) Uninstall(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("remove", name)
}

func (p *PackageManagerApt) AddRepo(sourceLine string, _ []byte) (stdout, stderr []byte, code int, err error) {
	return p.EnableRepo(sourceLine)
}

func (p *PackageManagerApt) RemoveRepo(sourceLine string) (stdout, stderr []byte, code int, err error) {
	return p.DisableRepo(sourceLine)
}

func (p *PackageManagerApt) EnableRepo(sourceLine string) (stdout, stderr []byte, code int, err error) {
	if err := os.MkdirAll(filepath.Base(sourcesFile), 0755); err != nil {
		return nil, nil, -1, fmt.Errorf("cannot create sources list directory: %w", err)
	}

	if err := writeLines(sourcesFile, []string{sourceLine + "\n"}, false); err != nil {
		return nil, nil, 1, fmt.Errorf("cannot write to sources list file: %w", err)
	}

	return nil, nil, 0, nil
}

func (p *PackageManagerApt) DisableRepo(sourceLine string) (stdout, stderr []byte, code int, err error) {
	lines, err := readLines(sourcesFile)
	if err != nil {
		return nil, nil, -1, fmt.Errorf("cannot read sources list file: %w", err)
	}

	for i, l := range lines {
		if l == sourceLine {
			copy(lines[i:], lines[i+1:])
			lines = lines[:len(lines)-1]
			break
		}
	}

	if err := writeLines(sourcesFile, lines, true); err != nil {
		return nil, nil, -1, fmt.Errorf("cannot write sources list file: %w", err)
	}

	return nil, nil, 0, nil
}

func (p *PackageManagerApt) run(command string, args ...string) (stdout, stderr []byte, code int, err error) {
	cmdargs := []string{"--assume-yes", command}
	cmdargs = append(cmdargs, args...)
	cmd := exec.Command("/usr/bin/apt-get", cmdargs...)
	stdout, stderr, code, err = run(cmd)
	return
}

package main

import (
	"os/exec"
)

type PackageManagerApt struct{}

func (p *PackageManagerApt) Install(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("install", name)
}

func (p *PackageManagerApt) Uninstall(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("remove", name)
}

func (p *PackageManagerApt) run(command, name string) (stdout, stderr []byte, code int, err error) {
	cmd := exec.Command("/usr/bin/apt", command, "--assumeyes", name)
	stdout, stderr, code, err = run(cmd)
	return
}

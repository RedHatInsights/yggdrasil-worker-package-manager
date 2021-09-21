package main

import (
	"os/exec"
)

type PackageManagerDnf struct{}

func (p *PackageManagerDnf) Install(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("install", name)
}

func (p *PackageManagerDnf) Uninstall(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("remove", name)
}

func (p *PackageManagerDnf) run(command, name string) (stdout, stderr []byte, code int, err error) {
	cmd := exec.Command("/usr/bin/dnf", command, "--assumeyes", name)
	stdout, stderr, code, err = run(cmd)
	return
}

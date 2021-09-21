package main

import (
	"os/exec"
)

type PackageManagerYum struct{}

func (p *PackageManagerYum) Install(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("install", name)
}

func (p *PackageManagerYum) Uninstall(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("remove", name)
}

func (p *PackageManagerYum) run(command, name string) (stdout, stderr []byte, code int, err error) {
	cmd := exec.Command("/usr/bin/yum", command, "--assumeyes", name)
	stdout, stderr, code, err = run(cmd)
	return
}

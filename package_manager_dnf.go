package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type PackageManagerDnf struct{}

func (p *PackageManagerDnf) Install(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("install", name)
}

func (p *PackageManagerDnf) Uninstall(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("remove", name)
}

func (p *PackageManagerDnf) AddRepo(name string, content []byte) (stdout, stderr []byte, code int, err error) {
	return nil, nil, -1, ioutil.WriteFile(filepath.Join("/etc/yum.repos.d/", canonicalizeRepoName(name, ".repo")), content, 0644)
}

func (p *PackageManagerDnf) RemoveRepo(name string) (stdout, stderr []byte, code int, err error) {
	return nil, nil, -1, os.Remove(filepath.Join("/etc/yum.repos.d/", canonicalizeRepoName(name, ".repo")))
}

func (p *PackageManagerDnf) EnableRepo(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("config-manager", "--enable", name)
}

func (p *PackageManagerDnf) DisableRepo(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("config-manager", "--disable", name)
}

func (p *PackageManagerDnf) run(command string, args ...string) (stdout, stderr []byte, code int, err error) {
	cmdargs := []string{"--assumeyes", command}
	cmdargs = append(cmdargs, args...)

	cmd := exec.Command("/usr/bin/dnf", cmdargs...)
	stdout, stderr, code, err = run(cmd)
	return
}

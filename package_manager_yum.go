package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type PackageManagerYum struct{}

func (p *PackageManagerYum) Install(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("install", name)
}

func (p *PackageManagerYum) Uninstall(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("remove", name)
}

func (p *PackageManagerYum) AddRepo(name string, content []byte) (stdout, stderr []byte, code int, err error) {
	return nil, nil, -1, ioutil.WriteFile(filepath.Join("/etc/yum.repos.d/", canonicalizeRepoName(name, ".repo")), content, 0644)
}

func (p *PackageManagerYum) RemoveRepo(name string) (stdout, stderr []byte, code int, err error) {
	return nil, nil, -1, os.Remove(filepath.Join("/etc/yum.repos.d/", canonicalizeRepoName(name, ".repo")))
}

func (p *PackageManagerYum) EnableRepo(name string) (stdout, stderr []byte, code int, err error) {
	cmd := exec.Command("/usr/bin/yum-config-manager", "--assumeyes", "--enable", name)
	stdout, stderr, code, err = run(cmd)
	return
}

func (p *PackageManagerYum) DisableRepo(name string) (stdout, stderr []byte, code int, err error) {
	cmd := exec.Command("/usr/bin/yum-config-manager", "--assumeyes", "--disable", name)
	stdout, stderr, code, err = run(cmd)
	return
}

func (p *PackageManagerYum) run(command string, args ...string) (stdout, stderr []byte, code int, err error) {
	cmdargs := []string{"--assumeyes", command}
	cmdargs = append(cmdargs, args...)

	cmd := exec.Command("/usr/bin/yum", cmdargs...)
	stdout, stderr, code, err = run(cmd)
	return
}

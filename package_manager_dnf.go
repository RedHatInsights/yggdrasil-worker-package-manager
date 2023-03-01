package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

type PackageManagerDnf struct {
	stdout chan []byte
	stderr chan []byte
}

func (p *PackageManagerDnf) Install(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("install", name)
}

func (p *PackageManagerDnf) Uninstall(name string) (stdout, stderr []byte, code int, err error) {
	return p.run("remove", name)
}

func (p *PackageManagerDnf) AddRepo(name string, content []byte) (stdout, stderr []byte, code int, err error) {
	return nil, nil, -1, os.WriteFile(filepath.Join("/etc/yum.repos.d/", canonicalizeRepoName(name, ".repo")), content, 0644)
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

func (p *PackageManagerDnf) Stdout() chan []byte {
	return p.stdout
}

func (p *PackageManagerDnf) Stderr() chan []byte {
	return p.stderr
}

func (p *PackageManagerDnf) run(command string, args ...string) (stdout, stderr []byte, code int, err error) {
	cmdargs := []string{"--assumeyes", command}
	cmdargs = append(cmdargs, args...)

	cmd := exec.Command("/usr/bin/dnf", cmdargs...)
	stdout, stderr, code, err = run(cmd, p.sendStdout, p.sendStderr)
	return
}

func (p *PackageManagerDnf) sendStdout(buf []byte) {
	p.stdout <- buf
}

func (p *PackageManagerDnf) sendStderr(buf []byte) {
	p.stderr <- buf
}

# yggdrasil-worker-package-manager

`yggdrasil-worker-package-manager` is a simple package manager yggd worker. It
knows how to install and remove packages, add, remove, enable and disable
repositories, and does rudimentary detection of the host it is running on to guess
the package manager to use. It only installs packages that match one of the
provided `allow-pattern` regular expressions.

# Installation

The easiest way to compile and install `yggdrasil-worker-package-manager` is
using `meson`. Because it runs as a bus-activatable D-Bus service, files must be
installed in specific directories.

Generally, it is recommended to follow your distribution's packaging guidelines
for compiling Go programs and installing projects using `meson`. What follows is a
generally acceptable set of steps to setup, compile and install yggdrasil using
`meson`.

```
# Set up the project according to distribution-specific directory locations
meson setup --prefix /usr/local --sysconfdir /etc --localstatedir /var builddir
# Compile
meson compile -C builddir
# Install
meson install -C builddir
```

`meson` includes an optional `--destdir` to its `install` subcommand to aid in
packaging.

# Configuration

Default configuration values are documented in `config.toml`. When installing
using meson, this file is installed into `/etc/yggdrasil-worker-package-manager`
by default.

# Usage

The worker will register itself as a handler for the "package_manager"
directive. It expects to receive messages with the following JSON schema:

```json
{
    "type": "object",
    "properties": {
        "command": { "enum": ["install", "remove", "enable-repo", "disable-repo", "add-repo", "remove-repo"] },
        "name": { "type": "string" },
        "content": { "type": "string" }
    },
    "required": ["command", "name" ]
}
```

## Examples

### Install `vim`

```json
{
    "command": "install",
    "name": "vim"
}
```

### Enable "updates-testing" repository

```json
{
    "command": "enable-repo",
    "name": "updates-testing"
}
```

### Add custom repository on a dnf or yum client

```json
{
    "command": "add-repo",
    "name": "my-custom-repo",
    "content": "[my-custom-repo]\nbaseurl=http://servername/path/to/repo\nenabled=1"
}
```

### Add custom repository on an apt client

```json
{
    "command": "add-repo",
    "name": "deb http://servername path component"
}
```

# Permitting operations

Before an operation on a package is permitted, the value of the `name` field is
matched against each regular expression specified in the `allow-pattern`
configuration value. **Only** if a package name matches one of the patterns in
that array, is it permitted to be installed.

For example, given the `allow-patterns` value:

```toml
allow-patterns = ["^vim.*"]
```

Only packages that begin with "vim" are allowed to be installed or removed.

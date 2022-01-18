# yggdrasil-worker-package-manager

`yggdrasil-worker-package-manager` is a simple package manager yggd worker. It
knows how to install and remove packages, add, remove, enable and disable
repositories, and does rudamentary detection of the host its running on to guess
the package manager to use. It only installs packages that match one of the
provided `allow-pattern` regular expressions.

# Installation

Compile the worker and install it into the `yggd` worker directory:

```
go build -o yggd-package-manager-worker .
install -D -m 755 yggd-package-manager-worker $(pkg-config --variable=workerexecdir yggdrasil)/
```

# Configuration

Default configuration values are documented in `config.toml`. To adjust any
configuration values, edit `config.toml` and copy into the `workerconfigdir`.

```
install -D -m 644 config.toml $(pkg-config --variable=workerconfdir)/package-manager.toml
```

# Usage

The worker will register itself as a handler for the "package-manager"
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

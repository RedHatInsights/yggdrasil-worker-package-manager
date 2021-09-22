# yggdrasil-worker-package-manager

`yggdrasil-worker-package-manager` is a simple package manager yggd worker. It
knows how to install and remove packages, and does rudamentary detection of the
host its running on to guess the package manager to use. It only installs
packages that match one of the provided `allow-pattern` regular expressions.

# Installation

Compile the worker and install it into the `yggd` worker directory:

```
go build ./ -o yggd-package-manager-worker
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
        "command": { "enum": ["install", "remove"] },
        "name": { "type": "string" }
    },
    "required": ["command", "name" ]
}
```

For example, to tell `yggd-package-manager-worker` to install "vim", send it:

```json
{
    "command": "install",
    "name": "vim"
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

# AGENTS.md

Read [`README.md`](README.md) first for project overview, build/install
instructions, D-Bus message schema with examples, and `allow-pattern`
configuration.

This file covers behavioral guidance for coding agents: what to watch out
for, how to test, and how to contribute.

## Key source files

| Path | Purpose |
|------|---------|
| `main.go` | Entry point: flag parsing, worker registration, message dispatch |
| `package_manager.go` | `PackageManager` interface and `run()` execution helper |
| `package_manager_dnf.go` | DNF backend |
| `package_manager_yum.go` | YUM backend |
| `package_manager_apt.go` | APT backend |
| `util.go` | File I/O helpers |
| `config.toml` | Default config: `allow-pattern` regexes and `log-level` |
| `data/` | D-Bus policy and systemd unit templates |
| `dist/` | RPM spec file template |
| `systemtest/` | FMF/TMT system tests (D-Bus ownership, help output) |

## Key dependencies

| Module | Purpose |
|--------|---------|
| `github.com/redhatinsights/yggdrasil/worker` | D-Bus worker SDK for registration and messaging |
| `github.com/zcalusic/sysinfo` | OS detection for backend selection |
| `github.com/peterbourgon/ff/v3` | Config file + env var + flag parsing |
| `github.com/subpop/go-log` | Leveled logging |
| `github.com/sgreben/flagvar` | Regex flag type for `allow-pattern` |
| `github.com/google/uuid` | Response message IDs |

## D-Bus message handling

See [`README.md`](README.md) for the JSON message schema and examples.

Key implementation files:
- `Message` struct in [`main.go`](main.go) defines the inbound payload
- `run()` helper in [`package_manager.go`](package_manager.go) executes
  commands and streams output via stdout/stderr channels
- `dataRx` handler in [`main.go`](main.go) reads those channels and emits
  D-Bus events during execution

## Testing

Before submitting changes, mirror the CI checks locally:

```bash
go build -v ./...
go test -v ./...
go vet -v ./...
golangci-lint run --verbose --timeout=3m
```

CI also runs **commitsar** (conventional commits) and **woke** (inclusive
language). See [`.github/workflows/`](.github/workflows/) for details.

System tests live in [`systemtest/`](systemtest/plans/main.fmf) and run via
TMT/Testing Farm on PRs through Packit.

## Security

The worker runs as a configurable system user (`worker_user` meson option) and
executes package manager commands (`dnf`, `yum`, `apt-get`) with `--assumeyes`.
When modifying or generating code:

- **Input validation matters.** Package names and repo names are matched against
  `allow-pattern` regexes before operations proceed. Never bypass this check.
- Avoid shell injection — commands are built with `exec.Command` argument lists,
  not shell strings. Keep it that way.
- Do not log secrets (certificates, keys, broker credentials, tokens).
- Repo file content (for `add-repo`) is written to disk — validate paths and
  sanitize filenames before file I/O.

## Code style and architecture

- **Go 1.21**, module `github.com/redhatinsights/yggdrasil-worker-package-manager`.
- Format with `gofmt` and `goimports` before committing.
- The codebase is intentionally small (~650 lines). New package manager backends
  should implement the `PackageManager` interface defined in
  `package_manager.go` and add OS detection logic in `detectPackageManager()`
  in `main.go`.
- The worker uses the `github.com/redhatinsights/yggdrasil/worker` SDK — follow
  existing patterns for D-Bus registration and message handling.

## Pull requests

- Ensure all CI checks pass (build, unit tests, vet, lint, woke) before
  requesting review.
- Add or update Go unit tests for changed logic.
- Update system tests when behavior visible to D-Bus clients changes.

## Git commits

Follow [Conventional Commits](https://www.conventionalcommits.org):

- Subject completes: "when applied, this commit will…"
- Subject and body are separated by a blank line.
- Body expands with relevant detail; detail items start with `* `.
- Prefix types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, etc.

Example:

```text
docs: add AGENTS.md for AI-assisted development

* Summarize build, test, CI, and contribution conventions for coding agents
```

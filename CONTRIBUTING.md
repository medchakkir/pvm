# Contributing to PVM

Thank you for your interest in contributing! PVM is a small, focused project and contributions of all sizes are welcome — bug fixes, new features, better error messages, or documentation improvements.

---

## Setting Up

```bash
git clone https://github.com/medchakkir/pvm.git
cd pvm
go mod tidy
go build -o pvm.exe .
```

You need **Go 1.21+** installed. That's the only dependency.

---

## Project Structure

```
pvm/
├── cmd/                  # One file per CLI command
│   ├── root.go           # Root command + global flags
│   ├── install.go
│   ├── use.go
│   ├── list.go
│   ├── list_remote.go
│   ├── current.go
│   ├── uninstall.go
│   └── bin.go
├── internal/
│   ├── config/           # PVM home directory management
│   ├── php/              # Version parsing, fetching, installing
│   └── env/              # Shim and PATH management
├── main.go
└── .goreleaser.yaml
```

Each command lives in its own file. If you're adding a new command, create a new file in `cmd/` following the same pattern — define the command, implement `RunE`, register in `init()`.

---

## Running Tests

```bash
go test ./...
```

Please add or update tests when changing logic in `internal/`. The `cmd/` layer doesn't need unit tests — focus on the packages doing the real work.

---

## Submitting a PR

1. Fork the repo and create a branch: `git checkout -b feat/my-feature`
2. Make your changes
3. Run `go test ./...` — all tests must pass
4. Run `go build -o pvm.exe .` — must compile cleanly
5. Open a pull request with a clear description of what you changed and why

---

## Good First Issues

Look for issues tagged [`good first issue`](https://github.com/medchakkir/pvm/issues?q=label%3A%22good+first+issue%22) — these are small, well-scoped tasks that are great for getting familiar with the codebase.

---

## Code Style

- Standard Go formatting — run `gofmt` before committing
- Error messages start with `✗`, success messages with `✓`
- Keep commands simple — business logic belongs in `internal/`, not in `cmd/`

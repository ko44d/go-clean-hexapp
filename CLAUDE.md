# CLAUDE.md

Go template repository implementing Hexagonal Architecture. Before adding or modifying any feature, read `docs/DESIGN.md` to understand the architecture constraints.

## Commands

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/domain/task/...
go test ./internal/usecase/task/...
go test ./internal/interface/handler/...

# Regenerate all mocks (after changing interfaces)
go generate ./...

# Build the binary
go build -o go-clean-hexapp ./cmd/server

# Run locally with Docker Compose (includes PostgreSQL)
docker-compose up
```

## Definition of Done

Before marking any work complete:
- `go test ./...` passes
- `go build -o go-clean-hexapp ./cmd/server` succeeds
- Commits are made in meaningful units

## Architecture (summary)

Dependency direction: `handler → usecase → domain ← repository`

- Domain layer owns the `Repository` interface — not the infrastructure layer
- Handler depends on `Interactor` interface — never on the concrete usecase struct
- No ORM. Raw SQL via `pgx` only.

For full architecture diagram, layer responsibilities, and API spec → `docs/DESIGN.md`

## Coding Style

Follow [Effective Go](https://go.dev/doc/effective_go). Key rules enforced in this project:

- Constructors are named `New()`, not `NewFoo()` — the package name provides the context
- Receiver names are abbreviated from the type (`t` for `*Task`); never `this` or `self`
- Errors are sentinel values (`var ErrFoo = errors.New("...")`); never compare error strings
- Wrap errors with `fmt.Errorf("context: %w", err)` so the failure site is always traceable
- Interfaces belong in the consuming package, not the implementing package
- Package names are singular and lowercase (`task`, not `tasks`)

## Testing Pattern

Tests use **Ginkgo v2 + Gomega** (BDD style). Mocks are generated with `go.uber.org/mock/mockgen`.

- Run `go generate ./...` whenever an interface changes
- Mock locations: `internal/domain/task/mocks/` and `internal/usecase/task/mocks/`
- Unit tests use mocks; no integration tests against a real database

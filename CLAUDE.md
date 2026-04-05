# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/domain/task/...
go test ./internal/usecase/task/...
go test ./internal/interface/handler/...

# Run with verbose output
go test -v ./...

# Regenerate all mocks (after changing interfaces)
go generate ./...

# Build the binary
go build -o go-clean-hexapp ./cmd/server

# Run locally with Docker Compose (includes PostgreSQL)
docker-compose up
```

## Definition of Done

- `go test ./...`
- `go build -o go-clean-hexapp ./cmd/server`
- Commits are made in meaningful units.

## Architecture

This project implements **Hexagonal Architecture (Ports & Adapters)** combined with Clean Architecture. Dependencies flow inward — outer layers depend on inner layers, never the reverse.

```
cmd/server/main.go
    ↓
config/ + internal/container/   ← Wires all dependencies
    ↓
internal/router/                ← Gin HTTP routing
    ↓
internal/interface/handler/     ← HTTP adapter (inbound)
    ↓
internal/usecase/task/          ← Application/business logic
    ↓
internal/domain/task/           ← Core entities + Repository interface (port)
    ↑
internal/interface/repository/  ← PostgreSQL adapter (outbound)
    ↑
internal/infrastructure/db/     ← pgx connection pool
```

### Layer Responsibilities

| Layer | Package | Responsibility |
|---|---|---|
| Domain | `internal/domain/task/` | Task entity, validation, domain errors, Repository **interface** |
| Usecase | `internal/usecase/task/` | Orchestrates domain + repository; defines Interactor **interface** |
| Interface | `internal/interface/handler/` | HTTP request/response handling, JSON mapping |
| Interface | `internal/interface/repository/` | PostgreSQL implementation of domain.Repository |
| Infrastructure | `internal/infrastructure/db/` | pgx connection pool |
| Container | `internal/container/` | Manual dependency injection — wires everything together |

## Coding Style

Follow [Effective Go](https://go.dev/doc/effective_go) as the baseline coding style.
Key conventions enforced in this project:
- Name constructors `New()` rather than `NewFoo()` — the package name already provides context
- Use abbreviated receiver names derived from the type (e.g., `t` for `*Task`); never use `this` or `self`
- Define errors as sentinel values (`var ErrFoo = errors.New("...")`); never compare error strings
- Wrap errors with `fmt.Errorf("context: %w", err)` so the failure site is always traceable
- Define interfaces in the consuming package (e.g., `Repository` is defined in the domain layer and implemented by the infrastructure layer)
- Package names are singular and lowercase (`task`, not `tasks`)

### Key Architectural Points

- The `Repository` interface is **defined in the domain layer** (`internal/domain/task/repository.go`), not in the infrastructure layer. This is the core of hexagonal architecture — the domain owns the port.
- The `Interactor` interface is defined in the usecase layer (`internal/usecase/task/interactor.go`), and the HTTP handler depends on it. This allows handler tests to use a mock interactor.
- No ORM — raw SQL via `pgx`.
- No Makefile — use `go` commands directly.

### Testing Pattern

Tests use **Ginkgo v2 + Gomega** (BDD style). Mocks are generated with `go.uber.org/mock/mockgen`.

Use `go generate ./...` to regenerate all mocks at once after interface changes.

- `go:generate` directives are defined on the interface source files: `internal/domain/task/repository.go` and `internal/usecase/task/interactor.go`
- `internal/domain/task/mocks/` — generated mocks for the domain-layer `Repository` interface (used in usecase tests)
- `internal/usecase/task/mocks/` — generated mocks for the usecase-layer `Interactor` interface (used in handler tests)

### API Endpoints

- `GET /tasks` — list all tasks
- `POST /tasks` — create task; body: `{"title": "..."}`
- `POST /tasks/complete?id=<uuid>` — mark task complete

### Configuration

All configuration is via environment variables (see `config/config.go`).

| Variable | Default |
|---|---|
| `POSTGRES_HOST` | `(required, no default)` |
| `POSTGRES_PORT` | `(required, no default)` |
| `POSTGRES_USER` | `(required, no default)` |
| `POSTGRES_PASSWORD` | `(required, no default)` |
| `POSTGRES_DB` | `(required, no default)` |
| `POSTGRES_SSLMODE` | `(required, no default)` |
| `PORT` | `8080` |

Refer to `.env.example` for a ready-to-use local configuration template.

### Task Status Values

Task status uses lowercase strings: `"todo"` and `"complete"`.

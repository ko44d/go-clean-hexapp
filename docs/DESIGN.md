# Design Document

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

### Key Architectural Decisions

- The `Repository` interface is **defined in the domain layer** (`internal/domain/task/repository.go`), not in the infrastructure layer. This is the core of hexagonal architecture — the domain owns the port.
- The `Interactor` interface is defined in the usecase layer (`internal/usecase/task/interactor.go`), and the HTTP handler depends on it. This allows handler tests to use a mock interactor.
- No ORM — raw SQL via `pgx`.
- No Makefile — use `go` commands directly.

## Testing

Tests use **Ginkgo v2 + Gomega** (BDD style). Mocks are generated with `go.uber.org/mock/mockgen`.

Use `go generate ./...` to regenerate all mocks at once after interface changes.

- `go:generate` directives are defined on the interface source files: `internal/domain/task/repository.go` and `internal/usecase/task/interactor.go`
- `internal/domain/task/mocks/` — generated mocks for the domain-layer `Repository` interface (used in usecase tests)
- `internal/usecase/task/mocks/` — generated mocks for the usecase-layer `Interactor` interface (used in handler tests)

## API Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/tasks` | List all tasks |
| POST | `/tasks` | Create task; body: `{"title": "..."}` |
| POST | `/tasks/complete?id=uuid` | Mark task complete |

## Configuration

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

## Domain Constraints

- Task status uses lowercase strings: `"todo"` and `"complete"`.

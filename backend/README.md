# Backend (Go)

Clean architecture scaffolding for the Go API.

## Run

```
go run ./cmd/api
```

Default address is `:54321`. Override with:

```
set EP_HTTP_ADDRESS=:54321
```

## Structure

- `cmd/api` entrypoint
- `internal/domain` entities and rules
- `internal/usecase` application services
- `internal/adapter/http` HTTP handlers and routing
- `internal/adapter/ws` websocket gateway (placeholder)
- `internal/infra` database, jobs, integrations

# Development Environment Guide

## Overview

This document describes the local development setup for EmpoweredPixels. The development environment uses Docker Compose with hot module replacement (HMR) for both frontend and backend.

## Prerequisites

- Docker & Docker Compose installed
- Make (optional, for convenience commands)
- Node.js not required (containerized)
- Go not required (containerized)

## Quick Start

1. **Copy local env template**
   ```bash
   cp .env.local.template .env.local
   # Edit .env.local if you need custom values
   ```

2. **Start dev environment**
   ```bash
   make dev
   # or: docker compose --profile dev up -d --build
   ```

3. **Access services**
   - Frontend (Vite dev server): http://localhost:5173
   - Backend API: http://localhost:49101
   - Database: localhost:5432 (user: `postgres`, password: `postgres`)

4. **View logs**
   ```bash
   make logs
   ```

5. **Run tests**
   ```bash
   make test
   # or: PLAYWRIGHT_BASE_URL=http://localhost:5173 npx playwright test
   ```

## Architecture

### Frontend

- **Vite dev server** with HMR on port 5173
- Source mounted as bind volume for instant updates
- `VITE_API_BASE_URL` set to `http://localhost:49101` (points to backend)
- Build production: `make build` or `docker-compose build frontend`

### Backend

- **Air** (`cosmtrek/air`) for live reload
- Go source mounted; binary rebuilt automatically on changes
- Structured JSON logging in dev (`LOG_FORMAT=json`, `LOG_LEVEL=debug`)
- Database migrations run automatically on startup (idempotent)

### Database

- PostgreSQL 16 Alpine
- Volume persisted: `ep_postgres_data`
- Migrations loaded from `backend/internal/infra/d/migrations`
- Health checks enabled

## Environment Files

- `.env` — Base configuration (committed)
- `.env.local` — Local overrides (gitignored, copy from template)
- `.env.production` — Production frontend build vars

## Development Workflow

1. Make code changes in `frontend/` or `backend/`
2. HMR updates instantly in browser
3. Backend restarts automatically (within ~1s)
4. Use `make logs` to monitor both services

## Testing

- **E2E**: Playwright tests in `frontend/tests/e2e/`
- Run against dev: `make test`
- Watch mode: `npx playwright test --watch`

## Cleaning Up

```bash
# Stop dev services only
make clean

# Full purge (removes all images, volumes, networks)
make purge
```

**Warning**: `make purge` deletes all Docker caches and volumes. Production data will be lost. Use only for development reset.

## Production vs Dev

| Aspect          | Dev                              | Production                     |
|-----------------|----------------------------------|--------------------------------|
| Frontend        | Vite dev server (5173)           | Nginx + built assets (49100)   |
| Backend         | Air live reload                  | Static binary                  |
| Database        | Local volume                     | External managed DB            |
| Logging         | JSON, debug level                | Structured, info+              |
| HMR             | Enabled                          | No                             |

## Troubleshooting

**Port conflicts**: Ensure ports 5173, 49101, 5432 are free.

**Stale containers**: `make purge` then `make dev`.

**Database connection refused**: Wait a few seconds after `make dev` for Postgres to become healthy.

**Air not reloading**: Ensure `.air.toml` exists in backend directory and source files are inside mounted `./backend`.

## Notes

- Dev environment does **not** use nginx; services accessed directly.
- Frontend dev server proxies API requests to `localhost:49101` via `VITE_API_BASE_URL`.
- Backend auto-reload works for `.go` file changes; static assets (templates, config) may require manual restart.

## Version Label

The version label (`v0.1.0`) is displayed in the frontend navigation bar on all screen sizes. It is sourced from `package.json`.

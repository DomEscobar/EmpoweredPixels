# Empowered Pixels

A monorepo containing the Go backend and Vue frontend for Empowered Pixels.

## Project Structure

- **`backend/`**: Go API following Clean Architecture conventions.
  - Entry point: `cmd/api/main.go`
  - Runs on port `:54321` by default.
- **`frontend/`**: Vue 3 + TypeScript + Tailwind CSS application.
  - Runs on port `:5173` (Vite default).
- **`EP_OLD/`**: Legacy codebase (reference).

## Getting Started

### Prerequisites

- **Go** (1.22+)
- **Node.js** (20+) & **npm**
- **PostgreSQL** running locally on port `5432` with default credentials (`postgres`/`postgres`) or configured via environment variables.

### Running the Project

You can start both the backend and frontend services using the provided PowerShell script:

```powershell
.\run-dev.ps1
```

This script will:
1. Check for required tools (Go, npm).
2. Install frontend dependencies if missing.
3. Launch the Backend and Frontend in separate terminal windows.

### Docker (production / VPS)

To run frontend, backend, and PostgreSQL together:

```bash
docker compose up -d --build
```

- Frontend: `http://<host>:49100`
- Backend API: `http://<host>:49101`

See [docs/docker-deploy.md](docs/docker-deploy.md) for ports, env vars, and options.

### Manual Startup

If you prefer running them manually:

**Backend:**
```bash
cd backend
go run ./cmd/api/main.go
```

**Frontend:**
```bash
cd frontend
npm install # first time only
npm run dev
```

## Documentation

- [Backend Readme](backend/README.md)
- [Frontend Readme](frontend/README.md)
- [Agent Conventions](AGENTS.md)

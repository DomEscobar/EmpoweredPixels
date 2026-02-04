# Docker deployment (VPS)

Run frontend, backend, and PostgreSQL with Docker Compose.

## Prerequisites

- Docker and Docker Compose on the VPS
- Ports **49100** (frontend) and **49101** (backend) free and open in firewall

## Build and run

From the repo root:

```bash
docker compose up -d --build
```

- Frontend: `http://<VPS_IP>:49100`
- Backend API: `http://<VPS_IP>:49101`

## Configuration

Optional env vars (set in `.env` in repo root or export before `docker compose up`):

| Variable | Description |
|----------|-------------|
| `VITE_API_BASE_URL` | URL the browser uses for API and WebSockets (default: `http://152.53.118.78:49101`) |
| `EP_JWT_SECRET` | Backend JWT signing secret (compose provides a default) |
| `EP_ENGINE_URL` | Match engine service URL (leave empty if not used) |

PostgreSQL data is stored in the named volume `ep_postgres_data`.

## Commands

- Stop: `docker compose down`
- Stop and remove volume: `docker compose down -v`
- View logs: `docker compose logs -f`

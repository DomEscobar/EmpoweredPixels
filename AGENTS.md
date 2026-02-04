# Agents

This file captures agent conventions and a running log of AI-driven changes.

## Conventions

- Keep changes incremental and clean-architecture aligned.
- Prefer explicit boundaries over implicit coupling.
- Record each change with date, scope, and files touched.

## Change Log

### 2026-02-04

- Docker VPS deployment: backend and frontend Dockerfiles, Nginx for static frontend, docker-compose with Postgres; ports 49100 (frontend) and 49101 (backend); deployment docs.
  - `backend/Dockerfile`
  - `frontend/Dockerfile`
  - `frontend/nginx.conf`
  - `docker-compose.yml`
  - `docs/docker-deploy.md`
  - `README.md`
- Tightened Match Viewer layout for viewport fit and mobile stacking.
  - `frontend/src/pages/MatchViewer.vue`

### 2026-02-03

- Overhauled Dashboard with "War Room" aesthetic and real-time data integration.
  - Implemented missing frontend feature stores/APIs for matches, rewards, and leagues.
  - Designed "Command Center" dashboard with KPI cards for active roster, campaigns, combat record, and pending rewards.
  - Added live operations feed (recent match history) and elite operative spotlight (top fighter).
  - Integrated Quick Actions grid for navigation.
  - `frontend/src/features/matches/api.ts`
  - `frontend/src/features/matches/store.ts`
  - `frontend/src/features/rewards/api.ts`
  - `frontend/src/features/rewards/store.ts`
  - `frontend/src/features/leagues/api.ts`
  - `frontend/src/features/leagues/store.ts`
  - `frontend/src/pages/Dashboard.vue`

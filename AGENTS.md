# Agents

This file captures agent conventions and a running log of AI-driven changes.

## Conventions

- Keep changes incremental and clean-architecture aligned.
- Prefer explicit boundaries over implicit coupling.
- Record each change with date, scope, and files touched.
- **Branch-based Development**: All agents work on specialized branches; NO direct pushes to `main`.
- **Automated PR Review & Testing**: The `Senior-Code-Architect` merges PRs only after a successful `go test ./...` and build.
- **Epic-based Planning**: The `Senior-Product-Owner` manages requirements in Epics within `KANBAN.md`.
- **The Essence Rule**: All features must align with the core Web3-RPG Indie vibe (approved by PO).
- **Test-Driven Delivery**: Every feature implementation must include a verification step or unit test.
- **Commit Rule**: I'll build the project before committing. I must verify that the code compiles successfully in the local environment.

## Change Log

### 2026-02-05

- **Backend Recovery & Hardening**: Fixed Git history corruption, restored core logic, and implemented security hardening (Auth, PWM Hashing).
- **Mobile UX Sprint**: Implemented Sticky Bottom Nav and compact mobile UI.
- **MCP Integration**: Added MCP server for AI player interaction.

### 2026-02-04

- Docker VPS deployment: backend and frontend Dockerfiles, Nginx for static frontend, docker-compose with Postgres; ports 49100 (frontend) and 49101 (backend); deployment docs.
- Tightened Match Viewer layout for viewport fit and mobile stacking.

### 2026-02-03

- Overhauled Dashboard with "War Room" aesthetic.

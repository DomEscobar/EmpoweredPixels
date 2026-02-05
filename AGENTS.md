# Agents

This file captures agent conventions and a running log of AI-driven changes.

## Agent Loop Architecture (v2 - Self-Healing)

**Current System:** 3-Agent Minimal Core with Heartbeat & Auto-Recovery

### Core Team
| Agent | Role | Responsibility |
|-------|------|----------------|
| **PO-Lead** | Product Owner | Prioritize tasks, assign to Architect, update KANBAN |
| **Architect-Lead** | Senior Code Architect | Implement features, commit code, signal QA |
| **QA-Lead** | QA Specialist | Verify implementations, test builds, report PASS/FAIL |

### Workflow (PO → Architect → QA → PO)
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   PO-Lead   │────→│ Architect   │────→│   QA-Lead   │
│  (Assign)   │     │ (Implement) │     │ (Verify)    │
└──────▲──────┘     └─────────────┘     └──────┬──────┘
       └─────────────────────────────────────────┘
                    (Report PASS/FAIL)
```

### Heartbeat Protocol
All agents MUST send heartbeat every 2-5 minutes:
```bash
/root/agent_heartbeat.sh "<Agent-Name>" "<current-task-status>"
```

### Persistent Loop Mode (Option D - Revised)
Agents run in **25-minute cycles** (1500s timeout) due to system 30min hard limit:
- **Cycle:** Work → Sleep 60s → Check for new tasks → Exit cleanly at 25min
- **Auto-respawn:** Watchdog respawns every 25min if heartbeats stale
- **Continuous:** Agents always alive via respawn cycles
- **Checkpoint:** State in heartbeat/task files survives respawn

### Auto-Recovery
- **Monitor:** `/root/agent_loop_monitor.sh` runs every 5 minutes
- **Detection:** Agent stale after 10 minutes → marked dead
- **Action:** Auto-respawn dead agents with last known task
- **Alert:** Telegram notification on recovery events

### State Files
- `/root/.openclaw/agent_state/*.heartbeat` - Last ping timestamp
- `/root/.openclaw/agent_state/*.task` - Current task description

## Conventions

- Keep changes incremental and clean-architecture aligned.
- Prefer explicit boundaries over implicit coupling.
- Record each change with date, scope, and files touched.
- **Branch-based Development**: All agents work on specialized branches; NO direct pushes to `main`.
- **Automated PR Review & Testing**: The `Senior-Code-Architect` merges PRs via `gh pr merge` ONLY after a successful `go test ./...` and local build verification. This ensures `main` never breaks.
- **Epic-based Planning**: The `Senior-Product-Owner` manages requirements in Epics within `KANBAN.md`.
- **The Essence Rule**: All features must align with the core Web3-RPG Indie vibe (approved by PO).
- **Test-Driven Delivery**: Every feature implementation must include a verification step or unit test.
- **Commit Rule**: I'll build the project before committing. I must verify that the code compiles successfully in the local environment.
- **Heartbeat Rule**: All persistent agents must heartbeat every 2-5 minutes or be considered dead.

## Change Log

### 2026-02-05

- **Agent Loop v2 - Self-Healing System**: Implemented persistent 3-agent core (PO-Lead, Architect-Lead, QA-Lead) with heartbeat protocol and auto-recovery. Agents respawn automatically if stale >10min. Loop monitor runs every 5min.
- **Backend Recovery & Hardening**: Fixed Git history corruption, restored core logic, and implemented security hardening (Auth, PWM Hashing).
- **Mobile UX Sprint**: Implemented Sticky Bottom Nav and compact mobile UI.
- **MCP Integration**: Added MCP server for AI player interaction (REST endpoints, rate limiting, audit logging).

### 2026-02-04

- Docker VPS deployment: backend and frontend Dockerfiles, Nginx for static frontend, docker-compose with Postgres; ports 49100 (frontend) and 49101 (backend); deployment docs.
- Tightened Match Viewer layout for viewport fit and mobile stacking.

### 2026-02-03

- Overhauled Dashboard with "War Room" aesthetic.

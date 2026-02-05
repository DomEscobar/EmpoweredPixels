# Agents

This file captures agent conventions and a running log of AI-driven changes.

## Agent Loop Architecture (v2 - Self-Healing)

**Current System:** 3-Agent Minimal Core with Heartbeat & Auto-Recovery

### Core Team
| Agent | Role | Responsibility |
|-------|------|----------------|
| **PO-Lead** | Product Owner | Prioritize tasks, assign to Architect, update KANBAN, manage Roadmap |
| **Architect-Lead** | Senior Code Architect | Implement features, commit code, signal QA |
| **QA-Lead** | QA Specialist | Verify implementations, test builds, report PASS/FAIL |
| **Senior-Game-Designer** | Game Design | Evaluate features, define mechanics, create user stories |

### Feature Development Workflow

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  Game Designer  │───→│    PO-Lead       │───→│  Architect-Lead │
│ (Feature Ideas) │    │ (Roadmap/Specs)  │    │  (Implement)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                         │
        ┌──────────────────────────────────────────────────┘
        ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   QA-Lead       │───→│  MCP-Live-Test   │───→│   Production    │
│ (Unit/Int/E2E)  │    │  (Final Verify)  │    │   Deploy        │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

#### 1. Feature Definition (PO-Lead)
- **New features** are transferred from Design into the **Roadmap** by PO-Lead
- Every feature has:
  - Clear user stories
  - Technical acceptance criteria
  - Definition of Done (DoD)
  - Estimated effort

#### 2. Test Coverage Requirements (QA-Lead)
Every new feature MUST have:
| Test Type | Coverage | Validated By |
|-----------|----------|--------------|
| **Unit Tests** | Core business logic | QA-Lead |
| **Integration Tests** | API endpoints, DB operations | QA-Lead |
| **E2E Tests** | Full user flows | QA-Lead |
| **MCP Live Tests** | External AI agent compatibility | QA-Lead via MCP |

**MCP Testing:** After implementation, QA-Lead tests the feature via MCP endpoints to ensure external AI agents can interact with it correctly.

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

### Git Push Safety (Anti-Stall)
**Problem:** Agent dies after commit but before push.
**Solutions:**
1. **Immediate Push:** Agent runs `git commit && git push` atomically
2. **Auto-Push Cron:** `/root/autopush_watchdog.sh` runs every 5min
3. **Exit Hook:** Agent ALWAYS pushes before `exit 0`
4. **Mama Monitor:** Main agent detects unpushed commits via `git log origin/main..HEAD`

**Rule:** Never leave commits unpushed >5 minutes.

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
- **Push Rule**: Agent MUST run `git push origin main` immediately after EVERY commit, before any other action.
- **Exit Rule**: Before exit, agent checks `git log origin/main..HEAD` and pushes if commits pending.

### Feature Development Rules

1. **PO-Lead owns the Roadmap**
   - New features come from Game Designer evaluation
   - PO-Lead transfers approved features to Roadmap/KANBAN
   - Features must have clear acceptance criteria before implementation

2. **Complete Test Coverage Required**
   - Unit Tests: Every function with business logic
   - Integration Tests: Every API endpoint
   - E2E Tests: Every user-facing flow
   - MCP Tests: Live verification via AI agent interface

3. **No Feature without Tests**
   - Architect implements feature + tests together
   - QA-Lead validates all test levels
   - MCP-Live-Test is the final gate before production

## Change Log

### 2026-02-05

- **Feature Development Workflow**: Added Senior-Game-Designer agent to core team. Defined complete workflow: Game Designer evaluates features → PO-Lead transfers to Roadmap → Architect implements with tests → QA validates (Unit/Integration/E2E/MCP) → Production.
- **Test Coverage Requirements**: All new features require Unit, Integration, E2E tests PLUS MCP live testing by QA-Lead before production.
- **Git Push Safety System**: Implemented 4-layer protection against unpushed commits - Immediate Push rule, Auto-Push Watchdog (5min cron), Exit Hook enforcement, and Mama Monitor detection. Prevents commit-without-push stall scenario.
- **Agent Loop v2 - Self-Healing System**: Implemented persistent 3-agent core (PO-Lead, Architect-Lead, QA-Lead) with heartbeat protocol and auto-recovery. Agents respawn automatically if stale >10min. Loop monitor runs every 5min.
- **Backend Recovery & Hardening**: Fixed Git history corruption, restored core logic, and implemented security hardening (Auth, PWM Hashing).
- **Mobile UX Sprint**: Implemented Sticky Bottom Nav and compact mobile UI.
- **MCP Integration**: Added MCP server for AI player interaction (REST endpoints, rate limiting, audit logging).

### 2026-02-04

- Docker VPS deployment: backend and frontend Dockerfiles, Nginx for static frontend, docker-compose with Postgres; ports 49100 (frontend) and 49101 (backend); deployment docs.
- Tightened Match Viewer layout for viewport fit and mobile stacking.

### 2026-02-03

- Overhauled Dashboard with "War Room" aesthetic.

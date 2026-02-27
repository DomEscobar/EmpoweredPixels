# AI Agency Guide for EmpoweredPixels

Two agency systems are available for automated development on this codebase.

---

## 1. LangGraph Agency (Recommended)

**Location:** `/root/FutureOfDev/langgraph-agency/`

LangGraph orchestrates a multi-agent pipeline. Each agent is a real opencode session with full tool access (file read/write/edit, grep, glob). The SOUL personality for each role is injected automatically.

### Pipeline

```
triage → architect → hammer → kpi_gate → checker → skeptic → medic → done
```

- **Architect** — Audits the codebase, writes `docs/ARCHITECTURE.md` and `.run/contract.md`
- **Hammer** — Implements the contract (writes code, tests)
- **KPI Gate** — Checks for red-test, green-test, contract.md; retries hammer if missing (max 5)
- **Checker** — Verifies red→green test transition
- **Skeptic** — Quality audit (security, performance, VETO_LOG)
- **Medic** — Fixes lint/build/regression issues

### Usage

```bash
cd /root/FutureOfDev/langgraph-agency

# Ad-hoc task
WORKSPACE=/root/EmpoweredPixels node src/orchestrator.js "Add pagination to the fighters roster API endpoint"

# From a predefined task file (opencode/tasks/<id>.json)
WORKSPACE=/root/EmpoweredPixels node src/orchestrator.js --task bench-001

# Skip KPI gate (faster, no red/green test enforcement)
WORKSPACE=/root/EmpoweredPixels BENCHMARK_MODE=1 node src/orchestrator.js "Fix the league leaderboard sorting"

# Limit hammer retries
WORKSPACE=/root/EmpoweredPixels AGENCY_MAX_HAMMER_RETRIES=2 node src/orchestrator.js "Add weapon stats tooltip"
```

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `WORKSPACE` | `/root/Erp_dev_bench-1` | Target repo for agents to work on |
| `AGENCY_HOME` | `../opencode` | Path to opencode roster/SOULs |
| `AGENCY_MAX_HAMMER_RETRIES` | `5` | Max implementation retries before moving on |
| `AGENCY_RECURSION_LIMIT` | `50` | LangGraph max graph steps |
| `BENCHMARK_MODE` | unset | Skip KPI gate when set |

### Telemetry

Progress is written to `opencode/.run/telemetry_state.json` and pushed to Telegram (if configured in `opencode/config.json`).

---

## 2. opencode Agency (CLI)

**Location:** `/root/FutureOfDev/opencode/`

The original orchestrator. Uses the same roster roles but runs via a CJS orchestrator and a CLI wrapper.

### Usage

```bash
# Link the CLI (one-time)
sudo ln -sf /root/FutureOfDev/opencode/agency.js /usr/local/bin/agency

# Ad-hoc task
cd /root/EmpoweredPixels
agency run "Implement daily reward claim endpoint"

# Formal benchmark task
agency run bench-001

# Check status
agency status

# List active agent roles
agency roster
```

### CLI Commands

| Command | Description |
|---|---|
| `agency run <task_id>` | Run a benchmark task (resets workspace to baseline) |
| `agency run "<prompt>"` | Run any natural language task |
| `agency status` | Show telemetry from last run |
| `agency roster` | List active agent SOULs |
| `agency init` | Bootstrap current directory for governance |

---

## Which to Use

| | LangGraph Agency | opencode Agency |
|---|---|---|
| Orchestration | LangGraph StateGraph with conditional edges | CJS sequential loop |
| Agent execution | Real opencode sessions (`opencode run`) | `dev-unit.cjs` wrapper |
| KPI retry logic | Configurable, graph-based | Linear |
| Telemetry | JSON + Telegram push | JSON + Telegram push |
| Best for | Complex tasks needing retry/branching | Quick ad-hoc runs |

---

## Agent Roles (Shared)

Both systems use the same roster at `/root/FutureOfDev/opencode/roster/`:

| Role | Agent Type | Permissions | Purpose |
|---|---|---|---|
| Architect | `dev-unit` | read, write, edit | Design contract, audit architecture |
| Hammer | `dev-unit` | read, write, edit | Implement the contract |
| Checker | `code-reviewer` | read only | Verify red→green tests |
| Skeptic | `code-reviewer` | read only | Quality/security audit |
| Medic | `dev-unit` | read, write, edit | Fix build/lint/test failures |

## EmpoweredPixels-Specific Notes

- **Backend:** Go 1.22+ with Clean Architecture (domain/usecase/adapter/infra), PostgreSQL, WebSocket
- **Frontend:** Vue 3 + TypeScript + Pinia + Tailwind CSS, Vite build
- **Tests:** Backend `go test ./...`, Frontend `npx vitest`, E2E `npx playwright test`
- The agents will discover the project structure automatically during the triage/brownfield phase

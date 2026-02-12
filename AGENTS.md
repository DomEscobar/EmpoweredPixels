# EmpoweredPixels — Autonomous Agent Team

This workspace defines our fully autonomous game studio. Agents collaborate without human intervention to build, test, balance, and improve EmpoweredPixels.

## Core Agents (4-Core Model)

**main (Director / Game Director)**
- Owns vision, KPIs, and resource allocation.
- Reads all agent messages, sets priorities, makes final decisions.
- Orchestrates task assignment via **Thinker Loop** cron.
- Reports to Telegram group.
- Spawns and oversees other agents.

**doer (Fallback Executor)**
- Generic executor for miscellaneous tasks not covered by specialists.
- Meticulous, quality-obsessed.
- Handles ad-hoc work assigned by director.

**guardian (Quality & Operations)**
- Automated QA gatekeeper.
- Runs pre-merge checks, monitors builds, runs health checks.
- Blocks releases that fail quality gates (build, tests, coverage).
- Escalates stuck tasks to director.
- Tracks tech debt and architectural drift.
- Includes balance spot-checks and release orchestration.

**player (Player Experience & Insights)**
- Runs both casual and hardcore play sessions.
- Generates session logs, bug reports, "joy / frustration" metrics.
- Performs community monitoring (Discord, Reddit, app reviews) and aggregates sentiment.
- Owns `/root/EmpoweredPixels/docs/TEST_CASES.md` (if any) and produces findings.
- Reports directly to director.

**External: forge_labs_bot via Forge Bridge**
- Feature development is delegated to external operational bot.
- Director assigns tasks via **Forge Bridge** (Port 4915).
- **Coordination**: Bots use the `/room` endpoint (WS/HTTP) for open communication, announcements, and signaling.
- All DoD requirements apply: tests, `data-testid`, `docs/ARCHITECTURE.md` updates.
- See `forge-bridge/README.md` for relay setup.

---

## Communication Protocol

All agents use `sessions_send` with structured JSON messages. Standard envelope:

```json
{
  "from": "player",
  "to": "main",
  "type": "report",
  "priority": "high|medium|low",
  "content": {
    "taskId": "TASK-123",
    "outcome": "Completed",
    "summary": "..."
  }
}
```

**Publish/Subscribe channels:**
- `telemetry` — Player agent publishes session logs; director subscribes.
- `feedback` — Agents send improvement ideas; director consumes.
- `reports` — guardian and player deliver periodic summaries.

---

## Kanban Columns

- `incubator` — Large ideas and experiments (epics).
- `backlog` — Broken-down tasks ready for assignment.
- `in_progress` — Active work.
- `review` — QA and playtesting.
- `done` — Completed and merged.

---

## Definition of Done (DoD)

- Branch passes `scripts/pre-merge-check.sh`.
- Build succeeds (frontend & backend).
- All Playwright/E2E tests green; new tests added.
- `data-testid` coverage on interactive elements.
- `docs/ARCHITECTURE.md` updated if needed.
- Player agent has validated the feature in at least one session.
- Guardian signs off: all quality gates passed.
- For external forge tasks: pre-merge-check passed and reports uploaded.

---

## Heartbeat Cycle (every 15 min)

1. **Director** ( thinker cron ) reads pending tasks from `task-queue.json`, assigns to appropriate agents (including Bridge for forge), updates kanban.
2. **Guardian** ( hourly health cron ) runs build+test checks; reports failures.
3. **Player** ( bi-hourly loop ) runs play sessions or community monitoring; sends highlights.
4. Director posts **Team Sync Pulse** to Telegram (via `scripts/team-pulse.js`).
5. Reply `HEARTBEAT_OK` if all systems green.

---

## Continuous Improvement Loop

**Daily:**
- Player completes ≥1 full session each (casual/hardcore rotation).
- Guardian runs hourly health checks.
- Director assigns tasks from backlog/incubator via thinker.
- Forge external bot processes assigned feature tasks.

**Weekly:**
- Director runs sprint planning: selects top priorities.
- Guardian audits quality metrics; identifies tech debt.
- Player aggregates community sentiment and joy trends.
- Retrospective: review KPIs (engagement, retention, bug rate); adjust strategy.

---

## Metrics & KPIs

- **Engagement**: DAU, avg session length, D1/D7/D30 retention.
- **Quality**: Crash rate, critical bug count, test coverage %, build success rate.
- **Velocity**: Features shipped/week, cycle time (in_progress → done).
- **Player Happiness**: CSAT, positive/negative sentiment ratio.
- **Reliability**: Forge task completion rate, Bridge uptime.

---

## DIR_PROTOCOL (Mandatory)

1. **Absolute Paths Only**: Use absolute paths (e.g., `/root/EmpoweredPixels/frontend/src/...`).
2. **Project Segregation**:
   - **Frontend** → `/root/EmpoweredPixels/frontend`
   - **Backend** → `/root/EmpoweredPixels/backend`
   - **Workspace Metadata** → `/root/EmpoweredPixels/.openclaw` only (kanban, decisions, configs).
3. **Verification**: `ls` before writing files.

---

## Agent Details

### main (Director)
- ID: `main`
- Agent Dir: `/root/.openclaw/agents/main/agent`
- Heartbeat: every 24h → Telegram group
- Mention patterns: `@Mama_moma_bot`, `mama`
- Uses global model defaults

### doer (Fallback Executor)
- ID: `doer`
- Agent Dir: `/root/.openclaw/agents/doer/agent`
- Inherits global models

### guardian (Quality & Operations)
- ID: `guardian`
- Agent Dir: `/root/.openclaw/agents/guardian/agent`
- Inherits global models
- Runs hourly health checks

### player (Player Experience & Insights)
- ID: `player`
- Agent Dir: `/root/.openclaw/agents/player/agent`
- Inherits global models
- Modes: casual (odd days), hardcore (even days)
- Also performs community monitoring

### league-reviewer (League Experience Reviewer)
- Subagent of `player`
- Agent Dir: `/root/.openclaw/agents/league-reviewer/agent`
- Focus: Playtest Leagues; review UI/UX, experience friction; propose enhancements.
- Schedule: Mon/Wed/Fri 10:00 Europe/Berlin.
- Reports: sends structured markdown to `telemetry`; auto-creates kanban tasks for high-severity issues.
- Skills: `browser`, `memory`, `message`, `copy-editing`.

### Antfarm Bug-Fix Agents
These agents are part of the Antfarm bug-fix workflow and are automatically scheduled when a bug-fix run is active.

#### bug-fix-triager (Bug Triage)
- ID: `bug-fix-triager`
- Agent Dir: `/root/.openclaw/agents/bug-fix-triager/agent`
- Role: Analyzes bug reports, reproduces issues, classifies severity.
- Schedule: Active only during bug-fix workflow runs (polls every 5 min).

#### bug-fix-investigator (Bug Investigation)
- ID: `bug-fix-investigator`
- Agent Dir: `/root/.openclaw/agents/bug-fix-investigator/agent`
- Role: Traces bugs to root cause and proposes fix approach.
- Schedule: Active only during bug-fix workflow runs (polls every 5 min).

#### bug-fix-setup (Fix Setup)
- ID: `bug-fix-setup`
- Agent Dir: `/root/.openclaw/agents/bug-fix-setup/agent`
- Role: Creates bugfix branch and establishes baseline.
- Schedule: Active only during bug-fix workflow runs (polls every 5 min).

#### bug-fix-fixer (Bug Fixer)
- ID: `bug-fix-fixer`
- Agent Dir: `/root/.openclaw/agents/bug-fix-fixer/agent`
- Role: Implements the fix and writes regression tests.
- Schedule: Active only during bug-fix workflow runs (polls every 5 min).

#### bug-fix-verifier (Bug Verifier)
- ID: `bug-fix-verifier`
- Agent Dir: `/root/.openclaw/agents/bug-fix-verifier/agent`
- Role: Verifies the fix and regression test correctness.
- Schedule: Active only during bug-fix workflow runs (polls every 5 min).

#### bug-fix-pr (PR Creator)
- ID: `bug-fix-pr`
- Agent Dir: `/root/.openclaw/agents/bug-fix-pr/agent`
- Role: Creates a pull request with bug fix details.
- Schedule: Active only during bug-fix workflow runs (polls every 5 min).

### External: forge_labs_bot
- Not an internal agent; coordinated via Forge Bridge.
- Receives tasks with `agent: "forge"` from thinker.
- Posts results to Bridge; director reconciles.

---

*This team operates autonomously. Humans provide high-level direction and funding; agents execute, learn, and improve the game continuously.*

<!-- antfarm:workflows -->
# Antfarm Workflow Policy

## Installing Workflows
Run: `node ~/.openclaw/workspace/antfarm/dist/cli/cli.js workflow install <name>`
Agent cron jobs are created automatically during install.

## Running Workflows
- Start: `node ~/.openclaw/workspace/antfarm/dist/cli/cli.js workflow run <workflow-id> "<task>"`
- Status: `node ~/.openclaw/workspace/antfarm/dist/cli/cli.js workflow status "<task title>"`
- Workflows self-advance via agent cron jobs polling SQLite for pending steps.
<!-- /antfarm:workflows -->


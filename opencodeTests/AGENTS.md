# OpenCode Agency — Project Context

## Build & Run

```bash
# Install dependencies
npm install

# Initialize directories
node cli.js init

# Start the agency (see RUNME.md for full instructions)
# Quick: start director and workers manually or via PM2
node director-advanced.js &
node agents/coder-advanced.js &
node agents/reviewer-advanced.js &
node agents/tester-advanced.js &
node agents/deployer-advanced.js &
```

## Agents

- **Director** – Orchestrates workflows, plans tasks, capability-based routing.
- **Coder** – Executes coding tasks via OpenCode CLI (`opencode --prompt ...`).
- **Reviewer** – Static analysis, security checks, ensures tests exist.
- **Tester** – Runs test suites, validates results.
- **Deployer** – Builds, migrates, deploys, smoke tests.

Each agent is a separate Node.js process with health endpoints and structured logging.

## Configuration

- `agency-config.json` – full configuration (queue, supervisor, agents, workflows).
- Environment:
  - `OPENROUTER_API_KEY` – for OpenCode to use OpenRouter
  - `OPENCODE_MODEL` – model override (default: `openrouter/auto`)
  - `AGENT_WORKDIR` – workspace directory (default: current)
  - `PROMETHEUS_METRICS=1` – enable /metrics endpoint
  - `LOG_LEVEL` – debug|info|warn|error

## Task Queue

SQLite-backed durable queue in `data/agency.db`:
- `tasks` table – pending, processing, completed, failed with retry count.
- `dlq` table – dead-letter after max retries.
- `events` table – audit log.

Agents poll the queue, mark tasks as `processing` on dequeue, and call `complete()` or `fail()`.

## Workflows

Defined in `agency-config.json`. Two built-in workflows:
- `feature`: coder → reviewer → tester → (deployer if approved)
- `bugfix`: coder → reviewer → tester

Start via CLI: `node cli.js workflow feature "Add login"`.

## Webhooks

Run `node webhook-server.js` (port 9091). Routes:
- `POST /webhook/github` – GitHub issue opened → task
- `POST /webhook/task` – manual task submission

Configure secret in `agency.webhooks.secret` for HMAC verification.

## Monitoring

- Director health: `http://localhost:9092/health`
- Agent health: each logs its port on start.
- Metrics: `http://localhost:<port>/metrics` if `PROMETHEUS_METRICS=1`.
- Queue status: `node cli.js status`
- DLQ: `node cli.js dlq`, requeue with `node cli.js requeue <id>`

## OpenCode Integration

Coder agent runs `opencode` as subprocess. Provide `--prompt` with task description, `--model` from config, and `--workdir`. Coder captures stdout/stderr and returns success/failure.

## Communication

Agents do NOT use stdin anymore. They all connect to the SQLite queue. Director enqueues subtasks; workers dequeue independently. Inter-agent data flow is via task payloads and database.

## Supervision

`supervisor.js` can monitor agent health and restart on crash. For production, use PM2 or systemd.

## Development

- `agent-base.js` – base class all agents extend (queue integration, logging, health server, metrics).
- `task-queue.js` – SQLite queue implementation.
- `director-advanced.js` – director agent.
- `agents/*-advanced.js` – worker implementations.

## TODO

- Complete supervisor auto-spawn of all agents
- Capability registry with dynamic discovery
- Cost tracking (token counting per task)
- Multi-host queue (Redis)
- State checkpointing for long tasks
- More sophisticated reviewer (gitleaks, trivy)

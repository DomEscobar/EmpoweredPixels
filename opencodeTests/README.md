# OpenCode Agency — Robust Multi-Agent Autonomous Swarm

Production-grade autonomous coding agency built on OpenCode CLI with durable queues, supervision, observability, and webhook integrations.

## Architecture

```
┌────────────┐
│   Client   │ (CLI / Webhook / GitHub)
└─────┬──────┘
      │ enqueue
      ▼
┌─────────────────┐
│    SQLite Queue │  durable, retries, DLQ
└────────┬────────┘
         │ dequeue
         ▼
┌────────────┐      ┌────────────┐      ┌────────────┐
│   Coder    │─────▶│  Reviewer  │─────▶│   Tester   │
│ (OpenCode) │      │            │      │            │
└────────────┘      └────────────┘      └────────────┘
      │                       │                   │
      └───────────────────────┼───────────────────┘
                              ▼
                       ┌────────────┐
                       │ Deployer   │ (optional)
                       └────────────┘
```

- **Director** – orchestrates workflows, routes tasks based on capabilities.
- **Agents** – coder, reviewer, tester, deployer. Each runs as separate process with health endpoints.
- **Supervisor** – monitors agent health, restarts on failure (Erlang/OTP style).
- **Queue** – SQLite-backed task queue with retries, exponential backoff, dead-letter queue.
- **Webhook Server** – receives GitHub issues, Discord commands, cron triggers and enqueues tasks.
- **Observability** – JSON structured logs, Prometheus metrics, health endpoints.

## Prerequisites

- Node.js 18+
- OpenCode CLI (`opencode`) installed and in PATH
- OpenRouter API key configured (via `OPENROUTER_API_KEY` or OpenCode config)

## Quick Start

```bash
# Install dependencies
npm install

# Create required directories
node cli.js init

# Start the agency (director + supervisor + all agents via supervisor)
# Currently, start agents individually or use PM2 for supervision.
# For simple demo:
 node director-advanced.js &   # director (port 9092)
 node agents/coder-advanced.js &   # coder
 node agents/reviewer-advanced.js &
 node agents/tester-advanced.js &
 node agents/deployer-advanced.js &
```

Or use PM2 to manage all processes (recommended for production).

## Submitting Tasks

### From CLI
```bash
# Simple task (goes through director.plan)
node cli.js task "Refactor utils.js into separate modules with tests"

# Workflow (feature or bugfix)
node cli.js workflow feature "Add user authentication with OAuth2"
```

### From GitHub Webhook
Configure a GitHub repository webhook to `http://localhost:9091/webhook/github` with secret. When an issue is opened, it becomes a task.

### Manual HTTP
```bash
curl -X POST http://localhost:9091/webhook/task \
  -H "Content-Type: application/json" \
  -d '{"type":"director.plan","payload":{"description":"Fix login bug"}}'
```

## Configuration

Edit `agency-config.json` to customize:

- `agency.queue` – SQLite path, retry settings
- `agency.supervisor` – restart policy, health check interval
- `agency.metrics` – Prometheus export port
- `agency.webhooks` – port, secret, routes
- `agents.*` – capabilities, max_concurrent, env vars

Each agent can have environment overrides (e.g., `OPENCODE_MODEL` for coder).

## Monitoring

### Health Endpoints
- Director: `http://localhost:9092/health`
- Other agents: start on random ports; log output shows port.

### Metrics (Prometheus)
Enable by setting `PROMETHEUS_METRICS=1` when starting agents. Then scrape `http://localhost:9090/metrics` (configurable).

### Queue Status
```bash
node cli.js status   # shows task counts by status, DLQ size
node cli.js dlq      # list dead-letter tasks
node cli.js requeue <dlq-id>  # retry a failed task
```

### Logs
Structured JSON logs in `./logs/agency.log` (daily rotation). Use `npm run logs` to tail.

## Task Queue

SQLite `data/agency.db` with tables:
- `tasks` – pending, processing, completed, failed with retry count
- `dlq` – dead-letter tasks after max retries
- `events` – audit log of all state changes

Agents poll the queue every second (`poll_interval_ms`). On success, they call `queue.complete(taskId)`. On failure, `queue.fail(taskId, error)` triggers retry/backoff or DLQ.

## Capability Routing

Agents declare `capabilities` in config. Director matches tasks:
- Direct assignment (`assigned_to` in task) bypasses routing.
- If not assigned, director selects agent advertising required capability (future enhancement; currently hardcoded mapping).

## Workflows

Defined in `agency-config.json` under `workflows`. Each step specifies agent and action, with optional `condition` for auto-advance.

Example `feature` workflow:
1. `coder.implement`
2. `reviewer.review`
3. `tester.validate`
4. `deployer.release` (if review approved)

## Production Considerations

- Use **PM2** or **systemd** to supervise agent processes.
- For multi-host scaling, replace SQLite queue with Redis/NATS.
- Secure webhook endpoints with secrets (GitHub signatures).
- Set up log aggregation (ELK, Grafana Loki) for logs.
- Configure alerts on DLQ size, agent crash loops.
- Use cost tracking: collect token usage from coder and multiply by model pricing.

## File Structure

```
opencodeTests/
├── agency-config.json     # full configuration
├── cli.js                 # control commands
├── supervisor.js          # process manager (optional)
├── director-advanced.js   # orchestrator agent
├── task-queue.js          # SQLite queue implementation
├── agent-base.js          # base class for agents
├── webhook-server.js      # external event ingress
├── agents/
│   ├── coder-advanced.js
│   ├── reviewer-advanced.js
│   ├── tester-advanced.js
│   └── deployer-advanced.js
├── logs/                  # JSON log files
├── data/
│   └── agency.db         # SQLite database
└── README.md
```

## Limitations & Roadmap

- [ ] Full capability registry with dynamic discovery
- [ ] Supervisor implementation (current CLI doesn't auto-restart)
- [ ] Distributed queue (Redis/NATS)
- [ ] Cost tracking dashboard
- [ ] Multi-repo support
- [ ] State checkpointing for long-running tasks

## References

Based on research into AI agent orchestration patterns (Microsoft, IBM ACP, Agentic Design). See memory notes for full report.

---

**Version:** 0.2.0 — Robust Edition

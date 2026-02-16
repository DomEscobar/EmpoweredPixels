# OpenCode Agency — Quick Start Guide

## Setup (First Time)

```bash
cd /root/EmpoweredPixels/opencodeTests

# Install dependencies
npm install

# Create logs and data directories
node cli.js init

# Verify config
cat agency-config.json
```

## Starting the Agency

### Option 1: Simple Demo (Manual Processes)

Open three terminals:

**Terminal 1 - Director:**
```bash
node director-advanced.js
```
Director will start health server on port 9092.

**Terminal 2 - Workers:**
```bash
node agents/coder-advanced.js &
node agents/reviewer-advanced.js &
node agents/tester-advanced.js &
node agents/deployer-advanced.js &
```

### Option 2: Supervisor (Auto-Restart)

The supervisor.js monitors agents and restarts them if they crash. Run:
```bash
node supervisor.js
```
Currently, supervisor is a separate process; you can enhance it to spawn all agents.

### Option 3: PM2 (Production)

```bash
npm install -g pm2
pm2 start director-advanced.js --name director
pm2 start agents/coder-advanced.js --name coder
pm2 start agents/reviewer-advanced.js --name reviewer
pm2 start agents/tester-advanced.js --name tester
pm2 start agents/deployer-advanced.js --name deployer
pm2 save
pm2 startup  # enable on boot
```

## Submitting Work

### Via CLI

```bash
# One-off coding task
node cli.js task "Create a utility function to format dates in YYYY-MM-DD format with tests"

# Feature workflow (sequential: coder → reviewer → tester → deployer if approved)
node cli.js workflow feature "Add user login with JWT tokens"

# Bug fix workflow
node cli.js workflow bugfix "Fix memory leak in match simulation"
```

### Via Webhook (GitHub)

1. Start webhook server:
```bash
PROMETHEUS_METRICS=1 node webhook-server.js &
```

2. Expose localhost:9091 with ngrok or similar:
```bash
ngrok http 9091
```

3. Configure GitHub repository webhook:
   - URL: `https://your-ngrok-subdomain.ngrok.io/webhook/github`
   - Secret: set `agency.webhooks.secret` in config
   - Events: Issues (opened)

When an issue is opened, it becomes a director.plan task.

### Via HTTP POST

```bash
curl -X POST http://localhost:9091/webhook/task \
  -H "Content-Type: application/json" \
  -d '{
    "type": "director.plan",
    "payload": {
      "description": "Add dark mode toggle to settings page",
      "priority": 3
    },
    "priority": 5
  }'
```

## Monitoring

### Check Queue Status
```bash
node cli.js status
# Output: { tasks_by_status: { pending: 3, processing: 1, completed: 10, failed: 0 }, dlq_count: 0, ... }
```

### View Dead-Letter Queue
```bash
node cli.js dlq
```

### Re-queue Failed Task
```bash
node cli.js requeue <dlq-id-from-list>
```

### Health Endpoints

- Director: `curl http://localhost:9092/health`
- Other agents: check logs for assigned port (currently random).

### Metrics

Enable Prometheus metrics for an agent:
```bash
PROMETHEUS_METRICS=1 node agents/coder-advanced.js &
```
Then scrape: `curl http://localhost:<port>/metrics`

### Logs

All agents write JSON logs to `logs/agency.log` (if configured). Use:
```bash
npm run logs   # tail -f logs/agency.log
```

Or rotate logs with external tool (logrotate, etc.).

## Configuration Reference

### agency-config.json

Key sections:

- `agency.queue.type` – "sqlite" (only supported now)
- `agency.queue.path` – SQLite file location (`./data/agency.db`)
- `agency.queue.max_retries` – default 3
- `agency.queue.retry_backoff_ms` – base for exponential backoff
- `agency.supervisor.enabled` – set false to disable supervisor
- `agency.supervisor.health_check_interval_ms` – 10000
- `agency.webhooks.port` – 9091
- `agency.webhooks.secret` – HMAC secret for GitHub signatures
- `agents.<name>.capabilities` – list of strings for routing
- `agents.<name>.max_concurrent` – how many tasks this agent can run in parallel

### Environment Variables

- `OPENROUTER_API_KEY` – required for OpenCode to call models
- `OPENCODE_MODEL` – default model (overridden in coder agent via config)
- `AGENT_WORKDIR` – workspace for all agents (default: current dir)
- `LOG_LEVEL` – debug, info, warn, error (default: info)
- `PROMETHEUS_METRICS` – set to "1" to enable /metrics endpoint
- `SQLITE_DEBUG` – set to see SQLite queries

## Testing the Agency

1. Start the system (director + workers).
2. Submit a simple task:
   ```bash
   node cli.js task "Create a file test.txt containing 'Hello World'"
   ```
3. Watch logs to see coder executing via OpenCode.
4. Check queue status:
   ```bash
   node cli.js status
   ```
5. If coder fails, the task will be retried up to max_retries, then go to DLQ.

## Troubleshooting

- **No tasks being picked up?** Ensure agents are running and their health endpoints respond. Check logs for connection errors.
- **OpenCode not found?** Make sure `opencode` is in PATH. Set `OPENCODE_PATH=/path/to/opencode`.
- **Rate limits?** OpenCode uses OpenRouter; check your quota. Consider adding backoff in coder agent.
- **Tasks stuck in processing?** Agent may have crashed. Check process list. Supervisor should restart; if not, manually restart.
- **DLQ growing?** Inspect `node cli.js dlq` to see errors; fix underlying cause then `node cli.js requeue <id>`.

## Extending

### Add a New Agent

1. Create `agents/myagent-advanced.js`:
   ```js
   const { AgentBase } = require('../agent-base');
   class MyAgent extends AgentBase {
     constructor(config, queue) { super(config, queue); }
     async execute(task) { /* ... */ }
   }
   module.exports = { MyAgent };
   ```
2. Add to `agency-config.json` under `agents`.
3. Add to `supervisor.js` start list (or start manually).
4. Update workflows if needed.

### Change Queue Backend

Swap `task-queue.js` with Redis or NATS implementation. Implement same interface: `enqueue(task)`, `dequeue(agentName)`, `complete(taskId,result)`, `fail(taskId,error)`, `getMetrics()`, etc.

## Performance & Scaling

- **SQLite** is fine for single-host, low-to-medium throughput (< 100 tasks/sec).
- For higher throughput or multi-host, replace with Redis (fast, pub/sub) or NATS (JetStream).
- Each agent can handle `max_concurrent` tasks; tune based on available CPU/GPU and LLM rate limits.
- Use capability routing to distribute load across multiple instances of the same agent type (e.g., 3 coders).

## Security

- Run agents in a sandboxed environment (Docker) if executing untrusted code.
- Never commit API keys; use env vars.
- Enable GitHub webhook secret verification (`agency.webhooks.secret`).
- Audit `logs/agency.log` for suspicious activity.

## License

MIT
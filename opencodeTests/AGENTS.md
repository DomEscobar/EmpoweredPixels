# OpenCode Agency — Project Context

## Build & Run

```bash
# Install dependencies (none beyond Node.js)
npm install

# Start the agency (director + coder)
npm start

# Or use CLI directly
node cli.js start coder
```

## Agents

- **Director**: Orchestrates tasks, assigns to workers
- **Coder**: Uses OpenCode CLI to implement features
- **Reviewer**: Static analysis, security, quality gate
- **Tester**: Generates and runs tests
- **Deployer**: CI/CD and release management

## Configuration

- `agency.json`: Agent definitions and workflows
- Environment:
  - `OPENCODE_PATH`: Path to opencode binary (default: `opencode`)
  - `OPENCODE_MODEL`: Model to use (default: `openrouter/auto`)
  - `AGENT_WORKDIR`: Working directory for agents

## OpenCode Integration

The Coder agent invokes the `opencode` CLI with `--prompt` and `--model` flags to execute coding tasks autonomously.

## Communication

Agents communicate via stdin/stdout pipes (current implementation). Director sends task assignments as `TASK: {...}` messages; agents respond with status logs.

## TODO

- Implement persistent task queue (SQLite)
- Add capability-based routing
- Webhook/task ingestion API
- Agent state persistence and recovery
- Better error handling and retries
- Supervisor to restart crashed agents

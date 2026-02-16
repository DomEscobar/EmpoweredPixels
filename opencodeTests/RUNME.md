# OpenCode Agency — Quick Start

## Prerequisites

- Node.js 18+
- OpenCode CLI installed and in PATH
- OpenRouter API key (or other provider) configured in OpenCode

## Installation

```bash
cd /root/EmpoweredPixels/opencodeTests
npm install  # not strictly needed but for completeness
```

## Starting the Agency

```bash
# Start director and coder agent
node cli.js start coder

# Or start all agents manually in separate terminals:
# Terminal 1: node director.js
# Terminal 2: node agents/coder.js
# Terminal 3: node agents/reviewer.js
# etc.
```

## Submitting Tasks

Currently, tasks are submitted by sending to director's stdin:

```bash
# In director terminal, type:
task: Create a calculator class with add/subtract/multiply/divide
```

Or programmatically:

```js
const { spawn } = require('child_process');
const director = spawn('node', ['director.js'], { stdio: 'pipe' });
director.stdin.write('task: Your task description here\n');
```

## How It Works

1. Director receives tasks and assigns to coder
2. Coder invokes `opencode --prompt "<task>" --model openrouter/auto`
3. OpenCode executes the task in the workspace
4. Reviewer runs static analysis (hook into git diff, linters)
5. Tester generates and runs tests
6. Deployer handles release if configured

## Configuration

Set these environment variables:

- `OPENCODE_PATH`: Path to opencode binary (default: `opencode`)
- `OPENCODE_MODEL`: Model string for OpenCode (default: `openrouter/auto`)
- `AGENT_WORKDIR`: Working directory (default: current)

Example:
```bash
export OPENCODE_MODEL="anthropic/claude-sonnet-4"
export AGENT_WORKDIR="/root/EmpoweredPixels/opencodeTests"
```

## Current Limitations

- Task queue is in-memory only (no persistence yet)
- Director uses simple round-robin assignment
- No web API for task submission
- Agents communicate via stdin; needs proper message bus
- No agent supervision/restart on crash (yet)

## Roadmap

- [ ] SQLite task queue with retry/backoff
- [ ] HTTP API for task submission
- [ ] Multi-repo support
- [ ] Agent state persistence
- [ ] Supervisor process (like PM2)
- [ ] Webhook integrations (GitHub, Discord)
- [ ] Cost tracking per agent/task

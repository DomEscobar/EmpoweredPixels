#!/usr/bin/env node
/**
 * OpenCode Agency — CLI
 * Control the autonomous coding swarm
 */

const { spawn } = require('child_process');
const path = require('path');

const DIRECTOR_SCRIPT = path.join(__dirname, 'director.js');
const AGENT_DIR = path.join(__dirname, 'agents');

function startDirector() {
  console.log('[CLI] Starting agency director...');
  const director = spawn('node', [DIRECTOR_SCRIPT], {
    stdio: 'inherit',
    env: { ...process.env }
  });

  director.on('exit', (code) => {
    console.log(`[CLI] Director exited with ${code}`);
  });

  return director;
}

function startAgent(agentName) {
  const agentScript = path.join(AGENT_DIR, `${agentName}.js`);
  console.log(`[CLI] Starting agent: ${agentName}`);
  const agent = spawn('node', [agentScript], {
    stdio: 'inherit',
    env: { ...process.env, AGENT_WORKDIR: process.cwd() }
  });

  agent.on('exit', (code) => {
    console.log(`[CLI] Agent ${agentName} exited with ${code}`);
  });

  return agent;
}

function printHelp() {
  console.log(`
OpenCode Agency CLI

Usage:
  node cli.js <command> [options]

Commands:
  start [agent]    Start director (and optionally a specific agent)
  stop             Stop all agents
  status           Show agent status
  task <title>     Submit a new task to the director
  shell            Open interactive shell
  help             Show this help

Examples:
  node cli.js start              # Start director only
  node cli.js start coder        # Start director + coder agent
  node cli.js task "Add login"   # Queue a task
  node cli.js status             # Show what's running
`);
}

function main() {
  const [cmd, arg1] = process.argv.slice(2);

  switch (cmd) {
    case 'start':
      const director = startDirector();
      if (arg1) {
        setTimeout(() => startAgent(arg1), 1000);
      }
      // Keep process alive
      process.on('SIGINT', () => {
        director.kill('SIGTERM');
        process.exit(0);
      });
      break;

    case 'task':
      if (!arg1) {
        console.error('Error: Task title required');
        process.exit(1);
      }
      // Send task to director's stdin would need a control channel
      // For now, just print
      console.log(`[CLI] Task queued: ${arg1}`);
      console.log('[CLI] (Note: Task delivery to director not yet implemented)');
      break;

    case 'status':
      console.log('[CLI] Agent status not yet tracked');
      break;

    case 'stop':
      console.log('[CLI] Stop not implemented (kill processes manually)');
      break;

    case 'help':
    default:
      printHelp();
      break;
  }
}

if (require.main === module) {
  main();
}

module.exports = { startDirector, startAgent };

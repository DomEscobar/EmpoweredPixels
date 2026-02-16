#!/usr/bin/env node
/**
 * OpenCode Agency — Director
 * Orchestrates multi-agent autonomous coding swarm
 */

const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');

const CONFIG_PATH = path.join(__dirname, 'agency.json');
const STATE_DIR = path.join(__dirname, '.agency-state');

// Ensure state directory exists
fs.mkdirSync(STATE_DIR, { recursive: true });

class Agent {
  constructor(name, config) {
    this.name = name;
    this.role = config.role;
    this.capabilities = config.capabilities;
    this.executor = config.executor || null;
    this.process = null;
    this.status = 'idle';
  }

  start() {
    if (this.process) return;

    const cmd = this.executor || 'crush';
    const args = ['agent-turn', `--agent=${this.name}`, '--message=Agent starting'];

    this.process = spawn(cmd, args, {
      stdio: 'inherit',
      env: { ...process.env, AGENT_ROLE: this.role }
    });

    this.process.on('exit', (code) => {
      console.log(`[Director] Agent ${this.name} exited with ${code}`);
      this.process = null;
      this.status = 'stopped';
    });

    this.status = 'running';
    console.log(`[Director] Started agent: ${this.name} (${this.role})`);
  }

  send(message) {
    if (!this.process) {
      console.warn(`[Director] Agent ${this.name} not running`);
      return;
    }
    this.process.stdin.write(message + '\n');
  }

  stop() {
    if (this.process) {
      this.process.kill('SIGTERM');
      this.process = null;
      this.status = 'stopped';
    }
  }
}

class Director {
  constructor(configPath) {
    this.config = JSON.parse(fs.readFileSync(configPath, 'utf8'));
    this.agents = {};
    this.queue = [];
    this.activeTask = null;
  }

  async initialize() {
    console.log('[Director] Initializing agency...');

    // Create agent instances
    for (const [name, agentConfig] of Object.entries(this.config.agents)) {
      this.agents[name] = new Agent(name, agentConfig);
    }

    console.log(`[Director] Loaded ${Object.keys(this.agents).length} agents`);
  }

  startAll() {
    console.log('[Director] Starting all agents...');
    for (const agent of Object.values(this.agents)) {
      agent.start();
    }
  }

  stopAll() {
    console.log('[Director] Stopping all agents...');
    for (const agent of Object.values(this.agents)) {
      agent.stop();
    }
  }

  assignTask(task) {
    // Simple round-robin for now; will enhance with capability matching
    const agentNames = Object.keys(this.agents);
    const coder = this.agents['coder'];
    if (coder && coder.status === 'running') {
      coder.send(`TASK: ${JSON.stringify(task)}`);
      this.activeTask = task;
      console.log(`[Director] Assigned task to coder: ${task.title}`);
    } else {
      console.error('[Director] No available coder agent');
    }
  }

  monitor() {
    // Periodic health check
    setInterval(() => {
      for (const [name, agent] of Object.entries(this.agents)) {
        if (agent.status !== 'running') {
          console.log(`[Director] Restarting agent ${name}`);
          agent.start();
        }
      }
    }, 30000);
  }
}

// Main
if (require.main === module) {
  const director = new Director(CONFIG_PATH);

  director.initialize().then(() => {
    director.startAll();
    director.monitor();

    // CLI input loop for manual task assignment
    process.stdin.resume();
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', (data) => {
      const cmd = data.trim();
      if (cmd.startsWith('task:')) {
        const task = { title: cmd.slice(5), created: Date.now() };
        director.assignTask(task);
      } else if (cmd === 'status') {
        console.log('[Director] Agent status:');
        for (const [name, agent] of Object.entries(director.agents)) {
          console.log(`  ${name}: ${agent.status}`);
        }
      } else if (cmd === 'stop') {
        director.stopAll();
        process.exit(0);
      }
    });

  }).catch((err) => {
    console.error('[Director] Failed to initialize:', err);
    process.exit(1);
  });
}

module.exports = { Director, Agent };

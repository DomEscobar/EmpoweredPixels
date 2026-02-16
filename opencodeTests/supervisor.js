#!/usr/bin/env node
/**
 * OpenCode Agency — Supervisor
 * Monitors agent health and restarts crashed processes
 */

const { spawn } = require('child_process');
const path = require('path');
const http = require('http');

class Supervisor {
  constructor(config, queue) {
    this.config = config;
    this.queue = queue;
    this.agents = new Map(); // name -> {proc, config, restartCount, lastRestart}
    this.logger = require('pino')({ level: process.env.LOG_LEVEL || 'info' });
    this.healthCheckInterval = null;
    this.running = false;
  }

  async start() {
    this.logger.info('Starting supervisor');
    this.running = true;

    // Start all agents defined in config
    for (const [name, agentConfig] of Object.entries(this.config.agents)) {
      await this.startAgent(name, agentConfig);
    }

    // Periodic health checks
    this.healthCheckInterval = setInterval(() => this.checkHealth(), this.config.supervisor.health_check_interval_ms);
  }

  async startAgent(name, config) {
    if (this.agents.has(name)) {
      throw new Error(`Agent ${name} already running`);
    }

    const scriptPath = path.join(__dirname, 'agents', `${name}-advanced.js`);
    const env = { ...process.env, ...config.env, AGENT_NAME: name, AGENT_WORKDIR: process.cwd() };

    const proc = spawn('node', [scriptPath], {
      stdio: 'pipe',
      env,
      detached: false
    });

    proc.stdout.on('data', (data) => this.logger.info({ agent: name }, data.toString().trim()));
    proc.stderr.on('data', (data) => this.logger.error({ agent: name }, data.toString().trim()));

    proc.on('exit', (code, signal) => {
      this.logger.warn({ agent: name, code, signal }, 'Agent exited');
      this.agents.delete(name);
      if (this.running && config.restart_on_failure) {
        this.restartAgent(name, config);
      }
    });

    this.agents.set(name, {
      proc,
      config,
      restartCount: 0,
      lastRestart: null
    });

    this.logger.info({ agent: name }, 'Agent started');
    // Give agent time to start health server
    await new Promise(resolve => setTimeout(resolve, 2000));
  }

  async restartAgent(name, config) {
    const agent = this.agents.get(name);
    const maxAttempts = this.config.supervisor.max_restart_attempts;
    const backoffMs = this.config.supervisor.restart_delay_ms;

    if (agent && agent.restartCount >= maxAttempts) {
      this.logger.error({ agent: name, attempts: agent.restartCount }, 'Max restart attempts reached, giving up');
      return;
    }

    // Exponential backoff if multiple restarts
    const delay = backoffMs * Math.pow(this.config.supervisor.backoff_multiplier, agent ? agent.restartCount : 0);
    this.logger.info({ agent: name, delayMs: delay, attempt: (agent?.restartCount || 0) + 1 }, 'Restarting agent');

    await new Promise(resolve => setTimeout(resolve, delay));

    try {
      await this.startAgent(name, config);
      if (agent) agent.restartCount = (agent.restartCount || 0) + 1;
      else this.agents.get(name).restartCount = 1;
      this.agents.get(name).lastRestart = Date.now();
    } catch (err) {
      this.logger.error({ agent: name, error: err.message }, 'Failed to restart agent');
    }
  }

  async checkHealth() {
    for (const [name, agent] of this.agents) {
      try {
        const healthUrl = `http://localhost:${agent.healthPort || 0}${agent.config.health_path}`;
        // We'll store healthPort when agent starts; for now just ping if we have it
        // In practice, agents should register their health port somewhere (shared memory/file)
        // For simplicity, we'll assume agents log their health port; we could also call agent's health endpoint if known
        // TODO: Track health ports in supervisor state
      } catch (err) {
        // ignore; agent might not have health server yet
      }
    }
  }

  async stop() {
    this.running = false;
    if (this.healthCheckInterval) clearInterval(this.healthCheckInterval);
    for (const [name, agent] of this.agents) {
      this.logger.info({ agent: name }, 'Stopping agent');
      agent.proc.kill('SIGTERM');
    }
  }
}

module.exports = { Supervisor };

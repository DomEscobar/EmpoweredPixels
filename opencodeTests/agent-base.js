#!/usr/bin/env node
/**
 * OpenCode Agency — Agent Base Class
 * Common functionality: queue consumption, health, metrics, logging
 */

const http = require('http');
const pino = require('pino');
const { v4: uuidv4 } = require('uuid');
const path = require('path');

class AgentBase {
  constructor(config, queue) {
    this.name = config.name;
    this.role = config.role;
    this.capabilities = config.capabilities || [];
    this.maxConcurrent = config.max_concurrent || 1;
    this.executor = config.executor;
    this.env = config.env || {};
    this.healthPath = config.health_path || '/health';
    this.queue = queue;

    // Setup logger
    this.logger = pino({
      level: process.env.LOG_LEVEL || 'info',
      transport: {
        target: 'pino-pretty',
        options: { destination: process.env.LOG_FILE || `./logs/${this.name}.log` }
      },
      timestamp: pino.stdTimeFunctions.isoTime
    });

    this.running = false;
    this.currentTasks = new Map(); // taskId -> startTime
    this.metrics = {
      tasks_started: 0,
      tasks_completed: 0,
      tasks_failed: 0,
      total_duration_ms: 0
    };
  }

  async start() {
    this.running = true;
    this.logger.info(`Starting agent ${this.name} (${this.role})`);
    // Start polling queue in background
    this.pollInterval = setInterval(() => this.processQueue(), this.queue.pollIntervalMs);
  }

  async stop() {
    this.running = false;
    if (this.pollInterval) clearInterval(this.pollInterval);
    this.logger.info(`Stopped agent ${this.name}`);
  }

  async processQueue() {
    if (!this.running) return;
    if (Object.keys(this.currentTasks).length >= this.maxConcurrent) return;

    const task = this.queue.dequeue(this.name);
    if (!task) return;

    this.logger.info({ taskId: task.id, taskType: task.type }, 'Dequeued task');
    this.currentTasks.set(task.id, Date.now());
    this.metrics.tasks_started++;

    try {
      const result = await this.execute(task);
      this.queue.complete(task.id, result);
      const duration = Date.now() - this.currentTasks.get(task.id);
      this.metrics.tasks_completed++;
      this.metrics.total_duration_ms += duration;
      this.logger.info(
        { taskId: task.id, durationMs: duration, result },
        'Task completed'
      );
    } catch (error) {
      this.logger.error({ taskId: task.id, error: error.message }, 'Task failed');
      this.metrics.tasks_failed++;
      this.queue.fail(task.id, error.message, this.name);
    } finally {
      this.currentTasks.delete(task.id);
    }
  }

  /**
   * Override in subclass
   */
  async execute(task) {
    throw new Error('execute() must be implemented by subclass');
  }

  /**
   * Health check endpoint server
   */
  startHealthServer(port = 0) {
    // auto-assign port if 0
    const server = http.createServer(async (req, res) => {
      if (req.url === this.healthPath) {
        if (this.running) {
          res.writeHead(200, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({
            status: 'ok',
            agent: this.name,
            role: this.role,
            uptime: process.uptime(),
            current_tasks: Object.keys(this.currentTasks).length,
            metrics: this.metrics
          }));
        } else {
          res.writeHead(503);
          res.end('Agent not running');
        }
      } else if (req.url === '/metrics' && process.env.PROMETHEUS_METRICS) {
        // Simple Prometheus text format
        const metrics = this.formatPrometheus();
        res.writeHead(200, { 'Content-Type': 'text/plain; version=0.0.4' });
        res.end(metrics);
      } else {
        res.writeHead(404);
        res.end();
      }
    });

    server.listen(port, 'localhost').on('listening', () => {
      const addr = server.address();
      this.logger.info({ port: addr.port }, 'Health server listening');
      this.healthPort = addr.port;
    });
    this.healthServer = server;
  }

  stopHealthServer() {
    if (this.healthServer) this.healthServer.close();
  }

  formatPrometheus() {
    const lines = [];
    lines.push(`# HELP agency_tasks_started Total tasks started`);
    lines.push(`# TYPE agency_tasks_started counter`);
    lines.push(`agency_tasks_started{agent="${this.name}"} ${this.metrics.tasks_started}`);
    lines.push(`# HELP agency_tasks_completed Total tasks completed`);
    lines.push(`# TYPE agency_tasks_completed counter`);
    lines.push(`agency_tasks_completed{agent="${this.name}"} ${this.metrics.tasks_completed}`);
    lines.push(`# HELP agency_tasks_failed Total tasks failed`);
    lines.push(`# TYPE agency_tasks_failed counter`);
    lines.push(`agency_tasks_failed{agent="${this.name}"} ${this.metrics.tasks_failed}`);
    lines.push(`# HELP agency_task_duration_ms Total task duration in ms`);
    lines.push(`# TYPE agency_task_duration_ms counter`);
    lines.push(`agency_task_duration_ms{agent="${this.name}"} ${this.metrics.total_duration_ms}`);
    lines.push(`# HELP agency_current_tasks Number of currently running tasks`);
    lines.push(`# TYPE agency_current_tasks gauge`);
    lines.push(`agency_current_tasks{agent="${this.name}"} ${Object.keys(this.currentTasks).length}`);
    return lines.join('\n') + '\n';
  }
}

module.exports = { AgentBase };

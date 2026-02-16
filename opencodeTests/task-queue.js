#!/usr/bin/env node
/**
 * OpenCode Agency — SQLite Task Queue (with in-memory fallback)
 * Durable task queue with retries, backoff, and dead-letter support
 */

let Database;
try {
  Database = require('better-sqlite3');
} catch (e) {
  console.warn('[Queue] better-sqlite3 not available, using in-memory fallback');
  Database = null;
}
const { v4: uuidv4 } = require('uuid');
const path = require('path');

// In-memory fallback implementation
class MemoryQueue {
  constructor() {
    this.tasks = new Map();
    this.dlq = [];
    this.events = [];
  }
  createTables() {} // noop
  enqueue(task) {
    const id = task.id || uuidv4();
    this.tasks.set(id, {
      id,
      type: task.type,
      payload: JSON.stringify(task.payload),
      priority: task.priority || 0,
      status: 'pending',
      assigned_to: null,
      locked_by: null,
      locked_at: null,
      retry_count: 0,
      max_retries: task.max_retries || 3,
      last_error: null,
      scheduled_at: task.scheduled_at || null,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    });
    this.events.push({ event_type: 'task_enqueued', task_id: id, timestamp: new Date().toISOString() });
    return id;
  }
  dequeue(agentName) {
    // Find a pending task not locked, highest priority, earliest scheduled
    let candidate = null;
    for (const [id, task] of this.tasks) {
      if (task.status === 'pending' && !task.locked_by && (!task.scheduled_at || new Date(task.scheduled_at) <= new Date())) {
        if (!candidate || task.priority > candidate.priority || (task.priority === candidate.priority && task.scheduled_at < candidate.scheduled_at)) {
          candidate = { id, ...task };
        }
      }
    }
    if (!candidate) return null;
    // Lock it
    const now = new Date().toISOString();
    candidate.locked_by = agentName;
    candidate.locked_at = now;
    candidate.status = 'processing';
    candidate.updated_at = now;
    this.tasks.set(candidate.id, candidate);
    return {
      id: candidate.id,
      type: candidate.type,
      payload: JSON.parse(candidate.payload)
    };
  }
  complete(taskId, result = {}) {
    const task = this.tasks.get(taskId);
    if (task) {
      task.status = 'completed';
      task.updated_at = new Date().toISOString();
      this.tasks.set(taskId, task);
      this.events.push({ event_type: 'task_completed', task_id: taskId, payload: result, timestamp: new Date().toISOString() });
    }
  }
  fail(taskId, error, agentName) {
    const task = this.tasks.get(taskId);
    if (!task) return;
    const newRetry = (task.retry_count || 0) + 1;
    if (newRetry >= task.max_retries) {
      this.moveToDlq(taskId, error, newRetry);
    } else {
      const backoffMs = 2000 * Math.pow(2, newRetry);
      const scheduledAt = new Date(Date.now() + backoffMs).toISOString();
      task.status = 'pending';
      task.locked_by = null;
      task.locked_at = null;
      task.retry_count = newRetry;
      task.last_error = error;
      task.scheduled_at = scheduledAt;
      task.updated_at = new Date().toISOString();
      this.tasks.set(taskId, task);
      this.events.push({ event_type: 'task_retry', task_id: taskId, agent: agentName, payload: { retry: newRetry, error, scheduled_at: scheduledAt }, timestamp: new Date().toISOString() });
    }
  }
  moveToDlq(taskId, error, retryCount) {
    const task = this.tasks.get(taskId);
    if (!task) return;
    this.dlq.push({
      id: uuidv4(),
      original_task_id: taskId,
      type: task.type,
      payload: task.payload,
      reason: error,
      failed_at: new Date().toISOString(),
      retry_count: retryCount
    });
    this.tasks.delete(taskId);
    this.events.push({ event_type: 'task_dlq', task_id: taskId, payload: { error, retry_count: retryCount }, timestamp: new Date().toISOString() });
  }
  getTask(taskId) {
    return this.tasks.get(taskId) || null;
  }
  listDlq() {
    return this.dlq;
  }
  requeueDlq(dlqId) {
    const idx = this.dlq.findIndex(item => item.id === dlqId);
    if (idx === -1) throw new Error('DLQ item not found');
    const item = this.dlq[idx];
    this.enqueue({
      id: item.original_task_id,
      type: item.type,
      payload: JSON.parse(item.payload)
    });
    this.dlq.splice(idx, 1);
    this.events.push({ event_type: 'task_requeued', task_id: item.original_task_id, payload: { from_dlq: dlqId }, timestamp: new Date().toISOString() });
  }
  getMetrics() {
    const statusCounts = { pending: 0, processing: 0, completed: 0, failed: 0 };
    for (const task of this.tasks.values()) {
      statusCounts[task.status] = (statusCounts[task.status] || 0) + 1;
    }
    const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000).toISOString();
    const eventsLastHour = this.events.filter(e => new Date(e.timestamp) >= new Date(oneHourAgo)).length;
    return { tasks_by_status: statusCounts, dlq_count: this.dlq.length, events_last_hour: eventsLastHour };
  }
  close() {}
}

class TaskQueue {
  constructor(config) {
    this.config = config;
    this.useMemory = false;
    this.db = null;
    this.pollIntervalMs = config.poll_interval_ms || 1000;
  }

  async initialize() {
    if (!Database) {
      console.log('[Queue] Using in-memory queue (better-sqlite3 not available)');
      this.queue = new MemoryQueue();
      this.useMemory = true;
      return;
    }
    try {
      this.db = new Database(this.config.path, { verbose: process.env.SQLITE_DEBUG ? console.log : undefined });
      this.db.pragma('journal_mode = WAL');
      this.db.pragma('synchronous = NORMAL');
      this.createTables();
      console.log('[Queue] Initialized SQLite at', this.config.path);
    } catch (e) {
      console.warn('[Queue] SQLite init failed, falling back to memory:', e.message);
      this.queue = new MemoryQueue();
      this.useMemory = true;
    }
  }

  createTables() {
    if (this.useMemory) return;
    // ... same as before, using this.db.exec ...
    this.db.exec(`
      CREATE TABLE IF NOT EXISTS ${this.config.table} (
        id TEXT PRIMARY KEY,
        type TEXT NOT NULL,
        payload TEXT NOT NULL,
        priority INTEGER DEFAULT 0,
        status TEXT DEFAULT 'pending',
        assigned_to TEXT,
        locked_by TEXT,
        locked_at TEXT,
        retry_count INTEGER DEFAULT 0,
        max_retries INTEGER DEFAULT ${this.config.max_retries || 3},
        last_error TEXT,
        scheduled_at TEXT DEFAULT (datetime('now')),
        created_at TEXT DEFAULT (datetime('now')),
        updated_at TEXT DEFAULT (datetime('now'))
      );
      CREATE INDEX IF NOT EXISTS idx_status_priority ON ${this.config.table} (status, priority DESC, scheduled_at);
      CREATE INDEX IF NOT EXISTS idx_locked_by ON ${this.config.table} (locked_by);
    `);
    this.db.exec(`
      CREATE TABLE IF NOT EXISTS ${this.config.dead_letter_table || 'dlq'} (
        id TEXT PRIMARY KEY,
        original_task_id TEXT,
        type TEXT NOT NULL,
        payload TEXT NOT NULL,
        reason TEXT,
        failed_at TEXT DEFAULT (datetime('now')),
        retry_count INTEGER DEFAULT 0
      );
    `);
    this.db.exec(`
      CREATE TABLE IF NOT EXISTS events (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        event_type TEXT NOT NULL,
        task_id TEXT,
        agent TEXT,
        payload TEXT,
        timestamp TEXT DEFAULT (datetime('now'))
      );
      CREATE INDEX IF NOT EXISTS idx_event_timestamp ON events (timestamp);
      CREATE INDEX IF NOT EXISTS idx_event_task_id ON events (task_id);
    `);
  }

  enqueue(task) {
    if (this.useMemory) {
      return this.queue.enqueue(task);
    }
    const id = task.id || uuidv4();
    const scheduledAt = task.scheduled_at || null;
    const stmt = this.db.prepare(`
      INSERT INTO ${this.config.table} (id, type, payload, priority, scheduled_at)
      VALUES (@id, @type, @payload, @priority, @scheduled_at)
    `);
    stmt.run({
      id,
      type: task.type,
      payload: JSON.stringify(task.payload),
      priority: task.priority || 0,
      scheduled_at: scheduledAt
    });
    this.logEvent('task_enqueued', id, null, { type: task.type, priority: task.priority });
    return id;
  }

  dequeue(agentName) {
    if (this.useMemory) {
      return this.queue.dequeue(agentName);
    }
    const now = new Date().toISOString();
    return this.db.transaction(() => {
      const select = this.db.prepare(`
        SELECT id, type, payload FROM ${this.config.table}
        WHERE status = 'pending' AND locked_by IS NULL
        AND (scheduled_at IS NULL OR scheduled_at <= datetime('now'))
        ORDER BY priority DESC, scheduled_at ASC
        LIMIT 1
      `);
      const task = select.get();
      if (!task) return null;
      const update = this.db.prepare(`
        UPDATE ${this.config.table}
        SET locked_by = ?, locked_at = ?, status = 'processing', updated_at = datetime('now')
        WHERE id = ? AND locked_by IS NULL
      `);
      const result = update.run(agentName, now, task.id);
      if (result.changes === 0) return null;
      return {
        id: task.id,
        type: task.type,
        payload: JSON.parse(task.payload)
      };
    })();
  }

  complete(taskId, result = {}) {
    if (this.useMemory) {
      this.queue.complete(taskId, result);
      return;
    }
    const stmt = this.db.prepare(`
      UPDATE ${this.config.table}
      SET status = 'completed', updated_at = datetime('now')
      WHERE id = ?
    `);
    stmt.run(taskId);
    this.logEvent('task_completed', taskId, null, result);
  }

  fail(taskId, error, agentName) {
    if (this.useMemory) {
      this.queue.fail(taskId, error, agentName);
      return;
    }
    const get = this.db.prepare(`SELECT retry_count, max_retries FROM ${this.config.table} WHERE id = ?`);
    const task = get.get(taskId);
    if (!task) return;
    const newRetry = task.retry_count + 1;
    if (newRetry >= task.max_retries) {
      this.moveToDlq(taskId, error, task.retry_count);
    } else {
      const backoffMs = this.config.retry_backoff_ms || 2000;
      const scheduledAt = new Date(Date.now() + backoffMs * Math.pow(2, newRetry)).toISOString();
      const stmt = this.db.prepare(`
        UPDATE ${this.config.table}
        SET status = 'pending', locked_by = NULL, locked_at = NULL,
            retry_count = ?, last_error = ?, scheduled_at = ?, updated_at = datetime('now')
        WHERE id = ?
      `);
      stmt.run(newRetry, error, scheduledAt, taskId);
      this.logEvent('task_retry', taskId, agentName, { retry: newRetry, error, scheduled_at: scheduledAt });
    }
  }

  moveToDlq(taskId, error, retryCount) {
    if (this.useMemory) {
      this.queue.moveToDlq(taskId, error, retryCount);
      return;
    }
    const get = this.db.prepare(`SELECT type, payload FROM ${this.config.table} WHERE id = ?`);
    const task = get.get(taskId);
    if (!task) return;
    const insertDlq = this.db.prepare(`
      INSERT INTO ${this.config.dead_letter_table || 'dlq'} (id, original_task_id, type, payload, reason, failed_at, retry_count)
      VALUES (@id, @original_id, @type, @payload, @reason, datetime('now'), @retry_count)
    `);
    insertDlq.run({
      id: uuidv4(),
      original_id: taskId,
      type: task.type,
      payload: task.payload,
      reason: error,
      retry_count: retryCount
    });
    this.db.prepare(`DELETE FROM ${this.config.table} WHERE id = ?`).run(taskId);
    this.logEvent('task_dlq', taskId, null, { error, retry_count: retryCount });
  }

  logEvent(type, taskId, agent, payload = {}) {
    if (this.useMemory) {
      this.queue.events.push({ event_type: type, task_id: taskId, agent, payload, timestamp: new Date().toISOString() });
      return;
    }
    const stmt = this.db.prepare(`
      INSERT INTO events (event_type, task_id, agent, payload)
      VALUES (?, ?, ?, ?)
    `);
    stmt.run(type, taskId, agent, JSON.stringify(payload));
  }

  getTask(taskId) {
    return this.useMemory ? this.queue.getTask(taskId) : (this.db.prepare(`SELECT * FROM ${this.config.table} WHERE id = ?`).get(taskId));
  }

  listDlq() {
    return this.useMemory ? this.queue.listDlq() : this.db.prepare(`SELECT * FROM ${this.config.dead_letter_table || 'dlq'} ORDER BY failed_at DESC`).all();
  }

  requeueDlq(dlqId) {
    if (this.useMemory) {
      this.queue.requeueDlq(dlqId);
      return;
    }
    const get = this.db.prepare(`SELECT * FROM ${this.config.dead_letter_table || 'dlq'} WHERE id = ?`);
    const item = get.get(dlqId);
    if (!item) throw new Error('DLQ item not found');
    this.enqueue({
      id: item.original_task_id,
      type: item.type,
      payload: JSON.parse(item.payload)
    });
    this.db.prepare(`DELETE FROM ${this.config.dead_letter_table || 'dlq'} WHERE id = ?`).run(dlqId);
    this.logEvent('task_requeued', item.original_task_id, 'admin', { from_dlq: dlqId });
  }

  getMetrics() {
    return this.useMemory ? this.queue.getMetrics() : (() => {
      const metrics = {};
      const statusCounts = this.db.prepare(`
        SELECT status, COUNT(*) as count FROM ${this.config.table} GROUP BY status
      `).all();
      metrics.tasks_by_status = statusCounts.reduce((acc, row) => ({ ...acc, [row.status]: row.count }), {});
      metrics.dlq_count = this.db.prepare(`SELECT COUNT(*) as count FROM ${this.config.dead_letter_table || 'dlq'}`).get().count;
      metrics.events_last_hour = this.db.prepare(`
        SELECT COUNT(*) as count FROM events
        WHERE timestamp >= datetime('now', '-1 hour')
      `).get().count;
      return metrics;
    })();
  }

  close() {
    if (this.db) {
      this.db.close();
      this.db = null;
    }
  }
}

module.exports = { TaskQueue };

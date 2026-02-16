#!/usr/bin/env node
/**
 * OpenCode Agency — SQLite Task Queue
 * Durable task queue with retries, backoff, and dead-letter support
 */

const Database = require('better-sqlite3');
const { v4: uuidv4 } = require('uuid');
const path = require('path');

class TaskQueue {
  constructor(config) {
    this.dbPath = config.path;
    this.table = config.table;
    this.dlqTable = config.dead_letter_table || 'dlq';
    this.maxRetries = config.max_retries || 3;
    this.retryBackoffMs = config.retry_backoff_ms || 2000;
    this.pollIntervalMs = config.poll_interval_ms || 1000;
    this.db = null;
  }

  async initialize() {
    this.db = new Database(this.dbPath, { verbose: process.env.SQLITE_DEBUG ? console.log : undefined });
    this.db.pragma('journal_mode = WAL');
    this.db.pragma('synchronous = NORMAL');
    this.createTables();
    console.log('[Queue] Initialized with SQLite at', this.dbPath);
  }

  createTables() {
    // Main tasks table
    this.db.exec(`
      CREATE TABLE IF NOT EXISTS ${this.table} (
        id TEXT PRIMARY KEY,
        type TEXT NOT NULL,
        payload TEXT NOT NULL,
        priority INTEGER DEFAULT 0,
        status TEXT DEFAULT 'pending',
        assigned_to TEXT,
        locked_by TEXT,
        locked_at TEXT,
        retry_count INTEGER DEFAULT 0,
        max_retries INTEGER DEFAULT ${this.maxRetries},
        last_error TEXT,
        scheduled_at TEXT DEFAULT (datetime('now')),
        created_at TEXT DEFAULT (datetime('now')),
        updated_at TEXT DEFAULT (datetime('now'))
      );
      CREATE INDEX IF NOT EXISTS idx_status_priority ON ${this.table} (status, priority DESC, scheduled_at);
      CREATE INDEX IF NOT EXISTS idx_locked_by ON ${this.table} (locked_by);
    `);

    // Dead-letter queue
    this.db.exec(`
      CREATE TABLE IF NOT EXISTS ${this.dlqTable} (
        id TEXT PRIMARY KEY,
        original_task_id TEXT,
        type TEXT NOT NULL,
        payload TEXT NOT NULL,
        reason TEXT,
        failed_at TEXT DEFAULT (datetime('now')),
        retry_count INTEGER DEFAULT 0
      );
    `);

    // Events/audit log
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
    const id = task.id || uuidv4();
    const stmt = this.db.prepare(`
      INSERT INTO ${this.table} (id, type, payload, priority, scheduled_at, created_at)
      VALUES (@id, @type, @payload, @priority, @scheduled_at, datetime('now'))
    `);
    stmt.run({
      id,
      type: task.type,
      payload: JSON.stringify(task.payload),
      priority: task.priority || 0,
      scheduled_at: task.scheduled_at || new Date().toISOString()
    });
    this.logEvent('task_enqueued', id, null, { type: task.type, priority: task.priority });
    return id;
  }

  /**
   * Dequeue the next pending task atomically, marking it as locked
   */
  dequeue(agentName) {
    const now = new Date().toISOString();
    const stmt = this.db.transaction(() => {
      // Find an unlocked pending task with highest priority, earliest scheduled_at
      const select = this.db.prepare(`
        SELECT id, type, payload FROM ${this.table}
        WHERE status = 'pending' AND locked_by IS NULL
        AND (scheduled_at IS NULL OR scheduled_at <= datetime('now'))
        ORDER BY priority DESC, scheduled_at ASC
        LIMIT 1
      `);
      const task = select.get();
      if (!task) return null;

      // Lock it
      const update = this.db.prepare(`
        UPDATE ${this.table}
        SET locked_by = ?, locked_at = ?, status = 'processing', updated_at = datetime('now')
        WHERE id = ? AND locked_by IS NULL
      `);
      const result = update.run(agentName, now, task.id);
      if (result.changes === 0) {
        // Lost race, another agent locked it
        return null;
      }
      return {
        id: task.id,
        type: task.type,
        payload: JSON.parse(task.payload)
      };
    });
    return stmt;
  }

  /**
   * Mark task as completed
   */
  complete(taskId, result = {}) {
    const stmt = this.db.prepare(`
      UPDATE ${this.table}
      SET status = 'completed', updated_at = datetime('now')
      WHERE id = ?
    `);
    stmt.run(taskId);
    this.logEvent('task_completed', taskId, null, result);
  }

  /**
   * Mark task as failed and schedule retry or send to DLQ
   */
  fail(taskId, error, agentName) {
    const get = this.db.prepare(`SELECT retry_count, max_retries FROM ${this.table} WHERE id = ?`);
    const task = get.get(taskId);
    if (!task) return;

    const newRetry = task.retry_count + 1;
    if (newRetry >= task.max_retries) {
      // Move to DLQ
      this.moveToDlq(taskId, error, task.retry_count);
    } else {
      // Update with retry info, requeue
      const scheduledAt = new Date(Date.now() + this.retryBackoffMs * Math.pow(2, newRetry)).toISOString();
      const stmt = this.db.prepare(`
        UPDATE ${this.table}
        SET status = 'pending', locked_by = NULL, locked_at = NULL,
            retry_count = ?, last_error = ?, scheduled_at = ?, updated_at = datetime('now')
        WHERE id = ?
      `);
      stmt.run(newRetry, error, scheduledAt, taskId);
      this.logEvent('task_retry', taskId, agentName, { retry: newRetry, error, scheduled_at: scheduledAt });
    }
  }

  moveToDlq(taskId, error, retryCount) {
    const get = this.db.prepare(`SELECT type, payload FROM ${this.table} WHERE id = ?`);
    const task = get.get(taskId);
    if (!task) return;

    const insertDlq = this.db.prepare(`
      INSERT INTO ${this.dlqTable} (id, original_task_id, type, payload, reason, failed_at, retry_count)
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

    // Remove from main queue
    this.db.prepare(`DELETE FROM ${this.table} WHERE id = ?`).run(taskId);
    this.logEvent('task_dlq', taskId, null, { error, retry_count: retryCount });
  }

  logEvent(type, taskId, agent, payload = {}) {
    const stmt = this.db.prepare(`
      INSERT INTO events (event_type, task_id, agent, payload)
      VALUES (?, ?, ?, ?)
    `);
    stmt.run(type, taskId, agent, JSON.stringify(payload));
  }

  /**
   * Get task by ID (for inspection)
   */
  getTask(taskId) {
    return this.db.prepare(`SELECT * FROM ${this.table} WHERE id = ?`).get(taskId);
  }

  /**
   * List DLQ tasks
   */
  listDlq() {
    return this.db.prepare(`SELECT * FROM ${this.dlqTable} ORDER BY failed_at DESC`).all();
  }

  /**
   * Requeue a DLQ task manually
   */
  requeueDlq(dlqId) {
    const get = this.db.prepare(`SELECT * FROM ${this.dlqTable} WHERE id = ?`);
    const item = get.get(dlqId);
    if (!item) throw new Error('DLQ item not found');

    // Re-enqueue original task
    this.enqueue({
      id: item.original_task_id,
      type: item.type,
      payload: JSON.parse(item.payload)
    });

    // Remove from DLQ
    this.db.prepare(`DELETE FROM ${this.dlqTable} WHERE id = ?`).run(dlqId);
    this.logEvent('task_requeued', item.original_task_id, 'admin', { from_dlq: dlqId });
  }

  /**
   * Get metrics for monitoring
   */
  getMetrics() {
    const metrics = {};
    // Count by status
    const statusCounts = this.db.prepare(`
      SELECT status, COUNT(*) as count FROM ${this.table} GROUP BY status
    `).all();
    metrics.tasks_by_status = statusCounts.reduce((acc, row) => ({ ...acc, [row.status]: row.count }), {});
    // DLQ count
    metrics.dlq_count = this.db.prepare(`SELECT COUNT(*) as count FROM ${this.dlqTable}`).get().count;
    // Total events (last hour)
    metrics.events_last_hour = this.db.prepare(`
      SELECT COUNT(*) as count FROM events
      WHERE timestamp >= datetime('now', '-1 hour')
    `).get().count;
    return metrics;
  }

  close() {
    if (this.db) this.db.close();
  }
}

module.exports = { TaskQueue };

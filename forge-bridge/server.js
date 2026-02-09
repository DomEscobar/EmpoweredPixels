#!/usr/bin/env node
// Forge Bridge — HTTP relay for bot-to-bot coordination
// Usage: node server.js [port]
// Default port: 3001

const http = require('http');
const fs = require('fs');
const path = require('path');

const PORT = process.argv[2] || 3002;

// In-memory state
let pendingTasks = [];
let results = [];

function sendJson(res, data, status = 200) {
  res.writeHead(status, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify(data, null, 2));
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    let body = '';
    req.on('data', chunk => body += chunk);
    req.on('end', () => resolve(body));
    req.on('error', reject);
  });
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://${HOST || 'localhost'}:${PORT}`);

  // Auth guard for mutating endpoints
  if (['POST'].includes(req.method) && ['/webhook/task', '/webhook/response'].includes(url.pathname)) {
    authGuard(req, res, () => {}); // will send response if fails
    if (res.writableEnded) return;
  }

  if (req.method === 'POST' && url.pathname === '/webhook/task') {
    try {
      const body = await readBody(req);
      const task = JSON.parse(body);
      if (!task.id || !task.recipient) throw new Error('missing fields');
      pendingTasks.push(task);
      results = results.filter(r => r.taskId !== task.id); // clear previous
      console.log(`[BRIDGE] Task queued: ${task.id} → ${task.recipient}`);
      sendJson(res, { ok: true, taskId: task.id, queued: pendingTasks.length });
    } catch (e) {
      sendJson(res, { error: e.message }, 400);
    }
    return;
  }

  if (req.method === 'POST' && url.pathname === '/webhook/response') {
    try {
      const body = await readBody(req);
      const resp = JSON.parse(body);
      if (!resp.taskId || resp.result === undefined) throw new Error('missing fields');
      results.push({ taskId: resp.taskId, result: resp.result, ts: Date.now() });
      // Remove from pending if present
      pendingTasks = pendingTasks.filter(t => t.id !== resp.taskId);
      console.log(`[BRIDGE] Result received: ${resp.taskId}`);
      sendJson(res, { ok: true });
    } catch (e) {
      sendJson(res, { error: e.message }, 400);
    }
    return;
  }

  if (req.method === 'GET' && url.pathname === '/queue') {
    sendJson(res, { pending: pendingTasks, results });
    return;
  }

  if (req.method === 'GET' && url.pathname === '/health') {
    sendJson(res, { status: 'ok', pending: pendingTasks.length, results: results.length });
    return;
  }

  // Alias for compatibility with forge bot polling
  if (req.method === 'GET' && url.pathname === '/tasks') {
    sendJson(res, { pending: pendingTasks, results });
    return;
  }

  sendJson(res, { error: 'not found' }, 404);
});

const HOST = process.env.FORGE_BRIDGE_HOST || '0.0.0.0';
const SECRET = process.env.FORGE_BRIDGE_SECRET || null;

function authGuard(req, res, next) {
  if (SECRET) {
    const auth = req.headers['authorization'] || '';
    if (!auth.startsWith('Bearer ')) {
      sendJson(res, { error: 'unauthorized' }, 401);
      return;
    }
    const token = auth.slice(7);
    if (token !== SECRET) {
      sendJson(res, { error: 'forbidden' }, 403);
      return;
    }
  }
  next();
}

server.listen(PORT, HOST, () => {
  console.log(`[BRIDGE] Forge Bridge listening on :${PORT}`);
  console.log(`  POST /webhook/task   { id, recipient, payload }`);
  console.log(`  POST /webhook/response { taskId, result }`);
  console.log(`  GET  /queue    — state dump`);
  console.log(`  GET  /health   — liveness`);
});

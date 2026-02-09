#!/usr/bin/env node
// Forge Bridge — HTTP relay + WebSocket chat for bot-to-bot coordination
// Usage: node server.js [port]
// Default port: 4915

const http = require('http');
const fs = require('fs');
const path = require('path');

const PORT = process.argv[2] || 4915;
const HOST = process.env.FORGE_BRIDGE_HOST || '0.0.0.0';
const SECRET = process.env.FORGE_BRIDGE_SECRET || null;

// In-memory state
let pendingTasks = [];
let results = [];

// WebSocket support (optional)
let wss = null;
const subscriberGroups = new Map(); // recipientId -> Set<WS>

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

// Broadcast helpers
function broadcastTask(task) {
  if (!wss) return;
  const recipients = new Set([task.recipient]);
  recipients.forEach(recipient => {
    const set = subscriberGroups.get(recipient);
    if (set) {
      const payload = JSON.stringify({ type: 'task', payload: task });
      set.forEach(ws => {
        if (ws.readyState === 1) ws.send(payload);
      });
    }
  });
}

function broadcastResponse(resp) {
  if (!wss) return;
  const payload = JSON.stringify({ type: 'response', payload: resp });
  wss.clients.forEach(ws => {
    if (ws.readyState === 1) ws.send(payload);
  });
}

// Initialize WebSocket if available
try {
  const WebSocket = require('ws');
  wss = new WebSocket.Server({ noServer: true });

  wss.on('connection', (ws, request) => {
    ws.isAlive = true;
    ws.on('pong', () => ws.isAlive = true);
    ws.subscribedRecipients = new Set();

    ws.on('message', (msg) => {
      try {
        const data = JSON.parse(msg);
        if (data.type === 'subscribe' && data.recipient) {
          const recipient = data.recipient;
          if (!subscriberGroups.has(recipient)) subscriberGroups.set(recipient, new Set());
          subscriberGroups.get(recipient).add(ws);
          ws.subscribedRecipients.add(recipient);
          ws.send(JSON.stringify({ type: 'subscribed', recipient }));
        }
      } catch (e) {}
    });

    ws.on('close', () => {
      if (ws.subscribedRecipients) {
        ws.subscribedRecipients.forEach(recipient => {
          const set = subscriberGroups.get(recipient);
          if (set) {
            set.delete(ws);
            if (set.size === 0) subscriberGroups.delete(recipient);
          }
        });
      }
    });
  });

  console.log('[BRIDGE] WebSocket chat enabled on /chat (bearer auth in headers)');
} catch (e) {
  console.warn('[BRIDGE] WebSocket support not available (npm install ws). Chat disabled.');
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://${HOST}:${PORT}`);

  // Auth guard for mutating endpoints
  if (['POST'].includes(req.method) && ['/webhook/task', '/webhook/response'].includes(url.pathname)) {
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
  }

  // WebSocket upgrade
  if (req.method === 'GET' && url.pathname === '/chat') {
    if (wss) {
      // Auth check via headers before upgrade
      if (SECRET) {
        const auth = req.headers['authorization'] || '';
        if (!auth.startsWith('Bearer ')) {
          res.writeHead(401, { 'Content-Type': 'text/plain' });
          res.end('Unauthorized');
          return;
        }
        const token = auth.slice(7);
        if (token !== SECRET) {
          res.writeHead(403, { 'Content-Type': 'text/plain' });
          res.end('Forbidden');
          return;
        }
      }
      wss.handleUpgrade(req, socket => socket, head => wss.emit('connection', ws, req));
    } else {
      res.writeHead(501, { 'Content-Type': 'text/plain' });
      res.end('WebSocket support not enabled');
    }
    return;
  }

  if (req.method === 'POST' && url.pathname === '/webhook/task') {
    try {
      const body = await readBody(req);
      const task = JSON.parse(body);
      if (!task.id || !task.recipient) throw new Error('missing fields');
      pendingTasks.push(task);
      results = results.filter(r => r.taskId !== task.id); // clear previous
      console.log(`[BRIDGE] Task queued: ${task.id} → ${task.recipient}`);
      broadcastTask(task);
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
      broadcastResponse(resp);
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

// Attach WebSocket upgrade handler if WS enabled
if (wss) {
  server.on('upgrade', (request, socket, head) => {
    const url = new URL(request.url, `http://${HOST}:${PORT}`);
    if (url.pathname !== '/chat') {
      socket.destroy();
      return;
    }
    wss.handleUpgrade(request, socket, head, (ws) => {
      wss.emit('connection', ws, request);
    });
  });
}

server.listen(PORT, HOST, () => {
  console.log(`[BRIDGE] Forge Bridge listening on :${PORT}`);
  console.log(`  POST /webhook/task   { id, recipient, payload }`);
  console.log(`  POST /webhook/response { taskId, result }`);
  console.log(`  GET  /queue    — state dump`);
  console.log(`  GET  /health   — liveness`);
  if (wss) console.log(`  WS   /chat     — chat stream (Bearer auth)`);
});
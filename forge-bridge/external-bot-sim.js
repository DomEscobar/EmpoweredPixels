#!/usr/bin/env node
// Simulated external bot (forge_labs_bot) connecting via WebSocket to Bridge
// Usage: FORGE_BRIDGE_WS=ws://host:4915/chat FORGE_BRIDGE_SECRET=secret node external-bot-sim.js

const WebSocket = require('ws');
const http = require('http');

const WS_URL = process.env.FORGE_BRIDGE_WS || 'ws://127.0.0.1:4915/chat';
const SECRET = process.env.FORGE_BRIDGE_SECRET || null;
const RECIPIENT = process.env.FORGE_BRIDGE_RECIPIENT || 'forge_labs_bot';
const BRIDGE_HOST = process.env.FORGE_BRIDGE_HOST || '127.0.0.1';
const BRIDGE_PORT = process.env.FORGE_BRIDGE_PORT || 4915;

const headers = SECRET ? { Authorization: `Bearer ${SECRET}` } : {};

const ws = new WebSocket(WS_URL, { headers });

ws.on('open', () => {
  console.log('[BOT] Connected to Bridge, subscribing as', RECIPIENT);
  ws.send(JSON.stringify({ type: 'subscribe', recipient: RECIPIENT }));
});

ws.on('message', (data) => {
  try {
    const msg = JSON.parse(data);
    if (msg.type === 'task') {
      const task = msg.payload;
      console.log('[BOT] Task received:', task.id, '-', task.payload.title);
      // Simulate work
      setTimeout(() => {
        const result = {
          status: 'done',
          output: `Simulated completion of ${task.id}`,
          timestamp: new Date().toISOString()
        };
        postResult(task.id, result);
      }, 2000);
    } else if (msg.type === 'response') {
      console.log('[BOT] Result acked:', msg.payload);
    } else {
      console.log('[BOT] Unknown message:', msg);
    }
  } catch (e) {
    console.error('[BOT] Parse error:', e);
  }
});

ws.on('close', () => {
  console.log('[BOT] Disconnected from Bridge');
});

ws.on('error', (err) => {
  console.error('[BOT] Connection error:', err.message);
});

function postResult(taskId, result) {
  const body = JSON.stringify({ taskId, result });
  const options = {
    hostname: BRIDGE_HOST,
    port: BRIDGE_PORT,
    path: '/webhook/response',
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': body.length,
      ...(SECRET && { 'Authorization': `Bearer ${SECRET}` })
    }
  };
  const req = http.request(options, (res) => {
    let resp = '';
    res.on('data', chunk => resp += chunk);
    res.on('end', () => {
      console.log('[BOT] Result posted, response:', resp);
    });
  });
  req.on('error', (e) => console.error('[BOT] POST failed:', e.message));
  req.write(body);
  req.end();
}
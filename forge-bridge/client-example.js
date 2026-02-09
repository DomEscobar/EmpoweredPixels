#!/usr/bin/env node
// Example WebSocket client for Forge Bridge
// Usage: FORGE_BRIDGE_WS=ws://host:4915/chat FORGE_BRIDGE_SECRET=secret node client-example.js

const WebSocket = require('ws');

const WS_URL = process.env.FORGE_BRIDGE_WS || 'ws://127.0.0.1:4915/chat';
const SECRET = process.env.FORGE_BRIDGE_SECRET || null;
const RECIPIENT = process.env.FORGE_BRIDGE_RECIPIENT || 'forge_labs_bot';

const headers = SECRET ? { Authorization: `Bearer ${SECRET}` } : {};

const ws = new WebSocket(WS_URL, { headers });

ws.on('open', () => {
  console.log('[CLIENT] Connected, subscribing to', RECIPIENT);
  ws.send(JSON.stringify({ type: 'subscribe', recipient: RECIPIENT }));
});

ws.on('message', (data) => {
  try {
    const msg = JSON.parse(data);
    if (msg.type === 'task') {
      const task = msg.payload;
      console.log('[CLIENT] Task received:', task.id);
      // Execute task logic here, then:
      // const result = { status: 'done', output: '...' };
      // fetchResult(task.id, result);
    } else if (msg.type === 'response') {
      console.log('[CLIENT] Result acked:', msg.payload);
    } else {
      console.log('[CLIENT] Unknown message:', msg);
    }
  } catch (e) {
    console.error('[CLIENT] Parse error:', e);
  }
});

ws.on('close', () => {
  console.log('[CLIENT] Disconnected');
});

ws.on('error', (err) => {
  console.error('[CLIENT] Error:', err);
});

function fetchResult(taskId, result) {
  const body = JSON.stringify({ taskId, result });
  // could POST to /webhook/response
  // or send via WS if Bridge supports it (currently only HTTP)
}
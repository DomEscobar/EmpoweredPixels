#!/usr/bin/env node
// Director helper: assign a task to forge via the bridge

const http = require('http');

const [taskId, recipient = 'forge_labs_bot', title, ...descParts] = process.argv.slice(2);
if (!taskId || !title || descParts.length === 0) {
  console.error('Usage: node assign.js <taskId> <recipient> <title> <description...>');
  process.exit(1);
}

const description = descParts.join(' ');

const payload = {
  id: taskId,
  recipient,
  payload: {
    taskId,
    title,
    description,
    from: 'director',
    timestamp: Date.now()
  }
};

const data = JSON.stringify(payload);

const BRIDGE_PORT = process.env.FORGE_BRIDGE_PORT || 3002;
http.request({
  hostname: '127.0.0.1',
  port: BRIDGE_PORT,
  path: '/webhook/task',
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Content-Length': data.length
  }
}, res => {
  let body = '';
  res.on('data', chunk => body += chunk);
  res.on('end', () => {
    console.log('Bridge response:', body);
  });
}).on('error', e => console.error('Request failed:', e)).end(data);

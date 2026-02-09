#!/usr/bin/env node
// Forge helper: poll queue and ack tasks
// Usage: node collect.js

const http = require('http');

const BRIDGE_PORT = process.env.FORGE_BRIDGE_PORT || 3002;
function getQueue(cb) {
  http.get(`http://127.0.0.1:${BRIDGE_PORT}/queue`, res => {
    let body = '';
    res.on('data', d => body += d);
    res.on('end', () => cb(null, JSON.parse(body)));
  }).on('error', cb);
}

function postResult(taskId, result) {
  const data = JSON.stringify({ taskId, result });
  http.request({
    hostname: '127.0.0.1',
    port: BRIDGE_PORT,
    path: '/webhook/response',
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': data.length
    }
  }, res => {
    let body = '';
    res.on('data', d => body += d);
    res.on('end', () => console.log('Result posted:', body));
  }).on('error', e => console.error(e)).end(data);
}

// Simple poll loop
(function poll() {
  getQueue((err, q) => {
    if (err) return console.error('poll error:', err);
    const [task] = q.pending;
    if (task) {
      console.log('New task:', task.id, '-', task.payload.title);
      // Here: execute task logic, then:
      // postResult(task.id, { status: 'done', output: '...' });
    } else {
      // No tasks
    }
    setTimeout(poll, 5000);
  });
})();

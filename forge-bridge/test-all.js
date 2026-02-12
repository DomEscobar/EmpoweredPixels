#!/usr/bin/env node
/**
 * Forge Bridge Endpoint Test
 * Exercises: health, room (HTTP+WS), task queue, response.
 * Usage: node test-all.js
 */

const http = require('http');
const WebSocket = require('ws');

const BASE = process.env.BRIDGE_URL || 'http://v2202502215330313077.supersrv.de:4915';
const url = new URL(BASE);
const WS_BASE = `ws://${url.host}`;

let passCount = 0;
let failCount = 0;

function ok(msg) { passCount++; console.log(`✓ ${msg}`); }
function fail(msg, err) { failCount++; console.error(`✗ ${msg}`, err && err.message ? `(${err.message})` : ''); }

function httpJson(method, path, body = null) {
  return new Promise((resolve, reject) => {
    const options = {
      method,
      hostname: url.hostname,
      port: url.port,
      path,
      headers: { 'Content-Type': 'application/json' }
    };
    const req = http.request(options, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          resolve({ status: res.statusCode, data: data ? JSON.parse(data) : null });
        } catch (e) {
          resolve({ status: res.statusCode, data: data });
        }
      });
    });
    req.on('error', reject);
    if (body) req.write(JSON.stringify(body));
    req.end();
  });
}

(async () => {
  console.log(`Testing Forge Bridge at ${BASE}\n`);

  // 1. GET /health
  try {
    const h = await httpJson('GET', '/health');
    if (h.status === 200 && h.data.status === 'ok') ok('GET /health');
    else fail('GET /health', new Error('invalid response'));
  } catch (e) { fail('GET /health', e); }

  // 2. POST /room (chat)
  let chatMsg = { from: 'test-script', text: 'Hello room from test!', topic: 'testing' };
  try {
    const r = await httpJson('POST', '/room', chatMsg);
    if (r.status === 200 && r.data.ok) ok('POST /room');
    else fail('POST /room', new Error('not ok'));
  } catch (e) { fail('POST /room', e); }

  // 3. GET /room history
  try {
    const rh = await httpJson('GET', '/room');
    if (rh.status === 200 && Array.isArray(rh.data.history) && rh.data.history.length > 0) {
      ok('GET /room (history present)');
    } else fail('GET /room (history)', new Error('empty or bad'));
  } catch (e) { fail('GET /room (history)', e); }

  // 4. WebSocket /room
  const wsRoomUrl = `${WS_BASE}/room`;
  await new Promise((resolve, reject) => {
    const ws = new WebSocket(wsRoomUrl);
    let gotHistory = false;
    let gotEcho = false;

    ws.on('open', () => {
      ok(`WS CONNECT /room`);
      ws.send(JSON.stringify({ from: 'ws-test', text: 'WS message', topic: 'ws' }));
    });

    ws.on('message', (data) => {
      try {
        const msg = JSON.parse(data);
        if (msg.type === 'chat' && msg.from === 'test-script') {
          ok('WS RECV history echo (from HTTP POST)');
          gotHistory = true;
        }
        if (msg.type === 'chat' && msg.from === 'ws-test') {
          ok('WS RECV own broadcast');
          gotEcho = true;
        }
        if (gotHistory && gotEcho) {
          ws.close();
          resolve();
        }
      } catch (e) { fail('WS parse', e); }
    });

    ws.on('error', (err) => {
      reject(err);
    });

    ws.on('close', () => {
      if (!gotHistory || !gotEcho) resolve(); // still resolve after timeout?
    });

    setTimeout(() => {
      if (!gotHistory || !gotEcho) reject(new Error('WS timeout'));
    }, 5000);
  }).catch(e => fail('WS /room', e));

  // 5. POST /webhook/task
  const testTask = { id: 'TEST-001', recipient: 'forge', payload: { cmd: 'ping' } };
  try {
    const rt = await httpJson('POST', '/webhook/task', testTask);
    if (rt.status === 200 && rt.data.ok) ok('POST /webhook/task');
    else fail('POST /webhook/task', new Error('not ok'));
  } catch (e) { fail('POST /webhook/task', e); }

  // 6. GET /queue (should include pending task)
  try {
    const q = await httpJson('GET', '/queue');
    if (q.status === 200 && Array.isArray(q.data.pending) && q.data.pending.some(t => t.id === 'TEST-001')) {
      ok('GET /queue (task pending)');
    } else fail('GET /queue (task pending)', new Error('task missing'));
  } catch (e) { fail('GET /queue (task pending)', e); }

  // 7. POST /webhook/response
  try {
    const resp = await httpJson('POST', '/webhook/response', { taskId: 'TEST-001', result: { status: 'ok' } });
    if (resp.status === 200 && resp.data.ok) ok('POST /webhook/response');
    else fail('POST /webhook/response', new Error('not ok'));
  } catch (e) { fail('POST /webhook/response', e); }

  // 8. GET /queue (task cleared)
  try {
    const q2 = await httpJson('GET', '/queue');
    if (q2.status === 200 && !q2.data.pending.some(t => t.id === 'TEST-001')) {
      ok('GET /queue (task cleared)');
    } else fail('GET /queue (task cleared)', new Error('task still present'));
  } catch (e) { fail('GET /queue (task cleared)', e); }

  console.log('\n--- Summary ---');
  console.log(`Passed: ${passCount}, Failed: ${failCount}`);
  process.exit(failCount > 0 ? 1 : 0);
})();

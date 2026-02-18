#!/usr/bin/env node
/**
 * Webhook End-to-End Test
 * Tests: POST /webhook/task and POST /webhook/github
 */

const http = require('http');
const { WebhookServer } = require('./webhook-server.js');
const { TaskQueue } = require('./task-queue.js');

const config = require('./agency-config.json');
const BASE_URL = `http://localhost:${config.agency.webhooks.port}`;

let server;
let queue;
let passCount = 0;
let failCount = 0;

function ok(msg) { passCount++; console.log(`✓ ${msg}`); }
function fail(msg, err) { failCount++; console.error(`✗ ${msg}`, err && err.message ? `(${err.message})` : ''); }

function httpPost(path, body) {
  return new Promise((resolve, reject) => {
    const url = new URL(path, BASE_URL);
    const options = {
      hostname: 'localhost',
      port: config.agency.webhooks.port,
      path: url.pathname,
      method: 'POST',
      headers: { 'Content-Type': 'application/json' }
    };
    const req = http.request(options, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          resolve({ status: res.statusCode, data: JSON.parse(data) });
        } catch (e) {
          resolve({ status: res.statusCode, data: data });
        }
      });
    });
    req.on('error', reject);
    req.write(JSON.stringify(body));
    req.end();
  });
}

async function run() {
  console.log('=== Webhook End-to-End Test ===\n');

  // Initialize in-memory queue (force by pointing to invalid path)
  queue = new TaskQueue({ 
    type: 'memory', 
    path: '/nonexistent/no-db-allowed-in-test',  // Force memory fallback
    table: 'tasks',
    dead_letter_table: 'dlq',
    max_retries: 3
  });
  await queue.initialize();
  queue.createTables();

  // Start webhook server
  server = new WebhookServer(config, queue);
  await server.start();
  console.log(`Webhook server started on port ${config.agency.webhooks.port}\n`);

  // Test 1: POST /webhook/task
  try {
    const taskPayload = {
      type: 'test.task',
      payload: { cmd: 'echo hello' },
      priority: 5
    };
    const res = await httpPost('/webhook/task', taskPayload);
    if (res.status === 200) {
      ok('POST /webhook/task returns 200');
    } else {
      fail('POST /webhook/task returns 200', new Error(`Got ${res.status}`));
    }

    // Verify task was queued
    const pendingTasks = Array.from(queue.queue.tasks.values()).filter(t => t.status === 'pending');
    if (pendingTasks.length > 0 && pendingTasks[0].type === 'test.task') {
      ok('Task enqueued to queue');
    } else {
      fail('Task enqueued to queue', new Error('Task not found'));
    }
  } catch (e) { fail('POST /webhook/task', e); }

  // Test 2: POST /webhook/github (issue opened)
  try {
    const githubPayload = {
      action: 'opened',
      issue: {
        number: 42,
        title: 'Test Issue',
        body: 'Issue description',
        html_url: 'https://github.com/test/repo/issues/42'
      }
    };
    const res = await httpPost('/webhook/github', githubPayload);
    if (res.status === 200) {
      ok('POST /webhook/github returns 200');
    } else {
      fail('POST /webhook/github returns 200', new Error(`Got ${res.status}`));
    }

    // Verify GitHub task was queued
    const pendingTasks2 = Array.from(queue.queue.tasks.values()).filter(t => t.status === 'pending');
    const githubTask = pendingTasks2.find(t => t.type === 'director.plan');
    if (githubTask) {
      const payload = JSON.parse(githubTask.payload);
      if (payload.description && payload.description.includes('42')) {
        ok('GitHub issue converted to director.plan task');
      } else {
        fail('GitHub issue converted to director.plan task', new Error('Invalid payload'));
      }
    } else {
      fail('GitHub issue converted to director.plan task', new Error('Task not found'));
    }
  } catch (e) { fail('POST /webhook/github', e); }

  // Test 3: Invalid route
  try {
    const res = await httpPost('/webhook/invalid', {});
    if (res.status === 404) {
      ok('Invalid route returns 404');
    } else {
      fail('Invalid route returns 404', new Error(`Got ${res.status}`));
    }
  } catch (e) { fail('Invalid route returns 404', e); }

  // Cleanup
  server.stop();
  console.log('\n--- Summary ---');
  console.log(`Passed: ${passCount}, Failed: ${failCount}`);
  process.exit(failCount > 0 ? 1 : 0);
}

run().catch(err => {
  console.error('Test failed:', err);
  if (server) server.stop();
  process.exit(1);
});

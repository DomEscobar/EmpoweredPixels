#!/usr/bin/env node
/**
 * OpenCode Agency — Single-Process Demo
 * Runs Director and Coder in the same process with a shared in-memory queue.
 * Demonstrates the full loop without SQLite multi-process complexity.
 */

const fs = require('fs');
if (!fs.existsSync('logs')) fs.mkdirSync('logs');
if (fs.existsSync('logs/coder.log')) fs.unlinkSync('logs/coder.log'); // clean

const { TaskQueue } = require('./task-queue');
const { CoderAgent } = require('./agents/coder-advanced');

(async () => {
  // Force memory queue
  const queueConfig = {
    type: 'memory',
    max_retries: 3,
    retry_backoff_ms: 2000,
    dead_letter_table: 'dlq'
  };
  const queue = new TaskQueue(queueConfig);
  await queue.initialize();

  const fullConfig = require('./agency-config.json');
  // const directorConfig = fullConfig.agents.director; // not needed
  const coderConfig = fullConfig.agents.coder;

  // Create agents with the same queue instance
  // const director = new DirectorAgent(directorConfig, queue);
  const coder = new CoderAgent(coderConfig, queue);

  // Start agents
  // await director.start();
  await coder.start();

  console.log('\n=== Demo agents started (coder only) ===\n');

  // Submit a task directly to coder for simplicity
  const taskId = queue.enqueue({
    id: 'demo-task-1',
    type: 'coder_task',
    payload: { description: 'Create AGENCY_DEMO.md in /root/EmpoweredPixels with text: Single-process demo succeeded.' },
    priority: 10,
    assigned_to: 'coder'
  });
  console.log(`Enqueued task ${taskId}`);

  // Wait for processing
  await new Promise(resolve => setTimeout(resolve, 8000));

  // Show coder logs
  const logPath = 'logs/coder.log';
  if (fs.existsSync(logPath)) {
    console.log('\n=== Coder Log (tail) ===');
    const logContent = fs.readFileSync(logPath, 'utf8');
    console.log(logContent);
  } else {
    console.log('\nNo coder log found');
  }

  // Check result
  console.log('\n=== Queue Metrics ===');
  console.log(queue.getMetrics());

  // Check if file created (look for AGENCY_TASK_*.txt)
  const pattern = /AGENCY_TASK_.*\.txt/;
  const files = fs.readdirSync('/root/EmpoweredPixels').filter(f => pattern.test(f));
  if (files.length > 0) {
    console.log('\n=== Demo Artifact Created ===');
    files.forEach(f => {
      console.log(`File: /root/EmpoweredPixels/${f}`);
      console.log(fs.readFileSync(`/root/EmpoweredPixels/${f}`, 'utf8'));
      fs.unlinkSync(`/root/EmpoweredPixels/${f}`); // cleanup
    });
  } else {
    console.log('\nDemo artifact NOT created');
  }

  // Stop agents
  await coder.stop();
  queue.close();

  console.log('\nDemo completed');
})().catch(err => {
  console.error('Demo failed:', err);
  process.exit(1);
});

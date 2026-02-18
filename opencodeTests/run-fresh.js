#!/usr/bin/env node
const { spawn } = require('child_process');
const cmd = '/root/.opencode/bin/opencode';
const args = [
  'run', '--dir', '/root/EmpoweredPixels',
  '--model', 'opencode/minimax-m2.5-free',
  'Create TEST_FRESH_20260217_2.md with content: fresh ok 2'
];
const env = {
  OPENROUTER_API_KEY: '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',
  OPENCODE_MODEL: 'opencode/minimax-m2.5-free'
};

console.log('Spawning opencode with 60s timeout...');
const child = spawn(cmd, args, { env, stdio: 'inherit' });

let timedOut = false;
const timer = setTimeout(() => {
  timedOut = true;
  process.kill(-child.pid, 'SIGTERM');
  console.log('Timed out, sent SIGTERM');
}, 60000);

child.on('exit', (code, signal) => {
  if (!timedOut) clearTimeout(timer);
  if (code === 0) {
    console.log('Task succeeded');
  } else {
    console.log(`Task exited code=${code} signal=${signal}`);
  }
});
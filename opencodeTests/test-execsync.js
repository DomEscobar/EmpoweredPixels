#!/usr/bin/env node
const { execSync } = require('child_process');
const cmd = `OPENROUTER_API_KEY=${process.env.OPENROUTER_API_KEY} OPENCODE_MODEL=${process.env.OPENCODE_MODEL} /root/.opencode/bin/opencode run --dir /root/EmpoweredPixels --model opencode/minimax-m2.5-free "Create TEST_FRESH_EXECSYNC.md with content: execSync ok"`;
console.log('Running command:', cmd);
try {
  const out = execSync(cmd, { stdio: 'inherit' });
  console.log('Command succeeded');
} catch (e) {
  console.error('Error:', e.status, e.message);
}
#!/usr/bin/env node
const { execSync } = require('child_process');
const cmd = `OPENROUTER_API_KEY=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 timeout 20s /root/.opencode/bin/opencode run --dir /root/EmpoweredPixels --model opencode/minimax-m2.5-free "Create SIMPLE_WEBHOOK.md with content: webhook ok"`;
try {
  console.log(execSync(cmd, { encoding: 'utf-8', stdio: 'inherit' }));
} catch (e) {
  console.error('Command failed:', e.message);
}
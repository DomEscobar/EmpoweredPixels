#!/usr/bin/env node
const { execFile } = require('child_process');
const desc = "Create TEST_FRESH_20260217.md with content: fresh ok";
const args = ['run','--dir','/root/EmpoweredPixels','--model','opencode/minimax-m2.5-free', desc];
const env = { ...process.env, OPENROUTER_API_KEY: process.env.OPENROUTER_API_KEY, NODE_ENV: 'production' };
console.log('Running:', '/root/.opencode/bin/opencode', args);
console.log('ENV:', { OPENROUTER_API_KEY: env.OPENROUTER_API_KEY ? 'set' : 'missing', OPENCODE_MODEL: env.OPENCODE_MODEL });
execFile('/root/.opencode/bin/opencode', args, { env }, (err, stdout, stderr) => {
  console.log('err:', err && err.code, err && err.message);
  console.log('stdout:', stdout.toString().slice(0, 500));
  console.log('stderr:', stderr.toString().slice(0, 500));
});
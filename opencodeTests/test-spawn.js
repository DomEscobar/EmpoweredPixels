#!/usr/bin/env node
const { spawn } = require('child_process');
const desc = "Create TEST_FRESH_SPAWN.md with content: fresh spawn ok";
const args = ['run','--dir','/root/EmpoweredPixels','--model','opencode/minimax-m2.5-free', desc];
const env = { ...process.env, OPENROUTER_API_KEY: process.env.OPENROUTER_API_KEY, NODE_ENV: 'production' };
console.log('Spawning opencode...');
const p = spawn('/root/.opencode/bin/opencode', args, { env, stdio: 'inherit' });
let out = '', err = '';
p.stdout.on('data', d => out += d);
p.stderr.on('data', d => err += d);
p.on('close', code => {
  console.log('Exit code:', code);
  console.log('stdout snippet:', out.slice(0, 300));
  console.log('stderr snippet:', err.slice(0, 300));
});
p.on('error', e => console.error('Spawn error:', e));
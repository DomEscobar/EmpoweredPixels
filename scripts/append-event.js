#!/usr/bin/env node
/**
 * Append-event.js — Event store append helper
 * Usage: node append-event.js <type> key1=value1 key2=value2 ...
 * Appends a JSON event to /root/.openclaw/events.jsonl
 */
const fs = require('fs');
const path = require('path');

const args = process.argv.slice(2);
if (args.length < 1) {
  console.error('Usage: append-event.js <type> [key=value ...]');
  process.exit(1);
}

const type = args[0];
const data = {};

for (let i = 1; i < args.length; i++) {
  const pair = args[i];
  const idx = pair.indexOf('=');
  if (idx === -1) {
    console.warn(`Skipping malformed arg: ${pair}`);
    continue;
  }
  const key = pair.slice(0, idx);
  const value = pair.slice(idx + 1);
  // Try to parse numbers/booleans if they look like it
  if (/^-?\d+$/.test(value)) data[key] = Number(value);
  else if (/^(true|false)$/.test(value)) data[key] = value === 'true';
  else data[key] = value;
}

const event = {
  type,
  timestamp: new Date().toISOString(),
  data
};

const logPath = path.resolve('/root/.openclaw/events.jsonl');
fs.appendFileSync(logPath, JSON.stringify(event) + '\n', { encoding: 'utf8' });
console.log(`Event appended: ${type}`);

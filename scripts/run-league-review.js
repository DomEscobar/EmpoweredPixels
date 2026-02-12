#!/usr/bin/env node
/**
 * Manual trigger for league-reviewer agent.
 * Usage: BRIDGE_URL=http://127.0.0.1:4915 node run-league-review.js
 */
const { sessions_spawn } = require('@openclaw/sdk'); // placeholder; we'll use openclaw CLI
const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

async function trigger() {
  const msg = 'Run full league review session';
  // Using openclaw agent spawn via gateway CLI
  const cmd = `openclaw agents spawn --agent player --label league-reviewer --message "${msg}"`;
  try {
    const { stdout } = await execPromise(cmd);
    console.log('Spawned league-reviewer:', stdout.trim());
  } catch (e) {
    console.error('Failed to spawn:', e.stderr || e.message);
    process.exit(1);
  }
}

trigger();

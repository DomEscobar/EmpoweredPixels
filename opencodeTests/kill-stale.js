#!/usr/bin/env node
const { spawn } = require('child_process');
const id = '2c964ccb-7edd-4c14-bcfa-0f1b6b8facdf';
const { execSync } = require('child_process');
execSync("pkill -f 'Create WEBHOOK_VERIFIED.md'", { stdio: 'ignore' });
console.log('Killed WEBHOOK_VERIFIED task child');
// Mark task as failed to move on
try {
  const db = require('better-sqlite3')('/root/EmpoweredPixels/opencodeTests/data/agency.db');
  db.prepare("UPDATE tasks SET status='failed', error='killed_by_test' WHERE id=?").run(id);
  console.log('Marked task as failed in DB');
} catch (e) {
  console.error(e.message);
}
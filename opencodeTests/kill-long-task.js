#!/usr/bin/env node
const { execSync } = require('child_process');
try {
  execSync("pkill -f 'Webhook test: verify end-to-end processing'", { stdio: 'ignore' });
  console.log('Killed long task');
} catch (e) {
  console.log('No matching process');
}
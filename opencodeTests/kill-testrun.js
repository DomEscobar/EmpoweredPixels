#!/usr/bin/env node
const { execSync } = require('child_process');
try {
  execSync("pkill -f 'Create TEST_RUN.md'", { stdio: 'ignore' });
  console.log('Killed TEST_RUN task');
} catch (e) {
  console.log('No matching process');
}
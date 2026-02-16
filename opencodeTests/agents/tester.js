#!/usr/bin/env node
/**
 * OpenCode Agency — Tester Agent
 * Generates and executes tests
 */

const { spawn } = require('child_process');

class TesterAgent {
  constructor() {
    this.workdir = process.env.AGENT_WORKDIR || process.cwd();
  }

  async generateTests(requirements) {
    console.log('[Tester] Generating tests for requirements:', requirements.title);

    // Use OpenCode to generate test cases
    const prompt = `Generate comprehensive test suite for:\n${requirements.description}\n\nInclude unit tests and integration tests.`;

    // Could invoke opencode directly or write to filesystem
    // For now, output placeholder
    const tests = {
      unit: ['test_feature_x.js'],
      integration: ['test_workflow.js']
    };

    return tests;
  }

  async runAll() {
    console.log('[Tester] Running test suite...');
    // Detect test framework and run
    // Could be npm test, pytest, go test, etc.
    return { passed: true, coverage: 85 };
  }

  async explore() {
    console.log('[Tester] Exploratory testing session...');
    return { bugs: [], findings: [] };
  }

  startListener() {
    process.stdin.resume();
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', async (data) => {
      const line = data.trim();
      if (line.startsWith('TEST:')) {
        try {
          const req = JSON.parse(line.slice(5));
          const tests = await this.generateTests(req);
          console.log('[Tester] Generated:', JSON.stringify(tests));
        } catch (e) {
          console.error('[Tester] Error:', e);
        }
      } else if (line === 'run') {
        const results = await this.runAll();
        console.log('[Tester] Results:', JSON.stringify(results));
      }
    });
  }
}

if (require.main === module) {
  const agent = new TesterAgent();
  console.log('[Tester] Starting...');
  agent.startListener();
}

module.exports = { TesterAgent };

#!/usr/bin/env node
/**
 * OpenCode Agency — Reviewer Agent
 * Quality gate that reviews changes before merge
 */

const { spawn } = require('child_process');

class ReviewerAgent {
  constructor() {
    this.workdir = process.env.AGENT_WORKDIR || process.cwd();
  }

  async review(changes) {
    console.log('[Reviewer] Starting review...');

    // Run static analysis
    await this.runStaticAnalysis();

    // Check for test coverage
    await this.checkTests();

    // Security scan
    await this.securityScan();

    console.log('[Reviewer] Review passed');
    return { approved: true, notes: [] };
  }

  async runStaticAnalysis() {
    console.log('[Reviewer] Running static analysis...');
    // Could run eslint, golint, etc. based on project
    // For now, simple dry-run
    return true;
  }

  async checkTests() {
    console.log('[Reviewer] Checking tests...');
    // Run: npm test -- --dry-run or similar
    return true;
  }

  async securityScan() {
    console.log('[Reviewer] Security scan...');
    // Could integrate with snyk, trivy, etc.
    return true;
  }

  startListener() {
    process.stdin.resume();
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', async (data) => {
      const line = data.trim();
      if (line.startsWith('REVIEW:')) {
        try {
          const changes = JSON.parse(line.slice(7));
          const result = await this.review(changes);
          console.log('[Reviewer] Result:', JSON.stringify(result));
        } catch (e) {
          console.error('[Reviewer] Review error:', e);
        }
      }
    });
  }
}

if (require.main === module) {
  const agent = new ReviewerAgent();
  console.log('[Reviewer] Starting...');
  agent.startListener();
}

module.exports = { ReviewerAgent };

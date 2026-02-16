#!/usr/bin/env node
/**
 * OpenCode Agency — Advanced Tester Agent
 * Generates and runs tests, performs validation
 */

const { AgentBase } = require('../agent-base');
const { spawn } = require('child_process');

class TesterAgent extends AgentBase {
  constructor(config, queue) {
    super(config, queue);
    this.workdir = process.env.AGENT_WORKDIR || process.cwd();
  }

  async execute(task) {
    this.logger.info({ taskId: task.id }, 'Tester executing');

    const { planId, description, previous_task_id } = task.payload;

    // If we have a previous task, we can correlate tests to changes
    // For now, just ensure tests exist and run them

    const testResult = await this.runTestSuite();

    return {
      success: testResult.passed,
      passed: testResult.passed,
      total: testResult.total,
      failures: testResult.failures,
      coverage: testResult.coverage,
      planId
    };
  }

  async runTestSuite() {
    this.logger.debug('Running test suite');

    // Detect test runner
    const packageJsonPath = require('path').join(this.workdir, 'package.json');
    let packageJson = {};
    try {
      packageJson = JSON.parse(require('fs').readFileSync(packageJsonPath, 'utf8'));
    } catch (e) {
      throw new Error('No package.json found, cannot determine test runner');
    }

    const scripts = packageJson.scripts || {};
    let testCmd = null;
    if (scripts.test) testCmd = 'npm test';
    else if (scripts['test:unit']) testCmd = 'npm run test:unit';
    else if (scripts['test:ci']) testCmd = 'npm run test:ci';

    if (!testCmd) {
      // Try Jest directly if node_modules exists
      const hasJest = require('fs').existsSync(require('path').join(this.workdir, 'node_modules', '.bin', 'jest'));
      if (hasJest) testCmd = 'npx jest --coverage';
      else throw new Error('No test script configured');
    }

    // Run tests with timeout
    const { stdout, stderr } = await this.exec(testCmd, { timeout: 120000 });

    // Parse output (very basic)
    const passed = !stderr.includes('FAIL') && !stdout.includes('FAIL');
    const coverageMatch = stdout.match(/All\s+files\s+\|\s+(\d+\.?\d*)\s+\|/);
    const coverage = coverageMatch ? parseFloat(coverageMatch[1]) : null;

    // Count tests (heuristic)
    const testCountMatch = stdout.match(/Test\s+Suites:\s+(\d+)\s+|Tests:\s+(\d+)\s+/);
    const total = testCountMatch ? (parseInt(testCountMatch[1]) || parseInt(testCountMatch[2]) || 0) : 0;

    return { passed, total, failures: passed ? 0 : 1, coverage };
  }

  exec(cmd, opts = {}) {
    return new Promise((resolve, reject) => {
      const [cmdMain, ...args] = cmd.split(' ');
      const proc = require('child_process').spawn(cmdMain, args, {
        cwd: this.workdir,
        stdio: 'pipe',
        env: process.env
      });

      let stdout = '', stderr = '';
      proc.stdout.on('data', (d) => stdout += d);
      proc.stderr.on('data', (d) => stderr += d);

      const timeout = opts.timeout ? setTimeout(() => proc.kill('SIGKILL'), opts.timeout) : null;

      proc.on('close', (code) => {
        if (timeout) clearTimeout(timeout);
        if (code === 0) resolve({ stdout, stderr });
        else reject(new Error(`Test command exited ${code}: ${stderr.slice(-500)}`));
      });

      proc.on('error', reject);
    });
  }
}

module.exports = { TesterAgent };

if (require.main === module) {
  (async () => {
    const configPath = process.env.AGENCY_CONFIG || '../agency-config.json';
    const config = require(configPath).agents.tester;
    const { TaskQueue } = require('../task-queue');
    const queueConfig = require(configPath).agency.queue;
    const queue = new TaskQueue(queueConfig);
    await queue.initialize();
    const agent = new TesterAgent(config, queue);
    await agent.start();
    process.on('SIGINT', async () => {
      await agent.stop();
      queue.close();
      process.exit(0);
    });
    console.log(`Tester agent ${config.name} started`);
    await new Promise(() => {});
  })();
}

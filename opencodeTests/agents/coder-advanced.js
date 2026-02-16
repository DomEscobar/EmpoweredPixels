#!/usr/bin/env node
/**
 * OpenCode Agency — Advanced Coder Agent
 * Executes coding tasks using OpenCode CLI with queue-based task consumption
 */

const { AgentBase } = require('./agent-base');
const { spawn } = require('child_process');

class CoderAgent extends AgentBase {
  constructor(config, queue) {
    super(config, queue);
    this.opencodePath = process.env.OPENCODE_PATH || 'opencode';
    this.model = this.env.OPENCODE_MODEL || 'openrouter/auto';
    this.workdir = process.env.AGENT_WORKDIR || process.cwd();
  }

  async execute(task) {
    this.logger.info({ taskId: task.id, description: task.payload.description }, 'Coder executing task');

    // MOCK MODE: For testing without OpenRouter API key
    if (process.env.OPENCODE_MOCK === 'true') {
      this.logger.info('Running in OPENCODE_MOCK mode – simulating work');
      const fs = require('fs');
      const path = require('path');
      const workdir = this.workdir;
      // Simulate some work
      await new Promise(resolve => setTimeout(resolve, 2000));
      // Write a simple artifact to demonstrate execution
      const outFile = path.join(workdir, `AGENCY_TASK_${task.id}.txt`);
      const content = `OpenCode Agency Mock Execution\nTask ID: ${task.id}\nDescription: ${task.payload.description}\nTimestamp: ${new Date().toISOString()}\n`;
      fs.writeFileSync(outFile, content, 'utf8');
      this.logger.info({ file: outFile }, 'Mock artifact created');
      return {
        success: true,
        output: `Mock execution completed – wrote ${outFile}`,
        files_changed: [outFile]
      };
    }

    // Real OpenCode execution
    const args = [
      '--prompt', task.payload.description,
      '--model', this.model,
      '--workdir', this.workdir,
      '--auto-approve' // assume OpenCode supports this
    ];

    return new Promise((resolve, reject) => {
      const proc = spawn(this.opencodePath, args, {
        stdio: 'pipe',
        env: { ...process.env, NODE_ENV: 'production' }
      });

      let stdout = '', stderr = '';
      proc.stdout.on('data', (d) => stdout += d);
      proc.stderr.on('data', (d) => stderr += d);

      proc.on('close', (code) => {
        if (code === 0) {
          this.logger.info({ taskId: task.id, output: stdout.slice(-1000) }, 'OpenCode completed');
          resolve({
            success: true,
            output: stdout,
            files_changed: this.extractChangedFiles(stdout)
          });
        } else {
          this.logger.error({ taskId: task.id, code, stderr }, 'OpenCode failed');
          reject(new Error(`OpenCode exited with ${code}: ${stderr.slice(-500)}`));
        }
      });

      proc.on('error', (err) => {
        reject(err);
      });
    });
  }

  extractChangedFiles(output) {
    // Simple heuristic: look for file paths in output
    const fileRegex = /(?:created|modified|deleted|changed):\s+([^\s]+)/gi;
    const files = [];
    let match;
    while ((match = fileRegex.exec(output)) !== null) {
      files.push(match[1]);
    }
    return files;
  }
}

module.exports = { CoderAgent };

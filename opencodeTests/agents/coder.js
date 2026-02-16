#!/usr/bin/env node
/**
 * OpenCode Agency — Coder Agent
 * Executes coding tasks using OpenCode CLI
 */

const { spawn } = require('child_process');

class CoderAgent {
  constructor() {
    this.opencodePath = process.env.OPENCODE_PATH || 'opencode';
    this.model = process.env.OPENCODE_MODEL || 'openrouter/auto';
    this.workdir = process.env.AGENT_WORKDIR || process.cwd();
    this.idle = true;
  }

  async execute(task) {
    console.log(`[Coder] Starting task: ${task.title}`);

    // Build OpenCode command
    // Using opencode with --prompt to execute specific task
    const args = [
      '--prompt', task.description || task.title,
      '--model', this.model,
      '--workdir', this.workdir
    ];

    return new Promise((resolve, reject) => {
      const proc = spawn(this.opencodePath, args, {
        stdio: 'inherit',
        env: { ...process.env }
      });

      proc.on('close', (code) => {
        this.idle = true;
        if (code === 0) {
          console.log(`[Coder] Task completed: ${task.title}`);
          resolve({ success: true, task });
        } else {
          console.error(`[Coder] Task failed with exit ${code}`);
          reject(new Error(`OpenCode exited with ${code}`));
        }
      });

      proc.on('error', (err) => {
        reject(err);
      });
    });
  }

  startListener() {
    // Listen for tasks from director via stdin
    process.stdin.resume();
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', async (data) => {
      const line = data.trim();
      if (line.startsWith('TASK:')) {
        try {
          const task = JSON.parse(line.slice(5));
          this.idle = false;
          await this.execute(task).catch((err) => {
            console.error('[Coder] Task error:', err);
            this.idle = true;
          });
        } catch (e) {
          console.error('[Coder] Failed to parse task:', e);
        }
      }
    });
  }
}

// Run as standalone agent
if (require.main === module) {
  const agent = new CoderAgent();
  console.log('[Coder] Starting OpenCode agent...');
  agent.startListener();
}

module.exports = { CoderAgent };

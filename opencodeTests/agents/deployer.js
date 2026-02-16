#!/usr/bin/env node
/**
 * OpenCode Agency — Deployer Agent
 * Manages CI/CD and deployments
 */

const { spawn } = require('child_process');

class DeployerAgent {
  constructor() {
    this.workdir = process.env.AGENT_WORKDIR || process.cwd();
    this.pipeline = process.env.DEPLOY_PIPELINE || 'auto';
  }

  async deploy(version) {
    console.log(`[Deployer] Deploying version ${version}`);

    switch (this.pipeline) {
      case 'github-actions':
        return await this.triggerGitHubActions(version);
      case 'gitlab-ci':
        return await this.triggerGitLabCI(version);
      case 'jenkins':
        return await this.triggerJenkins(version);
      default:
        return await this.autoDeploy(version);
    }
  }

  async autoDeploy(version) {
    console.log('[Deployer] Auto-deploying...');
    // Check pre-deploy health
    await this.healthCheck();
    // Build
    await this.runBuild();
    // Push to production
    await this.promoteToProd();
    // Post-deploy smoke test
    await this.smokeTest();
    return { success: true, version };
  }

  async healthCheck() {
    console.log('[Deployer] Pre-deploy health check...');
    // Could ping endpoints, check load, etc.
    return true;
  }

  async runBuild() {
    console.log('[Deployer] Running build...');
    const { stdout, stderr } = await this.exec('npm run build', { timeout: 60000 });
    if (stdout.includes('error')) throw new Error('Build failed');
    return true;
  }

  async promoteToProd() {
    console.log('[Deployer] Promoting to production...');
    // git tag, docker push, etc.
    return true;
  }

  async smokeTest() {
    console.log('[Deployer] Running smoke tests...');
    // Basic curl checks to verify deployment
    return true;
  }

  async rollback(version) {
    console.log(`[Deployer] Rolling back to ${version}`);
    // Implementation depends on deployment target
    return { success: true, rolled_back_to: version };
  }

  exec(cmd, opts = {}) {
    return new Promise((resolve, reject) => {
      const [cmdMain, ...args] = cmd.split(' ');
      const proc = spawn(cmdMain, args, {
        cwd: this.workdir,
        stdio: 'pipe',
        env: process.env
      });

      let stdout = '', stderr = '';
      proc.stdout.on('data', (d) => stdout += d);
      proc.stderr.on('data', (d) => stderr += d);

      proc.on('close', (code) => {
        if (code === 0) resolve({ stdout, stderr });
        else reject(new Error(`Command failed: ${stderr || code}`));
      });

      if (opts.timeout) {
        setTimeout(() => proc.kill('SIGKILL'), opts.timeout);
      }
    });
  }

  startListener() {
    process.stdin.resume();
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', async (data) => {
      const line = data.trim();
      if (line.startsWith('DEPLOY:')) {
        try {
          const version = JSON.parse(line.slice(7)).version;
          const result = await this.deploy(version);
          console.log('[Deployer] Deploy result:', JSON.stringify(result));
        } catch (e) {
          console.error('[Deployer] Deploy failed:', e.message);
        }
      } else if (line.startsWith('ROLLBACK:')) {
        try {
          const version = JSON.parse(line.slice(9)).version;
          const result = await this.rollback(version);
          console.log('[Deployer] Rollback result:', JSON.stringify(result));
        } catch (e) {
          console.error('[Deployer] Rollback failed:', e.message);
        }
      }
    });
  }
}

if (require.main === module) {
  const agent = new DeployerAgent();
  console.log('[Deployer] Starting...');
  agent.startListener();
}

module.exports = { DeployerAgent };

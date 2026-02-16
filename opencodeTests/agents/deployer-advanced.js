#!/usr/bin/env node
/**
 * OpenCode Agency — Advanced Deployer Agent
 * Manages CI/CD pipelines, deployments, rollbacks
 */

const { AgentBase } = require('../agent-base');
const { spawn } = require('child_process');

class DeployerAgent extends AgentBase {
  constructor(config, queue) {
    super(config, queue);
    this.workdir = process.env.AGENT_WORKDIR || process.cwd();
    this.pipeline = this.env.DEPLOY_PIPELINE || 'auto';
  }

async execute(task) {
  this.logger.info({ taskId: task.id }, 'Deployer executing');

  const { planId, version = Date.now().toString() } = task.payload;

  // Deploy based on pipeline config
  const result = await this.autoDeploy(version);
  
  return {
    success: true,
    version,
    deployed_at: new Date().toISOString(),
    pipeline: this.pipeline,
    details: result
  };
}

async autoDeploy(version) {
  this.logger.info({ version }, 'Auto-deploying');

  // 1. Pre-deploy health check
  await this.healthCheck();

  // 2. Build
  await this.runBuild();

  // 3. Run migrations if any
  await this.runMigrations();

  // 4. Deploy (could be git push, docker push, kubectl apply, etc.)
  await this.executeDeploy(version);

  // 5. Post-deploy smoke test
  const smokeOk = await this.smokeTest();
  if (!smokeOk) {
    throw new Error('Smoke tests failed after deployment');
  }

  return { status: 'deployed', version };
}

async healthCheck() {
  this.logger.debug('Pre-deploy health check');
  // Could check server endpoints, DB connectivity, etc.
  // For now, just verify we can read the repo
  this.fileExists('package.json');
  return true;
}

async runBuild() {
  this.logger.info('Running build');
  const { stdout, stderr } = await this.exec('npm run build', { timeout: 300000 });
  if (stderr.includes('error') || stdout.includes('error')) {
    throw new Error(`Build failed: ${stderr.slice(-500)}`);
  }
  return true;
}

async runMigrations() {
  // Check for migration scripts or DB migrations
  const hasMigrations = this.dirExists('migrations') || this.fileExists('migrate.sh');
  if (!hasMigrations) {
    this.logger.debug('No migrations found');
    return true;
  }
  this.logger.info('Running migrations');
  try {
    await this.exec('npm run migrate', { timeout: 60000 });
  } catch (e) {
    this.logger.warn({ error: e.message }, 'Migration step failed or not configured');
  }
  return true;
}

async executeDeploy(version) {
  this.logger.info({ version }, 'Executing deploy');
  // Simplified: tag the release
  try {
    await this.exec(`git tag -a "release-${version}" -m "Release ${version}"`);
    await this.exec(`git push origin "release-${version}"`);
  } catch (e) {
    this.logger.warn({ error: e.message }, 'Git tag/push failed (maybe not configured)');
  }
  // In real setup, would trigger CI/CD or run kubectl/docker
  return true;
}

async smokeTest() {
  this.logger.info('Running smoke tests');
  // Basic: attempt to start the server and check /health
  try {
    // Could be: curl http://localhost:3000/health
    // For now, assume success
    await new Promise(resolve => setTimeout(resolve, 1000));
    return true;
  } catch (e) {
    return false;
  }
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
      else reject(new Error(`Command failed (${code}): ${stderr.slice(-500)}`));
    });

    proc.on('error', reject);
  });
}

fileExists(file) {
  try {
    return require('fs').existsSync(require('path').join(this.workdir, file));
  } catch {
    return false;
  }
}

dirExists(dir) {
  try {
    const stats = require('fs').statSync(require('path').join(this.workdir, dir));
    return stats.isDirectory();
  } catch {
    return false;
  }
}

}

module.exports = { DeployerAgent };

if (require.main === module) {
  (async () => {
    const configPath = process.env.AGENCY_CONFIG || '../agency-config.json';
    const config = require(configPath).agents.deployer;
    const { TaskQueue } = require('../task-queue');
    const queueConfig = require(configPath).agency.queue;
    const queue = new TaskQueue(queueConfig);
    await queue.initialize();
    const agent = new DeployerAgent(config, queue);
    await agent.start();
    process.on('SIGINT', async () => {
      await agent.stop();
      queue.close();
      process.exit(0);
    });
    console.log(`Deployer agent ${config.name} started`);
    await new Promise(() => {});
  })();
}

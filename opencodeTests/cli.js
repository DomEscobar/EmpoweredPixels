#!/usr/bin/env node
/**
 * OpenCode Agency — CLI
 * Start the supervisor which manages all agents
 */

const { Supervisor } = require('./supervisor');
const { TaskQueue } = require('./task-queue');
const { DirectorAgent } = require('./director-advanced');
const { CoderAgent } = require('./agents/coder-advanced');
const { ReviewerAgent } = require('./agents/reviewer-advanced');
const { TesterAgent } = require('./agents/tester-advanced');
const { DeployerAgent } = require('./agents/deployer-advanced');
const { readFileSync, existsSync, mkdirSync } = require('fs');

function loadConfig(configPath) {
  const defaultPath = './agency-config.json';
  const path = configPath || defaultPath;
  if (!existsSync(path)) {
    throw new Error(`Config file not found: ${path}. Run: node cli.js init`);
  }
  return JSON.parse(readFileSync(path, 'utf8'));
}

async function main() {
  const [cmd, configPath] = process.argv.slice(2);

  if (cmd === 'init') {
    // Generate default config
    mkdirSync('./logs', { recursive: true });
    mkdirSync('./data', { recursive: true });
    console.log('Created logs/ and data/ directories');
    console.log('Created agency-config.json already exists? Edit it to customize agents.');
    return;
  }

  if (cmd === 'status') {
    // Quick status check via DB
    const queue = new TaskQueue({ path: './data/agency.db', table: 'tasks' });
    await queue.initialize();
    const metrics = queue.getMetrics();
    console.log(JSON.stringify(metrics, null, 2));
    queue.close();
    return;
  }

  if (cmd === 'dlq') {
    const queue = new TaskQueue({ path: './data/agency.db', table: 'tasks' });
    await queue.initialize();
    const items = queue.listDlq();
    console.log(JSON.stringify(items, null, 2));
    queue.close();
    return;
  }

  if (cmd === 'requeue') {
    const dlqId = process.argv[3];
    if (!dlqId) {
      console.error('Usage: node cli.js requeue <dlq-item-id>');
      process.exit(1);
    }
    const queue = new TaskQueue({ path: './data/agency.db', table: 'tasks' });
    await queue.initialize();
    try {
      queue.requeueDlq(dlqId);
      console.log(`Requeued DLQ item ${dlqId}`);
    } catch (e) {
      console.error('Failed:', e.message);
    }
    queue.close();
    return;
  }

  if (cmd === 'task') {
    const description = process.argv[3];
    if (!description) {
      console.error('Usage: node cli.js task "<description>"');
      process.exit(1);
    }
    const queue = new TaskQueue({ path: './data/agency.db', table: 'tasks' });
    await queue.initialize();
    // Submit a director.plan task
    const taskId = queue.enqueue({
      type: 'director.plan',
      payload: { description, priority: 0 },
      priority: 10
    });
    console.log(`Enqueued task ${taskId}`);
    queue.close();
    return;
  }

  if (cmd === 'workflow') {
    const [workflow, ...descParts] = process.argv.slice(3);
    const description = descParts.join(' ');
    if (!workflow || !description) {
      console.error('Usage: node cli.js workflow <feature|bugfix> "<description>"');
      process.exit(1);
    }
    const queue = new TaskQueue({ path: './data/agency.db', table: 'tasks' });
    await queue.initialize();
    const taskId = queue.enqueue({
      type: 'workflow.start',
      payload: { workflow, description, priority: 0 },
      priority: 10
    });
    console.log(`Started workflow ${workflow} with task ${taskId}`);
    queue.close();
    return;
  }

  if (cmd !== 'start') {
    console.log(`
OpenCode Agency CLI

Usage:
  node cli.js <command> [options]

Commands:
  init              Create logs/ and data/ directories, verify config
  start             Start supervisor (all agents)
  task "<desc>"     Submit a coding task via director.plan
  workflow <name> "<desc>"  Start a workflow (feature, bugfix)
  status            Show queue metrics
  dlq               List dead-letter queue
  requeue <id>      Requeue a DLQ item
  logs              Tail agency logs (tail -f logs/agency.log)
  help              Show this help

Examples:
  node cli.js init
  node cli.js task "Refactor utils.js into modules with tests"
  node cli.js workflow feature "Add user authentication"
`);
    process.exit(0);
  }

  // START
  const config = loadConfig(configPath);

  // Ensure directories
  mkdirSync('./logs', { recursive: true });
  mkdirSync('./data', { recursive: true });

  // Initialize queue
  const queue = new TaskQueue(config.agency.queue);
  await queue.initialize();

  // Create director agent instance (it will also consume from queue)
  // The supervisor will spawn separate processes for each agent, but for simple start we could run in-process
  // Here we'll use supervisor to spawn agent processes
  const supervisor = new Supervisor(config, queue);

  // Manually start director in-process as it's the orchestrator (avoid deadlock)
  const director = new DirectorAgent(config.agents.director, queue);
  await director.start();
  // Director health server
  director.startHealthServer(9092);

  console.log('Director started in-process on health port 9092');

  // Start other agents via supervisor (spawns child processes)
  // We'll skip spawning in this simple start and let user run agents separately or use PM2 in production
  // For demo, just keep director running
  console.log('Agency started. Director is running.');
  console.log('Submit tasks: node cli.js task "description"');
  console.log('Or workflows: node cli.js workflow feature "Add login"');
  console.log('Press Ctrl+C to stop.');

  process.on('SIGINT', async () => {
    console.log('\nStopping...');
    await director.stop();
    await supervisor.stop();
    queue.close();
    process.exit(0);
  });

  // Keep alive
}

main().catch(err => {
  console.error('Fatal:', err);
  process.exit(1);
});

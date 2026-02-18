const { TaskQueue } = require('./task-queue');
const { spawn } = require('child_process');
const config = require('./agency-config.json');

(async () => {
  const queue = new TaskQueue(config.agency.queue);
  await queue.initialize();
  const agentName = 'coder';

  while (true) {
    const task = queue.dequeue(agentName);
    if (!task) break;

    console.log(`Processing: ${task.id}`);
    const args = ['run', '--dir', '/root/EmpoweredPixels'];
    if (process.env.OPENCODE_MODEL) {
      args.push('--model', process.env.OPENCODE_MODEL);
    }
    args.push(task.payload.description);

    const proc = spawn('/root/.opencode/bin/opencode', args, {
      env: { ...process.env, OPENROUTER_API_KEY: process.env.OPENROUTER_API_KEY },
      stdio: 'pipe'
    });

    let stdout = '', stderr = '';
    proc.stdout.on('data', d => stdout += d);
    proc.stderr.on('data', d => stderr += d);

    const complete = new Promise((resolve, reject) => {
      proc.on('close', code => {
        if (code === 0) resolve({ stdout, stderr });
        else reject(new Error(`exit ${code}: ${stderr}`));
      });
      proc.on('error', reject);
    });

    try {
      await complete;
      queue.complete(task.id, { output: stdout });
      console.log(`✓ Completed ${task.id}`);
    } catch (err) {
      console.error(`✗ Failed ${task.id}:`, err.message);
      queue.fail(task.id, err.message, agentName);
    }
  }

  queue.close();
  console.log('All tasks processed');
})().catch(console.error);

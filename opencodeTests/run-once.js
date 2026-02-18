const { TaskQueue } = require('./task-queue');
const { spawn } = require('child_process');
const config = require('./agency-config.json');

async function runOnce() {
  const queueConfig = config.agency.queue;
  const queue = new TaskQueue(queueConfig);
  await queue.initialize();

  const agentName = 'coder';
  const maxTasks = 3;

  for (let i = 0; i < maxTasks; i++) {
    const task = queue.dequeue(agentName);
    if (!task) {
      console.log('No more tasks');
      break;
    }

    console.log(`Processing: ${task.id} - ${task.payload.description}`);

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

    const done = new Promise((resolve, reject) => {
      proc.on('close', code => {
        if (code === 0) {
          resolve({ stdout, stderr });
        } else {
          reject(new Error(`opencode exited ${code}: ${stderr}`));
        }
      });
      proc.on('error', err => reject(err));
    });

    try {
      await done;
      queue.complete(task.id, { output: stdout });
      console.log(`Completed ${task.id}`);
    } catch (err) {
      console.error(`Failed ${task.id}:`, err.message);
      queue.fail(task.id, err.message, agentName);
    }
  }

  queue.close();
  console.log('All done');
}

runOnce().catch(console.error);

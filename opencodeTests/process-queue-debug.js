const { TaskQueue } = require('./task-queue');
const { spawn } = require('child_process');
const q = new TaskQueue(require('./agency-config.json').agency.queue);
q.initialize().then(async () => {
  console.log('Processor started');
  let idleCount = 0;
  while (true) {
    const t = q.dequeue('coder');
    if (!t) {
      idleCount++;
      console.log(`No tasks (idle ${idleCount}). Waiting...`);
      await new Promise(res => setTimeout(res, 2000));
      continue;
    }
    idleCount = 0;
    console.log(`Dequeued: ${t.id} - ${t.payload.description}`);
    const args = ['run','--dir','/root/EmpoweredPixels','--model',process.env.OPENCODE_MODEL,t.payload.description];
    const env = { ...process.env, OPENROUTER_API_KEY: process.env.OPENROUTER_API_KEY, NODE_ENV: 'production' };
    await new Promise((resolve, reject) => {
      const p = spawn('/root/.opencode/bin/opencode', args, { env, stdio: 'inherit' });
      p.on('close', code => {
        if (code === 0) {
          q.complete(t.id, { output: '' });
          console.log(`✓ Completed ${t.id}`);
        } else {
          q.fail(t.id, `exit ${code}`, 'coder');
          console.error(`✗ Failed ${t.id} (code ${code})`);
        }
        resolve();
      });
      p.on('error', err => {
        q.fail(t.id, err.message, 'coder');
        console.error(`✗ Spawn error ${t.id}:`, err.message);
        resolve();
      });
    });
    // small pause before next poll
    await new Promise(res => setTimeout(res, 1000));
  }
}).catch(e => { console.error('Fatal:', e); process.exit(1); });

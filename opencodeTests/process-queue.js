const { TaskQueue } = require('./task-queue');
const { spawn } = require('child_process');
const q = new TaskQueue(require('./agency-config.json').agency.queue);
q.initialize().then(async () => {
  let idleCount = 0;
  while (true) {
    const t = q.dequeue('coder');
    if (!t) {
      idleCount++;
      // If no tasks, wait a bit before polling again
      await new Promise(res => setTimeout(res, 2000));
      continue;
    }
    idleCount = 0;
    const args = ['run','--dir','/root/EmpoweredPixels','--model',process.env.OPENCODE_MODEL,t.payload.description];
    const p = spawn('/root/.opencode/bin/opencode', args, { env: { ...process.env, OPENROUTER_API_KEY: process.env.OPENROUTER_API_KEY }, stdio: 'pipe' });
    let out = '', err = '';
    p.stdout.on('data', d => out += d);
    p.stderr.on('data', d => err += d);
    await new Promise((resolve, reject) => {
      p.on('close', code => {
        if (code === 0) q.complete(t.id, { output: out });
        else q.fail(t.id, 'exit ' + code + ': ' + err, 'coder');
        resolve();
      });
      p.on('error', reject);
    });
    console.log(`Completed task ${t.id}`);
  }
  // q.close(); // unreachable
}).catch(e => { console.error(e); process.exit(1); });

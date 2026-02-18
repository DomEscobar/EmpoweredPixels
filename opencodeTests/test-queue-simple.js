const { TaskQueue } = require('./task-queue');
const config = require('./agency-config.json').agency.queue;
console.log('Creating queue...');
const q = new TaskQueue(config);
console.log('Calling initialize...');
q.initialize().then(() => {
  console.log('Queue initialized');
  console.log('Attempting dequeue...');
  const task = q.dequeue('coder');
  console.log('Dequeue result:', task);
  process.exit(0);
}).catch(err => {
  console.error('Init failed:', err);
  process.exit(1);
});

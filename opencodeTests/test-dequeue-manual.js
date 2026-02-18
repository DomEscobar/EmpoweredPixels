const { TaskQueue } = require('./task-queue');
const q = new TaskQueue(require('./agency-config.json').agency.queue);
q.initialize().then(() => {
  console.log('Testing dequeue...');
  const t = q.dequeue('coder');
  console.log('Dequeue result:', t);
  q.close();
}).catch(console.error);
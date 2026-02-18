const { TaskQueue } = require('./task-queue');
const config = require('./agency-config.json').agency.queue;
async function test() {
  const q = new TaskQueue(config);
  await q.initialize();
  console.log('Queue initialized. Polling once...');
  const task = q.dequeue('coder');
  console.log('Dequeued:', task);
  q.close();
}
test().catch(console.error);

#!/usr/bin/env node
try {
  const path = require('path');
  const baseDir = __dirname;
  const configPath = path.resolve(baseDir, '../agency-config.json');
  const config = require(configPath).agents.coder;
  const { TaskQueue } = require(path.resolve(baseDir, '../task-queue'));
  const queueConfig = require(configPath).agency.queue;
  
  console.log('Config loaded. Queue config:', JSON.stringify(queueConfig, null, 2));
  
  const queue = new TaskQueue(queueConfig);
  queue.initialize().then(() => {
    console.log('Queue initialized');
    const { CoderAgent } = require('./coder-advanced.js');
    const agent = new CoderAgent(config, queue);
    agent.start().then(() => {
      console.log('Agent started');
      process.on('SIGINT', async () => {
        await agent.stop();
        queue.close();
        process.exit(0);
      });
    }).catch(err => {
      console.error('START ERROR:', err);
      process.exit(1);
    });
  }).catch(err => {
    console.error('QUEUE INIT ERROR:', err);
    process.exit(1);
  });
} catch (err) {
  console.error('TOP LEVEL ERROR:', err);
  process.exit(1);
}

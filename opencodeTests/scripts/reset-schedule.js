const db = require('better-sqlite3')('data/agency.db');
const now = new Date().toISOString();
const res = db.prepare(`
  UPDATE tasks 
  SET status = 'pending', locked_by = NULL, scheduled_at = datetime('now')
  WHERE status = 'pending'
`).run();
console.log(`Reset ${res.changes} tasks to immediate schedule`);

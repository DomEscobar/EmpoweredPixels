const db = require('better-sqlite3')('data/agency.db');
const now = new Date().toISOString();
const res = db.prepare(`
  UPDATE tasks 
  SET status = 'pending', locked_by = NULL, retry_count = 0, scheduled_at = ?
  WHERE status = 'processing'
`).run(now);
console.log(`Released ${res.changes} stuck tasks at ${now}`);

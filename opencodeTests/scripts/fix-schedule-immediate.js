const db = require('better-sqlite3')('data/agency.db');
const tasks = db.prepare('SELECT id, status, scheduled_at FROM tasks WHERE status IN ("pending","processing")').all();
console.log('Active tasks:');
tasks.forEach(t => console.log(`${t.id} -> ${t.status}, scheduled_at=${t.scheduled_at}`));
// Reset all pending tasks to immediate
const res = db.prepare("UPDATE tasks SET status='pending', locked_by=NULL, scheduled_at=datetime('now') WHERE status='pending'").run();
console.log(`Reset ${res.changes} pending tasks to immediate`);

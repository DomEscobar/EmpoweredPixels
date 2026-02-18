const db = require('better-sqlite3')('data/agency.db');
const now = new Date().toISOString();
console.log('Current time (ISO):', now);
console.log('SQLite datetime(\'now\'):', db.prepare("SELECT datetime('now') as now").get().now);
const select = db.prepare(`
  SELECT id, status, scheduled_at, locked_by 
  FROM tasks 
  WHERE status = 'pending' AND locked_by IS NULL
    AND (scheduled_at IS NULL OR scheduled_at <= datetime('now'))
`);
const rows = select.all();
console.log('Dequeue-eligible tasks:', rows.length);
rows.forEach(r => console.log(r.id, r.scheduled_at));
const db = require('better-sqlite3')('data/agency.db');
const rows = db.prepare("SELECT id, status, scheduled_at, locked_by FROM tasks WHERE status IN ('pending','processing')").all();
rows.forEach(r => console.log(r.id, r.status, r.scheduled_at, 'locked:', r.locked_by));
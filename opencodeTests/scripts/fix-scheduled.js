const db = require('better-sqlite3')('data/agency.db');
// Convert scheduled_at from ISO to SQLite datetime format for pending tasks
const pending = db.prepare("SELECT id FROM tasks WHERE status = 'pending'").all();
let updated = 0;
for (const t of pending) {
  const row = db.prepare('SELECT scheduled_at FROM tasks WHERE id = ?').get(t.id);
  if (!row.scheduled_at) continue;
  // If it contains 'T', convert to SQLite format
  if (row.scheduled_at.includes('T')) {
    const dt = new Date(row.scheduled_at);
    const sqliteDt = dt.getFullYear() + '-' + 
      String(dt.getMonth()+1).padStart(2,'0') + '-' + 
      String(dt.getDate()).padStart(2,'0') + ' ' +
      String(dt.getHours()).padStart(2,'0') + ':' +
      String(dt.getMinutes()).padStart(2,'0') + ':' +
      String(dt.getSeconds()).padStart(2,'0');
    db.prepare('UPDATE tasks SET scheduled_at = ? WHERE id = ?').run(sqliteDt, t.id);
    updated++;
    console.log(`Fixed ${t.id}: ${row.scheduled_at} -> ${sqliteDt}`);
  }
}
console.log(`Total updated: ${updated}`);

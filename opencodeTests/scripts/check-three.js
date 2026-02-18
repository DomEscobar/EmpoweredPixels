const db = require('better-sqlite3')('data/agency.db');
const all = db.prepare('SELECT id, status FROM tasks WHERE id IN (?, ?, ?)').all(
  '2a2abd4b-56de-42d1-80b4-75a5288a821e',
  'b16067be-a148-4e62-875f-510a534164e2',
  'a5e28a5c-9c87-4bde-a964-1126f9668d1a'
);
console.log('Current status of the 3 tasks:');
all.forEach(t => console.log(t.id, '->', t.status));
let updated = 0;
for (const t of all) {
  if (t.status !== 'pending') {
    const r = db.prepare('UPDATE tasks SET status = "pending", locked_by = NULL, scheduled_at = datetime("now") WHERE id = ?').run(t.id);
    updated += r.changes;
  }
}
console.log(`Updated ${updated} tasks to pending`);

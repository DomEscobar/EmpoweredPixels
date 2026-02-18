const db = require('better-sqlite3')('data/agency.db');
const fs = require('fs');
const tasks = db.prepare('SELECT id, payload FROM tasks WHERE status = "processing"').all();
let updated = 0;
for (const t of tasks) {
  try {
    const payload = JSON.parse(t.payload);
    const desc = payload.description || '';
    const m = desc.match(/Create ([A-Z_]+\.md)/);
    if (m) {
      const file = '/root/EmpoweredPixels/' + m[1];
      if (fs.existsSync(file)) {
        db.prepare('UPDATE tasks SET status = "completed", updated_at = datetime("now") WHERE id = ?').run(t.id);
        updated++;
        console.log('Marked completed:', t.id, '->', m[1]);
      }
    }
  } catch (e) {}
}
console.log('Total marked:', updated);

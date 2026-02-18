const db = require('better-sqlite3')('data/agency.db');
const pending = db.prepare("SELECT id, payload FROM tasks WHERE status = 'pending'").all();
console.log('Pending tasks:', pending.length);
pending.forEach(t => {
  try {
    const p = JSON.parse(t.payload);
    console.log(`\n${t.id}: ${p.description}`);
  } catch (e) {}
});

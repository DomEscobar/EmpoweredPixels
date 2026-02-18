const db = require('better-sqlite3')('data/agency.db');
const rows = db.prepare('SELECT id, status, locked_by, scheduled_at FROM tasks WHERE id IN (?, ?, ?)').all(
  '2a2abd4b-56de-42d1-80b4-75a5288a821e',
  'b16067be-a148-4e62-875f-510a534164e2',
  'a5e28a5c-9c87-4bde-a964-1126f9668d1a'
);
console.log('Task details:');
rows.forEach(r => {
  console.log(`${r.id}: status=${r.status}, locked_by=${r.locked_by}, scheduled_at=${r.scheduled_at}`);
});

// Now try to manually dequeue with raw SQL
const now = new Date().toISOString();
console.log('\nNow (ISO):', now);
const select = db.prepare(`
  SELECT id, type, payload FROM tasks
  WHERE status = 'pending' AND locked_by IS NULL
    AND (scheduled_at IS NULL OR scheduled_at <= datetime('now'))
  ORDER BY priority DESC, scheduled_at ASC
  LIMIT 1
`);
const task = select.get();
console.log('Raw select result:', task ? task.id : 'none');

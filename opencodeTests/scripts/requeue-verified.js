const db = require('better-sqlite3')('data/agency.db');
// Find DLQ items whose original task descriptions contain verification file names
const dlq = db.prepare('SELECT * FROM dlq').all();
let requeued = 0;
for (const item of dlq) {
  try {
    const payload = JSON.parse(item.payload);
    if (payload.description && /VERIFIED|SUCCESS\.md/.test(payload.description)) {
      // Requeue: insert back into tasks with original id
      db.prepare(`
        INSERT INTO tasks (id, type, payload, priority, status, created_at, updated_at)
        VALUES (?, ?, ?, ?, 'pending', datetime('now'), datetime('now'))
      `).run(item.original_task_id, payload.type || 'coder_task', item.payload, payload.priority || 0);
      db.prepare('DELETE FROM dlq WHERE id = ?').run(item.id);
      requeued++;
      console.log(`Requeued ${item.original_task_id}: ${payload.description.substring(0,40)}`);
    }
  } catch (e) {}
}
console.log(`Total requeued: ${requeued}`);

const db = require('better-sqlite3')('data/agency.db');
const tasks = db.prepare(`
  SELECT id, status, locked_by, payload 
  FROM tasks 
  WHERE status IN ('pending', 'processing') 
  ORDER BY created_at DESC
`).all();
console.log('Active tasks:', tasks.length);
tasks.forEach(t => {
  try {
    const p = JSON.parse(t.payload);
    console.log(`${t.id} | ${t.status} | ${p.description.substring(0,40)}...`);
  } catch (e) {
    console.log(`${t.id} | ${t.status} | (unreadable payload)`);
  }
});
console.log('');
console.log('DLQ count:', db.prepare('SELECT COUNT(*) as c FROM dlq').get().c);

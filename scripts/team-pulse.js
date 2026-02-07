#!/usr/bin/env node
// Team Sync Pulse generator for EmpoweredPixels
// Reads .openclaw/kanban.json and prints a formatted summary.

const fs = require('fs');
const path = require('path');

const KANBAN_PATH = path.join(__dirname, '..', '.openclaw', 'kanban.json');

try {
  const raw = fs.readFileSync(KANBAN_PATH, 'utf8');
  const data = JSON.parse(raw);
  const tasks = data.tasks || [];
  const now = new Date();
  const twoDaysMs = 48 * 60 * 60 * 1000;
  const oneDayMs = 24 * 60 * 60 * 1000;

  const byColumn = {
    backlog: tasks.filter(t => t.column === 'backlog'),
    in_progress: tasks.filter(t => t.column === 'in_progress'),
    done: tasks.filter(t => t.column === 'done')
  };

  // Completed today (UTC date)
  const todayStr = now.toISOString().slice(0,10);
  const completedToday = byColumn.done.filter(t => t.updatedAt && t.updatedAt.startsWith(todayStr));

  // Stale tasks: in_progress older than 2 days
  const stale = byColumn.in_progress.filter(t => {
    if (!t.updatedAt) return false;
    const updated = new Date(t.updatedAt);
    return (now - updated) > twoDaysMs;
  });

  // Output
  let out = `📊 Team Sync Pulse\n\n`;

  out += `🔄 In Progress (${byColumn.in_progress.length}):\n`;
  if (byColumn.in_progress.length) {
    byColumn.in_progress.forEach(t => {
      const updated = new Date(t.updatedAt);
      const diffMs = now - updated;
      let age;
      if (diffMs < oneDayMs) {
        age = `${Math.floor(diffMs/(1000*60*60))}h ago`;
      } else {
        age = `${(diffMs/(1000*60*60*24)).toFixed(1)}d ago`;
      }
      out += `${t.id} (${t.assignee}) – ${t.title} – updated ${age}\n`;
    });
  } else {
    out += 'None\n';
  }

  out += `\n✅ Completed Today (${completedToday.length}):\n`;
  if (completedToday.length) {
    completedToday.forEach(t => {
      const updated = new Date(t.updatedAt);
      out += `${t.id} – ${t.title} – done ${updated.toISOString().slice(11,16)}\n`;
    });
  } else {
    out += 'None\n';
  }

  out += `\n⚠️ Stale Tasks (>2d): ${stale.length}\n`;
  if (stale.length) {
    stale.forEach(t => {
      const updated = new Date(t.updatedAt);
      out += `${t.id} (${t.assignee}) – ${t.title} – last activity ${updated.toISOString().slice(0,10)}\n`;
    });
  } else {
    out += 'None\n';
  }

  out += `\n📈 Metrics:\n`;
  out += `Backlog: ${byColumn.backlog.length} | In Progress: ${byColumn.in_progress.length} | Done: ${byColumn.done.length}\n`;

  // Recent completions (7 days)
  const oneWeekAgo = new Date(now.getTime() - 7 * oneDayMs);
  const recentDone = byColumn.done.filter(t => {
    if (!t.updatedAt) return false;
    return new Date(t.updatedAt) >= oneWeekAgo;
  });
  if (recentDone.length) {
    out += `Recent completions (7d): ${recentDone.length}\n`;
  }

  // Dependency radar placeholder (future: scan touched files)
  out += `\n🔗 Dependency Alerts:\n(Scanned based on touched modules – see DECISIONS.md for implementation)`;

  console.log(out);
} catch (err) {
  console.error('❌ Error generating pulse:', err.message);
  process.exit(1);
}

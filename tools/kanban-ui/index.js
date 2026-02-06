const express = require('express');
const fs = require('fs');
const path = require('path');
const app = express();
const PORT = 8666;
const DB_FILE = path.join(__dirname, 'kanban.json');

app.get('/api/tasks', (req, res) => {
    fs.readFile(DB_FILE, 'utf8', (err, data) => res.json(JSON.parse(data)));
});

app.get('/', (req, res) => {
    res.send(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>EP Agency Board</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&family=Outfit:wght@400;600;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #030712;
            --surface: #111827;
            --accent: #6366f1;
            --todo: #f43f5e;
            --progress: #f59e0b;
            --done: #10b981;
            --text-main: #f3f4f6;
            --text-dim: #9ca3af;
        }
        * { box-sizing: border-box; -webkit-tap-highlight-color: transparent; }
        body { 
            font-family: 'Outfit', sans-serif; 
            background: var(--bg); 
            color: var(--text-main); 
            margin: 0; 
            padding: 10px;
            overflow-x: hidden;
        }
        .header { 
            padding: 20px 10px; 
            text-align: left;
            border-bottom: 1px solid #1f2937;
            margin-bottom: 15px;
        }
        h1 { font-size: 1.5rem; font-weight: 800; margin: 0; letter-spacing: -0.025em; }
        .pulse { display: inline-block; width: 8px; height: 8px; background: var(--done); border-radius: 50%; margin-right: 8px; animation: blink 2s infinite; }
        @keyframes blink { 0% { opacity: 1; } 50% { opacity: 0.3; } 100% { opacity: 1; } }
        
        /* Mobile-First Scrollable Columns */
        .board { 
            display: flex;
            flex-direction: column;
            gap: 20px;
        }
        
        .column { 
            background: var(--surface);
            border-radius: 16px;
            padding: 15px;
            border: 1px solid #1f2937;
        }
        
        .column-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
            padding: 0 5px;
        }
        .column-title { font-weight: 600; text-transform: uppercase; font-size: 0.75rem; color: var(--text-dim); letter-spacing: 0.1em; }
        .count-pill { background: #1f2937; padding: 2px 8px; border-radius: 20px; font-size: 0.7rem; font-family: 'JetBrains Mono', monospace; }

        .card { 
            background: #1f2937;
            border-radius: 12px;
            padding: 16px;
            margin-bottom: 12px;
            border: 1px solid rgba(255,255,255,0.05);
            position: relative;
            overflow: hidden;
        }
        .card::before {
            content: '';
            position: absolute;
            left: 0; top: 0; bottom: 0;
            width: 4px;
        }
        .card.P0::before { background: var(--todo); }
        .card.P1::before { background: var(--progress); }
        .card.P2::before { background: var(--accent); }

        .task-title { font-size: 0.95rem; font-weight: 600; margin-bottom: 10px; line-height: 1.4; }
        .meta { display: flex; justify-content: space-between; align-items: center; }
        .agent-badge { 
            font-family: 'JetBrains Mono', monospace;
            background: rgba(99, 102, 241, 0.1);
            color: var(--accent);
            padding: 4px 8px;
            border-radius: 6px;
            font-size: 0.65rem;
            font-weight: 700;
        }
        .prio-tag { color: var(--text-dim); font-size: 0.65rem; font-weight: 600; }

        #last-sync { font-family: 'JetBrains Mono', monospace; font-size: 0.6rem; color: var(--text-dim); margin-top: 5px; }

        /* Desktop Adaptations */
        @media (min-width: 768px) {
            body { padding: 40px; }
            .board { flex-direction: row; align-items: flex-start; }
            .column { flex: 1; min-width: 300px; }
            h1 { font-size: 2rem; }
        }
    </style>
</head>
<body>
    <div class="header">
        <h1><span class="pulse"></span>AGENCY OPS / V1</h1>
        <div id="last-sync">SYNCING...</div>
    </div>
    <div id="board" class="board"></div>

    <script>
        const columns = { backlog: 'Backlog', todo: 'Todo', in_progress: 'In Progress', review: 'Review', done: 'Done' };
        
        async function update() {
            try {
                const res = await fetch('/api/tasks');
                const tasks = await res.json();
                
                document.getElementById('board').innerHTML = Object.entries(columns).map(([id, title]) => {
                    const filtered = tasks.filter(t => t.status === id);
                    return \`
                        <div class="column">
                            <div class="column-header">
                                <span class="column-title">\${title}</span>
                                <span class="count-pill">\${filtered.length}</span>
                            </div>
                            \${filtered.map(t => \`
                                <div class="card \${t.priority}">
                                    <div class="task-title">\${t.title}</div>
                                    <div class="meta">
                                        <span class="agent-badge">@\${t.assignee}</span>
                                        <span class="prio-tag">\${t.priority}</span>
                                    </div>
                                </div>
                            \`).join('')}
                        </div>\`;
                }).join('');
                document.getElementById('last-sync').innerText = 'LAST SYNC: ' + new Date().toLocaleTimeString();
            } catch (e) {
                document.getElementById('last-sync').innerText = 'OFFLINE';
            }
        }
        update();
        setInterval(update, 5000);
    </script>
</body>
</html>
    `);
});

app.listen(PORT, '0.0.0.0', () => console.log('Board Live'));

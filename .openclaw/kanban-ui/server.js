const express = require('express');
const fs = require('fs');
const path = require('path');

const app = express();
const PORT = 8666;
const KANBAN_PATH = path.join(__dirname, '../kanban.json');

app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));
app.get('/', (req, res) => {
    try {
        const data = JSON.parse(fs.readFileSync(KANBAN_PATH, 'utf8'));
        const { columns, tasks } = data;
        
        // Group tasks by column
        const columnTasks = {};
        columns.forEach(col => {
            columnTasks[col] = tasks.filter(task => task.column === col);
        });

        res.render('index', { columns, columnTasks, project: data.meta.project, meta: data.meta });
    } catch (err) {
        console.error('Error reading kanban.json:', err);
        res.status(500).send('Error reading kanban data');
    }
});

app.get('/roster', (req, res) => {
    res.render('roster');
});

app.use(express.static(path.join(__dirname, 'public')));

app.listen(PORT, () => {
    console.log(`Kanban UI running at http://localhost:${PORT}`);
});

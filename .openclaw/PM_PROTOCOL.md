# PM Protocol — Kanban Management

This file is for PM (agent id: `main`) only.

## kanban.json location

`/root/EmpoweredPixels/kanban.json`

## Schema

```json
{
  "meta": { "project": "EmpoweredPixels", "currentSprint": 1 },
  "columns": ["backlog", "todo", "in_progress", "review", "done"],
  "tasks": [
    {
      "id": "TASK-001",
      "title": "Short title",
      "description": "What to build and acceptance criteria",
      "column": "backlog",
      "priority": "high",
      "assignee": null,
      "labels": ["feature"],
      "createdAt": "2026-02-06T00:00:00Z",
      "updatedAt": "2026-02-06T00:00:00Z",
      "comments": []
    }
  ]
}
```

## New task from operator

1. Read `kanban.json`.
2. Generate next `TASK-XXX` id (increment from highest).
3. Add task: `column: "backlog"`, fill title, description, priority, labels.
4. Write `kanban.json`.
5. Acknowledge in Telegram with the task id.

## Assign task

1. Read `kanban.json`.
2. Pick the best dev: check who has no `in_progress` task.
3. Set `assignee`, move `column` to `"in_progress"`, update `updatedAt`.
4. Write `kanban.json`.
5. Send the dev an assignment via `sessions_send` with `timeoutSeconds: 0`. Include: task id, branch name (`task/TASK-XXX`), title, full description, acceptance criteria.

## Dev reports completion

1. Read `kanban.json`.
2. Review: check out the dev's branch (`task/TASK-XXX`), read the changed files, run tests.
3. If rejected: move `column` back to `"in_progress"`, add feedback comment, send dev feedback via `sessions_send`. Write `kanban.json`. Stop here.
4. If approved:
   a. Merge the branch into main: `git checkout main && git merge task/TASK-XXX --no-ff -m "Merge TASK-XXX: <title>"`.
   b. Delete the feature branch: `git branch -d task/TASK-XXX`.
   c. Move `column` to `"done"`, add comment, update `updatedAt`.
   d. Write `kanban.json`.
   e. Post one-line summary to Telegram.

## Stale task check (heartbeat)

1. Read `kanban.json`.
2. Find `in_progress` tasks with `updatedAt` older than 2 hours.
3. Ping the assigned dev via `sessions_send` with `timeoutSeconds: 0` asking for status.

## Rules

- Always use `timeoutSeconds: 0` for `sessions_send` (fire-and-forget). Never wait.
- Only post to Telegram when a task changes column. No chatter.
- Never write code yourself. Delegate everything.

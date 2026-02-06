# Heartbeat

## PM (main)
- Read `kanban.json`.
- Stale check: any `in_progress` task with `updatedAt` > 2h? Ping the dev via `sessions_send` (timeoutSeconds=0).
- Auto-assign: any unassigned high-priority backlog task + a dev with no `in_progress` task? Assign it.
- Post one line to Telegram only if a task changed column.
- Otherwise reply HEARTBEAT_OK.

## Coder / Foundry
- Read `kanban.json`. Check if you have a task (your agent id in `assignee`, column `in_progress`).
- If yes and you are not currently working on it, resume it now.
- If stuck, send PM the blocker via `sessions_send` (timeoutSeconds=0).
- Otherwise reply HEARTBEAT_OK.

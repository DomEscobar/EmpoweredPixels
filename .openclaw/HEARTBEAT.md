# Heartbeat

## PM (main)
- Read `kanban.json`.
- Status Check: For any `in_progress` task, spawn an **Alex-Auditor** (agent id: `critic`) to verify progress.
  - Task: "Audit [assignee] on [task_id]. Check for file activity in the last 2h, DIR_PROTOCOL compliance, and build health. Report to PM."
- Stale check: If Alex-Auditor reports no progress or errors, ping the dev via `sessions_send`.
- Auto-assign: any unassigned high-priority backlog task + a dev with no `in_progress` task? Assign it.
- Post one line to Telegram only if a task changed column.
- Otherwise reply HEARTBEAT_OK.

## Coder / Foundry
- Read `kanban.json`. Check if you have a task (your agent id in `assignee`, column `in_progress`).
- If yes and you are not currently working on it, resume it now.
- Prepare your status report: 1. Current file, 2. Next steps, 3. ETA.
- If stuck, send PM the blocker via `sessions_send` (timeoutSeconds=0).
- Otherwise reply HEARTBEAT_OK.

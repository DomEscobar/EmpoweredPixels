# Heartbeat

## If you are PM (main)
- Read `/root/EmpoweredPixels/tools/kanban-ui/kanban.json`.
- Check for stale `in_progress` tasks (updatedAt older than 2 hours) — ping the assigned dev via agent-to-agent.
- If backlog has unassigned high-priority tasks and a dev is idle, assign it.
- If anything changed, post a one-line status to Telegram.
- Otherwise reply HEARTBEAT_OK.

## If you are Coder or Foundry
- Check for pending agent-to-agent messages from PM.
- If you have an active task, keep working. Send PM a brief progress update via agent-to-agent.
- If you finished a task, report completion to PM via agent-to-agent.
- If you are idle, tell PM you are available via agent-to-agent.
- Otherwise reply HEARTBEAT_OK.

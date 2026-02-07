# Heartbeat

## PM (main)
- Read `kanban.json`.
- For any `in_progress` task, spawn **Alex-Auditor** (`critic`) with expanded checks:
  - File activity in last 2h (git status/diff)
  - DIR_PROTOCOL compliance (absolute paths, frontend/backend segregation)
  - Build health: `npm run build` in affected directories
  - Test health: Playwright tests for assigned task passing
  - Repo stats: commits, files changed, lines added/removed
  - Stale flag if <5 lines changed in 2h
- Stale/error handling: If auditor reports no progress or failures, ping dev via `sessions_send`.
- Auto-assign: Unassigned high-priority backlog task to any dev with no `in_progress` task.
- **Team Sync Pulse:** Post a mini-standup to Telegram every heartbeat:
  - "Coder: TASK-042, X% done, ETA..."
  - "Foundry: TASK-045, status..."
  - "Blockers: none / [list]"
- **Kanban Health:**
  - Escalate tasks stuck >3 days in `in_progress`.
  - Before moving any task to `done`, verify:
    - Branch merged to main? (git log)
    - E2E tests exist? (check `/frontend/tests/e2e` for `[TASK-XXX]`)
  - Reject `done` if criteria unmet.
- **Dependency Radar:** If a task touches shared modules (e.g., `backend/handlers/league.go`, `frontend/src/components/Leagues.vue`), flag potential coupling and notify related task assignees.
- **Velocity Tracking:** Log actual vs estimated time per task; adjust future sizing based on historical multipliers (e.g., "high" takes 2× longer).
- **Gateway Health:** Check system resources (CPU, memory, disk); alert if thresholds exceeded; verify cron jobs running on schedule.
- **Configuration Validation:** Ensure critical files exist and are valid JSON (`.openclaw/kanban.json`, `config.json`). Auto-repair or alert on malformation.
- Reply `HEARTBEAT_OK` only if all health gates pass.

## Coder / Foundry
- Read `kanban.json`. Check your `assignee` tasks in `in_progress`.
- If assigned, resume work immediately and prepare detailed status:
  1. Current file being edited
  2. Next steps
  3. ETA (date/time)
- **Build & Test Gates:** Before reporting `done`, ensure:
  - `npm run build` passes in your domain (frontend/backend)
  - All Playwright E2E tests for the task are written and green
  - Tests include `data-testid` coverage where needed
- If stuck, send PM blocker via `sessions_send` (timeoutSeconds=0).
- On `HEARTBEAT_OK`, include brief progress note in your agent turn (auto-logged).
- Otherwise, reply `HEARTBEAT_OK`.

## Alex-Auditor (critic)
- Spawned by PM during heartbeat or on task completion.
- Audit scope:
  1. **Staling:** Check file modifications in last 2h (`git log -1`, `git diff`); if none, flag.
  2. **DIR_PROTOCOL:** Verify all file paths absolute and correctly segregated (frontend vs backend).
  3. **Build Health:** Run `npm run build` in relevant directories; report errors.
  4. **Validation:** Ask "Are you really finished?" and verify claims via `ls`, `git diff`, `cat`.
  5. **Repo Stats:** Count commits, files changed, lines added/removed; compare to threshold.
- Report findings immediately to PM with task id, assignee, and severity.
- Be skeptical — assume hallucination until proven.

## Senior Game Tester (tester)
- Maintain `/root/EmpoweredPixels/docs/TEST_CASES.md` (use cases) and `/root/EmpoweredPixels/docs/TEST_FINDINGS.md`.
- On every `done` task, run full test suite:
  - Playwright E2E for all user-facing pages
  - Manual browser checks for edge cases
- Report only **critical failures** to PM; include screenshots and repro steps.

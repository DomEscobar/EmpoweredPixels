# EmpoweredPixels — Technical Decisions

This log captures key architectural and process decisions, the rationale behind them, and any relevant context. It serves as a single source of truth for why we work a certain way.

---

## 2025-02-07: Agent Orchestration Framework

**Decision:** Adopt OpenClaw with custom agent roles (PM, Coder, Foundry, Alex-Auditor, Tester) and a strict heartbeat protocol.

**Rationale:**
- Automate project management and quality control.
- Ensure consistent quality via automated audits.
- Free human operators to focus on vision and creative direction.

**Consequences:**
- All tasks go through kanban (`.openclaw/kanban.json`).
- Agents follow `AGENTS.md` and `HEARTBEAT.md`.
- Strict DIR_PROTOCOL enforces clean separation of code.

---

## 2025-02-07: Testing Strategy

**Decision:** Playwright E2E tests are mandatory for every feature. Frontend components must include `data-testid` attributes.

**Rationale:**
- Game stability is critical; regressions must be caught early.
- `data-testid` provides reliable selectors for automation.
- E2E tests document real user workflows.

**Consequences:**
- No merge without passing Playwright suite.
- `data-testid` added to all interactive Vue components.
- Test failures block deployment.

---

## 2025-02-07: Repository Hygiene

**Decision:** Keep `.openclaw/kanban-ui/` and its `node_modules` untracked (local-only) while versioning only core configs (`kanban.json`, `AGENTS.md`, etc.).

**Rationale:**
- Visual kanban is a developer convenience, not part of the game codebase.
- Avoiding `node_modules` in git reduces repo bloat and merge conflicts.

**Consequences:**
- `.gitignore` excludes `.openclaw/kanban-ui/` and `.openclaw/node_modules/`.
- Kanban UI runs locally from `.openclaw/kanban-ui/`.
- Core kanban data remains tracked and shared via git.

---

## 2025-02-07: Agent Identity & Configuration Layout

**Decision:** Agent identity files (`SOUL.md`, `IDENTITY.md`, `USER.md`, `TOOLS.md`, `HEARTBEAT.md`, `PM_PROTOCOL.md`) live inside `.openclaw/`. `AGENTS.md` stays at workspace root.

**Rationale:**
- OpenClaw expects agent definitions at workspace root (`AGENTS.md`).
- `.openclaw/` is the designated directory for agent configuration, state, and metadata.
- Separation allows OpenClaw to load agent roles without scanning subdirectories.

**Consequences:**
- All persona and protocol files are in `.openclaw/`.
- Game code never resides in `.openclaw/`.
- Agents read their configs from `.openclaw/` on startup.

---

## 2025-02-07: Quality Gates via Pre-Merge Checks

**Decision:** Introduce `scripts/pre-merge-check.sh` that validates: branch up-to-date, builds pass, E2E tests present, `data-testid` coverage.

**Rationale:**
- Catch integration issues before they reach main.
- Automate the Definition of Done (DoD) enforcement.
- Provide immediate feedback to developers.

**Consequences:**
- Developers should run the script before marking a task `done`.
- PM can also run it during merge review.
- Builds and tests become non-negotiable gates.

---

## 2025-02-07: Metrics & Continuous Improvement

**Decision:** Track cycle time, escaped bugs, build success rate, and velocity. Use data to adjust future task sizing and heartbeat thresholds.

**Rationale:**
- Without measurement, we cannot improve.
- Historical velocity helps set realistic ETAs.
- Metrics expose bottlenecks (e.g., long cycle times, flaky builds).

**Consequences:**
- Alex-Auditor will verify metrics weekly.
- Heartbeat may eventually auto-tune thresholds based on trends.
- Transparency: metrics shared in Team Sync Pulse.

---

## 2025-02-07: Dependency Radar & Cross-Task Awareness

**Decision:** When a task touches shared modules (e.g., `backend/handlers/league.go`, `frontend/src/components/Leagues.vue`), the PM must notify assignees of related tasks.

**Rationale:**
- Changes in one area can break seemingly unrelated tasks.
- Early warning prevents wasted effort and merge conflicts.

**Consequences:**
- PM cross-references kanban labels (e.g., "leagues") when reviewing changes.
- Related task assignees receive `sessions_send` alerts.
- Encourages communication and coordination.

---

*This document is living. Update whenever a significant decision is made.*

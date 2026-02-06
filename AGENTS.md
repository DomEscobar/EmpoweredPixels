# EmpoweredPixels Agency - Sprint Protocol

The agency has transitioned to a **JSON-driven Sprint model**.

## ⚙️ Orchestration Core
- **Source of Truth:** `/root/EmpoweredPixels/tools/kanban-ui/kanban.json`
- **Visual Dashboard:** [http://v2202502215330313077.supersrv.de:8666](http://v2202502215330313077.supersrv.de:8666)
- **Status Lifecycle:** `todo` → `in_progress` → `done`

## 🤖 Active Agents
- **CEO (Main):** Orchestrates the high-level strategy and delegates.
- **Coder (@coder):** Builds game logic, combat simulators, and UI features.
- **Foundry (@foundry):** Manages DB architecture, infrastructure, and core refactoring.

## 📜 Execution Rules
1. **P0 First:** Tasks are prioritized by `priority` field in JSON.
2. **Branching:** Every `in_progress` task must have a corresponding git branch named after the task `id`.
3. **No Ghost Work:** If it's not on the JSON Board, it doesn't exist.
4. **Push Policy:** Immediate `git push` after every `safe_commit`.

---
*Updated: 2026-02-06 - Agency v1.1 Activation*

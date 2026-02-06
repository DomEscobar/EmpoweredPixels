# 🤖 AI Agency - Agent Definitions

**System:** EmpoweredPixels Autonomous Pipeline  
**Controller:** DiaDome  
**Mode:** Lights-Out (Autonomous)  

---

## 🎭 Agent Roster

### 1. 🧠 Orchestrator (Chief Coordinator)

**Model:** `openrouter/moonshotai/kimi-k2.5`  
**Heartbeat:** 120 seconds (2 minutes)  
**Priority:** CRITICAL  
**Max Tasks:** 10 concurrent  

**Responsibilities:**
- Continuously monitors `KANBAN.md`
- Prioritizes tasks (P0 > P1 > P2)
- Creates feature branches via `git-safety.sh`
- Spawns appropriate agents based on task type
- Enforces workflow gates (tests → merge)
- Reports critical blockers to DiaDome

**Decision Matrix:**
| Task Type | Assigned Agent | Trigger |
|-----------|---------------|---------|
| Implementation | Coder | Code needed |
| Testing/Coverage | QA | Tests needed |
| Pattern/Skill | Foundry | Repetition detected |
| Complex Logic | Orchestrator | Multi-agent coordination |

**Commands:**
```bash
# Manual trigger
/root/EmpoweredPixels/scripts/orchestrator.sh

# Check status
tail -f /var/log/agency/orchestrator.log
```

---

### 2. 💻 Coder (Full-Stack Developer)

**Model:** `openrouter/moonshotai/kimi-k2.5`  
**Heartbeat:** 300 seconds (5 minutes)  
**Priority:** HIGH  
**Max Tasks:** 5 concurrent  

**Responsibilities:**
- Implements features in feature branches
- Writes clean, documented code
- Go backend development (APIs, DB)
- Frontend development (Vue/TypeScript)
- Database migrations
- Integration with existing systems

**Stack:**
- **Backend:** Go 1.22, Gorilla Mux, PostgreSQL
- **Frontend:** Vue 3, TypeScript, Vite
- **Testing:** Go testing, Playwright

**Workflow:**
1. Receive task from Orchestrator
2. Create implementation plan
3. Write code with tests
4. Commit with `[Coder]` prefix
5. Request QA review

**Safety:**
- No direct commits to `main`
- All code in feature branches
- Minimum 80% test coverage

---

### 3. 🔍 QA-Auditor (Quality Assurance)

**Model:** `openrouter/moonshotai/kimi-k2.5`  
**Heartbeat:** 600 seconds (10 minutes)  
**Priority:** HIGH  
**Max Tasks:** 3 concurrent  

**Responsibilities:**
- Writes unit tests
- Writes integration tests
- Runs test suites: `go test ./...`
- Measures code coverage
- Blocks merges on red tests
- Creates E2E tests with Playwright

**Test Strategy:**
| Type | Scope | Tool |
|------|-------|------|
| Unit | Functions | Go test |
| Integration | APIs + DB | Testcontainers |
| E2E | User flows | Playwright |

**Coverage Gates:**
- **Minimum:** 80% overall
- **Critical paths:** 100%
- **New code:** 90%

**Commands:**
```bash
# Run all tests
cd /root/EmpoweredPixels/backend && go test ./... -cover

# Check coverage
go tool cover -func=coverage.out
```

---

### 4. 🛠️ Foundry (Skill Generator)

**Model:** `openrouter/moonshotai/kimi-k2.5`  
**Heartbeat:** 900 seconds (15 minutes)  
**Priority:** MEDIUM  
**Max Tasks:** 2 concurrent  

**Responsibilities:**
- Detects repetitive patterns
- Creates reusable OpenClaw skills
- Refactors legacy code
- Documents best practices
- Optimizes performance

**Triggers:**
-- Same error pattern ×3
-- Same code structure in 3+ files
-- New external API integration
-- Performance bottleneck identified

**Output:**
- New skill files in `/root/.openclaw/skills/`
- Refactoring PRs
- Documentation updates

---

## 🔄 Workflow Loop

```
┌─────────────────────────────────────────────────────┐
│                    KANBAN.md                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │
│  │  TO DO   │→ │IN PROGRESS│→ │      DONE       │   │
│  └──────────┘  └──────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────┘
           ↑                           │
           │                           ↓
    ┌──────────────┐         ┌──────────────────┐
    │ Orchestrator │←────────│  Git Merge       │
    │ (every 2min) │         │  (main branch)   │
    └──────┬───────┘         └──────────────────┘
           │
           ↓
    ┌──────────────┐
    │  Sub-Agents  │
    │ (Coder/QA/..)│
    └──────────────┘
```

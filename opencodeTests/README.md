# OpenCode Autonomous Agency — Swarm Architecture

## Overview

Multi-agent autonomous coding agency built on OpenCode CLI. Agents collaborate to plan, execute, review, and deploy code changes without human intervention.

## Core Agents

### director (Orchestrator)
- Receives tasks via CLI, webhook, or queue
- Plans task decomposition into subtasks
- Assigns subtasks to appropriate specialist agents
- Monitors progress, rebalances workload
- Escalates blockers
- Reports final status

### coder (Implementation)
- Primary executor using OpenCode
- Runs `opencode` to implement features/fixes
- Writes code, tests, documentation
- Handles multiple programming languages

### reviewer (Quality Gate)
- Reviews PRs/changes from coder
- Runs static analysis, security checks
- Validates tests pass
- Approves or requests revisions

### tester (QA Agent)
- Generates test cases from requirements
- Runs automated tests (unit, integration)
- Performs exploratory testing sessions
- Reports bugs back to director

### deployer (Release)
- Manages CI/CD pipelines
- Executes deployments on approval
- Monitors post-deploy health
- Rollback on failure

## Orchestration

- **Task Queue**: SQLite-backed queue for pending tasks
- **Communication**: Direct message passing via filesystem/DB
- **Supervision**: Director spawns/manages worker lifecycles
- **State**: `/var/lib/opencode-agency/` (configurable)

## Prerequisites

- OpenCode CLI installed and authenticated
- Git repository initialized
- Project-specific AGENTS.md for context

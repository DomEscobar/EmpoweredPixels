# EmpoweredPixels Architecture

This document describes the technical architecture and standards for the EmpoweredPixels project. It is maintained by **Coder** and **Foundry**.

## Tech Stack
- **Frontend**: Vue 3 (Composition API, Script Setup), Vite, Tailwind CSS v4, Pinia (Store), Playwright (E2E Testing).
- **Backend**: Node.js/TypeScript (Services), REST API.
- **Infrastructure**: Docker Compose, OpenClaw Gateway.

## Directory Structure
- `/frontend`: All Vue source code, assets, and UI logic.
- `/backend`: Server-side logic and database interactions.
- `/docs`: Documentation and architecture notes.
- `.openclaw`: Agent configurations and PM state (Metadata only).

## UI/UX Theme: "Ethereal Iron"
A darker, high-readability fusion of WoW structure, GW2 artistic splatters, and D4 moody aesthetics.
- Palette: Obsidian (#0d0d12), Iron Steel (#2a2a35), Gold Accent (#e2b349).
- Asset ID Strategy: Hardcoded UUID-like strings with standard pixel art prompt formatting.

## Feature Patterns

### Skill Mastery System
The Mastery Constellation allows for branched character progression.
- **Node Logic**: Sequential unlocking within branches. Unlocking requires "Soul-Shards" (available as part of the fighter's progression state).
- **Communication**: Components emit `unlock` and `reset` events to notify parent state/API handlers.
- **Visuals**: Uses CSS animations (`flow`) and SVG/CSS connective lines to denote path status. Connective lines are activated only when the target node is reachable or unlocked.

## Testing Standards
- **Mandatory E2E**: Every new feature or significant UI overhaul MUST include a Playwright test file in `/frontend/tests/e2e`.
- **Criteria**: Acceptance criteria in `kanban.json` are only met when these tests pass in the deployment environment.

# Empowered Pixels - Current State

Empowered Pixels is a game project featuring a Go backend and a Vue 3 frontend. It uses a clean architecture approach for both components.

## 🏗️ Architecture Overview

### Backend (Go)
- **Architecture**: Clean Architecture (Domain, Usecase, Adapter, Infra).
- **Core Features**:
    - **Identity/Auth**: Registration, login, and password management.
    - **Roster/Fighters**: Fighter management, stats, and XP.
    - **Combat/Matches**: Match history, simulator, and real-time (WS) hub.
    - **Progression**: Seasons, leagues, leaderboard, and rewards (daily/weekend).
    - **Economy**: Shop with gold packages, bundles, and purchase flows.
    - **Inventory**: Weapon and equipment management.
    - **Customization**: Attunements and skills.
- **Tech Stack**:
    - Go (Golang)
    - PostgreSQL (DB)
    - WebSocket (Match notifications)

### Frontend (Vue 3)
- **Tech Stack**: Vue 3, Pinia, TypeScript, Tailwind CSS, Vite.
- **Pages/Views**:
    - `Register.vue` / `Login.vue`: User onboarding.
    - `Dashboard.vue` / `Home.vue`: Main overview.
    - `Roster.vue`: Fighter management.
    - `Inventory.vue`: Equipment and weapons.
    - `Matches.vue` / `MatchViewer.vue`: Matchmaking and battle replay.
    - `Leagues.vue` / `Leaderboard.vue`: Competitive view.
    - `Shop.vue`: Economy and purchases.
    - `Attunement.vue`: Fighter customization.
- **UI & Shared**:
    - Voxel-based fighter generation (`voxelGenerator.ts`).
    - Clean separation between features (api, store, components) and pages.

## 📁 Repository Structure

- `backend/`: Go source code, Dockerfile, migrations.
- `frontend/`: Vue application, Vite config, E2E tests (Playwright).
- `assets/`: Game assets (weapons, etc.).
- `docs/`: Additional documentation.

## 🛠️ Tech Stack & Dev Tools

- **Backend**: Go 1.22+, PostgreSQL.
- **Frontend**: Node 20+, Vue 3 (Composition API), Pinia (State), Tailwind CSS (UI), Vite (Build).
- **Automation/Dev**: `run-dev.ps1` (launcher), `docker-compose.yml`.

## 📈 Recent Progress (from file structure)
- Implemented **Shop** system with bundles and gold packages.
- Added **Daily Rewards** and **Weekend Events**.
- Built **Leagues** and **Seasons** logic.
- Implemented **Fighter XP** and **Leaderboards**.
- Integrated **Attunement** system.
- E2E testing with Playwright for matchmaking and battling.

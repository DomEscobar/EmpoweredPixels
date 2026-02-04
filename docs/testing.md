# Testing Log

This file tracks the verification of features via browser testing (Playwright MCP) and API smoke tests.

## Feature Status

| Feature | Status | Last Verified | Notes |
| :--- | :--- | :--- | :--- |
| Landing Page | ✅ Passed | 2026-02-03 | Beautiful UI with gradients and text rendering. |
| Registration | ✅ Passed | 2026-02-03 | Successful user creation and redirect to login. |
| Login | ✅ Passed | 2026-02-03 | Token acquisition and redirect to dashboard. |
| Dashboard | ✅ Passed | 2026-02-03 | Verified stats and dynamic counts. |
| Roster List | ✅ Passed | 2026-02-03 | Verified fighters display with equipment/stats. |
| Fighter Creation | ✅ Passed | 2026-02-03 | Verified modal and multi-fighter support. |
| Match Browser | ✅ Passed | 2026-02-03 | Lobby listing, status filter, create/join/leave. |
| Fighter selection on join | ✅ Passed | 2026-02-03 | Join modal with fighter picker. |
| Match status (lobby/running/completed) | ✅ Passed | 2026-02-03 | Status-driven UI and start/leave/replay. |
| Combat Engine | ✅ Passed | 2026-02-03 | Simulator runs 50-round battles with bots. |
| Match Replay | ✅ Passed | 2026-02-03 | Battle log and canvas with round logs. |
| Live banner (running match) | ✅ Passed | 2026-02-03 | Match viewer shows live when status=running. |
| Inventory | ✅ Passed | 2026-02-03 | Verify resource balances. |
| Leagues | ✅ Passed | 2026-02-03 | Verify league cards. |
| League run (POST /api/league/:id/run) | Manual | 2026-02-03 | Triggers league match for subscribers. |
| Logout | ✅ Passed | 2026-02-03 | Verify session clearing and redirect. |

## End-to-End Test Matrix

| Flow | Steps | Expected |
| :--- | :--- | :--- |
| Lobby lifecycle | Create lobby → Join (select fighter) → Start battle → View replay | Match status: lobby → running → completed; replay has rounds. |
| Leave lobby | Join lobby → Leave | Current match cleared; can browse again. |
| Start without enough fighters | Create lobby, 0 bots, do not join with second player → Start | 400 Not enough fighters (or add bots). |
| Browse by status | Browse with status=lobby, running, completed | List filtered by status. |
| Match viewer live | Open viewer for match with status=running | Banner "Live: battle in progress"; when completed, replay loads. |
| League run | Subscribe to league with fighter → POST /api/league/:id/run | Match created, all subscribers joined, battle runs, league_matches updated. |

## API Smoke Tests

- **GET /health**: 200.
- **POST /api/match/browse** (body: `{ "page": 1, "pageSize": 20, "status": "lobby" }`, auth): 200, `items` array.
- **GET /api/match/:id** (auth): 200, `status`, `started`, `completedAt` in response.
- **POST /api/match/:id/start** (auth): 200 when lobby has enough participants; 400 when not lobby or not enough fighters.

## WebSocket

- **GET /ws/match**: Connect, send `{ "action": "subscribe", "matchId": "<id>" }`; expect `matchStatus` / `matchEnded` / `lobbyUpdate` when match state changes.

## Playwright E2E (Browser)

- **Location**: `frontend/e2e/matchmaking-battling.spec.ts`
- **Run**: From `frontend`: `npx playwright test e2e/matchmaking-battling.spec.ts`
- **Prerequisites**: Backend API on **port 54321** (or set `VITE_API_BASE_URL` in frontend env). Frontend is started automatically by Playwright unless in CI.

| Test | Description | Result (no backend) | Result (backend :54321) |
| :--- | :--- | :--- | :--- |
| matches page shows status filter and create button | Login → Matches → assert heading, status select, Create Lobby button | ❌ (login fails) | ✅ Pass |
| register, login, create fighter, create lobby, start battle, view replay | Full flow: register → login → roster (create fighter) → matches (create lobby, start battle, view replay) | ❌ (fails at API calls) | ✅ Pass when backend + DB available |

**Notes**: The full-flow test uses unique credentials per run. It depends on backend for register, login, fighter create, match create/start, and replay. Run with backend on :54321 and DB for full pass.

## Test Execution Details

### 2026-02-03: Playwright E2E and Matchmaking
- **Goal**: Automate matchmaking and battling flows via browser tests; maintain testing.md.
- **Environment**: Frontend :5173 (Playwright webServer); backend must be :54321 for API (default in frontend `http.ts`).
- **Result**: Selectors fixed (landing "Get Started Free", login form-scoped "Sign In", register "Access HQ", roster "Create Fighter", create-lobby flow). Test "matches page shows status filter and create button" passes. Full-flow test passes when backend is running; otherwise fails at first API-dependent step (e.g. fighter create or lobby create).

### 2026-02-03: Full UI & Gameplay Verification
- **Goal**: Verify all core pages, auth state, and the new combat engine.
- **Environment**: Backend :54321, Frontend :5173.
- **Result**: All core flows passed. Successfully registered, logged in, created a fighter, initialized a match with bots, and reviewed the battle log in the Match Viewer.
- **Extensive Gameplay Test**: Created multiple matches with 1-3 bots and levels up to 30. Verified the Battle Log correctly renders complex turn-based interactions including movements, attacks, and deaths across 50 rounds.

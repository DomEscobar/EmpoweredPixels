# Matchmaking Flow

This document describes the end-to-end matchmaking and battle execution flow in EmpoweredPixels.

## 1. Discovery & Lobby Creation
- **UI:** `Matches.vue` header and browse controls. `showCreate` toggle for modal.
- **State:** `browseStatus` (lobby/running/completed), `matches` array, `options` object for creation.
- **Backend:** `POST /api/match/browse` for filtering by status; `PUT /api/match/create` for persistence.
- **Improvements:**
  - **Auto-Join:** Creating a lobby now automatically registers the creator's primary fighter in the backend.
  - **Status Banner:** A persistent, high-visibility banner at the top of the Matches page shows the user's current match status.

## 2. Joining a Match
- **UI:** "Join" button in match cards -> `showJoinModal` (standardized `BaseModal`). Fighter selection list with level and power context.
- **State:** `selectedFighterId`, `currentMatchId`, `currentMatch` (full object).
- **Backend:** `POST /api/match/join` (requires `matchId` and `fighterId`).
- **WebSocket:**
  - Client sends `{ "action": "subscribe", "matchId": "[ID]" }`.
  - Server broadcasts `lobbyUpdate`.
- **Improvements:**
  - **Current Match Endpoint:** Reconnection now uses a dedicated `GET /api/match/current` endpoint, replacing the previous inefficient $O(N)$ scanning logic.

## 3. Battle Initiation
- **UI:** "Start Battle" button in the persistent status banner.
- **State:** `isStarting` (loading state), `currentMatchStatus` (computed).
- **Backend:** `POST /api/match/[ID]/start`. Transitions DB status `lobby` -> `running`.
- **Broadcast:** Server sends `matchStatus: running`.
- **Improvements:**
  - **Standardized Modals:** All lobby interactions now use `BaseModal` for a consistent look and feel.

## 4. Execution & Live Monitoring
- **UI:** `MatchViewer.vue` "Live" banner. `Matches.vue` "Watch Live" button in the status banner.
- **State:** `matchStatus` (polled in viewer), `isLoading`.
- **Backend:** `ExecuteMatch` service runs combat logic.
- **Broadcast:** `matchEnded` when done.
- **Improvements:**
  - **Visual Feedback:** The "Watch Live" button is prominently displayed in the status banner as soon as a match starts.

## 5. Battlefield Visualization (Replay)
- **UI:** Canvas with enhanced visuals, Round selection buttons, Timeline slider, Interactive Round Logs.
- **State:** `selectedRound`, `isPlaying`, `playbackSpeed`, `roundStateMap`.
- **Backend:** `GET /api/match/[ID]/roundticks` returns the full JSON blob of events.
- **Improvements:**
  - **Health Bars:** Dynamic HP bars are rendered above each entity on the canvas.
  - **Entity Highlighting:** The user's own fighters have a distinct blue aura/ring for easy identification.
  - **Damage Numbers:** Floating red/gold numbers appear when damage is dealt.
  - **Interactive Logs:** Clicking any entry in the round logs sidebar automatically seeks the playback to that round.

## 6. Technical Notes & Observations
- **Synchronous Execution:** The `ExecuteMatch` call is currently synchronous. For future scalability, moving this to a background worker is recommended.
- **WebSocket Usage:** WebSockets handle lobby updates and status transitions. 
- **Optimized Reconnection:** The frontend efficiently restores match state on mount via `/api/match/current`.
- **Standardized UI:** The entire matchmaking and battle viewing experience now uses standardized components (`BaseModal`, `BaseButton`, `BaseCard`).

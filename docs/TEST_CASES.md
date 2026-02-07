# Empowered Pixels - Test Cases (Senior Game Tester Edition)

This document defines the core use cases and validation points for **Empowered Pixels**. It serves as the source of truth for QA automation and manual testing.

---

## 1. Authentication System
*Target: `/register`, `/login`, `AuthMiddleware`*

### 1.1 User Registration
- **Success Case**: Enter unique email, valid username, and secure password.
    - **Validation**: User created in `identity` table, salt generated, password hashed.
- **Duplicate User**: Attempt to register with existing email or username.
    - **Validation**: Return `409 Conflict`.
- **Validation**: Invalid email format or password too short.
    - **Validation**: Return `400 Bad Request`.

### 1.2 User Login
- **Success Case**: Correct credentials provided.
    - **Validation**: Return JWT and Refresh Token. `LastLogin` timestamp updated in DB.
- **Wrong Password**: Correct email, wrong password.
    - **Validation**: Return `401 Unauthorized`. Side channel attack protection (consistent response time).
- **Banned User**: Attempt login with a banned account.
    - **Validation**: Check `Banned` field in DB. Return `403 Forbidden` with reason.

---

## 2. Inventory Management
*Target: `InventoryGrid.vue`, `EquipmentCard.vue`, `inventory_postgres.go`*

### 2.1 Filtering & Sorting
- **Rarity Filter**: View only "Legendary" items.
    - **Validation**: Frontend list refreshes based on `Rarity` enum.
- **Type Filter**: Filter by Weapons vs. Consumables.
- **Sorting**: Sort by "Recent", "Level", or "Rarity".

### 2.2 Tooltips & Item Details
- **Hover/Click Interaction**: Trigger `EquipmentCard` detail view.
    - **Validation**: Accurate display of Stats (Power, Vitality, etc.), Enhancement Level, and "Favorite" status.
- **Fighter Assignment**: Check which item is equipped to which Fighter via `FighterID` foreign key.

---

## 3. Combat Core (Battle Simulation)
*Target: `BattleSimulator.go`, `MatchViewer.vue`, `match_hub.go`*

### 3.1 Battle Initialization
- **Spawn Logic**: Random spawn points within `MapSize`.
    - **Validation**: `spawn` events included in Round 0 ticks.
- **Initiative Order**: Turn order calculated via `Speed + Agility + RNG`.

### 3.2 Combat Resolution
- **Skill Execution**: Range check before execution.
    - **Validation**: Skill triggers `damage`, `heal`, or `buff` events.
- **Momentum System**: Combo count increases on consecutive hits; resets on movement.
- **Death & Scoring**: When `CurrentHP <= 0`, emit `died` event.
    - **Validation**: `Scores` table updated (Kills/Deaths/Damage).

### 3.3 Logs & Viewer
- **Ticking System**: Match data structure is a stream of JSON-encoded ticks.
- **Replay Accuracy**: Playback in `MatchViewer` must match the server-side simulation exactly.

---

## 4. Mastery Constellation (Skill Unlock Logic)
*Target: `MasteryConstellation.vue`, `CanAllocate` logic in `models.go`*

### 4.1 Skill Allocation
- **Branch Dependency**: Unlocking Tier 2 in "Offense" requires 2 points in Tier 1 of the same branch.
- **Rank Scaling**: Allocating points (1-3) must scale `EffectValue` by 25% per rank.
- **Level Gate**: One point granted per Level. Total points cannot exceed `FighterLevel`.

### 4.2 Loadout Configuration
- **Active Skills Limit**: Maximum 2 active skills in `Loadout`.
- **Validation**: Cannot set a "Passive" skill into an active loadout slot.

---

## 5. Squad System
*Target: `squad_postgres.go`, `squad_service.go`, `Roster.vue`*

### 5.1 Squad Composition (Current state: Logic/DB ready)
- **Active Squad**: User can have multiple squads, but only one `isActive`.
- **Slot Management**: Squad contains up to 5 members.
- **Member Validation**: Ensure selected Fighter belongs to the logged-in User.

### 5.2 Persistence
- **State Change**: Switching members triggers immediate DB update in `squad_members` table.
- **Integrity**: Deleting a Fighter should remove them from all Squads (Cascade/Manual cleanup).

---

🛡️🧪 *End of Document. Written by Senior Game Tester Agent.*

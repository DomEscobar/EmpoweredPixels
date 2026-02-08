# Elemental Resonance Squad Synergy System

## Overview

The **Elemental Resonance Squad Synergy System** is a game mechanic that calculates squad harmony based on fighter attunements and applies bonuses/debuffs to combat stats.

### Core Mechanic

- **Harmony Score**: 0-100 calculated from fighter attunement composition
- **Harmonic Pairs** (synergy): Fire↔Air, Water↔Earth, Light↔Dark (+8-12 points each)
- **Dissonant Pairs** (conflict): Fire↔Water, Air↔Earth (-5 penalty)
- **Bonuses/Debuffs**: Up to +12% damage/defense for high harmony, -5% for dissonance
- **Cosmetic Auras**: Visual effects in match viewer based on harmony tier
- **Seasonal Achievements**: Unlock "Resonance Master" for 80+ harmony + 100+ matches

---

## Implementation Details

### Score Formula

```
Base: 25 points
+ (Harmonic Pairs Count × 10)
- (Dissonant Pairs Count × 5)
= Final Score (clamped 0-100)
```

**Example**: Fire + Air + Water
- Harmonic: Fire↔Air (+10)
- Dissonant: Fire↔Water (-5)
- Score: 25 + 10 - 5 = **30** (Aligned tier)

### Tiers & Bonuses

| Tier | Score Range | Damage Bonus | Defense Bonus | Aura Color |
|------|---|---|---|---|
| Resonant | 76-100 | 1.12 (+12%) | 1.06 (+6%) | #FFD700 (Gold) |
| Harmonized | 51-75 | 1.08 (+8%) | 1.04 (+4%) | #FFD700 (Gold) |
| Aligned | 26-50 | 1.0 (+0%) | 1.0 (+0%) | #4169E1 (Blue) |
| Discordant | 0-25 | 0.95 (-5%) | 0.95 (-5%) | #808080 (Gray) |

---

## Code Architecture

### Backend

#### Domain Layer (`internal/domain/attunement/`)
- **`resonance.go`**: Core harmony calculation
  - `CalculateHarmony(elements []Element) (score int, pattern *ResonancePattern)`
  - `GetBonusMultiplier(score int) (dmg, def float64)`
  - `GetTierName(score int) string`
  - `GetDissonanceWarning(pattern *ResonancePattern) string`

#### Usecase Layer (`internal/usecase/resonance/`)
- **`service.go`**: Business logic
  - `CalculateSquadResonance(squadID) (*ResonanceState, error)`
  - `ApplySquadBonuses(fighters []Fighter, state ResonanceState) []Fighter`
  - WebSocket/HTTP integration

#### Adapter Layer
- **HTTP Handler** (`internal/adapter/http/handlers/resonance/`)
  - `GET /api/v1/squads/:squadID/resonance` → Returns `ResonanceResponse`
  - `POST /api/v1/squads/:squadID/resonance/prefetch` → Triggers calculation

- **WebSocket** (`internal/adapter/ws/`)
  - Broadcasts `match.resonance_state` message after match starts
  - Payload includes harmonies for each squad

#### Database
- **Migrations**:
  - `0019_add_squad_resonance.up.sql`: Adds `resonance_score`, `resonance_pattern`, `last_calculated_at` to `squads` table
  - `0020_add_resonance_achievements.up.sql`: Creates `resonance_achievements` table for tracking
  
- **Repositories**:
  - `ResonanceAchievementRepository`: Tracks achievement progress and unlock status

#### Match Integration
- **Match Service** (`internal/usecase/matches/service.go`)
  - Calls `CalculateSquadResonance()` for each user's squad
  - Applies bonuses to fighters before simulator runs
  - Broadcasts resonance state via WebSocket
  - Updates achievements on match completion

### Frontend

#### Components
- **`ResonancePreview.vue`**: Squad harmony display
  - Attunement circles with element icons
  - Harmony/dissonance connection lines
  - Score gauge with tier indicator
  - Bonus table and aura preview

- **`AuraOverlay.vue`**: Visual aura effect
  - Drop-shadow + glow filters
  - Pulse animations based on tier
  - Tier badge and bonus indicators

#### Composables
- **`useMatchResonance.ts`**: WebSocket integration
  - `onResonanceMessage()`: Handle incoming resonance state
  - `getAuraStyle()`: Compute CSS filters for aura effects
  - `getAuraIntensity()`: Calculate animation intensity per tier

#### E2E Tests
- **`e2e/resonance.spec.ts`**: End-to-end test scenarios
  - Squad creation with known attunements
  - API endpoint validation
  - Match viewer aura rendering
  - WebSocket message capture
  - Tier display and colors

---

## API Reference

### GET /api/v1/squads/{squadID}/resonance

**Response**:
```json
{
  "squadID": "550e8400-e29b-41d4-a716-446655440000",
  "harmonyScore": 75,
  "tierName": "Harmonized",
  "harmonicElements": ["Fire", "Air"],
  "dissonantElements": ["Water"],
  "bonuses": {
    "damage": 1.08,
    "defense": 1.04
  },
  "auraColor": "#FFD700"
}
```

### WebSocket: match.resonance_state

**Message Type**: `match.resonance_state`

**Payload**:
```json
{
  "type": "match.resonance_state",
  "resonances": {
    "fighter-id-1": {
      "squadID": "...",
      "harmonyScore": 80,
      "tierName": "Resonant",
      "harmonicElements": ["Fire", "Air", "Light"],
      "dissonantElements": [],
      "bonuses": {"damage": 1.12, "defense": 1.06},
      "auraColor": "#FFD700"
    },
    "fighter-id-2": { ... }
  }
}
```

---

## Usage Examples

### For Game Designers

1. **Create Harmonic Squad**:
   - Select fighters with Fire, Air, Light attunements
   - Expected harmony score: ~35-45 (Aligned tier)
   - Fighters gain +0-8% damage/defense

2. **Create Dissonant Squad**:
   - Select fighters with Fire, Water, Earth attunements
   - Expected harmony score: ~20 (Discordant tier)
   - Fighters suffer -5% damage/defense penalty

### For Frontend Developers

1. **Display Resonance Preview**:
   ```vue
   <ResonancePreview :resonanceState="squad.resonance" />
   ```

2. **Listen to WebSocket Resonance**:
   ```typescript
   const { onResonanceMessage, getAuraStyle } = useMatchResonance()
   
   websocket.addEventListener('message', (event) => {
     onResonanceMessage(JSON.parse(event.data))
   })
   
   // Apply aura to fighter element
   const auraStyle = getAuraStyle(fighterId)
   ```

### For Backend Developers

1. **Calculate Resonance**:
   ```go
   elements := []attunement.Element{attunement.Fire, attunement.Air}
   score, pattern := attunement.CalculateHarmony(elements)
   // score = 35, pattern.TierName = "Aligned"
   ```

2. **Apply Bonuses**:
   ```go
   dmg, def := attunement.GetBonusMultiplier(score)
   fighter.Power = int(float64(fighter.Power) * dmg)
   fighter.Armor = int(float64(fighter.Armor) * def)
   ```

---

## Testing

### Unit Tests

**Domain Tests** (`internal/domain/attunement/resonance_test.go`):
- Test harmony calculation for various combos
- Verify tier names and aura colors
- Check bonus multipliers per score range

**Run**: `go test ./internal/domain/attunement/...`

**Usecase Tests** (`internal/usecase/resonance/service_test.go`):
- Mock repositories to test service logic
- Verify bonus application to fighters
- Test resonance state structure

**Run**: `go test ./internal/usecase/resonance/...`

### E2E Tests

**File**: `frontend/e2e/resonance.spec.ts`

**Scenarios**:
1. Squad creation with harmonic elements
2. API endpoint response validation
3. Match viewer aura rendering
4. WebSocket message capture
5. Dissonance warning display
6. Tier color verification

**Run**: `npm run test:e2e -- resonance.spec.ts`

---

## Deployment Checklist

- [x] Domain logic implemented and tested
- [x] Database migrations created (up/down)
- [x] API endpoint implemented and routed
- [x] WebSocket broadcast added
- [x] Match simulator integration
- [x] Achievement tracking added
- [x] Frontend components created
- [x] E2E tests written
- [ ] Staging deployment verified
- [ ] Performance tested (large squads)
- [ ] Rollback procedure documented

See `RESONANCE_SMOKE_TEST.md` for detailed smoke test procedures.

---

## Future Enhancements

1. **Squad Synergy Combos**: Special bonuses for specific 3-element combos (e.g., "Triangle of Power" for Fire+Water+Earth)
2. **Dynamic Harmony**: Adjust harmony mid-match based on combat actions
3. **Harmony Artifacts**: Equippable items that boost or adjust harmony
4. **Leaderboard**: "Harmony Masters" ranked by highest harmony matches won
5. **Seasonal Resonance**: Different element bonuses per season
6. **Resonance Spells**: Ultimate abilities unlocked at high harmony tiers

---

## Files Modified

### Backend
- `cmd/api/main.go` — Added resonance service initialization
- `internal/domain/attunement/resonance.go` — Core logic
- `internal/domain/attunement/resonance_test.go` — Unit tests
- `internal/usecase/resonance/service.go` — Business logic
- `internal/usecase/resonance/service_test.go` — Service tests
- `internal/usecase/matches/service.go` — Match integration
- `internal/adapter/http/handlers/resonance/resonance_handler.go` — HTTP endpoint
- `internal/adapter/http/router.go` — Route registration
- `internal/infra/db/migrations/0019_add_squad_resonance.up.sql` — Schema
- `internal/infra/db/migrations/0020_add_resonance_achievements.up.sql` — Achievements
- `internal/infra/db/repositories/resonance_achievements_postgres.go` — Repository

### Frontend
- `src/features/squads/components/ResonancePreview.vue` — Squad preview component
- `src/components/AuraOverlay.vue` — Aura effect component
- `src/composables/useMatchResonance.ts` — WebSocket integration
- `e2e/resonance.spec.ts` — E2E tests

### Documentation
- `RESONANCE_FEATURE_README.md` — This file
- `RESONANCE_SMOKE_TEST.md` — Deployment checklist

---

## Support

For questions or issues:
1. Check the `RESONANCE_SMOKE_TEST.md` for troubleshooting
2. Review test cases in `resonance_test.go` for expected behavior
3. Inspect WebSocket messages in browser DevTools (Network → WS tab)

# Elemental Resonance Squad Synergy - Smoke Test Checklist

## Pre-Deployment Verification

### Backend Build & Tests
- [ ] Build backend locally: `cd backend && go build ./cmd/...`
  - Expected: No compilation errors
- [ ] Run domain tests: `go test ./internal/domain/attunement/...`
  - Expected: All tests pass, including resonance calculation tests
- [ ] Run usecase tests: `go test ./internal/usecase/resonance/...`
  - Expected: ResonanceService tests pass

### Database Migration
- [ ] Review migration files:
  - `backend/internal/infra/db/migrations/0019_add_squad_resonance.up.sql`
  - `backend/internal/infra/db/migrations/0020_add_resonance_achievements.up.sql`
- [ ] Test migrations locally:
  ```bash
  psql -U postgres -d empoweredpixels -f backend/internal/infra/db/migrations/0019_add_squad_resonance.up.sql
  psql -U postgres -d empoweredpixels -f backend/internal/infra/db/migrations/0020_add_resonance_achievements.up.sql
  ```
  - Expected: No errors, tables created with correct columns and indexes

### Code Review Checklist
- [ ] Domain model (`resonance.go`): Harmonic/dissonant pair logic correct
- [ ] Score formula: `(harmonic_count × 10) + (dissonant_count × -5) + base 25`
- [ ] Tier names match score ranges: Discordant (0-25), Aligned (26-50), Harmonized (51-75), Resonant (76-100)
- [ ] Bonus multipliers: 51-75 = 1.08/1.04, 76-100 = 1.12/1.06
- [ ] ResonanceService properly imports and uses repositories
- [ ] API handler validates inputs and returns expected JSON structure
- [ ] Match simulator integration applies bonuses before battle runs
- [ ] WebSocket broadcasts resonance state with correct message format
- [ ] Achievement repository tracks matches and harmony scores

---

## Staging Deployment

### Deploy Steps
1. Run migrations on staging database:
   ```bash
   go run ./cmd/migrations/main.go --env staging
   ```
   
2. Deploy backend to staging
3. Deploy frontend to staging

### API Endpoint Smoke Tests

**Test 1: Load squad and call resonance endpoint**
```bash
# Get a test squad ID (adjust as needed)
SQUAD_ID="<test-squad-uuid>"

# Call resonance endpoint
curl -X GET "http://staging-api.example.com/api/v1/squads/$SQUAD_ID/resonance" \
  -H "Authorization: Bearer <token>"

# Expected response:
{
  "squadID": "...",
  "harmonyScore": 30-75,  // varies by squad
  "tierName": "Aligned|Harmonized|Resonant|Discordant",
  "harmonicElements": [...],
  "dissonantElements": [...],
  "bonuses": {
    "damage": 0.95-1.12,
    "defense": 0.95-1.06
  },
  "auraColor": "#808080|#4169E1|#FFD700"
}
```

**Test 2: Verify score calculation is correct**
- Create or find squad with known attunements
- Call endpoint
- Manually verify score formula:
  - Fire + Air = +10 (harmonic)
  - Fire + Water = -5 (dissonant)
  - Expected score: 25 + 10 - 5 = 30 ✓

**Test 3: Start practice match and observe bonuses applied**
- Create match with test squad (Fire + Air + Water attunements)
- Observe fighters' initial stats (e.g., Power: 100, Armor: 50)
- Expected: Stats multiplied by bonus (e.g., Power: 108, Armor: 52 if score is 75)
- Run simulation and verify damage/defense values in combat reflect bonuses

**Test 4: Check match viewer aura renders**
- Start match and open in browser
- Open DevTools Console → Network → WS
- Subscribe to match WebSocket
- Expected message received:
  ```json
  {
    "type": "match.resonance_state",
    "resonances": {
      "fighter-id-1": {
        "harmonyScore": 75,
        "tierName": "Harmonized",
        "auraColor": "#FFD700"
      }
    }
  }
  ```
- Visual check: Fighters on canvas should have glow/aura effect
- Verify aura color matches `auraColor` in message

**Test 5: WebSocket message logged in browser console**
```javascript
// In browser console:
// Watch for messages in Network tab, or add debug handler:
ws.addEventListener('message', (event) => {
  const data = JSON.parse(event.data);
  if (data.type === 'match.resonance_state') {
    console.log('✅ Resonance state received:', data);
  }
});
```
- Expected: At least one message with `match.resonance_state` type

---

## Rollback Procedure

If critical issues found during smoke tests:

1. **Revert migrations** (if database issues):
   ```bash
   psql -U postgres -d empoweredpixels -f backend/internal/infra/db/migrations/0020_add_resonance_achievements.down.sql
   psql -U postgres -d empoweredpixels -f backend/internal/infra/db/migrations/0019_add_squad_resonance.down.sql
   ```

2. **Redeploy previous backend version** (if code issues)

3. **Verify squads table still queries correctly** after rollback:
   ```bash
   psql -U postgres -d empoweredpixels -c "SELECT * FROM squads LIMIT 1;"
   ```
   - Expected: Query succeeds, no errors about missing columns

---

## Sign-Off Checklist

- [ ] All tests pass locally
- [ ] Migrations apply without errors
- [ ] API endpoint returns correct response structure
- [ ] Harmony scores calculated correctly for test squads
- [ ] Match simulator applies stat bonuses before combat
- [ ] WebSocket broadcasts resonance state message
- [ ] Match viewer renders aura effects correctly
- [ ] Aura colors match tier (gray/blue/gold/vibrant)
- [ ] Achievement table created and tracked
- [ ] No nil pointer panics in logs
- [ ] No schema validation errors
- [ ] Rollback procedure tested and verified

---

## Notes

- Resonance scores are cached in `squads` table: `resonance_score`, `resonance_pattern`, `last_calculated_at`
- WebSocket broadcast is sent during `ExecuteMatch()` before simulator runs
- Bonuses applied to fighters: `Power`, `ConditionPower` (damage) and `Armor` (defense)
- Achievement unlocks when `harmony >= 80 AND matches >= 100`
- Frontend components: `ResonancePreview.vue`, `AuraOverlay.vue`, `useMatchResonance.ts`

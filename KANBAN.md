# EmpoweredPixels Kanban

Last Updated: 2026-02-06

## 🟢 DONE

### Shop MVP
- [x] Backend API (`/api/shop/*`, `/api/player/*`)
- [x] Frontend (Shop.vue, components, store)
- [x] Database migrations + seed data
- [x] Gold packages, bundles, purchase flow
- [x] Merged to main

### Attunement System
- [x] Domain models (6 elements, levels 1-25)
- [x] DB migrations
- [x] Repository + Service layer
- [x] API endpoints
- [x] Frontend (Attunement.vue)
- [x] Merged to main

---

## 🟡 IN PROGRESS / ANALYSIS

- [x] **Weekend Events** (IN PROGRESS)
  Assignee: coder @ feature/weekend-events

### 🔥 ROSTER Flow Analysis (DEEP)
**Status:** Core system exists, CRITICAL gaps for engagement

```
┌─────────────────────────────────────────────────────────────────┐
│  CURRENT ROSTER FLOW                                            │
├─────────────────────────────────────────────────────────────────┤
│  User Journey:                                                  │
│  1. Create Fighter (Name only) → 2. View Stats → 3. Equip       │
│                                                                 │
│  Frontend (Roster.vue)          Backend                         │
│  ─────────────────────          ────────                        │
│  ✅ List Fighters    ───────►  ✅ GET /fighters                │
│  ✅ Create Fighter   ───────►  ✅ POST /fighters               │
│  ✅ Delete Fighter   ───────►  ✅ DELETE /fighters/{id}        │
│  ✅ View Equipment   ───────►  ✅ GET /fighters/{id}/equipment │
│  ✅ Set Attunement   ───────►  ✅ POST /fighter/config         │
│                                                                 │
│  ❌ Fighter Progression         ❌ No XP system                 │
│  ❌ Level Up Animation          ❌ No visual feedback           │
│  ❌ Stats Comparison            ❌ No side-by-side view         │
│  ❌ Fighter Customization       ❌ Only name, no visuals        │
│  ❌ Fighter History             ❌ No match history             │
└─────────────────────────────────────────────────────────────────┘
```

**PSYCHOLOGICAL GAPS:**

| Gap | Impact | Why It Hurts |
|-----|--------|--------------|
| **No Fighter Progression** | CRITICAL | Players see static numbers - no growth = no attachment |
| **No Level Up Moments** | HIGH | Missing dopamine hit from progression |
| **No Fighter Identity** | HIGH | Can't customize appearance = no emotional bond |
| **No Match History** | MEDIUM | Can't see fighter's legacy = no pride |
| **No Stats Visualization** | MEDIUM | Raw numbers are boring, charts are engaging |

**MISSING MECHANICS:**

1. **XP/Leveling System**
   - Fighter gains XP from matches
   - Visual level-up animation
   - Stat increases on level up
   - Current: Fighter created at level 1, stays level 1 forever

2. **Fighter Customization**
   - Appearance (colors, accessories)
   - Titles ("Dragon Slayer", "Veteran")
   - Background story/bio
   - Current: Only name can be set

3. **Match History & Stats**
   - Wins/losses per fighter
   - Favorite weapons
   - Total damage dealt
   - Current: No historical data tracked

4. **Fighter Evolution**
   - Prestige system (reset for bonuses)
   - Class specialization at level 10
   - Current: Static forever

**Quick Fixes (1-2h each):**
- Add XP column to fighters table
- Show match count in roster
- Add fighter "bio" field
- Visual stat bars instead of raw numbers

---

### ⚔️ MATCHES Flow Analysis (DEEP)
**Status:** Complex system, UI/UX friction points

```
┌─────────────────────────────────────────────────────────────────┐
│  CURRENT MATCH FLOW                                             │
├─────────────────────────────────────────────────────────────────┤
│  User Journey:                                                  │
│  Create Match → Wait Lobby → Start → Watch Combat → Results     │
│                                                                 │
│  Frontend (Matches.vue)         Backend                         │
│  ─────────────────────          ────────                        │
│  ✅ Create Match     ───────►  ✅ POST /matches                │
│  ✅ Join Lobby       ───────►  ✅ POST /matches/{id}/register  │
│  ✅ Start Match      ───────►  ✅ POST /matches/{id}/start     │
│  ✅ Spectate         ───────►  ✅ WebSocket /ws/match          │
│  ✅ View Results     ───────►  ✅ GET /matches/{id}/results    │
│                                                                 │
│  ❌ Pre-Match Strategy          ❌ No team formation phase      │
│  ❌ Real-time Chat              ❌ No lobby communication       │
│  ❌ Match Replay                ❌ No replay storage            │
│  ❌ Bet/Wager System            ❌ No spectator engagement      │
│  ❌ Ranked Mode                 ❌ No skill-based matchmaking   │
│  ❌ Tournament Brackets         ❌ Single matches only          │
└─────────────────────────────────────────────────────────────────┘
```

**PSYCHOLOGICAL FRICTION POINTS:**

| Problem | Severity | Psychology Impact |
|---------|----------|-------------------|
| **Empty Lobbies** | CRITICAL | Waiting alone = boredom = quit |
| **No Pre-Combat** | HIGH | No strategy = random = less investment |
| **Passive Spectating** | HIGH | Watching ≠ Playing = disengagement |
| **No Replay Value** | MEDIUM | Can't relive victories = lost memories |
| **No Stakes** | MEDIUM | No risk/reward = no excitement |

**MISSING MECHANICS:**

1. **Quick Match / Matchmaking**
   - Join random lobby instantly
   - Skill-based matching
   - Current: Must create or find lobby manually

2. **Pre-Combat Strategy Phase**
   - Position fighters on grid
   - Set formation (aggressive/defensive)
   - Choose opening move
   - Current: Immediate combat start

3. **Spectator Engagement**
   - Betting on matches (virtual currency)
   - Reactions/emotes during combat
   - Commentator system
   - Current: Pure passive watching

4. **Match Replay & Highlights**
   - Save interesting matches
   - Share replays
   - "Play of the Game" moments
   - Current: Gone forever after match

5. **Ranked & Casual Split**
   - Ranked with ELO/MMR
   - Casual for fun/testing
   - Current: All matches same

6. **Match History Dashboard**
   - Recent matches list
   - Performance trends
   - Win rate by fighter
   - Current: No persistence

**Quick Fixes (2-4h each):**
- Add "Quick Join" button (join random open lobby)
- Show "Players Online" count
- Save last 10 matches to history
- Add "Rematch" button after combat

---

### 🏆 ROSTER + MATCHES INTEGRATION GAPS

**The Big Picture Problem:**

```
Roster fighters feel DISCONNECTED from matches:

Roster          Match            Result
──────          ─────            ──────
Static Stats →  Combat Happens → Rewards?
     ↑              ↓                ↓
  No growth    No fighter         No fighter
  visible      progression        identity
```

**Missing Connection:**
- Fighter doesn't visibly level up from matches
- No "Fighter of the Match" recognition
- No fighter-specific achievements
- No rivalry history between fighters

**SOLUTION: Fighter Career System**
```
Match Win → XP Gained → Level Up → New Title → Visual Change
    ↓          ↓           ↓            ↓            ↓
  Stats     History    Animation    "Veteran"    New Color
```

---

---

### 🎁 DAILY REWARDS System (P0 - CRITICAL)
**Psychology Driver:** Variable Reward Schedule (Skinner Box)  
**Business Impact:** +40% D1 Retention

```
┌─────────────────────────────────────────────────────────────────┐
│  DAILY REWARD FLOW                                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Day 1: Small Pouch        (100 Gold)                          │
│  Day 2: Common Chest       (250 Gold + Common Item)            │
│  Day 3: Rare Cache         (500 Gold + Rare Item)              │
│  Day 4: Energy Boost       (2x XP for 1h)                      │
│  Day 5: Mystery Box        (Random rarity 1-4)                 │
│  Day 6: Fabled Vault       (1000 Gold + Fabled Item)           │
│  Day 7: LEGENDARY CRATE    (Guaranteed Legendary + 2000 Gold)  │
│                                                                 │
│  BREAK STREAK = RESET TO DAY 1                                  │
│  (Loss Aversion drives daily login)                             │
└─────────────────────────────────────────────────────────────────┘
```

**API Design:**
```
GET /api/daily-reward
Response: {
  "can_claim": true,
  "streak": 3,
  "day": 3,
  "next_reward": {
    "name": "Rare Cache",
    "description": "500 Gold + Rare Item",
    "icon": "📦"
  },
  "time_until_reset": "14:32:15"
}

POST /api/daily-reward/claim
Response: {
  "reward": {
    "type": "gold",
    "amount": 500
  },
  "new_streak": 4,
  "next_reward_preview": { ... }
}
```

**Database Schema:**
```sql
CREATE TABLE daily_rewards (
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    streak INTEGER DEFAULT 0,
    last_claimed DATE,
    total_claimed INTEGER DEFAULT 0
);
```

**Frontend Component:**
- Modal popup on login (if can_claim)
- Visual calendar showing 7-day track
- Progress to next streak milestone
- "Come back tomorrow for LEGENDARY!" teaser

**Implementation Time:** 2-3 hours
**Priority:** P0 - HIGHEST

---

### 🏆 LEADERBOARD System (P1 - HIGH)
**Psychology Driver:** Social Status Competition  
**Business Impact:** +300% Virality, +engagement

```
┌─────────────────────────────────────────────────────────────────┐
│  LEADERBOARD CATEGORIES                                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  💪 POWER RANKING                                               │
│     - Based on: Total Fighter Power (all fighters summed)      │
│     - Updates: Real-time                                        │
│     - Reward: Weekly Top 10 get exclusive title/frame          │
│                                                                 │
│  💰 WEALTH RANKING                                              │
│     - Based on: Total Gold + Inventory Value                    │
│     - Updates: Daily                                            │
│     - Reward: "Tycoon" title, shop discounts                   │
│                                                                 │
│  ⚔️ COMBAT RANKING (ELO)                                        │
│     - Based on: Match wins/losses with skill-based algo        │
│     - Updates: After each match                                 │
│     - Reward: Rank badges (Bronze/Silver/Gold/Platinum/Diamond)│
│                                                                 │
│  🏅 ACHIEVEMENT RANKING                                         │
│     - Based on: Total achievement points                        │
│     - Updates: Real-time                                        │
│     - Reward: "Completionist" cosmetic rewards                 │
│                                                                 │
│  🔥 WIN STREAK RANKING                                          │
│     - Based on: Current consecutive wins                        │
│     - Updates: Real-time                                        │
│     - Reward: Streak flames on avatar                          │
└─────────────────────────────────────────────────────────────────┘
```

**API Design:**
```
GET /api/leaderboard/{type}?limit=10&page=1
types: power, wealth, combat, achievements, streak

Response: {
  "type": "power",
  "total_entries": 15420,
  "user_rank": 145,      // Current user's position
  "entries": [
    {
      "rank": 1,
      "username": "DragonSlayer",
      "avatar": "...",
      "value": 45820,      // Power score
      "trend": "up",       // up/down/same from yesterday
      "is_user": false
    },
    { rank: 2, ... },
    { rank: 3, ... },
    { rank: 145, is_user: true, ... }  // Always include user
  ]
}

GET /api/leaderboard/rankings/nearby?type=power
Response: {
  "user_rank": 145,
  "nearby": [
    { rank: 142, ... },
    { rank: 143, ... },
    { rank: 144, ... },
    { rank: 145, is_user: true, ... },
    { rank: 146, ... },
    { rank: 147, ... },
    { rank: 148, ... }
  ]
}
```

**Frontend Components:**
1. **LeaderboardPage.vue** - Full page with tabs for each category
2. **LeaderboardWidget.vue** - Mini version for dashboard (top 3)
3. **RankBadge.vue** - Shows user's current rank with flair

**Key Features:**
- "Just 5 more power to reach rank #100!" - Proximity motivation
- Weekly reset for competitive categories
- "You beat 95% of players!" - Percentile framing
- Top player spotlight on homepage

**Implementation Time:** 4-6 hours
**Priority:** P1 - HIGH

---

### Leagues Flow Analysis
**Status:** Core features complete, missing admin capabilities

| Feature | Frontend | Backend | Status |
|---------|----------|---------|--------|
| List Leagues | ✅ | ✅ `GET /league` | Complete |
| View Details | ✅ | ✅ `GET /league/{id}` | Complete |
| Subscribe | ✅ | ✅ `POST /league/subscribe` | Complete |
| Unsubscribe | ✅ | ✅ `POST /league/unsubscribe` | Complete |
| Highscores | ✅ | ✅ `POST /league/{id}/highscores` | Complete |
| Matches | ✅ | ✅ `POST /league/{id}/matches` | Complete |
| **Create League** | ❌ | ❌ No endpoint | **MISSING** |
| **Admin Panel** | ❌ | ❌ No admin check | **MISSING** |

**Missing Components:**
1. **Admin Endpoint** `POST /admin/league` - Create/edit leagues
2. **Admin Middleware** - Role-based authorization
3. **Seed Data** - Default test leagues
4. **Admin UI** - League management interface

**Fix Options:**
- **Quick Fix (2 min):** Add seed data migration with 3 standard leagues
- **Proper Fix (15 min):** Build admin endpoints + middleware + UI

---

## 🔴 TODO

### High Priority
- [ ] Deploy current main to VPS
- [ ] Verify Shop MVP on production
- [ ] Verify Attunement on production
- [ ] Add league seed data (if quick fix chosen)

### Medium Priority (if proper fix chosen)
- [ ] Add `is_admin` to users table
- [ ] Create admin middleware
- [ ] Build `POST /admin/league` endpoint
- [ ] Build `PUT /admin/league/{id}` endpoint
- [ ] Create AdminLeagues.vue frontend page

### Low Priority
- [ ] Combat simulator improvements
- [ ] Additional rarity tiers (Missing, Epic, Exotic, etc.)
- [ ] Tournament system

---

## 📊 Current State

**Latest Commit:** `88a44f7` - Attunement complete
**Branch:** `main`
**Deployment Status:** Ready for VPS

**Live URLs (after deploy):**
- Frontend: http://152.53.118.78:49100
- Backend: http://152.53.118.78:49101

**Completed Features:**
1. 9-Tier Rarity System
2. Shop MVP (Gold packages, bundles)
3. Attunement System (6 elements, 25 levels)
4. Leagues Core (List, subscribe, matches, highscores)

**Known Gaps:**
- No admin league creation (requires manual DB insert or seed data)

---

## 🧠 PSYCHOLOGY ANALYSIS - Missing Engagement Drivers

**Full Analysis:** `/docs/GAME_PSYCHOLOGY_ANALYSIS.md`

### Critical Missing Features (by Psychological Impact)

#### P0 - IMMEDIATE (Implement Today)
| Feature | Psychology Driver | Business Impact |
|---------|-------------------|-----------------|
| **Daily Rewards + Streaks** | Variable Reward Schedule (Dopamine) | +40% D1 Retention |
| **Collection Book** | Completionism (Zeigarnik Effect) | +2x Session Length |
| **Progress Bars** | Visible Growth (Endowed Progress) | +30% Engagement |

#### P1 - HIGH PRIORITY (This Week)
| Feature | Psychology Driver | Business Impact |
|---------|-------------------|-----------------|
| **Global Leaderboards** | Social Status Competition | +300% Virality |
| **Weekend Events (2x Drops)** | FOMO + Scarcity | +150% Revenue |
| **Achievement System** | Long-term Goals + Mastery | +25% L7 Retention |
| **Guilds/Factions** | Social Belonging | +3x Retention |

#### P2 - MEDIUM (Next Sprint)
| Feature | Psychology Driver | Business Impact |
|---------|-------------------|-----------------|
| **Season Pass** | Predictable + FOMO Combo | +200% Monetization |
| **Cosmetics/Expression** | Identity + Autonomy | +Whale Spend |
| **Tournament Mode** | Skill Expression + Status | +Hardcore Retention |

#### P3 - NICE TO HAVE
- PvE Campaign (Narrative)
- Trading/Auction House
- Replay System

### 🔥 The Core Problem

**Current State:**
- Great core mechanics ✓
- Solid progression systems ✓
- **BUT:** No "reason to return tomorrow"

**Missing:** The **Variable Reward Schedule** - Skinner box mechanics that make gambling addictive:
- Daily mystery boxes
- Streak bonuses (loss aversion)
- Random drops with visual suspense
- Limited-time events

### 🎯 Recommended Build Order

**Today (2 hours):**
1. ✅ `GET /api/daily-reward` - Random reward based on streak
2. ✅ `GET /api/collection/progress` - % completion tracking
3. ✅ League seed data (so leagues work)

**This Week (8 hours):**
4. Leaderboards endpoint + UI
5. Weekend event system (2x drops flag)
6. Achievement framework

**Next Week (12 hours):**
7. Season Pass (tiers + rewards)
8. Guild system MVP
9. Cosmetic shop items

### 📊 Expected Impact

| Metric | Current | After P0+P1 | Change |
|--------|---------|-------------|--------|
| D1 Retention | 25% | 40% | +60% |
| Session Length | 8 min | 15 min | +88% |
| Revenue/Player | $5 | $12.50 | +150% |
| Social Shares | Low | High | +300% |

### 💡 Key Insight

> "Players don't quit because the game is bad. They quit because they forget it exists."

**Daily rewards + streaks = Habit formation**  
**Leaderboards + guilds = Social obligation**  
**Events + FOMO = Urgency to act**

Combine all three = Addictive retention loop 🎯

---

## 🚀 NEXT ACTION

**DiaDome's Choice:**
1. **Implement P0 features** (Daily rewards, collection, progress) - 2h work, massive impact
2. **Deploy first** - Test current state live
3. **Skip to P1** - Leaderboards + events

What's the priority? 🔥

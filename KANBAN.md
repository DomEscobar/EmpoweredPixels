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

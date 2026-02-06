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

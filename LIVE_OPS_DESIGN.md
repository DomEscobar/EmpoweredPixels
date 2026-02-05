# Live Ops Design - EmpoweredPixels

## Overview
This document defines the live operations systems designed to maximize player engagement, retention, and monetization for EmpoweredPixels Studio.

---

## Daily Quest System

### Quest Types
| Type | Example | Reward |
|------|---------|--------|
| Combat | Win 3 matches | 150 Gold, 300 XP |
| Collection | Equip rare+ weapon | 100 Gold, 200 XP |
| Social | Join guild | 200 Gold, 400 XP |
| Daily Login | Log in today | 50 Gold, 100 XP |
| Hard Mode | Win without taking damage | 400 Gold, 800 XP, Small Lootbox |
| Exploration | Visit 3 different arenas | 75 Gold, 150 XP |

### Quest Pool (20 Total)

**Combat Quests (8):**
1. Win 1 match - 50 Gold, 100 XP
2. Win 3 matches - 150 Gold, 300 XP
3. Win 5 matches - 250 Gold, 500 XP
4. Deal 1000 total damage - 100 Gold, 200 XP
5. Get 3 critical hits - 75 Gold, 150 XP
6. Win 2 matches in a row - 200 Gold, 400 XP
7. Defeat a higher-level player - 300 Gold, 600 XP
8. Win using only basic attacks - 400 Gold, 800 XP, Small Lootbox

**Collection Quests (6):**
9. Equip a rare weapon - 100 Gold, 200 XP
10. Equip 3 different weapons - 150 Gold, 300 XP
11. Upgrade any weapon once - 100 Gold, 200 XP
12. Collect 5 items today - 75 Gold, 150 XP
13. Equip epic+ armor - 200 Gold, 400 XP
14. Max enhance a weapon - 500 Gold, 1000 XP, Medium Lootbox

**Social Quests (6):**
15. Join a guild - 200 Gold, 400 XP
16. Send 3 friend requests - 100 Gold, 200 XP
17. Play 2 matches with friends - 150 Gold, 300 XP
18. Chat in global once - 50 Gold, 100 XP
19. Invite a friend to the game - 300 Gold, 600 XP
20. Participate in guild event - 250 Gold, 500 XP

### Rotation Rules
- 3 quests per day per player
- Reset: 00:00 UTC daily
- Quests selected randomly from pool of 20
- Guaranteed: 1 Combat + 1 Collection + 1 Social
- Hard quests appear 20% of the time (replaces one random quest)
- Duplicate prevention: Same quest won't appear 2 days in a row

### Streak System
| Streak | Bonus Multiplier |
|--------|-----------------|
| 3 days | 1.1x |
| 7 days | 1.25x |
| 14 days | 1.5x |
| 30 days | 2.0x + Epic Lootbox |

---

## Weekend Events Calendar

### Monthly Schedule Template

| Weekend | Event Type | Details | Boost |
|---------|-----------|---------|-------|
| Week 1 | Double XP | All XP gains doubled | 2x XP |
| Week 2 | Drop Boost | Weapon drop rates +50% | +50% drops |
| Week 3 | Tournament | Competitive bracket event | Special rewards |
| Week 4 | Boss Rush | Special boss fights | Exclusive loot |

### Event Details

**Double XP Weekend**
- Duration: Friday 00:00 UTC to Sunday 23:59 UTC
- All XP gains doubled (matches, quests, achievements)
- Special weekend quests with 3x XP
- Visual indicator: XP particles glow gold

**Drop Boost Weekend**
- Weapon drop rates increased by 50%
- Rare+ drop chance doubled
- Special "Lucky" aura for active players
- Guaranteed rare drop every 10 matches

**Tournament Mode**
- 64-player single elimination bracket
- Entry fee: 100 Gold
- Prizes: 
  - 1st: 5000 Gold + Legendary Weapon
  - 2nd: 2500 Gold + Epic Weapon
  - 3rd-4th: 1000 Gold + Rare Weapon
  - 5th-8th: 500 Gold
- Runs Saturday 12:00-20:00 UTC

**Boss Rush Event**
- Special PvE boss fights every 2 hours
- 3 difficulty tiers: Normal, Hard, Nightmare
- Boss drops exclusive cosmetic items
- Guild leaderboards for total damage dealt
- Top guild gets exclusive guild banner

### Seasonal Events

**Spring Festival (March)**
- Cherry blossom arena theme
- Limited "Sakura" weapon skins
- Special quests with flower-themed rewards

**Summer Heat (July)**
- Beach arena variant
- Water-themed weapons
- Double Gold weekend

**Halloween Horror (October)**
- Spooky arena decorations
- Zombie survival mode
- Exclusive horror-themed cosmetics

**Winter Celebration (December)**
- Snow-covered arenas
- Gift-giving system between players
- Limited "Frost" legendary weapons

---

## Achievement System (50 Total)

### Combat Achievements (20)

| # | Name | Requirement | Reward |
|---|------|-------------|--------|
| 1 | First Blood | Win your first match | 100 Gold, Title: "Newbie" |
| 2 | Warrior | Win 10 matches | 500 Gold, 500 XP |
| 3 | Veteran | Win 50 matches | 2000 Gold, Title: "Veteran" |
| 4 | Champion | Win 100 matches | 5000 Gold, Epic Weapon |
| 5 | Legend | Win 500 matches | 20000 Gold, Title: "Legend", Legendary Frame |
| 6 | Unstoppable | Win 10 matches in a row | 3000 Gold, Aura: "Flame" |
| 7 | Dominator | Win 50 matches in a row | 10000 Gold, Title: "Dominator" |
| 8 | Crusher | Deal 10,000 total damage | 1000 Gold |
| 9 | Destroyer | Deal 100,000 total damage | 5000 Gold, Title: "Destroyer" |
| 10 | Annihilator | Deal 1,000,000 total damage | 20000 Gold, Legendary Weapon Skin |
| 11 | Critical Strike | Land 100 critical hits | 500 Gold |
| 12 | Precision Master | Land 1000 critical hits | 3000 Gold, Title: "Sharpshooter" |
| 13 | Comeback Kid | Win with <10% HP | 1000 Gold |
| 14 | Perfect Victory | Win without taking damage | 5000 Gold, Title: "Untouchable" |
| 15 | Speed Demon | Win in under 30 seconds | 2000 Gold |
| 16 | Marathon | Win a match lasting 5+ minutes | 1500 Gold |
| 17 | Giant Slayer | Defeat player 10+ levels higher | 2000 Gold |
| 18 | Ranked Warrior | Reach Bronze rank | 1000 Gold |
| 19 | Ranked Elite | Reach Diamond rank | 10000 Gold, Title: "Elite" |
| 20 | Top 100 | Reach top 100 on leaderboard | 25000 Gold, Exclusive Avatar |

### Collection Achievements (15)

| # | Name | Requirement | Reward |
|---|------|-------------|--------|
| 21 | First Weapon | Collect your first weapon | 100 Gold |
| 22 | Collector | Collect 10 weapons | 500 Gold |
| 23 | Hoarder | Collect 50 weapons | 3000 Gold, Storage Expansion |
| 24 | Completionist | Collect 100 weapons | 10000 Gold, Title: "Collector" |
| 25 | Rare Find | Collect first rare weapon | 300 Gold |
| 26 | Epic Loot | Collect first epic weapon | 1000 Gold |
| 27 | Legendary | Collect first legendary weapon | 5000 Gold, Aura: "Gold" |
| 28 | Enhancer | Upgrade a weapon once | 200 Gold |
| 29 | Enhancer+ | Upgrade a weapon to +5 | 1000 Gold |
| 30 | Max Power | Max enhance a weapon (+10) | 5000 Gold, Title: "Blacksmith" |
| 31 | Fashion Forward | Equip a full cosmetic set | 1000 Gold |
| 32 | Weapon Master | Own one of each weapon type | 3000 Gold |
| 33 | Rich | Accumulate 10,000 Gold | 1000 Gold |
| 34 | Wealthy | Accumulate 100,000 Gold | 5000 Gold, Title: "Tycoon" |
| 35 | Millionaire | Accumulate 1,000,000 Gold | 20000 Gold, Exclusive Mount |

### Social Achievements (10)

| # | Name | Requirement | Reward |
|---|------|-------------|--------|
| 36 | Socialite | Make 5 friends | 500 Gold |
| 37 | Popular | Make 25 friends | 2000 Gold, Title: "Friendly" |
| 38 | Guild Member | Join a guild | 300 Gold |
| 39 | Guild Officer | Become a guild officer | 1000 Gold |
| 40 | Guild Leader | Create/lead a guild | 2000 Gold, Title: "Leader" |
| 41 | Team Player | Play 10 matches with friends | 1000 Gold |
| 42 | Recruiter | Invite 5 friends to the game | 5000 Gold |
| 43 | Mentor | Help 3 new players win matches | 2000 Gold, Title: "Mentor" |
| 44 | Chatty | Send 100 chat messages | 500 Gold |
| 45 | Famous | Get 1000 profile views | 3000 Gold, Title: "Celebrity" |

### Secret Achievements (5)

| # | Name | Requirement | Reward |
|---|------|-------------|--------|
| 46 | Easter Egg | Find the hidden room in arena | 5000 Gold, Title: "Explorer" |
| 47 | ??? | Lose 10 matches in a row | 1000 Gold, Title: "Persistent" |
| 48 | ??? | Stay logged in for 24 hours straight | 2000 Gold, Title: "Dedicated" |
| 49 | ??? | Click your profile 100 times | 1000 Gold, Title: "Narcissist" |
| 50 | ??? | Win with 1 HP remaining | 5000 Gold, Title: "Lucky" |

---

## Analytics KPIs

### Core Retention Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| DAU (Daily Active Users) | Growth 5%/week | <10% of MAU |
| MAU (Monthly Active Users) | Growth 10%/month | Decline 2 weeks |
| D1 Retention | >40% | <35% |
| D7 Retention | >20% | <15% |
| D30 Retention | >10% | <8% |

### Engagement Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| Avg Session Length | >15 minutes | <10 minutes |
| Sessions per User/Day | >2.5 | <1.5 |
| Quest Completion Rate | >70% | <50% |
| Weekend Event Participation | >30% of DAU | <15% |
| Achievement Unlock Rate | >5 achievements/player | <2 |

### Economy Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| Gold per Player/Day | 500-1000 | <300 or >2000 |
| XP per Player/Day | 1000-3000 | <500 or >5000 |
| Lootbox Open Rate | >60% of earned | <30% |
| Marketplace Activity | >20% list daily | <5% |
| Premium Purchase Rate | >5% of MAU | <2% |

### Social Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| Guild Participation | >40% in guilds | <25% |
| Friend Requests/Day | >1000 total | <100 |
| Chat Messages/Day | >10000 total | <1000 |
| Tournament Signups | >15% of eligible | <5% |

### Monitoring Dashboards

**Real-time (5-min refresh):**
- Current online players
- Active matches
- Quest completions (last hour)
- Error rates

**Daily (00:00 UTC):**
- DAU, new registrations, churn
- Quest completion breakdown
- Event participation summary
- Revenue

**Weekly (Monday 00:00 UTC):**
- Retention cohorts (D1, D7, D30)
- Achievement unlock rates
- Economy health check
- Leaderboard snapshots

---

## Economy Monitoring Framework

### Currency Flow Tracking

**Gold Sources:**
| Source | Daily Cap | % of Economy |
|--------|-----------|--------------|
| Match Wins | 500/player | 40% |
| Quests | 400/player | 30% |
| Achievements | Unlimited | 15% |
| Events | Variable | 15% |

**Gold Sinks:**
| Sink | Cost Range | Target Drain |
|------|-----------|--------------|
| Weapon Upgrades | 100-10000 | 25% |
| Marketplace Tax | 10% | 20% |
| Tournament Entry | 100-500 | 15% |
| Cosmetics | 500-5000 | 25% |
| Lootboxes | 100-300 | 15% |

### Inflation Control
- Monthly economy rebalance review
- Adjust quest rewards if inflation >5%/month
- Introduce new gold sinks quarterly
- Special "Gold Drain" events when reserves high

### Anti-Farming Measures
- Daily caps on match rewards (prevent botting)
- Quest variety prevents repetitive farming
- Diminishing returns on repeated activities

---

## Leaderboard System

### Categories

**Global Leaderboards:**
1. **All-Time Wins** - Total match wins (lifetime)
2. **Season Rating** - Competitive ranked score (resets monthly)
3. **Weekly Kills** - Total eliminations (resets Monday 00:00 UTC)
4. **Guild Power** - Total guild member ratings

### Rewards

**Weekly Top 100:**
| Rank | Reward |
|------|--------|
| 1 | 10000 Gold + Legendary Weapon + Title "#1" |
| 2-3 | 5000 Gold + Epic Weapon + Title "Elite" |
| 4-10 | 2500 Gold + Rare Weapon |
| 11-50 | 1000 Gold |
| 51-100 | 500 Gold |

**Season Rewards (Monthly):**
| Rank Tier | Reward |
|-----------|--------|
| Top 10 | Exclusive Seasonal Skin + 25000 Gold |
| Top 100 | Epic Skin + 10000 Gold |
| Top 1000 | Rare Skin + 5000 Gold |
| Top 10% | 2000 Gold |

### Display Rules
- Leaderboards update every 15 minutes
- Player must have played in last 7 days to appear
- Anti-cheat verification before finalizing weekly rewards
- Guild leaderboards require 5+ members

---

## Implementation Priority

### Phase 1 (Week 1-2)
- [ ] Daily Quest System backend
- [ ] Quest UI/UX
- [ ] Basic Analytics tracking

### Phase 2 (Week 3-4)
- [ ] Achievement System (all 50)
- [ ] Achievement UI
- [ ] First Weekend Event (Double XP)

### Phase 3 (Week 5-6)
- [ ] Leaderboard System
- [ ] Full Event Calendar
- [ ] Advanced Analytics Dashboard

### Phase 4 (Week 7-8)
- [ ] Economy monitoring tools
- [ ] Seasonal event content
- [ ] Secret achievements activation

---

## Success Metrics

**30-Day Targets:**
- DAU increase: +25%
- D7 retention: >20%
- Quest completion: >65%
- Weekend event participation: >25%
- Average achievements unlocked: 8 per player

**90-Day Targets:**
- D30 retention: >10%
- MAU growth: +50%
- Average session length: >20 minutes
- 40% of players in guilds
- 5+ tournament participants per event

---

*Document Version: 1.0*
*Created: 2026-02-05*
*Owner: Live Ops Manager*

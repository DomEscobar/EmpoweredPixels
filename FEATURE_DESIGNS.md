# 🎮 Feature Designs - EmpoweredPixels

## 1. Weapon System

### User Stories
- As a player, I want to collect different weapons with varying rarities, so that I can optimize my fighter's combat style
- As a player, I want to upgrade my weapons through enhancement, so that I can progress without relying only on drops
- As a player, I want to see my weapon collection, so that I can feel accomplishment

### Acceptance Criteria
- [ ] Weapon Domain with Type, Rarity, Stats, Level
- [ ] 5 Rarities: Common (white), Rare (blue), Epic (purple), Legendary (gold), Mythic (red)
- [ ] 5 Weapon Types: Sword (balanced), Bow (range), Staff (magic), Dagger (speed), Axe (power)
- [ ] Each weapon has: Damage (10-100), Attack Speed (0.8-1.5), Crit Chance (5-30%), **Weapon Skill/Proc**
- [ ] **Weapon Sets: Equip 2+ weapons from same "Forge" for set bonuses**
- [ ] **Salvage System: Convert unwanted weapons into Enhancement Dust**
- [ ] **Favorites System: Mark weapons to prevent accidental salvage/enhancement**
- [ ] Enhancement system: +1 to +10, each level +10% base stats
- [ ] Enhancement failure chance: +1-3 (0%), +4-6 (15%), +7-9 (35%), +10 (50%)
- [ ] **On failure at +4-6: No penalty** | **On failure at +7-9: Drop to +5** | **On failure at +10: Drop to +7 (NOT +0!)**
- [ ] **Protection Scrolls: Optional consumable that prevents downgrade on failure (rare drop/craftable)**
- [ ] Weapon inventory: **100 slots base**, expandable to 200 for premium
- [ ] Weapon equipping: 1 weapon per fighter, switchable between matches
- [ ] **Compare Mode: Side-by-side stat comparison when selecting new weapon**
- [ ] Visual: Different sprites per rarity using vibemedia.space API

### Balancing Table
| Rarity | Base Damage | Speed | Crit % | **Weapon Skill** | Drop Rate | Enhance Cost | Salvage Dust |
|--------|-------------|-------|--------|------------------|-----------|--------------|--------------|
| Common | 10-20 | 1.0 | 5% | None | 55% | 100 gold | 10 |
| Rare | 20-35 | 1.1 | 10% | Minor Proc (5% on hit) | 28% | 250 gold | 25 |
| Epic | 35-55 | 1.2 | 15% | Skill Slot 1 | 14% | 500 gold | 60 |
| Legendary | 55-80 | 1.3 | 22% | Skill Slot 1+2 | 2.9% | 1000 gold | 150 |
| Mythic | 80-100 | 1.5 | 30% | **Unique Passive** + Skill Slots | 0.1% | 2000 gold | 400 |

### Weapon Skill Examples (Procs)
| Type | Skill Name | Proc Chance | Effect |
|------|------------|-------------|--------|
| Sword | Cleave | 15% | Hit 2 adjacent enemies |
| Bow | Piercing Shot | 15% | Arrow pierces through target |
| Staff | Mana Surge | 20% | Restore 15 mana on hit |
| Dagger | Poison Blade | 20% | 8 damage/s for 4s |
| Axe | Crushing Blow | 12% | Stun target for 1.5s |

### Weapon Set System (Collection Depth)
| Set Name | Weapons Required | Bonus |
|----------|------------------|-------|
| Iron Legion | 3 Iron-forged | +5% damage |
| Dragon Slayer | 2 Dragon weapons | +10% vs bosses |
| Shadow Assassin | 4 Shadow weapons | +15% crit chance |
| Ancient Relic | 5 Ancient | Unlock hidden dungeon |

### MVP Scope
- 20 weapons (4 per type, distributed across rarities)
- Basic enhancement UI with risk indicator
- Inventory grid view
- Equip/unequip functionality
- Visual sprites for each rarity

### Dependencies
- Requires: Equipment system (exists), Inventory system (exists)
- Blocks: Shop (weapons need pricing), Crafting (needs materials)
- Assets: vibemedia.space for weapon sprites

### API Endpoints Needed
```
GET /api/weapons - List inventory
POST /api/weapons/equip - Equip weapon
POST /api/weapons/enhance - Enhance weapon (+1 level)
GET /api/weapons/forge - Preview enhancement odds
```

---

## 2. Skill System

### User Stories
- As a player, I want to unlock skills as I level up, so that my fighter becomes more powerful
- As a player, I want to choose which skills to activate, so that I can customize my playstyle
- As a player, I want skills to have cooldowns, so that combat requires strategy

### Acceptance Criteria
- [ ] Skill Tree with 3 branches: Offense (damage), Defense (survival), Utility (support)
- [ ] 5 tiers per branch, 2 skills per tier = 30 total skills
- [ ] **Skill Prerequisites: Tier 2+ requires 2 points in previous tier**
- [ ] **Skill Ranks: Each skill has 3 ranks (improved by spending additional points)**
- [ ] Skill points: 1 per level, max 100 at level cap, **+10 bonus points from achievements**
- [ ] Max 3 active skills per match (loadout system), **+1 slot at level 50 (Ultimate Slot)**
- [ ] Skills have: Mana cost, Cooldown (5-60s), Effect duration, Damage/heal values
- [ ] Ultimate skill: Charges through combat, releases for massive effect
- [ ] **Ultimate Charge: Builds from damage dealt (1% per 50 damage) AND damage taken (1% per 25 damage)**
- [ ] Skill synergies: Using Skill A then B within 5s = combo effect
- [ ] **Quick-Swap Loadouts: 3 preset loadouts, switchable before match**
- [ ] **Skill Preview: Test skills in training dummy arena before committing points**
- [ ] Reset: Can reset skill tree for gold cost, **first reset free**

### Skill Examples (with 3 Ranks Each)
**Offense Branch:**
| Tier | Skill | Mana | CD | Rank 1 | Rank 2 | Rank 3 |
|------|-------|------|-----|--------|--------|--------|
| 1 | Power Strike | 20 | 8s | +40% damage | +55% damage | +70% damage |
| 2 | Bleed | 25 | 12s | 5 dmg/s x 4s | 8 dmg/s x 5s | 12 dmg/s x 6s |
| 3 | Whirlwind | 40 | 15s | 80% weapon dmg | 100% weapon dmg | 120% weapon dmg |
| 4 | Berserk | 50 | 30s | +20% dmg, -15% armor | +30% dmg, -10% armor | +40% dmg, -5% armor |
| 5 | Execute | 70 | 45s | Kill below 15% HP | Kill below 20% HP | Kill below 25% HP |

**Defense Branch:**
| Tier | Skill | Mana | CD | Rank 1 | Rank 2 | Rank 3 |
|------|-------|------|-----|--------|--------|--------|
| 1 | Block | 15 | 8s | Block 1 attack | Block 2 attacks | Block 3 attacks |
| 2 | Heal | 30 | 18s | 15% HP | 25% HP | 35% HP |
| 3 | Shield Wall | 45 | 25s | +40% armor 6s | +50% armor 8s | +60% armor 10s |
| 4 | Reflect | 40 | 22s | Return 20% dmg | Return 30% dmg | Return 40% dmg |
| 5 | Immortal | 90 | 75s | 3s duration | 4s duration | 5s duration |

**Ultimate (Universal) - Unlocked at Level 50:**
| Skill | Charge Rate | Effect |
|-------|-------------|--------|
| Meteor Strike | 1% per 40 dmg dealt | 250 AOE damage to all enemies |
| Divine Protection | 1% per 20 dmg taken | Full heal + invincible 4s |
| Time Warp | 1% per 60 dmg dealt | Reset cooldowns, +40% speed 10s |

### Skill Combo System (Synergy)
| First Skill | Second Skill | Combo Effect | Window |
|-------------|--------------|--------------|--------|
| Power Strike | Execute | Execute threshold +10% | 5s |
| Bleed | Whirlwind | Whirlwind applies bleed | 5s |
| Block | Reflect | Next block also reflects | 4s |
| Heal | Shield Wall | Shield Wall heals 10% HP | 4s |

### MVP Scope
- 15 skills (5 per branch, tier 1-3 only)
- 2 active skill slots (expandable to 3)
- Basic skill tree UI
- Cooldown timers in combat

### Dependencies
- Requires: Combat system (exists), Fighter level system (exists)
- Blocks: None
- Assets: Skill icons from vibemedia.space

---

## 3. Attunement System

### User Stories
- As a player, I want to choose an elemental attunement, so that I can specialize my fighter
- As a player, I want elements to have strengths and weaknesses, so that strategy matters
- As a player, I want to level up my attunement, so that I feel progression

### Acceptance Criteria
- [ ] 6 Attunements: Fire, Water, Earth, Wind, Light, Dark
- [ ] Element cycle: Fire > Wind > Earth > Water > Fire, Light <> Dark
- [ ] **Dual Attunement: Unlock at Level 50 to equip secondary attunement (50% weaker bonuses)**
- [ ] Attunement Level: 1-50, XP gained through combat
- [ ] Each attunement has: **3 Passive tiers**, **2 Active abilities** (basic + ultimate), Elemental damage type
- [ ] Strong against: +25% damage, Weak against: -25% damage, **Neutral: No modifier**
- [ ] **Attunement Resonance: Matching weapon element to attunement = +10% weapon damage**
- [ ] Change attunement: **12h cooldown** or 300 gold instant, **First change free**
- [ ] **Attunement Mastery: Max all 6 to unlock "Avatar" title + special cosmetic**
- [ ] Visual effects: Particle aura matching attunement, **intensity scales with level**

### Attunement Details (3 Passive Tiers, 2 Active Abilities)
| Attunement | Passive Lv1 | Passive Lv25 | Passive Lv50 | Active (Basic) | Active (Ultimate) | Strong vs | Weak vs |
|------------|-------------|--------------|--------------|----------------|-------------------|-----------|---------|
| Fire | +5% dmg | +10% dmg, +5% crit | +15% dmg, +10% crit, burn spreads | Fireball (50 dmg) | **Inferno** (150 AOE + 10s burn) | Wind | Water |
| Water | +5% healing | +10% heal, +5% mana regen | +15% heal, +10% mana, cleanse debuffs | Healing Wave (15% HP) | **Tidal Wave** (30% HP team heal) | Fire | Earth |
| Earth | +8% armor | +15% armor, +5% HP | +20% armor, +10% HP, thorns dmg | Stone Shield (block 2 hits) | **Earthquake** (stun all enemies 3s) | Water | Wind |
| Wind | +8% speed | +15% speed, +5% dodge | +20% speed, +10% dodge, move while casting | Gust (knockback + 30 dmg) | **Tempest** (dodge all attacks 5s) | Earth | Fire |
| Light | +5% crit | +10% crit, +5% healing | +15% crit, +10% heal, revive once per match | Holy Light (80 dmg + heal 10%) | **Ascension** (invincible 3s + full heal) | Dark | - |
| Dark | +8% lifesteal | +15% lifesteal, +5% dmg | +20% lifesteal, +10% dmg, hp drain aura | Shadow Bolt (60 dmg + fear 1s) | **Apocalypse** (200 dmg + life drain 50%) | Light | - |

### XP Table (Adjusted - More Rewarding Curve)
| Level | XP Required | Cumulative | Unlock |
|-------|-------------|------------|--------|
| 1-10 | 80 per level | 800 | Basic Active |
| 11-25 | 150 per level | 2,900 | Passive Tier 2 at 15 |
| 26-40 | 300 per level | 7,400 | Ultimate Active at 30 |
| 41-50 | 600 per level | 13,400 | Passive Tier 3 at 50, Dual Attunement unlocked |

### Weapon-Attunement Synergy
| Weapon Type | Best Attunement | Synergy Bonus |
|-------------|-----------------|---------------|
| Sword | Fire | Power Strike burns target |
| Bow | Wind | +20% projectile speed |
| Staff | Water | Mana costs reduced 15% |
| Dagger | Dark | Lifesteal applies to bleeds |
| Axe | Earth | Stunning Blow also roots |

### MVP Scope
- 4 Attunements (Fire, Water, Earth, Wind)
- Level cap 25
- Basic passive bonuses only
- Selection UI before match

### Dependencies
- Requires: Combat system (exists), Fighter system (exists)
- Blocks: None
- Assets: Attunement aura effects from vibemedia.space

---

## 4. Shop System

### User Stories
- As a player, I want to buy items with in-game currency, so that I can progress faster
- As a player, I want cosmetic options, so that I can customize my fighter's appearance
- As a player, I want fair pricing, so that the game doesn't feel pay-to-win

### Acceptance Criteria
- [ ] 2 Currencies: Gold (earned), Gems (premium)
- [ ] Shop categories: Weapons, Consumables, Cosmetics, Boosts, **Bundles**
- [ ] Daily rotation: 6 items refresh every 24h at 00:00 UTC, **1 guaranteed Epic+**
- [ ] **Featured Deals: Weekly spotlight item at 20-50% discount**
- [ ] **Wishlist: Mark items, get notified when on sale or in rotation**
- [ ] **Preview System: Try cosmetics on your fighter before buying**
- [ ] **Buyback: Undo accidental purchases within 5 minutes**
- [ ] F2P friendly: Everything gameplay-relevant is buyable with gold
- [ ] Premium convenience: XP boosts, instant upgrades, cosmetics
- [ ] No loot boxes: Direct purchase only
- [ ] **Salvage Shop: Convert unwanted items to shop credit (50% value)**
- [ ] Price transparency: No hidden costs, **show "time to earn" for gold items**

### Pricing
**Weapons (Gold):**
| Rarity | Price Range | Est. Farm Time* |
|--------|-------------|-----------------|
| Common | 500-1000 | 1-2 matches |
| Rare | 2000-3500 | 4-7 matches |
| Epic | 8000-15000 | 16-30 matches |
| Legendary | 40000-60000 | 80-120 matches |
| Mythic | **Not sold** | Boss drop only |

*Based on average 500 gold per match

**Bundles (Weekly Rotation):**
| Bundle | Contents | Gold Price | Discount |
|--------|----------|------------|----------|
| Starter Pack | 3 Rare weapons, 10 potions, 500 dust | 4,000 | 25% off |
| Enhancer's Kit | 5 Scrolls, 2000 dust, 1 Protection Scroll | 3,500 | 30% off |
| Collector's Set | 1 Epic weapon + matching cosmetic | 12,000 | 20% off |
| Ultimate Arsenal | 1 Legendary + 2 Epic weapons | 50,000 | 15% off |

**Consumables:**
| Item | Gold | Effect |
|------|------|--------|
| Health Potion | 50 | Heal 50% HP |
| Mana Potion | 50 | Restore 50 mana |
| Enhancement Scroll | 200 | +5% enhance success |
| XP Boost (1h) | 100 | +50% XP gain |

**Cosmetics (Gems):**
| Item | Gems | Effect |
|------|------|--------|
| Fighter Skin | 500 | Visual change |
| Weapon Skin | 300 | Visual change |
| Aura Effect | 200 | Visual effect |
| Title | 100 | Display title |

**Premium (Gems):**
| Item | Gems | Effect |
|------|------|--------|
| Inventory Expand | 200 | +10 slots |
| Instant Attunement Change | 100 | Skip 24h cooldown |
| Skill Reset | 150 | Free skill tree reset |
| Battle Pass | 1000 | Seasonal rewards |

### MVP Scope
- Gold-only purchases
- Weapons + Consumables only
- Static shop (no daily rotation)
- No premium currency yet

### Dependencies
- Requires: Weapon system, Inventory system
- Blocks: None (but needs other systems first)
- Assets: Shop UI, item icons

### Monetization Principles
1. **No P2W:** Premium only buys convenience/cosmetics - NO direct power purchases
2. **Fair Grind:** F2P players can achieve everything with reasonable time investment (Legendary weapon = ~2 weeks casual play)
3. **Transparent:** All prices visible, no gambling, drop rates public
4. **Value:** Premium feels worth it but not required - $5-10/month optional subscription for QoL
5. **Respect Player Time:** No artificial wait timers on core gameplay
6. **Meaningful Progress:** Every session feels rewarding, even short ones
7. **Catch-Up Mechanics:** New players can reach endgame in 30-60 days

---

## Implementation Priority

1. **Weapon System** - P1 (blocks others, high impact)
2. **Skill System** - P2 (adds depth, medium effort)
3. **Attunement System** - P2 (strategic layer, medium effort)
4. **Shop** - P3 (monetization, implement last)

## Recommended Start Order

**Week 1:** Weapon System MVP
- Basic weapon domain
- 20 weapons with stats
- Inventory UI
- Equip/unequip

**Week 2:** Weapon Enhancement + Skill System MVP
- Enhancement system
- 15 basic skills
- Skill tree UI

**Week 3:** Attunement System + Polish
- 4 attunements
- Elemental effects
- Balance tuning

**Week 4:** Shop MVP
- Basic shop UI
- Gold purchases
- Daily deals

---

---

## 🎮 Design Review Summary - Enhanced by Senior Game Designer

### Key Improvements Made

#### 1. Better Balancing
- **Weapon Enhancement:** Changed brutal +0 failure to graduated penalties (+7→+5, +10→+7)
- **Attunement XP:** Smoother curve, faster early levels, more rewarding milestones
- **Skill System:** Added ranks (3 per skill) for meaningful progression choices
- **Shop Pricing:** Added "time to earn" estimates for transparency

#### 2. Missing Depth Added
- **Weapon Skills/Procs:** Every weapon now has special abilities beyond stats
- **Weapon Sets:** Collection bonuses encourage building diverse arsenals
- **Salvage System:** Unwanted weapons become crafting materials
- **Skill Combos:** Synergy system rewards tactical play
- **Dual Attunement:** Endgame unlock adds build diversity
- **Attunement Mastery:** Long-term goal for completionists
- **Bundles:** Shop now has value-focused package deals

#### 3. UX Improvements
- **Favorites System:** Prevent accidental salvage of cherished weapons
- **Compare Mode:** Side-by-side weapon stat comparison
- **Skill Preview:** Test before you invest points
- **Quick-Swap Loadouts:** 3 preset builds for different situations
- **Wishlist:** Never miss a sale on wanted items
- **Preview System:** Try before you buy cosmetics
- **Buyback:** 5-minute undo for accidental purchases
- **First Reset Free:** Remove fear of experimentation
- **Protection Scrolls:** Optional safety net for enhancement

### Risk Assessment
| Feature | Risk Level | Mitigation |
|---------|------------|------------|
| Enhancement Failure | Medium | Graduated penalties + protection scrolls |
| Dual Attunement | Low | Late-game unlock only (level 50+) |
| Skill Tree Reset Cost | Low | First reset free, gold cost scales with level |
| Mythic Drop Rate (0.1%) | Low | Expected 1 per ~1000 drops, chase item |

### Recommended Next Steps
1. **Prototype Weapon Skills First** - High impact, core to combat feel
2. **A/B Test Enhancement Rates** - Player sentiment on risk vs reward
3. **Soft-Launch with Reduced XP** - Monitor progression speed, adjust before full launch

---

*Design Date: 2026-02-05*
*Original Designer: Senior Game Designer (via Mama)*
*Enhancement Review: Senior Game Designer (Proactive Review)*

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
- [ ] Each weapon has: Damage (10-100), Attack Speed (0.8-1.5), Crit Chance (5-30%), Durability
- [ ] Enhancement system: +1 to +10, each level +10% base stats
- [ ] Enhancement failure chance: +1-3 (0%), +4-6 (15%), +7-9 (35%), +10 (50%)
- [ ] On failure above +5: Weapon drops to +0 (risk mechanic)
- [ ] Weapon inventory: 50 slots, expandable to 100 for premium
- [ ] Weapon equipping: 1 weapon per fighter, switchable between matches
- [ ] Visual: Different sprites per rarity using vibemedia.space API

### Balancing Table
| Rarity | Base Damage | Speed | Crit % | Drop Rate | Enhance Cost |
|--------|-------------|-------|--------|-----------|--------------|
| Common | 10-20 | 1.0 | 5% | 60% | 100 gold |
| Rare | 20-35 | 1.1 | 10% | 25% | 250 gold |
| Epic | 35-55 | 1.2 | 15% | 12% | 500 gold |
| Legendary | 55-80 | 1.3 | 22% | 2.8% | 1000 gold |
| Mythic | 80-100 | 1.5 | 30% | 0.2% | 2000 gold |

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
- [ ] Skill points: 1 per level, max 100 at level cap
- [ ] Max 3 active skills per match (loadout system)
- [ ] Skills have: Mana cost, Cooldown (5-60s), Effect duration, Damage/heal values
- [ ] Ultimate skill: Charges through combat, releases for massive effect
- [ ] Skill synergies: Using Skill A then B within 5s = combo effect
- [ ] Reset: Can reset skill tree for gold cost

### Skill Examples
**Offense Branch:**
| Tier | Skill | Mana | CD | Effect |
|------|-------|------|-----|--------|
| 1 | Power Strike | 20 | 8s | +50% damage next attack |
| 2 | Bleed | 25 | 12s | Target takes 5 damage/s for 6s |
| 3 | Whirlwind | 40 | 15s | AOE damage to all adjacent |
| 4 | Berserk | 60 | 30s | +30% damage, -20% armor for 10s |
| 5 | Execute | 80 | 45s | Instantly kill target below 20% HP |

**Defense Branch:**
| Tier | Skill | Mana | CD | Effect |
|------|-------|------|-----|--------|
| 1 | Block | 15 | 6s | Negate next attack |
| 2 | Heal | 30 | 15s | Restore 20% HP |
| 3 | Shield Wall | 50 | 25s | +50% armor for 8s |
| 4 | Reflect | 45 | 20s | Return 30% damage to attacker for 5s |
| 5 | Immortal | 100 | 60s | Cannot die for 5s |

**Ultimate (Universal):**
| Skill | Charge | Effect |
|-------|--------|--------|
| Meteor Strike | 100 hits | 200 AOE damage to all enemies |
| Divine Protection | 100 hits | Full heal + invincible 5s |
| Time Warp | 100 hits | Reset all cooldowns, +50% speed 10s |

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
- [ ] Attunement Level: 1-50, XP gained through combat
- [ ] Each attunement has: Passive bonus, Active ability, Elemental damage type
- [ ] Strong against: +25% damage, Weak against: -25% damage
- [ ] Change attunement: 24h cooldown or 500 gold instant
- [ ] Visual effects: Particle aura matching attunement

### Attunement Details
| Attunement | Passive | Active | Strong vs | Weak vs |
|------------|---------|--------|-----------|---------|
| Fire | +10% damage | Burn (DoT 10/s for 5s) | Wind | Water |
| Water | +10% healing | Heal (20% HP) | Fire | Earth |
| Earth | +15% armor | Stun (2s) | Water | Wind |
| Wind | +15% speed | Dodge next 3 attacks | Earth | Fire |
| Light | +10% crit | Divine Shield (block 1 hit) | Dark | - |
| Dark | +15% lifesteal | Fear (enemy flees 3s) | Light | - |

### XP Table
| Level | XP Required | Cumulative |
|-------|-------------|------------|
| 1-10 | 100 per level | 1,000 |
| 11-25 | 200 per level | 4,000 |
| 26-40 | 400 per level | 10,000 |
| 41-50 | 800 per level | 18,000 |

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
- [ ] Shop categories: Weapons, Consumables, Cosmetics, Boosts
- [ ] Daily rotation: 6 items refresh every 24h at 00:00 UTC
- [ ] F2P friendly: Everything gameplay-relevant is buyable with gold
- [ ] Premium convenience: XP boosts, instant upgrades, cosmetics
- [ ] No loot boxes: Direct purchase only
- [ ] Price transparency: No hidden costs

### Pricing
**Weapons (Gold):**
| Rarity | Price Range |
|--------|-------------|
| Common | 500-1000 |
| Rare | 2000-3500 |
| Epic | 8000-15000 |
| Legendary | 40000-60000 |

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
1. **No P2W:** Premium only buys convenience/cosmetics
2. **Fair Grind:** F2P players can achieve everything with time
3. **Transparent:** All prices visible, no gambling
4. **Value:** Premium feels worth it but not required

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

*Design Date: 2026-02-05*
*Designer: Senior Game Designer (via Mama)*

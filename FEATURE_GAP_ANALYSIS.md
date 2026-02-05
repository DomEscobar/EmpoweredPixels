# 🔍 Feature Gap Analysis - Attunements, Skills, Weapons, Shop

## Executive Summary

| Feature | Implementiert | Fehlend | Priorität |
|---------|---------------|---------|-----------|
| **Attunements** | 20% | Element-System (5+ Attunements), Selection UI, Balance | P2 |
| **Skills** | 30% | Skill-Tree, Skill-Selection, 20+ Skills, Effects | P2 |
| **Weapons** | 25% | Weapon-Domain, Stats, Inventory, Rarity, 20+ Weapons | P2 |
| **Shop** | 0% | Alles | P3 |

---

## 1. Attunements System 🟡

### Current State
```go
// backend/internal/domain/combat/attunements.go
- Interface definiert
- FireAttunement implementiert (10% Burn Chance)
- In DB gespeichert (attunement_id)
- Nur 1 von 5+ geplanten Attunements
```

### What's Missing
- [ ] **Attunement Selection**: UI + API für Spieler-Auswahl
- [ ] **Element Cycle**: Fire > Wind > Earth > Water > Fire
- [ ] **5+ Attunements**: Fire, Water, Wind, Earth, Light, Dark
- [ ] **Attunement Effects**: 
  - Fire: Burn DoT
  - Water: Heal/Slow
  - Wind: Speed/Dodge
  - Earth: Armor/Stun
  - Light: Crit/Divine Shield
  - Dark: Life steal/Fear
- [ ] **Attunement Leveling**: XP pro Attunement
- [ ] **Visual Effects**: Particles je nach Attunement

### Acceptance Criteria
- Spieler kann Attunement vor Match wählen
- Attunements haben klare Stärken/Schwächen
- Attunement-Wechsel hat Cooldown/Kosten
- Combat zeigt Attunement-Effects

---

## 2. Skills System 🟡

### Current State
```go
// backend/internal/domain/combat/skills.go
- Interface definiert
- 4 Basic Skills (BowShot, DaggerSlice, GlaiveSwing, GreatswordBlow)
- Skills sind hardcoded pro Waffe
- Keine Skill-Progression
```

### What's Missing
- [ ] **Skill-Tree**: 3 Branches pro Klasse (Offense, Defense, Utility)
- [ ] **Skill Points**: Pro Level-Up Skill-Punkte
- [ ] **20+ Skills**: 
  - Warrior: Berserk, Shield Bash, Whirlwind
  - Mage: Fireball, Frost Nova, Teleport
  - Rogue: Backstab, Smoke Bomb, Poison Blade
  - Ranger: Multishot, Trap, Camouflage
- [ ] **Skill Selection**: Spieler wählt 3 Skills für Match
- [ ] **Skill Cooldowns**: Balanced ability rotation
- [ ] **Skill Combos**: Skill A → Skill B = Combo Effect
- [ ] **Ultimate Skills**: Charged super-abilities

### Acceptance Criteria
- Jeder Fighter hat Skill-Tree
- Skill-Punkte können verteilt werden
- Max 3 Skills pro Match
- Skills haben Cooldowns + Mana/Kosten
- Ultimate lädt sich auf durch Combat

---

## 3. Weapons System 🟡

### Current State
```go
// backend/internal/domain/combat/skills.go:103
- WeaponID String in Fighter
- 4 Hardcoded Waffen in Switch-Statement
- Keine Weapon-Domain
- Keine Weapon-Stats
```

### What's Missing
- [ ] **Weapon Domain**: Type, Rarity, Stats, Level
- [ ] **Weapon Inventory**: Spieler sammelt Waffen
- [ ] **Weapon Stats**: Damage, Speed, Crit, Range
- [ ] **20+ Weapons**:
  - Common: Rusty Sword, Wooden Bow
  - Rare: Enchanted Blade, Elven Longbow
  - Epic: Dragon Slayer, Void Staff
  - Legendary: Excalibur, Infinity Bow
- [ ] **Weapon Enhancement**: +1 bis +10 Upgrade
- [ ] **Weapon Crafting**: Materials → Weapon
- [ ] **Weapon Trading**: Marketplace zwischen Spielern

### Acceptance Criteria
- Waffen haben Stats die Combat beeinflussen
- Rarity-System (Common → Legendary)
- Enhancement-System mit Risiko/Belohnung
- Weapon-Skins für Visuals
- Max 1 Waffe pro Fighter

---

## 4. Shop System 🔴

### Current State
```
NICHT IMPLEMENTIERT
```

### What's Missing
- [ ] **Currency System**: Gold, Gems, EP-Tokens
- [ ] **Item Categories**: Weapons, Armor, Consumables, Skins
- [ ] **Pricing**: Balance zwischen Grind und Convenience
- [ ] **Daily Deals**: Rotierende Angebote
- [ ] **Premium Currency**: Gems für Echtgeld
- [ ] **Battle Pass**: Saisonale Progression
- [ ] **Loot Boxes**: Random Rewards (reguliert)

### Acceptance Criteria
- Shop UI klar und intuitiv
- F2P-freundlich (kein P2W)
- Kosmetika kaufbar
- Convenience-Items (XP-Boosts)
- Transparente Preise

---

## Recommended Priority

1. **Weapons** (P2) - Direkter Combat-Impact, leicht zu balancen
2. **Skills** (P2) - Tiefe für Combat, längere Entwicklung
3. **Attunements** (P2) - Strategische Tiefe, nach Weapons/Skills
4. **Shop** (P3) - Monetarisierung, erst nach Gameplay-Festigung

---

## Next Steps

1. **Game Designer**: Detailed Design für Weapon System
2. **PO-Lead**: Priorisierung in Roadmap
3. **Architect**: Implementation mit Tests
4. **QA-Lead**: Unit/Integration/E2E/MCP Tests

---

*Analysis Date: 2026-02-05*
*Analyst: Code-Architect-Agent*

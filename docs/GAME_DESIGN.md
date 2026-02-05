# EmpoweredPixels - Game Design Document
## Vision: Das ultimative AI-native Web3 RPG

**Version:** 1.0  
**Stand:** 2026-02-05  
**Zielgruppe:** Web3-Native Gamer, RPG-Fans, Crypto-Enthusiasten  

---

## 🎯 Core Vision

EmpoweredPixels ist ein **taktisches RPG** mit rundenbasierten Kämpfen, das durch KI-gesteuerte Gegner und dynamische Kampfsysteme überzeugt. Spieler sammeln, verbessern und kämpfen mit einzigartigen Kämpfern in einer lebendigen, spielergesteuerten Ökonomie.

**Unique Selling Points:**
- 9-Stufen-Raritätssystem für maximale Sammlerfreude
- KI-gesteuerte dynamische Kämpfe (keine statischen Gegner)
- Spielergetriebene Wirtschaft mit Gold- und Token-System
- Tiefes Skill- und Attunement-System für strategische Tiefe

---

## 👤 User Personas

### 1. Der Sammler ("The Collector")
**Motivation:** Seltene Gegenstände sammeln, Completionismus  
**Verhalten:** Spielt täglich für Loot, verfolgt Season-Rewards  
**Pain Points:** Zu niedrige Drop-Raten, unintuitive Upgrade-Pfade  

### 2. Der Competitor ("The Grinder")
**Motivation:** Leaderboards dominieren, PvP-Erfolge  
**Verhalten:** Optimiert Builds, studiert Meta, nimmt an Turnieren teil  
**Pain Points:** Unbalanced Matchmaking, Pay-to-Win-Verdacht  

### 3. Der Strategist ("The Tactician")
**Motivation:** Komplexe Systeme meistern, Theorie-Crafting  
**Verhalten:** Experimentiert mit Skill-Combos, analysiert Kämpfe  
**Pain Points:** Mangelnde Transparenz bei Formeln, unintuitive UI  

### 4. Der Investor ("The Whale")
**Motivation:** Seltene Assets besitzen, Wertsteigerung  
**Verhalten:** Kauft Bundles, handelt auf Secondary Market  
**Pain Points:** Unklarer Wert, fehlende Prestige-Objekte  

---

## 🎮 Core Gameplay Loops

### Loop 1: Kampf & Belohnung (Core Loop)
```
Kampf starten → Strategie wählen → Kämpfen → Loot erhalten → Items verbessern → Nächster Kampf
```
**Frequency:** 5-10 Minuten pro Zyklus  
**Reward:** Sofortige Befriedigung durch visuelle Loot-Drops  

### Loop 2: Progression & Mastery (Mid-Term)
```
XP sammeln → Level Up → Skill Points verteilen → Neue Fähigkeiten freischalten → Schwierigere Inhalte
```
**Frequency:** Tägliche Sessions über Wochen  
**Reward:** Sinnvolle Charakter-Entwicklung sichtbar  

### Loop 3: Sammlung & Prestige (Long-Term)
```
Seltene Items farmen → Sammlung vervollständigen → Prestige-Ränge erreichen → Season-Rewards
```
**Frequency:** Über Monate/Seasons  
**Reward:** Status, Anerkennung, exklusive Kosmetika  

### Loop 4: Wirtschaft & Handel (Meta)
```
Ressourcen sammeln → Im Shop verkaufen/traden → Gold/Gems akkumulieren → Premium-Inhalte kaufen
```
**Frequency:** Kontinuierlich im Hintergrund  
**Reward:** Ökonomische Macht, Seltene Assets  

---

## 🎲 Spielsysteme (User-Facing)

### 1. Kämpfer-System

**Charakter-Erstellung:**
- Name, Aussehen, Klasse wählbar
- 3 Attribute: Stärke (DMG), Vitalität (HP), Geschick (Speed/Crit)
- Max Level: 50

**Progression:**
- XP durch Kämpfe gewinnen
- Pro Level-Up: +5 Attribute, +1 Skill Point
- Respeccing möglich (mit steigenden Kosten)

**Visualisierung:**
- Charakter-Portrait mit Ausrüstung
- XP-Bar mit nächstem Level-Ziel
- Stat-Vergleich beim Equip-Wechsel

### 2. Raritäts-System (9 Tiers)

| Tier | Name | Drop-Rate | Multiplier | Farbe |
|------|------|-----------|------------|-------|
| 1 | Broken | 5% | 0.5x | Grau |
| 2 | Common | 30% | 1.0x | Weiß |
| 3 | Uncommon | 25% | 1.5x | Grün |
| 4 | Rare | 18% | 2.0x | Blau |
| 5 | Epic | 12% | 3.0x | Lila |
| 6 | Legendary | 6% | 5.0x | Orange |
| 7 | Mythic | 3% | 7.0x | Pink |
| 8 | Divine | 0.8% | 10.0x | Gold |
| 9 | Unique | 0.2% | Variabel | Rainbow |

**User Experience:**
- Farbcodierung sofort erkennbar
- Rarität beeinflusst Basis-Stats und Enhancement-Potenzial
- Einzigartige Items (Unique) haben spezielle Effekte
- "Rarity Upgrades" durch Crafting möglich

### 3. Waffen-System

**Waffentypen:**
- **Swords:** Balanced, hohe Crit-Chance
- **Staves:** Magie-basiert, Area Damage
- **Bows:** (Geplant) Hohe Reichweite
- **Axes:** (Geplant) Hoher Schaden, langsam

**Enhancement (+1 bis +10):**
- Jeder Level erhöht Stats um 10%
- Material-Kosten steigen exponentiell
- Risiko des Scheiterns ab +6 (ohne Schutzitems)
- Visual Glow-Effect bei hohen Levels

**User Journey:**
1. Weapon finden/droppen
2. Mit besserer Weapon vergleichen (Stat-Vergleich UI)
3. Alte Weapon salvagen (Materialien zurückgewinnen)
4. Neue Weapon equipen oder alte upgraden

### 4. Shop-System (MVP)

**Verfügbare Kategorien:**

#### Gold-Pakete
- Small: 100 Gold ($0.99)
- Medium: 550 Gold ($4.99) + 10% Bonus
- Large: 1.200 Gold ($9.99) + 20% Bonus
- Whale: 6.500 Gold ($49.99) + 30% Bonus

#### Rarity-Bundles
- **Starter Bundle:** Common Weapon + 200 Gold ($2.99)
- **Epic Hunter:** 5x Epic-Drop-Boost + 500 Gold ($9.99)
- **Legendary Crate:** Zufällige Legendary Weapon ($19.99)
- **Mythic Ascension:** Mythic Weapon Choice + 2.000 Gold ($49.99)

#### Season Pass (Geplant)
- 30-Tage Season mit Daily Rewards
- Exklusive Kosmetika
- XP-Boosts
- Premium-Währung: "Gems"

**Payment Flow:**
1. User wählt Produkt
2. Preis in Fiat ($) angezeigt
3. Mock-Zahlung (Phase 1) → Später: Crypto-Integration
4. Sofortige Auslieferung ins Inventar
5. Bestätigung: "Purchase Successful!"

### 5. Skill-System

**Skill-Bäume (3 Pfade):**
- **Offense:** Schaden, Crit, Attack Speed
- **Defense:** HP, Armor, Regeneration  
- **Utility:** Speed, Dodge, Special Effects

**Mechanics:**
- 15 Skills pro Baum (5 pro Tier)
- 3 Ultimates (1 pro Pfad, endgame)
- Skill Points durch Level-Up
- Loadout-System: 5 aktive Skills im Kampf

**User Experience:**
- Skill-Tree Visualisierung (tech-tree style)
- Tooltips mit genauen Zahlen
- Respec-Kosten transparent anzeigen
- Loadout-Vergleich: Damage vs. Survivability

### 6. Attunement-System

**6 Elemente:**
- Fire, Water, Earth, Air, Light, Dark

**Leveling 1-25:**
- XP durch Element-spezifische Kämpfe
- Jede Stufe: +2% Element-DMG, +1% Resistenz
- Level 10/20: Major Bonus (z.B. Fire Proc)
- Level 25: Ultimate Passive

**User Impact:**
- Visuelle Effekte (Fire Weapons brennen, etc.)
- Synergien mit Skill-Builds
- Strategische Entscheidung: Fokus vs. Balance

### 7. Daily Quests (Geplant)

**Quest-Typen:**
- Win 3 Matches
- Equip a Divine rarity item
- Win without taking damage
- Complete 5 League matches

**Rewards:**
- Gold scaling mit Streak
- Daily: 100 Gold
- 7-Day Streak: +500 Bonus + Rare Drop

**Retention-Mechanik:**
- Push-Notifications: "Your daily rewards are waiting!"
- FOMO: Quests resetzen täglich
- Visible streak counter

---

## 🏆 Progression & Rewards

### Level-Belohnungen
| Level | Belohnung |
|-------|-----------|
| 5 | First Rare Weapon |
| 10 | Skill Point Bundle (5) |
| 25 | Epic Weapon Choice |
| 50 | Legendary Crate |
| Every 10 | Title + Cosmetic |

### Achievement-System (Geplant)
- "First Blood" - Erster Kampf gewonnen
- "Hoarder" - 100 Items gesammelt
- "Perfectionist" - +10 Weapon crafted
- "Undefeated" - 20 Wins Streak

### Season-Rewards
- 30-Tage Seasons
- Battle Pass mit Free/Premium Track
- Ranked Leaderboards
- Exclusive Cosmetics (nur diese Season)

---

## 💰 Wirtschaftssystem

### Währungen

**1. Gold (Soft Currency)**
- Quellen: Kämpfe, Quests, Salvaging
- Verwendung: Shop-Käufe, Upgrades, Respecs
- Keine Cap

**2. Gems (Hard Currency)**
- Quellen: Real Money, Seltene Drops, Season Pass
- Verwendung: Premium Shop, Bundles, Convenience

**3. Token (Geplant - Web3)**
- Quellen: High-level Content, Trading
- Verwendung: Governance, Staking, Premium Items

### Drop-Ökonomie
- Jeder Kampf hat garantierte Drops (Gold + Material)
- Weapon-Drops: Skalieren mit Kampf-Schwierigkeit
- Raritäts-Chancen: Transparent anzeigen
- Bad Luck Protection: Nach X Kämpfen höhere Chance

### Inflations-Kontrolle
- Gold-Sinks: Upgrades, Respecs, Shop
- Item-Decay: (Optional) Items könnten "wear" erhalten
- Seasonal Resets: Soft-Reset für Leaderboards

---

## 🎨 UI/UX-Prinzipien

### Mobile-First Design
- Sticky Bottom Navigation
- Thumb-reachable Controls
- Swipe-Gestures für Inventar
- Bottom Sheets für Modals

### Visual Feedback
- Loot-Drops: Screen shake + Sound + Particle Effects
- Level-Up: Full-screen celebration
- Rarity: Glow effects entsprechend Tier
- Progress: Animated bars, checkmarks

### Information Architecture
- Dashboard: Wichtigste Stats auf einen Blick
- Inventory: Grid-View mit Filter/Sort
- Shop: Kategorien, Featured, Bundles
- Fighter: Stats, Equipment, Skills Tabs

---

## 📊 Erfolgs-Metriken (KPIs)

### Engagement
- DAU/MAU Ratio (Ziel: >30%)
- Average Session Length (Ziel: 15+ Min)
- Retention D1/D7/D30 (Ziel: 60%/30%/15%)

### Monetarisierung
- ARPDAU (Average Revenue Per Daily Active User)
- Conversion Rate (Free → Paying)
- Average Purchase Value
- LTV (Lifetime Value)

### Game Health
- Churn Rate (Ziel: <5% monthly)
- Support Tickets per User
- NPS Score (Ziel: >50)

---

## 🚀 Roadmap

### Phase 1: MVP (Jetzt)
- ✅ 9-Tier Rarity System
- ✅ Basic Combat
- 🔄 Shop MVP
- 🔄 Skill System

### Phase 2: Engagement (Q2)
- Daily Quests
- Season Pass
- Achievements
- Leaderboards

### Phase 3: Social (Q3)
- Guilds
- PvP Arena
- Trading
- Chat

### Phase 4: Web3 (Q4)
- Token Integration
- NFT Weapons
- Staking
- DAO Governance

---

## 🎭 Tone of Voice

**Brand Personality:**
- Epic, aber approachable
- Crypto-native, aber nicht exklusiv
- Strategisch-taktisch, nicht mindless grind
- Community-driven

**Sprache:**
- Englisch primär
- Fachbegriffe: "Rarity", "Attunement", "Enhancement"
- Kein "Pay-to-Win", sondern "Accelerate Your Progress"

---

## 📝 Glossar

**Attunement:** Elementare Bindung eines Kämpfers (Fire/Water/etc.)  
**Enhancement:** Verbesserung von +0 auf +10  
**Rarity:** Seltenheitsgrad (Broken → Unique)  
**Salvage:** Item zerstören für Materialien  
**Respec:** Skill-Points zurücksetzen  
**Loadout:** Aktive Skill-Auswahl für Kampf  

---

*Dokument gepflegt vom Game Master Agent*  
*Nächste Überprüfung: Bei jeder major Feature-Änderung*

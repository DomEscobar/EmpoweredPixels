package attunements

import "time"

// Element represents the 6 elemental attunements
type Element int

const (
	Fire Element = iota
	Water
	Earth
	Wind
	Light
	Dark
)

func (e Element) String() string {
	switch e {
	case Fire:
		return "Fire"
	case Water:
		return "Water"
	case Earth:
		return "Earth"
	case Wind:
		return "Wind"
	case Light:
		return "Light"
	case Dark:
		return "Dark"
	default:
		return "Unknown"
	}
}

// ElementFromString converts string to Element
func ElementFromString(s string) (Element, bool) {
	switch s {
	case "Fire":
		return Fire, true
	case "Water":
		return Water, true
	case "Earth":
		return Earth, true
	case "Wind":
		return Wind, true
	case "Light":
		return Light, true
	case "Dark":
		return Dark, true
	default:
		return Fire, false
	}
}

// GetIconURL returns the icon URL for an element
func (e Element) GetIconURL() string {
	baseURL := "https://vibemedia.space"
	switch e {
	case Fire:
		return baseURL + "/att_fire_001.png?prompt=fire%20element%20icon%20with%20flames%20and%20ember&style=pixel_game_asset&key=NOGON"
	case Water:
		return baseURL + "/att_water_001.png?prompt=water%20element%20icon%20with%20waves%20and%20droplets&style=pixel_game_asset&key=NOGON"
	case Earth:
		return baseURL + "/att_earth_001.png?prompt=earth%20element%20icon%20with%20stones%20and%20mountains&style=pixel_game_asset&key=NOGON"
	case Wind:
		return baseURL + "/att_wind_001.png?prompt=wind%20element%20icon%20with%20swirls%20and%20clouds&style=pixel_game_asset&key=NOGON"
	case Light:
		return baseURL + "/att_light_001.png?prompt=light%20element%20icon%20with%20rays%20and%20holy%20glow&style=pixel_game_asset&key=NOGON"
	case Dark:
		return baseURL + "/att_dark_001.png?prompt=dark%20element%20icon%20with%20shadows%20and%20void&style=pixel_game_asset&key=NOGON"
	default:
		return ""
	}
}

// GetColor returns the hex color for an element
func (e Element) GetColor() string {
	switch e {
	case Fire:
		return "#FF4500"
	case Water:
		return "#1E90FF"
	case Earth:
		return "#8B4513"
	case Wind:
		return "#32CD32"
	case Light:
		return "#FFD700"
	case Dark:
		return "#4B0082"
	default:
		return "#808080"
	}
}

// IsStrongAgainst returns true if this element is strong against target
// Element cycle: Fire > Wind > Earth > Water > Fire
// Light <> Dark (strong against each other)
func (e Element) IsStrongAgainst(target Element) bool {
	switch e {
	case Fire:
		return target == Wind
	case Wind:
		return target == Earth
	case Earth:
		return target == Water
	case Water:
		return target == Fire
	case Light:
		return target == Dark
	case Dark:
		return target == Light
	}
	return false
}

// IsWeakAgainst returns true if this element is weak against target
func (e Element) IsWeakAgainst(target Element) bool {
	switch e {
	case Fire:
		return target == Water
	case Water:
		return target == Earth
	case Earth:
		return target == Wind
	case Wind:
		return target == Fire
	case Light:
		return target == Dark
	case Dark:
		return target == Light
	}
	return false
}

// DamageModifier returns the damage modifier against a target element
// Strong: +25% damage, Weak: -25% damage, Neutral: 0% modifier
const (
	StrongModifier = 1.25
	WeakModifier   = 0.75
	NeutralModifier = 1.0
)

func (e Element) DamageModifierAgainst(target Element) float64 {
	if e.IsStrongAgainst(target) {
		return StrongModifier
	}
	if e.IsWeakAgainst(target) {
		return WeakModifier
	}
	return NeutralModifier
}

// PassiveBonus represents the passive bonuses at different level tiers
type PassiveBonus struct {
	DamageBonus   float64 // percentage (e.g., 0.10 = +10%)
	HealingBonus  float64
	ArmorBonus    float64
	SpeedBonus    float64
	CritBonus     float64
	LifestealBonus float64
}

// GetPassiveBonus returns passive bonuses based on attunement level (1-25 for MVP)
func GetPassiveBonus(element Element, level int) PassiveBonus {
	// Level 1-10: Tier 1, Level 11-25: Tier 2
	tier := 1
	if level >= 11 {
		tier = 2
	}

	switch element {
	case Fire:
		if tier == 1 {
			return PassiveBonus{DamageBonus: 0.05}
		}
		return PassiveBonus{DamageBonus: 0.10}
	case Water:
		if tier == 1 {
			return PassiveBonus{HealingBonus: 0.05}
		}
		return PassiveBonus{HealingBonus: 0.10}
	case Earth:
		if tier == 1 {
			return PassiveBonus{ArmorBonus: 0.08}
		}
		return PassiveBonus{ArmorBonus: 0.15}
	case Wind:
		if tier == 1 {
			return PassiveBonus{SpeedBonus: 0.08}
		}
		return PassiveBonus{SpeedBonus: 0.15}
	case Light:
		if tier == 1 {
			return PassiveBonus{CritBonus: 0.05}
		}
		return PassiveBonus{CritBonus: 0.10}
	case Dark:
		if tier == 1 {
			return PassiveBonus{LifestealBonus: 0.08}
		}
		return PassiveBonus{LifestealBonus: 0.15}
	}
	return PassiveBonus{}
}

// ActiveAbility represents an active ability for an attunement
type ActiveAbility struct {
	Name        string
	Description string
	Cooldown    int // seconds
	ManaCost    int
	EffectType  string // damage, heal, buff, crowd_control
	EffectValue int
}

// GetActiveAbility returns the active ability for an element
func GetActiveAbility(element Element) ActiveAbility {
	switch element {
	case Fire:
		return ActiveAbility{
			Name:        "Fireball",
			Description: "Launch a ball of fire dealing 50 damage and applying Burn",
			Cooldown:    10,
			ManaCost:    25,
			EffectType:  "damage",
			EffectValue: 50,
		}
	case Water:
		return ActiveAbility{
			Name:        "Healing Wave",
			Description: "Restore 15% of maximum health to target",
			Cooldown:    15,
			ManaCost:    30,
			EffectType:  "heal",
			EffectValue: 15, // percentage
		}
	case Earth:
		return ActiveAbility{
			Name:        "Stone Shield",
			Description: "Block the next 2 incoming attacks",
			Cooldown:    20,
			ManaCost:    35,
			EffectType:  "buff",
			EffectValue: 2, // blocks
		}
	case Wind:
		return ActiveAbility{
			Name:        "Gust",
			Description: "Knockback enemy and deal 30 damage",
			Cooldown:    12,
			ManaCost:    20,
			EffectType:  "crowd_control",
			EffectValue: 30,
		}
	case Light:
		return ActiveAbility{
			Name:        "Divine Shield",
			Description: "Become invulnerable for 3 seconds",
			Cooldown:    45,
			ManaCost:    50,
			EffectType:  "buff",
			EffectValue: 3, // seconds
		}
	case Dark:
		return ActiveAbility{
			Name:        "Shadow Bolt",
			Description: "Deal 60 damage and apply Fear for 1 second",
			Cooldown:    18,
			ManaCost:    40,
			EffectType:  "crowd_control",
			EffectValue: 60,
		}
	}
	return ActiveAbility{}
}

// FighterAttunement represents a fighter's attunement state
type FighterAttunement struct {
	FighterID       string
	Element         Element
	Level           int // 1-25 for MVP
	XP              int
	SelectedAt      time.Time
	LastChangedAt   time.Time
	ChangeCount     int
}

// AllElements returns all 6 elements
func AllElements() []Element {
	return []Element{Fire, Water, Earth, Wind, Light, Dark}
}

// XP Requirements for MVP (Level 1-25)
// Level 1-10: 100 XP per level
// Level 11-25: 200 XP per level
// Total to 25: 4,000 XP
func XPNeededForLevel(level int) int {
	if level < 1 || level > 25 {
		return 0
	}
	if level <= 10 {
		return 100
	}
	return 200
}

// TotalXPForLevel returns cumulative XP needed to reach a level
func TotalXPForLevel(level int) int {
	if level < 1 || level > 25 {
		return 0
	}
	total := 0
	for i := 1; i < level; i++ {
		total += XPNeededForLevel(i)
	}
	return total
}

// MaxXPForLevel returns the XP cap for a given level
func MaxXPForLevel(level int) int {
	return TotalXPForLevel(level + 1)
}

// AddXP adds XP and returns new level, remaining XP, and whether leveled up
func (fa *FighterAttunement) AddXP(amount int) (newLevel int, leveledUp bool) {
	fa.XP += amount
	leveledUp = false

	for fa.Level < 25 {
		xpNeeded := XPNeededForLevel(fa.Level)
		if fa.XP >= xpNeeded {
			fa.XP -= xpNeeded
			fa.Level++
			leveledUp = true
		} else {
			break
		}
	}

	return fa.Level, leveledUp
}

// ProgressPercent returns the progress percentage to next level
func (fa *FighterAttunement) ProgressPercent() float64 {
	if fa.Level >= 25 {
		return 100.0
	}
	xpNeeded := XPNeededForLevel(fa.Level)
	return float64(fa.XP) / float64(xpNeeded) * 100.0
}

// CanChange returns true if attunement can be changed
// Free for first change, then 24h cooldown or instant for gold
func (fa *FighterAttunement) CanChange() (bool, string) {
	if fa.ChangeCount == 0 {
		return true, "First change is free"
	}

	cooldownEnd := fa.LastChangedAt.Add(24 * time.Hour)
	if time.Now().After(cooldownEnd) {
		return true, "Cooldown completed"
	}

	remaining := time.Until(cooldownEnd)
	return false, remaining.String()
}

// ChangeCostGold returns gold cost for instant change
const InstantChangeCost = 300

func ChangeCostGold(changeCount int) int {
	if changeCount == 0 {
		return 0
	}
	return InstantChangeCost
}

// ChangeAttunement changes the attunement element
func (fa *FighterAttunement) ChangeAttunement(newElement Element, useGold bool) error {
	canChange, reason := fa.CanChange()
	if !canChange && !useGold {
		return &AttunementError{Message: "Cannot change: " + reason}
	}

	cost := ChangeCostGold(fa.ChangeCount)
	if useGold && cost == 0 {
		// First change is free, no gold needed
		useGold = false
	}

	fa.Element = newElement
	fa.Level = 1
	fa.XP = 0
	fa.LastChangedAt = time.Now()
	fa.ChangeCount++

	return nil
}

// AttunementError represents an attunement-related error
type AttunementError struct {
	Message string
}

func (e *AttunementError) Error() string {
	return e.Message
}

// GetDescription returns the lore description for an element
func GetDescription(element Element) string {
	switch element {
	case Fire:
		return "Fire attunement grants destructive power and burning damage over time."
	case Water:
		return "Water attunement enhances healing abilities and mana regeneration."
	case Earth:
		return "Earth attunement provides unmatched durability and defensive prowess."
	case Wind:
		return "Wind attunement offers superior speed and evasion capabilities."
	case Light:
		return "Light attunement blesses fighters with critical strikes and divine protection."
	case Dark:
		return "Dark attunement drains life from enemies and instills fear."
	}
	return ""
}

// GetStrengths returns elements this element is strong against
func (e Element) GetStrengths() []Element {
	var strengths []Element
	for _, el := range AllElements() {
		if e.IsStrongAgainst(el) {
			strengths = append(strengths, el)
		}
	}
	return strengths
}

// GetWeaknesses returns elements this element is weak against
func (e Element) GetWeaknesses() []Element {
	var weaknesses []Element
	for _, el := range AllElements() {
		if e.IsWeakAgainst(el) {
			weaknesses = append(weaknesses, el)
		}
	}
	return weaknesses
}

// AttunementInfo provides complete info about an element
type AttunementInfo struct {
	Element         Element
	Name            string
	Description     string
	Color           string
	IconURL         string
	Strengths       []Element
	Weaknesses      []Element
	PassiveBonus    PassiveBonus
	ActiveAbility   ActiveAbility
}

// GetAttunementInfo returns complete info for an element at a given level
func GetAttunementInfo(element Element, level int) AttunementInfo {
	return AttunementInfo{
		Element:       element,
		Name:          element.String(),
		Description:   GetDescription(element),
		Color:         element.GetColor(),
		IconURL:       element.GetIconURL(),
		Strengths:     element.GetStrengths(),
		Weaknesses:    element.GetWeaknesses(),
		PassiveBonus:  GetPassiveBonus(element, level),
		ActiveAbility: GetActiveAbility(element),
	}
}

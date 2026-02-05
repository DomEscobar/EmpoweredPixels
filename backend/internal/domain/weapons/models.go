package weapons

import "time"

// WeaponType represents the type of weapon
type WeaponType int

const (
	Sword WeaponType = iota
	Bow
	Staff
	Dagger
	Axe
)

func (t WeaponType) String() string {
	switch t {
	case Sword:
		return "Sword"
	case Bow:
		return "Bow"
	case Staff:
		return "Staff"
	case Dagger:
		return "Dagger"
	case Axe:
		return "Axe"
	default:
		return "Unknown"
	}
}

// Rarity represents weapon rarity level
type Rarity int

const (
	Common Rarity = iota
	Rare
	Epic
	Legendary
	Mythic
)

func (r Rarity) String() string {
	switch r {
	case Common:
		return "Common"
	case Rare:
		return "Rare"
	case Epic:
		return "Epic"
	case Legendary:
		return "Legendary"
	case Mythic:
		return "Mythic"
	default:
		return "Unknown"
	}
}

func (r Rarity) Color() string {
	switch r {
	case Common:
		return "#9CA3AF" // gray-400
	case Rare:
		return "#3B82F6" // blue-500
	case Epic:
		return "#A855F7" // purple-500
	case Legendary:
		return "#F59E0B" // amber-500
	case Mythic:
		return "#EF4444" // red-500
	default:
		return "#6B7280"
	}
}

// Weapon represents a weapon definition (static data)
type Weapon struct {
	ID            string
	Name          string
	Type          WeaponType
	Rarity        Rarity
	BaseDamage    int
	AttackSpeed   float64
	CritChance    int // percentage
	Durability    int
	IconURL       string
	Description   string
}

// UserWeapon represents a weapon owned by a user (instance data)
type UserWeapon struct {
	ID            string
	UserID        int64
	WeaponID      string
	Enhancement   int // +1 to +10
	Durability    int
	IsEquipped    bool
	FighterID     *string
	Created       time.Time
}

// WeaponStats represents calculated weapon stats after enhancement
type WeaponStats struct {
	Damage        int
	AttackSpeed   float64
	CritChance    int
}

// EnhancementResult represents the result of an enhancement attempt
type EnhancementResult struct {
	Success       bool
	NewLevel      int
	PreviousLevel int
	Destroyed     bool // true if weapon breaks (drops to +0)
}

// InventorySlot represents a slot in the weapon inventory
type InventorySlot struct {
	SlotNumber    int
	UserWeapon    *UserWeapon
	Weapon        *Weapon
}

// Constants
const (
	MaxInventorySlots = 50
	MaxEnhancement    = 10
)

// EnhancementFailureChance returns the failure chance for a given enhancement level
func EnhancementFailureChance(level int) float64 {
	switch {
	case level <= 3:
		return 0.0
	case level <= 6:
		return 0.15
	case level <= 9:
		return 0.35
	case level == 10:
		return 0.50
	default:
		return 1.0
	}
}

// EnhancementCost returns the gold cost for enhancement based on rarity
func EnhancementCost(rarity Rarity, level int) int {
	baseCost := 0
	switch rarity {
	case Common:
		baseCost = 100
	case Rare:
		baseCost = 250
	case Epic:
		baseCost = 500
	case Legendary:
		baseCost = 1000
	case Mythic:
		baseCost = 2000
	}
	// Cost increases with level
	return baseCost * (1 + level)
}

// CalculateStats calculates weapon stats after enhancement
func CalculateStats(weapon *Weapon, enhancement int) WeaponStats {
	multiplier := 1.0 + (float64(enhancement) * 0.1)
	return WeaponStats{
		Damage:      int(float64(weapon.BaseDamage) * multiplier),
		AttackSpeed: weapon.AttackSpeed,
		CritChance:  weapon.CritChance + (enhancement * 2), // +2% crit per enhancement
	}
}

// CanEnhance checks if a weapon can be enhanced further
func CanEnhance(enhancement int) bool {
	return enhancement < MaxEnhancement
}

// ApplyEnhancement applies enhancement result
func ApplyEnhancement(weapon *UserWeapon, success bool) EnhancementResult {
	prevLevel := weapon.Enhancement
	result := EnhancementResult{
		Success:       success,
		PreviousLevel: prevLevel,
	}

	if success {
		weapon.Enhancement++
		result.NewLevel = weapon.Enhancement
	} else {
		// Failure above +5 drops to +0
		if prevLevel > 5 {
			weapon.Enhancement = 0
			result.Destroyed = true
			result.NewLevel = 0
		} else {
			result.NewLevel = prevLevel
		}
	}

	return result
}
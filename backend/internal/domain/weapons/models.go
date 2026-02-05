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
	Broken Rarity = iota    // Tier 1 - Gray - 0.5x power
	Common                  // Tier 2 - White - 1.0x power
	Uncommon                // Tier 3 - Green - 1.3x power
	Rare                    // Tier 4 - Blue - 1.6x power
	Epic                    // Tier 5 - Purple - 2.0x power
	Legendary               // Tier 6 - Gold - 2.5x power
	Mythic                  // Tier 7 - Red - 3.5x power
	Divine                  // Tier 8 - Rainbow - 5.0x power
	Unique                  // Tier 9 - Animated - 10.0x power (Event-only)
)

func (r Rarity) String() string {
	switch r {
	case Broken:
		return "Broken"
	case Common:
		return "Common"
	case Uncommon:
		return "Uncommon"
	case Rare:
		return "Rare"
	case Epic:
		return "Epic"
	case Legendary:
		return "Legendary"
		case Mythic:
		return "Mythic"
	case Divine:
		return "Divine"
	case Unique:
		return "Unique"
	default:
		return "Unknown"
	}
}

func (r Rarity) Color() string {
	switch r {
	case Broken:
		return "#6B7280" // gray-500
	case Common:
		return "#9CA3AF" // gray-400
	case Uncommon:
		return "#22C55E" // green-500
	case Rare:
		return "#3B82F6" // blue-500
	case Epic:
		return "#A855F7" // purple-500
	case Legendary:
		return "#F59E0B" // amber-500
	case Mythic:
		return "#EF4444" // red-500
	case Divine:
		return "#E879F9" // fuchsia-400 (rainbow-like)
	case Unique:
		return "#FACC15" // yellow-400 (animated/gold)
	default:
		return "#6B7280"
	}
}

// PowerMultiplier returns the power scaling factor for this rarity
func (r Rarity) PowerMultiplier() float64 {
	switch r {
	case Broken:
		return 0.5
	case Common:
		return 1.0
	case Uncommon:
		return 1.3
	case Rare:
		return 1.6
	case Epic:
		return 2.0
	case Legendary:
		return 2.5
	case Mythic:
		return 3.5
	case Divine:
		return 5.0
	case Unique:
		return 10.0
	default:
		return 1.0
	}
}

// DropRate returns the drop rate percentage for this rarity
func (r Rarity) DropRate() float64 {
	switch r {
	case Broken:
		return 15.0
	case Common:
		return 50.0
	case Uncommon:
		return 20.0
	case Rare:
		return 10.0
	case Epic:
		return 4.0
	case Legendary:
		return 0.9
	case Mythic:
		return 0.1
	case Divine:
		return 0.01
	case Unique:
		return 0.0 // Event-only
	default:
		return 0.0
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
	case Broken:
		baseCost = 25
	case Common:
		baseCost = 100
	case Uncommon:
		baseCost = 175
	case Rare:
		baseCost = 250
	case Epic:
		baseCost = 500
	case Legendary:
		baseCost = 1000
	case Mythic:
		baseCost = 2000
	case Divine:
		baseCost = 5000
	case Unique:
		baseCost = 10000
	}
	// Cost increases with level
	return baseCost * (1 + level)
}

// CalculateStats calculates weapon stats after enhancement
func CalculateStats(weapon *Weapon, enhancement int) WeaponStats {
	// Apply rarity power multiplier first
	rarityMultiplier := weapon.Rarity.PowerMultiplier()
	
	// Apply enhancement multiplier
	enhancementMultiplier := 1.0 + (float64(enhancement) * 0.1)
	
	// Combined multiplier
	totalMultiplier := rarityMultiplier * enhancementMultiplier
	
	return WeaponStats{
		Damage:      int(float64(weapon.BaseDamage) * totalMultiplier),
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

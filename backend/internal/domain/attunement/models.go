package attunement

// Element represents one of the 6 attunement elements
type Element string

const (
	Fire  Element = "fire"
	Water Element = "water"
	Earth Element = "earth"
	Air   Element = "air"
	Light Element = "light"
	Dark  Element = "dark"
)

// AllElements returns all 6 elements
var AllElements = []Element{Fire, Water, Earth, Air, Light, Dark}

// Attunement represents a player's progress in one element
type Attunement struct {
	Element   Element `json:"element" db:"element"`
	Level     int     `json:"level" db:"level"`
	CurrentXP int     `json:"current_xp" db:"current_xp"`
	TotalXP   int     `json:"total_xp" db:"total_xp"` // Lifetime XP earned
}

// PlayerAttunements holds all 6 attunements for a player
type PlayerAttunements struct {
	UserID      int          `json:"user_id"`
	Attunements []Attunement `json:"attunements"`
	TotalLevel  int          `json:"total_level"` // Sum of all levels
}

// AttunementBonus represents the bonuses granted at each level
type AttunementBonus struct {
	Level     int     `json:"level"`
	Power     float64 `json:"power"`     // % damage increase
	Defense   float64 `json:"defense"`   // % damage reduction
	Speed     float64 `json:"speed"`     // % speed/turn order
	Precision float64 `json:"precision"` // % crit chance
}

// XP required for each level (1-25)
// Formula: base + (level * multiplier)
var XPRequirements = []int{
	0,      // Level 1 (starting)
	100,    // Level 2
	250,    // Level 3
	450,    // Level 4
	700,    // Level 5
	1000,   // Level 6
	1350,   // Level 7
	1750,   // Level 8
	2200,   // Level 9
	2700,   // Level 10
	3250,   // Level 11
	3850,   // Level 12
	4500,   // Level 13
	5200,   // Level 14
	5950,   // Level 15
	6750,   // Level 16
	7600,   // Level 17
	8500,   // Level 18
	9450,   // Level 19
	10450,  // Level 20
	11500,  // Level 21
	12600,  // Level 22
	13750,  // Level 23
	14950,  // Level 24
	16200,  // Level 25 (max)
}

// GetBonus returns the bonus values for a given level
func GetBonus(element Element, level int) AttunementBonus {
	if level < 1 {
		level = 1
	}
	if level > 25 {
		level = 25
	}

	// Each element has different bonus focus
	switch element {
	case Fire:
		return AttunementBonus{
			Level:     level,
			Power:     float64(level) * 0.5,  // 0.5% per level (12.5% max)
			Defense:   float64(level) * 0.1,
			Speed:     float64(level) * 0.2,
			Precision: float64(level) * 0.2,
		}
	case Water:
		return AttunementBonus{
			Level:     level,
			Power:     float64(level) * 0.2,
			Defense:   float64(level) * 0.4,  // 0.4% per level (10% max)
			Speed:     float64(level) * 0.2,
			Precision: float64(level) * 0.1,
		}
	case Earth:
		return AttunementBonus{
			Level:     level,
			Power:     float64(level) * 0.2,
			Defense:   float64(level) * 0.5,  // 0.5% per level (12.5% max)
			Speed:     float64(level) * 0.1,
			Precision: float64(level) * 0.1,
		}
	case Air:
		return AttunementBonus{
			Level:     level,
			Power:     float64(level) * 0.2,
			Defense:   float64(level) * 0.1,
			Speed:     float64(level) * 0.5,  // 0.5% per level (12.5% max)
			Precision: float64(level) * 0.1,
		}
	case Light:
		return AttunementBonus{
			Level:     level,
			Power:     float64(level) * 0.3,
			Defense:   float64(level) * 0.2,
			Speed:     float64(level) * 0.2,
			Precision: float64(level) * 0.3,  // 0.3% per level (7.5% max)
		}
	case Dark:
		return AttunementBonus{
			Level:     level,
			Power:     float64(level) * 0.4,  // 0.4% per level (10% max)
			Defense:   float64(level) * 0.1,
			Speed:     float64(level) * 0.2,
			Precision: float64(level) * 0.2,
		}
	}
	return AttunementBonus{Level: level}
}

// GetXPRequired returns XP needed to reach next level
func GetXPRequired(currentLevel int) int {
	if currentLevel < 1 || currentLevel >= 25 {
		return 0
	}
	return XPRequirements[currentLevel]
}

// XPSource defines how XP can be earned
type XPSource struct {
	Source string `json:"source"`
	XP     int    `json:"xp"`
}

// XPReward represents XP gained from an activity
const (
	XPWinMatch     = 50   // Win a match
	XPCompleteTask = 30   // Complete daily task
	XPUseElement   = 10   // Use element in combat (per use)
	XPEquipItem    = 100  // Equip attuned item
)

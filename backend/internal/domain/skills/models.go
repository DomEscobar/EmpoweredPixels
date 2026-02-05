package skills

// SkillBranch represents the three skill branches
type SkillBranch int

const (
	Offense SkillBranch = iota
	Defense
	Utility
)

func (b SkillBranch) String() string {
	switch b {
	case Offense:
		return "Offense"
	case Defense:
		return "Defense"
	case Utility:
		return "Utility"
	default:
		return "Unknown"
	}
}

// SkillType indicates if skill is active or passive
type SkillType int

const (
	Passive SkillType = iota
	Active
	Ultimate
)

func (t SkillType) String() string {
	switch t {
	case Passive:
		return "Passive"
	case Active:
		return "Active"
	case Ultimate:
		return "Ultimate"
	default:
		return "Unknown"
	}
}

// Skill represents a skill definition (static data)
type Skill struct {
	ID          string
	Name        string
	Branch      SkillBranch
	Tier        int      // 1-5
	Type        SkillType
	Description string
	ManaCost    int
	Cooldown    int // seconds
	EffectValue int // damage, heal, etc.
	Duration    int // seconds, 0 for instant
	IconURL     string
}

// FighterSkills represents a fighter's skill allocations and loadout
type FighterSkills struct {
	FighterID      string
	SkillPoints    int // available points
	AllocatedPoints map[string]int // skillID -> points spent (1-3 ranks)
	Loadout        []string // skill IDs for active skills (max 2)
	UltimateCharge int      // 0-100
}

// SkillAllocation represents points allocated to a skill
type SkillAllocation struct {
	SkillID string
	Points  int // 1-3 (ranks)
}

// LoadoutSlot represents an active skill slot
type LoadoutSlot struct {
	SlotNumber int
	SkillID    string
}

// CooldownStatus represents current cooldown state
type CooldownStatus struct {
	SkillID   string
	Remaining int // seconds remaining
}

// Constants
const (
	MaxSkillPoints    = 100
	MaxActiveSkills   = 2
	UltimateThreshold = 100
	PointsPerLevel    = 1
)

// ResetCost returns gold cost to reset skill tree (scales with level)
func ResetCost(fighterLevel int) int {
	return fighterLevel * 50
}

// CanAllocate checks if skill points can be allocated
func CanAllocate(currentAllocated map[string]int, skillID string, fighterLevel int) bool {
	// Check if skill exists and prerequisites met
	skill, found := GetSkillByID(skillID)
	if !found {
		return false
	}

	// Count total allocated points
	totalAllocated := 0
	for _, points := range currentAllocated {
		totalAllocated += points
	}

	// Check available points
	availablePoints := fighterLevel * PointsPerLevel
	if totalAllocated >= availablePoints {
		return false
	}

	// Check current allocation for this skill (max 3 ranks)
	currentPoints := currentAllocated[skillID]
	if currentPoints >= 3 {
		return false
	}

	// Check tier prerequisites: Tier 2+ requires 2 points in previous tier of same branch
	if skill.Tier > 1 {
		pointsInPrevTier := 0
		for id, points := range currentAllocated {
			if s, ok := GetSkillByID(id); ok && s.Branch == skill.Branch && s.Tier == skill.Tier-1 {
				pointsInPrevTier += points
			}
		}
		if pointsInPrevTier < 2 {
			return false
		}
	}

	return true
}

// CalculateSkillEffect calculates the effect value based on skill rank
func CalculateSkillEffect(skill *Skill, rank int) int {
	if rank <= 0 || rank > 3 {
		return 0
	}
	// Effect increases by 25% per rank
	multiplier := 1.0 + (float64(rank-1) * 0.25)
	return int(float64(skill.EffectValue) * multiplier)
}

// CanSetLoadout checks if loadout is valid
func CanSetLoadout(loadout []string, allocated map[string]int) bool {
	if len(loadout) > MaxActiveSkills {
		return false
	}

	for _, skillID := range loadout {
		// Must have points allocated
		if allocated[skillID] == 0 {
			return false
		}

		// Must be active skill type
		if skill, found := GetSkillByID(skillID); found {
			if skill.Type != Active {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// AddUltimateCharge adds charge to ultimate
func AddUltimateCharge(current, damageDealt, damageTaken int) int {
	// 1% per 50 damage dealt + 1% per 25 damage taken
	charge := (damageDealt / 50) + (damageTaken / 25)
	newCharge := current + charge
	if newCharge > UltimateThreshold {
		return UltimateThreshold
	}
	return newCharge
}

// SkillProgress represents skill tree progress
type SkillProgress struct {
	Branch           SkillBranch
	PointsAllocated  int
	SkillsUnlocked   int
	MaxTierUnlocked  int
}

// GetSkillProgress calculates progress for each branch
func GetSkillProgress(allocated map[string]int) []SkillProgress {
	progress := []SkillProgress{
		{Branch: Offense},
		{Branch: Defense},
		{Branch: Utility},
	}

	for id, points := range allocated {
		if skill, found := GetSkillByID(id); found {
			idx := int(skill.Branch)
			progress[idx].PointsAllocated += points
			progress[idx].SkillsUnlocked++
			if skill.Tier > progress[idx].MaxTierUnlocked {
				progress[idx].MaxTierUnlocked = skill.Tier
			}
		}
	}

	return progress
}
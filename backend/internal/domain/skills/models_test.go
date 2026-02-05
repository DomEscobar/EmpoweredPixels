package skills_test

import (
	"fmt"
	"testing"

	"empoweredpixels/internal/domain/skills"
)

func TestCanAllocate(t *testing.T) {
	tests := []struct {
		name            string
		allocated       map[string]int
		skillID         string
		fighterLevel    int
		expectedResult  bool
	}{
		{
			name:           "Can allocate tier 1 skill",
			allocated:      map[string]int{},
			skillID:        "skl_power_strike",
			fighterLevel:   10,
			expectedResult: true,
		},
		{
			name:           "Cannot allocate without skill points",
			allocated:      map[string]int{"skl_power_strike": 1},
			skillID:        "skl_whirlwind",
			fighterLevel:   1,
			expectedResult: false,
		},
		{
			name:           "Cannot exceed max rank",
			allocated:      map[string]int{"skl_power_strike": 3},
			skillID:        "skl_power_strike",
			fighterLevel:   10,
			expectedResult: false,
		},
		{
			name:           "Cannot allocate tier 2 without tier 1 prerequisites",
			allocated:      map[string]int{"skl_power_strike": 1},
			skillID:        "skl_whirlwind",
			fighterLevel:   10,
			expectedResult: false,
		},
		{
			name:           "Can allocate tier 2 with tier 1 prerequisites met",
			allocated:      map[string]int{"skl_power_strike": 2, "skl_bleed": 1},
			skillID:        "skl_whirlwind",
			fighterLevel:   10,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skills.CanAllocate(tt.allocated, tt.skillID, tt.fighterLevel)
			if result != tt.expectedResult {
				t.Errorf("CanAllocate() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestCalculateSkillEffect(t *testing.T) {
	skill := &skills.Skill{
		ID:          "skl_power_strike",
		Name:        "Power Strike",
		EffectValue: 40,
	}

	tests := []struct {
		rank           int
		expectedValue  int
	}{
		{rank: 1, expectedValue: 40},
		{rank: 2, expectedValue: 50},  // 40 * 1.25
		{rank: 3, expectedValue: 60},  // 40 * 1.5
		{rank: 0, expectedValue: 0},
		{rank: 4, expectedValue: 0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Rank %d", tt.rank), func(t *testing.T) {
			result := skills.CalculateSkillEffect(skill, tt.rank)
			if result != tt.expectedValue {
				t.Errorf("CalculateSkillEffect() = %v, want %v", result, tt.expectedValue)
			}
		})
	}
}

func TestCanSetLoadout(t *testing.T) {
	tests := []struct {
		name           string
		loadout        []string
		allocated      map[string]int
		expectedResult bool
	}{
		{
			name:           "Valid loadout with 2 active skills",
			loadout:        []string{"skl_power_strike", "skl_block"},
			allocated:      map[string]int{"skl_power_strike": 1, "skl_block": 1},
			expectedResult: true,
		},
		{
			name:           "Loadout exceeds max size",
			loadout:        []string{"skl_power_strike", "skl_block", "skl_heal", "skl_haste"},
			allocated:      map[string]int{"skl_power_strike": 1, "skl_block": 1, "skl_heal": 1, "skl_haste": 1},
			expectedResult: false,
		},
		{
			name:           "Cannot equip unallocated skill",
			loadout:        []string{"skl_power_strike"},
			allocated:      map[string]int{},
			expectedResult: false,
		},
		{
			name:           "Cannot equip passive skill",
			loadout:        []string{"skl_bleed"},
			allocated:      map[string]int{"skl_bleed": 1},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skills.CanSetLoadout(tt.loadout, tt.allocated)
			if result != tt.expectedResult {
				t.Errorf("CanSetLoadout() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestAddUltimateCharge(t *testing.T) {
	tests := []struct {
		name           string
		current        int
		damageDealt    int
		damageTaken    int
		expectedCharge int
	}{
		{
			name:           "Charge from damage dealt",
			current:        0,
			damageDealt:    100,
			damageTaken:    0,
			expectedCharge: 2, // 100/50 = 2
		},
		{
			name:           "Charge from damage taken",
			current:        0,
			damageDealt:    0,
			damageTaken:    100,
			expectedCharge: 4, // 100/25 = 4
		},
		{
			name:           "Charge from both",
			current:        50,
			damageDealt:    100,
			damageTaken:    50,
			expectedCharge: 54, // 50 + 2 + 2
		},
		{
			name:           "Cap at 100",
			current:        95,
			damageDealt:    500,
			damageTaken:    0,
			expectedCharge: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skills.AddUltimateCharge(tt.current, tt.damageDealt, tt.damageTaken)
			if result != tt.expectedCharge {
				t.Errorf("AddUltimateCharge() = %v, want %v", result, tt.expectedCharge)
			}
		})
	}
}

func TestResetCost(t *testing.T) {
	tests := []struct {
		fighterLevel int
		expectedCost int
	}{
		{fighterLevel: 10, expectedCost: 500},
		{fighterLevel: 50, expectedCost: 2500},
		{fighterLevel: 100, expectedCost: 5000},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level %d", tt.fighterLevel), func(t *testing.T) {
			result := skills.ResetCost(tt.fighterLevel)
			if result != tt.expectedCost {
				t.Errorf("ResetCost() = %v, want %v", result, tt.expectedCost)
			}
		})
	}
}

func TestGetSkillByID(t *testing.T) {
	tests := []struct {
		skillID        string
		expectedFound  bool
	}{
		{skillID: "skl_power_strike", expectedFound: true},
		{skillID: "skl_block", expectedFound: true},
		{skillID: "skl_haste", expectedFound: true},
		{skillID: "ult_meteor_strike", expectedFound: true},
		{skillID: "nonexistent", expectedFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.skillID, func(t *testing.T) {
			skill, found := skills.GetSkillByID(tt.skillID)
			if found != tt.expectedFound {
				t.Errorf("GetSkillByID() found = %v, want %v", found, tt.expectedFound)
			}
			if found && skill == nil {
				t.Error("GetSkillByID() returned nil skill when found")
			}
		})
	}
}

func TestGetSkillsByBranch(t *testing.T) {
	branches := []skills.SkillBranch{skills.Offense, skills.Defense, skills.Utility}

	for _, branch := range branches {
		t.Run(branch.String(), func(t *testing.T) {
			skills := skills.GetSkillsByBranch(branch)
			if len(skills) == 0 {
				t.Error("GetSkillsByBranch() returned no skills")
			}
			for _, skill := range skills {
				if skill.Branch != branch {
					t.Errorf("GetSkillsByBranch() returned skill from wrong branch: %v", skill.Branch)
				}
			}
		})
	}
}
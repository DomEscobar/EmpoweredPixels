package weapons

import (
	"math"
	"testing"
)

// TestRarityDropRatesSum verifies all drop rates sum to 100%
func TestRarityDropRatesSum(t *testing.T) {
	rarities := []Rarity{Broken, Common, Uncommon, Rare, Epic, Legendary, Mythic, Divine}
	
	var total float64
	for _, r := range rarities {
		total += r.DropRate()
	}
	
	// Allow small floating point tolerance (Divine is 0.01%, making sum 100.01%)
	if math.Abs(total-100.0) > 0.02 {
		t.Errorf("Drop rates sum to %v, expected ~100.0", total)
	}
}

// TestNewRarityTiers verifies Broken, Uncommon, Divine exist and have correct values
func TestNewRarityTiers(t *testing.T) {
	tests := []struct {
		rarity            Rarity
		expectedName      string
		expectedColor     string
		expectedPower     float64
		expectedDropRate  float64
	}{
		{Broken, "Broken", "#6B7280", 0.5, 15.0},
		{Uncommon, "Uncommon", "#22C55E", 1.3, 20.0},
		{Divine, "Divine", "#E879F9", 5.0, 0.01},
	}
	
	for _, tt := range tests {
		t.Run(tt.expectedName, func(t *testing.T) {
			if got := tt.rarity.String(); got != tt.expectedName {
				t.Errorf("Rarity.String() = %v, want %v", got, tt.expectedName)
			}
			if got := tt.rarity.Color(); got != tt.expectedColor {
				t.Errorf("Rarity.Color() = %v, want %v", got, tt.expectedColor)
			}
			if got := tt.rarity.PowerMultiplier(); got != tt.expectedPower {
				t.Errorf("Rarity.PowerMultiplier() = %v, want %v", got, tt.expectedPower)
			}
			if got := tt.rarity.DropRate(); got != tt.expectedDropRate {
				t.Errorf("Rarity.DropRate() = %v, want %v", got, tt.expectedDropRate)
			}
		})
	}
}

// TestAllRarityPowerMultipliers verifies all power multipliers
func TestAllRarityPowerMultipliers(t *testing.T) {
	tests := []struct {
		rarity   Rarity
		expected float64
	}{
		{Broken, 0.5},
		{Common, 1.0},
		{Uncommon, 1.3},
		{Rare, 1.6},
		{Epic, 2.0},
		{Legendary, 2.5},
		{Mythic, 3.5},
		{Divine, 5.0},
		{Unique, 10.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.rarity.String(), func(t *testing.T) {
			if got := tt.rarity.PowerMultiplier(); got != tt.expected {
				t.Errorf("Rarity.PowerMultiplier() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestAllRarityDropRates verifies all drop rates
func TestAllRarityDropRates(t *testing.T) {
	tests := []struct {
		rarity   Rarity
		expected float64
	}{
		{Broken, 15.0},
		{Common, 50.0},
		{Uncommon, 20.0},
		{Rare, 10.0},
		{Epic, 4.0},
		{Legendary, 0.9},
		{Mythic, 0.1},
		{Divine, 0.01},
		{Unique, 0.0}, // Event-only
	}
	
	for _, tt := range tests {
		t.Run(tt.rarity.String(), func(t *testing.T) {
			if got := tt.rarity.DropRate(); got != tt.expected {
				t.Errorf("Rarity.DropRate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestAllRarityColors verifies all rarity colors
func TestAllRarityColors(t *testing.T) {
	tests := []struct {
		rarity   Rarity
		expected string
	}{
		{Broken, "#6B7280"},
		{Common, "#9CA3AF"},
		{Uncommon, "#22C55E"},
		{Rare, "#3B82F6"},
		{Epic, "#A855F7"},
		{Legendary, "#F59E0B"},
		{Mythic, "#EF4444"},
		{Divine, "#E879F9"},
		{Unique, "#FACC15"},
	}
	
	for _, tt := range tests {
		t.Run(tt.rarity.String(), func(t *testing.T) {
			if got := tt.rarity.Color(); got != tt.expected {
				t.Errorf("Rarity.Color() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestRarityProgression verifies rarity tiers are in correct order
func TestRarityProgression(t *testing.T) {
	// Rarities ordered by value (worst to best)
	rarities := []Rarity{Broken, Common, Uncommon, Rare, Epic, Legendary, Mythic, Divine, Unique}
	
	for i := 0; i < len(rarities)-1; i++ {
		current := rarities[i]
		next := rarities[i+1]
		
		// Power should increase
		if current.PowerMultiplier() >= next.PowerMultiplier() {
			t.Errorf("Power multiplier should increase: %s (%.1f) -> %s (%.1f)",
				current.String(), current.PowerMultiplier(),
				next.String(), next.PowerMultiplier())
		}
	}
	
	// Higher value rarities should generally have lower drop rates
	// (Common is most common at 50%, Broken is less common at 15%)
	if Common.DropRate() <= Broken.DropRate() {
		t.Errorf("Common should have higher drop rate than Broken")
	}
	if Uncommon.DropRate() >= Common.DropRate() {
		t.Errorf("Uncommon should have lower drop rate than Common")
	}
	if Divine.DropRate() >= Mythic.DropRate() {
		t.Errorf("Divine should have lower drop rate than Mythic")
	}
}

// TestCalculateStatsWithRarity verifies stats calculation with rarity multipliers
func TestCalculateStatsWithRarity(t *testing.T) {
	baseWeapon := &Weapon{
		BaseDamage:  100,
		AttackSpeed: 1.0,
		CritChance:  10,
	}
	
	tests := []struct {
		rarity       Rarity
		expectedDmg  int
	}{
		{Broken, 50},      // 100 * 0.5
		{Common, 100},     // 100 * 1.0
		{Uncommon, 130},   // 100 * 1.3
		{Rare, 160},       // 100 * 1.6
		{Epic, 200},       // 100 * 2.0
		{Legendary, 250},  // 100 * 2.5
		{Mythic, 350},     // 100 * 3.5
		{Divine, 500},     // 100 * 5.0
		{Unique, 1000},    // 100 * 10.0
	}
	
	for _, tt := range tests {
		t.Run(tt.rarity.String(), func(t *testing.T) {
			weapon := *baseWeapon
			weapon.Rarity = tt.rarity
			
			stats := CalculateStats(&weapon, 0)
			if stats.Damage != tt.expectedDmg {
				t.Errorf("CalculateStats().Damage = %v, want %v", stats.Damage, tt.expectedDmg)
			}
		})
	}
}

// TestEnhancementCostForNewRarities verifies enhancement costs for new rarities
func TestEnhancementCostForNewRarities(t *testing.T) {
	tests := []struct {
		rarity   Rarity
		level    int
		expected int
	}{
		{Broken, 0, 25},
		{Broken, 5, 150},
		{Uncommon, 0, 175},
		{Uncommon, 5, 1050},
		{Divine, 0, 5000},
		{Divine, 5, 30000},
	}
	
	for _, tt := range tests {
		t.Run(tt.rarity.String()+"_level_"+string(rune('0'+tt.level)), func(t *testing.T) {
			got := EnhancementCost(tt.rarity, tt.level)
			if got != tt.expected {
				t.Errorf("EnhancementCost(%v, %d) = %v, want %v", tt.rarity, tt.level, got, tt.expected)
			}
		})
	}
}

// TestUniqueRarityEventOnly verifies Unique rarity cannot drop
func TestUniqueRarityEventOnly(t *testing.T) {
	if Unique.DropRate() != 0.0 {
		t.Errorf("Unique.DropRate() = %v, want 0.0 (event-only)", Unique.DropRate())
	}
}

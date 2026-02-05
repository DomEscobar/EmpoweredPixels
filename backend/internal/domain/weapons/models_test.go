package weapons

import (
	"testing"
)

func TestWeaponType_String(t *testing.T) {
	tests := []struct {
		weaponType WeaponType
		want       string
	}{
		{Sword, "Sword"},
		{Bow, "Bow"},
		{Staff, "Staff"},
		{Dagger, "Dagger"},
		{Axe, "Axe"},
		{WeaponType(99), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.weaponType.String(); got != tt.want {
				t.Errorf("WeaponType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRarity_String(t *testing.T) {
	tests := []struct {
		rarity Rarity
		want   string
	}{
		{Common, "Common"},
		{Rare, "Rare"},
		{Epic, "Epic"},
		{Legendary, "Legendary"},
		{Mythic, "Mythic"},
		{Rarity(99), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.rarity.String(); got != tt.want {
				t.Errorf("Rarity.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRarity_Color(t *testing.T) {
	tests := []struct {
		rarity Rarity
		want   string
	}{
		{Common, "#9CA3AF"},
		{Rare, "#3B82F6"},
		{Epic, "#A855F7"},
		{Legendary, "#F59E0B"},
		{Mythic, "#EF4444"},
	}
	for _, tt := range tests {
		t.Run(tt.rarity.String(), func(t *testing.T) {
			if got := tt.rarity.Color(); got != tt.want {
				t.Errorf("Rarity.Color() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnhancementFailureChance(t *testing.T) {
	tests := []struct {
		level int
		want  float64
	}{
		{0, 0.0},
		{1, 0.0},
		{2, 0.0},
		{3, 0.0},
		{4, 0.15},
		{5, 0.15},
		{6, 0.15},
		{7, 0.35},
		{8, 0.35},
		{9, 0.35},
		{10, 0.50},
		{11, 1.0},
	}
	for _, tt := range tests {
		t.Run(string(rune('0'+tt.level)), func(t *testing.T) {
			if got := EnhancementFailureChance(tt.level); got != tt.want {
				t.Errorf("EnhancementFailureChance(%d) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestEnhancementCost(t *testing.T) {
	tests := []struct {
		rarity Rarity
		level  int
		min    int
		max    int
	}{
		{Common, 0, 100, 100},
		{Common, 5, 600, 600},
		{Rare, 0, 250, 250},
		{Epic, 0, 500, 500},
		{Legendary, 0, 1000, 1000},
		{Mythic, 0, 2000, 2000},
	}
	for _, tt := range tests {
		t.Run(tt.rarity.String(), func(t *testing.T) {
			got := EnhancementCost(tt.rarity, tt.level)
			if got < tt.min || got > tt.max {
				t.Errorf("EnhancementCost(%v, %d) = %v, want between %v and %v", tt.rarity, tt.level, got, tt.min, tt.max)
			}
		})
	}
}

func TestCalculateStats(t *testing.T) {
	weapon := &Weapon{
		BaseDamage:  100,
		AttackSpeed: 1.0,
		CritChance:  10,
	}

	tests := []struct {
		enhancement int
		wantDamage  int
		wantCrit    int
	}{
		{0, 100, 10},
		{1, 110, 12},
		{5, 150, 20},
		{10, 200, 30},
	}
	for _, tt := range tests {
		t.Run(string(rune('0'+tt.enhancement)), func(t *testing.T) {
			got := CalculateStats(weapon, tt.enhancement)
			if got.Damage != tt.wantDamage {
				t.Errorf("CalculateStats().Damage = %v, want %v", got.Damage, tt.wantDamage)
			}
			if got.CritChance != tt.wantCrit {
				t.Errorf("CalculateStats().CritChance = %v, want %v", got.CritChance, tt.wantCrit)
			}
			if got.AttackSpeed != weapon.AttackSpeed {
				t.Errorf("CalculateStats().AttackSpeed = %v, want %v", got.AttackSpeed, weapon.AttackSpeed)
			}
		})
	}
}

func TestCanEnhance(t *testing.T) {
	tests := []struct {
		enhancement int
		want        bool
	}{
		{0, true},
		{5, true},
		{9, true},
		{10, false},
		{11, false},
	}
	for _, tt := range tests {
		t.Run(string(rune('0'+tt.enhancement)), func(t *testing.T) {
			if got := CanEnhance(tt.enhancement); got != tt.want {
				t.Errorf("CanEnhance(%d) = %v, want %v", tt.enhancement, got, tt.want)
			}
		})
	}
}

func TestApplyEnhancement(t *testing.T) {
	tests := []struct {
		name           string
		initialLevel   int
		success        bool
		wantNewLevel   int
		wantDestroyed  bool
	}{
		{"success_low_level", 2, true, 3, false},
		{"failure_low_level", 2, false, 2, false},
		{"success_high_level", 7, true, 8, false},
		{"failure_high_level", 7, false, 0, true},
		{"success_max", 9, true, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &UserWeapon{Enhancement: tt.initialLevel}
			result := ApplyEnhancement(weapon, tt.success)

			if result.NewLevel != tt.wantNewLevel {
				t.Errorf("ApplyEnhancement().NewLevel = %v, want %v", result.NewLevel, tt.wantNewLevel)
			}
			if result.Destroyed != tt.wantDestroyed {
				t.Errorf("ApplyEnhancement().Destroyed = %v, want %v", result.Destroyed, tt.wantDestroyed)
			}
			if result.PreviousLevel != tt.initialLevel {
				t.Errorf("ApplyEnhancement().PreviousLevel = %v, want %v", result.PreviousLevel, tt.initialLevel)
			}
			if result.Success != tt.success {
				t.Errorf("ApplyEnhancement().Success = %v, want %v", result.Success, tt.success)
			}
		})
	}
}

func TestGetWeaponByID(t *testing.T) {
	// Test with a known weapon from the database
	weapon, found := GetWeaponByID("wpn_sword_excalibur_005")
	if !found {
		t.Error("Expected to find Excalibur")
	}
	if weapon.Name != "Excalibur" {
		t.Errorf("Expected weapon name 'Excalibur', got '%s'", weapon.Name)
	}
	if weapon.Rarity != Legendary {
		t.Errorf("Expected Legendary rarity, got %v", weapon.Rarity)
	}

	// Test with unknown ID
	_, found = GetWeaponByID("unknown_weapon")
	if found {
		t.Error("Expected not to find unknown weapon")
	}
}

func TestGetWeaponsByType(t *testing.T) {
	swords := GetWeaponsByType(Sword)
	if len(swords) != 5 {
		t.Errorf("Expected 5 swords, got %d", len(swords))
	}

	bows := GetWeaponsByType(Bow)
	if len(bows) != 4 {
		t.Errorf("Expected 4 bows, got %d", len(bows))
	}

	daggers := GetWeaponsByType(Dagger)
	if len(daggers) != 3 {
		t.Errorf("Expected 3 daggers, got %d", len(daggers))
	}
}

func TestGetWeaponsByRarity(t *testing.T) {
	common := GetWeaponsByRarity(Common)
	if len(common) < 4 {
		t.Errorf("Expected at least 4 common weapons, got %d", len(common))
	}

	mythic := GetWeaponsByRarity(Mythic)
	if len(mythic) != 1 {
		t.Errorf("Expected 1 mythic weapon, got %d", len(mythic))
	}
}

func TestWeaponDatabase_Count(t *testing.T) {
	if len(WeaponDatabase) != 20 {
		t.Errorf("Expected 20 weapons in database, got %d", len(WeaponDatabase))
	}
}
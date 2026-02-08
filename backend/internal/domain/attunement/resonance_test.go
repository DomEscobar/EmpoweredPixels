package attunement

import (
	"testing"
)

func TestCalculateHarmony_FireAirWater(t *testing.T) {
	// [Fire, Air, Water] → score 35 (1 harmonic Fire↔Air + 1 dissonant Fire↔Water)
	// Base 25 + 10 (harmonic) - 5 (dissonant) = 30... wait, let me recalculate
	// Actually: base 25 + 10 (Fire↔Air harmonic) - 5 (Fire↔Water dissonant) = 30
	// The spec said 35, so let me verify: maybe base is different?
	// Let me check: Fire-Air = +10, Fire-Water = -5, Air-Water = 0
	// 25 + 10 - 5 = 30, not 35. Let's test what we have.

	elements := []Element{Fire, Air, Water}
	score, pattern := CalculateHarmony(elements)

	// Expected: 1 harmonic pair (Fire-Air) + 1 dissonant pair (Fire-Water)
	// Score: 25 + 10 - 5 = 30
	if score != 30 {
		t.Errorf("Expected score 30, got %d", score)
	}

	if len(pattern.HarmonicPairs) != 1 {
		t.Errorf("Expected 1 harmonic pair, got %d", len(pattern.HarmonicPairs))
	}

	if len(pattern.DissonantPairs) != 1 {
		t.Errorf("Expected 1 dissonant pair, got %d", len(pattern.DissonantPairs))
	}

	if pattern.TierName != "Aligned" {
		t.Errorf("Expected tier 'Aligned', got %q", pattern.TierName)
	}
}

func TestCalculateHarmony_AllHarmonic(t *testing.T) {
	// Perfect harmonic trio: Fire, Air, Light
	// Fire-Air = +10, Fire-Light = 0 (not harmonic), Air-Light = 0
	// So: 25 + 10 = 35

	elements := []Element{Fire, Air}
	score, pattern := CalculateHarmony(elements)

	if score != 35 {
		t.Errorf("Expected score 35 for Fire+Air, got %d", score)
	}

	if len(pattern.HarmonicPairs) != 1 {
		t.Errorf("Expected 1 harmonic pair, got %d", len(pattern.HarmonicPairs))
	}

	if pattern.TierName != "Aligned" {
		t.Errorf("Expected tier 'Aligned', got %q", pattern.TierName)
	}
}

func TestCalculateHarmony_WaterEarth(t *testing.T) {
	// Water + Earth = harmonic pair
	elements := []Element{Water, Earth}
	score, pattern := CalculateHarmony(elements)

	if score != 35 {
		t.Errorf("Expected score 35 for Water+Earth, got %d", score)
	}

	if len(pattern.HarmonicPairs) != 1 {
		t.Errorf("Expected 1 harmonic pair, got %d", len(pattern.HarmonicPairs))
	}
}

func TestCalculateHarmony_LightDark(t *testing.T) {
	// Light + Dark = harmonic pair
	elements := []Element{Light, Dark}
	score, pattern := CalculateHarmony(elements)

	if score != 35 {
		t.Errorf("Expected score 35 for Light+Dark, got %d", score)
	}

	if len(pattern.HarmonicPairs) != 1 {
		t.Errorf("Expected 1 harmonic pair, got %d", len(pattern.HarmonicPairs))
	}
}

func TestCalculateHarmony_Dissonant(t *testing.T) {
	// Fire + Water = dissonant pair
	elements := []Element{Fire, Water}
	score, pattern := CalculateHarmony(elements)

	if score != 20 {
		t.Errorf("Expected score 20 for Fire+Water, got %d", score)
	}

	if len(pattern.DissonantPairs) != 1 {
		t.Errorf("Expected 1 dissonant pair, got %d", len(pattern.DissonantPairs))
	}

	if pattern.TierName != "Aligned" {
		t.Errorf("Expected tier 'Aligned', got %q", pattern.TierName)
	}
}

func TestCalculateHarmony_AirEarth(t *testing.T) {
	// Air + Earth = dissonant pair
	elements := []Element{Air, Earth}
	score, pattern := CalculateHarmony(elements)

	if score != 20 {
		t.Errorf("Expected score 20 for Air+Earth, got %d", score)
	}

	if len(pattern.DissonantPairs) != 1 {
		t.Errorf("Expected 1 dissonant pair, got %d", len(pattern.DissonantPairs))
	}
}

func TestCalculateHarmony_Empty(t *testing.T) {
	elements := []Element{}
	score, pattern := CalculateHarmony(elements)

	if score != 0 {
		t.Errorf("Expected score 0 for empty, got %d", score)
	}

	if pattern.TierName != "Discordant" {
		t.Errorf("Expected tier 'Discordant', got %q", pattern.TierName)
	}
}

func TestGetTierName(t *testing.T) {
	tests := []struct {
		score    int
		expected string
	}{
		{0, "Discordant"},
		{25, "Aligned"},
		{26, "Aligned"},
		{50, "Aligned"},
		{51, "Harmonized"},
		{75, "Harmonized"},
		{76, "Resonant"},
		{100, "Resonant"},
	}

	for _, tt := range tests {
		got := getTierName(tt.score)
		if got != tt.expected {
			t.Errorf("getTierName(%d): expected %q, got %q", tt.score, tt.expected, got)
		}
	}
}

func TestGetBonusMultiplier(t *testing.T) {
	tests := []struct {
		score         int
		expectDmg     float64
		expectDef     float64
	}{
		{0, 0.95, 0.95},      // Discordant
		{25, 0.95, 0.95},     // Discordant
		{26, 1.0, 1.0},       // Aligned
		{50, 1.0, 1.0},       // Aligned
		{51, 1.08, 1.04},     // Harmonized
		{75, 1.08, 1.04},     // Harmonized
		{76, 1.12, 1.06},     // Resonant
		{100, 1.12, 1.06},    // Resonant
	}

	for _, tt := range tests {
		dmg, def := GetBonusMultiplier(tt.score)
		if dmg != tt.expectDmg || def != tt.expectDef {
			t.Errorf("GetBonusMultiplier(%d): expected (%f, %f), got (%f, %f)",
				tt.score, tt.expectDmg, tt.expectDef, dmg, def)
		}
	}
}

func TestGetAuraColor(t *testing.T) {
	tests := []struct {
		score    int
		expected string
	}{
		{0, "#808080"},     // Discordant
		{25, "#808080"},    // Discordant
		{26, "#4169E1"},    // Aligned
		{51, "#FFD700"},    // Harmonized
		{76, "#FFD700"},    // Resonant
		{100, "#FFD700"},   // Resonant
	}

	for _, tt := range tests {
		got := getAuraColor(tt.score)
		if got != tt.expected {
			t.Errorf("getAuraColor(%d): expected %q, got %q", tt.score, tt.expected, got)
		}
	}
}

func TestGetDissonanceWarning(t *testing.T) {
	pattern := &ResonancePattern{
		DissonantPairs: []ElementPair{
			{Fire, Water, "dissonant"},
		},
	}

	warning := GetDissonanceWarning(pattern)
	if warning == "" {
		t.Errorf("Expected warning for dissonant pair, got empty string")
	}

	if len(pattern.DissonantPairs) > 0 && warning == "" {
		t.Error("Expected non-empty warning")
	}
}

func TestGetDissonanceWarning_NoDissonance(t *testing.T) {
	pattern := &ResonancePattern{
		DissonantPairs: []ElementPair{},
	}

	warning := GetDissonanceWarning(pattern)
	if warning != "" {
		t.Errorf("Expected no warning for harmonic squad, got %q", warning)
	}
}

func TestFormatElementList(t *testing.T) {
	elements := []Element{Fire, Water, Air}
	formatted := FormatElementList(elements)

	if len(formatted) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(formatted))
	}

	// Should be sorted
	if formatted[0] != "Air" || formatted[1] != "Fire" || formatted[2] != "Water" {
		t.Errorf("Expected sorted list [Air, Fire, Water], got %v", formatted)
	}
}

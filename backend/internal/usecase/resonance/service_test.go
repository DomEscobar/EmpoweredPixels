package resonance

import (
	"context"
	"testing"
	"time"

	"empoweredpixels/internal/domain/attunement"
	"empoweredpixels/internal/domain/roster"
)

// Mock implementations for testing

type mockSquadRepo struct {
	squad *roster.Squad
}

func (m *mockSquadRepo) GetByID(ctx context.Context, squadID string) (*roster.Squad, error) {
	return m.squad, nil
}

func (m *mockSquadRepo) UpdateResonanceScore(ctx context.Context, squadID string, score int, pattern string, timestamp time.Time) error {
	return nil
}

type mockFighterRepo struct {
	fighters []roster.Fighter
}

func (m *mockFighterRepo) GetByID(ctx context.Context, fighterID string) (*roster.Fighter, error) {
	for i, f := range m.fighters {
		if f.ID == fighterID {
			return &m.fighters[i], nil
		}
	}
	return nil, nil
}

func (m *mockFighterRepo) GetMultiple(ctx context.Context, fighterIDs []string) ([]roster.Fighter, error) {
	return m.fighters, nil
}

type mockAttunementRepo struct {
	playerAttunements []attunement.Attunement
}

func (m *mockAttunementRepo) GetPlayerAttunements(ctx context.Context, userID int) ([]attunement.Attunement, error) {
	return m.playerAttunements, nil
}

func TestCalculateSquadResonance_FireAirCombo(t *testing.T) {
	squad := &roster.Squad{
		ID:     "squad-1",
		UserID: 1,
		Members: []roster.Member{
			{FighterID: "fighter-1"},
			{FighterID: "fighter-2"},
		},
	}

	fighters := []roster.Fighter{
		{ID: "fighter-1", Power: 100},
		{ID: "fighter-2", Power: 100},
	}

	playerAttunements := []attunement.Attunement{
		{Element: attunement.Fire, Level: 1},
		{Element: attunement.Air, Level: 1},
		{Element: attunement.Water, Level: 0},
		{Element: attunement.Earth, Level: 0},
		{Element: attunement.Light, Level: 0},
		{Element: attunement.Dark, Level: 0},
	}

	squadRepo := &mockSquadRepo{squad: squad}
	fighterRepo := &mockFighterRepo{fighters: fighters}
	attunementRepo := &mockAttunementRepo{playerAttunements: playerAttunements}

	service := NewResonanceService(squadRepo, fighterRepo, attunementRepo)
	state, err := service.CalculateSquadResonance(context.Background(), "squad-1")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if state.HarmonyScore != 35 {
		t.Errorf("Expected harmony score 35, got %d", state.HarmonyScore)
	}

	if state.TierName != "Aligned" {
		t.Errorf("Expected tier 'Aligned', got %q", state.TierName)
	}

	if len(state.HarmonicElements) != 2 {
		t.Errorf("Expected 2 harmonic elements, got %d", len(state.HarmonicElements))
	}
}

func TestApplySquadBonuses(t *testing.T) {
	fighters := []roster.Fighter{
		{ID: "fighter-1", Power: 100, Armor: 50, ConditionPower: 80},
		{ID: "fighter-2", Power: 100, Armor: 50, ConditionPower: 80},
	}

	state := &ResonanceState{
		HarmonyScore: 75,
		TierName:     "Harmonized",
	}

	squadRepo := &mockSquadRepo{}
	fighterRepo := &mockFighterRepo{}
	attunementRepo := &mockAttunementRepo{}

	service := NewResonanceService(squadRepo, fighterRepo, attunementRepo)
	modified := service.ApplySquadBonuses(fighters, state)

	// At score 75, multiplier is 1.08 for damage, 1.04 for defense
	expectedPower := int(float64(100) * 1.08)       // 108
	expectedArmor := int(float64(50) * 1.04)        // 52
	expectedCondition := int(float64(80) * 1.08)    // 86

	for _, f := range modified {
		if f.Power != expectedPower {
			t.Errorf("Expected power %d, got %d", expectedPower, f.Power)
		}
		if f.Armor != expectedArmor {
			t.Errorf("Expected armor %d, got %d", expectedArmor, f.Armor)
		}
		if f.ConditionPower != expectedCondition {
			t.Errorf("Expected condition power %d, got %d", expectedCondition, f.ConditionPower)
		}
	}
}

func TestGetBonusMultiplier_Resonant(t *testing.T) {
	squadRepo := &mockSquadRepo{}
	fighterRepo := &mockFighterRepo{}
	attunementRepo := &mockAttunementRepo{}

	service := NewResonanceService(squadRepo, fighterRepo, attunementRepo)
	dmg, def := service.GetBonusMultiplier(80)

	if dmg != 1.12 {
		t.Errorf("Expected damage multiplier 1.12, got %f", dmg)
	}
	if def != 1.06 {
		t.Errorf("Expected defense multiplier 1.06, got %f", def)
	}
}

func TestGetTierName(t *testing.T) {
	squadRepo := &mockSquadRepo{}
	fighterRepo := &mockFighterRepo{}
	attunementRepo := &mockAttunementRepo{}

	service := NewResonanceService(squadRepo, fighterRepo, attunementRepo)

	tests := []struct {
		score    int
		expected string
	}{
		{0, "Discordant"},
		{50, "Aligned"},
		{60, "Harmonized"},
		{80, "Resonant"},
	}

	for _, tt := range tests {
		got := service.GetTierName(tt.score)
		if got != tt.expected {
			t.Errorf("GetTierName(%d): expected %q, got %q", tt.score, tt.expected, got)
		}
	}
}

func TestCalculateSquadResonance_NoSquad(t *testing.T) {
	squadRepo := &mockSquadRepo{squad: nil}
	fighterRepo := &mockFighterRepo{}
	attunementRepo := &mockAttunementRepo{}

	service := NewResonanceService(squadRepo, fighterRepo, attunementRepo)
	state, err := service.CalculateSquadResonance(context.Background(), "nonexistent")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if state.HarmonyScore != 0 {
		t.Errorf("Expected harmony score 0 for nonexistent squad, got %d", state.HarmonyScore)
	}
}

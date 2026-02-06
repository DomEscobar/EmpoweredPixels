package matches

import (
	"testing"

	"empoweredpixels/internal/domain/inventory"
	"empoweredpixels/internal/domain/roster"
)

func TestComboMomentumSystem(t *testing.T) {
	simulator := NewSimulator()

	// Create test fighters
	fighters := []roster.Fighter{
		{
			ID:        "fighter1",
			Name:      "Test Fighter 1",
			Level:     10,
			Power:     20,
			Accuracy:  15,
			Speed:     10,
			Armor:     10,
			Vitality:  10,
		},
		{
			ID:        "fighter2",
			Name:      "Test Fighter 2",
			Level:     10,
			Power:     18,
			Accuracy:  15,
			Speed:     12,
			Armor:     12,
			Vitality:  10,
		},
	}

	equipment := make(map[string][]inventory.Equipment)
	options := MatchOptions{}

	result, err := simulator.Run("test-match", fighters, equipment, options)
	if err != nil {
		t.Fatalf("simulator failed: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}

	// Verify match result structure
	if result.MatchID != "test-match" {
		t.Errorf("expected match ID 'test-match', got %s", result.MatchID)
	}

	// Check that we have rounds
	if len(result.RoundTicks) == 0 {
		t.Error("expected at least one round")
	}

	t.Logf("Match completed with %d rounds", len(result.RoundTicks))
}

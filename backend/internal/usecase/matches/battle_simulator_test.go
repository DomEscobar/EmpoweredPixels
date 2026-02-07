package matches

import (
	"testing"

	"empoweredpixels/internal/domain/roster"
	"github.com/google/uuid"
)

func TestBattleSimulator_Run(t *testing.T) {
	sim := NewBattleSimulator()
	matchID := uuid.NewString()

	fighters := []roster.Fighter{
		{
			ID:       uuid.NewString(),
			Name:     "Warrior",
			Level:    10,
			Power:    15,
			Armor:    10,
			Vitality: 12,
			Speed:    5,
		},
		{
			ID:       uuid.NewString(),
			Name:     "Ranger",
			Level:    10,
			Power:    8,
			Precision: 18,
			Agility:  15,
			Speed:    12,
			Vitality: 8,
		},
	}

	options := BattleOptions{
		MaxRounds: 50,
		MapSize:   20.0,
	}

	result, err := sim.Run(matchID, fighters, options)
	if err != nil {
		t.Fatalf("Failed to run simulation: %v", err)
	}

	if result.MatchID != matchID {
		t.Errorf("Expected match ID %s, got %s", matchID, result.MatchID)
	}

	if len(result.RoundTicks) == 0 {
		t.Error("Expected at least one round tick")
	}

	t.Logf("Simulated %d rounds", len(result.RoundTicks))
	for _, round := range result.RoundTicks {
		for _, tick := range round.Ticks {
			t.Logf("Round %d: %s", round.Round, tick.Type)
		}
	}
}

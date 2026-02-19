package matches

import (
	"testing"

	"empoweredpixels/internal/domain/roster"
)

func TestBattleSimulator_TurnPriority(t *testing.T) {
	sim := NewBattleSimulator()

	team1 := "Team1"
	team2 := "Team2"

	// High Agility vs Low Agility
	fighters := []roster.Fighter{
		{
			ID:      "fast",
			Name:    "Fast Rogue",
			Agility: 100,
			Speed:   10,
			Power:   10,
			TeamID:  &team1,
		},
		{
			ID:      "slow",
			Name:    "Slow Tank",
			Agility: 0,
			Speed:   10,
			Power:   10,
			TeamID:  &team2,
		},
	}

	result, err := sim.Run("test-match", fighters, BattleOptions{MaxRounds: 100})
	if err != nil {
		t.Fatalf("failed to run match: %v", err)
	}

	fastTurns := 0
	slowTurns := 0

	for _, rt := range result.RoundTicks {
		for _, tick := range rt.Ticks {
			// We track attacks or moves initiated by the fighters
			// The tick payload contains the fighter ID
			if tick.Type == "attack" || tick.Type == "move" {
				// Rough parsing of ID from JSON payload (simplest for test)
				// For move: {"fighterId":"..."}
				// For attack: {"attackerId":"..."}
				payload := string(tick.Payload)
				if contains(payload, "fast") {
					fastTurns++
				} else if contains(payload, "slow") {
					slowTurns++
				}
			}
		}
	}

	t.Logf("Fast turns: %d, Slow turns: %d", fastTurns, slowTurns)

	if fastTurns <= slowTurns {
		t.Errorf("Expected Fast Rogue to have more turns than Slow Tank, got %d vs %d", fastTurns, slowTurns)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	}()
}

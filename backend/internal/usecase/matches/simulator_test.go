package matches

import (
	"testing"

	"empoweredpixels/internal/domain/combat"
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

	// Verify momentum events are generated
	momentumEventFound := false
	flurryEventFound := false
	sunderEventFound := false

	for _, round := range result.RoundTicks {
		for _, tick := range round.Ticks {
			switch tick.Type {
			case "momentum":
				momentumEventFound = true
			case "flurry":
				flurryEventFound = true
			case "sunder":
				sunderEventFound = true
			}
		}
	}

	if !momentumEventFound {
		t.Error("expected momentum events to be generated")
	}

	t.Logf("Match completed with %d rounds", len(result.RoundTicks))
	t.Logf("Momentum events: %v, Flurry events: %v, Sunder events: %v", 
		momentumEventFound, flurryEventFound, sunderEventFound)
}

func TestMomentumOnHit(t *testing.T) {
	state := &combat.ComboMomentumState{
		FighterID:        "test1",
		Momentum:         0,
		ConsecutiveHits:  0,
		CurrentTargetID:  "",
		SunderStacks:     0,
		FlurryActive:     false,
		RoundsSinceHit:   0,
	}

	// Test first hit on target
	ticks, flurryActivated := updateMomentumOnHit(state, "target1")
	
	if state.Momentum != MomentumPerHit {
		t.Errorf("expected momentum %d, got %d", MomentumPerHit, state.Momentum)
	}
	
	if state.ConsecutiveHits != 1 {
		t.Errorf("expected consecutive hits 1, got %d", state.ConsecutiveHits)
	}
	
	if flurryActivated {
		t.Error("flurry should not activate on first hit")
	}

	// Test consecutive hits
	for i := 0; i < 5; i++ {
		updateMomentumOnHit(state, "target1")
	}

	if state.ConsecutiveHits != 5 {
		t.Errorf("expected consecutive hits 5, got %d", state.ConsecutiveHits)
	}

	if state.Momentum != 60 {
		t.Errorf("expected momentum 60, got %d", state.Momentum)
	}

	// Test flurry activation at >50 momentum
	if !state.FlurryActive {
		t.Error("flurry should be active at >50 momentum")
	}

	// Test target change resets combo
	updateMomentumOnHit(state, "target2")
	
	if state.ConsecutiveHits != 1 {
		t.Errorf("expected consecutive hits to reset to 1, got %d", state.ConsecutiveHits)
	}

	if len(ticks) == 0 {
		t.Error("expected momentum event ticks")
	}
}

func TestSunderDebuff(t *testing.T) {
	attacker := &combat.ComboMomentumState{
		FighterID:        "attacker1",
		ConsecutiveHits:  3, // 3 hits = 2 sunder stacks
		CurrentTargetID:  "target1",
	}
	
	target := &combat.ComboMomentumState{
		FighterID:       "target1",
		SunderStacks:    0,
	}
	
	targetEntity := &combat.Entity{
		ID: "target1",
		Stats: combat.Stats{
			Armor: 100,
		},
	}
	
	originalArmor := targetEntity.Stats.Armor
	
	ticks := applySunderDebuff(attacker, target, targetEntity)
	
	if target.SunderStacks != 2 {
		t.Errorf("expected 2 sunder stacks, got %d", target.SunderStacks)
	}
	
	// 2 stacks = 10% reduction = 10 armor
	expectedArmor := originalArmor - int(float64(originalArmor)*SunderArmorReduction*2)
	if targetEntity.Stats.Armor != expectedArmor {
		t.Errorf("expected armor %d, got %d", expectedArmor, targetEntity.Stats.Armor)
	}
	
	if len(ticks) == 0 {
		t.Error("expected sunder event ticks")
	}
}

func TestMomentumDecay(t *testing.T) {
	states := map[string]*combat.ComboMomentumState{
		"fighter1": {
			FighterID:       "fighter1",
			Momentum:        80,
			FlurryActive:    true,
			RoundsSinceHit:  0,
		},
	}
	
	// decayMomentum expects RoundsSinceHit to be already incremented
	// and only applies decay if RoundsSinceHit > 0 and Momentum > 0
	states["fighter1"].RoundsSinceHit = 1
	
	ticks := decayMomentum(states)
	
	// Should decay by MomentumDecay (5) each call
	expectedMomentum := 80 - MomentumDecay
	if states["fighter1"].Momentum != expectedMomentum {
		t.Errorf("expected momentum %d, got %d", expectedMomentum, states["fighter1"].Momentum)
	}
	
	// Flurry should still be active (75 > 50)
	if !states["fighter1"].FlurryActive {
		t.Error("flurry should still be active at 75 momentum")
	}
	
	// Apply more decay to drop below threshold
	states["fighter1"].RoundsSinceHit = 1
	decayMomentum(states)
	states["fighter1"].RoundsSinceHit = 1
	decayMomentum(states)
	states["fighter1"].RoundsSinceHit = 1
	decayMomentum(states)
	states["fighter1"].RoundsSinceHit = 1
	decayMomentum(states)
	
	// Now momentum should be 55, still above threshold
	if states["fighter1"].Momentum != 55 {
		t.Errorf("expected momentum 55, got %d", states["fighter1"].Momentum)
	}
	
	// One more decay to drop below threshold
	states["fighter1"].RoundsSinceHit = 1
	decayMomentum(states)
	
	if states["fighter1"].Momentum != 50 {
		t.Errorf("expected momentum 50, got %d", states["fighter1"].Momentum)
	}
	
	// Flurry should deactivate at 50 (threshold is >50)
	if states["fighter1"].FlurryActive {
		t.Error("flurry should deactivate when momentum drops to 50")
	}
	
	if len(ticks) == 0 {
		t.Error("expected decay event ticks")
	}
}
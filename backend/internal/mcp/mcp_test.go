package mcp

import (
	"context"
	"testing"
	"time"

	"empoweredpixels/internal/domain/matches"
)

func TestFairnessFilter(t *testing.T) {
	// Limit: 2 actions per 1 second
	filter := NewFairnessFilter(2, 1*time.Second)
	key := "test-ai-key-123456"

	// 1st action: OK
	if ok, _ := filter.Allow(key); !ok {
		t.Error("Should allow 1st action")
	}

	// 2nd action: OK
	if ok, _ := filter.Allow(key); !ok {
		t.Error("Should allow 2nd action")
	}

	// 3rd action: Should fail
	if ok, _ := filter.Allow(key); ok {
		t.Error("Should block 3rd action within same window")
	}

	// Wait for window to pass
	time.Sleep(1100 * time.Millisecond)

	// After wait: OK
	if ok, _ := filter.Allow(key); !ok {
		t.Error("Should allow action after window expired")
	}
}

func TestDefaultFairnessFilter(t *testing.T) {
	filter := DefaultFairnessFilter()
	key := "test-ai-key"

	// Should allow 100 requests per minute
	for i := 0; i < 100; i++ {
		if ok, _ := filter.Allow(key); !ok {
			t.Errorf("Should allow request %d", i+1)
		}
	}

	// 101st should fail
	if ok, _ := filter.Allow(key); ok {
		t.Error("Should block 101st request")
	}
}

func TestAuditLogger(t *testing.T) {
	logger, err := NewAuditLogger("")
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	ctx := context.Background()
	entry := AuditLogEntry{
		Action:   "test_action",
		AgentID:  "test-agent",
		UserID:   123,
		Endpoint: "/test",
		Success:  true,
	}

	// Should not panic
	logger.Log(ctx, entry)

	// Should return recent entries
	recent := logger.GetRecent(10)
	if len(recent) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(recent))
	}
}

func TestGameStateStructs(t *testing.T) {
	// Test that GameState struct works correctly
	state := &GameState{
		Matches: []MatchSummary{
			{
				ID:           "match-1",
				Status:       matches.MatchStatusLobby,
				FighterCount: 2,
				IsJoinable:   true,
			},
		},
		TotalLobbies:  1,
		ActiveMatches: 0,
	}

	if len(state.Matches) != 1 {
		t.Error("Expected 1 match")
	}

	if state.Matches[0].ID != "match-1" {
		t.Error("Expected match ID to be 'match-1'")
	}
}

func TestPlayerStatsStructs(t *testing.T) {
	stats := &PlayerStats{
		FighterID: "fighter-1",
		Name:      "Test Fighter",
		Level:     5,
		Power:     100,
		Attributes: map[string]int{
			"precision": 10,
			"armor":     20,
		},
	}

	if stats.Name != "Test Fighter" {
		t.Error("Expected fighter name to be 'Test Fighter'")
	}

	if stats.Attributes["precision"] != 10 {
		t.Error("Expected precision to be 10")
	}
}

func TestSubmitActionRequest(t *testing.T) {
	req := SubmitActionRequest{
		ActionType: "join_match",
		MatchID:    "match-1",
		FighterID:  "fighter-1",
		Parameters: map[string]interface{}{
			"team_id": "team-1",
		},
	}

	if req.ActionType != "join_match" {
		t.Error("Expected action type to be 'join_match'")
	}
}

func TestSubmitActionResult(t *testing.T) {
	result := SubmitActionResult{
		Success: true,
		Message: "Action completed",
		Data: map[string]string{
			"match_id": "match-1",
		},
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.Message != "Action completed" {
		t.Error("Expected message to be 'Action completed'")
	}
}

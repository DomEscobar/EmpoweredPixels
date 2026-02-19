package repositories

import (
	"context"
	"testing"
	"time"

	"empoweredpixels/internal/domain/rewards"
	"github.com/google/uuid"
)

func TestRewardRepository_SourceID(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewRewardRepository(pool)
	ctx := context.Background()

	userID := int64(999)
	matchID := "match-" + uuid.NewString()
	rewardID := uuid.NewString()

	reward := &rewards.Reward{
		ID:           rewardID,
		UserID:       userID,
		RewardPoolID: "test_pool",
		SourceID:     &matchID,
		Created:      time.Now().Truncate(time.Microsecond),
	}

	// Test Create
	if err := repo.Create(ctx, reward); err != nil {
		t.Fatalf("failed to create reward: %v", err)
	}

	// Test ListUnclaimed
	unclaimed, err := repo.ListUnclaimed(ctx, userID)
	if err != nil {
		t.Fatalf("failed to list unclaimed: %v", err)
	}

	found := false
	for _, r := range unclaimed {
		if r.ID == rewardID {
			found = true
			if r.SourceID == nil || *r.SourceID != matchID {
				t.Errorf("expected source ID %s, got %v", matchID, r.SourceID)
			}
		}
	}
	if !found {
		t.Error("created reward not found in unclaimed list")
	}

	// Test GetUnclaimed
	r, err := repo.GetUnclaimed(ctx, userID, rewardID, "test_pool")
	if err != nil {
		t.Fatalf("failed to get unclaimed: %v", err)
	}
	if r == nil {
		t.Fatal("expected reward, got nil")
	}
	if r.SourceID == nil || *r.SourceID != matchID {
		t.Errorf("expected source ID %s, got %v", matchID, r.SourceID)
	}
}

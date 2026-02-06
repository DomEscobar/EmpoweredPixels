package daily

import (
	"context"
	"fmt"
	"math/rand"

	"empoweredpixels/internal/domain/daily"
	"empoweredpixels/internal/infra/db/repositories"
)

// Service handles daily reward business logic
type Service struct {
	repo     repositories.DailyRewardRepository
	goldRepo repositories.PlayerGoldRepository
}

// NewService creates a new daily reward service
func NewService(
	repo repositories.DailyRewardRepository,
	goldRepo repositories.PlayerGoldRepository,
) *Service {
	return &Service{
		repo:     repo,
		goldRepo: goldRepo,
	}
}

// GetStatus returns user's daily reward status
func (s *Service) GetStatus(ctx context.Context, userID int) (*daily.UserDailyReward, error) {
	return s.repo.GetUserDailyReward(ctx, userID)
}

// Claim processes a reward claim
func (s *Service) Claim(ctx context.Context, userID int) (*daily.ClaimResult, error) {
	// Get current status
	status, err := s.repo.GetUserDailyReward(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	// Check if already claimed today
	if !status.CanClaim {
		return nil, fmt.Errorf("already claimed today")
	}

	// Check if streak is broken
	newStreak := status.Streak + 1
	if daily.HasStreakBroken(status.LastClaimed) {
		newStreak = 1 // Reset to day 1
	}

	// Get reward for this day
	reward := daily.GetRewardForDay(newStreak)

	// Process reward based on type
	var rewardValue int
	switch reward.Type {
	case "gold":
		// Add gold to user
		if err := s.goldRepo.AddGold(ctx, userID, reward.Value); err != nil {
			return nil, fmt.Errorf("failed to add gold: %w", err)
		}
		rewardValue = reward.Value

	case "item":
		// For now, just give gold equivalent (item system integration later)
		if err := s.goldRepo.AddGold(ctx, userID, reward.Value); err != nil {
			return nil, fmt.Errorf("failed to add gold: %w", err)
		}
		rewardValue = reward.Value

	case "boost":
		// Store boost in user session or separate table (simplified for now)
		rewardValue = reward.Value

	case "mystery":
		// Random reward between 100-1000 gold or random item
		mysteryGold := rand.Intn(900) + 100 // 100-1000
		if err := s.goldRepo.AddGold(ctx, userID, mysteryGold); err != nil {
			return nil, fmt.Errorf("failed to add mystery gold: %w", err)
		}
		rewardValue = mysteryGold
	}

	// Save claim
	if err := s.repo.ClaimReward(ctx, userID, newStreak); err != nil {
		return nil, fmt.Errorf("failed to save claim: %w", err)
	}

	// Get next reward preview
	nextDay := newStreak + 1
	nextReward := daily.GetRewardForDay(nextDay)

	return &daily.ClaimResult{
		Success:     true,
		Reward:      reward,
		RewardValue: rewardValue,
		NewStreak:   newStreak,
		Day:         newStreak,
		NextReward:  nextReward,
	}, nil
}

package leaderboard

import (
	"context"
	"fmt"
	"sort"

	"empoweredpixels/internal/domain/leaderboard"
	"empoweredpixels/internal/infra/db/repositories"
)

// Service handles leaderboard business logic
type Service struct {
	repo       repositories.LeaderboardRepository
	achieveRepo repositories.AchievementRepository
	userRepo   *repositories.UserRepository
	fighterRepo *repositories.FighterRepository
	goldRepo   repositories.PlayerGoldRepository
}

// NewService creates a new leaderboard service
func NewService(
	repo repositories.LeaderboardRepository,
	achieveRepo repositories.AchievementRepository,
	userRepo *repositories.UserRepository,
	fighterRepo *repositories.FighterRepository,
	goldRepo repositories.PlayerGoldRepository,
) *Service {
	return &Service{
		repo:       repo,
		achieveRepo: achieveRepo,
		userRepo:   userRepo,
		fighterRepo: fighterRepo,
		goldRepo:   goldRepo,
	}
}

// GetLeaderboard retrieves a leaderboard category
func (s *Service) GetLeaderboard(ctx context.Context, category string, userID int, limit int, offset int) (*leaderboard.ListResponse, error) {
	entries, err := s.repo.GetByCategory(ctx, category, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	totalCount, err := s.repo.GetTotalCount(ctx, category)
	if err != nil {
		totalCount = len(entries)
	}

	// Get user's rank
	userEntry, err := s.repo.GetUserRank(ctx, category, userID)
	if err != nil {
		userEntry = nil
	}

	return &leaderboard.ListResponse{
		Category:   category,
		TotalCount: totalCount,
		UserRank:   userEntry.Rank,
		UserEntry:  userEntry,
		Entries:    entries,
	}, nil
}

// GetNearbyRanks retrieves ranks near the user
func (s *Service) GetNearbyRanks(ctx context.Context, category string, userID int, rangeSize int) (*leaderboard.ListResponse, error) {
	entries, err := s.repo.GetNearbyRanks(ctx, category, userID, rangeSize)
	if err != nil {
		return nil, err
	}

	userEntry, _ := s.repo.GetUserRank(ctx, category, userID)

	return &leaderboard.ListResponse{
		Category:  category,
		UserRank:  userEntry.Rank,
		UserEntry: userEntry,
		Entries:   entries,
	}, nil
}

// RecalculateAll updates all leaderboard categories
func (s *Service) RecalculateAll(ctx context.Context) error {
	if err := s.recalculatePower(ctx); err != nil {
		return fmt.Errorf("power leaderboard failed: %w", err)
	}
	if err := s.recalculateWealth(ctx); err != nil {
		return fmt.Errorf("wealth leaderboard failed: %w", err)
	}
	if err := s.recalculateCombat(ctx); err != nil {
		return fmt.Errorf("combat leaderboard failed: %w", err)
	}
	return nil
}

// recalculatePower calculates power rankings based on total fighter power
func (s *Service) recalculatePower(ctx context.Context) error {
	// Get all users with their total fighter power
	users, err := s.userRepo.ListAll(ctx)
	if err != nil {
		return err
	}

	type userPower struct {
		userID int
		power  int64
	}

	var rankings []userPower
	for _, u := range users {
		fighters, err := s.fighterRepo.ListByUser(ctx, int64(u.ID))
		if err != nil {
			continue
		}
		var totalPower int64
		for _, f := range fighters {
			totalPower += int64(f.Power)
		}
		rankings = append(rankings, userPower{userID: int(u.ID), power: totalPower})
	}

	// Sort by power descending
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].power > rankings[j].power
	})

	// Update leaderboard entries
	for rank, r := range rankings {
		if r.power == 0 {
			continue
		}
		entry := &leaderboard.Entry{
			Category: leaderboard.CategoryPower,
			UserID:   int(r.userID),
			Rank:     rank + 1,
			Score:    r.power,
		}
		if err := s.repo.UpsertEntry(ctx, entry); err != nil {
			continue
		}
	}

	return nil
}

// recalculateWealth calculates wealth rankings based on gold balance
func (s *Service) recalculateWealth(ctx context.Context) error {
	balances, err := s.goldRepo.ListAllBalances(ctx)
	if err != nil {
		return err
	}

	// Sort by balance descending
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].Balance > balances[j].Balance
	})

	for rank, b := range balances {
		if b.Balance == 0 {
			continue
		}
		entry := &leaderboard.Entry{
			Category: leaderboard.CategoryWealth,
			UserID:   int(b.UserID),
			Rank:     rank + 1,
			Score:    int64(b.Balance),
		}
		if err := s.repo.UpsertEntry(ctx, entry); err != nil {
			continue
		}
	}

	return nil
}

// recalculateCombat calculates combat rankings based on matches won
func (s *Service) recalculateCombat(ctx context.Context) error {
	// This would ideally use the match_history table
	// For now, use fighter match statistics
	users, err := s.userRepo.ListAll(ctx)
	if err != nil {
		return err
	}

	type userCombat struct {
		userID int
		wins   int64
	}

	var rankings []userCombat
	for _, u := range users {
		fighters, err := s.fighterRepo.ListByUser(ctx, int64(u.ID))
		if err != nil {
			continue
		}
		var totalWins int64
		for _, f := range fighters {
			totalWins += int64(f.MatchesWon)
		}
		rankings = append(rankings, userCombat{userID: int(u.ID), wins: totalWins})
	}

	// Sort by wins descending
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].wins > rankings[j].wins
	})

	for rank, r := range rankings {
		if r.wins == 0 {
			continue
		}
		entry := &leaderboard.Entry{
			Category: leaderboard.CategoryCombat,
			UserID:   r.userID,
			Rank:     rank + 1,
			Score:    r.wins,
		}
		if err := s.repo.UpsertEntry(ctx, entry); err != nil {
			continue
		}
	}

	return nil
}

// GetAchievements retrieves all achievements
func (s *Service) GetAchievements(ctx context.Context) ([]leaderboard.Achievement, error) {
	return s.achieveRepo.ListAchievements(ctx)
}

// GetPlayerAchievements retrieves a player's achievements
func (s *Service) GetPlayerAchievements(ctx context.Context, userID int) ([]leaderboard.PlayerAchievement, error) {
	return s.achieveRepo.GetPlayerAchievements(ctx, userID)
}

// UpdateAchievementProgress updates a player's achievement progress
func (s *Service) UpdateAchievementProgress(ctx context.Context, userID int, achievementKey string, progress int) error {
	return s.achieveRepo.UpdateAchievementProgress(ctx, userID, achievementKey, progress)
}

// ClaimAchievementReward claims an achievement reward
func (s *Service) ClaimAchievementReward(ctx context.Context, userID int, achievementID string) error {
	return s.achieveRepo.ClaimAchievementReward(ctx, userID, achievementID)
}

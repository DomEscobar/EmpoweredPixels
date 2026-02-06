package repositories

import (
	"context"
	"fmt"
	"time"

	"empoweredpixels/internal/domain/daily"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DailyRewardRepository defines daily reward operations
type DailyRewardRepository interface {
	GetUserDailyReward(ctx context.Context, userID int) (*daily.UserDailyReward, error)
	ClaimReward(ctx context.Context, userID int, streak int) error
	ResetStreak(ctx context.Context, userID int) error
}

// DailyRewardPostgres implements DailyRewardRepository
type DailyRewardPostgres struct {
	db *pgxpool.Pool
}

// NewDailyRewardRepository creates a new daily reward repository
func NewDailyRewardRepository(db *pgxpool.Pool) DailyRewardRepository {
	return &DailyRewardPostgres{db: db}
}

// GetUserDailyReward retrieves a user's daily reward status
func (r *DailyRewardPostgres) GetUserDailyReward(ctx context.Context, userID int) (*daily.UserDailyReward, error) {
	query := `
		SELECT user_id, streak, last_claimed, total_claimed, updated
		FROM daily_rewards
		WHERE user_id = $1
	`

	var dr daily.UserDailyReward
	var lastClaimed *time.Time
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&dr.UserID, &dr.Streak, &lastClaimed, &dr.TotalClaimed, &dr.Updated,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Return empty record if not exists
			return &daily.UserDailyReward{
				UserID:   userID,
				Streak:   0,
				CanClaim: true,
			}, nil
		}
		return nil, fmt.Errorf("failed to get daily reward: %w", err)
	}

	dr.LastClaimed = lastClaimed

	// Check if streak is broken
	if daily.HasStreakBroken(lastClaimed) {
		dr.Streak = 0
	}

	// Check if can claim today
	dr.CanClaim = daily.CanClaimToday(lastClaimed)

	// Calculate next reward
	nextDay := dr.Streak + 1
	if dr.CanClaim {
		nextDay = dr.Streak + 1
	} else {
		nextDay = dr.Streak // Show current day's reward if already claimed
	}
	if nextDay < 1 {
		nextDay = 1
	}
	dr.NextReward = daily.GetRewardForDay(nextDay)

	// Calculate time until reset
	if !dr.CanClaim {
		dr.TimeUntilReset = formatDuration(daily.TimeUntilMidnight())
	}

	return &dr, nil
}

// ClaimReward records a reward claim
func (r *DailyRewardPostgres) ClaimReward(ctx context.Context, userID int, streak int) error {
	query := `
		INSERT INTO daily_rewards (user_id, streak, last_claimed, total_claimed, updated)
		VALUES ($1, $2, CURRENT_DATE, 1, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			streak = EXCLUDED.streak,
			last_claimed = CURRENT_DATE,
			total_claimed = daily_rewards.total_claimed + 1,
			updated = NOW()
	`

	_, err := r.db.Exec(ctx, query, userID, streak)
	if err != nil {
		return fmt.Errorf("failed to claim reward: %w", err)
	}

	return nil
}

// ResetStreak resets user's streak to 0
func (r *DailyRewardPostgres) ResetStreak(ctx context.Context, userID int) error {
	query := `
		UPDATE daily_rewards
		SET streak = 0, updated = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to reset streak: %w", err)
	}

	return nil
}

// formatDuration formats duration as HH:MM:SS
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

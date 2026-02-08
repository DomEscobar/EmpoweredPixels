package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ResonanceAchievementRepository handles persistence of resonance achievements
type ResonanceAchievementRepository struct {
	pool *pgxpool.Pool
}

// ResonanceAchievement represents an achievement record
type ResonanceAchievement struct {
	ID               string
	UserID           int64
	AchievementType  string
	ResonanceScoreMax int
	MatchesCount     int
	UnlockedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

const (
	AchievementResonanceMaster    = "RESONANCE_MASTER"
	AchievementResonanceCollector = "RESONANCE_COLLECTOR"
	AchievementHarmonyPerfectionist = "HARMONY_PERFECTIONIST"
)

// NewResonanceAchievementRepository creates a new repository
func NewResonanceAchievementRepository(pool *pgxpool.Pool) *ResonanceAchievementRepository {
	return &ResonanceAchievementRepository{pool: pool}
}

// GetOrCreate gets an achievement or creates it if it doesn't exist
func (r *ResonanceAchievementRepository) GetOrCreate(ctx context.Context, userID int64, achievementType string) (*ResonanceAchievement, error) {
	achievement, err := r.GetByUserAndType(ctx, userID, achievementType)
	if err != nil {
		return nil, err
	}
	if achievement != nil {
		return achievement, nil
	}

	// Create new achievement
	return r.Create(ctx, userID, achievementType)
}

// Create creates a new achievement
func (r *ResonanceAchievementRepository) Create(ctx context.Context, userID int64, achievementType string) (*ResonanceAchievement, error) {
	const query = `
		INSERT INTO resonance_achievements (user_id, achievement_type, resonance_score_max, matches_count, created_at, updated_at)
		VALUES ($1, $2, 0, 0, NOW(), NOW())
		ON CONFLICT (user_id, achievement_type) DO UPDATE
		SET updated_at = NOW()
		RETURNING id, user_id, achievement_type, resonance_score_max, matches_count, unlocked_at, created_at, updated_at
	`

	var a ResonanceAchievement
	err := r.pool.QueryRow(ctx, query, userID, achievementType).Scan(
		&a.ID, &a.UserID, &a.AchievementType, &a.ResonanceScoreMax, &a.MatchesCount, &a.UnlockedAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create achievement: %w", err)
	}

	return &a, nil
}

// GetByUserAndType gets an achievement by user ID and type
func (r *ResonanceAchievementRepository) GetByUserAndType(ctx context.Context, userID int64, achievementType string) (*ResonanceAchievement, error) {
	const query = `
		SELECT id, user_id, achievement_type, resonance_score_max, matches_count, unlocked_at, created_at, updated_at
		FROM resonance_achievements
		WHERE user_id = $1 AND achievement_type = $2
	`

	var a ResonanceAchievement
	err := r.pool.QueryRow(ctx, query, userID, achievementType).Scan(
		&a.ID, &a.UserID, &a.AchievementType, &a.ResonanceScoreMax, &a.MatchesCount, &a.UnlockedAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get achievement: %w", err)
	}

	return &a, nil
}

// UpdateMatchCount increments match count and checks for unlock
func (r *ResonanceAchievementRepository) UpdateMatchCount(ctx context.Context, userID int64, achievementType string, harmonyScore int) error {
	const query = `
		UPDATE resonance_achievements
		SET matches_count = matches_count + 1,
		    resonance_score_max = GREATEST(resonance_score_max, $3),
		    unlocked_at = CASE
		        WHEN achievement_type = $2 AND matches_count + 1 >= 100 AND $3 >= 80 THEN NOW()
		        ELSE unlocked_at
		    END,
		    updated_at = NOW()
		WHERE user_id = $1 AND achievement_type = $2
	`

	_, err := r.pool.Exec(ctx, query, userID, achievementType, harmonyScore)
	if err != nil {
		return fmt.Errorf("failed to update achievement: %w", err)
	}

	return nil
}

// GetAllByUser gets all achievements for a user
func (r *ResonanceAchievementRepository) GetAllByUser(ctx context.Context, userID int64) ([]ResonanceAchievement, error) {
	const query = `
		SELECT id, user_id, achievement_type, resonance_score_max, matches_count, unlocked_at, created_at, updated_at
		FROM resonance_achievements
		WHERE user_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query achievements: %w", err)
	}
	defer rows.Close()

	var achievements []ResonanceAchievement
	for rows.Next() {
		var a ResonanceAchievement
		if err := rows.Scan(&a.ID, &a.UserID, &a.AchievementType, &a.ResonanceScoreMax, &a.MatchesCount, &a.UnlockedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %w", err)
		}
		achievements = append(achievements, a)
	}

	return achievements, rows.Err()
}

// GetUnlocked gets all unlocked achievements for a user
func (r *ResonanceAchievementRepository) GetUnlocked(ctx context.Context, userID int64) ([]ResonanceAchievement, error) {
	const query = `
		SELECT id, user_id, achievement_type, resonance_score_max, matches_count, unlocked_at, created_at, updated_at
		FROM resonance_achievements
		WHERE user_id = $1 AND unlocked_at IS NOT NULL
		ORDER BY unlocked_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query achievements: %w", err)
	}
	defer rows.Close()

	var achievements []ResonanceAchievement
	for rows.Next() {
		var a ResonanceAchievement
		if err := rows.Scan(&a.ID, &a.UserID, &a.AchievementType, &a.ResonanceScoreMax, &a.MatchesCount, &a.UnlockedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %w", err)
		}
		achievements = append(achievements, a)
	}

	return achievements, rows.Err()
}

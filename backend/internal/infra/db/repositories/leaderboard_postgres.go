package repositories

import (
	"context"
	"fmt"

	"empoweredpixels/internal/domain/leaderboard"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LeaderboardRepository defines leaderboard operations
type LeaderboardRepository interface {
	GetByCategory(ctx context.Context, category string, limit int, offset int) ([]leaderboard.Entry, error)
	GetUserRank(ctx context.Context, category string, userID int) (*leaderboard.Entry, error)
	UpsertEntry(ctx context.Context, entry *leaderboard.Entry) error
	GetTotalCount(ctx context.Context, category string) (int, error)
	GetNearbyRanks(ctx context.Context, category string, userID int, rangeSize int) ([]leaderboard.Entry, error)
}

// AchievementRepository defines achievement operations
type AchievementRepository interface {
	ListAchievements(ctx context.Context) ([]leaderboard.Achievement, error)
	GetPlayerAchievements(ctx context.Context, userID int) ([]leaderboard.PlayerAchievement, error)
	UpdateAchievementProgress(ctx context.Context, userID int, achievementKey string, progress int) error
	ClaimAchievementReward(ctx context.Context, userID int, achievementID string) error
}

// LeaderboardPostgres implements LeaderboardRepository
type LeaderboardPostgres struct {
	db *pgxpool.Pool
}

// AchievementPostgres implements AchievementRepository
type AchievementPostgres struct {
	db *pgxpool.Pool
}

// NewLeaderboardRepository creates a new leaderboard repository
func NewLeaderboardRepository(db *pgxpool.Pool) LeaderboardRepository {
	return &LeaderboardPostgres{db: db}
}

// NewAchievementRepository creates a new achievement repository
func NewAchievementRepository(db *pgxpool.Pool) AchievementRepository {
	return &AchievementPostgres{db: db}
}

// GetByCategory retrieves leaderboard entries for a category
func (r *LeaderboardPostgres) GetByCategory(ctx context.Context, category string, limit int, offset int) ([]leaderboard.Entry, error) {
	query := `
		SELECT l.id, l.category, l.user_id, u.username, l.rank, l.score, l.previous_rank, l.trend, l.updated_at
		FROM leaderboard_entries l
		JOIN users u ON u.id = l.user_id
		WHERE l.category = $1
		ORDER BY l.rank ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, category, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []leaderboard.Entry
	for rows.Next() {
		var e leaderboard.Entry
		if err := rows.Scan(&e.ID, &e.Category, &e.UserID, &e.Username, &e.Rank, &e.Score, &e.PreviousRank, &e.Trend, &e.UpdatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// GetUserRank retrieves a user's rank in a category
func (r *LeaderboardPostgres) GetUserRank(ctx context.Context, category string, userID int) (*leaderboard.Entry, error) {
	query := `
		SELECT l.id, l.category, l.user_id, u.username, l.rank, l.score, l.previous_rank, l.trend, l.updated_at
		FROM leaderboard_entries l
		JOIN users u ON u.id = l.user_id
		WHERE l.category = $1 AND l.user_id = $2
	`

	var e leaderboard.Entry
	err := r.db.QueryRow(ctx, query, category, userID).Scan(
		&e.ID, &e.Category, &e.UserID, &e.Username, &e.Rank, &e.Score, &e.PreviousRank, &e.Trend, &e.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// UpsertEntry creates or updates a leaderboard entry
func (r *LeaderboardPostgres) UpsertEntry(ctx context.Context, entry *leaderboard.Entry) error {
	query := `
		INSERT INTO leaderboard_entries (category, user_id, rank, score, previous_rank, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (category, user_id)
		DO UPDATE SET rank = EXCLUDED.rank, score = EXCLUDED.score, previous_rank = leaderboard_entries.rank, updated_at = NOW()
	`

	_, err := r.db.Exec(ctx, query, entry.Category, entry.UserID, entry.Rank, entry.Score, entry.PreviousRank)
	return err
}

// GetTotalCount returns total entries in a category
func (r *LeaderboardPostgres) GetTotalCount(ctx context.Context, category string) (int, error) {
	query := `SELECT COUNT(*) FROM leaderboard_entries WHERE category = $1`

	var count int
	err := r.db.QueryRow(ctx, query, category).Scan(&count)
	return count, err
}

// GetNearbyRanks retrieves entries near a user's rank
func (r *LeaderboardPostgres) GetNearbyRanks(ctx context.Context, category string, userID int, rangeSize int) ([]leaderboard.Entry, error) {
	// First get user's rank
	userEntry, err := r.GetUserRank(ctx, category, userID)
	if err != nil || userEntry == nil {
		return nil, err
	}

	query := `
		SELECT l.id, l.category, l.user_id, u.username, l.rank, l.score, l.previous_rank, l.trend, l.updated_at
		FROM leaderboard_entries l
		JOIN users u ON u.id = l.user_id
		WHERE l.category = $1 AND l.rank BETWEEN $2 AND $3
		ORDER BY l.rank ASC
	`

	minRank := userEntry.Rank - rangeSize
	if minRank < 1 {
		minRank = 1
	}
	maxRank := userEntry.Rank + rangeSize

	rows, err := r.db.Query(ctx, query, category, minRank, maxRank)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []leaderboard.Entry
	for rows.Next() {
		var e leaderboard.Entry
		if err := rows.Scan(&e.ID, &e.Category, &e.UserID, &e.Username, &e.Rank, &e.Score, &e.PreviousRank, &e.Trend, &e.UpdatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// ListAchievements retrieves all achievements
func (r *AchievementPostgres) ListAchievements(ctx context.Context) ([]leaderboard.Achievement, error) {
	query := `
		SELECT id, key, name, description, icon, category, requirement_type, requirement_value, reward_gold, reward_title, hidden, created_at
		FROM achievements
		ORDER BY category, requirement_value
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []leaderboard.Achievement
	for rows.Next() {
		var a leaderboard.Achievement
		if err := rows.Scan(&a.ID, &a.Key, &a.Name, &a.Description, &a.Icon, &a.Category, &a.RequirementType, &a.RequirementValue, &a.RewardGold, &a.RewardTitle, &a.Hidden, &a.CreatedAt); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}

	return achievements, rows.Err()
}

// GetPlayerAchievements retrieves a player's achievements
func (r *AchievementPostgres) GetPlayerAchievements(ctx context.Context, userID int) ([]leaderboard.PlayerAchievement, error) {
	query := `
		SELECT pa.id, pa.user_id, pa.achievement_id, pa.progress, pa.completed, pa.completed_at, pa.claimed, pa.claimed_at,
		       a.id, a.key, a.name, a.description, a.icon, a.category, a.requirement_type, a.requirement_value, a.reward_gold, a.reward_title, a.hidden, a.created_at
		FROM player_achievements pa
		JOIN achievements a ON a.id = pa.achievement_id
		WHERE pa.user_id = $1
		ORDER BY pa.completed, a.requirement_value
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []leaderboard.PlayerAchievement
	for rows.Next() {
		var pa leaderboard.PlayerAchievement
		var a leaderboard.Achievement
		if err := rows.Scan(
			&pa.ID, &pa.UserID, &pa.AchievementID, &pa.Progress, &pa.Completed, &pa.CompletedAt, &pa.Claimed, &pa.ClaimedAt,
			&a.ID, &a.Key, &a.Name, &a.Description, &a.Icon, &a.Category, &a.RequirementType, &a.RequirementValue, &a.RewardGold, &a.RewardTitle, &a.Hidden, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		pa.Achievement = &a
		achievements = append(achievements, pa)
	}

	return achievements, rows.Err()
}

// UpdateAchievementProgress updates a player's achievement progress
func (r *AchievementPostgres) UpdateAchievementProgress(ctx context.Context, userID int, achievementKey string, progress int) error {
	query := `
		INSERT INTO player_achievements (user_id, achievement_id, progress, completed, completed_at)
		SELECT $1, a.id, $3, $3 >= a.requirement_value, CASE WHEN $3 >= a.requirement_value THEN NOW() ELSE NULL END
		FROM achievements a
		WHERE a.key = $2
		ON CONFLICT (user_id, achievement_id)
		DO UPDATE SET
			progress = EXCLUDED.progress,
			completed = EXCLUDED.completed,
			completed_at = COALESCE(player_achievements.completed_at, EXCLUDED.completed_at)
	`

	_, err := r.db.Exec(ctx, query, userID, achievementKey, progress)
	return err
}

// ClaimAchievementReward marks an achievement reward as claimed
func (r *AchievementPostgres) ClaimAchievementReward(ctx context.Context, userID int, achievementID string) error {
	query := `
		UPDATE player_achievements
		SET claimed = true, claimed_at = NOW()
		WHERE user_id = $1 AND achievement_id = $2 AND completed = true AND claimed = false
	`

	_, err := r.db.Exec(ctx, query, userID, achievementID)
	return err
}

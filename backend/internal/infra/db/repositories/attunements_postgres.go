package repositories

import (
	"context"
	"errors"
	"time"

	"empoweredpixels/internal/domain/attunements"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAttunementNotFound = errors.New("attunement not found")
)

type AttunementRepository struct {
	pool *pgxpool.Pool
}

func NewAttunementRepository(pool *pgxpool.Pool) *AttunementRepository {
	return &AttunementRepository{pool: pool}
}

// PlayerAttunement represents a player's progress in a single element
type PlayerAttunement struct {
	ID            int64
	UserID        int64
	Element       attunements.Element
	Level         int
	XP            int
	TotalXPEarned int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// GetPlayerAttunements returns all 6 attunements for a user
func (r *AttunementRepository) GetPlayerAttunements(ctx context.Context, userID int64) ([]PlayerAttunement, error) {
	const query = `
		SELECT id, user_id, element, level, xp, total_xp_earned, created_at, updated_at
		FROM player_attunements
		WHERE user_id = $1
		ORDER BY element ASC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []PlayerAttunement
	for rows.Next() {
		var a PlayerAttunement
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Element,
			&a.Level,
			&a.XP,
			&a.TotalXPEarned,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, rows.Err()
}

// GetAttunement returns a specific attunement for a user and element
func (r *AttunementRepository) GetAttunement(ctx context.Context, userID int64, element attunements.Element) (*PlayerAttunement, error) {
	const query = `
		SELECT id, user_id, element, level, xp, total_xp_earned, created_at, updated_at
		FROM player_attunements
		WHERE user_id = $1 AND element = $2
	`

	var a PlayerAttunement
	err := r.pool.QueryRow(ctx, query, userID, int(element)).Scan(
		&a.ID,
		&a.UserID,
		&a.Element,
		&a.Level,
		&a.XP,
		&a.TotalXPEarned,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAttunementNotFound
	}
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// AddXP adds XP to an attunement, handling level-ups, and logs to history
func (r *AttunementRepository) AddXP(ctx context.Context, userID int64, element attunements.Element, xpAmount int, source string) (*PlayerAttunement, bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback(ctx)

	// Get current attunement
	const selectQuery = `
		SELECT id, user_id, element, level, xp, total_xp_earned, created_at, updated_at
		FROM player_attunements
		WHERE user_id = $1 AND element = $2
		FOR UPDATE
	`

	var a PlayerAttunement
	err = tx.QueryRow(ctx, selectQuery, userID, int(element)).Scan(
		&a.ID,
		&a.UserID,
		&a.Element,
		&a.Level,
		&a.XP,
		&a.TotalXPEarned,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, false, ErrAttunementNotFound
	}
	if err != nil {
		return nil, false, err
	}

	oldLevel := a.Level
	a.XP += xpAmount
	a.TotalXPEarned += xpAmount
	leveledUp := false

	// Process level-ups
	for a.Level < 25 {
		xpNeeded := attunements.XPNeededForLevel(a.Level)
		if a.XP >= xpNeeded {
			a.XP -= xpNeeded
			a.Level++
			leveledUp = true
		} else {
			break
		}
	}

	// Cap XP at max level
	if a.Level >= 25 {
		a.XP = 0
	}

	// Update attunement
	const updateQuery = `
		UPDATE player_attunements
		SET level = $1, xp = $2, total_xp_earned = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err = tx.Exec(ctx, updateQuery, a.Level, a.XP, a.TotalXPEarned, a.ID)
	if err != nil {
		return nil, false, err
	}

	// Log XP history
	const historyQuery = `
		INSERT INTO attunement_xp_history (user_id, element, xp_amount, source, old_level, new_level, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err = tx.Exec(ctx, historyQuery, userID, int(element), xpAmount, source, oldLevel, a.Level)
	if err != nil {
		return nil, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, false, err
	}

	a.UpdatedAt = time.Now()
	return &a, leveledUp, nil
}

// CreateInitialAttunements creates all 6 elements at level 1 for a new user
func (r *AttunementRepository) CreateInitialAttunements(ctx context.Context, userID int64) error {
	const query = `
		INSERT INTO player_attunements (user_id, element, level, xp, total_xp_earned, created_at, updated_at)
		VALUES ($1, $2, 1, 0, 0, NOW(), NOW())
		ON CONFLICT (user_id, element) DO NOTHING
	`

	for _, element := range attunements.AllElements() {
		_, err := r.pool.Exec(ctx, query, userID, int(element))
		if err != nil {
			return err
		}
	}

	return nil
}

// HasAttunements checks if a user has attunements initialized
func (r *AttunementRepository) HasAttunements(ctx context.Context, userID int64) (bool, error) {
	const query = `SELECT COUNT(*) FROM player_attunements WHERE user_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 6, nil
}

// GetXPHistory returns XP history for a user
func (r *AttunementRepository) GetXPHistory(ctx context.Context, userID int64, limit int) ([]XPHistoryEntry, error) {
	const query = `
		SELECT id, user_id, element, xp_amount, source, old_level, new_level, created_at
		FROM attunement_xp_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []XPHistoryEntry
	for rows.Next() {
		var h XPHistoryEntry
		err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.Element,
			&h.XPAmount,
			&h.Source,
			&h.OldLevel,
			&h.NewLevel,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, h)
	}

	return result, rows.Err()
}

// XPHistoryEntry represents a single XP gain event
type XPHistoryEntry struct {
	ID        int64
	UserID    int64
	Element   attunements.Element
	XPAmount  int
	Source    string
	OldLevel  int
	NewLevel  int
	CreatedAt time.Time
}

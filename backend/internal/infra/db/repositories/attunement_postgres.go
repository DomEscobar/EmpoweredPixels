package repositories

import (
	"context"
	"fmt"

	"empoweredpixels/internal/domain/attunement"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AttunementPostgres implements your AttunementRepository interface
type AttunementPostgres struct {
	pool *pgxpool.Pool
}

// NewAttunementRepository creates a new attunement repository
func NewAttunementRepository(pool *pgxpool.Pool) *AttunementPostgres {
	return &AttunementPostgres{pool: pool}
}

// GetPlayerAttunements retrieves all 6 attunements for a player
func (r *AttunementPostgres) GetPlayerAttunements(ctx context.Context, userID int) ([]attunement.Attunement, error) {
	query := `
		SELECT element, level, current_xp, total_xp
		FROM player_attunements
		WHERE user_id = $1
		ORDER BY element ASC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query attunements: %w", err)
	}
	defer rows.Close()

	var attunements []attunement.Attunement
	for rows.Next() {
		var a attunement.Attunement
		if err := rows.Scan(&a.Element, &a.Level, &a.CurrentXP, &a.TotalXP); err != nil {
			return nil, fmt.Errorf("failed to scan attunement: %w", err)
		}
		attunements = append(attunements, a)
	}

	return attunements, rows.Err()
}

// GetAttunement retrieves a specific element attunement for a player
func (r *AttunementPostgres) GetAttunement(ctx context.Context, userID int, element attunement.Element) (*attunement.Attunement, error) {
	query := `
		SELECT element, level, current_xp, total_xp
		FROM player_attunements
		WHERE user_id = $1 AND element = $2
	`

	var a attunement.Attunement
	err := r.pool.QueryRow(ctx, query, userID, string(element)).Scan(
		&a.Element, &a.Level, &a.CurrentXP, &a.TotalXP,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get attunement: %w", err)
	}

	return &a, nil
}

// AddXP adds XP to an element and handles level-ups
func (r *AttunementPostgres) AddXP(ctx context.Context, userID int, element attunement.Element, xpAmount int, source string) (levelUp bool, newLevel int, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get current attunement
	var currentLevel, currentXP, totalXP int
	query := `
		SELECT level, current_xp, total_xp
		FROM player_attunements
		WHERE user_id = $1 AND element = $2
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, query, userID, string(element)).Scan(&currentLevel, &currentXP, &totalXP)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Create initial attunement if not exists
			insertQuery := `
				INSERT INTO player_attunements (user_id, element, level, current_xp, total_xp)
				VALUES ($1, $2, 1, 0, 0)
				RETURNING level, current_xp, total_xp
			`
			err = tx.QueryRow(ctx, insertQuery, userID, string(element)).Scan(&currentLevel, &currentXP, &totalXP)
			if err != nil {
				return false, 0, fmt.Errorf("failed to create initial attunement: %w", err)
			}
		} else {
			return false, 0, fmt.Errorf("failed to get current attunement: %w", err)
		}
	}

	// Calculate new XP and level
	newXP := currentXP + xpAmount
	newTotalXP := totalXP + xpAmount
	newLevel = currentLevel
	levelUp = false

	// Check for level-ups (max level 25)
	for newLevel < 25 {
		xpRequired := attunement.GetXPRequired(newLevel)
		if xpRequired == 0 || newXP < xpRequired {
			break
		}
		newXP -= xpRequired
		newLevel++
		levelUp = true
	}

	// Update attunement
	updateQuery := `
		UPDATE player_attunements
		SET level = $3, current_xp = $4, total_xp = $5, updated = NOW()
		WHERE user_id = $1 AND element = $2
	`
	_, err = tx.Exec(ctx, updateQuery, userID, string(element), newLevel, newXP, newTotalXP)
	if err != nil {
		return false, 0, fmt.Errorf("failed to update attunement: %w", err)
	}

	// Record XP history
	historyQuery := `
		INSERT INTO attunement_xp_history (user_id, element, xp_gained, source, new_level)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, historyQuery, userID, string(element), xpAmount, source, newLevel)
	if err != nil {
		return false, 0, fmt.Errorf("failed to record xp history: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return false, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return levelUp, newLevel, nil
}

// CreateInitialAttunements creates all 6 elements at level 1 for a new player
func (r *AttunementPostgres) CreateInitialAttunements(ctx context.Context, userID int) error {
	elements := []attunement.Element{
		attunement.Fire, attunement.Water, attunement.Earth,
		attunement.Air, attunement.Light, attunement.Dark,
	}

	query := `
		INSERT INTO player_attunements (user_id, element, level, current_xp, total_xp)
		VALUES ($1, $2, 1, 0, 0)
		ON CONFLICT (user_id, element) DO NOTHING
	`

	for _, element := range elements {
		_, err := r.pool.Exec(ctx, query, userID, string(element))
		if err != nil {
			return fmt.Errorf("failed to create attunement for %s: %w", element, err)
		}
	}

	return nil
}

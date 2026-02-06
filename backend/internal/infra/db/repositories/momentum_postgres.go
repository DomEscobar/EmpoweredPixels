package repositories

import (
	"context"
	"database/sql"
	"time"

	"empoweredpixels/internal/domain/momentum"
	momentumusecase "empoweredpixels/internal/usecase/momentum"
)

type MomentumPostgres struct {
	db *sql.DB
}

func NewMomentumPostgres(db *sql.DB) *MomentumPostgres {
	return &MomentumPostgres{db: db}
}

func (r *MomentumPostgres) GetStakedMomentum(ctx context.Context, fighterID string) (*momentum.StakedMomentum, error) {
	var sm momentum.StakedMomentum
	sm.FighterID = fighterID

	query := `SELECT current_value, last_updated_at FROM fighter_momentum WHERE fighter_id = $1`
	err := r.db.QueryRowContext(ctx, query, fighterID).Scan(&sm.CurrentValue, &sm.LastUpdatedAt)
	
	if err == sql.ErrNoRows {
		sm.CurrentValue = 0
		sm.LastUpdatedAt = time.Now()
		return &sm, nil
	}
	
	if err != nil {
		return nil, err
	}

	return &sm, nil
}

func (r *MomentumPostgres) UpdateStakedMomentum(ctx context.Context, sm *momentum.StakedMomentum) error {
	query := `
		INSERT INTO fighter_momentum (fighter_id, current_value, last_updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (fighter_id) DO UPDATE
		SET current_value = EXCLUDED.current_value,
		    last_updated_at = EXCLUDED.last_updated_at
	`
	_, err := r.db.ExecContext(ctx, query, sm.FighterID, sm.CurrentValue, sm.LastUpdatedAt)
	return err
}

// Ensure interface compliance
var _ momentumusecase.MomentumRepository = (*MomentumPostgres)(nil)

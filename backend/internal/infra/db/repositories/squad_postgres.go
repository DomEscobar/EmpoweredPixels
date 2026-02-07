package repositories

import (
	"context"
	"errors"

	"empoweredpixels/internal/domain/roster"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SquadRepository struct {
	pool *pgxpool.Pool
}

func NewSquadRepository(pool *pgxpool.Pool) *SquadRepository {
	return &SquadRepository{pool: pool}
}

func (r *SquadRepository) Create(ctx context.Context, squad *roster.Squad) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const squadQuery = `
		insert into squads (id, user_id, name, is_active)
		values ($1, $2, $3, $4)`
	
	_, err = tx.Exec(ctx, squadQuery, squad.ID, squad.UserID, squad.Name, squad.IsActive)
	if err != nil {
		return err
	}

	const memberQuery = `
		insert into squad_members (squad_id, fighter_id, slot_index)
		values ($1, $2, $3)`
	
	for _, m := range squad.Members {
		_, err = tx.Exec(ctx, memberQuery, squad.ID, m.FighterID, m.SlotIndex)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *SquadRepository) GetActiveByUserID(ctx context.Context, userID int64) (*roster.Squad, error) {
	const squadQuery = `
		select id, user_id, name, is_active, created_at, updated_at
		from squads
		where user_id = $1 and is_active = true`

	var squad roster.Squad
	err := r.pool.QueryRow(ctx, squadQuery, userID).Scan(
		&squad.ID, &squad.UserID, &squad.Name, &squad.IsActive, &squad.CreatedAt, &squad.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	const memberQuery = `
		select fighter_id, slot_index
		from squad_members
		where squad_id = $1
		order by slot_index asc`
	
	rows, err := r.pool.Query(ctx, memberQuery, squad.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m roster.Member
		if err := rows.Scan(&m.FighterID, &m.SlotIndex); err != nil {
			return nil, err
		}
		squad.Members = append(squad.Members, m)
	}

	return &squad, nil
}

func (r *SquadRepository) DeactivateAll(ctx context.Context, userID int64) error {
	const query = `update squads set is_active = false where user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

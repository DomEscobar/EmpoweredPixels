package repositories

import (
	"context"
	"encoding/json"

	"empoweredpixels/internal/domain/combat"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CombatRepository struct {
	pool *pgxpool.Pool
}

func NewCombatRepository(pool *pgxpool.Pool) *CombatRepository {
	return &CombatRepository{pool: pool}
}

func (r *CombatRepository) SaveLogs(ctx context.Context, matchID string, roundTicks []combat.RoundTick) error {
	const query = `
		insert into combat_logs (match_id, round, tick, event_type, payload)
		values ($1, $2, $3, $4, $5)`

	batch := &pgx.Batch{}
	for _, rt := range roundTicks {
		for i, tick := range rt.Ticks {
			batch.Queue(query, matchID, rt.Round, i, tick.Type, tick.Payload)
		}
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

type BattleSummary struct {
	WinnerID    *string               `json:"winnerId"`
	TotalRounds int                   `json:"totalRounds"`
	Scores      []combat.FighterScore `json:"scores"`
}

func (r *CombatRepository) SaveSummary(ctx context.Context, matchID string, winnerID *string, totalRounds int, scores []combat.FighterScore) error {
	summary := BattleSummary{
		WinnerID:    winnerID,
		TotalRounds: totalRounds,
		Scores:      scores,
	}
	summaryJSON, _ := json.Marshal(summary)

	const query = `
		insert into battle_details (match_id, winner_id, total_rounds, summary)
		values ($1, $2, $3, $4)
		on conflict (match_id) do update set
			winner_id = excluded.winner_id,
			total_rounds = excluded.total_rounds,
			summary = excluded.summary`

	_, err := r.pool.Exec(ctx, query, matchID, winnerID, totalRounds, summaryJSON)
	return err
}

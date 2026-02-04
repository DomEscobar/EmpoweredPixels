package repositories

import (
	"context"

	"empoweredpixels/internal/domain/seasons"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeasonSummaryRepository struct {
	pool *pgxpool.Pool
}

func NewSeasonSummaryRepository(pool *pgxpool.Pool) *SeasonSummaryRepository {
	return &SeasonSummaryRepository{pool: pool}
}

func (r *SeasonSummaryRepository) ListByUser(ctx context.Context, userID int64, limit int, offset int) ([]seasons.SeasonSummary, error) {
	const query = `
		select id, user_id, season_id, position
		from season_summaries
		where user_id = $1
		order by season_id desc
		limit $2 offset $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []seasons.SeasonSummary
	for rows.Next() {
		var summary seasons.SeasonSummary
		if err := rows.Scan(&summary.ID, &summary.UserID, &summary.SeasonID, &summary.Position); err != nil {
			return nil, err
		}
		result = append(result, summary)
	}
	return result, rows.Err()
}

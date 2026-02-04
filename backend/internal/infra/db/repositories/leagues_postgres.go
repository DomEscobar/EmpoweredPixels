package repositories

import (
	"context"
	"errors"
	"time"

	"empoweredpixels/internal/domain/leagues"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeagueRepository struct {
	pool *pgxpool.Pool
}

type LeagueSubscriptionRepository struct {
	pool *pgxpool.Pool
}

type LeagueMatchRepository struct {
	pool *pgxpool.Pool
}

func NewLeagueRepository(pool *pgxpool.Pool) *LeagueRepository {
	return &LeagueRepository{pool: pool}
}

func NewLeagueSubscriptionRepository(pool *pgxpool.Pool) *LeagueSubscriptionRepository {
	return &LeagueSubscriptionRepository{pool: pool}
}

func NewLeagueMatchRepository(pool *pgxpool.Pool) *LeagueMatchRepository {
	return &LeagueMatchRepository{pool: pool}
}

func (r *LeagueRepository) List(ctx context.Context) ([]leagues.League, error) {
	const query = `
		select id, name, options, is_deactivated
		from leagues
		where is_deactivated = false
		order by id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []leagues.League
	for rows.Next() {
		var league leagues.League
		if err := rows.Scan(&league.ID, &league.Name, &league.Options, &league.IsDeactivated); err != nil {
			return nil, err
		}
		result = append(result, league)
	}
	return result, rows.Err()
}

func (r *LeagueRepository) GetByID(ctx context.Context, id int) (*leagues.League, error) {
	const query = `
		select id, name, options, is_deactivated
		from leagues
		where id = $1`

	var league leagues.League
	err := r.pool.QueryRow(ctx, query, id).Scan(&league.ID, &league.Name, &league.Options, &league.IsDeactivated)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &league, nil
}

func (r *LeagueSubscriptionRepository) ListByLeague(ctx context.Context, leagueID int) ([]leagues.LeagueSubscription, error) {
	const query = `
		select league_id, fighter_id, created
		from league_subscriptions
		where league_id = $1`

	rows, err := r.pool.Query(ctx, query, leagueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []leagues.LeagueSubscription
	for rows.Next() {
		var sub leagues.LeagueSubscription
		if err := rows.Scan(&sub.LeagueID, &sub.FighterID, &sub.Created); err != nil {
			return nil, err
		}
		result = append(result, sub)
	}
	return result, rows.Err()
}

func (r *LeagueSubscriptionRepository) ListByLeagueAndUser(ctx context.Context, leagueID int, userID int64) ([]leagues.LeagueSubscription, error) {
	const query = `
		select ls.league_id, ls.fighter_id, ls.created
		from league_subscriptions ls
		join fighters f on f.id = ls.fighter_id
		where ls.league_id = $1 and f.user_id = $2`

	rows, err := r.pool.Query(ctx, query, leagueID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []leagues.LeagueSubscription
	for rows.Next() {
		var sub leagues.LeagueSubscription
		if err := rows.Scan(&sub.LeagueID, &sub.FighterID, &sub.Created); err != nil {
			return nil, err
		}
		result = append(result, sub)
	}
	return result, rows.Err()
}

func (r *LeagueMatchRepository) Create(ctx context.Context, leagueID int, matchID string) error {
	const query = `
		insert into league_matches (league_id, match_id, started)
		values ($1, $2, null)
		on conflict (league_id, match_id) do nothing`
	_, err := r.pool.Exec(ctx, query, leagueID, matchID)
	return err
}

func (r *LeagueMatchRepository) UpdateStarted(ctx context.Context, leagueID int, matchID string, started *time.Time) error {
	const query = `update league_matches set started = $3 where league_id = $1 and match_id = $2`
	_, err := r.pool.Exec(ctx, query, leagueID, matchID, started)
	return err
}

func (r *LeagueSubscriptionRepository) Create(ctx context.Context, subscription *leagues.LeagueSubscription) error {
	const query = `
		insert into league_subscriptions (league_id, fighter_id, created)
		values ($1, $2, $3)
		on conflict (league_id, fighter_id) do nothing`

	_, err := r.pool.Exec(ctx, query, subscription.LeagueID, subscription.FighterID, subscription.Created)
	return err
}

func (r *LeagueSubscriptionRepository) Delete(ctx context.Context, subscription *leagues.LeagueSubscription) error {
	const query = `
		delete from league_subscriptions
		where league_id = $1 and fighter_id = $2`

	result, err := r.pool.Exec(ctx, query, subscription.LeagueID, subscription.FighterID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *LeagueMatchRepository) ListByLeague(ctx context.Context, leagueID int, limit int, offset int) ([]leagues.LeagueMatch, error) {
	const query = `
		select league_id, match_id, started
		from league_matches
		where league_id = $1
		order by started desc nulls last
		limit $2 offset $3`

	rows, err := r.pool.Query(ctx, query, leagueID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []leagues.LeagueMatch
	for rows.Next() {
		var match leagues.LeagueMatch
		if err := rows.Scan(&match.LeagueID, &match.MatchID, &match.Started); err != nil {
			return nil, err
		}
		result = append(result, match)
	}
	return result, rows.Err()
}

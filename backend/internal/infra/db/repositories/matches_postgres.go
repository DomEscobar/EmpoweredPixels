package repositories

import (
	"context"
	"errors"

	"empoweredpixels/internal/domain/matches"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository struct {
	pool *pgxpool.Pool
}

type MatchTeamRepository struct {
	pool *pgxpool.Pool
}

type MatchRegistrationRepository struct {
	pool *pgxpool.Pool
}

type MatchResultRepository struct {
	pool *pgxpool.Pool
}

type MatchScoreRepository struct {
	pool *pgxpool.Pool
}

func NewMatchRepository(pool *pgxpool.Pool) *MatchRepository {
	return &MatchRepository{pool: pool}
}

func NewMatchTeamRepository(pool *pgxpool.Pool) *MatchTeamRepository {
	return &MatchTeamRepository{pool: pool}
}

func NewMatchRegistrationRepository(pool *pgxpool.Pool) *MatchRegistrationRepository {
	return &MatchRegistrationRepository{pool: pool}
}

func NewMatchResultRepository(pool *pgxpool.Pool) *MatchResultRepository {
	return &MatchResultRepository{pool: pool}
}

func NewMatchScoreRepository(pool *pgxpool.Pool) *MatchScoreRepository {
	return &MatchScoreRepository{pool: pool}
}

func (r *MatchRepository) Create(ctx context.Context, match *matches.Match) error {
	const query = `
		insert into matches (id, creator_user_id, created, started, completed_at, cancelled_at, status, options)
		values ($1, $2, $3, $4, $5, $6, $7, $8)`
	status := match.Status
	if status == "" {
		status = matches.MatchStatusLobby
	}
	_, err := r.pool.Exec(ctx, query,
		match.ID, match.CreatorUserID, match.Created, match.Started,
		match.CompletedAt, match.CancelledAt, status, match.Options)
	return err
}

func (r *MatchRepository) GetByID(ctx context.Context, id string) (*matches.Match, error) {
	const query = `
		select id, creator_user_id, created, started, completed_at, cancelled_at, status, options
		from matches
		where id = $1`

	var match matches.Match
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&match.ID, &match.CreatorUserID, &match.Created, &match.Started,
		&match.CompletedAt, &match.CancelledAt, &match.Status, &match.Options,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) ListOpen(ctx context.Context, limit int, offset int) ([]matches.Match, error) {
	return r.ListByStatus(ctx, matches.MatchStatusLobby, limit, offset)
}

func (r *MatchRepository) ListByStatus(ctx context.Context, status string, limit int, offset int) ([]matches.Match, error) {
	const query = `
		select id, creator_user_id, created, started, completed_at, cancelled_at, status, options
		from matches
		where status = $1
		order by created desc
		limit $2 offset $3`

	rows, err := r.pool.Query(ctx, query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []matches.Match
	for rows.Next() {
		var match matches.Match
		if err := rows.Scan(&match.ID, &match.CreatorUserID, &match.Created, &match.Started,
			&match.CompletedAt, &match.CancelledAt, &match.Status, &match.Options); err != nil {
			return nil, err
		}
		result = append(result, match)
	}

	return result, rows.Err()
}

func (r *MatchRepository) GetCurrentMatch(ctx context.Context, userID int64) (*matches.Match, error) {
	const query = `
		select m.id, m.creator_user_id, m.created, m.started, m.completed_at, m.cancelled_at, m.status, m.options
		from matches m
		join match_registrations mr on mr.match_id = m.id
		join fighters f on f.id = mr.fighter_id
		where f.user_id = $1 and m.status in ('lobby', 'running')
		order by m.created desc
		limit 1`

	var match matches.Match
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&match.ID, &match.CreatorUserID, &match.Created, &match.Started,
		&match.CompletedAt, &match.CancelledAt, &match.Status, &match.Options,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) Update(ctx context.Context, match *matches.Match) error {
	const query = `
		update matches
		set started = $2, completed_at = $3, cancelled_at = $4, status = $5
		where id = $1`
	_, err := r.pool.Exec(ctx, query,
		match.ID, match.Started, match.CompletedAt, match.CancelledAt, match.Status)
	return err
}

func (r *MatchRepository) ListStaleLobbies(ctx context.Context, olderThanMinutes int) ([]matches.Match, error) {
	const query = `
		select id, creator_user_id, created, started, completed_at, cancelled_at, status, options
		from matches
		where status = 'lobby' and created < now() - interval '1 minute' * $1
		order by created asc`

	rows, err := r.pool.Query(ctx, query, olderThanMinutes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []matches.Match
	for rows.Next() {
		var match matches.Match
		if err := rows.Scan(&match.ID, &match.CreatorUserID, &match.Created, &match.Started,
			&match.CompletedAt, &match.CancelledAt, &match.Status, &match.Options); err != nil {
			return nil, err
		}
		result = append(result, match)
	}

	return result, rows.Err()
}

func (r *MatchTeamRepository) Create(ctx context.Context, team *matches.MatchTeam) error {
	const query = `
		insert into match_teams (id, match_id, password)
		values ($1, $2, $3)`

	_, err := r.pool.Exec(ctx, query, team.ID, team.MatchID, team.Password)
	return err
}

func (r *MatchTeamRepository) ListByMatch(ctx context.Context, matchID string) ([]matches.MatchTeam, error) {
	const query = `
		select id, match_id, password
		from match_teams
		where match_id = $1`

	rows, err := r.pool.Query(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []matches.MatchTeam
	for rows.Next() {
		var team matches.MatchTeam
		if err := rows.Scan(&team.ID, &team.MatchID, &team.Password); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, rows.Err()
}

func (r *MatchTeamRepository) GetByID(ctx context.Context, id string) (*matches.MatchTeam, error) {
	const query = `
		select id, match_id, password
		from match_teams
		where id = $1`

	var team matches.MatchTeam
	err := r.pool.QueryRow(ctx, query, id).Scan(&team.ID, &team.MatchID, &team.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *MatchRegistrationRepository) Upsert(ctx context.Context, registration *matches.MatchRegistration) error {
	const query = `
		insert into match_registrations (match_id, fighter_id, team_id, date)
		values ($1, $2, $3, $4)
		on conflict (match_id, fighter_id)
		do update set team_id = excluded.team_id, date = excluded.date`

	_, err := r.pool.Exec(ctx, query,
		registration.MatchID,
		registration.FighterID,
		registration.TeamID,
		registration.Date,
	)
	return err
}

func (r *MatchRegistrationRepository) Delete(ctx context.Context, matchID string, fighterID string) error {
	const query = `
		delete from match_registrations
		where match_id = $1 and fighter_id = $2`

	_, err := r.pool.Exec(ctx, query, matchID, fighterID)
	return err
}

func (r *MatchRegistrationRepository) GetByMatchAndFighter(ctx context.Context, matchID string, fighterID string) (*matches.MatchRegistration, error) {
	const query = `
		select match_id, fighter_id, team_id, date
		from match_registrations
		where match_id = $1 and fighter_id = $2`

	var registration matches.MatchRegistration
	err := r.pool.QueryRow(ctx, query, matchID, fighterID).Scan(
		&registration.MatchID, &registration.FighterID, &registration.TeamID, &registration.Date,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &registration, nil
}

func (r *MatchRegistrationRepository) CountByMatchAndUser(ctx context.Context, matchID string, userID int64) (int, error) {
	const query = `
		select count(*)
		from match_registrations mr
		join fighters f on f.id = mr.fighter_id
		where mr.match_id = $1 and f.user_id = $2`

	var count int
	err := r.pool.QueryRow(ctx, query, matchID, userID).Scan(&count)
	return count, err
}

func (r *MatchRegistrationRepository) ListByMatch(ctx context.Context, matchID string) ([]matches.MatchRegistration, error) {
	const query = `
		select match_id, fighter_id, team_id, date
		from match_registrations
		where match_id = $1
		order by date asc`

	rows, err := r.pool.Query(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []matches.MatchRegistration
	for rows.Next() {
		var reg matches.MatchRegistration
		if err := rows.Scan(&reg.MatchID, &reg.FighterID, &reg.TeamID, &reg.Date); err != nil {
			return nil, err
		}
		registrations = append(registrations, reg)
	}
	return registrations, rows.Err()
}

func (r *MatchResultRepository) GetByMatch(ctx context.Context, matchID string) (*matches.MatchResult, error) {
	const query = `
		select id, match_id, round_ticks
		from match_results
		where match_id = $1`

	var result matches.MatchResult
	err := r.pool.QueryRow(ctx, query, matchID).Scan(&result.ID, &result.MatchID, &result.RoundTicks)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *MatchResultRepository) Upsert(ctx context.Context, result *matches.MatchResult) error {
	const query = `
		insert into match_results (id, match_id, round_ticks)
		values ($1, $2, $3)
		on conflict (match_id)
		do update set round_ticks = excluded.round_ticks`

	_, err := r.pool.Exec(ctx, query, result.ID, result.MatchID, result.RoundTicks)
	return err
}

func (r *MatchScoreRepository) ListByMatch(ctx context.Context, matchID string) ([]matches.MatchScoreFighter, error) {
	const query = `
		select match_id, fighter_id, total_kills, total_deaths, total_assists
		from match_score_fighters
		where match_id = $1`

	rows, err := r.pool.Query(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []matches.MatchScoreFighter
	for rows.Next() {
		var score matches.MatchScoreFighter
		if err := rows.Scan(&score.MatchID, &score.FighterID, &score.TotalKills, &score.TotalDeaths, &score.TotalAssists); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	return scores, rows.Err()
}

func (r *MatchScoreRepository) Upsert(ctx context.Context, scores []matches.MatchScoreFighter) error {
	const query = `
		insert into match_score_fighters (match_id, fighter_id, total_kills, total_deaths, total_assists)
		values ($1, $2, $3, $4, $5)
		on conflict (match_id, fighter_id)
		do update set total_kills = excluded.total_kills,
					  total_deaths = excluded.total_deaths,
					  total_assists = excluded.total_assists`

	batch := &pgx.Batch{}
	for _, score := range scores {
		batch.Queue(query, score.MatchID, score.FighterID, score.TotalKills, score.TotalDeaths, score.TotalAssists)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	_, err := br.Exec()
	return err
}

package matches

import (
	"context"

	"empoweredpixels/internal/domain/matches"
	"empoweredpixels/internal/domain/roster"
)

type MatchRepository interface {
	Create(ctx context.Context, match *matches.Match) error
	GetByID(ctx context.Context, id string) (*matches.Match, error)
	Update(ctx context.Context, match *matches.Match) error
	ListOpen(ctx context.Context, limit int, offset int) ([]matches.Match, error)
	ListByStatus(ctx context.Context, status string, limit int, offset int) ([]matches.Match, error)
	GetCurrentMatch(ctx context.Context, userID int64) (*matches.Match, error)
	ListStaleLobbies(ctx context.Context, olderThanMinutes int) ([]matches.Match, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team *matches.MatchTeam) error
	ListByMatch(ctx context.Context, matchID string) ([]matches.MatchTeam, error)
	GetByID(ctx context.Context, id string) (*matches.MatchTeam, error)
}

type RegistrationRepository interface {
	Upsert(ctx context.Context, registration *matches.MatchRegistration) error
	Delete(ctx context.Context, matchID string, fighterID string) error
	GetByMatchAndFighter(ctx context.Context, matchID string, fighterID string) (*matches.MatchRegistration, error)
	CountByMatchAndUser(ctx context.Context, matchID string, userID int64) (int, error)
	ListByMatch(ctx context.Context, matchID string) ([]matches.MatchRegistration, error)
}

type ResultRepository interface {
	GetByMatch(ctx context.Context, matchID string) (*matches.MatchResult, error)
	Upsert(ctx context.Context, result *matches.MatchResult) error
}

type ScoreRepository interface {
	ListByMatch(ctx context.Context, matchID string) ([]matches.MatchScoreFighter, error)
	Upsert(ctx context.Context, scores []matches.MatchScoreFighter) error
}

type FighterRepository interface {
	GetByID(ctx context.Context, id string) (*roster.Fighter, error)
	GetByUserAndID(ctx context.Context, userID int64, id string) (*roster.Fighter, error)
	ListByUser(ctx context.Context, userID int64) ([]roster.Fighter, error)
	ListByMatch(ctx context.Context, matchID string) ([]roster.Fighter, error)
}

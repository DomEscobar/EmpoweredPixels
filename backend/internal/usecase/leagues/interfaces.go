package leagues

import (
	"context"

	"empoweredpixels/internal/domain/leagues"
	"empoweredpixels/internal/domain/roster"
)

type LeagueRepository interface {
	List(ctx context.Context) ([]leagues.League, error)
	GetByID(ctx context.Context, id int) (*leagues.League, error)
	Create(ctx context.Context, league *leagues.League) error
	Update(ctx context.Context, league *leagues.League) error
	Delete(ctx context.Context, id int) error
}

type SubscriptionRepository interface {
	ListByLeague(ctx context.Context, leagueID int) ([]leagues.LeagueSubscription, error)
	ListByLeagueAndUser(ctx context.Context, leagueID int, userID int64) ([]leagues.LeagueSubscription, error)
	Create(ctx context.Context, subscription *leagues.LeagueSubscription) error
	Delete(ctx context.Context, subscription *leagues.LeagueSubscription) error
}

type LeagueMatchRepository interface {
	ListByLeague(ctx context.Context, leagueID int, limit int, offset int) ([]leagues.LeagueMatch, error)
	CountByLeague(ctx context.Context, leagueID int) (int, error)
	GetLastWinner(ctx context.Context, leagueID int) (*leagues.LeagueWinner, error)
	GetHighScores(ctx context.Context, leagueID int, lastMatches int) ([]leagues.LeagueHighscore, error)
	Create(ctx context.Context, leagueID int, matchID string) error
	UpdateStarted(ctx context.Context, leagueID int, matchID string, started *time.Time) error
	Delete(ctx context.Context, leagueID int, matchID string) error
}

type FighterRepository interface {
	GetByUserAndID(ctx context.Context, userID int64, id string) (*roster.Fighter, error)
	GetByID(ctx context.Context, id string) (*roster.Fighter, error)
}

type AchievementRepository interface {
	UpdateAchievementProgress(ctx context.Context, userID int, achievementKey string, progress int) error
}

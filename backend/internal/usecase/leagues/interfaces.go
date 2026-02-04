package leagues

import (
	"context"

	"empoweredpixels/internal/domain/leagues"
	"empoweredpixels/internal/domain/roster"
)

type LeagueRepository interface {
	List(ctx context.Context) ([]leagues.League, error)
	GetByID(ctx context.Context, id int) (*leagues.League, error)
}

type SubscriptionRepository interface {
	ListByLeague(ctx context.Context, leagueID int) ([]leagues.LeagueSubscription, error)
	ListByLeagueAndUser(ctx context.Context, leagueID int, userID int64) ([]leagues.LeagueSubscription, error)
	Create(ctx context.Context, subscription *leagues.LeagueSubscription) error
	Delete(ctx context.Context, subscription *leagues.LeagueSubscription) error
}

type LeagueMatchRepository interface {
	ListByLeague(ctx context.Context, leagueID int, limit int, offset int) ([]leagues.LeagueMatch, error)
}

type FighterRepository interface {
	GetByUserAndID(ctx context.Context, userID int64, id string) (*roster.Fighter, error)
}

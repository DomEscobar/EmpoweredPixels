package leagues

import (
	"context"
	"errors"
	"time"

	"empoweredpixels/internal/domain/leagues"
)

var (
	ErrInvalidLeague       = errors.New("invalid league")
	ErrInvalidFighter      = errors.New("invalid fighter")
	ErrInvalidSubscription = errors.New("invalid league subscription")
)

type Service struct {
	leagues       LeagueRepository
	subscriptions SubscriptionRepository
	matches       LeagueMatchRepository
	fighters      FighterRepository
	now           func() time.Time
}

func NewService(
	leagues LeagueRepository,
	subscriptions SubscriptionRepository,
	matches LeagueMatchRepository,
	fighters FighterRepository,
	now func() time.Time,
) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		leagues:       leagues,
		subscriptions: subscriptions,
		matches:       matches,
		fighters:      fighters,
		now:           now,
	}
}

func (s *Service) List(ctx context.Context) ([]leagues.League, error) {
	return s.leagues.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int) (*leagues.League, error) {
	return s.leagues.GetByID(ctx, id)
}

func (s *Service) Subscribe(ctx context.Context, userID int64, leagueID int, fighterID string) error {
	league, err := s.leagues.GetByID(ctx, leagueID)
	if err != nil {
		return err
	}
	if league == nil {
		return ErrInvalidLeague
	}

	fighter, err := s.fighters.GetByUserAndID(ctx, userID, fighterID)
	if err != nil {
		return err
	}
	if fighter == nil {
		return ErrInvalidFighter
	}

	subscription := &leagues.LeagueSubscription{
		LeagueID:  leagueID,
		FighterID: fighterID,
		Created:   s.now(),
	}

	return s.subscriptions.Create(ctx, subscription)
}

func (s *Service) Unsubscribe(ctx context.Context, userID int64, leagueID int, fighterID string) error {
	league, err := s.leagues.GetByID(ctx, leagueID)
	if err != nil {
		return err
	}
	if league == nil {
		return ErrInvalidLeague
	}

	fighter, err := s.fighters.GetByUserAndID(ctx, userID, fighterID)
	if err != nil {
		return err
	}
	if fighter == nil {
		return ErrInvalidFighter
	}

	subscription := &leagues.LeagueSubscription{
		LeagueID:  leagueID,
		FighterID: fighterID,
	}

	if err := s.subscriptions.Delete(ctx, subscription); err != nil {
		return ErrInvalidSubscription
	}

	return nil
}

func (s *Service) Subscriptions(ctx context.Context, leagueID int) ([]leagues.LeagueSubscription, error) {
	return s.subscriptions.ListByLeague(ctx, leagueID)
}

func (s *Service) SubscriptionsForUser(ctx context.Context, leagueID int, userID int64) ([]leagues.LeagueSubscription, error) {
	return s.subscriptions.ListByLeagueAndUser(ctx, leagueID, userID)
}

func (s *Service) Matches(ctx context.Context, leagueID int, page int, pageSize int) ([]leagues.LeagueMatch, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.matches.ListByLeague(ctx, leagueID, pageSize, offset)
}

package leagues

import (
	"context"
	"encoding/json"
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

	if err := s.subscriptions.Create(ctx, subscription); err != nil {
		return err
	}

	return nil
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

func (s *Service) Matches(ctx context.Context, leagueID int, page int, pageSize int) ([]leagues.LeagueMatch, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	items, err := s.matches.ListByLeague(ctx, leagueID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	totalCount, err := s.matches.CountByLeague(ctx, leagueID)
	if err != nil {
		// If count fails, return items with totalCount = len(items) as fallback
		totalCount = len(items)
	}
	return items, totalCount, nil
}

func (s *Service) GetLastWinner(ctx context.Context, leagueID int) (*leagues.LeagueWinner, error) {
	return s.matches.GetLastWinner(ctx, leagueID)
}

func (s *Service) GetHighScores(ctx context.Context, leagueID int, lastMatches int) ([]leagues.LeagueHighscore, error) {
	return s.matches.GetHighScores(ctx, leagueID, lastMatches)
}

func (s *Service) CreateLeague(ctx context.Context, name string, options []byte, isDeactivated bool) (*leagues.League, error) {
	if name == "" {
		return nil, errors.New("league name is required")
	}
	// Validate options JSON if present
	if len(options) > 0 {
		var opts map[string]interface{}
		if err := json.Unmarshal(options, &opts); err != nil {
			return nil, errors.New("invalid options JSON")
		}
		// Validate tier if present
		if tier, ok := opts["tier"]; ok {
			validTiers := map[string]bool{
				"standard":  true,
				"epic":      true,
				"rare":      true,
				"legendary": true,
				"mythic":    true,
			}
			if tierStr, ok := tier.(string); !ok || !validTiers[tierStr] {
				return nil, errors.New("invalid tier")
			}
		}
	}

	league := &leagues.League{
		Name:          name,
		Options:       options,
		IsDeactivated: isDeactivated,
	}

	if err := s.leagues.Create(ctx, league); err != nil {
		return nil, err
	}

	return league, nil
}

func (s *Service) UpdateLeague(ctx context.Context, id int, name string, options []byte, isDeactivated bool) (*leagues.League, error) {
	// Check if league exists
	existing, err := s.leagues.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrInvalidLeague
	}

	if name == "" {
		return nil, errors.New("league name is required")
	}

	// Validate options JSON if present
	if len(options) > 0 {
		var opts map[string]interface{}
		if err := json.Unmarshal(options, &opts); err != nil {
			return nil, errors.New("invalid options JSON")
		}
		// Validate tier if present
		if tier, ok := opts["tier"]; ok {
			validTiers := map[string]bool{
				"standard":  true,
				"epic":      true,
				"rare":      true,
				"legendary": true,
				"mythic":    true,
			}
			if tierStr, ok := tier.(string); !ok || !validTiers[tierStr] {
				return nil, errors.New("invalid tier")
			}
		}
	}

	league := &leagues.League{
		ID:            id,
		Name:          name,
		Options:       options,
		IsDeactivated: isDeactivated,
	}

	if err := s.leagues.Update(ctx, league); err != nil {
		return nil, err
	}

	return league, nil
}

func (s *Service) DeleteLeague(ctx context.Context, id int) error {
	// Verify league exists
	existing, err := s.leagues.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrInvalidLeague
	}

	return s.leagues.Delete(ctx, id)
}

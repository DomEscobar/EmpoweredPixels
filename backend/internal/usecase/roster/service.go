package roster

import (
	"context"
	"errors"
	"time"

	"empoweredpixels/internal/domain/roster"
	"github.com/google/uuid"
)

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidFighter    = errors.New("invalid fighter")
	ErrFighterExists     = errors.New("fighter already exists")
	ErrFighterNameExists = errors.New("fighter name already exists")
)

type Service struct {
	fighters       FighterRepository
	experiences    ExperienceRepository
	configurations ConfigurationRepository
	now            func() time.Time
}

func NewService(
	fighters FighterRepository,
	experiences ExperienceRepository,
	configurations ConfigurationRepository,
	now func() time.Time,
) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		fighters:       fighters,
		experiences:    experiences,
		configurations: configurations,
		now:            now,
	}
}

func (s *Service) List(ctx context.Context, userID int64) ([]roster.Fighter, error) {
	return s.fighters.ListByUser(ctx, userID)
}

func (s *Service) Get(ctx context.Context, userID int64, id string) (*roster.Fighter, error) {
	return s.fighters.GetByUserAndID(ctx, userID, id)
}

func (s *Service) GetByID(ctx context.Context, id string) (*roster.Fighter, error) {
	return s.fighters.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, userID int64, name string) (*roster.Fighter, error) {
	if name == "" {
		return nil, ErrInvalidFighter
	}

	exists, err := s.fighters.NameExists(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrFighterNameExists
	}

	hasFighter, err := s.fighters.UserHasFighter(ctx, userID)
	if err != nil {
		return nil, err
	}
	if hasFighter {
		return nil, ErrFighterExists
	}

	fighter := &roster.Fighter{
		ID:      uuid.NewString(),
		UserID:  userID,
		Name:    name,
		Level:   1,
		Created: s.now(),
	}

	if err := s.fighters.Create(ctx, fighter); err != nil {
		return nil, err
	}

	return fighter, nil
}

func (s *Service) Delete(ctx context.Context, userID int64, id string) error {
	return s.fighters.SoftDelete(ctx, userID, id)
}

func (s *Service) GetExperience(ctx context.Context, fighterID string) (*roster.FighterExperience, error) {
	exp, err := s.experiences.GetByFighterID(ctx, fighterID)
	if err != nil {
		return nil, err
	}
	if exp == nil {
		return &roster.FighterExperience{
			FighterID:  fighterID,
			Experience: 0,
		}, nil
	}
	return exp, nil
}

func (s *Service) GetConfiguration(ctx context.Context, fighterID string) (*roster.FighterConfiguration, error) {
	config, err := s.configurations.GetByFighterID(ctx, fighterID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return &roster.FighterConfiguration{
			FighterID: fighterID,
		}, nil
	}
	return config, nil
}

func (s *Service) UpdateConfiguration(ctx context.Context, configuration *roster.FighterConfiguration) error {
	return s.configurations.Upsert(ctx, configuration)
}

func (s *Service) UpdateExperience(ctx context.Context, experience *roster.FighterExperience) error {
	return s.experiences.Upsert(ctx, experience)
}

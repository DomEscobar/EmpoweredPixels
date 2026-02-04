package roster

import (
	"context"

	"empoweredpixels/internal/domain/roster"
)

type FighterRepository interface {
	ListByUser(ctx context.Context, userID int64) ([]roster.Fighter, error)
	GetByUserAndID(ctx context.Context, userID int64, id string) (*roster.Fighter, error)
	GetByID(ctx context.Context, id string) (*roster.Fighter, error)
	NameExists(ctx context.Context, name string) (bool, error)
	UserHasFighter(ctx context.Context, userID int64) (bool, error)
	Create(ctx context.Context, fighter *roster.Fighter) error
	Update(ctx context.Context, fighter *roster.Fighter) error
	SoftDelete(ctx context.Context, userID int64, id string) error
}

type ExperienceRepository interface {
	GetByFighterID(ctx context.Context, fighterID string) (*roster.FighterExperience, error)
	Upsert(ctx context.Context, experience *roster.FighterExperience) error
}

type ConfigurationRepository interface {
	GetByFighterID(ctx context.Context, fighterID string) (*roster.FighterConfiguration, error)
	Upsert(ctx context.Context, configuration *roster.FighterConfiguration) error
}

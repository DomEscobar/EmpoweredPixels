package attunement

import (
	"context"

	"empoweredpixels/internal/domain/attunement"
)

type AttunementRepository interface {
	GetPlayerAttunements(ctx context.Context, userID int) ([]attunement.Attunement, error)
	GetAttunement(ctx context.Context, userID int, element attunement.Element) (*attunement.Attunement, error)
	AddXP(ctx context.Context, userID int, element attunement.Element, xpAmount int, source string) (levelUp bool, newLevel int, err error)
	CreateInitialAttunements(ctx context.Context, userID int) error
}

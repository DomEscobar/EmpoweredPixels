package seasons

import (
	"context"

	"empoweredpixels/internal/domain/seasons"
)

type SummaryRepository interface {
	ListByUser(ctx context.Context, userID int64, limit int, offset int) ([]seasons.SeasonSummary, error)
}

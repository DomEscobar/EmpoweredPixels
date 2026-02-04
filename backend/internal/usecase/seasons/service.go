package seasons

import (
	"context"

	"empoweredpixels/internal/domain/seasons"
)

type Service struct {
	summaries SummaryRepository
}

func NewService(summaries SummaryRepository) *Service {
	return &Service{summaries: summaries}
}

func (s *Service) SummaryPage(ctx context.Context, userID int64, page int, pageSize int) ([]seasons.SeasonSummary, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.summaries.ListByUser(ctx, userID, pageSize, offset)
}

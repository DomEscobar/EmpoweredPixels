package events

import (
	"context"
	"empoweredpixels/internal/domain/events"
	"empoweredpixels/internal/infra/db/repositories"
)

type Service struct {
	repo *repositories.EventRepository
}

func NewService(repo *repositories.EventRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCurrentEvents(ctx context.Context) ([]events.ActiveEvent, error) {
	return s.repo.GetActiveEvents(ctx)
}

func (s *Service) GetStatus(ctx context.Context) (*events.EventStatus, error) {
	return s.repo.GetStatus(ctx)
}

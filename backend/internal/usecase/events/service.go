package events

import (
	"context"
	"fmt"
	"time"

	"empoweredpixels/internal/domain/events"
	"empoweredpixels/internal/infra/db/repositories"
)

// Service handles event business logic
type Service struct {
	repo repositories.EventRepository
}

// NewService creates a new events service
func NewService(repo repositories.EventRepository) *Service {
	return &Service{repo: repo}
}

// GetCurrentEvents returns all currently active events
func (s *Service) GetCurrentEvents(ctx context.Context) ([]events.ActiveEvent, error) {
	return s.repo.GetActiveEvents(ctx)
}

// GetEventStatus returns the current event status
func (s *Service) GetEventStatus(ctx context.Context, userID int) (*events.EventStatus, error) {
	return s.repo.GetEventStatus(ctx, userID)
}

// CheckAndActivateEvents checks for events that should be active and activates them
func (s *Service) CheckAndActivateEvents(ctx context.Context) error {
	// Deactivate expired events first
	if err := s.repo.DeactivateExpiredEvents(ctx); err != nil {
		return fmt.Errorf("failed to deactivate expired events: %w", err)
	}

	// Get all configured events
	allEvents, err := s.repo.GetAllEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}

	now := time.Now()

	// Check each event
	for _, e := range allEvents {
		if e.IsEventActive(now) {
			// Calculate end time
			endTime := time.Date(now.Year(), now.Month(), now.Day(), e.EndHour, 0, 0, 0, now.Location())
			if now.Hour() > e.EndHour {
				endTime = endTime.AddDate(0, 0, 1)
			}

			// Activate event
			if err := s.repo.CreateActiveEvent(ctx, e.ID, endTime); err != nil {
				continue // Don't fail on single event error
			}
		}
	}

	return nil
}

// GetEventMultiplier returns the current multiplier for a specific event type
func (s *Service) GetEventMultiplier(ctx context.Context, eventType string) float64 {
	activeEvents, err := s.repo.GetActiveEvents(ctx)
	if err != nil {
		return 1.0
	}

	for _, ae := range activeEvents {
		if ae.Event != nil && ae.Event.EventType == eventType {
			return ae.Event.Multiplier
		}
	}

	return 1.0
}

// GetNextEventInfo returns information about the next upcoming event
func (s *Service) GetNextEventInfo(ctx context.Context) (*events.WeekendEvent, time.Duration, error) {
	allEvents, err := s.repo.GetAllEvents(ctx)
	if err != nil {
		return nil, 0, err
	}

	now := time.Now()
	var nextEvent *events.WeekendEvent
	var shortestWait time.Duration = -1

	for i := range allEvents {
		e := &allEvents[i]
		wait := e.GetTimeUntilNextEvent(now)
		
		if shortestWait < 0 || wait < shortestWait {
			shortestWait = wait
			nextEvent = e
		}
	}

	return nextEvent, shortestWait, nil
}

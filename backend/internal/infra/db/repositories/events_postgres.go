package repositories

import (
	"context"
	"fmt"
	"time"

	"empoweredpixels/internal/domain/events"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EventRepository defines event operations
type EventRepository interface {
	GetAllEvents(ctx context.Context) ([]events.WeekendEvent, error)
	GetActiveEvents(ctx context.Context) ([]events.ActiveEvent, error)
	CreateActiveEvent(ctx context.Context, eventID string, endsAt time.Time) error
	DeactivateExpiredEvents(ctx context.Context) error
	GetEventStatus(ctx context.Context, userID int) (*events.EventStatus, error)
}

// EventPostgres implements EventRepository
type EventPostgres struct {
	db *pgxpool.Pool
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *pgxpool.Pool) EventRepository {
	return &EventPostgres{db: db}
}

// GetAllEvents retrieves all weekend events
func (r *EventPostgres) GetAllEvents(ctx context.Context) ([]events.WeekendEvent, error) {
	query := `
		SELECT id, name, description, event_type, multiplier, start_day, end_day, start_hour, end_hour, is_active, created_at
		FROM weekend_events
		WHERE is_active = true
		ORDER BY start_day, start_hour
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()

	var eventList []events.WeekendEvent
	for rows.Next() {
		var e events.WeekendEvent
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.EventType, &e.Multiplier, 
			&e.StartDay, &e.EndDay, &e.StartHour, &e.EndHour, &e.IsActive, &e.CreatedAt); err != nil {
			return nil, err
		}
		eventList = append(eventList, e)
	}

	return eventList, rows.Err()
}

// GetActiveEvents retrieves currently active events
func (r *EventPostgres) GetActiveEvents(ctx context.Context) ([]events.ActiveEvent, error) {
	query := `
		SELECT ae.id, ae.event_id, ae.started_at, ae.ends_at, ae.is_active,
		       e.id, e.name, e.description, e.event_type, e.multiplier, e.start_day, e.end_day, e.start_hour, e.end_hour, e.is_active, e.created_at
		FROM active_events ae
		JOIN weekend_events e ON e.id = ae.event_id
		WHERE ae.is_active = true AND ae.ends_at > NOW()
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activeEvents []events.ActiveEvent
	for rows.Next() {
		var ae events.ActiveEvent
		var e events.WeekendEvent
		if err := rows.Scan(&ae.ID, &ae.EventID, &ae.StartedAt, &ae.EndsAt, &ae.IsActive,
			&e.ID, &e.Name, &e.Description, &e.EventType, &e.Multiplier, &e.StartDay, &e.EndDay, &e.StartHour, &e.EndHour, &e.IsActive, &e.CreatedAt); err != nil {
			return nil, err
		}
		ae.Event = &e
		activeEvents = append(activeEvents, ae)
	}

	return activeEvents, rows.Err()
}

// CreateActiveEvent creates a new active event
func (r *EventPostgres) CreateActiveEvent(ctx context.Context, eventID string, endsAt time.Time) error {
	query := `
		INSERT INTO active_events (event_id, started_at, ends_at, is_active)
		VALUES ($1, NOW(), $2, true)
		ON CONFLICT DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, eventID, endsAt)
	return err
}

// DeactivateExpiredEvents marks expired events as inactive
func (r *EventPostgres) DeactivateExpiredEvents(ctx context.Context) error {
	query := `
		UPDATE active_events
		SET is_active = false
		WHERE is_active = true AND ends_at <= NOW()
	`

	_, err := r.db.Exec(ctx, query)
	return err
}

// GetEventStatus gets the current event status for a user
func (r *EventPostgres) GetEventStatus(ctx context.Context, userID int) (*events.EventStatus, error) {
	// Check for active events
	activeEvents, err := r.GetActiveEvents(ctx)
	if err != nil {
		return nil, err
	}

	if len(activeEvents) == 0 {
		return &events.EventStatus{
			HasActiveEvent: false,
			Multiplier:     1.0,
		}, nil
	}

	// Return the first active event (could be extended for multiple events)
	ae := activeEvents[0]
	timeRemaining := ae.EndsAt.Sub(time.Now())

	return &events.EventStatus{
		HasActiveEvent: true,
		ActiveEvent:    &ae,
		TimeRemaining:  formatDurationEvent(timeRemaining),
		Multiplier:     ae.Event.Multiplier,
		Type:           ae.Event.EventType,
	}, nil
}

// formatDurationEvent formats duration as HH:MM:SS
func formatDurationEvent(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

package repositories

import (
	"context"
	"fmt"
	"time"

	"empoweredpixels/internal/domain/events"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetActiveEvents(ctx context.Context) ([]events.ActiveEvent, error) {
	query := `
		SELECT ae.id, ae.event_id, ae.started_at, ae.ends_at, ae.is_active,
		       e.id, e.name, e.description, e.event_type, e.multiplier
		FROM active_events ae
		JOIN weekend_events e ON e.id = ae.event_id
		WHERE ae.is_active = true AND ae.ends_at > NOW()
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []events.ActiveEvent
	for rows.Next() {
		var ae events.ActiveEvent
		var e events.WeekendEvent
		if err := rows.Scan(&ae.ID, &ae.EventID, &ae.StartedAt, &ae.EndsAt, &ae.IsActive,
			&e.ID, &e.Name, &e.Description, &e.EventType, &e.Multiplier); err != nil {
			return nil, err
		}
		ae.Event = &e
		result = append(result, ae)
	}
	return result, nil
}

func (r *EventRepository) GetStatus(ctx context.Context) (*events.EventStatus, error) {
	active, err := r.GetActiveEvents(ctx)
	if err != nil {
		return nil, err
	}

	status := &events.EventStatus{
		HasActiveEvent: false,
		Multiplier:     1.0,
	}

	if len(active) > 0 {
		ae := active[0]
		status.HasActiveEvent = true
		status.ActiveEvent = &ae
		status.Multiplier = ae.Event.Multiplier
		status.Type = ae.Event.EventType
		
		diff := ae.EndsAt.Sub(time.Now())
		status.TimeRemaining = fmt.Sprintf("%02d:%02d:%02d", int(diff.Hours()), int(diff.Minutes())%60, int(diff.Seconds())%60)
	}

	return status, nil
}

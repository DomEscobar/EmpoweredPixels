package events

import "time"

type WeekendEvent struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	EventType   string    `json:"event_type" db:"event_type"`
	Multiplier  float64   `json:"multiplier" db:"multiplier"`
	StartDay    int       `json:"start_day" db:"start_day"`
	EndDay      int       `json:"end_day" db:"end_day"`
	StartHour   int       `json:"start_hour" db:"start_hour"`
	EndHour     int       `json:"end_hour" db:"end_hour"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ActiveEvent struct {
	ID          string          `json:"id" db:"id"`
	EventID     string          `json:"event_id" db:"event_id"`
	Event       *WeekendEvent   `json:"event,omitempty"`
	StartedAt   time.Time       `json:"started_at" db:"started_at"`
	EndsAt      time.Time       `json:"ends_at" db:"ends_at"`
	IsActive    bool            `json:"is_active" db:"is_active"`
}

type EventStatus struct {
	HasActiveEvent bool           `json:"has_active_event"`
	ActiveEvent    *ActiveEvent   `json:"active_event,omitempty"`
	TimeRemaining  string         `json:"time_remaining,omitempty"`
	Multiplier     float64        `json:"multiplier"`
	Type           string         `json:"type,omitempty"`
}

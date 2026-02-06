package events

import "time"

// WeekendEvent represents an event configuration
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

// ActiveEvent represents a currently running event
type ActiveEvent struct {
	ID          string          `json:"id" db:"id"`
	EventID     string          `json:"event_id" db:"event_id"`
	Event       *WeekendEvent   `json:"event,omitempty"`
	StartedAt   time.Time       `json:"started_at" db:"started_at"`
	EndsAt      time.Time       `json:"ends_at" db:"ends_at"`
	IsActive    bool            `json:"is_active" db:"is_active"`
}

// EventStatus represents the current event status for a user
type EventStatus struct {
	HasActiveEvent bool           `json:"has_active_event"`
	ActiveEvent    *ActiveEvent   `json:"active_event,omitempty"`
	TimeRemaining  string         `json:"time_remaining,omitempty"`
	Multiplier     float64        `json:"multiplier"`
	Type           string         `json:"type,omitempty"`
}

// Event types
const (
	EventTypeDoubleDrops = "double_drops"
	EventTypeDoubleXP    = "double_xp"
	EventTypeBonusGold   = "bonus_gold"
	EventTypeHalfPrice   = "half_price"
)

// IsEventActive checks if an event should be active based on current time
func (e *WeekendEvent) IsEventActive(now time.Time) bool {
	if !e.IsActive {
		return false
	}

	currentDay := int(now.Weekday())
	currentHour := now.Hour()

	// Check if current day is within event days
	if e.StartDay <= e.EndDay {
		if currentDay < e.StartDay || currentDay > e.EndDay {
			return false
		}
	} else {
		// Event spans across week boundary (e.g., Sat to Sun)
		if currentDay < e.StartDay && currentDay > e.EndDay {
			return false
		}
	}

	// Check hour
	if currentHour < e.StartHour || currentHour > e.EndHour {
		return false
	}

	return true
}

// GetTimeUntilNextEvent calculates time until next event starts
func (e *WeekendEvent) GetTimeUntilNextEvent(now time.Time) time.Duration {
	currentDay := int(now.Weekday())
	
	// Find next occurrence of start day
	daysUntil := (e.StartDay - currentDay + 7) % 7
	if daysUntil == 0 && now.Hour() >= e.StartHour {
		daysUntil = 7 // Next week
	}

	nextEvent := now.AddDate(0, 0, daysUntil)
	nextEvent = time.Date(nextEvent.Year(), nextEvent.Month(), nextEvent.Day(), 
		e.StartHour, 0, 0, 0, nextEvent.Location())

	return nextEvent.Sub(now)
}

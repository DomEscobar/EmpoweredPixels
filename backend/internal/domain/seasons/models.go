package seasons

import "time"

type Season struct {
	ID        int64
	SeasonID  int
	StartDate time.Time
	EndDate   time.Time
}

type SeasonSummary struct {
	ID       int64
	UserID   int64
	SeasonID int
	Position int
}

type SeasonProgress struct {
	ID         int64
	UserID     int64
	SeasonID   int
	IsComplete bool
}

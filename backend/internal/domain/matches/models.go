package matches

import "time"

const (
	MatchStatusLobby    = "lobby"
	MatchStatusRunning  = "running"
	MatchStatusCompleted = "completed"
	MatchStatusCancelled = "cancelled"
)

type Match struct {
	ID            string
	CreatorUserID *int64
	Created       time.Time
	Started       *time.Time
	CompletedAt   *time.Time
	CancelledAt   *time.Time
	Status        string
	Options       []byte
}

type MatchTeam struct {
	ID       string
	MatchID  string
	Password *string
}

type MatchRegistration struct {
	MatchID   string
	FighterID string
	TeamID    *string
	Date      time.Time
}

type MatchResult struct {
	ID         string
	MatchID    string
	RoundTicks []byte
}

type MatchScoreFighter struct {
	MatchID      string
	FighterID    string
	TotalKills   int
	TotalDeaths  int
	TotalAssists int
}

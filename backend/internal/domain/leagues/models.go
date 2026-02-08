package leagues

import "time"

type League struct {
	ID            int
	Name          string
	Options       []byte
	IsDeactivated bool
}

type LeagueSubscription struct {
	LeagueID  int
	FighterID string
	Created   time.Time
}

type LeagueMatch struct {
	LeagueID int
	MatchID  string
	Started  *time.Time
}

type LeagueHighscore struct {
	FighterID   string
	FighterName string
	Username    string
	Score       int
}

type LeagueWinner struct {
	FighterID   string
	FighterName string
	Username    string
	MatchID     string
}

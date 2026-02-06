package roster

import "time"

type Fighter struct {
	ID             string
	UserID         int64
	Name           string
	Level          int
	XP             int
	XPToNextLevel  int
	Power          int
	ConditionPower int
	Precision      int
	Ferocity       int
	Accuracy       int
	Agility        int
	Armor          int
	Vitality       int
	ParryChance    int
	HealingPower   int
	Speed          int
	Vision         int
	WeaponID       *string
	AttunementID   *string
	// Match Statistics
	MatchesWon     int
	MatchesLost    int
	TotalMatches   int
	TotalDamageDealt int64
	TotalDamageTaken int64
	Created        time.Time
	IsDeleted      bool
}

type FighterExperience struct {
	ID         int64
	FighterID  string
	Experience int
}

type FighterConfiguration struct {
	FighterID    string
	AttunementID *string
}

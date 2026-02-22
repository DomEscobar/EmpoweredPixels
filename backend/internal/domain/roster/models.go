package roster

import "time"

type Fighter struct {
	ID             string  `json:"id"`
	UserID         int64   `json:"userId"`
	Name           string  `json:"name"`
	Level          int     `json:"level"`
	XP             int     `json:"xp"`
	XPToNextLevel  int     `json:"xpToNextLevel"`
	Power          int     `json:"power"`
	ConditionPower int     `json:"conditionPower"`
	Precision      int     `json:"precision"`
	Ferocity       int     `json:"ferocity"`
	Accuracy       int     `json:"accuracy"`
	Agility        int     `json:"agility"`
	Armor          int     `json:"armor"`
	Vitality       int     `json:"vitality"`
	ParryChance    int     `json:"parryChance"`
	HealingPower   int     `json:"healingPower"`
	Speed          int     `json:"speed"`
	Vision         int     `json:"vision"`
	WeaponID       *string `json:"weaponId"`
	TeamID         *string `json:"teamId"` // Transient field for battle team tracking
	// Match Statistics
	MatchesWon       int       `json:"matchesWon"`
	MatchesLost      int       `json:"matchesLost"`
	TotalMatches     int       `json:"totalMatches"`
	TotalDamageDealt int64     `json:"totalDamageDealt"`
	TotalDamageTaken int64     `json:"totalDamageTaken"`
	Created          time.Time `json:"created"`
	IsDeleted        bool      `json:"isDeleted"`
}

type FighterExperience struct {
	ID         int64  `json:"id"`
	FighterID  string `json:"fighterId"`
	Experience int    `json:"experience"`
}

type FighterConfiguration struct {
	FighterID string `json:"fighterId"`
}

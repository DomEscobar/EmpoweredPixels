package combat

import (
	"encoding/json"
)

type Entity struct {
	ID           string
	Name         string
	Level        int
	MaxHP        int
	CurrentHP    int
	X            float64
	Y            float64
	TeamID       *string
	AttunementID *string
	Stats        Stats
}

type Stats struct {
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
}

type MatchResult struct {
	MatchID    string
	RoundTicks []RoundTick
	Scores     []FighterScore
}

type RoundTick struct {
	Round int    `json:"round"`
	Ticks []Tick `json:"ticks"`
}

type Tick struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type FighterScore struct {
	FighterID string
	Kills     int
	Deaths    int
	Assists   int
}

type EventSpawn struct {
	FighterID string  `json:"fighterId"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	HP        int     `json:"hp"`
}

type EventMove struct {
	FighterID string  `json:"fighterId"`
	FromX     float64 `json:"fromX"`
	FromY     float64 `json:"fromY"`
	ToX       float64 `json:"toX"`
	ToY       float64 `json:"toY"`
}

type EventAttack struct {
	AttackerID string `json:"attackerId"`
	TargetID   string `json:"targetId"`
	SkillID    string `json:"skillId"`
	Damage     int    `json:"damage"`
	IsCritical bool   `json:"isCritical"`
	IsParried  bool   `json:"isParried"`
}

type EventHeal struct {
	HealerID string `json:"healerId"`
	TargetID string `json:"targetId"`
	Amount   int    `json:"amount"`
}

type EventDied struct {
	FighterID string `json:"fighterId"`
	KillerID  string `json:"killerId"`
}

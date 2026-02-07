package roster

import (
	"time"
)

type Squad struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"userId"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	Members   []Member  `json:"members"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Member struct {
	FighterID string `json:"fighterId"`
	SlotIndex int    `json:"slotIndex"`
}

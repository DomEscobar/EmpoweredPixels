package inventory

import "time"

type Item struct {
	ID      string
	UserID  int64
	ItemID  string
	Rarity  int
	Created time.Time
}

type Equipment struct {
	ID          string
	UserID      int64
	FighterID   *string
	ItemID      string
	Level       int
	Rarity      int
	Enhancement int
	Created     time.Time
}

type EquipmentOption struct {
	EquipmentID string
	IsFavorite  bool
}

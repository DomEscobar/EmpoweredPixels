package shop

import (
	"time"
)

// Rarity constants
const (
	RarityBasic     = 0
	RarityCommon    = 1
	RarityRare      = 2
	RarityFabled    = 3
	RarityMythic    = 4
	RarityLegendary = 5
)

// Item type constants
const (
	ItemTypeGoldPackage = "gold_package"
	ItemTypeBundle      = "bundle"
	ItemTypeEquipment   = "equipment"
	ItemTypeConsumable  = "consumable"
)

// Currency constants
const (
	CurrencyUSD       = "usd"
	CurrencyGold      = "gold"
	CurrencyParticles = "particles"
)

// Transaction status constants
const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusRefunded  = "refunded"
)

// Shop represents a shop in the game
type Shop struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Currency    string    `json:"currency"`
	IsActive    bool      `json:"is_active"`
	Created     time.Time `json:"created"`
}

// ShopItem represents an item available for purchase
type ShopItem struct {
	ID            string                 `json:"id"`
	ShopID        string                 `json:"shop_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	ItemType      string                 `json:"item_type"`
	PriceAmount   int                    `json:"price_amount"`
	PriceCurrency string                 `json:"price_currency"`
	GoldAmount    *int                   `json:"gold_amount,omitempty"`
	Rarity        int                    `json:"rarity"`
	ImageURL      *string                `json:"image_url,omitempty"`
	Metadata      map[string]interface{} `json:"metadata"`
	IsActive      bool                   `json:"is_active"`
	SortOrder     int                    `json:"sort_order"`
	Created       time.Time              `json:"created"`
}

// PlayerGold represents a player's gold balance
type PlayerGold struct {
	UserID         int64     `json:"user_id"`
	Balance        int       `json:"balance"`
	LifetimeEarned int       `json:"lifetime_earned"`
	LifetimeSpent  int       `json:"lifetime_spent"`
	Updated        time.Time `json:"updated"`
}

// Transaction represents a purchase transaction
type Transaction struct {
	ID            string                 `json:"id"`
	UserID        int64                  `json:"user_id"`
	ShopItemID    *string                `json:"shop_item_id,omitempty"`
	ItemType      string                 `json:"item_type"`
	ItemName      string                 `json:"item_name"`
	PriceAmount   int                    `json:"price_amount"`
	PriceCurrency string                 `json:"price_currency"`
	GoldChange    int                    `json:"gold_change"`
	Status        string                 `json:"status"`
	Metadata      map[string]interface{} `json:"metadata"`
	Created       time.Time              `json:"created"`
}

// PurchaseRequest represents a purchase request
type PurchaseRequest struct {
	ItemID string `json:"item_id"`
}

// PurchaseResult represents the result of a purchase
type PurchaseResult struct {
	Transaction   *Transaction `json:"transaction"`
	NewBalance    int          `json:"new_balance"`
	ItemsReceived []string     `json:"items_received,omitempty"`
}

// RarityName returns the display name for a rarity level
func RarityName(rarity int) string {
	switch rarity {
	case RarityBasic:
		return "Basic"
	case RarityCommon:
		return "Common"
	case RarityRare:
		return "Rare"
	case RarityFabled:
		return "Fabled"
	case RarityMythic:
		return "Mythic"
	case RarityLegendary:
		return "Legendary"
	default:
		return "Unknown"
	}
}

// RarityColor returns the CSS color class for a rarity level
func RarityColor(rarity int) string {
	switch rarity {
	case RarityBasic:
		return "gray"
	case RarityCommon:
		return "green"
	case RarityRare:
		return "blue"
	case RarityFabled:
		return "purple"
	case RarityMythic:
		return "orange"
	case RarityLegendary:
		return "yellow"
	default:
		return "gray"
	}
}

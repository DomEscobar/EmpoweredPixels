package shop

import (
	"time"
)

// Shop represents a shop category/type
type Shop struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	ShopType    string    `json:"shop_type" db:"shop_type"` // gold, bundles, premium
	Currency    string    `json:"currency" db:"currency"`   // usd, gold
	IsActive    bool      `json:"is_active" db:"is_active"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	Created     time.Time `json:"created" db:"created"`
}

// ShopItem represents a product in the shop
type ShopItem struct {
	ID            int                    `json:"id" db:"id"`
	ShopID        int                    `json:"shop_id" db:"shop_id"`
	Name          string                 `json:"name" db:"name"`
	Description   string                 `json:"description" db:"description"`
	ItemType      string                 `json:"item_type" db:"item_type"` // gold_package, bundle, equipment
	PriceAmount   int                    `json:"price_amount" db:"price_amount"`
	PriceCurrency string                 `json:"price_currency" db:"price_currency"` // usd, gold, particles
	GoldAmount    *int                   `json:"gold_amount,omitempty" db:"gold_amount"`
	Rarity        int                    `json:"rarity" db:"rarity"` // 0=Basic, 1=Common, 2=Rare, 3=Fabled, 4=Mythic, 5=Legendary
	ImageURL      *string                `json:"image_url,omitempty" db:"image_url"`
	Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
	IsActive      bool                   `json:"is_active" db:"is_active"`
	SortOrder     int                    `json:"sort_order" db:"sort_order"`
	Created       time.Time              `json:"created" db:"created"`
}

// PlayerGold represents a player's gold balance
type PlayerGold struct {
	UserID         int       `json:"user_id" db:"user_id"`
	Balance        int       `json:"balance" db:"balance"`
	LifetimeEarned int       `json:"lifetime_earned" db:"lifetime_earned"`
	LifetimeSpent  int       `json:"lifetime_spent" db:"lifetime_spent"`
	Updated        time.Time `json:"updated" db:"updated"`
}

// Transaction represents a purchase or gold change
type Transaction struct {
	ID            int                    `json:"id" db:"id"`
	UserID        int                    `json:"user_id" db:"user_id"`
	ShopItemID    *int                   `json:"shop_item_id,omitempty" db:"shop_item_id"`
	ItemType      string                 `json:"item_type" db:"item_type"`
	ItemName      string                 `json:"item_name" db:"item_name"`
	PriceAmount   int                    `json:"price_amount" db:"price_amount"`
	PriceCurrency string                 `json:"price_currency" db:"price_currency"`
	GoldChange    int                    `json:"gold_change" db:"gold_change"` // negative for purchases, positive for rewards
	Status        string                 `json:"status" db:"status"`           // pending, completed, failed, refunded
	Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
	Created       time.Time              `json:"created" db:"created"`
}

// PurchaseRequest represents a purchase attempt
type PurchaseRequest struct {
	ItemID int `json:"item_id"`
}

// PurchaseResponse represents the result of a purchase
type PurchaseResponse struct {
	Success        bool   `json:"success"`
	TransactionID  int    `json:"transaction_id"`
	NewBalance     int    `json:"new_balance"`
	ItemsReceived  []string `json:"items_received,omitempty"`
	Message        string `json:"message"`
}

// GoldPackage represents predefined gold packages
type GoldPackage struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	GoldAmount  int    `json:"gold_amount"`
	PriceCents  int    `json:"price_cents"` // Price in USD cents
	BonusPercent int   `json:"bonus_percent"`
}

// Rarity constants for shop items
const (
	RarityBasic     = 0
	RarityCommon    = 1
	RarityRare      = 2
	RarityFabled    = 3
	RarityMythic    = 4
	RarityLegendary = 5
)

// ItemType constants
const (
	ItemTypeGoldPackage = "gold_package"
	ItemTypeBundle      = "bundle"
	ItemTypeEquipment   = "equipment"
	ItemTypeConsumable  = "consumable"
)

// Currency constants
const (
	CurrencyUSD      = "usd"
	CurrencyGold     = "gold"
	CurrencyParticles = "particles"
)

// ShopType constants
const (
	ShopTypeGold    = "gold"
	ShopTypeBundles = "bundles"
	ShopTypePremium = "premium"
)

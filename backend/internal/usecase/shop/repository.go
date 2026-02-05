package shop

import (
	"context"

	"empoweredpixels/internal/domain/shop"
)

// ShopRepository defines the interface for shop data access
type ShopRepository interface {
	GetShop(ctx context.Context, id string) (*shop.Shop, error)
	ListShops(ctx context.Context) ([]shop.Shop, error)
}

// ShopItemRepository defines the interface for shop item data access
type ShopItemRepository interface {
	GetShopItem(ctx context.Context, id string) (*shop.ShopItem, error)
	ListShopItems(ctx context.Context, shopID string) ([]shop.ShopItem, error)
	ListAllActiveItems(ctx context.Context) ([]shop.ShopItem, error)
	ListItemsByType(ctx context.Context, itemType string) ([]shop.ShopItem, error)
}

// PlayerGoldRepository defines the interface for player gold data access
type PlayerGoldRepository interface {
	GetPlayerGold(ctx context.Context, userID int64) (*shop.PlayerGold, error)
	AddGold(ctx context.Context, userID int64, amount int) error
	SpendGold(ctx context.Context, userID int64, amount int) error
	SetGold(ctx context.Context, userID int64, balance int) error
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *shop.Transaction) error
	GetTransaction(ctx context.Context, id string) (*shop.Transaction, error)
	ListTransactions(ctx context.Context, userID int64, limit, offset int) ([]shop.Transaction, error)
}

// Repository combines all shop-related repositories
type Repository interface {
	ShopRepository
	ShopItemRepository
	PlayerGoldRepository
	TransactionRepository
}

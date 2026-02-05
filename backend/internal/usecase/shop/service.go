package shop

import (
	"context"
	"errors"
	"time"

	"empoweredpixels/internal/domain/shop"

	"github.com/google/uuid"
)

var (
	ErrItemNotFound       = errors.New("item not found")
	ErrItemNotActive      = errors.New("item not available for purchase")
	ErrInsufficientGold   = errors.New("insufficient gold balance")
	ErrInvalidPurchase    = errors.New("invalid purchase request")
	ErrTransactionFailed  = errors.New("transaction failed")
)

// Service handles shop business logic
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService creates a new shop service
func NewService(repo Repository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{repo: repo, now: now}
}

// GetPlayerGold returns the player's gold balance
func (s *Service) GetPlayerGold(ctx context.Context, userID int64) (*shop.PlayerGold, error) {
	return s.repo.GetPlayerGold(ctx, userID)
}

// ListShops returns all active shops
func (s *Service) ListShops(ctx context.Context) ([]shop.Shop, error) {
	return s.repo.ListShops(ctx)
}

// ListShopItems returns all items for a specific shop
func (s *Service) ListShopItems(ctx context.Context, shopID string) ([]shop.ShopItem, error) {
	return s.repo.ListShopItems(ctx, shopID)
}

// ListAllItems returns all active items across all shops
func (s *Service) ListAllItems(ctx context.Context) ([]shop.ShopItem, error) {
	return s.repo.ListAllActiveItems(ctx)
}

// ListGoldPackages returns all gold package items
func (s *Service) ListGoldPackages(ctx context.Context) ([]shop.ShopItem, error) {
	return s.repo.ListItemsByType(ctx, shop.ItemTypeGoldPackage)
}

// ListBundles returns all bundle items
func (s *Service) ListBundles(ctx context.Context) ([]shop.ShopItem, error) {
	return s.repo.ListItemsByType(ctx, shop.ItemTypeBundle)
}

// GetShopItem returns a specific shop item
func (s *Service) GetShopItem(ctx context.Context, itemID string) (*shop.ShopItem, error) {
	item, err := s.repo.GetShopItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}
	return item, nil
}

// Purchase processes a purchase request
func (s *Service) Purchase(ctx context.Context, userID int64, itemID string) (*shop.PurchaseResult, error) {
	// Get the item
	item, err := s.repo.GetShopItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}
	if !item.IsActive {
		return nil, ErrItemNotActive
	}

	// Handle purchase based on currency type
	switch item.PriceCurrency {
	case shop.CurrencyGold:
		return s.purchaseWithGold(ctx, userID, item)
	case shop.CurrencyUSD:
		return s.purchaseWithUSD(ctx, userID, item)
	default:
		return nil, ErrInvalidPurchase
	}
}

// purchaseWithGold handles purchases made with in-game gold
func (s *Service) purchaseWithGold(ctx context.Context, userID int64, item *shop.ShopItem) (*shop.PurchaseResult, error) {
	// Check balance
	playerGold, err := s.repo.GetPlayerGold(ctx, userID)
	if err != nil {
		return nil, err
	}
	if playerGold.Balance < item.PriceAmount {
		return nil, ErrInsufficientGold
	}

	// Deduct gold
	if err := s.repo.SpendGold(ctx, userID, item.PriceAmount); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := &shop.Transaction{
		ID:            uuid.NewString(),
		UserID:        userID,
		ShopItemID:    &item.ID,
		ItemType:      item.ItemType,
		ItemName:      item.Name,
		PriceAmount:   item.PriceAmount,
		PriceCurrency: item.PriceCurrency,
		GoldChange:    -item.PriceAmount,
		Status:        shop.StatusCompleted,
		Metadata:      make(map[string]interface{}),
		Created:       s.now(),
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		// Try to refund
		_ = s.repo.AddGold(ctx, userID, item.PriceAmount)
		return nil, ErrTransactionFailed
	}

	// Get new balance
	newBalance, err := s.repo.GetPlayerGold(ctx, userID)
	if err != nil {
		newBalance = &shop.PlayerGold{Balance: playerGold.Balance - item.PriceAmount}
	}

	// TODO: Grant bundle items if applicable

	return &shop.PurchaseResult{
		Transaction: tx,
		NewBalance:  newBalance.Balance,
	}, nil
}

// purchaseWithUSD handles purchases made with real money (gold packages)
func (s *Service) purchaseWithUSD(ctx context.Context, userID int64, item *shop.ShopItem) (*shop.PurchaseResult, error) {
	// For MVP, we simulate successful USD purchase
	// In production, this would integrate with Stripe/PayPal/etc.

	if item.GoldAmount == nil || *item.GoldAmount <= 0 {
		return nil, ErrInvalidPurchase
	}

	// Add gold to player
	if err := s.repo.AddGold(ctx, userID, *item.GoldAmount); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := &shop.Transaction{
		ID:            uuid.NewString(),
		UserID:        userID,
		ShopItemID:    &item.ID,
		ItemType:      item.ItemType,
		ItemName:      item.Name,
		PriceAmount:   item.PriceAmount,
		PriceCurrency: item.PriceCurrency,
		GoldChange:    *item.GoldAmount,
		Status:        shop.StatusCompleted,
		Metadata:      map[string]interface{}{"simulated": true},
		Created:       s.now(),
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		// Try to refund
		_ = s.repo.SpendGold(ctx, userID, *item.GoldAmount)
		return nil, ErrTransactionFailed
	}

	// Get new balance
	newBalance, err := s.repo.GetPlayerGold(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &shop.PurchaseResult{
		Transaction: tx,
		NewBalance:  newBalance.Balance,
	}, nil
}

// AddGold adds gold to a player's balance (admin function or reward)
func (s *Service) AddGold(ctx context.Context, userID int64, amount int) error {
	return s.repo.AddGold(ctx, userID, amount)
}

// GetTransactionHistory returns a player's transaction history
func (s *Service) GetTransactionHistory(ctx context.Context, userID int64, page, pageSize int) ([]shop.Transaction, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.ListTransactions(ctx, userID, pageSize, offset)
}

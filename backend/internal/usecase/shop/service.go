package shop

import (
	"context"
	"fmt"

	"empoweredpixels/internal/domain/shop"
	"empoweredpixels/internal/infra/db/repositories"
)

// Service handles shop business logic
type Service struct {
	shopRepo        repositories.ShopRepository
	goldRepo        repositories.PlayerGoldRepository
	transactionRepo repositories.TransactionRepository
}

// NewService creates a new shop service
func NewService(
	shopRepo repositories.ShopRepository,
	goldRepo repositories.PlayerGoldRepository,
	transactionRepo repositories.TransactionRepository,
) *Service {
	return &Service{
		shopRepo:        shopRepo,
		goldRepo:        goldRepo,
		transactionRepo: transactionRepo,
	}
}

// GetGoldPackages returns all gold package items
func (s *Service) GetGoldPackages(ctx context.Context) ([]shop.ShopItem, error) {
	shopType := shop.ItemTypeGoldPackage
	return s.shopRepo.GetShopItems(ctx, nil, &shopType)
}

// GetBundles returns all bundle items
func (s *Service) GetBundles(ctx context.Context) ([]shop.ShopItem, error) {
	shopType := shop.ItemTypeBundle
	return s.shopRepo.GetShopItems(ctx, nil, &shopType)
}

// GetShopItems returns all active shop items
func (s *Service) GetShopItems(ctx context.Context) ([]shop.ShopItem, error) {
	return s.shopRepo.GetShopItems(ctx, nil, nil)
}

// GetShopItemByID returns a shop item by ID
func (s *Service) GetShopItemByID(ctx context.Context, id int) (*shop.ShopItem, error) {
	return s.shopRepo.GetShopItemByID(ctx, id)
}

// GetPlayerGold returns a player's gold balance
func (s *Service) GetPlayerGold(ctx context.Context, userID int) (*shop.PlayerGold, error) {
	return s.goldRepo.GetPlayerGold(ctx, userID)
}

// GetTransactions returns a player's transaction history
func (s *Service) GetTransactions(ctx context.Context, userID int, limit int) ([]shop.Transaction, error) {
	return s.transactionRepo.GetTransactionsByUser(ctx, userID, limit)
}

// PurchaseItem handles item purchase logic
func (s *Service) PurchaseItem(ctx context.Context, userID int, itemID int) (*shop.PurchaseResponse, error) {
	// Get the item
	item, err := s.shopRepo.GetShopItemByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	if item == nil {
		return &shop.PurchaseResponse{
			Success: false,
			Message: "Item not found",
		}, nil
	}

	if !item.IsActive {
		return &shop.PurchaseResponse{
			Success: false,
			Message: "Item is no longer available",
		}, nil
	}

	// Create transaction record
	tx := &shop.Transaction{
		UserID:        userID,
		ShopItemID:    &item.ID,
		ItemType:      item.ItemType,
		ItemName:      item.Name,
		PriceAmount:   item.PriceAmount,
		PriceCurrency: item.PriceCurrency,
		Status:        "pending",
	}

	// Handle different purchase types
	switch item.PriceCurrency {
	case shop.CurrencyGold:
		// Check and deduct gold
		playerGold, err := s.goldRepo.GetPlayerGold(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get player gold: %w", err)
		}

		if playerGold.Balance < item.PriceAmount {
			return &shop.PurchaseResponse{
				Success: false,
				Message: "Insufficient gold",
			}, nil
		}

		// Deduct gold
		if err := s.goldRepo.SpendGold(ctx, userID, item.PriceAmount); err != nil {
			return &shop.PurchaseResponse{
				Success: false,
				Message: "Failed to process payment",
			}, nil
		}

		tx.GoldChange = -item.PriceAmount

	case shop.CurrencyUSD:
		// For USD purchases (mock for now - would integrate with payment provider)
		// In production, this would create a payment intent and wait for webhook
		tx.GoldChange = 0 // USD purchases don't change gold directly

	default:
		return &shop.PurchaseResponse{
			Success: false,
			Message: "Unsupported currency",
		}, nil
	}

	// Create transaction record
	txID, err := s.transactionRepo.CreateTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Deliver items based on type
	var itemsReceived []string

	switch item.ItemType {
	case shop.ItemTypeGoldPackage:
		if item.GoldAmount != nil && *item.GoldAmount > 0 {
			if err := s.goldRepo.AddGold(ctx, userID, *item.GoldAmount); err != nil {
				return nil, fmt.Errorf("failed to add gold: %w", err)
			}
			itemsReceived = append(itemsReceived, fmt.Sprintf("%d Gold", *item.GoldAmount))
		}

	case shop.ItemTypeBundle:
		// Add gold bonus if present
		if item.GoldAmount != nil && *item.GoldAmount > 0 {
			if err := s.goldRepo.AddGold(ctx, userID, *item.GoldAmount); err != nil {
				return nil, fmt.Errorf("failed to add bonus gold: %w", err)
			}
			itemsReceived = append(itemsReceived, fmt.Sprintf("%d Gold (Bonus)", *item.GoldAmount))
		}

		// Equipment would be added to inventory here
		if metadata, ok := item.Metadata["equipment_count"]; ok {
			itemsReceived = append(itemsReceived, fmt.Sprintf("%v Equipment Items", metadata))
		}

		// Drop boosts
		if metadata, ok := item.Metadata["drop_boosts"]; ok {
			itemsReceived = append(itemsReceived, fmt.Sprintf("%v Drop Boosts", metadata))
		}
	}

	// Mark transaction as completed
	if err := s.transactionRepo.UpdateTransactionStatus(ctx, txID, "completed"); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Get updated balance
	playerGold, err := s.goldRepo.GetPlayerGold(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated balance: %w", err)
	}

	return &shop.PurchaseResponse{
		Success:        true,
		TransactionID:  txID,
		NewBalance:     playerGold.Balance,
		ItemsReceived:  itemsReceived,
		Message:        fmt.Sprintf("Successfully purchased %s", item.Name),
	}, nil
}

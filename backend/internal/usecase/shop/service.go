package shop

import (
	"context"
	"fmt"
	"math/rand"

	"empoweredpixels/internal/domain/shop"
	"empoweredpixels/internal/domain/weapons"
	"empoweredpixels/internal/infra/db/repositories"
)

// Service handles shop business logic
type Service struct {
	shopRepo        repositories.ShopRepository
	goldRepo        repositories.PlayerGoldRepository
	transactionRepo repositories.TransactionRepository
	weaponService   WeaponService
	paymentProvider PaymentProvider
}

// WeaponService defines the required interface for weapon delivery
type WeaponService interface {
	AddWeaponToInventory(ctx context.Context, userID int64, weaponDefID string) (*weapons.UserWeapon, error)
}

// NewService creates a new shop service
func NewService(
	shopRepo repositories.ShopRepository,
	goldRepo repositories.PlayerGoldRepository,
	transactionRepo repositories.TransactionRepository,
	weaponService WeaponService,
	paymentProvider PaymentProvider,
) *Service {
	return &Service{
		shopRepo:        shopRepo,
		goldRepo:        goldRepo,
		transactionRepo: transactionRepo,
		weaponService:   weaponService,
		paymentProvider: paymentProvider,
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
		// Process payment through the provider
		providerTxID, err := s.paymentProvider.ProcessPayment(ctx, userID, item.PriceAmount, item.PriceCurrency)
		if err != nil {
			return &shop.PurchaseResponse{
				Success: false,
				Message: fmt.Sprintf("Payment failed: %v", err),
			}, nil
		}

		tx.GoldChange = 0 // USD purchases don't change gold directly
		if tx.Metadata == nil {
			tx.Metadata = make(map[string]interface{})
		}
		tx.Metadata["provider_transaction_id"] = providerTxID

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
		if countVal, ok := item.Metadata["equipment_count"]; ok {
			count, _ := countVal.(float64) // JSON numbers are float64 in Go maps
			if count == 0 {
				count = 1 // Default to 1 if not specified
			}

			guaranteedRarity := weapons.Common
			if rarityVal, ok := item.Metadata["guaranteed_rarity"]; ok {
				guaranteedRarity = weapons.Rarity(int(rarityVal.(float64)))
			}

			for i := 0; i < int(count); i++ {
				// Determine rarity for this drop
				rarity := s.rollRarity(guaranteedRarity)

				// Logging for debugging
				fmt.Printf("Rolling for rarity: guaranteed=%v, rolled=%v\n", guaranteedRarity, rarity)

				// Pick a random weapon of that rarity
				weaponPool := weapons.GetWeaponsByRarity(rarity)
				if len(weaponPool) == 0 {
					// Fallback to common if no weapons found for rarity
					weaponPool = weapons.GetWeaponsByRarity(weapons.Common)
				}

				if len(weaponPool) > 0 {
					weapon := weaponPool[rand.Intn(len(weaponPool))]
					_, err := s.weaponService.AddWeaponToInventory(ctx, int64(userID), weapon.ID)
					if err != nil {
						return nil, fmt.Errorf("failed to grant equipment: %w", err)
					}
					itemsReceived = append(itemsReceived, fmt.Sprintf("%s (%s)", weapon.Name, weapon.Rarity.String()))
				}
			}
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

// rollRarity determines the rarity for a drop, respecting a minimum guaranteed rarity
func (s *Service) rollRarity(guaranteed weapons.Rarity) weapons.Rarity {
	roll := rand.Float64() * 100.0
	cumulative := 0.0

	// Order of rarities to check (highest to lowest)
	rarities := []weapons.Rarity{
		weapons.Unique,
		weapons.Divine,
		weapons.Mythic,
		weapons.Legendary,
		weapons.Epic,
		weapons.Rare,
		weapons.Uncommon,
		weapons.Common,
		weapons.Broken,
	}

	for _, r := range rarities {
		cumulative += r.DropRate()
		if roll <= cumulative {
			if r < guaranteed {
				return guaranteed
			}
			return r
		}
	}

	return guaranteed
}

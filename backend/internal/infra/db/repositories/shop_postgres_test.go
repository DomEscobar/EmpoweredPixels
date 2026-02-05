package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShopRepository_GetShops(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewShopRepository(pool)
	ctx := context.Background()

	shops, err := repo.GetShops(ctx)
	require.NoError(t, err)
	assert.NotNil(t, shops)
	assert.GreaterOrEqual(t, len(shops), 1)

	// Verify seed data exists
	foundGold := false
	for _, shop := range shops {
		if shop.ShopType == "gold" {
			foundGold = true
			assert.True(t, shop.IsActive)
		}
	}
	assert.True(t, foundGold, "Gold shop should exist in seed data")
}

func TestShopRepository_GetShopItems(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewShopRepository(pool)
	ctx := context.Background()

	// Get gold packages
	itemType := "gold_package"
	items, err := repo.GetShopItems(ctx, nil, &itemType)
	require.NoError(t, err)
	assert.NotEmpty(t, items)

	// Verify gold package structure
	for _, item := range items {
		assert.Equal(t, "gold_package", item.ItemType)
		assert.True(t, item.IsActive)
		assert.Greater(t, item.PriceAmount, 0)
		if item.GoldAmount != nil {
			assert.Greater(t, *item.GoldAmount, 0)
		}
	}
}

func TestShopRepository_GetShopItemByID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewShopRepository(pool)
	ctx := context.Background()

	// Get existing item (ID 1 should be seed data)
	item, err := repo.GetShopItemByID(ctx, 1)
	require.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, 1, item.ID)

	// Get non-existent item
	item, err = repo.GetShopItemByID(ctx, 99999)
	require.NoError(t, err)
	assert.Nil(t, item)
}

func TestPlayerGoldRepository_GetPlayerGold(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewPlayerGoldRepository(pool)
	ctx := context.Background()

	// Get gold for new user (should return 0 balance)
	gold, err := repo.GetPlayerGold(ctx, 999999)
	require.NoError(t, err)
	assert.NotNil(t, gold)
	assert.Equal(t, 0, gold.Balance)
}

func TestPlayerGoldRepository_AddGold(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewPlayerGoldRepository(pool)
	ctx := context.Background()

	userID := 12345

	// Add gold to new user
	err := repo.AddGold(ctx, userID, 1000)
	require.NoError(t, err)

	// Verify balance
	gold, err := repo.GetPlayerGold(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 1000, gold.Balance)
	assert.Equal(t, 1000, gold.LifetimeEarned)

	// Add more gold
	err = repo.AddGold(ctx, userID, 500)
	require.NoError(t, err)

	gold, err = repo.GetPlayerGold(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 1500, gold.Balance)
	assert.Equal(t, 1500, gold.LifetimeEarned)
}

func TestPlayerGoldRepository_SpendGold(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewPlayerGoldRepository(pool)
	ctx := context.Background()

	userID := 12346

	// Setup: Add gold first
	err := repo.AddGold(ctx, userID, 1000)
	require.NoError(t, err)

	// Spend gold
	err = repo.SpendGold(ctx, userID, 300)
	require.NoError(t, err)

	// Verify balance
	gold, err := repo.GetPlayerGold(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 700, gold.Balance)
	assert.Equal(t, 300, gold.LifetimeSpent)

	// Try to spend more than available
	err = repo.SpendGold(ctx, userID, 1000)
	assert.Error(t, err)
}

func TestTransactionRepository_CreateAndGet(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTransactionRepository(pool)
	ctx := context.Background()

	userID := 12347

	// Create transaction
	tx := &shop.Transaction{
		UserID:        userID,
		ItemType:      "gold_package",
		ItemName:      "Small Pouch",
		PriceAmount:   99,
		PriceCurrency: "usd",
		GoldChange:    100,
		Status:        "completed",
	}

	id, err := repo.CreateTransaction(ctx, tx)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	// Get transactions
	transactions, err := repo.GetTransactionsByUser(ctx, userID, 10)
	require.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, "Small Pouch", transactions[0].ItemName)
}

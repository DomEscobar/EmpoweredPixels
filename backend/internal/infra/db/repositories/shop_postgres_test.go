package repositories

import (
	"context"
	"testing"

	"empoweredpixels/internal/domain/shop"

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
}

func TestShopRepository_GetShopItems(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewShopRepository(pool)
	ctx := context.Background()

	items, err := repo.GetShopItems(ctx, nil, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, items)
}

func TestPlayerGoldRepository_AddAndSpend(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewPlayerGoldRepository(pool)
	ctx := context.Background()
	userID := 12345

	// Add gold
	err := repo.AddGold(ctx, userID, 1000)
	require.NoError(t, err)

	// Verify
	gold, err := repo.GetPlayerGold(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 1000, gold.Balance)

	// Spend
	err = repo.SpendGold(ctx, userID, 300)
	require.NoError(t, err)

	gold, err = repo.GetPlayerGold(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 700, gold.Balance)
}

func TestTransactionRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTransactionRepository(pool)
	ctx := context.Background()

	tx := &shop.Transaction{
		UserID:        123,
		ItemType:      "gold_package",
		ItemName:      "Test Package",
		PriceAmount:   99,
		PriceCurrency: "usd",
		GoldChange:    100,
		Status:        "completed",
	}

	id, err := repo.CreateTransaction(ctx, tx)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
}

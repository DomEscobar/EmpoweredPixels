package shop

import (
	"context"
	"testing"

	"empoweredpixels/internal/domain/shop"
	"empoweredpixels/internal/domain/weapons"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWeaponService struct {
	mock.Mock
}

func (m *mockWeaponService) AddWeaponToInventory(ctx context.Context, userID int64, weaponDefID string) (*weapons.UserWeapon, error) {
	args := m.Called(ctx, userID, weaponDefID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*weapons.UserWeapon), args.Error(1)
}

// Mock repositories
type mockShopRepo struct {
	mock.Mock
}

func (m *mockShopRepo) GetShops(ctx context.Context) ([]shop.Shop, error) {
	args := m.Called(ctx)
	return args.Get(0).([]shop.Shop), args.Error(1)
}

func (m *mockShopRepo) GetShopByID(ctx context.Context, id int) (*shop.Shop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*shop.Shop), args.Error(1)
}

func (m *mockShopRepo) GetShopItems(ctx context.Context, shopID *int, itemType *string) ([]shop.ShopItem, error) {
	args := m.Called(ctx, shopID, itemType)
	return args.Get(0).([]shop.ShopItem), args.Error(1)
}

func (m *mockShopRepo) GetShopItemByID(ctx context.Context, id int) (*shop.ShopItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*shop.ShopItem), args.Error(1)
}

type mockGoldRepo struct {
	mock.Mock
}

func (m *mockGoldRepo) GetPlayerGold(ctx context.Context, userID int) (*shop.PlayerGold, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*shop.PlayerGold), args.Error(1)
}

func (m *mockGoldRepo) AddGold(ctx context.Context, userID int, amount int) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *mockGoldRepo) SpendGold(ctx context.Context, userID int, amount int) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *mockGoldRepo) CreatePlayerGold(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockGoldRepo) ListAllBalances(ctx context.Context) ([]shop.PlayerGold, error) {
	args := m.Called(ctx)
	return args.Get(0).([]shop.PlayerGold), args.Error(1)
}

type mockTxRepo struct {
	mock.Mock
}

func (m *mockTxRepo) CreateTransaction(ctx context.Context, tx *shop.Transaction) (int, error) {
	args := m.Called(ctx, tx)
	return args.Int(0), args.Error(1)
}

func (m *mockTxRepo) GetTransactionsByUser(ctx context.Context, userID int, limit int) ([]shop.Transaction, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]shop.Transaction), args.Error(1)
}

func (m *mockTxRepo) UpdateTransactionStatus(ctx context.Context, id int, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestService_GetGoldPackages(t *testing.T) {
	shopRepo := new(mockShopRepo)
	goldRepo := new(mockGoldRepo)
	txRepo := new(mockTxRepo)
	weaponService := new(mockWeaponService)

	service := NewService(shopRepo, goldRepo, txRepo, weaponService)

	expected := []shop.ShopItem{
		{ID: 1, Name: "Small Pouch", ItemType: "gold_package", PriceAmount: 99},
		{ID: 2, Name: "Treasure Chest", ItemType: "gold_package", PriceAmount: 999},
	}

	shopRepo.On("GetShopItems", mock.Anything, (*int)(nil), mock.AnythingOfType("*string")).
		Return(expected, nil)

	items, err := service.GetGoldPackages(context.Background())

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "Small Pouch", items[0].Name)
	shopRepo.AssertExpectations(t)
}

func TestService_PurchaseItem_Success(t *testing.T) {
	shopRepo := new(mockShopRepo)
	goldRepo := new(mockGoldRepo)
	txRepo := new(mockTxRepo)
	weaponService := new(mockWeaponService)

	service := NewService(shopRepo, goldRepo, txRepo, weaponService)

	itemID := 1
	userID := 123
	price := 100

	// Setup mocks
	item := &shop.ShopItem{
		ID:            itemID,
		Name:          "Test Item",
		ItemType:      "equipment",
		PriceAmount:   price,
		PriceCurrency: shop.CurrencyGold,
		IsActive:      true,
	}

	shopRepo.On("GetShopItemByID", mock.Anything, itemID).Return(item, nil)
	goldRepo.On("GetPlayerGold", mock.Anything, userID).Return(&shop.PlayerGold{Balance: 500}, nil).Once()
	goldRepo.On("SpendGold", mock.Anything, userID, price).Return(nil)
	txRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*shop.Transaction")).Return(1, nil)
	txRepo.On("UpdateTransactionStatus", mock.Anything, 1, "completed").Return(nil)
	goldRepo.On("GetPlayerGold", mock.Anything, userID).Return(&shop.PlayerGold{Balance: 400}, nil).Once()

	result, err := service.PurchaseItem(context.Background(), userID, itemID)

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 400, result.NewBalance)
	assert.Contains(t, result.Message, "Successfully purchased")
}

func TestService_PurchaseItem_InsufficientGold(t *testing.T) {
	shopRepo := new(mockShopRepo)
	goldRepo := new(mockGoldRepo)
	txRepo := new(mockTxRepo)
	weaponService := new(mockWeaponService)

	service := NewService(shopRepo, goldRepo, txRepo, weaponService)

	itemID := 1
	userID := 123
	price := 1000

	item := &shop.ShopItem{
		ID:            itemID,
		Name:          "Expensive Item",
		ItemType:      "equipment",
		PriceAmount:   price,
		PriceCurrency: shop.CurrencyGold,
		IsActive:      true,
	}

	shopRepo.On("GetShopItemByID", mock.Anything, itemID).Return(item, nil)
	goldRepo.On("GetPlayerGold", mock.Anything, userID).Return(&shop.PlayerGold{Balance: 500}, nil)

	result, err := service.PurchaseItem(context.Background(), userID, itemID)

	assert.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "Insufficient gold", result.Message)
}

func TestService_PurchaseItem_ItemNotFound(t *testing.T) {
	shopRepo := new(mockShopRepo)
	goldRepo := new(mockGoldRepo)
	txRepo := new(mockTxRepo)
	weaponService := new(mockWeaponService)

	service := NewService(shopRepo, goldRepo, txRepo, weaponService)

	shopRepo.On("GetShopItemByID", mock.Anything, 999).Return(nil, nil)

	result, err := service.PurchaseItem(context.Background(), 123, 999)

	assert.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "Item not found", result.Message)
}

func TestService_PurchaseItem_InactiveItem(t *testing.T) {
	shopRepo := new(mockShopRepo)
	goldRepo := new(mockGoldRepo)
	txRepo := new(mockTxRepo)
	weaponService := new(mockWeaponService)

	service := NewService(shopRepo, goldRepo, txRepo, weaponService)

	item := &shop.ShopItem{
		ID:       1,
		Name:     "Inactive Item",
		IsActive: false,
	}

	shopRepo.On("GetShopItemByID", mock.Anything, 1).Return(item, nil)

	result, err := service.PurchaseItem(context.Background(), 123, 1)

	assert.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, "Item is no longer available", result.Message)
}

func TestService_GetPlayerGold(t *testing.T) {
	shopRepo := new(mockShopRepo)
	goldRepo := new(mockGoldRepo)
	txRepo := new(mockTxRepo)
	weaponService := new(mockWeaponService)

	service := NewService(shopRepo, goldRepo, txRepo, weaponService)

	expected := &shop.PlayerGold{
		UserID:  123,
		Balance: 1000,
	}

	goldRepo.On("GetPlayerGold", mock.Anything, 123).Return(expected, nil)

	gold, err := service.GetPlayerGold(context.Background(), 123)

	assert.NoError(t, err)
	assert.Equal(t, 1000, gold.Balance)
	goldRepo.AssertExpectations(t)
}

func TestService_GetTransactions(t *testing.T) {
	shopRepo := new(mockShopRepo)
	goldRepo := new(mockGoldRepo)
	txRepo := new(mockTxRepo)
	weaponService := new(mockWeaponService)

	service := NewService(shopRepo, goldRepo, txRepo, weaponService)

	expected := []shop.Transaction{
		{ID: 1, ItemName: "Item 1"},
		{ID: 2, ItemName: "Item 2"},
	}

	txRepo.On("GetTransactionsByUser", mock.Anything, 123, 50).Return(expected, nil)

	transactions, err := service.GetTransactions(context.Background(), 123, 50)

	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	txRepo.AssertExpectations(t)
}

func TestService_PurchaseItem_BundleEquipment(t *testing.T) {
	shopRepo := new(mockShopRepo)
	goldRepo := new(mockGoldRepo)
	txRepo := new(mockTxRepo)
	weaponService := new(mockWeaponService)

	service := NewService(shopRepo, goldRepo, txRepo, weaponService)

	itemID := 1
	userID := 123
	price := 500

	item := &shop.ShopItem{
		ID:            itemID,
		Name:          "Starter Bundle",
		ItemType:      shop.ItemTypeBundle,
		PriceAmount:   price,
		PriceCurrency: shop.CurrencyGold,
		IsActive:      true,
		Metadata: map[string]interface{}{
			"equipment_count":   2.0,
			"guaranteed_rarity": float64(weapons.Rare),
		},
	}

	shopRepo.On("GetShopItemByID", mock.Anything, itemID).Return(item, nil)
	goldRepo.On("GetPlayerGold", mock.Anything, userID).Return(&shop.PlayerGold{Balance: 1000}, nil).Once()
	goldRepo.On("SpendGold", mock.Anything, userID, price).Return(nil)
	txRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*shop.Transaction")).Return(1, nil)

	// Expect 2 weapons to be added
	weaponService.On("AddWeaponToInventory", mock.Anything, int64(userID), mock.AnythingOfType("string")).
		Return(&weapons.UserWeapon{}, nil).Times(2)

	txRepo.On("UpdateTransactionStatus", mock.Anything, 1, "completed").Return(nil)
	goldRepo.On("GetPlayerGold", mock.Anything, userID).Return(&shop.PlayerGold{Balance: 500}, nil).Once()

	result, err := service.PurchaseItem(context.Background(), userID, itemID)

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Len(t, result.ItemsReceived, 2)
	weaponService.AssertExpectations(t)
}

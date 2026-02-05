package shop

import (
	"context"
	"errors"
	"testing"
	"time"

	"empoweredpixels/internal/domain/shop"
)

// MockRepository implements Repository for testing
type MockRepository struct {
	shops        map[string]*shop.Shop
	items        map[string]*shop.ShopItem
	playerGold   map[int64]*shop.PlayerGold
	transactions map[string]*shop.Transaction
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		shops:        make(map[string]*shop.Shop),
		items:        make(map[string]*shop.ShopItem),
		playerGold:   make(map[int64]*shop.PlayerGold),
		transactions: make(map[string]*shop.Transaction),
	}
}

// Shop methods
func (m *MockRepository) GetShop(ctx context.Context, id string) (*shop.Shop, error) {
	return m.shops[id], nil
}

func (m *MockRepository) ListShops(ctx context.Context) ([]shop.Shop, error) {
	var result []shop.Shop
	for _, s := range m.shops {
		if s.IsActive {
			result = append(result, *s)
		}
	}
	return result, nil
}

// ShopItem methods
func (m *MockRepository) GetShopItem(ctx context.Context, id string) (*shop.ShopItem, error) {
	return m.items[id], nil
}

func (m *MockRepository) ListShopItems(ctx context.Context, shopID string) ([]shop.ShopItem, error) {
	var result []shop.ShopItem
	for _, item := range m.items {
		if item.ShopID == shopID && item.IsActive {
			result = append(result, *item)
		}
	}
	return result, nil
}

func (m *MockRepository) ListAllActiveItems(ctx context.Context) ([]shop.ShopItem, error) {
	var result []shop.ShopItem
	for _, item := range m.items {
		if item.IsActive {
			result = append(result, *item)
		}
	}
	return result, nil
}

func (m *MockRepository) ListItemsByType(ctx context.Context, itemType string) ([]shop.ShopItem, error) {
	var result []shop.ShopItem
	for _, item := range m.items {
		if item.ItemType == itemType && item.IsActive {
			result = append(result, *item)
		}
	}
	return result, nil
}

// PlayerGold methods
func (m *MockRepository) GetPlayerGold(ctx context.Context, userID int64) (*shop.PlayerGold, error) {
	if pg, ok := m.playerGold[userID]; ok {
		return pg, nil
	}
	return &shop.PlayerGold{
		UserID:         userID,
		Balance:        0,
		LifetimeEarned: 0,
		LifetimeSpent:  0,
		Updated:        time.Now(),
	}, nil
}

func (m *MockRepository) AddGold(ctx context.Context, userID int64, amount int) error {
	if pg, ok := m.playerGold[userID]; ok {
		pg.Balance += amount
		pg.LifetimeEarned += amount
	} else {
		m.playerGold[userID] = &shop.PlayerGold{
			UserID:         userID,
			Balance:        amount,
			LifetimeEarned: amount,
			LifetimeSpent:  0,
			Updated:        time.Now(),
		}
	}
	return nil
}

func (m *MockRepository) SpendGold(ctx context.Context, userID int64, amount int) error {
	pg, ok := m.playerGold[userID]
	if !ok || pg.Balance < amount {
		return errors.New("insufficient gold")
	}
	pg.Balance -= amount
	pg.LifetimeSpent += amount
	return nil
}

func (m *MockRepository) SetGold(ctx context.Context, userID int64, balance int) error {
	if pg, ok := m.playerGold[userID]; ok {
		pg.Balance = balance
	} else {
		m.playerGold[userID] = &shop.PlayerGold{
			UserID:  userID,
			Balance: balance,
			Updated: time.Now(),
		}
	}
	return nil
}

// Transaction methods
func (m *MockRepository) CreateTransaction(ctx context.Context, tx *shop.Transaction) error {
	m.transactions[tx.ID] = tx
	return nil
}

func (m *MockRepository) GetTransaction(ctx context.Context, id string) (*shop.Transaction, error) {
	return m.transactions[id], nil
}

func (m *MockRepository) ListTransactions(ctx context.Context, userID int64, limit, offset int) ([]shop.Transaction, error) {
	var result []shop.Transaction
	for _, tx := range m.transactions {
		if tx.UserID == userID {
			result = append(result, *tx)
		}
	}
	return result, nil
}

// Test fixtures
func setupTestService() (*Service, *MockRepository) {
	repo := NewMockRepository()
	now := func() time.Time { return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) }
	service := NewService(repo, now)
	return service, repo
}

func TestGetPlayerGold_NewUser(t *testing.T) {
	service, _ := setupTestService()
	ctx := context.Background()

	gold, err := service.GetPlayerGold(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gold.Balance != 0 {
		t.Errorf("expected balance 0, got %d", gold.Balance)
	}
	if gold.UserID != 1 {
		t.Errorf("expected userID 1, got %d", gold.UserID)
	}
}

func TestGetPlayerGold_ExistingUser(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	repo.playerGold[1] = &shop.PlayerGold{
		UserID:         1,
		Balance:        500,
		LifetimeEarned: 1000,
		LifetimeSpent:  500,
	}

	gold, err := service.GetPlayerGold(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gold.Balance != 500 {
		t.Errorf("expected balance 500, got %d", gold.Balance)
	}
}

func TestListBundles(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	repo.items["1"] = &shop.ShopItem{
		ID:       "1",
		Name:     "Starter Bundle",
		ItemType: shop.ItemTypeBundle,
		IsActive: true,
	}
	repo.items["2"] = &shop.ShopItem{
		ID:       "2",
		Name:     "Gold Pack",
		ItemType: shop.ItemTypeGoldPackage,
		IsActive: true,
	}

	bundles, err := service.ListBundles(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bundles) != 1 {
		t.Errorf("expected 1 bundle, got %d", len(bundles))
	}
	if bundles[0].Name != "Starter Bundle" {
		t.Errorf("expected 'Starter Bundle', got '%s'", bundles[0].Name)
	}
}

func TestListGoldPackages(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	repo.items["1"] = &shop.ShopItem{
		ID:       "1",
		Name:     "Small Gold Pack",
		ItemType: shop.ItemTypeGoldPackage,
		IsActive: true,
	}
	repo.items["2"] = &shop.ShopItem{
		ID:       "2",
		Name:     "Large Gold Pack",
		ItemType: shop.ItemTypeGoldPackage,
		IsActive: true,
	}
	repo.items["3"] = &shop.ShopItem{
		ID:       "3",
		Name:     "Bundle",
		ItemType: shop.ItemTypeBundle,
		IsActive: true,
	}

	packages, err := service.ListGoldPackages(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(packages) != 2 {
		t.Errorf("expected 2 gold packages, got %d", len(packages))
	}
}

func TestGetShopItem_NotFound(t *testing.T) {
	service, _ := setupTestService()
	ctx := context.Background()

	_, err := service.GetShopItem(ctx, "nonexistent")
	if err != ErrItemNotFound {
		t.Errorf("expected ErrItemNotFound, got %v", err)
	}
}

func TestGetShopItem_Found(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	repo.items["test-item"] = &shop.ShopItem{
		ID:       "test-item",
		Name:     "Test Item",
		ItemType: shop.ItemTypeBundle,
		IsActive: true,
	}

	item, err := service.GetShopItem(ctx, "test-item")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Name != "Test Item" {
		t.Errorf("expected 'Test Item', got '%s'", item.Name)
	}
}

func TestPurchase_ItemNotFound(t *testing.T) {
	service, _ := setupTestService()
	ctx := context.Background()

	_, err := service.Purchase(ctx, 1, "nonexistent")
	if err != ErrItemNotFound {
		t.Errorf("expected ErrItemNotFound, got %v", err)
	}
}

func TestPurchase_ItemNotActive(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	repo.items["inactive"] = &shop.ShopItem{
		ID:       "inactive",
		Name:     "Inactive Item",
		IsActive: false,
	}

	_, err := service.Purchase(ctx, 1, "inactive")
	if err != ErrItemNotActive {
		t.Errorf("expected ErrItemNotActive, got %v", err)
	}
}

func TestPurchase_InsufficientGold(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	repo.items["bundle"] = &shop.ShopItem{
		ID:            "bundle",
		Name:          "Expensive Bundle",
		ItemType:      shop.ItemTypeBundle,
		PriceAmount:   1000,
		PriceCurrency: shop.CurrencyGold,
		IsActive:      true,
	}
	repo.playerGold[1] = &shop.PlayerGold{
		UserID:  1,
		Balance: 100, // Not enough
	}

	_, err := service.Purchase(ctx, 1, "bundle")
	if err != ErrInsufficientGold {
		t.Errorf("expected ErrInsufficientGold, got %v", err)
	}
}

func TestPurchase_WithGold_Success(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	repo.items["bundle"] = &shop.ShopItem{
		ID:            "bundle",
		Name:          "Test Bundle",
		ItemType:      shop.ItemTypeBundle,
		PriceAmount:   100,
		PriceCurrency: shop.CurrencyGold,
		IsActive:      true,
	}
	repo.playerGold[1] = &shop.PlayerGold{
		UserID:  1,
		Balance: 500,
	}

	result, err := service.Purchase(ctx, 1, "bundle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NewBalance != 400 {
		t.Errorf("expected balance 400, got %d", result.NewBalance)
	}
	if result.Transaction.GoldChange != -100 {
		t.Errorf("expected gold change -100, got %d", result.Transaction.GoldChange)
	}
	if result.Transaction.Status != shop.StatusCompleted {
		t.Errorf("expected status 'completed', got '%s'", result.Transaction.Status)
	}
}

func TestPurchase_WithUSD_Success(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	goldAmount := 550
	repo.items["gold-pack"] = &shop.ShopItem{
		ID:            "gold-pack",
		Name:          "Gold Pack",
		ItemType:      shop.ItemTypeGoldPackage,
		PriceAmount:   499,
		PriceCurrency: shop.CurrencyUSD,
		GoldAmount:    &goldAmount,
		IsActive:      true,
	}

	result, err := service.Purchase(ctx, 1, "gold-pack")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NewBalance != 550 {
		t.Errorf("expected balance 550, got %d", result.NewBalance)
	}
	if result.Transaction.GoldChange != 550 {
		t.Errorf("expected gold change 550, got %d", result.Transaction.GoldChange)
	}
}

func TestAddGold(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	err := service.AddGold(ctx, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.playerGold[1].Balance != 100 {
		t.Errorf("expected balance 100, got %d", repo.playerGold[1].Balance)
	}

	err = service.AddGold(ctx, 1, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.playerGold[1].Balance != 150 {
		t.Errorf("expected balance 150, got %d", repo.playerGold[1].Balance)
	}
}

func TestGetTransactionHistory(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	repo.transactions["tx1"] = &shop.Transaction{
		ID:       "tx1",
		UserID:   1,
		ItemName: "Bundle 1",
	}
	repo.transactions["tx2"] = &shop.Transaction{
		ID:       "tx2",
		UserID:   1,
		ItemName: "Bundle 2",
	}
	repo.transactions["tx3"] = &shop.Transaction{
		ID:       "tx3",
		UserID:   2, // Different user
		ItemName: "Bundle 3",
	}

	txs, err := service.GetTransactionHistory(ctx, 1, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(txs))
	}
}

func TestListAllItems(t *testing.T) {
	service, repo := setupTestService()
	ctx := context.Background()

	repo.items["1"] = &shop.ShopItem{ID: "1", Name: "Item 1", IsActive: true}
	repo.items["2"] = &shop.ShopItem{ID: "2", Name: "Item 2", IsActive: true}
	repo.items["3"] = &shop.ShopItem{ID: "3", Name: "Item 3", IsActive: false}

	items, err := service.ListAllItems(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 active items, got %d", len(items))
	}
}

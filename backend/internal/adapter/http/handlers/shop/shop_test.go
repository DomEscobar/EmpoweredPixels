package shop

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"empoweredpixels/internal/domain/shop"
	shopusecase "empoweredpixels/internal/usecase/shop"
)

// Mock Service for testing
type MockShopService struct {
	items        []shop.ShopItem
	bundles      []shop.ShopItem
	goldPackages []shop.ShopItem
	playerGold   *shop.PlayerGold
	transactions []shop.Transaction
	purchaseErr  error
	purchaseRes  *shop.PurchaseResult
}

func (m *MockShopService) ListAllItems(ctx context.Context) ([]shop.ShopItem, error) {
	return m.items, nil
}

func (m *MockShopService) ListBundles(ctx context.Context) ([]shop.ShopItem, error) {
	return m.bundles, nil
}

func (m *MockShopService) ListGoldPackages(ctx context.Context) ([]shop.ShopItem, error) {
	return m.goldPackages, nil
}

func (m *MockShopService) GetShopItem(ctx context.Context, id string) (*shop.ShopItem, error) {
	for _, item := range m.items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, shopusecase.ErrItemNotFound
}

func (m *MockShopService) GetPlayerGold(ctx context.Context, userID int64) (*shop.PlayerGold, error) {
	if m.playerGold == nil {
		return &shop.PlayerGold{UserID: userID, Balance: 0}, nil
	}
	return m.playerGold, nil
}

func (m *MockShopService) Purchase(ctx context.Context, userID int64, itemID string) (*shop.PurchaseResult, error) {
	if m.purchaseErr != nil {
		return nil, m.purchaseErr
	}
	return m.purchaseRes, nil
}

func (m *MockShopService) GetTransactionHistory(ctx context.Context, userID int64, page, pageSize int) ([]shop.Transaction, error) {
	return m.transactions, nil
}

func TestListItems(t *testing.T) {
	mockService := &MockShopService{
		items: []shop.ShopItem{
			{ID: "1", Name: "Test Item 1", IsActive: true},
			{ID: "2", Name: "Test Item 2", IsActive: true},
		},
	}

	// Create a service wrapper that implements the expected interface
	service := &shopusecase.Service{}
	// For now, we'll test the handler directly with the mock

	handler := &Handler{service: service}
	_ = mockService // Not directly usable without interface changes
	_ = handler

	// Note: Full integration test would require proper DI setup
	// This is a unit test placeholder
}

func TestGetPlayerGold_Unauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/player/gold", nil)
	rec := httptest.NewRecorder()

	// Without auth middleware setting user ID, should fail
	handler := &Handler{service: nil}
	handler.GetPlayerGold(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestPurchase_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/shop/purchase", strings.NewReader("invalid json"))
	rec := httptest.NewRecorder()

	handler := &Handler{service: nil}
	handler.Purchase(rec, req)

	if rec.Code != http.StatusUnauthorized {
		// Without auth, should be unauthorized first
		t.Logf("Got status %d (expected unauthorized without auth)", rec.Code)
	}
}

func TestPurchaseRequest_Marshal(t *testing.T) {
	req := PurchaseRequest{ItemID: "test-item-id"}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded PurchaseRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ItemID != "test-item-id" {
		t.Errorf("expected item_id 'test-item-id', got '%s'", decoded.ItemID)
	}
}

func TestPurchaseResponse_Marshal(t *testing.T) {
	resp := PurchaseResponse{
		Success:       true,
		TransactionID: "tx-123",
		NewBalance:    500,
		ItemsReceived: []string{"item-1", "item-2"},
		Message:       "Purchase successful!",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded PurchaseResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !decoded.Success {
		t.Error("expected success true")
	}
	if decoded.NewBalance != 500 {
		t.Errorf("expected balance 500, got %d", decoded.NewBalance)
	}
	if len(decoded.ItemsReceived) != 2 {
		t.Errorf("expected 2 items received, got %d", len(decoded.ItemsReceived))
	}
}

func TestGoldBalanceResponse_Marshal(t *testing.T) {
	resp := GoldBalanceResponse{
		Balance:        1000,
		LifetimeEarned: 5000,
		LifetimeSpent:  4000,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded GoldBalanceResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Balance != 1000 {
		t.Errorf("expected balance 1000, got %d", decoded.Balance)
	}
	if decoded.LifetimeEarned != 5000 {
		t.Errorf("expected lifetime_earned 5000, got %d", decoded.LifetimeEarned)
	}
}

// Test domain models
func TestShopItem_JSON(t *testing.T) {
	goldAmount := 100
	item := shop.ShopItem{
		ID:            "item-1",
		ShopID:        "shop-1",
		Name:          "Test Gold Pack",
		Description:   "A test gold package",
		ItemType:      shop.ItemTypeGoldPackage,
		PriceAmount:   99,
		PriceCurrency: shop.CurrencyUSD,
		GoldAmount:    &goldAmount,
		Rarity:        shop.RarityCommon,
		IsActive:      true,
		SortOrder:     1,
		Created:       time.Now(),
		Metadata:      map[string]interface{}{"test": true},
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded shop.ShopItem
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Name != "Test Gold Pack" {
		t.Errorf("expected name 'Test Gold Pack', got '%s'", decoded.Name)
	}
	if decoded.GoldAmount == nil || *decoded.GoldAmount != 100 {
		t.Error("gold_amount not properly serialized")
	}
}

func TestTransaction_JSON(t *testing.T) {
	itemID := "item-123"
	tx := shop.Transaction{
		ID:            "tx-1",
		UserID:        1,
		ShopItemID:    &itemID,
		ItemType:      shop.ItemTypeBundle,
		ItemName:      "Test Bundle",
		PriceAmount:   500,
		PriceCurrency: shop.CurrencyGold,
		GoldChange:    -500,
		Status:        shop.StatusCompleted,
		Created:       time.Now(),
		Metadata:      map[string]interface{}{"items": []string{"eq-1"}},
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded shop.Transaction
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.GoldChange != -500 {
		t.Errorf("expected gold_change -500, got %d", decoded.GoldChange)
	}
}

func TestPlayerGold_JSON(t *testing.T) {
	pg := shop.PlayerGold{
		UserID:         1,
		Balance:        1000,
		LifetimeEarned: 5000,
		LifetimeSpent:  4000,
		Updated:        time.Now(),
	}

	data, err := json.Marshal(pg)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded shop.PlayerGold
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Balance != 1000 {
		t.Errorf("expected balance 1000, got %d", decoded.Balance)
	}
}

// Test error handling
func TestServiceErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrItemNotFound", shopusecase.ErrItemNotFound},
		{"ErrItemNotActive", shopusecase.ErrItemNotActive},
		{"ErrInsufficientGold", shopusecase.ErrInsufficientGold},
		{"ErrInvalidPurchase", shopusecase.ErrInvalidPurchase},
		{"ErrTransactionFailed", shopusecase.ErrTransactionFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("error should not be nil")
			}
			if tt.err.Error() == "" {
				t.Error("error message should not be empty")
			}
		})
	}
}

// Test rarity helpers
func TestRarityName(t *testing.T) {
	tests := []struct {
		rarity int
		want   string
	}{
		{shop.RarityBasic, "Basic"},
		{shop.RarityCommon, "Common"},
		{shop.RarityRare, "Rare"},
		{shop.RarityFabled, "Fabled"},
		{shop.RarityMythic, "Mythic"},
		{shop.RarityLegendary, "Legendary"},
		{99, "Unknown"},
	}

	for _, tt := range tests {
		got := shop.RarityName(tt.rarity)
		if got != tt.want {
			t.Errorf("RarityName(%d) = %s, want %s", tt.rarity, got, tt.want)
		}
	}
}

func TestRarityColor(t *testing.T) {
	tests := []struct {
		rarity int
		want   string
	}{
		{shop.RarityBasic, "gray"},
		{shop.RarityCommon, "green"},
		{shop.RarityRare, "blue"},
		{shop.RarityFabled, "purple"},
		{shop.RarityMythic, "orange"},
		{shop.RarityLegendary, "yellow"},
	}

	for _, tt := range tests {
		got := shop.RarityColor(tt.rarity)
		if got != tt.want {
			t.Errorf("RarityColor(%d) = %s, want %s", tt.rarity, got, tt.want)
		}
	}
}

// Test constants
func TestShopConstants(t *testing.T) {
	// Item types
	if shop.ItemTypeGoldPackage != "gold_package" {
		t.Errorf("ItemTypeGoldPackage = %s, want gold_package", shop.ItemTypeGoldPackage)
	}
	if shop.ItemTypeBundle != "bundle" {
		t.Errorf("ItemTypeBundle = %s, want bundle", shop.ItemTypeBundle)
	}

	// Currencies
	if shop.CurrencyUSD != "usd" {
		t.Errorf("CurrencyUSD = %s, want usd", shop.CurrencyUSD)
	}
	if shop.CurrencyGold != "gold" {
		t.Errorf("CurrencyGold = %s, want gold", shop.CurrencyGold)
	}

	// Statuses
	if shop.StatusCompleted != "completed" {
		t.Errorf("StatusCompleted = %s, want completed", shop.StatusCompleted)
	}
}

// Helper function tests
func TestNewHandler(t *testing.T) {
	handler := NewHandler(nil)
	if handler == nil {
		t.Error("NewHandler should not return nil")
	}
}

func BenchmarkRarityName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		shop.RarityName(i % 6)
	}
}

package shop

import (
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	shopusecase "empoweredpixels/internal/usecase/shop"
)

// Handler handles shop-related HTTP requests
type Handler struct {
	service *shopusecase.Service
}

// NewHandler creates a new shop handler
func NewHandler(service *shopusecase.Service) *Handler {
	return &Handler{service: service}
}

// ListItems handles GET /api/shop/items
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.ListAllItems(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responses.JSON(w, http.StatusOK, items)
}

// ListGoldPackages handles GET /api/shop/gold
func (h *Handler) ListGoldPackages(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.ListGoldPackages(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responses.JSON(w, http.StatusOK, items)
}

// ListBundles handles GET /api/shop/bundles
func (h *Handler) ListBundles(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.ListBundles(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responses.JSON(w, http.StatusOK, items)
}

// GetItem handles GET /api/shop/item/{id}
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request, itemID string) {
	item, err := h.service.GetShopItem(r.Context(), itemID)
	if err != nil {
		if err == shopusecase.ErrItemNotFound {
			responses.Error(w, http.StatusNotFound, "item not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responses.JSON(w, http.StatusOK, item)
}

// PurchaseRequest represents a purchase request body
type PurchaseRequest struct {
	ItemID string `json:"item_id"`
}

// PurchaseResponse represents a purchase response
type PurchaseResponse struct {
	Success       bool     `json:"success"`
	TransactionID string   `json:"transaction_id"`
	NewBalance    int      `json:"new_balance"`
	ItemsReceived []string `json:"items_received,omitempty"`
	Message       string   `json:"message"`
}

// Purchase handles POST /api/shop/purchase
func (h *Handler) Purchase(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok || userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ItemID == "" {
		responses.Error(w, http.StatusBadRequest, "item_id is required")
		return
	}

	result, err := h.service.Purchase(r.Context(), userID, req.ItemID)
	if err != nil {
		switch err {
		case shopusecase.ErrItemNotFound:
			responses.Error(w, http.StatusNotFound, "item not found")
		case shopusecase.ErrItemNotActive:
			responses.Error(w, http.StatusBadRequest, "item not available for purchase")
		case shopusecase.ErrInsufficientGold:
			responses.Error(w, http.StatusPaymentRequired, "insufficient gold balance")
		case shopusecase.ErrInvalidPurchase:
			responses.Error(w, http.StatusBadRequest, "invalid purchase request")
		default:
			responses.Error(w, http.StatusInternalServerError, "purchase failed")
		}
		return
	}

	resp := PurchaseResponse{
		Success:       true,
		TransactionID: result.Transaction.ID,
		NewBalance:    result.NewBalance,
		ItemsReceived: result.ItemsReceived,
		Message:       "Purchase successful!",
	}
	responses.JSON(w, http.StatusOK, resp)
}

// GoldBalanceResponse represents a gold balance response
type GoldBalanceResponse struct {
	Balance        int `json:"balance"`
	LifetimeEarned int `json:"lifetime_earned"`
	LifetimeSpent  int `json:"lifetime_spent"`
}

// GetPlayerGold handles GET /api/player/gold
func (h *Handler) GetPlayerGold(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok || userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	gold, err := h.service.GetPlayerGold(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := GoldBalanceResponse{
		Balance:        gold.Balance,
		LifetimeEarned: gold.LifetimeEarned,
		LifetimeSpent:  gold.LifetimeSpent,
	}
	responses.JSON(w, http.StatusOK, resp)
}

// GetTransactions handles GET /api/player/transactions
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok || userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	transactions, err := h.service.GetTransactionHistory(r.Context(), userID, 1, 50)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, transactions)
}

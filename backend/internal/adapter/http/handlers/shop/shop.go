package shophandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"empoweredpixels/internal/domain/shop"
	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"

	"github.com/gorilla/mux"
)

// ShopService defines the shop service interface
type ShopService interface {
	GetGoldPackages(ctx context.Context) ([]shop.ShopItem, error)
	GetBundles(ctx context.Context) ([]shop.ShopItem, error)
	GetShopItems(ctx context.Context) ([]shop.ShopItem, error)
	GetShopItemByID(ctx context.Context, id int) (*shop.ShopItem, error)
	GetPlayerGold(ctx context.Context, userID int) (*shop.PlayerGold, error)
	GetTransactions(ctx context.Context, userID int, limit int) ([]shop.Transaction, error)
	PurchaseItem(ctx context.Context, userID int, itemID int) (*shop.PurchaseResponse, error)
}

// Handler handles shop HTTP requests
type Handler struct {
	service ShopService
}

// NewHandler creates a new shop handler
func NewHandler(service ShopService) *Handler {
	return &Handler{service: service}
}

// GetGoldPackages handles GET /api/shop/gold
func (h *Handler) GetGoldPackages(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := h.service.GetGoldPackages(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, items)
}

// GetBundles handles GET /api/shop/bundles
func (h *Handler) GetBundles(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := h.service.GetBundles(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, items)
}

// GetShopItems handles GET /api/shop/items
func (h *Handler) GetShopItems(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := h.service.GetShopItems(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, items)
}

// GetShopItem handles GET /api/shop/item/{id}
func (h *Handler) GetShopItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		responses.Error(w, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := h.service.GetShopItemByID(r.Context(), id)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if item == nil {
		responses.Error(w, http.StatusNotFound, "item not found")
		return
	}

	responses.JSON(w, http.StatusOK, item)
}

// Purchase handles POST /api/shop/purchase
func (h *Handler) Purchase(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req shop.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ItemID <= 0 {
		responses.Error(w, http.StatusBadRequest, "invalid item id")
		return
	}

	result, err := h.service.PurchaseItem(r.Context(), userID, req.ItemID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !result.Success {
		responses.JSON(w, http.StatusBadRequest, result)
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

// GetPlayerGold handles GET /api/player/gold
func (h *Handler) GetPlayerGold(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	gold, err := h.service.GetPlayerGold(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, gold)
}

// GetTransactions handles GET /api/player/transactions
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get limit from query param, default to 50
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	transactions, err := h.service.GetTransactions(r.Context(), userID, limit)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, transactions)
}

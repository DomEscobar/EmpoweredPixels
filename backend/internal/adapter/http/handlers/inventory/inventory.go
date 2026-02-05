package inventory

import (
	"encoding/json"
	"log"
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/domain/inventory"
	inventoryusecase "empoweredpixels/internal/usecase/inventory"
)

type Handler struct {
	service inventoryusecase.Service
}

func NewHandler(service inventoryusecase.Service) *Handler {
	return &Handler{service: service}
}

type currencyBalanceDto struct {
	ItemID  string `json:"itemId"`
	Balance int    `json:"balance"`
}

type equipmentDto struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	UserID      int64   `json:"userId"`
	FighterID   *string `json:"fighterId"`
	IsFavorite  bool    `json:"isFavorite"`
	Level       int     `json:"level"`
	Rarity      int     `json:"rarity"`
	Enhancement int     `json:"enhancement"`
}

type enhanceDto struct {
	Equipment          equipmentDto `json:"equipment"`
	DesiredEnhancement int          `json:"desiredEnhancement"`
}

type itemDto struct {
	ID     string `json:"id"`
	ItemID string `json:"itemId"`
	Rarity int    `json:"rarity"`
}

type pagingOptions struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type pageDto[T any] struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	Items      []T `json:"items"`
}

func (h *Handler) BalanceParticles(w http.ResponseWriter, r *http.Request) {
	h.balance(w, r, inventory.EmpoweredParticleID)
}

func (h *Handler) BalanceTokenCommon(w http.ResponseWriter, r *http.Request) {
	h.balance(w, r, inventory.EquipmentTokenCommonID)
}

func (h *Handler) BalanceTokenRare(w http.ResponseWriter, r *http.Request) {
	h.balance(w, r, inventory.EquipmentTokenRareID)
}

func (h *Handler) BalanceTokenFabled(w http.ResponseWriter, r *http.Request) {
	h.balance(w, r, inventory.EquipmentTokenFabledID)
}

func (h *Handler) BalanceTokenMythic(w http.ResponseWriter, r *http.Request) {
	h.balance(w, r, inventory.EquipmentTokenMythicID)
}

func (h *Handler) balance(w http.ResponseWriter, r *http.Request, itemID string) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	balance, err := h.service.Balance(r.Context(), userID, itemID)
	if err != nil {
		log.Printf("inventory balance error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, currencyBalanceDto{
		ItemID:  itemID,
		Balance: balance,
	})
}

func (h *Handler) GetEquipment(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	equip, option, err := h.service.GetEquipment(r.Context(), userID, id)
	if err != nil {
		if err == inventoryusecase.ErrInvalidEquipment {
			responses.Error(w, http.StatusBadRequest, "invalid equipment")
			return
		}
		log.Printf("inventory get equipment error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, mapEquipment(equip, option))
}

func (h *Handler) EnhancementCost(w http.ResponseWriter, r *http.Request) {
	var payload enhanceDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	cost := h.service.EnhancementCost(payload.Equipment.Enhancement, payload.DesiredEnhancement)
	responses.JSON(w, http.StatusOK, cost)
}

func (h *Handler) Enhance(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload enhanceDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	equip, err := h.service.Enhance(r.Context(), userID, payload.Equipment.ID, payload.DesiredEnhancement)
	if err != nil {
		if err == inventoryusecase.ErrInvalidEquipment {
			responses.Error(w, http.StatusBadRequest, "invalid equipment")
			return
		}
		log.Printf("inventory enhance error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	option := &inventory.EquipmentOption{EquipmentID: equip.ID, IsFavorite: false}
	responses.JSON(w, http.StatusOK, mapEquipment(equip, option))
}

func (h *Handler) Salvage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload equipmentDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	items, err := h.service.Salvage(r.Context(), userID, payload.ID)
	if err != nil {
		if err == inventoryusecase.ErrInvalidEquipment {
			responses.Error(w, http.StatusBadRequest, "invalid equipment")
			return
		}
		log.Printf("inventory salvage error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, mapItems(items))
}

func (h *Handler) SalvageInventory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	items, err := h.service.SalvageInventory(r.Context(), userID)
	if err != nil {
		log.Printf("inventory salvage inventory error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, mapItems(items))
}

func (h *Handler) InventoryPage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload pagingOptions
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		payload = pagingOptions{Page: 1, PageSize: 20}
	}

	items, err := h.service.InventoryPage(r.Context(), userID, payload.Page, payload.PageSize)
	if err != nil {
		log.Printf("inventory page error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	result := make([]equipmentDto, 0, len(items))
	for _, equip := range items {
		result = append(result, mapEquipment(&equip, &inventory.EquipmentOption{EquipmentID: equip.ID, IsFavorite: false}))
	}

	responses.JSON(w, http.StatusOK, pageDto[equipmentDto]{
		Page:       payload.Page,
		PageSize:   payload.PageSize,
		TotalCount: len(items),
		Items:      result,
	})
}

func (h *Handler) ListByFighter(w http.ResponseWriter, r *http.Request, fighterID string) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	items, err := h.service.ListByFighter(r.Context(), userID, fighterID)
	if err != nil {
		log.Printf("inventory list by fighter error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	result := make([]equipmentDto, 0, len(items))
	for _, equip := range items {
		result = append(result, mapEquipment(&equip, nil))
	}

	responses.JSON(w, http.StatusOK, result)
}

func (h *Handler) SetFavorite(w http.ResponseWriter, r *http.Request, id string, favorite bool) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	option, err := h.service.SetFavorite(r.Context(), userID, id, favorite)
	if err != nil {
		if err == inventoryusecase.ErrInvalidEquipment {
			responses.Error(w, http.StatusBadRequest, "invalid equipment")
			return
		}
		log.Printf("inventory set favorite error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, mapEquipment(&inventory.Equipment{ID: id}, option))
}

func mapEquipment(equip *inventory.Equipment, option *inventory.EquipmentOption) equipmentDto {
	isFavorite := false
	if option != nil {
		isFavorite = option.IsFavorite
	}
	if equip == nil {
		equip = &inventory.Equipment{}
	}
	return equipmentDto{
		ID:          equip.ID,
		Type:        equip.ItemID,
		UserID:      equip.UserID,
		FighterID:   equip.FighterID,
		IsFavorite:  isFavorite,
		Level:       equip.Level,
		Rarity:      equip.Rarity,
		Enhancement: equip.Enhancement,
	}
}

func mapItems(items []inventory.Item) []itemDto {
	result := make([]itemDto, 0, len(items))
	for _, item := range items {
		result = append(result, itemDto{
			ID:     item.ID,
			ItemID: item.ItemID,
			Rarity: item.Rarity,
		})
	}
	return result
}

package rewards

import (
	"encoding/json"
	"log"
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	rewardsusecase "empoweredpixels/internal/usecase/rewards"
)

type Handler struct {
	service *rewardsusecase.Service
}

func NewHandler(service *rewardsusecase.Service) *Handler {
	return &Handler{service: service}
}

type rewardDto struct {
	ID     string `json:"id"`
	PoolID string `json:"poolId"`
}

type itemDto struct {
	ID     string `json:"id"`
	ItemID string `json:"itemId"`
	Rarity int    `json:"rarity"`
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

type rewardContentDto struct {
	Items     []itemDto      `json:"items"`
	Equipment []equipmentDto `json:"equipment"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	rewards, err := h.service.List(r.Context(), userID)
	if err != nil {
		log.Printf("rewards list error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	result := make([]rewardDto, 0, len(rewards))
	for _, reward := range rewards {
		result = append(result, rewardDto{
			ID:     reward.ID,
			PoolID: reward.RewardPoolID,
		})
	}

	responses.JSON(w, http.StatusOK, result)
}

func (h *Handler) Claim(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload rewardDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	content, err := h.service.Claim(r.Context(), userID, payload.ID, payload.PoolID)
	if err != nil {
		if err == rewardsusecase.ErrInvalidReward {
			responses.Error(w, http.StatusBadRequest, "invalid reward")
			return
		}
		log.Printf("rewards claim error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, mapContent(content))
}

func (h *Handler) ClaimAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	content, err := h.service.ClaimAll(r.Context(), userID)
	if err != nil {
		if err == rewardsusecase.ErrInvalidReward {
			responses.Error(w, http.StatusBadRequest, "invalid reward")
			return
		}
		log.Printf("rewards claim all error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	responses.JSON(w, http.StatusOK, mapContent(content))
}

func mapContent(content *rewardsusecase.RewardContent) rewardContentDto {
	if content == nil {
		return rewardContentDto{}
	}
	items := make([]itemDto, 0, len(content.Items))
	for _, item := range content.Items {
		items = append(items, itemDto{
			ID:     item.ID,
			ItemID: item.ItemID,
			Rarity: item.Rarity,
		})
	}
	equipment := make([]equipmentDto, 0, len(content.Equipment))
	for _, equip := range content.Equipment {
		equipment = append(equipment, equipmentDto{
			ID:          equip.ID,
			Type:        equip.ItemID,
			UserID:      equip.UserID,
			FighterID:   equip.FighterID,
			IsFavorite:  false,
			Level:       equip.Level,
			Rarity:      equip.Rarity,
			Enhancement: equip.Enhancement,
		})
	}
	return rewardContentDto{
		Items:     items,
		Equipment: equipment,
	}
}

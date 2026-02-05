package attunementhandlers

import (
	"context"
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/domain/attunement"
	attunementusecase "empoweredpixels/internal/usecase/attunement"
	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"

	"github.com/gorilla/mux"
)

// AttunementService defines the attunement service interface
type AttunementService interface {
	GetAttunements(ctx context.Context, userID int) (*attunement.PlayerAttunements, error)
	GetAttunementWithBonuses(ctx context.Context, userID int, element attunement.Element) (*attunementusecase.AttunementWithBonuses, error)
	AwardXP(ctx context.Context, userID int, element attunement.Element, source string) (levelUp bool, newLevel int, xpAwarded int, err error)
	GetAllElementsBonus(ctx context.Context, userID int) (*attunementusecase.AggregatedBonuses, error)
}

// Handler handles attunement HTTP requests
type Handler struct {
	service AttunementService
}

// NewHandler creates a new attunement handler
func NewHandler(service AttunementService) *Handler {
	return &Handler{service: service}
}

func getUserID(r *http.Request) int {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		return 0
	}
	return int(userID)
}

// GetAttunements handles GET /api/attunements
func (h *Handler) GetAttunements(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.service.GetAttunements(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

// GetAttunement handles GET /api/attunement/{element}
func (h *Handler) GetAttunement(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)
	elementStr := vars["element"]
	element := attunement.Element(elementStr)

	// Validate element
	valid := false
	for _, e := range attunement.AllElements {
		if e == element {
			valid = true
			break
		}
	}
	if !valid {
		responses.Error(w, http.StatusBadRequest, "invalid element")
		return
	}

	result, err := h.service.GetAttunementWithBonuses(r.Context(), userID, element)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

// AwardXP handles POST /api/attunement/award-xp
func (h *Handler) AwardXP(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Element string `json:"element"`
		Source  string `json:"source"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	element := attunement.Element(req.Element)
	valid := false
	for _, e := range attunement.AllElements {
		if e == element {
			valid = true
			break
		}
	}
	if !valid {
		responses.Error(w, http.StatusBadRequest, "invalid element")
		return
	}

	levelUp, newLevel, xpAwarded, err := h.service.AwardXP(r.Context(), userID, element, req.Source)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"level_up":   levelUp,
		"new_level":  newLevel,
		"xp_awarded": xpAwarded,
	})
}

// GetBonuses handles GET /api/attunements/bonuses
func (h *Handler) GetBonuses(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.service.GetAllElementsBonus(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

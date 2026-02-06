package attunementhandlers

import (
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/domain/attunements"
	attunementusecase "empoweredpixels/internal/usecase/attunements"
)

type Handler struct {
	service *attunementusecase.Service
}

func NewHandler(service *attunementusecase.Service) *Handler {
	return &Handler{service: service}
}

// List returns all attunements for the authenticated user
// GET /api/attunements
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	atts, err := h.service.GetAttunements(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"attunements": atts,
		"count":       len(atts),
	})
}

// Get returns a specific attunement with bonuses
// GET /api/attunement/{element}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, elementStr string) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	element, valid := attunements.ElementFromString(elementStr)
	if !valid {
		responses.Error(w, http.StatusBadRequest, "invalid element: must be Fire, Water, Earth, Wind, Light, or Dark")
		return
	}

	att, err := h.service.GetAttunementWithBonuses(r.Context(), userID, element)
	if err == attunementusecase.ErrAttunementNotFound {
		responses.Error(w, http.StatusNotFound, "attunement not found")
		return
	}
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, att)
}

// AwardXP awards XP to an attunement (for testing/admin)
// POST /api/attunement/award-xp
func (h *Handler) AwardXP(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Element string `json:"element"`
		Amount  int    `json:"amount"`
		Source  string `json:"source"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	element, valid := attunements.ElementFromString(req.Element)
	if !valid {
		responses.Error(w, http.StatusBadRequest, "invalid element: must be Fire, Water, Earth, Wind, Light, or Dark")
		return
	}

	if req.Amount <= 0 {
		responses.Error(w, http.StatusBadRequest, "amount must be positive")
		return
	}

	source := req.Source
	if source == "" {
		source = attunementusecase.XPSourceAdmin
	}

	att, leveledUp, err := h.service.AwardXPAmount(r.Context(), userID, element, req.Amount, source)
	if err == attunementusecase.ErrAttunementNotFound {
		responses.Error(w, http.StatusNotFound, "attunement not found")
		return
	}
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"attunement": att,
		"leveledUp":  leveledUp,
		"xpAwarded":  req.Amount,
		"source":     source,
	})
}

// GetBonuses returns aggregated bonuses from all elements
// GET /api/attunements/bonuses
func (h *Handler) GetBonuses(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bonuses, err := h.service.GetAllElementsBonus(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, bonuses)
}

// GetXPSources returns available XP sources and their amounts
// GET /api/attunements/xp-sources
func (h *Handler) GetXPSources(w http.ResponseWriter, _ *http.Request) {
	sources := []map[string]interface{}{
		{"source": attunementusecase.XPSourceMatchWin, "amount": attunementusecase.XPAmounts[attunementusecase.XPSourceMatchWin], "description": "Winning a match"},
		{"source": attunementusecase.XPSourceMatchLoss, "amount": attunementusecase.XPAmounts[attunementusecase.XPSourceMatchLoss], "description": "Losing a match"},
		{"source": attunementusecase.XPSourceMatchDraw, "amount": attunementusecase.XPAmounts[attunementusecase.XPSourceMatchDraw], "description": "Drawing a match"},
		{"source": attunementusecase.XPSourceDailyLogin, "amount": attunementusecase.XPAmounts[attunementusecase.XPSourceDailyLogin], "description": "Daily login bonus"},
		{"source": attunementusecase.XPSourceQuest, "amount": attunementusecase.XPAmounts[attunementusecase.XPSourceQuest], "description": "Completing a quest"},
		{"source": attunementusecase.XPSourceSkillUse, "amount": attunementusecase.XPAmounts[attunementusecase.XPSourceSkillUse], "description": "Using an attunement skill in combat"},
		{"source": attunementusecase.XPSourceBossDefeat, "amount": attunementusecase.XPAmounts[attunementusecase.XPSourceBossDefeat], "description": "Defeating a boss"},
		{"source": attunementusecase.XPSourceLeagueWin, "amount": attunementusecase.XPAmounts[attunementusecase.XPSourceLeagueWin], "description": "Winning a league match"},
		{"source": attunementusecase.XPSourceAchievement, "amount": attunementusecase.XPAmounts[attunementusecase.XPSourceAchievement], "description": "Unlocking an achievement"},
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"sources": sources,
	})
}

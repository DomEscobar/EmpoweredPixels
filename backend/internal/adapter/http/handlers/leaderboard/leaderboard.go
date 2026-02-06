package leaderboard

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/usecase/leaderboard"
)

// Handler handles leaderboard HTTP requests
type Handler struct {
	service *leaderboard.Service
}

// NewHandler creates a new leaderboard handler
func NewHandler(service *leaderboard.Service) *Handler {
	return &Handler{service: service}
}

func getUserID(r *http.Request) int {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		return 0
	}
	return int(userID)
}

// GetLeaderboard handles GET /api/leaderboard/{category}
func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)
	category := vars["category"]
	if category == "" {
		responses.Error(w, http.StatusBadRequest, "category required")
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	result, err := h.service.GetLeaderboard(r.Context(), category, userID, limit, offset)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

// GetNearbyRanks handles GET /api/leaderboard/{category}/nearby
func (h *Handler) GetNearbyRanks(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)
	category := vars["category"]
	if category == "" {
		responses.Error(w, http.StatusBadRequest, "category required")
		return
	}

	rangeSize := 5
	if rs := r.URL.Query().Get("range"); rs != "" {
		if parsed, err := strconv.Atoi(rs); err == nil && parsed > 0 {
			rangeSize = parsed
		}
	}

	result, err := h.service.GetNearbyRanks(r.Context(), category, userID, rangeSize)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

// GetAchievements handles GET /api/achievements
func (h *Handler) GetAchievements(w http.ResponseWriter, r *http.Request) {
	achievements, err := h.service.GetAchievements(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, achievements)
}

// GetPlayerAchievements handles GET /api/player/achievements
func (h *Handler) GetPlayerAchievements(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	achievements, err := h.service.GetPlayerAchievements(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, achievements)
}

// ClaimAchievement handles POST /api/achievement/{id}/claim
func (h *Handler) ClaimAchievement(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)
	achievementID := vars["id"]
	if achievementID == "" {
		responses.Error(w, http.StatusBadRequest, "achievement id required")
		return
	}

	if err := h.service.ClaimAchievementReward(r.Context(), userID, achievementID); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]bool{"success": true})
}

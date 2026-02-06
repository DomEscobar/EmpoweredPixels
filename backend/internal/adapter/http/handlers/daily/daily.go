package dailyhandlers

import (
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/usecase/daily"
)

// Handler handles daily reward HTTP requests
type Handler struct {
	service *daily.Service
}

// NewHandler creates a new daily reward handler
func NewHandler(service *daily.Service) *Handler {
	return &Handler{service: service}
}

func getUserID(r *http.Request) int {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		return 0
	}
	return int(userID)
}

// GetStatus handles GET /api/daily-reward
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, err := h.service.GetStatus(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, status)
}

// Claim handles POST /api/daily-reward/claim
func (h *Handler) Claim(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.service.Claim(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

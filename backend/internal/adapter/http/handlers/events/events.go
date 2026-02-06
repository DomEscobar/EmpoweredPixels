package events

import (
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/usecase/events"
)

// Handler handles event HTTP requests
type Handler struct {
	service *events.Service
}

// NewHandler creates a new events handler
func NewHandler(service *events.Service) *Handler {
	return &Handler{service: service}
}

func getUserID(r *http.Request) int {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		return 0
	}
	return int(userID)
}

// GetCurrentEvents handles GET /api/events/current
func (h *Handler) GetCurrentEvents(w http.ResponseWriter, r *http.Request) {
	activeEvents, err := h.service.GetCurrentEvents(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, activeEvents)
}

// GetEventStatus handles GET /api/events/status
func (h *Handler) GetEventStatus(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, err := h.service.GetEventStatus(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, status)
}

// GetNextEvent handles GET /api/events/next
func (h *Handler) GetNextEvent(w http.ResponseWriter, r *http.Request) {
	nextEvent, wait, err := h.service.GetNextEventInfo(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]any{
		"event":    nextEvent,
		"wait_ms":  wait.Milliseconds(),
		"wait_str": wait.String(),
	})
}

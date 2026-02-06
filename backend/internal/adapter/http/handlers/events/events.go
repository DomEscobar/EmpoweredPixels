package eventhandlers

import (
	"net/http"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/usecase/events"
)

type Handler struct {
	service *events.Service
}

func NewHandler(service *events.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.service.GetStatus(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responses.JSON(w, http.StatusOK, status)
}

func (h *Handler) GetCurrentEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetCurrentEvents(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responses.JSON(w, http.StatusOK, events)
}

func (h *Handler) GetEventStatus(w http.ResponseWriter, r *http.Request) {
	h.GetStatus(w, r)
}

func (h *Handler) GetNextEvent(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, map[string]string{"message": "Next event coming soon"})
}

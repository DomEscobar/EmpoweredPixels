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

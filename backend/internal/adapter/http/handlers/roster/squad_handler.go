package roster

import (
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	rosterusecase "empoweredpixels/internal/usecase/roster"
)

type SquadHandler struct {
	service *rosterusecase.SquadService
}

func NewSquadHandler(service *rosterusecase.SquadService) *SquadHandler {
	return &SquadHandler{service: service}
}

func (h *SquadHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	squad, err := h.service.GetActiveSquad(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to get active squad")
		return
	}
	if squad == nil {
		responses.Error(w, http.StatusNotFound, "No active squad found")
		return
	}
	responses.JSON(w, http.StatusOK, squad)
}

type SetSquadRequest struct {
	Name       string   `json:"name"`
	FighterIDs []string `json:"fighterIds"`
}

func (h *SquadHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var req SetSquadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	squad, err := h.service.SetActiveSquad(r.Context(), userID, req.Name, req.FighterIDs)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to set active squad")
		return
	}

	responses.JSON(w, http.StatusOK, squad)
}

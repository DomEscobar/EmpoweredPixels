package guilds

import (
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/usecase/guilds"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *guilds.Service
}

func NewHandler(service *guilds.Service) *Handler {
	return &Handler{service: service}
}

type CreateGuildRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Assume user ID is available in context (from auth middleware)
	fighterID := r.Context().Value("user_id").(string)

	var req CreateGuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	guild, err := h.service.CreateGuild(r.Context(), fighterID, req.Name, req.Description)
	if err != nil {
		responses.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responses.JSON(w, http.StatusCreated, guild)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	guilds, err := h.service.ListGuilds(r.Context())
	if err != nil {
		responses.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	responses.JSON(w, http.StatusOK, guilds)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// Note: You would typically want a GetByID method in service too, 
	// for now let's just use what's available or imply service expansion
	responses.JSON(w, http.StatusNotImplemented, map[string]string{"message": "GET single guild TBD"})
}

func (h *Handler) RequestJoin(w http.ResponseWriter, r *http.Request) {
	fighterID := r.Context().Value("user_id").(string)
	guildID := chi.URLParam(r, "id")

	if err := h.service.JoinGuild(r.Context(), fighterID, guildID); err != nil {
		responses.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "request submitted"})
}

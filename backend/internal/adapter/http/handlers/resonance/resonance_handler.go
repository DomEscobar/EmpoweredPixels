package resonance

import (
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/usecase/resonance"
	"github.com/go-chi/chi/v5"
)

// ResonanceResponse represents the API response for squad resonance
type ResonanceResponse struct {
	SquadID            string   `json:"squadID"`
	HarmonyScore       int      `json:"harmonyScore"`
	TierName           string   `json:"tierName"`
	HarmonicElements   []string `json:"harmonicElements"`
	DissonantElements  []string `json:"dissonantElements"`
	Bonuses            BonusInfo `json:"bonuses"`
	AuraColor          string   `json:"auraColor"`
}

// BonusInfo contains damage and defense bonus multipliers
type BonusInfo struct {
	Damage  float64 `json:"damage"`
	Defense float64 `json:"defense"`
}

// ResonanceHandler handles resonance-related HTTP requests
type ResonanceHandler struct {
	service *resonance.ResonanceService
}

// NewResonanceHandler creates a new resonance handler
func NewResonanceHandler(service *resonance.ResonanceService) *ResonanceHandler {
	return &ResonanceHandler{service: service}
}

// GetSquadResonance returns the resonance state for a specific squad
// GET /api/v1/squads/:squadID/resonance
func (h *ResonanceHandler) GetSquadResonance(w http.ResponseWriter, r *http.Request) {
	squadID := chi.URLParam(r, "squadID")
	if squadID == "" {
		responses.Error(w, http.StatusBadRequest, "Missing squad ID")
		return
	}

	state, err := h.service.CalculateSquadResonance(r.Context(), squadID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to calculate resonance")
		return
	}

	if state == nil {
		responses.Error(w, http.StatusNotFound, "Squad not found")
		return
	}

	resp := ResonanceResponse{
		SquadID:           state.SquadID,
		HarmonyScore:      state.HarmonyScore,
		TierName:          state.TierName,
		HarmonicElements:  state.HarmonicElements,
		DissonantElements: state.DissonantElements,
		Bonuses: BonusInfo{
			Damage:  state.BonusDamage,
			Defense: state.BonusDefense,
		},
		AuraColor: state.AuraColor,
	}

	responses.JSON(w, http.StatusOK, resp)
}

// GetUserActiveSquadResonance returns the resonance state for the user's active squad
// GET /api/v1/user/resonance
func (h *ResonanceHandler) GetUserActiveSquadResonance(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		responses.Error(w, http.StatusForbidden, "Unauthorized")
		return
	}

	// For now, return a placeholder. This would require the squad service integration.
	responses.Error(w, http.StatusNotImplemented, "Endpoint requires squad service integration")
}

// PrefetchResonance is a POST endpoint to manually trigger resonance calculation and cache it
// POST /api/v1/squads/:squadID/resonance/prefetch
func (h *ResonanceHandler) PrefetchResonance(w http.ResponseWriter, r *http.Request) {
	squadID := chi.URLParam(r, "squadID")
	if squadID == "" {
		responses.Error(w, http.StatusBadRequest, "Missing squad ID")
		return
	}

	state, err := h.service.CalculateSquadResonance(r.Context(), squadID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to calculate resonance")
		return
	}

	if state == nil {
		responses.Error(w, http.StatusNotFound, "Squad not found")
		return
	}

	resp := ResonanceResponse{
		SquadID:           state.SquadID,
		HarmonyScore:      state.HarmonyScore,
		TierName:          state.TierName,
		HarmonicElements:  state.HarmonicElements,
		DissonantElements: state.DissonantElements,
		Bonuses: BonusInfo{
			Damage:  state.BonusDamage,
			Defense: state.BonusDefense,
		},
		AuraColor: state.AuraColor,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

package skillhandlers

import (
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/adapter/http/responses"
	skillsusecase "empoweredpixels/internal/usecase/skills"
)

type Handler struct {
	service *skillsusecase.Service
}

func NewHandler(service *skillsusecase.Service) *Handler {
	return &Handler{service: service}
}

// GetSkillTree returns the full skill tree
func (h *Handler) GetSkillTree(w http.ResponseWriter, r *http.Request) {
	tree := h.service.GetSkillTree(r.Context())
	responses.JSON(w, http.StatusOK, tree)
}

// GetFighterSkills returns skill state for a specific fighter
func (h *Handler) GetFighterSkills(w http.ResponseWriter, r *http.Request, fighterID string) {
	state, err := h.service.GetFighterSkillState(r.Context(), fighterID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	responses.JSON(w, http.StatusOK, state)
}

// AllocateSkill allocates a skill point
func (h *Handler) AllocateSkill(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	_ = userID // For authorization logging

	var req struct {
		FighterID string `json:"fighterId"`
		SkillID   string `json:"skillId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FighterID == "" || req.SkillID == "" {
		responses.Error(w, http.StatusBadRequest, "fighterId and skillId are required")
		return
	}

	if err := h.service.AllocateSkillPoint(r.Context(), req.FighterID, req.SkillID); err != nil {
		switch err {
		case skillsusecase.ErrSkillNotFound:
			responses.Error(w, http.StatusNotFound, "skill not found")
		case skillsusecase.ErrNoSkillPoints:
			responses.Error(w, http.StatusBadRequest, "no skill points available")
		case skillsusecase.ErrSkillMaxRank:
			responses.Error(w, http.StatusBadRequest, "skill already at maximum rank")
		case skillsusecase.ErrPrerequisitesNotMet:
			responses.Error(w, http.StatusBadRequest, "prerequisites not met for this skill")
		default:
			responses.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SetLoadout sets the active skills loadout
func (h *Handler) SetLoadout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	_ = userID

	var req struct {
		FighterID string   `json:"fighterId"`
		Loadout   []string `json:"loadout"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FighterID == "" {
		responses.Error(w, http.StatusBadRequest, "fighterId is required")
		return
	}

	if err := h.service.SetLoadout(r.Context(), req.FighterID, req.Loadout); err != nil {
		switch err {
		case skillsusecase.ErrLoadoutTooLarge:
			responses.Error(w, http.StatusBadRequest, "loadout exceeds maximum active skills (2)")
		case skillsusecase.ErrSkillNotAllocated:
			responses.Error(w, http.StatusBadRequest, "cannot equip skill that hasn't been allocated")
		case skillsusecase.ErrNotActiveSkill:
			responses.Error(w, http.StatusBadRequest, "only active skills can be in loadout")
		default:
			responses.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetLoadout returns current loadout for a fighter
func (h *Handler) GetLoadout(w http.ResponseWriter, r *http.Request, fighterID string) {
	state, err := h.service.GetFighterSkillState(r.Context(), fighterID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"fighterId": fighterID,
		"loadout":   state.Loadout,
		"maxSlots":  2,
	})
}

// ResetSkills resets all skill allocations
func (h *Handler) ResetSkills(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	_ = userID

	var req struct {
		FighterID string `json:"fighterId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FighterID == "" {
		responses.Error(w, http.StatusBadRequest, "fighterId is required")
		return
	}

	// Get reset cost first
	cost, err := h.service.GetResetCost(r.Context(), req.FighterID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.service.ResetSkills(r.Context(), req.FighterID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"cost":      cost,
		"message":   "Skill tree reset successfully",
	})
}

// GetResetCost returns the cost to reset skills
func (h *Handler) GetResetCost(w http.ResponseWriter, r *http.Request, fighterID string) {
	cost, err := h.service.GetResetCost(r.Context(), fighterID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"fighterId": fighterID,
		"resetCost": cost,
	})
}
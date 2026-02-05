package weaponhandlers

import (
	"encoding/json"
	"net/http"

	"empoweredpixels/internal/adapter/http/responses"
	"empoweredpixels/internal/domain/weapons"
	weaponusecase "empoweredpixels/internal/usecase/weapons"
)

type Handler struct {
	service *weaponusecase.Service
}

func NewHandler(service *weaponusecase.Service) *Handler {
	return &Handler{service: service}
}

// List returns all weapons in user's inventory
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	userWeapons, err := h.service.ListUserWeapons(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var result []WeaponResponse
	for _, uw := range userWeapons {
		weaponDef, found := weapons.GetWeaponByID(uw.WeaponID)
		if !found {
			continue
		}
		stats := weapons.CalculateStats(weaponDef, uw.Enhancement)
		result = append(result, toWeaponResponse(&uw, weaponDef, &stats))
	}

	responses.JSON(w, http.StatusOK, result)
}

// Get returns details for a specific weapon
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, weaponID string) {
	userID := r.Context().Value("userID").(int64)

	uw, weaponDef, stats, err := h.service.GetWeaponDetails(r.Context(), userID, weaponID)
	if err != nil {
		if err == weaponusecase.ErrWeaponNotFound {
			responses.Error(w, http.StatusNotFound, "weapon not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, toWeaponResponse(uw, weaponDef, stats))
}

// Equip equips a weapon to a fighter
func (h *Handler) Equip(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	var req struct {
		WeaponID  string `json:"weaponId"`
		FighterID string `json:"fighterId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.WeaponID == "" || req.FighterID == "" {
		responses.Error(w, http.StatusBadRequest, "weaponId and fighterId are required")
		return
	}

	if err := h.service.EquipWeapon(r.Context(), userID, req.WeaponID, req.FighterID); err != nil {
		switch err {
		case weaponusecase.ErrWeaponNotFound:
			responses.Error(w, http.StatusNotFound, "weapon not found")
		case weaponusecase.ErrWeaponAlreadyEquipped:
			responses.Error(w, http.StatusConflict, "weapon already equipped")
		case weaponusecase.ErrFighterHasWeapon:
			responses.Error(w, http.StatusConflict, "fighter already has a weapon equipped")
		default:
			responses.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Unequip removes a weapon from a fighter
func (h *Handler) Unequip(w http.ResponseWriter, r *http.Request, weaponID string) {
	userID := r.Context().Value("userID").(int64)

	if err := h.service.UnequipWeapon(r.Context(), userID, weaponID); err != nil {
		if err == weaponusecase.ErrWeaponNotFound {
			responses.Error(w, http.StatusNotFound, "weapon not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Enhance attempts to enhance a weapon
func (h *Handler) Enhance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	var req struct {
		WeaponID string `json:"weaponId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.WeaponID == "" {
		responses.Error(w, http.StatusBadRequest, "weaponId is required")
		return
	}

	result, err := h.service.EnhanceWeapon(r.Context(), userID, req.WeaponID)
	if err != nil {
		switch err {
		case weaponusecase.ErrWeaponNotFound:
			responses.Error(w, http.StatusNotFound, "weapon not found")
		case weaponusecase.ErrMaxEnhancement:
			responses.Error(w, http.StatusBadRequest, "weapon is at maximum enhancement")
		default:
			responses.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	responses.JSON(w, http.StatusOK, EnhancementResponse{
		Success:       result.Success,
		NewLevel:      result.NewLevel,
		PreviousLevel: result.PreviousLevel,
		Destroyed:     result.Destroyed,
	})
}

// ForgePreview returns enhancement odds without modifying anything
func (h *Handler) ForgePreview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	var req struct {
		WeaponID string `json:"weaponId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.WeaponID == "" {
		responses.Error(w, http.StatusBadRequest, "weaponId is required")
		return
	}

	weaponDef, nextLevel, successChance, cost, err := h.service.PreviewEnhancement(r.Context(), userID, req.WeaponID)
	if err != nil {
		if err == weaponusecase.ErrWeaponNotFound {
			responses.Error(w, http.StatusNotFound, "weapon not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, ForgePreviewResponse{
		WeaponID:      req.WeaponID,
		WeaponName:    weaponDef.Name,
		CurrentLevel:  nextLevel - 1,
		NextLevel:     nextLevel,
		SuccessChance: successChance,
		Cost:          cost,
	})
}

// Database returns all weapon definitions (for client reference)
func (h *Handler) Database(w http.ResponseWriter, r *http.Request) {
	var result []WeaponDefinitionResponse
	for _, w := range weapons.WeaponDatabase {
		result = append(result, WeaponDefinitionResponse{
			ID:          w.ID,
			Name:        w.Name,
			Type:        w.Type.String(),
			Rarity:      w.Rarity.String(),
			BaseDamage:  w.BaseDamage,
			AttackSpeed: w.AttackSpeed,
			CritChance:  w.CritChance,
			Durability:  w.Durability,
			IconURL:     w.IconURL,
			Description: w.Description,
		})
	}
	responses.JSON(w, http.StatusOK, result)
}

// Response types
type WeaponResponse struct {
	ID            string  `json:"id"`
	WeaponID      string  `json:"weaponId"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Rarity        string  `json:"rarity"`
	Enhancement   int     `json:"enhancement"`
	Durability    int     `json:"durability"`
	IsEquipped    bool    `json:"isEquipped"`
	FighterID     *string `json:"fighterId,omitempty"`
	Damage        int     `json:"damage"`
	AttackSpeed   float64 `json:"attackSpeed"`
	CritChance    int     `json:"critChance"`
	IconURL       string  `json:"iconUrl"`
	Description   string  `json:"description"`
}

type EnhancementResponse struct {
	Success       bool `json:"success"`
	NewLevel      int  `json:"newLevel"`
	PreviousLevel int  `json:"previousLevel"`
	Destroyed     bool `json:"destroyed"`
}

type ForgePreviewResponse struct {
	WeaponID      string  `json:"weaponId"`
	WeaponName    string  `json:"weaponName"`
	CurrentLevel  int     `json:"currentLevel"`
	NextLevel     int     `json:"nextLevel"`
	SuccessChance float64 `json:"successChance"`
	Cost          int     `json:"cost"`
}

type WeaponDefinitionResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Rarity      string  `json:"rarity"`
	BaseDamage  int     `json:"baseDamage"`
	AttackSpeed float64 `json:"attackSpeed"`
	CritChance  int     `json:"critChance"`
	Durability  int     `json:"durability"`
	IconURL     string  `json:"iconUrl"`
	Description string  `json:"description"`
}

func toWeaponResponse(uw *weapons.UserWeapon, def *weapons.Weapon, stats *weapons.WeaponStats) WeaponResponse {
	return WeaponResponse{
		ID:          uw.ID,
		WeaponID:    uw.WeaponID,
		Name:        def.Name,
		Type:        def.Type.String(),
		Rarity:      def.Rarity.String(),
		Enhancement: uw.Enhancement,
		Durability:  uw.Durability,
		IsEquipped:  uw.IsEquipped,
		FighterID:   uw.FighterID,
		Damage:      stats.Damage,
		AttackSpeed: stats.AttackSpeed,
		CritChance:  stats.CritChance,
		IconURL:     def.IconURL,
		Description: def.Description,
	}
}
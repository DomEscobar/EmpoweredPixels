package roster

import (
	"encoding/json"
	"log"
	"net/http"

	"empoweredpixels/internal/adapter/http/middleware"
	"empoweredpixels/internal/adapter/http/responses"
	rosterusecase "empoweredpixels/internal/usecase/roster"
)

type FighterHandler struct {
	service *rosterusecase.Service
}

func NewFighterHandler(service *rosterusecase.Service) *FighterHandler {
	return &FighterHandler{service: service}
}

type fighterDto struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Level          int     `json:"level"`
	CurrentExp     int     `json:"currentExp"`
	LevelExp       int     `json:"levelExp"`
	Power          int     `json:"power"`
	ConditionPower int     `json:"conditionPower"`
	Precision      int     `json:"precision"`
	Ferocity       int     `json:"ferocity"`
	Accuracy       int     `json:"accuracy"`
	Agility        int     `json:"agility"`
	Armor          int     `json:"armor"`
	Vitality       int     `json:"vitality"`
	ParryChance    int     `json:"parryChance"`
	HealingPower   int     `json:"healingPower"`
	Speed          int     `json:"speed"`
	Vision         int     `json:"vision"`
	WeaponID       *string `json:"weaponId"`
	Created        string  `json:"created"`
}

type fighterNameDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type fighterExperienceDto struct {
	Level      int `json:"level"`
	CurrentExp int `json:"currentExp"`
	LevelExp   int `json:"levelExp"`
}

func (h *FighterHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fighters, err := h.service.List(r.Context(), userID)
	if err != nil {
		log.Printf("roster list error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	result := make([]fighterDto, 0, len(fighters))
	for _, fighter := range fighters {
		exp, _ := h.service.GetExperience(r.Context(), fighter.ID)
		level, current, next := calculateLevel(exp.Experience)
		result = append(result, fighterDto{
			ID:             fighter.ID,
			Name:           fighter.Name,
			Level:          level,
			CurrentExp:     current,
			LevelExp:       next,
			Power:          fighter.Power,
			ConditionPower: fighter.ConditionPower,
			Precision:      fighter.Precision,
			Ferocity:       fighter.Ferocity,
			Accuracy:       fighter.Accuracy,
			Agility:        fighter.Agility,
			Armor:          fighter.Armor,
			Vitality:       fighter.Vitality,
			ParryChance:    fighter.ParryChance,
			HealingPower:   fighter.HealingPower,
			Speed:          fighter.Speed,
			Vision:         fighter.Vision,
			WeaponID:       fighter.WeaponID,
			Created:        fighter.Created.Format(timeLayout),
		})
	}

	responses.JSON(w, http.StatusOK, result)
}

func (h *FighterHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fighter, err := h.service.Get(r.Context(), userID, id)
	if err != nil {
		log.Printf("roster get error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}
	if fighter == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	exp, _ := h.service.GetExperience(r.Context(), fighter.ID)
	level, current, next := calculateLevel(exp.Experience)
	responses.JSON(w, http.StatusOK, fighterDto{
		ID:             fighter.ID,
		Name:           fighter.Name,
		Level:          level,
		CurrentExp:     current,
		LevelExp:       next,
		Power:          fighter.Power,
		ConditionPower: fighter.ConditionPower,
		Precision:      fighter.Precision,
		Ferocity:       fighter.Ferocity,
		Accuracy:       fighter.Accuracy,
		Agility:        fighter.Agility,
		Armor:          fighter.Armor,
		Vitality:       fighter.Vitality,
		ParryChance:    fighter.ParryChance,
		HealingPower:   fighter.HealingPower,
		Speed:          fighter.Speed,
		Vision:         fighter.Vision,
		WeaponID:       fighter.WeaponID,
		Created:        fighter.Created.Format(timeLayout),
	})
}

func (h *FighterHandler) GetName(w http.ResponseWriter, r *http.Request, id string) {
	fighter, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("roster get name error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}
	if fighter == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responses.JSON(w, http.StatusOK, fighterNameDto{
		ID:   fighter.ID,
		Name: fighter.Name,
	})
}

func (h *FighterHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	fighter, err := h.service.Create(r.Context(), userID, payload.Name)
	if err != nil {
		switch err {
		case rosterusecase.ErrFighterNameExists, rosterusecase.ErrFighterExists, rosterusecase.ErrInvalidFighter:
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		default:
			log.Printf("roster create error: %v", err)
			responses.Error(w, http.StatusInternalServerError, "server error")
			return
		}
	}

	exp, _ := h.service.GetExperience(r.Context(), fighter.ID)
	level, current, next := calculateLevel(exp.Experience)
	responses.JSON(w, http.StatusOK, fighterDto{
		ID:             fighter.ID,
		Name:           fighter.Name,
		Level:          level,
		CurrentExp:     current,
		LevelExp:       next,
		Power:          fighter.Power,
		ConditionPower: fighter.ConditionPower,
		Precision:      fighter.Precision,
		Ferocity:       fighter.Ferocity,
		Accuracy:       fighter.Accuracy,
		Agility:        fighter.Agility,
		Armor:          fighter.Armor,
		Vitality:       fighter.Vitality,
		ParryChance:    fighter.ParryChance,
		HealingPower:   fighter.HealingPower,
		Speed:          fighter.Speed,
		Vision:         fighter.Vision,
		WeaponID:       fighter.WeaponID,
		Created:        fighter.Created.Format(timeLayout),
	})
}

func (h *FighterHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := h.service.Delete(r.Context(), userID, id); err != nil {
		log.Printf("roster delete error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FighterHandler) GetExperience(w http.ResponseWriter, r *http.Request, id string) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fighter, err := h.service.Get(r.Context(), userID, id)
	if err != nil || fighter == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	exp, err := h.service.GetExperience(r.Context(), id)
	if err != nil {
		log.Printf("roster get experience error: %v", err)
		responses.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	level, current, next := calculateLevel(exp.Experience)
	responses.JSON(w, http.StatusOK, fighterExperienceDto{
		Level:      level,
		CurrentExp: current,
		LevelExp:   next,
	})
}

const timeLayout = "2006-01-02T15:04:05Z07:00"

func calculateLevel(exp int) (level int, current int, next int) {
	level = 1
	next = 100
	current = exp
	for current >= next {
		current -= next
		level++
		next = 100 * level
	}
	return level, current, next
}

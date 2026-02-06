package attunements

import (
	"context"
	"errors"

	"empoweredpixels/internal/domain/attunements"
	"empoweredpixels/internal/infra/db/repositories"
)

var (
	ErrAttunementNotFound = errors.New("attunement not found")
	ErrInvalidElement     = errors.New("invalid element")
	ErrInvalidXPSource    = errors.New("invalid XP source")
)

// XP Source constants define how much XP each activity grants
const (
	XPSourceMatchWin      = "match_win"
	XPSourceMatchLoss     = "match_loss"
	XPSourceMatchDraw     = "match_draw"
	XPSourceDailyLogin    = "daily_login"
	XPSourceQuest         = "quest"
	XPSourceSkillUse      = "skill_use"
	XPSourceBossDefeat    = "boss_defeat"
	XPSourceLeagueWin     = "league_win"
	XPSourceAchievement   = "achievement"
	XPSourceAdmin         = "admin"
)

// XP amounts for each source
var XPAmounts = map[string]int{
	XPSourceMatchWin:    50,
	XPSourceMatchLoss:   15,
	XPSourceMatchDraw:   25,
	XPSourceDailyLogin:  10,
	XPSourceQuest:       100,
	XPSourceSkillUse:    5,
	XPSourceBossDefeat:  200,
	XPSourceLeagueWin:   150,
	XPSourceAchievement: 75,
	XPSourceAdmin:       0, // Admin can specify any amount
}

type AttunementRepository interface {
	GetPlayerAttunements(ctx context.Context, userID int64) ([]repositories.PlayerAttunement, error)
	GetAttunement(ctx context.Context, userID int64, element attunements.Element) (*repositories.PlayerAttunement, error)
	AddXP(ctx context.Context, userID int64, element attunements.Element, xpAmount int, source string) (*repositories.PlayerAttunement, bool, error)
	CreateInitialAttunements(ctx context.Context, userID int64) error
	HasAttunements(ctx context.Context, userID int64) (bool, error)
	GetXPHistory(ctx context.Context, userID int64, limit int) ([]repositories.XPHistoryEntry, error)
}

type Service struct {
	repo AttunementRepository
}

func NewService(repo AttunementRepository) *Service {
	return &Service{repo: repo}
}

// AttunementResponse represents an attunement with calculated bonuses
type AttunementResponse struct {
	Element       string  `json:"element"`
	ElementID     int     `json:"elementId"`
	Level         int     `json:"level"`
	XP            int     `json:"xp"`
	XPToNextLevel int     `json:"xpToNextLevel"`
	Progress      float64 `json:"progress"`
	TotalXPEarned int     `json:"totalXpEarned"`
	IconURL       string  `json:"iconUrl"`
	Color         string  `json:"color"`
	Description   string  `json:"description"`
}

// AttunementWithBonuses includes passive bonuses
type AttunementWithBonuses struct {
	AttunementResponse
	Bonuses       BonusInfo       `json:"bonuses"`
	ActiveAbility AbilityInfo     `json:"activeAbility"`
	Strengths     []string        `json:"strengths"`
	Weaknesses    []string        `json:"weaknesses"`
}

// BonusInfo represents the stat bonuses from an attunement
type BonusInfo struct {
	DamageBonus    float64 `json:"damageBonus"`
	HealingBonus   float64 `json:"healingBonus"`
	ArmorBonus     float64 `json:"armorBonus"`
	SpeedBonus     float64 `json:"speedBonus"`
	CritBonus      float64 `json:"critBonus"`
	LifestealBonus float64 `json:"lifestealBonus"`
}

// AbilityInfo represents an active ability
type AbilityInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cooldown    int    `json:"cooldown"`
	ManaCost    int    `json:"manaCost"`
	EffectType  string `json:"effectType"`
	EffectValue int    `json:"effectValue"`
}

// AggregatedBonuses represents combined bonuses from all elements
type AggregatedBonuses struct {
	TotalDamageBonus    float64 `json:"totalDamageBonus"`
	TotalHealingBonus   float64 `json:"totalHealingBonus"`
	TotalArmorBonus     float64 `json:"totalArmorBonus"`
	TotalSpeedBonus     float64 `json:"totalSpeedBonus"`
	TotalCritBonus      float64 `json:"totalCritBonus"`
	TotalLifestealBonus float64 `json:"totalLifestealBonus"`
}

// GetAttunements returns all attunements for a user
func (s *Service) GetAttunements(ctx context.Context, userID int64) ([]AttunementResponse, error) {
	// Ensure user has attunements initialized
	has, err := s.repo.HasAttunements(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !has {
		if err := s.repo.CreateInitialAttunements(ctx, userID); err != nil {
			return nil, err
		}
	}

	playerAttunements, err := s.repo.GetPlayerAttunements(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []AttunementResponse
	for _, a := range playerAttunements {
		result = append(result, s.toAttunementResponse(a))
	}

	return result, nil
}

// GetAttunementWithBonuses returns a specific attunement with all bonuses calculated
func (s *Service) GetAttunementWithBonuses(ctx context.Context, userID int64, element attunements.Element) (*AttunementWithBonuses, error) {
	// Ensure user has attunements initialized
	has, err := s.repo.HasAttunements(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !has {
		if err := s.repo.CreateInitialAttunements(ctx, userID); err != nil {
			return nil, err
		}
	}

	a, err := s.repo.GetAttunement(ctx, userID, element)
	if errors.Is(err, repositories.ErrAttunementNotFound) {
		return nil, ErrAttunementNotFound
	}
	if err != nil {
		return nil, err
	}

	return s.toAttunementWithBonuses(*a), nil
}

// AwardXP awards XP based on a source constant
func (s *Service) AwardXP(ctx context.Context, userID int64, element attunements.Element, source string) (*AttunementResponse, bool, error) {
	xpAmount, exists := XPAmounts[source]
	if !exists {
		return nil, false, ErrInvalidXPSource
	}

	// Ensure user has attunements initialized
	has, err := s.repo.HasAttunements(ctx, userID)
	if err != nil {
		return nil, false, err
	}
	if !has {
		if err := s.repo.CreateInitialAttunements(ctx, userID); err != nil {
			return nil, false, err
		}
	}

	a, leveledUp, err := s.repo.AddXP(ctx, userID, element, xpAmount, source)
	if err != nil {
		return nil, false, err
	}

	resp := s.toAttunementResponse(*a)
	return &resp, leveledUp, nil
}

// AwardXPAmount awards a specific amount of XP (for admin/custom use)
func (s *Service) AwardXPAmount(ctx context.Context, userID int64, element attunements.Element, amount int, source string) (*AttunementResponse, bool, error) {
	if source == "" {
		source = XPSourceAdmin
	}

	// Ensure user has attunements initialized
	has, err := s.repo.HasAttunements(ctx, userID)
	if err != nil {
		return nil, false, err
	}
	if !has {
		if err := s.repo.CreateInitialAttunements(ctx, userID); err != nil {
			return nil, false, err
		}
	}

	a, leveledUp, err := s.repo.AddXP(ctx, userID, element, amount, source)
	if err != nil {
		return nil, false, err
	}

	resp := s.toAttunementResponse(*a)
	return &resp, leveledUp, nil
}

// GetAllElementsBonus returns aggregated bonuses from all elements
func (s *Service) GetAllElementsBonus(ctx context.Context, userID int64) (*AggregatedBonuses, error) {
	// Ensure user has attunements initialized
	has, err := s.repo.HasAttunements(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !has {
		if err := s.repo.CreateInitialAttunements(ctx, userID); err != nil {
			return nil, err
		}
	}

	playerAttunements, err := s.repo.GetPlayerAttunements(ctx, userID)
	if err != nil {
		return nil, err
	}

	agg := &AggregatedBonuses{}
	for _, a := range playerAttunements {
		bonus := attunements.GetPassiveBonus(a.Element, a.Level)
		agg.TotalDamageBonus += bonus.DamageBonus
		agg.TotalHealingBonus += bonus.HealingBonus
		agg.TotalArmorBonus += bonus.ArmorBonus
		agg.TotalSpeedBonus += bonus.SpeedBonus
		agg.TotalCritBonus += bonus.CritBonus
		agg.TotalLifestealBonus += bonus.LifestealBonus
	}

	return agg, nil
}

// EnsureAttunements initializes attunements for a user if not already done
func (s *Service) EnsureAttunements(ctx context.Context, userID int64) error {
	has, err := s.repo.HasAttunements(ctx, userID)
	if err != nil {
		return err
	}
	if !has {
		return s.repo.CreateInitialAttunements(ctx, userID)
	}
	return nil
}

// Helper to convert repository model to response
func (s *Service) toAttunementResponse(a repositories.PlayerAttunement) AttunementResponse {
	xpToNext := attunements.XPNeededForLevel(a.Level)
	progress := 0.0
	if a.Level < 25 && xpToNext > 0 {
		progress = float64(a.XP) / float64(xpToNext) * 100.0
	} else if a.Level >= 25 {
		progress = 100.0
		xpToNext = 0
	}

	return AttunementResponse{
		Element:       a.Element.String(),
		ElementID:     int(a.Element),
		Level:         a.Level,
		XP:            a.XP,
		XPToNextLevel: xpToNext,
		Progress:      progress,
		TotalXPEarned: a.TotalXPEarned,
		IconURL:       a.Element.GetIconURL(),
		Color:         a.Element.GetColor(),
		Description:   attunements.GetDescription(a.Element),
	}
}

// Helper to convert repository model to full response with bonuses
func (s *Service) toAttunementWithBonuses(a repositories.PlayerAttunement) *AttunementWithBonuses {
	resp := s.toAttunementResponse(a)
	bonus := attunements.GetPassiveBonus(a.Element, a.Level)
	ability := attunements.GetActiveAbility(a.Element)

	// Get strengths and weaknesses as strings
	var strengths, weaknesses []string
	for _, el := range a.Element.GetStrengths() {
		strengths = append(strengths, el.String())
	}
	for _, el := range a.Element.GetWeaknesses() {
		weaknesses = append(weaknesses, el.String())
	}

	return &AttunementWithBonuses{
		AttunementResponse: resp,
		Bonuses: BonusInfo{
			DamageBonus:    bonus.DamageBonus,
			HealingBonus:   bonus.HealingBonus,
			ArmorBonus:     bonus.ArmorBonus,
			SpeedBonus:     bonus.SpeedBonus,
			CritBonus:      bonus.CritBonus,
			LifestealBonus: bonus.LifestealBonus,
		},
		ActiveAbility: AbilityInfo{
			Name:        ability.Name,
			Description: ability.Description,
			Cooldown:    ability.Cooldown,
			ManaCost:    ability.ManaCost,
			EffectType:  ability.EffectType,
			EffectValue: ability.EffectValue,
		},
		Strengths:  strengths,
		Weaknesses: weaknesses,
	}
}

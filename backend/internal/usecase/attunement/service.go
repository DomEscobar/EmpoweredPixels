package attunement

import (
	"context"
	"fmt"

	"empoweredpixels/internal/domain/attunement"
)

// Service handles attunement business logic
type Service struct {
	repo AttunementRepository
}

// NewService creates a new attunement service
func NewService(repo AttunementRepository) *Service {
	return &Service{repo: repo}
}

// GetAttunements returns all 6 attunements for a player
func (s *Service) GetAttunements(ctx context.Context, userID int) (*attunement.PlayerAttunements, error) {
	attunements, err := s.repo.GetPlayerAttunements(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attunements: %w", err)
	}

	// If no attunements exist, create initial ones
	if len(attunements) == 0 {
		if err := s.repo.CreateInitialAttunements(ctx, userID); err != nil {
			return nil, fmt.Errorf("failed to create initial attunements: %w", err)
		}
		attunements, err = s.repo.GetPlayerAttunements(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get attunements after creation: %w", err)
		}
	}

	totalLevel := 0
	for _, a := range attunements {
		totalLevel += a.Level
	}

	return &attunement.PlayerAttunements{
		UserID:      userID,
		Attunements: attunements,
		TotalLevel:  totalLevel,
	}, nil
}

// GetAttunementWithBonuses returns a specific attunement with calculated bonuses
func (s *Service) GetAttunementWithBonuses(ctx context.Context, userID int, element attunement.Element) (*AttunementWithBonuses, error) {
	a, err := s.repo.GetAttunement(ctx, userID, element)
	if err != nil {
		return nil, fmt.Errorf("failed to get attunement: %w", err)
	}

	if a == nil {
		// Create initial attunement
		if err := s.repo.CreateInitialAttunements(ctx, userID); err != nil {
			return nil, fmt.Errorf("failed to create initial attunements: %w", err)
		}
		a, err = s.repo.GetAttunement(ctx, userID, element)
		if err != nil {
			return nil, fmt.Errorf("failed to get attunement after creation: %w", err)
		}
	}

	bonus := attunement.GetBonus(a.Element, a.Level)
	xpRequired := attunement.GetXPRequired(a.Level)

	return &AttunementWithBonuses{
		Attunement: *a,
		Bonuses:    bonus,
		XPRequired: xpRequired,
		Progress:   calculateProgress(a.CurrentXP, xpRequired),
	}, nil
}

// AttunementWithBonuses combines attunement data with bonuses
type AttunementWithBonuses struct {
	attunement.Attunement
	Bonuses    attunement.AttunementBonus `json:"bonuses"`
	XPRequired int                        `json:"xp_required"`
	Progress   float64                    `json:"progress"` // 0.0 to 1.0
}

// AwardXP grants XP to an element based on source
func (s *Service) AwardXP(ctx context.Context, userID int, element attunement.Element, source string) (levelUp bool, newLevel int, xpAwarded int, err error) {
	// XP amounts based on source
	xpAmount := getXPAward(source)
	if xpAmount == 0 {
		return false, 0, 0, fmt.Errorf("unknown xp source: %s", source)
	}

	levelUp, newLevel, err = s.repo.AddXP(ctx, userID, element, xpAmount, source)
	if err != nil {
		return false, 0, 0, fmt.Errorf("failed to add xp: %w", err)
	}

	return levelUp, newLevel, xpAmount, nil
}

// GetAllElementsBonus returns aggregated bonuses across all elements
func (s *Service) GetAllElementsBonus(ctx context.Context, userID int) (*AggregatedBonuses, error) {
	attunements, err := s.GetAttunements(ctx, userID)
	if err != nil {
		return nil, err
	}

	aggregated := &AggregatedBonuses{
		ByElement: make(map[string]attunement.AttunementBonus),
	}

	for _, a := range attunements.Attunements {
		bonus := attunement.GetBonus(a.Element, a.Level)
		aggregated.ByElement[string(a.Element)] = bonus
		aggregated.TotalPower += bonus.Power
		aggregated.TotalDefense += bonus.Defense
		aggregated.TotalSpeed += bonus.Speed
		aggregated.TotalPrecision += bonus.Precision
	}

	return aggregated, nil
}

// AggregatedBonuses combines all element bonuses
type AggregatedBonuses struct {
	ByElement      map[string]attunement.AttunementBonus `json:"by_element"`
	TotalPower     float64                               `json:"total_power"`
	TotalDefense   float64                               `json:"total_defense"`
	TotalSpeed     float64                               `json:"total_speed"`
	TotalPrecision float64                               `json:"total_precision"`
}

// getXPAward returns XP amount for a given source
func getXPAward(source string) int {
	switch source {
	case "match_win":
		return attunement.XPWinMatch
	case "daily_task":
		return attunement.XPCompleteTask
	case "element_use":
		return attunement.XPUseElement
	case "equip_item":
		return attunement.XPEquipItem
	case "test": // For testing/admin
		return 100
	default:
		return 0
	}
}

// calculateProgress returns progress percentage (0.0 to 1.0)
func calculateProgress(currentXP, xpRequired int) float64 {
	if xpRequired <= 0 {
		return 1.0 // Max level
	}
	progress := float64(currentXP) / float64(xpRequired)
	if progress > 1.0 {
		return 1.0
	}
	return progress
}

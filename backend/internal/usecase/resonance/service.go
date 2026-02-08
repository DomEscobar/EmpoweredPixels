package resonance

import (
	"context"
	"time"

	"empoweredpixels/internal/domain/attunement"
	"empoweredpixels/internal/domain/roster"
)

// ResonanceState contains calculated resonance information for a squad
type ResonanceState struct {
	SquadID          string
	HarmonyScore     int
	TierName         string
	HarmonicElements []string
	DissonantElements []string
	BonusDamage      float64
	BonusDefense     float64
	AuraColor        string
	Pattern          *attunement.ResonancePattern
}

// SquadRepository defines repository interface for squads
type SquadRepository interface {
	GetByID(ctx context.Context, squadID string) (*roster.Squad, error)
	UpdateResonanceScore(ctx context.Context, squadID string, score int, pattern string, timestamp time.Time) error
}

// FighterRepository defines repository interface for fighters
type FighterRepository interface {
	GetByID(ctx context.Context, fighterID string) (*roster.Fighter, error)
	GetMultiple(ctx context.Context, fighterIDs []string) ([]roster.Fighter, error)
}

// AttunementRepository defines repository interface for attunements
type AttunementRepository interface {
	GetPlayerAttunements(ctx context.Context, userID int) ([]attunement.Attunement, error)
}

// ResonanceService encapsulates resonance calculations
type ResonanceService struct {
	squadRepo       SquadRepository
	fighterRepo     FighterRepository
	attunementRepo  AttunementRepository
}

// NewResonanceService creates a new resonance service
func NewResonanceService(
	squadRepo SquadRepository,
	fighterRepo FighterRepository,
	attunementRepo AttunementRepository,
) *ResonanceService {
	return &ResonanceService{
		squadRepo:      squadRepo,
		fighterRepo:    fighterRepo,
		attunementRepo: attunementRepo,
	}
}

// CalculateSquadResonance calculates the resonance state for a squad
func (s *ResonanceService) CalculateSquadResonance(ctx context.Context, squadID string) (*ResonanceState, error) {
	// Get squad
	squad, err := s.squadRepo.GetByID(ctx, squadID)
	if err != nil {
		return nil, err
	}
	if squad == nil {
		return &ResonanceState{SquadID: squadID, HarmonyScore: 0, TierName: "Discordant"}, nil
	}

	// Get player attunements by UserID
	userAttunements, err := s.attunementRepo.GetPlayerAttunements(ctx, int(squad.UserID))
	if err != nil {
		return nil, err
	}

	// Extract top 3 attunements by level (for squad of 3 fighters)
	// This creates a mapping of attunement elements
	var elements []attunement.Element
	for _, att := range userAttunements {
		// Only include attunements with level >= 1
		if att.Level >= 1 {
			elements = append(elements, att.Element)
		}
	}

	// Calculate harmony
	score, pattern := attunement.CalculateHarmony(elements)

	// Get multipliers
	dmgMult, defMult := attunement.GetBonusMultiplier(score)

	// Build harmonicElements and dissonantElements lists
	harmonicElements := []string{}
	for _, pair := range pattern.HarmonicPairs {
		if !contains(harmonicElements, string(pair.Element1)) {
			harmonicElements = append(harmonicElements, string(pair.Element1))
		}
		if !contains(harmonicElements, string(pair.Element2)) {
			harmonicElements = append(harmonicElements, string(pair.Element2))
		}
	}

	dissonantElements := []string{}
	for _, pair := range pattern.DissonantPairs {
		if !contains(dissonantElements, string(pair.Element1)) {
			dissonantElements = append(dissonantElements, string(pair.Element1))
		}
		if !contains(dissonantElements, string(pair.Element2)) {
			dissonantElements = append(dissonantElements, string(pair.Element2))
		}
	}

	state := &ResonanceState{
		SquadID:           squadID,
		HarmonyScore:      score,
		TierName:          pattern.TierName,
		HarmonicElements:  harmonicElements,
		DissonantElements: dissonantElements,
		BonusDamage:       dmgMult,
		BonusDefense:      defMult,
		AuraColor:         pattern.AuraColor,
		Pattern:           pattern,
	}

	// Update database with calculated resonance
	err = s.squadRepo.UpdateResonanceScore(ctx, squadID, score, pattern.TierName, time.Now())
	if err != nil {
		// Log error but don't fail the operation
		_ = err
	}

	return state, nil
}

// ApplySquadBonuses applies resonance bonuses to a squad's fighters
func (s *ResonanceService) ApplySquadBonuses(fighters []roster.Fighter, state *ResonanceState) []roster.Fighter {
	if state == nil || state.HarmonyScore == 0 {
		return fighters
	}

	dmgMult, defMult := attunement.GetBonusMultiplier(state.HarmonyScore)

	// Apply bonuses to all fighters in the squad
	for i := range fighters {
		fighters[i].Power = int(float64(fighters[i].Power) * dmgMult)
		fighters[i].Armor = int(float64(fighters[i].Armor) * defMult)
		fighters[i].ConditionPower = int(float64(fighters[i].ConditionPower) * dmgMult)
	}

	return fighters
}

// GetBonusMultiplier returns damage and defense multipliers for a harmony score
func (s *ResonanceService) GetBonusMultiplier(score int) (dmg, def float64) {
	return attunement.GetBonusMultiplier(score)
}

// GetTierName returns the tier name for a harmony score
func (s *ResonanceService) GetTierName(score int) string {
	switch {
	case score >= 76:
		return "Resonant"
	case score >= 51:
		return "Harmonized"
	case score >= 26:
		return "Aligned"
	default:
		return "Discordant"
	}
}

// GetDissonanceWarning returns a warning message for dissonant squads
func (s *ResonanceService) GetDissonanceWarning(pattern *attunement.ResonancePattern) string {
	return attunement.GetDissonanceWarning(pattern)
}

// contains is a helper to check if a string is in a slice
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

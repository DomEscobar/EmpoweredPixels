package momentum

import (
	"context"
	"time"

	"empoweredpixels/internal/domain/momentum"
)

type MomentumRepository interface {
	GetStakedMomentum(ctx context.Context, fighterID string) (*momentum.StakedMomentum, error)
	UpdateStakedMomentum(ctx context.Context, sm *momentum.StakedMomentum) error
}

type Service struct {
	repo MomentumRepository
}

func NewService(repo MomentumRepository) *Service {
	return &Service{repo: repo}
}

// StakeMomentum converts in-combat momentum from a finished match to staked momentum
func (s *Service) StakeMomentum(ctx context.Context, fighterID string, earnedMomentum float64) error {
	sm, err := s.repo.GetStakedMomentum(ctx, fighterID)
	if err != nil {
		return err
	}

	now := time.Now()
	// Apply decay first if any
	decay := sm.CalculateDecay(now)
	sm.CurrentValue -= decay
	if sm.CurrentValue < 0 {
		sm.CurrentValue = 0
	}

	// Add new momentum (diminishing returns or capped)
	sm.CurrentValue += (earnedMomentum * 0.5) // Conversion rate: match momentum is worth 50% as permanent stake
	if sm.CurrentValue > momentum.MaxStakedMomentum {
		sm.CurrentValue = momentum.MaxStakedMomentum
	}

	sm.LastUpdatedAt = now
	return s.repo.UpdateStakedMomentum(ctx, sm)
}

// GetActiveBonuses returns currently active bonuses for a fighter
func (s *Service) GetActiveBonuses(ctx context.Context, fighterID string) ([]momentum.MomentumBonus, error) {
	sm, err := s.repo.GetStakedMomentum(ctx, fighterID)
	if err != nil {
		return nil, err
	}

	var active []momentum.MomentumBonus
	for _, m := range momentum.DefaultMilestones {
		if sm.CurrentValue >= m.Threshold {
			m.IsUnlocked = true
			active = append(active, m)
		}
	}
	return active, nil
}

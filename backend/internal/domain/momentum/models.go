package momentum

import (
	"time"
)

// Constants for Momentum Logic
const (
	MaxStakedMomentum = 10.0
	BaseDecayRate     = 0.05 // 5% decay per hour of inactivity
	StakingInterval   = 1 * time.Hour
)

// MomentumBonusType defines what kind of passive benefit staking provides
type MomentumBonusType string

const (
	BonusGoldGain      MomentumBonusType = "gold_gain"
	BonusXPGain        MomentumBonusType = "xp_gain"
	BonusMagicFind     MomentumBonusType = "magic_find"
	BonusStatBoost     MomentumBonusType = "stat_boost"
)

// StakedMomentum represents the persistent momentum a fighter has "banked"
// for passive bonuses, as opposed to the volatile in-combat momentum.
type StakedMomentum struct {
	FighterID     string
	CurrentValue  float64
	LastUpdatedAt time.Time
	ActiveBonuses []MomentumBonus
}

// MomentumBonus defines an individual bonus granted by staked momentum
type MomentumBonus struct {
	Type       MomentumBonusType
	Value      float64 // Multiplier or flat value
	IsUnlocked bool
	Threshold  float64 // Staked amount required to activate
}

// CalculateDecay calculates how much momentum is lost based on time passed
func (s *StakedMomentum) CalculateDecay(now time.Time) float64 {
	if s.LastUpdatedAt.IsZero() {
		return 0
	}
	hoursPassed := now.Sub(s.LastUpdatedAt).Hours()
	if hoursPassed < 1 {
		return 0
	}
	
	// Example Decay: Exponential or linear decay based on hours
	decay := s.CurrentValue * (BaseDecayRate * hoursPassed)
	return decay
}

// GetEffectiveMultiplier returns the reward multiplier based on staked momentum
func (s *StakedMomentum) GetEffectiveMultiplier() float64 {
	// Formula: 1.0 + (Momentum / 10.0) * MaxBonus
	// e.g., 5.0 Stake -> 1.0 + 0.5 * 0.2 = 1.1x (10% bonus)
	const MaxBonus = 0.25 
	return 1.0 + (s.CurrentValue/MaxStakedMomentum)*MaxBonus
}

// Define available staking tiers/milestones
var DefaultMilestones = []MomentumBonus{
	{Type: BonusGoldGain, Value: 1.05, Threshold: 2.0},  // +5% Gold at 2.0 Momentum
	{Type: BonusXPGain, Value: 1.10, Threshold: 5.0},    // +10% XP at 5.0 Momentum
	{Type: BonusMagicFind, Value: 1.15, Threshold: 8.0}, // +15% Rare Drop at 8.0 Momentum
	{Type: BonusStatBoost, Value: 5.0, Threshold: 10.0}, // +5 All Stats at Max Momentum
}

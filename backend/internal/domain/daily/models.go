package daily

import (
	"time"
)

// Reward represents a daily reward configuration
type Reward struct {
	Day         int    `json:"day"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Type        string `json:"type"` // "gold", "item", "boost", "mystery"
	Value       int    `json:"value,omitempty"` // Gold amount or boost multiplier
	Rarity      int    `json:"rarity,omitempty"` // For item rewards
}

// UserDailyReward represents a user's daily reward status
type UserDailyReward struct {
	UserID        int        `json:"user_id" db:"user_id"`
	Streak        int        `json:"streak" db:"streak"`
	LastClaimed   *time.Time `json:"last_claimed" db:"last_claimed"`
	TotalClaimed  int        `json:"total_claimed" db:"total_claimed"`
	Updated       time.Time  `json:"updated" db:"updated"`
	CanClaim      bool       `json:"can_claim"`
	NextReward    Reward     `json:"next_reward"`
	TimeUntilReset string    `json:"time_until_reset,omitempty"`
}

// ClaimResult represents the result of claiming a reward
type ClaimResult struct {
	Success        bool      `json:"success"`
	Reward         Reward    `json:"reward"`
	RewardValue    int       `json:"reward_value,omitempty"` // Actual gold amount or item ID
	NewStreak      int       `json:"new_streak"`
	Day            int       `json:"day"`
	NextReward     Reward    `json:"next_reward"`
}

// RewardSchedule defines the 7-day reward cycle
var RewardSchedule = []Reward{
	{Day: 1, Name: "Small Pouch", Description: "100 Gold", Icon: "🪙", Type: "gold", Value: 100},
	{Day: 2, Name: "Common Chest", Description: "250 Gold + Common Item", Icon: "📦", Type: "gold", Value: 250, Rarity: 1},
	{Day: 3, Name: "Rare Cache", Description: "500 Gold + Rare Item", Icon: "💎", Type: "gold", Value: 500, Rarity: 2},
	{Day: 4, Name: "Energy Boost", Description: "2x XP for 1 hour", Icon: "⚡", Type: "boost", Value: 2},
	{Day: 5, Name: "Mystery Box", Description: "Random item (Common to Mythic)", Icon: "🎁", Type: "mystery"},
	{Day: 6, Name: "Fabled Vault", Description: "1000 Gold + Fabled Item", Icon: "🏆", Type: "gold", Value: 1000, Rarity: 3},
	{Day: 7, Name: "Legendary Crate", Description: "2000 Gold + Guaranteed Legendary", Icon: "👑", Type: "gold", Value: 2000, Rarity: 5},
}

// GetRewardForDay returns the reward for a specific day
func GetRewardForDay(day int) Reward {
	if day < 1 {
		day = 1
	}
	// Cycle through rewards (day 8 = day 1, etc.)
	index := (day - 1) % len(RewardSchedule)
	reward := RewardSchedule[index]
	reward.Day = day
	return reward
}

// HasStreakBroken checks if the streak should be reset
func HasStreakBroken(lastClaimed *time.Time) bool {
	if lastClaimed == nil {
		return false
	}
	now := time.Now()
	// Streak breaks if more than 48 hours since last claim
	// (Gives players a grace period of almost 2 days)
	deadline := lastClaimed.Add(48 * time.Hour)
	return now.After(deadline)
}

// CanClaimToday checks if user can claim reward today
func CanClaimToday(lastClaimed *time.Time) bool {
	if lastClaimed == nil {
		return true
	}
	now := time.Now()
	// Can claim if last claim was before today (midnight)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastClaim := time.Date(lastClaimed.Year(), lastClaimed.Month(), lastClaimed.Day(), 0, 0, 0, 0, lastClaimed.Location())
	return lastClaim.Before(today)
}

// TimeUntilMidnight returns time until next midnight
func TimeUntilMidnight() time.Duration {
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return tomorrow.Sub(now)
}

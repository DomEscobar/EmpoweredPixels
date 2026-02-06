package leaderboard

import "time"

// Entry represents a single leaderboard entry
type Entry struct {
	ID           string    `json:"id" db:"id"`
	Category     string    `json:"category" db:"category"`
	UserID       int       `json:"user_id" db:"user_id"`
	Username     string    `json:"username" db:"username"`
	Avatar       string    `json:"avatar,omitempty" db:"avatar"`
	Rank         int       `json:"rank" db:"rank"`
	Score        int64     `json:"score" db:"score"`
	PreviousRank int       `json:"previous_rank" db:"previous_rank"`
	Trend        string    `json:"trend" db:"trend"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ListResponse represents a leaderboard list response
type ListResponse struct {
	Category    string  `json:"category"`
	TotalCount  int     `json:"total_count"`
	UserRank    int     `json:"user_rank"`
	UserEntry   *Entry  `json:"user_entry,omitempty"`
	Entries     []Entry `json:"entries"`
}

// Category definitions
const (
	CategoryPower       = "power"
	CategoryWealth      = "wealth"
	CategoryCombat      = "combat"
	CategoryAchievements = "achievements"
	CategoryStreak      = "streak"
)

// Trend types
const (
	TrendUp   = "up"
	TrendDown = "down"
	TrendSame = "same"
)

// Achievement represents an achievement definition
type Achievement struct {
	ID               string    `json:"id" db:"id"`
	Key              string    `json:"key" db:"key"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	Icon             string    `json:"icon" db:"icon"`
	Category         string    `json:"category" db:"category"`
	RequirementType  string    `json:"requirement_type" db:"requirement_type"`
	RequirementValue int       `json:"requirement_value" db:"requirement_value"`
	RewardGold       int       `json:"reward_gold" db:"reward_gold"`
	RewardTitle      string    `json:"reward_title,omitempty" db:"reward_title"`
	Hidden           bool      `json:"hidden" db:"hidden"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// PlayerAchievement represents a player's progress on an achievement
type PlayerAchievement struct {
	ID            string     `json:"id" db:"id"`
	UserID        int        `json:"user_id" db:"user_id"`
	AchievementID string     `json:"achievement_id" db:"achievement_id"`
	Achievement   *Achievement `json:"achievement,omitempty"`
	Progress      int        `json:"progress" db:"progress"`
	Completed     bool       `json:"completed" db:"completed"`
	CompletedAt   *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	Claimed       bool       `json:"claimed" db:"claimed"`
	ClaimedAt     *time.Time `json:"claimed_at,omitempty" db:"claimed_at"`
}

// AchievementCategories
const (
	AchievementCategoryCombat      = "combat"
	AchievementCategoryCollection = "collection"
	AchievementCategoryProgression = "progression"
	AchievementCategorySocial      = "social"
)

package rewards

import "time"

type Reward struct {
	ID           string
	UserID       int64
	RewardPoolID string
	SourceID     *string // ID of the match or event that triggered the reward
	Claimed      *time.Time
	Created      time.Time
}

package rewards

import "time"

type Reward struct {
	ID           string
	UserID       int64
	RewardPoolID string
	Claimed      *time.Time
	Created      time.Time
}

package jobs

import "context"

type LoginRewardJob struct {
}

func NewLoginRewardJob() *LoginRewardJob {
	return &LoginRewardJob{}
}

func (j *LoginRewardJob) CreateLoginRewards(ctx context.Context) error {
	return nil
}

package jobs

import "context"

type SeasonInitiatorJob struct {
}

func NewSeasonInitiatorJob() *SeasonInitiatorJob {
	return &SeasonInitiatorJob{}
}

func (j *SeasonInitiatorJob) Init(ctx context.Context) error {
	return nil
}

package repositories

import (
	"context"
	"errors"
	"time"

	"empoweredpixels/internal/domain/rewards"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RewardRepository struct {
	pool *pgxpool.Pool
}

func NewRewardRepository(pool *pgxpool.Pool) *RewardRepository {
	return &RewardRepository{pool: pool}
}

func (r *RewardRepository) Create(ctx context.Context, reward *rewards.Reward) error {
	const query = `
		insert into rewards (id, user_id, reward_pool_id, claimed, created)
		values ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, query, reward.ID, reward.UserID, reward.RewardPoolID, reward.Claimed, reward.Created)
	return err
}

func (r *RewardRepository) ListUnclaimed(ctx context.Context, userID int64) ([]rewards.Reward, error) {
	const query = `
		select id, user_id, reward_pool_id, claimed, created
		from rewards
		where user_id = $1 and claimed is null`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []rewards.Reward
	for rows.Next() {
		var reward rewards.Reward
		if err := rows.Scan(&reward.ID, &reward.UserID, &reward.RewardPoolID, &reward.Claimed, &reward.Created); err != nil {
			return nil, err
		}
		result = append(result, reward)
	}
	return result, rows.Err()
}

func (r *RewardRepository) GetUnclaimed(ctx context.Context, userID int64, rewardID string, poolID string) (*rewards.Reward, error) {
	const query = `
		select id, user_id, reward_pool_id, claimed, created
		from rewards
		where id = $1 and reward_pool_id = $2 and user_id = $3 and claimed is null`

	var reward rewards.Reward
	err := r.pool.QueryRow(ctx, query, rewardID, poolID, userID).Scan(
		&reward.ID, &reward.UserID, &reward.RewardPoolID, &reward.Claimed, &reward.Created,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reward, nil
}

func (r *RewardRepository) ListAllUnclaimed(ctx context.Context, userID int64) ([]rewards.Reward, error) {
	return r.ListUnclaimed(ctx, userID)
}

func (r *RewardRepository) MarkClaimed(ctx context.Context, rewardID string, claimedAt time.Time) error {
	const query = `
		update rewards
		set claimed = $1
		where id = $2`

	_, err := r.pool.Exec(ctx, query, claimedAt, rewardID)
	return err
}

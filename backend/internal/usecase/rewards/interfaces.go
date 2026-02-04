package rewards

import (
	"context"
	"time"

	"empoweredpixels/internal/domain/inventory"
	"empoweredpixels/internal/domain/rewards"
)

type RewardRepository interface {
	Create(ctx context.Context, reward *rewards.Reward) error
	ListUnclaimed(ctx context.Context, userID int64) ([]rewards.Reward, error)
	GetUnclaimed(ctx context.Context, userID int64, rewardID string, poolID string) (*rewards.Reward, error)
	ListAllUnclaimed(ctx context.Context, userID int64) ([]rewards.Reward, error)
	MarkClaimed(ctx context.Context, rewardID string, claimedAt time.Time) error
}

type ItemRepository interface {
	CreateMany(ctx context.Context, items []inventory.Item) error
}

type EquipmentRepository interface {
	Create(ctx context.Context, equipment *inventory.Equipment) error
}

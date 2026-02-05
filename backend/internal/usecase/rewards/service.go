package rewards

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"empoweredpixels/internal/domain/inventory"
	"empoweredpixels/internal/domain/rewards"

	"github.com/google/uuid"
)

var (
	ErrInvalidReward = errors.New("invalid reward")
)

type Service struct {
	rewards   RewardRepository
	items     ItemRepository
	equipment EquipmentRepository
	now       func() time.Time
}

func NewService(rewards RewardRepository, items ItemRepository, equipment EquipmentRepository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{
		rewards:   rewards,
		items:     items,
		equipment: equipment,
		now:       now,
	}
}

type RewardContent struct {
	Items     []inventory.Item
	Equipment []inventory.Equipment
}

func (s *Service) List(ctx context.Context, userID int64) ([]rewards.Reward, error) {
	return s.rewards.ListUnclaimed(ctx, userID)
}

func (s *Service) IssueReward(ctx context.Context, userID int64, poolID string) (*rewards.Reward, error) {
	reward := &rewards.Reward{
		ID:           uuid.NewString(),
		UserID:       userID,
		RewardPoolID: poolID,
		Created:      s.now(),
	}

	if err := s.rewards.Create(ctx, reward); err != nil {
		return nil, err
	}

	return reward, nil
}

func (s *Service) Claim(ctx context.Context, userID int64, rewardID string, poolID string) (*RewardContent, error) {
	reward, err := s.rewards.GetUnclaimed(ctx, userID, rewardID, poolID)
	if err != nil {
		return nil, err
	}
	if reward == nil {
		return nil, ErrInvalidReward
	}

	if err := s.rewards.MarkClaimed(ctx, rewardID, s.now()); err != nil {
		return nil, err
	}

	content := s.generateRewards(userID, poolID)
	if err := s.items.CreateMany(ctx, content.Items); err != nil {
		return nil, err
	}
	for i := range content.Equipment {
		if err := s.equipment.Create(ctx, &content.Equipment[i]); err != nil {
			return nil, err
		}
	}
	return &content, nil
}

func (s *Service) ClaimAll(ctx context.Context, userID int64) (*RewardContent, error) {
	rewardsList, err := s.rewards.ListAllUnclaimed(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(rewardsList) == 0 {
		return nil, ErrInvalidReward
	}

	for _, reward := range rewardsList {
		if err := s.rewards.MarkClaimed(ctx, reward.ID, s.now()); err != nil {
			return nil, err
		}
	}

	var all RewardContent
	for _, reward := range rewardsList {
		content := s.generateRewards(userID, reward.RewardPoolID)
		all.Items = append(all.Items, content.Items...)
		all.Equipment = append(all.Equipment, content.Equipment...)
	}
	if err := s.items.CreateMany(ctx, all.Items); err != nil {
		return nil, err
	}
	for i := range all.Equipment {
		if err := s.equipment.Create(ctx, &all.Equipment[i]); err != nil {
			return nil, err
		}
	}
	return &all, nil
}

func (s *Service) generateRewards(userID int64, poolID string) RewardContent {
	items := make([]inventory.Item, 0, 10)
	equipment := make([]inventory.Equipment, 0)

	// Base reward for everyone: 20 particles
	for i := 0; i < 20; i++ {
		items = append(items, inventory.Item{
			ID:      uuid.NewString(),
			UserID:  userID,
			ItemID:  inventory.EmpoweredParticleID,
			Rarity:  inventory.ItemRarityBasic,
			Created: s.now(),
		})
	}

	// Winner bonus
	if poolID == "match_win" {
		// 100 more particles for the winner
		for i := 0; i < 100; i++ {
			items = append(items, inventory.Item{
				ID:      uuid.NewString(),
				UserID:  userID,
				ItemID:  inventory.EmpoweredParticleID,
				Rarity:  inventory.ItemRarityBasic,
				Created: s.now(),
			})
		}
		// Guaranteed Common Token for a win
		items = append(items, inventory.Item{
			ID:      uuid.NewString(),
			UserID:  userID,
			ItemID:  inventory.EquipmentTokenCommonID,
			Rarity:  inventory.ItemRarityCommon,
			Created: s.now(),
		})
	} else if poolID == "match_participation" {
		// Small chance (20%) for a Common Token even if you lose
		if rand.Float32() < 0.2 {
			items = append(items, inventory.Item{
				ID:      uuid.NewString(),
				UserID:  userID,
				ItemID:  inventory.EquipmentTokenCommonID,
				Rarity:  inventory.ItemRarityCommon,
				Created: s.now(),
			})
		}
	} else if poolID == "starter_pack" {
		// Starter Pack: Guaranteed Weapon and Armor
		equipment = append(equipment, inventory.Equipment{
			ID:      uuid.NewString(),
			UserID:  userID,
			ItemID:  "9C46EB15-4D04-4F90-B8A1-D4BD0A5A82B1", // Basic Sword
			Level:   1,
			Rarity:  inventory.ItemRarityCommon,
			Created: s.now(),
		}, inventory.Equipment{
			ID:      uuid.NewString(),
			UserID:  userID,
			ItemID:  "E98E3A82-C2B1-4A2D-A8C6-29BAE0D6A5A6", // Basic Vest
			Level:   1,
			Rarity:  inventory.ItemRarityCommon,
			Created: s.now(),
		})
	}

	return RewardContent{
		Items:     items,
		Equipment: equipment,
	}
}

func tokenIDForRarity(rarity int) string {
	switch rarity {
	case inventory.ItemRarityCommon:
		return inventory.EquipmentTokenCommonID
	case inventory.ItemRarityRare:
		return inventory.EquipmentTokenRareID
	case inventory.ItemRarityFabled:
		return inventory.EquipmentTokenFabledID
	case inventory.ItemRarityMythic:
		return inventory.EquipmentTokenMythicID
	case inventory.ItemRarityLegendary:
		return inventory.EquipmentTokenLegendaryID
	default:
		return inventory.EquipmentTokenCommonID
	}
}

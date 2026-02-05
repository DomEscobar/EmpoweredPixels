package inventory

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"empoweredpixels/internal/domain/inventory"
	"github.com/google/uuid"
)

var (
	ErrInvalidEquipment = errors.New("invalid equipment")
)

type Service struct {
	items     ItemRepository
	equipment EquipmentRepository
	options   EquipmentOptionRepository
	now       func() time.Time
}

func NewService(items ItemRepository, equipment EquipmentRepository, options EquipmentOptionRepository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		items:     items,
		equipment: equipment,
		options:   options,
		now:       now,
	}
}

func (s *Service) Balance(ctx context.Context, userID int64, itemID string) (int, error) {
	return s.items.CountByUserAndItemID(ctx, userID, itemID)
}

func (s *Service) GetEquipment(ctx context.Context, userID int64, id string) (*inventory.Equipment, *inventory.EquipmentOption, error) {
	equip, err := s.equipment.GetByID(ctx, userID, id)
	if err != nil {
		return nil, nil, err
	}
	if equip == nil {
		return nil, nil, ErrInvalidEquipment
	}

	option, err := s.options.GetByEquipmentID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if option == nil {
		option = &inventory.EquipmentOption{EquipmentID: id, IsFavorite: false}
		if err := s.options.Upsert(ctx, option); err != nil {
			return nil, nil, err
		}
	}

	return equip, option, nil
}

func (s *Service) EnhancementCost(current int, desired int) int {
	if desired <= current {
		return 0
	}
	sum := desired * (desired + 1) / 2
	return sum * 25
}

func (s *Service) Enhance(ctx context.Context, userID int64, equipmentID string, desired int) (*inventory.Equipment, error) {
	equip, err := s.equipment.GetByID(ctx, userID, equipmentID)
	if err != nil {
		return nil, err
	}
	if equip == nil {
		return nil, ErrInvalidEquipment
	}

	cost := s.EnhancementCost(equip.Enhancement, desired)
	if cost > 0 {
		removed, err := s.items.DeleteByUserAndItemID(ctx, userID, inventory.EmpoweredParticleID, cost)
		if err != nil {
			return nil, err
		}
		if removed != cost {
			return nil, ErrInvalidEquipment
		}
	}

	if desired > equip.Enhancement {
		equip.Enhancement = desired
		if err := s.equipment.UpdateEnhancement(ctx, equip.ID, desired); err != nil {
			return nil, err
		}
	}

	return equip, nil
}

func (s *Service) Salvage(ctx context.Context, userID int64, equipmentID string) ([]inventory.Item, error) {
	equip, err := s.equipment.GetByID(ctx, userID, equipmentID)
	if err != nil {
		return nil, err
	}
	if equip == nil {
		return nil, ErrInvalidEquipment
	}

	items := s.buildSalvageItems(userID, equip.Rarity)
	if err := s.items.CreateMany(ctx, items); err != nil {
		return nil, err
	}
	if err := s.equipment.Delete(ctx, equip.ID); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) SalvageInventory(ctx context.Context, userID int64) ([]inventory.Item, error) {
	equipmentList, err := s.equipment.ListInventoryAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	var items []inventory.Item
	for _, equip := range equipmentList {
		items = append(items, s.buildSalvageItems(userID, equip.Rarity)...)
		if err := s.equipment.Delete(ctx, equip.ID); err != nil {
			return nil, err
		}
	}

	if len(items) > 0 {
		if err := s.items.CreateMany(ctx, items); err != nil {
			return nil, err
		}
	}

	return items, nil
}

func (s *Service) InventoryPage(ctx context.Context, userID int64, page int, pageSize int) ([]inventory.Equipment, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.equipment.ListInventory(ctx, userID, pageSize, offset)
}

func (s *Service) ListByFighter(ctx context.Context, userID int64, fighterID string) ([]inventory.Equipment, error) {
	return s.equipment.ListByFighter(ctx, userID, fighterID)
}

func (s *Service) Equip(ctx context.Context, userID int64, equipmentID string, fighterID *string) error {
	equip, err := s.equipment.GetByID(ctx, userID, equipmentID)
	if err != nil {
		return err
	}
	if equip == nil {
		return ErrInvalidEquipment
	}

	return s.equipment.UpdateFighter(ctx, equipmentID, fighterID)
}

func (s *Service) SetFavorite(ctx context.Context, userID int64, equipmentID string, favorite bool) (*inventory.EquipmentOption, error) {
	equip, err := s.equipment.GetByID(ctx, userID, equipmentID)
	if err != nil {
		return nil, err
	}
	if equip == nil {
		return nil, ErrInvalidEquipment
	}

	option := &inventory.EquipmentOption{
		EquipmentID: equipmentID,
		IsFavorite:  favorite,
	}
	if err := s.options.Upsert(ctx, option); err != nil {
		return nil, err
	}
	return option, nil
}

func (s *Service) buildSalvageItems(userID int64, rarity int) []inventory.Item {
	if rarity == inventory.ItemRarityBasic {
		return []inventory.Item{}
	}

	value := rarity * 40
	items := make([]inventory.Item, 0, value+1)
	for i := 0; i < value; i++ {
		items = append(items, inventory.Item{
			ID:      uuid.NewString(),
			UserID:  userID,
			ItemID:  inventory.EmpoweredParticleID,
			Rarity:  inventory.ItemRarityBasic,
			Created: s.now(),
		})
	}

	if rarity != inventory.ItemRarityBasic && rand.Intn(100) < 33 {
		items = append(items, inventory.Item{
			ID:      uuid.NewString(),
			UserID:  userID,
			ItemID:  tokenIDForRarity(rarity),
			Rarity:  rarity,
			Created: s.now(),
		})
	}

	return items
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

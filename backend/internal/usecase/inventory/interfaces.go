package inventory

import (
	"context"

	"empoweredpixels/internal/domain/inventory"
)

type ItemRepository interface {
	CountByUserAndItemID(ctx context.Context, userID int64, itemID string) (int, error)
	CreateMany(ctx context.Context, items []inventory.Item) error
	DeleteByUserAndItemID(ctx context.Context, userID int64, itemID string, limit int) (int, error)
}

type EquipmentRepository interface {
	GetByID(ctx context.Context, userID int64, id string) (*inventory.Equipment, error)
	ListInventory(ctx context.Context, userID int64, limit int, offset int) ([]inventory.Equipment, error)
	ListInventoryAll(ctx context.Context, userID int64) ([]inventory.Equipment, error)
	ListByFighter(ctx context.Context, userID int64, fighterID string) ([]inventory.Equipment, error)
	UpdateEnhancement(ctx context.Context, equipmentID string, enhancement int) error
	Delete(ctx context.Context, equipmentID string) error
}

type EquipmentOptionRepository interface {
	GetByEquipmentID(ctx context.Context, equipmentID string) (*inventory.EquipmentOption, error)
	Upsert(ctx context.Context, option *inventory.EquipmentOption) error
}

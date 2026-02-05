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
	UpdateFighter(ctx context.Context, equipmentID string, fighterID *string) error
	Delete(ctx context.Context, equipmentID string) error
}

type Service interface {
	Balance(ctx context.Context, userID int64, itemID string) (int, error)
	GetEquipment(ctx context.Context, userID int64, id string) (*inventory.Equipment, *inventory.EquipmentOption, error)
	Enhance(ctx context.Context, userID int64, equipmentID string, desired int) (*inventory.Equipment, error)
	Salvage(ctx context.Context, userID int64, equipmentID string) ([]inventory.Item, error)
	SalvageInventory(ctx context.Context, userID int64) ([]inventory.Item, error)
	InventoryPage(ctx context.Context, userID int64, page int, pageSize int) ([]inventory.Equipment, error)
	ListByFighter(ctx context.Context, userID int64, fighterID string) ([]inventory.Equipment, error)
	SetFavorite(ctx context.Context, userID int64, equipmentID string, favorite bool) (*inventory.EquipmentOption, error)
	Equip(ctx context.Context, userID int64, equipmentID string, fighterID *string) error
}

type EquipmentOptionRepository interface {
	GetByEquipmentID(ctx context.Context, equipmentID string) (*inventory.EquipmentOption, error)
	Upsert(ctx context.Context, option *inventory.EquipmentOption) error
}

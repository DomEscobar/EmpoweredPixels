package repositories

import (
	"context"
	"errors"

	"empoweredpixels/internal/domain/inventory"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository struct {
	pool *pgxpool.Pool
}

type EquipmentRepository struct {
	pool *pgxpool.Pool
}

type EquipmentOptionRepository struct {
	pool *pgxpool.Pool
}

func NewItemRepository(pool *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{pool: pool}
}

func NewEquipmentRepository(pool *pgxpool.Pool) *EquipmentRepository {
	return &EquipmentRepository{pool: pool}
}

func NewEquipmentOptionRepository(pool *pgxpool.Pool) *EquipmentOptionRepository {
	return &EquipmentOptionRepository{pool: pool}
}

func (r *ItemRepository) CountByUserAndItemID(ctx context.Context, userID int64, itemID string) (int, error) {
	const query = `
		select count(*)
		from items
		where user_id = $1 and item_id = $2`

	var count int
	err := r.pool.QueryRow(ctx, query, userID, itemID).Scan(&count)
	return count, err
}

func (r *ItemRepository) CreateMany(ctx context.Context, items []inventory.Item) error {
	const query = `
		insert into items (id, user_id, item_id, rarity, created)
		values ($1, $2, $3, $4, $5)`

	batch := &pgx.Batch{}
	for _, item := range items {
		batch.Queue(query, item.ID, item.UserID, item.ItemID, item.Rarity, item.Created)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	_, err := br.Exec()
	return err
}

func (r *ItemRepository) DeleteByUserAndItemID(ctx context.Context, userID int64, itemID string, limit int) (int, error) {
	const query = `
		delete from items
		where id in (
			select id from items
			where user_id = $1 and item_id = $2
			limit $3
		)`

	result, err := r.pool.Exec(ctx, query, userID, itemID, limit)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}

func (r *EquipmentRepository) GetByID(ctx context.Context, userID int64, id string) (*inventory.Equipment, error) {
	const query = `
		select id, user_id, fighter_id, item_id, level, rarity, enhancement, created
		from equipment
		where id = $1 and user_id = $2`

	var equip inventory.Equipment
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&equip.ID, &equip.UserID, &equip.FighterID, &equip.ItemID, &equip.Level, &equip.Rarity, &equip.Enhancement, &equip.Created,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &equip, nil
}

func (r *EquipmentRepository) ListInventory(ctx context.Context, userID int64, limit int, offset int) ([]inventory.Equipment, error) {
	const query = `
		select id, user_id, fighter_id, item_id, level, rarity, enhancement, created
		from equipment
		where user_id = $1 and fighter_id is null
		order by level desc, rarity desc
		limit $2 offset $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipment []inventory.Equipment
	for rows.Next() {
		var equip inventory.Equipment
		if err := rows.Scan(
			&equip.ID, &equip.UserID, &equip.FighterID, &equip.ItemID, &equip.Level, &equip.Rarity, &equip.Enhancement, &equip.Created,
		); err != nil {
			return nil, err
		}
		equipment = append(equipment, equip)
	}
	return equipment, rows.Err()
}

func (r *EquipmentRepository) ListInventoryAll(ctx context.Context, userID int64) ([]inventory.Equipment, error) {
	const query = `
		select id, user_id, fighter_id, item_id, level, rarity, enhancement, created
		from equipment
		where user_id = $1 and fighter_id is null`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipment []inventory.Equipment
	for rows.Next() {
		var equip inventory.Equipment
		if err := rows.Scan(
			&equip.ID, &equip.UserID, &equip.FighterID, &equip.ItemID, &equip.Level, &equip.Rarity, &equip.Enhancement, &equip.Created,
		); err != nil {
			return nil, err
		}
		equipment = append(equipment, equip)
	}
	return equipment, rows.Err()
}

func (r *EquipmentRepository) ListByFighter(ctx context.Context, userID int64, fighterID string) ([]inventory.Equipment, error) {
	const query = `
		select id, user_id, fighter_id, item_id, level, rarity, enhancement, created
		from equipment
		where user_id = $1 and fighter_id = $2
		order by item_id`

	rows, err := r.pool.Query(ctx, query, userID, fighterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipment []inventory.Equipment
	for rows.Next() {
		var equip inventory.Equipment
		if err := rows.Scan(
			&equip.ID, &equip.UserID, &equip.FighterID, &equip.ItemID, &equip.Level, &equip.Rarity, &equip.Enhancement, &equip.Created,
		); err != nil {
			return nil, err
		}
		equipment = append(equipment, equip)
	}
	return equipment, rows.Err()
}

func (r *EquipmentRepository) UpdateEnhancement(ctx context.Context, equipmentID string, enhancement int) error {
	const query = `
		update equipment
		set enhancement = $1
		where id = $2`

	_, err := r.pool.Exec(ctx, query, enhancement, equipmentID)
	return err
}

func (r *EquipmentRepository) UpdateFighter(ctx context.Context, equipmentID string, fighterID *string) error {
	const query = `
		update equipment
		set fighter_id = $1
		where id = $2`

	_, err := r.pool.Exec(ctx, query, fighterID, equipmentID)
	return err
}

func (r *EquipmentRepository) Delete(ctx context.Context, equipmentID string) error {
	const query = `
		delete from equipment
		where id = $1`

	_, err := r.pool.Exec(ctx, query, equipmentID)
	return err
}

func (r *EquipmentRepository) Create(ctx context.Context, equipment *inventory.Equipment) error {
	const query = `
		insert into equipment (id, user_id, fighter_id, item_id, level, rarity, enhancement, created)
		values ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.pool.Exec(ctx, query,
		equipment.ID,
		equipment.UserID,
		equipment.FighterID,
		equipment.ItemID,
		equipment.Level,
		equipment.Rarity,
		equipment.Enhancement,
		equipment.Created,
	)
	return err
}

func (r *EquipmentOptionRepository) GetByEquipmentID(ctx context.Context, equipmentID string) (*inventory.EquipmentOption, error) {
	const query = `
		select equipment_id, is_favorite
		from equipment_options
		where equipment_id = $1`

	var option inventory.EquipmentOption
	err := r.pool.QueryRow(ctx, query, equipmentID).Scan(&option.EquipmentID, &option.IsFavorite)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &option, nil
}

func (r *EquipmentOptionRepository) Upsert(ctx context.Context, option *inventory.EquipmentOption) error {
	const query = `
		insert into equipment_options (equipment_id, is_favorite)
		values ($1, $2)
		on conflict (equipment_id)
		do update set is_favorite = excluded.is_favorite`

	_, err := r.pool.Exec(ctx, query, option.EquipmentID, option.IsFavorite)
	return err
}

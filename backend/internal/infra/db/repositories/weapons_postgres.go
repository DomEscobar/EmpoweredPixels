package repositories

import (
	"context"
	"errors"

	"empoweredpixels/internal/domain/weapons"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WeaponRepository struct {
	pool *pgxpool.Pool
}

func NewWeaponRepository(pool *pgxpool.Pool) *WeaponRepository {
	return &WeaponRepository{pool: pool}
}

func (r *WeaponRepository) GetByID(ctx context.Context, userID int64, id string) (*weapons.UserWeapon, error) {
	const query = `
		select id, user_id, weapon_id, enhancement, durability, fighter_id, created
		from user_weapons
		where id = $1 and user_id = $2`

	var uw weapons.UserWeapon
	var fighterID *string
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&uw.ID, &uw.UserID, &uw.WeaponID, &uw.Enhancement, &uw.Durability, &fighterID, &uw.Created,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	uw.FighterID = fighterID
	uw.IsEquipped = fighterID != nil
	return &uw, nil
}

func (r *WeaponRepository) ListByUser(ctx context.Context, userID int64, limit int, offset int) ([]weapons.UserWeapon, error) {
	const query = `
		select id, user_id, weapon_id, enhancement, durability, fighter_id, created
		from user_weapons
		where user_id = $1
		order by created desc
		limit $2 offset $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []weapons.UserWeapon
	for rows.Next() {
		var uw weapons.UserWeapon
		var fighterID *string
		if err := rows.Scan(&uw.ID, &uw.UserID, &uw.WeaponID, &uw.Enhancement, &uw.Durability, &fighterID, &uw.Created); err != nil {
			return nil, err
		}
		uw.FighterID = fighterID
		uw.IsEquipped = fighterID != nil
		result = append(result, uw)
	}
	return result, rows.Err()
}

func (r *WeaponRepository) ListByUserAll(ctx context.Context, userID int64) ([]weapons.UserWeapon, error) {
	const query = `
		select id, user_id, weapon_id, enhancement, durability, fighter_id, created
		from user_weapons
		where user_id = $1
		order by created desc`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []weapons.UserWeapon
	for rows.Next() {
		var uw weapons.UserWeapon
		var fighterID *string
		if err := rows.Scan(&uw.ID, &uw.UserID, &uw.WeaponID, &uw.Enhancement, &uw.Durability, &fighterID, &uw.Created); err != nil {
			return nil, err
		}
		uw.FighterID = fighterID
		uw.IsEquipped = fighterID != nil
		result = append(result, uw)
	}
	return result, rows.Err()
}

func (r *WeaponRepository) CountByUser(ctx context.Context, userID int64) (int, error) {
	const query = `
		select count(*)
		from user_weapons
		where user_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *WeaponRepository) Create(ctx context.Context, userWeapon *weapons.UserWeapon) error {
	const query = `
		insert into user_weapons (id, user_id, weapon_id, enhancement, durability, fighter_id, created)
		values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		userWeapon.ID,
		userWeapon.UserID,
		userWeapon.WeaponID,
		userWeapon.Enhancement,
		userWeapon.Durability,
		userWeapon.FighterID,
		userWeapon.Created,
	)
	return err
}

func (r *WeaponRepository) UpdateEnhancement(ctx context.Context, id string, enhancement int) error {
	const query = `
		update user_weapons
		set enhancement = $1
		where id = $2`

	_, err := r.pool.Exec(ctx, query, enhancement, id)
	return err
}

func (r *WeaponRepository) UpdateFighter(ctx context.Context, id string, fighterID *string) error {
	const query = `
		update user_weapons
		set fighter_id = $1
		where id = $2`

	_, err := r.pool.Exec(ctx, query, fighterID, id)
	return err
}

func (r *WeaponRepository) Delete(ctx context.Context, id string) error {
	const query = `
		delete from user_weapons
		where id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *WeaponRepository) GetEquippedByFighter(ctx context.Context, fighterID string) (*weapons.UserWeapon, error) {
	const query = `
		select id, user_id, weapon_id, enhancement, durability, fighter_id, created
		from user_weapons
		where fighter_id = $1`

	var uw weapons.UserWeapon
	var fid *string
	err := r.pool.QueryRow(ctx, query, fighterID).Scan(
		&uw.ID, &uw.UserID, &uw.WeaponID, &uw.Enhancement, &uw.Durability, &fid, &uw.Created,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	uw.FighterID = fid
	uw.IsEquipped = true
	return &uw, nil
}

// InventorySlotRepository manages inventory slot tracking
type WeaponInventoryRepository struct {
	pool *pgxpool.Pool
}

func NewWeaponInventoryRepository(pool *pgxpool.Pool) *WeaponInventoryRepository {
	return &WeaponInventoryRepository{pool: pool}
}

func (r *WeaponInventoryRepository) GetOrCreate(ctx context.Context, userID int64) (*weapons.InventorySlot, error) {
	const query = `
		insert into weapon_inventory (user_id, slot_count, used_slots)
		values ($1, 50, 0)
		on conflict (user_id) do update set slot_count = weapon_inventory.slot_count
		returning slot_count, used_slots`

	var slot weapons.InventorySlot
	err := r.pool.QueryRow(ctx, query, userID).Scan(&slot.SlotNumber, &slot.SlotNumber)
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *WeaponInventoryRepository) UpdateUsedSlots(ctx context.Context, userID int64, usedSlots int) error {
	const query = `
		insert into weapon_inventory (user_id, slot_count, used_slots)
		values ($1, 50, $2)
		on conflict (user_id) do update set used_slots = $2`

	_, err := r.pool.Exec(ctx, query, userID, usedSlots)
	return err
}
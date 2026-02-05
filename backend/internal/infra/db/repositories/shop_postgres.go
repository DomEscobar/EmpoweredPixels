package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"empoweredpixels/internal/domain/shop"
)

type ShopPostgres struct {
	db *sql.DB
}

func NewShopPostgres(db *sql.DB) *ShopPostgres {
	return &ShopPostgres{db: db}
}

// Shop methods

func (r *ShopPostgres) GetShop(ctx context.Context, id string) (*shop.Shop, error) {
	query := `SELECT id, name, description, currency, is_active, created FROM shops WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var s shop.Shop
	err := row.Scan(&s.ID, &s.Name, &s.Description, &s.Currency, &s.IsActive, &s.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShopPostgres) ListShops(ctx context.Context) ([]shop.Shop, error) {
	query := `SELECT id, name, description, currency, is_active, created FROM shops WHERE is_active = true ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shops []shop.Shop
	for rows.Next() {
		var s shop.Shop
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Currency, &s.IsActive, &s.Created); err != nil {
			return nil, err
		}
		shops = append(shops, s)
	}
	return shops, rows.Err()
}

// ShopItem methods

func (r *ShopPostgres) GetShopItem(ctx context.Context, id string) (*shop.ShopItem, error) {
	query := `SELECT id, shop_id, name, description, item_type, price_amount, price_currency, 
              gold_amount, rarity, image_url, metadata, is_active, sort_order, created 
              FROM shop_items WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var item shop.ShopItem
	var goldAmount sql.NullInt64
	var imageURL sql.NullString
	var metadataJSON []byte

	err := row.Scan(&item.ID, &item.ShopID, &item.Name, &item.Description, &item.ItemType,
		&item.PriceAmount, &item.PriceCurrency, &goldAmount, &item.Rarity, &imageURL,
		&metadataJSON, &item.IsActive, &item.SortOrder, &item.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if goldAmount.Valid {
		ga := int(goldAmount.Int64)
		item.GoldAmount = &ga
	}
	if imageURL.Valid {
		item.ImageURL = &imageURL.String
	}
	if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
		item.Metadata = make(map[string]interface{})
	}

	return &item, nil
}

func (r *ShopPostgres) ListShopItems(ctx context.Context, shopID string) ([]shop.ShopItem, error) {
	query := `SELECT id, shop_id, name, description, item_type, price_amount, price_currency, 
              gold_amount, rarity, image_url, metadata, is_active, sort_order, created 
              FROM shop_items WHERE shop_id = $1 AND is_active = true ORDER BY sort_order, name`
	rows, err := r.db.QueryContext(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanShopItems(rows)
}

func (r *ShopPostgres) ListAllActiveItems(ctx context.Context) ([]shop.ShopItem, error) {
	query := `SELECT id, shop_id, name, description, item_type, price_amount, price_currency, 
              gold_amount, rarity, image_url, metadata, is_active, sort_order, created 
              FROM shop_items WHERE is_active = true ORDER BY sort_order, name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanShopItems(rows)
}

func (r *ShopPostgres) ListItemsByType(ctx context.Context, itemType string) ([]shop.ShopItem, error) {
	query := `SELECT id, shop_id, name, description, item_type, price_amount, price_currency, 
              gold_amount, rarity, image_url, metadata, is_active, sort_order, created 
              FROM shop_items WHERE item_type = $1 AND is_active = true ORDER BY sort_order, name`
	rows, err := r.db.QueryContext(ctx, query, itemType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanShopItems(rows)
}

func (r *ShopPostgres) scanShopItems(rows *sql.Rows) ([]shop.ShopItem, error) {
	var items []shop.ShopItem
	for rows.Next() {
		var item shop.ShopItem
		var goldAmount sql.NullInt64
		var imageURL sql.NullString
		var metadataJSON []byte

		if err := rows.Scan(&item.ID, &item.ShopID, &item.Name, &item.Description, &item.ItemType,
			&item.PriceAmount, &item.PriceCurrency, &goldAmount, &item.Rarity, &imageURL,
			&metadataJSON, &item.IsActive, &item.SortOrder, &item.Created); err != nil {
			return nil, err
		}

		if goldAmount.Valid {
			ga := int(goldAmount.Int64)
			item.GoldAmount = &ga
		}
		if imageURL.Valid {
			item.ImageURL = &imageURL.String
		}
		if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
			item.Metadata = make(map[string]interface{})
		}

		items = append(items, item)
	}
	return items, rows.Err()
}

// PlayerGold methods

func (r *ShopPostgres) GetPlayerGold(ctx context.Context, userID int64) (*shop.PlayerGold, error) {
	query := `SELECT user_id, balance, lifetime_earned, lifetime_spent, updated FROM player_gold WHERE user_id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)

	var pg shop.PlayerGold
	err := row.Scan(&pg.UserID, &pg.Balance, &pg.LifetimeEarned, &pg.LifetimeSpent, &pg.Updated)
	if err == sql.ErrNoRows {
		// Return default balance of 0
		return &shop.PlayerGold{
			UserID:         userID,
			Balance:        0,
			LifetimeEarned: 0,
			LifetimeSpent:  0,
			Updated:        time.Now(),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &pg, nil
}

func (r *ShopPostgres) AddGold(ctx context.Context, userID int64, amount int) error {
	query := `INSERT INTO player_gold (user_id, balance, lifetime_earned, lifetime_spent, updated)
              VALUES ($1, $2, $2, 0, $3)
              ON CONFLICT (user_id) DO UPDATE SET 
                balance = player_gold.balance + $2,
                lifetime_earned = player_gold.lifetime_earned + $2,
                updated = $3`
	_, err := r.db.ExecContext(ctx, query, userID, amount, time.Now())
	return err
}

func (r *ShopPostgres) SpendGold(ctx context.Context, userID int64, amount int) error {
	query := `UPDATE player_gold SET 
              balance = balance - $2,
              lifetime_spent = lifetime_spent + $2,
              updated = $3
              WHERE user_id = $1 AND balance >= $2`
	result, err := r.db.ExecContext(ctx, query, userID, amount, time.Now())
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrInsufficientGold
	}
	return nil
}

func (r *ShopPostgres) SetGold(ctx context.Context, userID int64, balance int) error {
	query := `INSERT INTO player_gold (user_id, balance, lifetime_earned, lifetime_spent, updated)
              VALUES ($1, $2, 0, 0, $3)
              ON CONFLICT (user_id) DO UPDATE SET balance = $2, updated = $3`
	_, err := r.db.ExecContext(ctx, query, userID, balance, time.Now())
	return err
}

// Transaction methods

func (r *ShopPostgres) CreateTransaction(ctx context.Context, tx *shop.Transaction) error {
	metadataJSON, err := json.Marshal(tx.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := `INSERT INTO transactions (id, user_id, shop_item_id, item_type, item_name, 
              price_amount, price_currency, gold_change, status, metadata, created)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = r.db.ExecContext(ctx, query, tx.ID, tx.UserID, tx.ShopItemID, tx.ItemType,
		tx.ItemName, tx.PriceAmount, tx.PriceCurrency, tx.GoldChange, tx.Status, metadataJSON, tx.Created)
	return err
}

func (r *ShopPostgres) GetTransaction(ctx context.Context, id string) (*shop.Transaction, error) {
	query := `SELECT id, user_id, shop_item_id, item_type, item_name, price_amount, 
              price_currency, gold_change, status, metadata, created 
              FROM transactions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var tx shop.Transaction
	var shopItemID sql.NullString
	var metadataJSON []byte

	err := row.Scan(&tx.ID, &tx.UserID, &shopItemID, &tx.ItemType, &tx.ItemName,
		&tx.PriceAmount, &tx.PriceCurrency, &tx.GoldChange, &tx.Status, &metadataJSON, &tx.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if shopItemID.Valid {
		tx.ShopItemID = &shopItemID.String
	}
	if err := json.Unmarshal(metadataJSON, &tx.Metadata); err != nil {
		tx.Metadata = make(map[string]interface{})
	}

	return &tx, nil
}

func (r *ShopPostgres) ListTransactions(ctx context.Context, userID int64, limit, offset int) ([]shop.Transaction, error) {
	query := `SELECT id, user_id, shop_item_id, item_type, item_name, price_amount, 
              price_currency, gold_change, status, metadata, created 
              FROM transactions WHERE user_id = $1 ORDER BY created DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []shop.Transaction
	for rows.Next() {
		var tx shop.Transaction
		var shopItemID sql.NullString
		var metadataJSON []byte

		if err := rows.Scan(&tx.ID, &tx.UserID, &shopItemID, &tx.ItemType, &tx.ItemName,
			&tx.PriceAmount, &tx.PriceCurrency, &tx.GoldChange, &tx.Status, &metadataJSON, &tx.Created); err != nil {
			return nil, err
		}

		if shopItemID.Valid {
			tx.ShopItemID = &shopItemID.String
		}
		if err := json.Unmarshal(metadataJSON, &tx.Metadata); err != nil {
			tx.Metadata = make(map[string]interface{})
		}

		transactions = append(transactions, tx)
	}
	return transactions, rows.Err()
}

// Custom errors
var ErrInsufficientGold = &InsufficientGoldError{}

type InsufficientGoldError struct{}

func (e *InsufficientGoldError) Error() string {
	return "insufficient gold balance"
}

package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"empoweredpixels/internal/domain/shop"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ShopRepository defines shop-related database operations
type ShopRepository interface {
	GetShops(ctx context.Context) ([]shop.Shop, error)
	GetShopByID(ctx context.Context, id int) (*shop.Shop, error)
	GetShopItems(ctx context.Context, shopID *int, itemType *string) ([]shop.ShopItem, error)
	GetShopItemByID(ctx context.Context, id int) (*shop.ShopItem, error)
}

// PlayerGoldRepository defines player gold operations
type PlayerGoldRepository interface {
	GetPlayerGold(ctx context.Context, userID int) (*shop.PlayerGold, error)
	AddGold(ctx context.Context, userID int, amount int) error
	SpendGold(ctx context.Context, userID int, amount int) error
	CreatePlayerGold(ctx context.Context, userID int) error
	ListAllBalances(ctx context.Context) ([]shop.PlayerGold, error)
}

// TransactionRepository defines transaction operations
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *shop.Transaction) (int, error)
	GetTransactionsByUser(ctx context.Context, userID int, limit int) ([]shop.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id int, status string) error
}

// ShopPostgres implements ShopRepository
type ShopPostgres struct {
	db *pgxpool.Pool
}

// NewShopRepository creates a new shop repository
func NewShopRepository(db *pgxpool.Pool) ShopRepository {
	return &ShopPostgres{db: db}
}

// GetShops retrieves all active shops
func (r *ShopPostgres) GetShops(ctx context.Context) ([]shop.Shop, error) {
	query := `
		SELECT id, name, description, shop_type, currency, is_active, sort_order, created
		FROM shops
		WHERE is_active = true
		ORDER BY sort_order ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query shops: %w", err)
	}
	defer rows.Close()

	var shops []shop.Shop
	for rows.Next() {
		var s shop.Shop
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Description, &s.ShopType,
			&s.Currency, &s.IsActive, &s.SortOrder, &s.Created,
		); err != nil {
			return nil, fmt.Errorf("failed to scan shop: %w", err)
		}
		shops = append(shops, s)
	}

	return shops, rows.Err()
}

// GetShopByID retrieves a shop by ID
func (r *ShopPostgres) GetShopByID(ctx context.Context, id int) (*shop.Shop, error) {
	query := `
		SELECT id, name, description, shop_type, currency, is_active, sort_order, created
		FROM shops
		WHERE id = $1
	`

	var s shop.Shop
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.Description, &s.ShopType,
		&s.Currency, &s.IsActive, &s.SortOrder, &s.Created,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	return &s, nil
}

// GetShopItems retrieves shop items with optional filtering
func (r *ShopPostgres) GetShopItems(ctx context.Context, shopID *int, itemType *string) ([]shop.ShopItem, error) {
	query := `
		SELECT id, shop_id, name, description, item_type, price_amount, price_currency,
		       gold_amount, rarity, image_url, metadata, is_active, sort_order, created
		FROM shop_items
		WHERE is_active = true
	`
	args := []interface{}{}
	argIdx := 1

	if shopID != nil {
		query += fmt.Sprintf(" AND shop_id = $%d", argIdx)
		args = append(args, *shopID)
		argIdx++
	}

	if itemType != nil {
		query += fmt.Sprintf(" AND item_type = $%d", argIdx)
		args = append(args, *itemType)
		argIdx++
	}

	query += " ORDER BY sort_order ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query shop items: %w", err)
	}
	defer rows.Close()

	var items []shop.ShopItem
	for rows.Next() {
		var item shop.ShopItem
		var metadata []byte
		if err := rows.Scan(
			&item.ID, &item.ShopID, &item.Name, &item.Description, &item.ItemType,
			&item.PriceAmount, &item.PriceCurrency, &item.GoldAmount, &item.Rarity,
			&item.ImageURL, &metadata, &item.IsActive, &item.SortOrder, &item.Created,
		); err != nil {
			return nil, fmt.Errorf("failed to scan shop item: %w", err)
		}
		if len(metadata) > 0 {
			json.Unmarshal(metadata, &item.Metadata)
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// GetShopItemByID retrieves a shop item by ID
func (r *ShopPostgres) GetShopItemByID(ctx context.Context, id int) (*shop.ShopItem, error) {
	query := `
		SELECT id, shop_id, name, description, item_type, price_amount, price_currency,
		       gold_amount, rarity, image_url, metadata, is_active, sort_order, created
		FROM shop_items
		WHERE id = $1
	`

	var item shop.ShopItem
	var metadata []byte
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.ShopID, &item.Name, &item.Description, &item.ItemType,
		&item.PriceAmount, &item.PriceCurrency, &item.GoldAmount, &item.Rarity,
		&item.ImageURL, &metadata, &item.IsActive, &item.SortOrder, &item.Created,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get shop item: %w", err)
	}

	if len(metadata) > 0 {
		json.Unmarshal(metadata, &item.Metadata)
	}

	return &item, nil
}

// PlayerGoldPostgres implements PlayerGoldRepository
type PlayerGoldPostgres struct {
	db *pgxpool.Pool
}

// NewPlayerGoldRepository creates a new player gold repository
func NewPlayerGoldRepository(db *pgxpool.Pool) PlayerGoldRepository {
	return &PlayerGoldPostgres{db: db}
}

// GetPlayerGold retrieves a player's gold balance
func (r *PlayerGoldPostgres) GetPlayerGold(ctx context.Context, userID int) (*shop.PlayerGold, error) {
	query := `
		SELECT user_id, balance, lifetime_earned, lifetime_spent, updated
		FROM player_gold
		WHERE user_id = $1
	`

	var pg shop.PlayerGold
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&pg.UserID, &pg.Balance, &pg.LifetimeEarned, &pg.LifetimeSpent, &pg.Updated,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Return zero balance if not found
			return &shop.PlayerGold{
				UserID:  userID,
				Balance: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get player gold: %w", err)
	}

	return &pg, nil
}

// AddGold adds gold to a player's balance
func (r *PlayerGoldPostgres) AddGold(ctx context.Context, userID int, amount int) error {
	query := `
		INSERT INTO player_gold (user_id, balance, lifetime_earned, lifetime_spent, updated)
		VALUES ($1, $2, $2, 0, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			balance = player_gold.balance + EXCLUDED.balance,
			lifetime_earned = player_gold.lifetime_earned + EXCLUDED.balance,
			updated = NOW()
	`

	_, err := r.db.Exec(ctx, query, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to add gold: %w", err)
	}

	return nil
}

// SpendGold subtracts gold from a player's balance
func (r *PlayerGoldPostgres) SpendGold(ctx context.Context, userID int, amount int) error {
	query := `
		UPDATE player_gold
		SET balance = balance - $2,
		    lifetime_spent = lifetime_spent + $2,
		    updated = NOW()
		WHERE user_id = $1 AND balance >= $2
	`

	result, err := r.db.Exec(ctx, query, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to spend gold: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("insufficient gold")
	}

	return nil
}

// CreatePlayerGold initializes a player's gold record
func (r *PlayerGoldPostgres) CreatePlayerGold(ctx context.Context, userID int) error {
	query := `
		INSERT INTO player_gold (user_id, balance, lifetime_earned, lifetime_spent, updated)
		VALUES ($1, 0, 0, 0, NOW())
		ON CONFLICT (user_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to create player gold: %w", err)
	}

	return nil
}

// ListAllBalances retrieves all player gold balances
func (r *PlayerGoldPostgres) ListAllBalances(ctx context.Context) ([]shop.PlayerGold, error) {
	query := `
		SELECT user_id, balance, lifetime_earned, lifetime_spent, updated
		FROM player_gold
		WHERE balance > 0
		ORDER BY balance DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list gold balances: %w", err)
	}
	defer rows.Close()

	var balances []shop.PlayerGold
	for rows.Next() {
		var g shop.PlayerGold
		if err := rows.Scan(&g.UserID, &g.Balance, &g.LifetimeEarned, &g.LifetimeSpent, &g.Updated); err != nil {
			return nil, err
		}
		balances = append(balances, g)
	}

	return balances, rows.Err()
}

// TransactionPostgres implements TransactionRepository
type TransactionPostgres struct {
	db *pgxpool.Pool
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &TransactionPostgres{db: db}
}

// CreateTransaction creates a new transaction record
func (r *TransactionPostgres) CreateTransaction(ctx context.Context, tx *shop.Transaction) (int, error) {
	query := `
		INSERT INTO transactions (user_id, shop_item_id, item_type, item_name, price_amount, price_currency, gold_change, status, metadata, created)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id
	`

	metadata, _ := json.Marshal(tx.Metadata)

	var id int
	err := r.db.QueryRow(ctx, query,
		tx.UserID, tx.ShopItemID, tx.ItemType, tx.ItemName,
		tx.PriceAmount, tx.PriceCurrency, tx.GoldChange, tx.Status, metadata,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}

	return id, nil
}

// GetTransactionsByUser retrieves transactions for a user
func (r *TransactionPostgres) GetTransactionsByUser(ctx context.Context, userID int, limit int) ([]shop.Transaction, error) {
	query := `
		SELECT id, user_id, shop_item_id, item_type, item_name, price_amount, price_currency, gold_change, status, metadata, created
		FROM transactions
		WHERE user_id = $1
		ORDER BY created DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []shop.Transaction
	for rows.Next() {
		var tx shop.Transaction
		var metadata []byte
		if err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.ShopItemID, &tx.ItemType, &tx.ItemName,
			&tx.PriceAmount, &tx.PriceCurrency, &tx.GoldChange, &tx.Status,
			&metadata, &tx.Created,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		if len(metadata) > 0 {
			json.Unmarshal(metadata, &tx.Metadata)
		}
		transactions = append(transactions, tx)
	}

	return transactions, rows.Err()
}

// UpdateTransactionStatus updates a transaction's status
func (r *TransactionPostgres) UpdateTransactionStatus(ctx context.Context, id int, status string) error {
	query := `
		UPDATE transactions
		SET status = $2
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}

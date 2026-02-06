-- Migration: Create shop tables
-- Created: 2026-02-05

-- Shops table (shop categories)
CREATE TABLE IF NOT EXISTS shops (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    shop_type VARCHAR(50) NOT NULL DEFAULT 'bundles',
    currency VARCHAR(20) NOT NULL DEFAULT 'gold',
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_shop_type CHECK (shop_type IN ('gold', 'bundles', 'premium')),
    CONSTRAINT chk_shop_currency CHECK (currency IN ('usd', 'gold', 'particles'))
);

-- Shop items table (products)
CREATE TABLE IF NOT EXISTS shop_items (
    id SERIAL PRIMARY KEY,
    shop_id INTEGER NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    item_type VARCHAR(50) NOT NULL DEFAULT 'bundle',
    price_amount INTEGER NOT NULL DEFAULT 0,
    price_currency VARCHAR(20) NOT NULL DEFAULT 'gold',
    gold_amount INTEGER,
    rarity INTEGER NOT NULL DEFAULT 0,
    image_url VARCHAR(500),
    metadata JSONB DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_item_type CHECK (item_type IN ('gold_package', 'bundle', 'equipment', 'consumable')),
    CONSTRAINT chk_price_currency CHECK (price_currency IN ('usd', 'gold', 'particles')),
    CONSTRAINT chk_rarity CHECK (rarity >= 0 AND rarity <= 5),
    CONSTRAINT chk_price_positive CHECK (price_amount >= 0),
    UNIQUE (name, shop_id)
);

-- Player gold balance table
CREATE TABLE IF NOT EXISTS player_gold (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance INTEGER NOT NULL DEFAULT 0,
    lifetime_earned INTEGER NOT NULL DEFAULT 0,
    lifetime_spent INTEGER NOT NULL DEFAULT 0,
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_balance_non_negative CHECK (balance >= 0),
    CONSTRAINT chk_lifetime_earned_non_negative CHECK (lifetime_earned >= 0),
    CONSTRAINT chk_lifetime_spent_non_negative CHECK (lifetime_spent >= 0)
);

-- Transactions table (purchase history)
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shop_item_id INTEGER REFERENCES shop_items(id) ON DELETE SET NULL,
    item_type VARCHAR(50) NOT NULL,
    item_name VARCHAR(200) NOT NULL,
    price_amount INTEGER NOT NULL,
    price_currency VARCHAR(20) NOT NULL,
    gold_change INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    metadata JSONB DEFAULT '{}',
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_transaction_status CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    CONSTRAINT chk_transaction_price_positive CHECK (price_amount >= 0)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_shop_items_shop_id ON shop_items(shop_id);
CREATE INDEX IF NOT EXISTS idx_shop_items_active ON shop_items(is_active);
CREATE INDEX IF NOT EXISTS idx_shop_items_type ON shop_items(item_type);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_created ON transactions(user_id, created DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);

-- Seed data: Shops
INSERT INTO shops (name, description, shop_type, currency, is_active, sort_order) VALUES
('Gold Emporium', 'Purchase gold for upgrades and items', 'gold', 'usd', true, 1),
('Equipment Bundles', 'Curated bundles with guaranteed rarity items', 'bundles', 'gold', true, 2),
('Premium Store', 'Exclusive items and cosmetics', 'premium', 'usd', true, 3)
ON CONFLICT (name) DO NOTHING;

-- Seed data: Gold Packages (prices in cents)
INSERT INTO shop_items (shop_id, name, description, item_type, price_amount, price_currency, gold_amount, rarity, is_active, sort_order) VALUES
((SELECT id FROM shops WHERE name = 'Gold Emporium'), 'Small Pouch', '100 Gold for small upgrades', 'gold_package', 99, 'usd', 100, 0, true, 1),
((SELECT id FROM shops WHERE name = 'Gold Emporium'), 'Merchant Satchel', '550 Gold with 10% bonus - Most Popular!', 'gold_package', 499, 'usd', 550, 1, true, 2),
((SELECT id FROM shops WHERE name = 'Gold Emporium'), 'Treasure Chest', '1,200 Gold with 20% bonus', 'gold_package', 999, 'usd', 1200, 2, true, 3),
((SELECT id FROM shops WHERE name = 'Gold Emporium'), 'Dragon Hoard', '6,500 Gold with 30% bonus - Best Value!', 'gold_package', 4999, 'usd', 6500, 4, true, 4)
ON CONFLICT (name, shop_id) DO NOTHING;

-- Seed data: Equipment Bundles
INSERT INTO shop_items (shop_id, name, description, item_type, price_amount, price_currency, rarity, metadata, is_active, sort_order) VALUES
((SELECT id FROM shops WHERE name = 'Equipment Bundles'), 'Starter Bundle', 'Common weapon + 200 Gold for new fighters', 'bundle', 299, 'usd', 1, '{"equipment_count": 1, "gold_bonus": 200}', true, 1),
((SELECT id FROM shops WHERE name = 'Equipment Bundles'), 'Epic Hunter Pack', '5 Epic Drop Boosts + 500 Gold', 'bundle', 999, 'usd', 3, '{"drop_boosts": 5, "gold_bonus": 500}', true, 2),
((SELECT id FROM shops WHERE name = 'Equipment Bundles'), 'Legendary Crate', 'Guaranteed Legendary Weapon', 'bundle', 1999, 'usd', 5, '{"equipment_count": 1, "guaranteed_rarity": 5}', true, 3),
((SELECT id FROM shops WHERE name = 'Equipment Bundles'), 'Mythic Ascension', 'Mythic Weapon of Choice + 2,000 Gold', 'bundle', 4999, 'usd', 4, '{"equipment_count": 1, "guaranteed_rarity": 4, "gold_bonus": 2000}', true, 4)
ON CONFLICT (name, shop_id) DO NOTHING;

-- Initialize player_gold for existing users (optional migration helper)
-- INSERT INTO player_gold (user_id, balance, lifetime_earned, lifetime_spent)
-- SELECT id, 0, 0, 0 FROM users 
-- WHERE id NOT IN (SELECT user_id FROM player_gold)
-- ON CONFLICT DO NOTHING;

-- Shop MVP: shops, shop_items, player_gold, transactions

-- Shops table (shop types, currencies)
create table if not exists shops (
  id uuid primary key,
  name text not null unique,
  description text not null default '',
  currency text not null default 'gold',
  is_active boolean not null default true,
  created timestamptz not null default now()
);

-- Shop items table (products, prices, rarities)
create table if not exists shop_items (
  id uuid primary key,
  shop_id uuid not null references shops(id) on delete cascade,
  name text not null,
  description text not null default '',
  item_type text not null, -- 'gold_package', 'bundle', 'equipment', 'consumable'
  price_amount int not null,
  price_currency text not null default 'usd', -- 'usd', 'gold', 'particles'
  gold_amount int null, -- for gold packages
  rarity int not null default 0,
  image_url text null,
  metadata jsonb not null default '{}',
  is_active boolean not null default true,
  sort_order int not null default 0,
  created timestamptz not null default now()
);

-- Player gold balances
create table if not exists player_gold (
  user_id bigint primary key references users(id) on delete cascade,
  balance int not null default 0,
  lifetime_earned int not null default 0,
  lifetime_spent int not null default 0,
  updated timestamptz not null default now()
);

-- Transactions table (purchase history)
create table if not exists transactions (
  id uuid primary key,
  user_id bigint not null references users(id) on delete cascade,
  shop_item_id uuid references shop_items(id) on delete set null,
  item_type text not null,
  item_name text not null,
  price_amount int not null,
  price_currency text not null,
  gold_change int not null default 0,
  status text not null default 'completed', -- 'pending', 'completed', 'failed', 'refunded'
  metadata jsonb not null default '{}',
  created timestamptz not null default now()
);

-- Indexes for performance
create index if not exists idx_shop_items_shop_id on shop_items(shop_id);
create index if not exists idx_shop_items_active on shop_items(is_active, sort_order);
create index if not exists idx_transactions_user_id on transactions(user_id);
create index if not exists idx_transactions_created on transactions(created desc);

-- Seed default shops
insert into shops (id, name, description, currency) values
  ('11111111-1111-1111-1111-111111111111', 'Gold Store', 'Purchase gold packages with real money', 'usd'),
  ('22222222-2222-2222-2222-222222222222', 'Equipment Shop', 'Buy equipment bundles with gold', 'gold')
on conflict (name) do nothing;

-- Seed gold packages (100, 550, 1200, 6500)
insert into shop_items (id, shop_id, name, description, item_type, price_amount, price_currency, gold_amount, rarity, sort_order) values
  ('aaaa1111-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'Handful of Gold', 'A small pouch of gold coins', 'gold_package', 99, 'usd', 100, 0, 1),
  ('aaaa1111-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'Bag of Gold', 'A generous bag of gold coins (+10% bonus)', 'gold_package', 499, 'usd', 550, 0, 2),
  ('aaaa1111-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'Chest of Gold', 'A treasure chest full of gold (+20% bonus)', 'gold_package', 999, 'usd', 1200, 0, 3),
  ('aaaa1111-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'Vault of Gold', 'A legendary vault of riches (+30% bonus)', 'gold_package', 4999, 'usd', 6500, 1, 4)
on conflict do nothing;

-- Seed rarity bundles (Starter, Epic, Legendary, Mythic)
insert into shop_items (id, shop_id, name, description, item_type, price_amount, price_currency, gold_amount, rarity, metadata, sort_order) values
  ('bbbb2222-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'Starter Bundle', 'Perfect for new fighters. Contains 3 common equipment pieces.', 'bundle', 100, 'gold', null, 1, '{"equipment_count": 3, "min_rarity": 1, "max_rarity": 1}', 1),
  ('bbbb2222-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'Epic Bundle', 'For aspiring champions. Contains 2 rare equipment pieces.', 'bundle', 500, 'gold', null, 2, '{"equipment_count": 2, "min_rarity": 2, "max_rarity": 2}', 2),
  ('bbbb2222-0000-0000-0000-000000000003', '22222222-2222-2222-2222-222222222222', 'Legendary Bundle', 'Forge your legend. Contains 1 fabled equipment piece.', 'bundle', 1500, 'gold', null, 3, '{"equipment_count": 1, "min_rarity": 3, "max_rarity": 3}', 3),
  ('bbbb2222-0000-0000-0000-000000000004', '22222222-2222-2222-2222-222222222222', 'Mythic Bundle', 'The rarest treasures await. Contains 1 mythic equipment piece.', 'bundle', 5000, 'gold', null, 4, '{"equipment_count": 1, "min_rarity": 4, "max_rarity": 4}', 4)
on conflict do nothing;

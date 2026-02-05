create table if not exists user_weapons (
  id uuid primary key,
  user_id bigint not null references users(id) on delete cascade,
  weapon_id varchar(50) not null,
  enhancement int not null default 0,
  durability int not null default 100,
  fighter_id uuid null references fighters(id) on delete set null,
  created timestamptz not null default now()
);

create index if not exists idx_user_weapons_user_id on user_weapons(user_id);
create index if not exists idx_user_weapons_fighter_id on user_weapons(fighter_id);

-- Weapon inventory slot tracking (max 50 slots)
create table if not exists weapon_inventory (
  user_id bigint primary key references users(id) on delete cascade,
  slot_count int not null default 50,
  used_slots int not null default 0
);
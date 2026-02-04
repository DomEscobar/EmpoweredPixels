create table if not exists items (
  id uuid primary key,
  user_id bigint not null references users(id) on delete cascade,
  item_id uuid not null,
  rarity int not null default 0,
  created timestamptz not null
);

create table if not exists equipment (
  id uuid primary key,
  user_id bigint not null references users(id) on delete cascade,
  fighter_id uuid null references fighters(id) on delete set null,
  item_id uuid not null,
  level int not null default 1,
  rarity int not null default 0,
  enhancement int not null default 0,
  created timestamptz not null
);

create table if not exists equipment_options (
  equipment_id uuid primary key references equipment(id) on delete cascade,
  is_favorite boolean not null default false
);

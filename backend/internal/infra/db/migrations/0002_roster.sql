create table if not exists fighters (
  id uuid primary key,
  user_id bigint not null references users(id) on delete cascade,
  name text not null unique,
  level int not null default 1,
  power int not null default 0,
  condition_power int not null default 0,
  precision int not null default 0,
  ferocity int not null default 0,
  accuracy int not null default 0,
  agility int not null default 0,
  armor int not null default 0,
  vitality int not null default 0,
  parry_chance int not null default 0,
  healing_power int not null default 0,
  speed int not null default 0,
  vision int not null default 0,
  weapon_id text null,
  attunement_id text null,
  created timestamptz not null,
  is_deleted boolean not null default false
);

create table if not exists fighter_experiences (
  id bigserial primary key,
  fighter_id uuid not null unique references fighters(id) on delete cascade,
  experience int not null default 0
);

create table if not exists fighter_configurations (
  fighter_id uuid primary key references fighters(id) on delete cascade,
  attunement_id uuid null
);

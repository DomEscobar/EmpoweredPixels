create table if not exists leagues (
  id serial primary key,
  name text null,
  options jsonb not null default '{}'::jsonb,
  is_deactivated boolean not null default false
);

create table if not exists league_subscriptions (
  league_id int not null references leagues(id) on delete cascade,
  fighter_id uuid not null references fighters(id) on delete cascade,
  created timestamptz not null,
  primary key (league_id, fighter_id)
);

create table if not exists league_matches (
  league_id int not null references leagues(id) on delete cascade,
  match_id uuid not null references matches(id) on delete cascade,
  started timestamptz null,
  primary key (league_id, match_id)
);

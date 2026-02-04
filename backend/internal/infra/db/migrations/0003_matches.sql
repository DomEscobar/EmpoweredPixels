create table if not exists matches (
  id uuid primary key,
  creator_user_id bigint null references users(id) on delete set null,
  created timestamptz not null,
  started timestamptz null,
  options jsonb not null
);

create table if not exists match_teams (
  id uuid primary key,
  match_id uuid not null references matches(id) on delete cascade,
  password text null
);

create table if not exists match_registrations (
  match_id uuid not null references matches(id) on delete cascade,
  fighter_id uuid not null references fighters(id) on delete cascade,
  team_id uuid null references match_teams(id) on delete set null,
  date timestamptz not null,
  primary key (match_id, fighter_id)
);

create table if not exists match_results (
  id uuid primary key,
  match_id uuid not null unique references matches(id) on delete cascade,
  round_ticks jsonb not null
);

create table if not exists match_score_fighters (
  match_id uuid not null references matches(id) on delete cascade,
  fighter_id uuid not null references fighters(id) on delete cascade,
  total_kills int not null default 0,
  total_deaths int not null default 0,
  total_assists int not null default 0,
  primary key (match_id, fighter_id)
);

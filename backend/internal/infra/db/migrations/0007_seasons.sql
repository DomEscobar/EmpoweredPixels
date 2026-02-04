create table if not exists seasons (
  id bigserial primary key,
  season_id int not null,
  start_date timestamptz not null,
  end_date timestamptz not null
);

create table if not exists season_summaries (
  id bigserial primary key,
  user_id bigint not null references users(id) on delete cascade,
  season_id int not null,
  position int not null
);

create table if not exists season_progresses (
  id bigserial primary key,
  user_id bigint not null references users(id) on delete cascade,
  season_id int not null,
  is_complete boolean not null default false
);

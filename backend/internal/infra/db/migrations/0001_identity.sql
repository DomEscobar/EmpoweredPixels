create table if not exists users (
  id bigserial primary key,
  name text not null unique,
  email text not null unique,
  password text not null,
  salt text not null,
  is_verified boolean not null default false,
  created timestamptz not null,
  last_login timestamptz not null,
  banned timestamptz null
);

create table if not exists tokens (
  id uuid primary key,
  user_id bigint not null unique references users(id) on delete cascade,
  value text not null,
  refresh_value text not null,
  issued timestamptz not null
);

create table if not exists verifications (
  id uuid primary key,
  user_id bigint not null references users(id) on delete cascade
);

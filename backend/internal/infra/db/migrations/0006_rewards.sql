create table if not exists rewards (
  id uuid primary key,
  user_id bigint not null references users(id) on delete cascade,
  reward_pool_id uuid not null,
  claimed timestamptz null,
  created timestamptz not null
);

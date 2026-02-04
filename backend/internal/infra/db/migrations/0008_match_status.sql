alter table matches
  add column if not exists status text not null default 'lobby',
  add column if not exists completed_at timestamptz null,
  add column if not exists cancelled_at timestamptz null;

update matches
set status = 'completed',
  completed_at = started
where started is not null;

create index if not exists idx_matches_status on matches(status);

ALTER TABLE rewards ADD COLUMN IF NOT EXISTS source_id TEXT;
CREATE INDEX IF NOT EXISTS idx_rewards_source_id ON rewards(source_id);

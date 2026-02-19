ALTER TABLE rewards ADD COLUMN source_id TEXT;
CREATE INDEX idx_rewards_source_id ON rewards(source_id);

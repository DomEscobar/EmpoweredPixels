-- Migration: Add XP and Stats to Fighters
-- Created: 2026-02-06

-- Add XP column for fighter progression
ALTER TABLE fighters ADD COLUMN IF NOT EXISTS xp INTEGER NOT NULL DEFAULT 0;
ALTER TABLE fighters ADD COLUMN IF NOT EXISTS xp_to_next_level INTEGER NOT NULL DEFAULT 100;

-- Add match tracking statistics
ALTER TABLE fighters ADD COLUMN IF NOT EXISTS matches_won INTEGER NOT NULL DEFAULT 0;
ALTER TABLE fighters ADD COLUMN IF NOT EXISTS matches_lost INTEGER NOT NULL DEFAULT 0;
ALTER TABLE fighters ADD COLUMN IF NOT EXISTS total_matches INTEGER NOT NULL DEFAULT 0;
ALTER TABLE fighters ADD COLUMN IF NOT EXISTS total_damage_dealt BIGINT NOT NULL DEFAULT 0;
ALTER TABLE fighters ADD COLUMN IF NOT EXISTS total_damage_taken BIGINT NOT NULL DEFAULT 0;

-- Index for leaderboards
CREATE INDEX IF NOT EXISTS idx_fighters_xp ON fighters(xp DESC);
CREATE INDEX IF NOT EXISTS idx_fighters_total_matches ON fighters(total_matches DESC);
CREATE INDEX IF NOT EXISTS idx_fighters_matches_won ON fighters(matches_won DESC);

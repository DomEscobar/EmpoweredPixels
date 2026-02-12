-- Migration: Daily Rewards System
-- Created: 2026-02-06

CREATE TABLE IF NOT EXISTS daily_rewards (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    streak INTEGER NOT NULL DEFAULT 0,
    last_claimed DATE,
    total_claimed INTEGER NOT NULL DEFAULT 0,
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index for cleanup jobs
CREATE INDEX IF NOT EXISTS idx_daily_rewards_last_claimed ON daily_rewards(last_claimed);

-- Reward definitions (in code, but documented here)
-- Day 1: 100 Gold
-- Day 2: 250 Gold + Common Item
-- Day 3: 500 Gold + Rare Item
-- Day 4: 2x XP Boost (1 hour)
-- Day 5: Mystery Box (Random rarity 1-4)
-- Day 6: 1000 Gold + Fabled Item
-- Day 7: Legendary Crate (Guaranteed Legendary + 2000 Gold)

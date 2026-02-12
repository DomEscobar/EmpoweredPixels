-- Migration: Leaderboard System
-- Created: 2026-02-06

-- Player rankings cache (updated periodically)
CREATE TABLE IF NOT EXISTS leaderboard_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(50) NOT NULL, -- 'power', 'wealth', 'combat', 'achievements', 'streak'
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rank INTEGER NOT NULL,
    score BIGINT NOT NULL DEFAULT 0,
    previous_rank INTEGER,
    trend VARCHAR(10) GENERATED ALWAYS AS (
        CASE 
            WHEN previous_rank IS NULL THEN 'same'
            WHEN rank < previous_rank THEN 'up'
            WHEN rank > previous_rank THEN 'down'
            ELSE 'same'
        END
    ) STORED,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(category, user_id)
);

-- Indexes for leaderboard queries
CREATE INDEX IF NOT EXISTS idx_leaderboard_category_rank ON leaderboard_entries(category, rank);
CREATE INDEX IF NOT EXISTS idx_leaderboard_user ON leaderboard_entries(user_id);

-- Player achievements tracking
CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    icon VARCHAR(50) NOT NULL DEFAULT '🏆',
    category VARCHAR(50) NOT NULL, -- 'combat', 'collection', 'progression', 'social'
    requirement_type VARCHAR(50) NOT NULL, -- 'wins', 'matches', 'level', 'gold', etc.
    requirement_value INTEGER NOT NULL DEFAULT 1,
    reward_gold INTEGER NOT NULL DEFAULT 0,
    reward_title VARCHAR(100),
    hidden BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Player achievement progress
CREATE TABLE IF NOT EXISTS player_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    progress INTEGER NOT NULL DEFAULT 0,
    completed BOOLEAN NOT NULL DEFAULT false,
    completed_at TIMESTAMP WITH TIME ZONE,
    claimed BOOLEAN NOT NULL DEFAULT false,
    claimed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, achievement_id)
);

CREATE INDEX IF NOT EXISTS idx_player_achievements_user ON player_achievements(user_id);
CREATE INDEX IF NOT EXISTS idx_player_achievements_completed ON player_achievements(user_id, completed);

-- Insert default achievements
INSERT INTO achievements (key, name, description, icon, category, requirement_type, requirement_value, reward_gold, reward_title) VALUES
('first_blood', 'First Blood', 'Win your first match', '🩸', 'combat', 'wins', 1, 100, 'Rookie'),
('veteran', 'Veteran', 'Win 10 matches', '⚔️', 'combat', 'wins', 10, 500, 'Veteran'),
('champion', 'Champion', 'Win 50 matches', '👑', 'combat', 'wins', 50, 2000, 'Champion'),
('legend', 'Legend', 'Win 100 matches', '🔥', 'combat', 'wins', 100, 5000, 'Legend'),
('gladiator', 'Gladiator', 'Complete 10 matches', '🛡️', 'combat', 'matches', 10, 200, 'Gladiator'),
('warlord', 'Warlord', 'Complete 100 matches', '⚡', 'combat', 'matches', 100, 1000, 'Warlord'),
('collector', 'Collector', 'Own 10 different equipment pieces', '📦', 'collection', 'equipment', 10, 300, 'Collector'),
('hoarder', 'Hoarder', 'Own 50 different equipment pieces', '💎', 'collection', 'equipment', 50, 1000, 'Hoarder'),
('millionaire', 'Millionaire', 'Accumulate 1,000,000 gold', '💰', 'progression', 'gold', 1000000, 0, 'Millionaire'),
('daily_dedication', 'Daily Dedication', 'Claim 7 daily rewards in a row', '📅', 'progression', 'daily_streak', 7, 500, 'Dedicated')
ON CONFLICT (key) DO NOTHING;

-- Migration: Create attunement tables
-- Created: 2026-02-05

-- Player attunements table (6 elements, levels 1-25)
CREATE TABLE IF NOT EXISTS player_attunements (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    element VARCHAR(20) NOT NULL,
    level INTEGER NOT NULL DEFAULT 1,
    current_xp INTEGER NOT NULL DEFAULT 0,
    total_xp INTEGER NOT NULL DEFAULT 0,
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, element),
    CONSTRAINT chk_element CHECK (element IN ('fire', 'water', 'earth', 'air', 'light', 'dark')),
    CONSTRAINT chk_level CHECK (level >= 1 AND level <= 25),
    CONSTRAINT chk_xp_non_negative CHECK (current_xp >= 0 AND total_xp >= 0)
);

-- Attunement XP history (audit trail)
CREATE TABLE IF NOT EXISTS attunement_xp_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    element VARCHAR(20) NOT NULL,
    xp_gained INTEGER NOT NULL,
    source VARCHAR(100) NOT NULL, -- match_win, daily_task, element_use, etc.
    source_id INTEGER, -- optional: match_id, task_id, etc.
    new_level INTEGER,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_xp_positive CHECK (xp_gained > 0),
    CONSTRAINT chk_element_valid CHECK (element IN ('fire', 'water', 'earth', 'air', 'light', 'dark'))
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_player_attunements_user ON player_attunements(user_id);
CREATE INDEX IF NOT EXISTS idx_attunement_xp_user ON attunement_xp_history(user_id);
CREATE INDEX IF NOT EXISTS idx_attunement_xp_element ON attunement_xp_history(element);
CREATE INDEX IF NOT EXISTS idx_attunement_xp_created ON attunement_xp_history(created DESC);

-- Initialize attunements for existing users (optional)
-- All 6 elements at level 1 with 0 XP
INSERT INTO player_attunements (user_id, element, level, current_xp, total_xp)
SELECT 
    u.id,
    e.element,
    1,
    0,
    0
FROM users u
CROSS JOIN (VALUES ('fire'), ('water'), ('earth'), ('air'), ('light'), ('dark')) AS e(element)
WHERE NOT EXISTS (
    SELECT 1 FROM player_attunements pa 
    WHERE pa.user_id = u.id AND pa.element = e.element
)
ON CONFLICT DO NOTHING;

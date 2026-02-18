-- Resonance achievement tracking
CREATE TABLE IF NOT EXISTS resonance_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_type VARCHAR(100) NOT NULL,
    resonance_score_max INT DEFAULT 0,
    matches_count INT DEFAULT 0,
    unlocked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_achievement_type CHECK (achievement_type IN ('RESONANCE_MASTER', 'RESONANCE_COLLECTOR', 'HARMONY_PERFECTIONIST')),
    UNIQUE (user_id, achievement_type)
);

-- Indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_resonance_achievements_user_id ON resonance_achievements(user_id);
CREATE INDEX IF NOT EXISTS idx_resonance_achievements_unlocked ON resonance_achievements(unlocked_at) WHERE unlocked_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_resonance_achievements_type ON resonance_achievements(achievement_type);

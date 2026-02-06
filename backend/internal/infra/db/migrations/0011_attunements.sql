-- Player attunements table (6 elements per player)
CREATE TABLE IF NOT EXISTS player_attunements (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    element INT NOT NULL, -- 0=Fire, 1=Water, 2=Earth, 3=Wind, 4=Light, 5=Dark
    level INT NOT NULL DEFAULT 1,
    xp INT NOT NULL DEFAULT 0,
    total_xp_earned INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT player_attunements_element_check CHECK (element >= 0 AND element <= 5),
    CONSTRAINT player_attunements_level_check CHECK (level >= 1 AND level <= 25),
    CONSTRAINT unique_user_element UNIQUE (user_id, element)
);

-- XP history for tracking sources
CREATE TABLE IF NOT EXISTS attunement_xp_history (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    element INT NOT NULL,
    xp_amount INT NOT NULL,
    source VARCHAR(50) NOT NULL, -- 'match_win', 'match_loss', 'daily_login', 'quest', 'skill_use', etc.
    old_level INT NOT NULL,
    new_level INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_player_attunements_user_id ON player_attunements(user_id);
CREATE INDEX IF NOT EXISTS idx_attunement_xp_history_user_id ON attunement_xp_history(user_id);
CREATE INDEX IF NOT EXISTS idx_attunement_xp_history_created ON attunement_xp_history(created_at);

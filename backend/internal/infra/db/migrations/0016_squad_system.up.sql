-- Squad System: Managing 3-fighter squads
CREATE TABLE IF NOT EXISTS squads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Ensure only one active squad per user
CREATE UNIQUE INDEX idx_unique_active_squad_per_user ON squads(user_id) WHERE (is_active = true);

CREATE TABLE IF NOT EXISTS squad_members (
    squad_id UUID NOT NULL REFERENCES squads(id) ON DELETE CASCADE,
    fighter_id UUID NOT NULL REFERENCES fighters(id) ON DELETE CASCADE,
    slot_index INT NOT NULL CHECK (slot_index BETWEEN 0 AND 2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (squad_id, slot_index),
    -- Ensure a fighter is not in the same squad multiple times
    UNIQUE (squad_id, fighter_id)
);

CREATE INDEX idx_squads_user_id ON squads(user_id);

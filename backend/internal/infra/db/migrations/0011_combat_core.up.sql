-- Combat Logs and Detailed Battle History
CREATE TABLE IF NOT EXISTS combat_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL,
    round INT NOT NULL,
    tick INT NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_combat_logs_match_id ON combat_logs(match_id);

CREATE TABLE IF NOT EXISTS battle_details (
    match_id UUID PRIMARY KEY,
    winner_id UUID,
    total_rounds INT,
    summary JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

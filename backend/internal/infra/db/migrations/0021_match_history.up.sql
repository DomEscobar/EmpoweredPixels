-- Migration: Match History & Player Stats
-- Created: 2026-02-06

-- Match history for players
CREATE TABLE IF NOT EXISTS match_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    fighter_id UUID NOT NULL REFERENCES fighters(id) ON DELETE CASCADE,
    result TEXT NOT NULL CHECK (result IN ('win', 'loss', 'draw')),
    damage_dealt BIGINT NOT NULL DEFAULT 0,
    damage_taken BIGINT NOT NULL DEFAULT 0,
    kills INTEGER NOT NULL DEFAULT 0,
    deaths INTEGER NOT NULL DEFAULT 0,
    played_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index for quick lookups
CREATE INDEX IF NOT EXISTS idx_match_history_user_id ON match_history(user_id);
CREATE INDEX IF NOT EXISTS idx_match_history_played_at ON match_history(played_at DESC);
CREATE INDEX IF NOT EXISTS idx_match_history_match_id ON match_history(match_id);

-- Online players tracking (for "Players Online" counter)
CREATE TABLE IF NOT EXISTS player_sessions (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    last_activity TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_online BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_player_sessions_online ON player_sessions(is_online, last_activity);

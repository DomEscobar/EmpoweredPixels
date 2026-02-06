-- Migration: Create Leagues Table
-- Created: 2026-02-05

CREATE TABLE IF NOT EXISTS leagues (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    options JSONB NOT NULL DEFAULT '{}',
    is_deactivated BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS league_subscriptions (
    league_id INT REFERENCES leagues(id) ON DELETE CASCADE,
    fighter_id UUID REFERENCES fighters(id) ON DELETE CASCADE,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (league_id, fighter_id)
);

CREATE TABLE IF NOT EXISTS league_matches (
    league_id INT REFERENCES leagues(id) ON DELETE CASCADE,
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    started TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (league_id, match_id)
);

-- Seed Initial Leagues
INSERT INTO leagues (name, options) VALUES
    ('Bronze League', '{"min_level": 1, "max_level": 10, "entry_fee": 0}'),
    ('Silver League', '{"min_level": 11, "max_level": 20, "entry_fee": 100}'),
    ('Gold League', '{"min_level": 21, "max_level": null, "entry_fee": 500}')
ON CONFLICT (name) DO NOTHING;

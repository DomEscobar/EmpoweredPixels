-- Migration: Weekend Events System
-- Created: 2026-02-06

CREATE TABLE IF NOT EXISTS weekend_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- 'double_drops', 'double_xp', 'bonus_gold'
    multiplier DECIMAL(3,2) NOT NULL DEFAULT 2.00,
    start_day INTEGER NOT NULL CHECK (start_day BETWEEN 0 AND 6), -- 0=Sun, 6=Sat
    end_day INTEGER NOT NULL CHECK (end_day BETWEEN 0 AND 6),
    start_hour INTEGER NOT NULL DEFAULT 0,
    end_hour INTEGER NOT NULL DEFAULT 23,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS active_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES weekend_events(id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ends_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

INSERT INTO weekend_events (name, description, event_type, multiplier, start_day, end_day, start_hour, end_hour) VALUES
('Double Drop Weekend', 'All equipment drops have 2x chance!', 'double_drops', 2.00, 5, 6, 0, 23),
('XP Frenzy', 'Double experience from all matches!', 'double_xp', 2.00, 5, 6, 0, 23),
('Gold Rush', 'Earn 50% more gold from matches!', 'bonus_gold', 1.50, 0, 0, 0, 23)
ON CONFLICT (name) DO NOTHING;

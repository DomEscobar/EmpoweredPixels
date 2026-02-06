-- Migration: Weekend Events System
-- Created: 2026-02-06

-- Weekend events configuration
CREATE TABLE IF NOT EXISTS weekend_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- 'double_drops', 'double_xp', 'half_price', 'bonus_gold'
    multiplier DECIMAL(3,2) NOT NULL DEFAULT 2.00,
    start_day INTEGER NOT NULL CHECK (start_day BETWEEN 0 AND 6), -- 0=Sunday, 6=Saturday
    end_day INTEGER NOT NULL CHECK (end_day BETWEEN 0 AND 6),
    start_hour INTEGER NOT NULL DEFAULT 0 CHECK (start_hour BETWEEN 0 AND 23),
    end_hour INTEGER NOT NULL DEFAULT 23 CHECK (end_hour BETWEEN 0 AND 23),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Active event tracking
CREATE TABLE IF NOT EXISTS active_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES weekend_events(id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ends_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Event participation tracking
CREATE TABLE IF NOT EXISTS event_participation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES weekend_events(id) ON DELETE CASCADE,
    participation_count INTEGER NOT NULL DEFAULT 0,
    total_bonus_earned INTEGER NOT NULL DEFAULT 0,
    last_participation TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, event_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_weekend_events_active ON weekend_events(is_active);
CREATE INDEX IF NOT EXISTS idx_active_events_active ON active_events(is_active, ends_at);
CREATE INDEX IF NOT EXISTS idx_event_participation_user ON event_participation(user_id);

-- Seed default weekend events
INSERT INTO weekend_events (name, description, event_type, multiplier, start_day, end_day, start_hour, end_hour) VALUES
('Double Drop Weekend', 'All equipment drops have 2x chance!', 'double_drops', 2.00, 5, 6, 0, 23),
('XP Frenzy', 'Double experience from all matches!', 'double_xp', 2.00, 5, 6, 0, 23),
('Gold Rush', 'Earn 50% more gold from matches!', 'bonus_gold', 1.50, 0, 0, 0, 23)
ON CONFLICT DO NOTHING;

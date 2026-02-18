-- Add squad resonance tracking
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'squads' AND column_name = 'resonance_score') THEN
        ALTER TABLE squads ADD COLUMN resonance_score INT DEFAULT 0;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'squads' AND column_name = 'resonance_pattern') THEN
        ALTER TABLE squads ADD COLUMN resonance_pattern VARCHAR(50);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'squads' AND column_name = 'last_calculated_at') THEN
        ALTER TABLE squads ADD COLUMN last_calculated_at TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Create index for efficient resonance lookups
CREATE INDEX IF NOT EXISTS idx_squads_resonance_score ON squads(resonance_score);
CREATE INDEX IF NOT EXISTS idx_squads_resonance_pattern ON squads(resonance_pattern);

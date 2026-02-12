-- Add squad resonance tracking
ALTER TABLE squads ADD COLUMN resonance_score INT DEFAULT 0;
ALTER TABLE squads ADD COLUMN resonance_pattern VARCHAR(50);
ALTER TABLE squads ADD COLUMN last_calculated_at TIMESTAMP WITH TIME ZONE;

-- Create index for efficient resonance lookups
CREATE INDEX idx_squads_resonance_score ON squads(resonance_score);
CREATE INDEX idx_squads_resonance_pattern ON squads(resonance_pattern);

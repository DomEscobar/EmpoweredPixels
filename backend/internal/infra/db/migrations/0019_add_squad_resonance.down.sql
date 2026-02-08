-- Rollback squad resonance tracking
DROP INDEX IF EXISTS idx_squads_resonance_pattern;
DROP INDEX IF EXISTS idx_squads_resonance_score;

ALTER TABLE squads DROP COLUMN last_calculated_at;
ALTER TABLE squads DROP COLUMN resonance_pattern;
ALTER TABLE squads DROP COLUMN resonance_score;

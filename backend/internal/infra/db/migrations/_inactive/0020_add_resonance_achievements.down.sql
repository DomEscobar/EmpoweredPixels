-- Rollback resonance achievements
DROP INDEX IF EXISTS idx_resonance_achievements_type;
DROP INDEX IF EXISTS idx_resonance_achievements_unlocked;
DROP INDEX IF EXISTS idx_resonance_achievements_user_id;
DROP TABLE IF EXISTS resonance_achievements;

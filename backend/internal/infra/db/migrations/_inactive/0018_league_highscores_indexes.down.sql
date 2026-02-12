-- Rollback indexes for league highscores

DROP INDEX IF EXISTS idx_match_score_fighters_fighter_id;
DROP INDEX IF EXISTS idx_league_matches_league_id_started;

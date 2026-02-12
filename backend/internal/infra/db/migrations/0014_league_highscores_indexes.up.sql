-- Add indexes for league highscores performance

-- Index on match_score_fighters.fighter_id to speed up grouping by fighter
CREATE INDEX IF NOT EXISTS idx_match_score_fighters_fighter_id ON match_score_fighters(fighter_id);

-- Index on league_matches(league_id, started) to efficiently fetch recent matches for a league
CREATE INDEX IF NOT EXISTS idx_league_matches_league_id_started ON league_matches(league_id, started DESC);

package repositories

import (
	"context"
	"errors"

	"empoweredpixels/internal/domain/roster"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FighterRepository struct {
	pool *pgxpool.Pool
}

type ExperienceRepository struct {
	pool *pgxpool.Pool
}

type ConfigurationRepository struct {
	pool *pgxpool.Pool
}

func NewFighterRepository(pool *pgxpool.Pool) *FighterRepository {
	return &FighterRepository{pool: pool}
}

func NewExperienceRepository(pool *pgxpool.Pool) *ExperienceRepository {
	return &ExperienceRepository{pool: pool}
}

func NewConfigurationRepository(pool *pgxpool.Pool) *ConfigurationRepository {
	return &ConfigurationRepository{pool: pool}
}

func (r *FighterRepository) ListByUser(ctx context.Context, userID int64) ([]roster.Fighter, error) {
	const query = `
		select id, user_id, name, level, xp, xp_to_next_level, power, condition_power, precision, ferocity, accuracy, agility, armor, vitality, parry_chance, healing_power, speed, vision, weapon_id, attunement_id, matches_won, matches_lost, total_matches, total_damage_dealt, total_damage_taken, created, is_deleted
		from fighters
		where user_id = $1 and is_deleted = false
		order by created`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fighters []roster.Fighter
	for rows.Next() {
		var fighter roster.Fighter
		if err := rows.Scan(
			&fighter.ID, &fighter.UserID, &fighter.Name, &fighter.Level, &fighter.XP, &fighter.XPToNextLevel,
			&fighter.Power, &fighter.ConditionPower, &fighter.Precision, &fighter.Ferocity,
			&fighter.Accuracy, &fighter.Agility, &fighter.Armor, &fighter.Vitality,
			&fighter.ParryChance, &fighter.HealingPower, &fighter.Speed, &fighter.Vision,
			&fighter.WeaponID, &fighter.AttunementID,
			&fighter.MatchesWon, &fighter.MatchesLost, &fighter.TotalMatches,
			&fighter.TotalDamageDealt, &fighter.TotalDamageTaken,
			&fighter.Created, &fighter.IsDeleted,
		); err != nil {
			return nil, err
		}
		fighters = append(fighters, fighter)
	}

	return fighters, rows.Err()
}

func (r *FighterRepository) GetByUserAndID(ctx context.Context, userID int64, id string) (*roster.Fighter, error) {
	const query = `
		select id, user_id, name, level, xp, xp_to_next_level, power, condition_power, precision, ferocity, accuracy, agility, armor, vitality, parry_chance, healing_power, speed, vision, weapon_id, attunement_id, matches_won, matches_lost, total_matches, total_damage_dealt, total_damage_taken, created, is_deleted
		from fighters
		where user_id = $1 and id = $2 and is_deleted = false`

	var fighter roster.Fighter
	err := r.pool.QueryRow(ctx, query, userID, id).Scan(
		&fighter.ID, &fighter.UserID, &fighter.Name, &fighter.Level, &fighter.XP, &fighter.XPToNextLevel,
		&fighter.Power, &fighter.ConditionPower, &fighter.Precision, &fighter.Ferocity,
		&fighter.Accuracy, &fighter.Agility, &fighter.Armor, &fighter.Vitality,
		&fighter.ParryChance, &fighter.HealingPower, &fighter.Speed, &fighter.Vision,
		&fighter.WeaponID, &fighter.AttunementID,
		&fighter.MatchesWon, &fighter.MatchesLost, &fighter.TotalMatches,
		&fighter.TotalDamageDealt, &fighter.TotalDamageTaken,
		&fighter.Created, &fighter.IsDeleted,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &fighter, nil
}

func (r *FighterRepository) ListByMatch(ctx context.Context, matchID string) ([]roster.Fighter, error) {
	const query = `
		select f.id, f.user_id, f.name, f.level, f.xp, f.xp_to_next_level, f.power, f.condition_power, f.precision, f.ferocity, f.accuracy, f.agility, f.armor, f.vitality, f.parry_chance, f.healing_power, f.speed, f.vision, f.weapon_id, f.attunement_id, f.matches_won, f.matches_lost, f.total_matches, f.total_damage_dealt, f.total_damage_taken, f.created, f.is_deleted
		from fighters f
		join match_registrations mr on mr.fighter_id = f.id
		where mr.match_id = $1 and f.is_deleted = false`

	rows, err := r.pool.Query(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fighters []roster.Fighter
	for rows.Next() {
		var fighter roster.Fighter
		if err := rows.Scan(
			&fighter.ID, &fighter.UserID, &fighter.Name, &fighter.Level, &fighter.XP, &fighter.XPToNextLevel,
			&fighter.Power, &fighter.ConditionPower, &fighter.Precision, &fighter.Ferocity,
			&fighter.Accuracy, &fighter.Agility, &fighter.Armor, &fighter.Vitality,
			&fighter.ParryChance, &fighter.HealingPower, &fighter.Speed, &fighter.Vision,
			&fighter.WeaponID, &fighter.AttunementID,
			&fighter.MatchesWon, &fighter.MatchesLost, &fighter.TotalMatches,
			&fighter.TotalDamageDealt, &fighter.TotalDamageTaken,
			&fighter.Created, &fighter.IsDeleted,
		); err != nil {
			return nil, err
		}
		fighters = append(fighters, fighter)
	}

	return fighters, rows.Err()
}

func (r *FighterRepository) GetByID(ctx context.Context, id string) (*roster.Fighter, error) {
	const query = `
		select id, user_id, name, level, xp, xp_to_next_level, power, condition_power, precision, ferocity, accuracy, agility, armor, vitality, parry_chance, healing_power, speed, vision, weapon_id, attunement_id, matches_won, matches_lost, total_matches, total_damage_dealt, total_damage_taken, created, is_deleted
		from fighters
		where id = $1`

	var fighter roster.Fighter
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&fighter.ID, &fighter.UserID, &fighter.Name, &fighter.Level, &fighter.XP, &fighter.XPToNextLevel,
		&fighter.Power, &fighter.ConditionPower, &fighter.Precision, &fighter.Ferocity,
		&fighter.Accuracy, &fighter.Agility, &fighter.Armor, &fighter.Vitality,
		&fighter.ParryChance, &fighter.HealingPower, &fighter.Speed, &fighter.Vision,
		&fighter.WeaponID, &fighter.AttunementID,
		&fighter.MatchesWon, &fighter.MatchesLost, &fighter.TotalMatches,
		&fighter.TotalDamageDealt, &fighter.TotalDamageTaken,
		&fighter.Created, &fighter.IsDeleted,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &fighter, nil
}

func (r *FighterRepository) NameExists(ctx context.Context, name string) (bool, error) {
	const query = `select 1 from fighters where name = $1`
	var exists int
	err := r.pool.QueryRow(ctx, query, name).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *FighterRepository) UserHasFighter(ctx context.Context, userID int64) (bool, error) {
	const query = `select 1 from fighters where user_id = $1 and is_deleted = false limit 1`
	var exists int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *FighterRepository) Create(ctx context.Context, fighter *roster.Fighter) error {
	const query = `
		insert into fighters (id, user_id, name, level, xp, xp_to_next_level, power, condition_power, precision, ferocity, accuracy, agility, armor, vitality, parry_chance, healing_power, speed, vision, weapon_id, attunement_id, matches_won, matches_lost, total_matches, total_damage_dealt, total_damage_taken, created, is_deleted)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`

	_, err := r.pool.Exec(ctx, query,
		fighter.ID,
		fighter.UserID,
		fighter.Name,
		fighter.Level,
		fighter.XP,
		fighter.XPToNextLevel,
		fighter.Power,
		fighter.ConditionPower,
		fighter.Precision,
		fighter.Ferocity,
		fighter.Accuracy,
		fighter.Agility,
		fighter.Armor,
		fighter.Vitality,
		fighter.ParryChance,
		fighter.HealingPower,
		fighter.Speed,
		fighter.Vision,
		fighter.WeaponID,
		fighter.AttunementID,
		fighter.MatchesWon,
		fighter.MatchesLost,
		fighter.TotalMatches,
		fighter.TotalDamageDealt,
		fighter.TotalDamageTaken,
		fighter.Created,
		fighter.IsDeleted,
	)
	return err
}

func (r *FighterRepository) SoftDelete(ctx context.Context, userID int64, id string) error {
	const query = `
		update fighters
		set is_deleted = true
		where id = $1 and user_id = $2`

	_, err := r.pool.Exec(ctx, query, id, userID)
	return err
}

func (r *FighterRepository) Update(ctx context.Context, fighter *roster.Fighter) error {
	const query = `
		update fighters
		set level = $1, power = $2, condition_power = $3, precision = $4, ferocity = $5,
		    accuracy = $6, agility = $7, armor = $8, vitality = $9, parry_chance = $10,
		    healing_power = $11, speed = $12, vision = $13
		where id = $14`

	_, err := r.pool.Exec(ctx, query,
		fighter.Level, fighter.Power, fighter.ConditionPower, fighter.Precision, fighter.Ferocity,
		fighter.Accuracy, fighter.Agility, fighter.Armor, fighter.Vitality, fighter.ParryChance,
		fighter.HealingPower, fighter.Speed, fighter.Vision,
		fighter.ID,
	)
	return err
}

func (r *ExperienceRepository) GetByFighterID(ctx context.Context, fighterID string) (*roster.FighterExperience, error) {
	const query = `
		select id, fighter_id, experience
		from fighter_experiences
		where fighter_id = $1`

	var exp roster.FighterExperience
	err := r.pool.QueryRow(ctx, query, fighterID).Scan(&exp.ID, &exp.FighterID, &exp.Experience)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

func (r *ExperienceRepository) Upsert(ctx context.Context, experience *roster.FighterExperience) error {
	const query = `
		insert into fighter_experiences (fighter_id, experience)
		values ($1, $2)
		on conflict (fighter_id)
		do update set experience = excluded.experience`

	_, err := r.pool.Exec(ctx, query, experience.FighterID, experience.Experience)
	return err
}

func (r *ConfigurationRepository) GetByFighterID(ctx context.Context, fighterID string) (*roster.FighterConfiguration, error) {
	const query = `
		select fighter_id, attunement_id
		from fighter_configurations
		where fighter_id = $1`

	var config roster.FighterConfiguration
	err := r.pool.QueryRow(ctx, query, fighterID).Scan(&config.FighterID, &config.AttunementID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ConfigurationRepository) Upsert(ctx context.Context, configuration *roster.FighterConfiguration) error {
	const query = `
		insert into fighter_configurations (fighter_id, attunement_id)
		values ($1, $2)
		on conflict (fighter_id)
		do update set attunement_id = excluded.attunement_id`

	_, err := r.pool.Exec(ctx, query, configuration.FighterID, configuration.AttunementID)
	return err
}

package guilds

import (
	"context"
	"database/sql"
	"empoweredpixels/internal/domain/guilds"
	"github.com/google/uuid"
	"time"
)

type GuildRepository struct {
	pool *sql.DB
}

func NewGuildRepository(pool *sql.DB) *GuildRepository {
	return &GuildRepository{pool: pool}
}

func (r *GuildRepository) Create(ctx context.Context, g *guilds.Guild) error {
	if g.ID == "" {
		g.ID = uuid.NewString()
	}
	query := `
		INSERT INTO guilds (id, name, description, leader_id, level, experience, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, g.ID, g.Name, g.Description, g.LeaderID, g.Level, g.Experience, now, now)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*guilds.Guild, error) {
	query := `SELECT id, name, description, leader_id, level, experience, created_at, updated_at FROM guilds WHERE id = $1`
	g := &guilds.Guild{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&g.ID, &g.Name, &g.Description, &g.LeaderID, &g.Level, &g.Experience, &g.CreatedAt, &g.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return g, err
}

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (*guilds.Guild, error) {
	query := `SELECT id, name, description, leader_id, level, experience, created_at, updated_at FROM guilds WHERE name = $1`
	g := &guilds.Guild{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(&g.ID, &g.Name, &g.Description, &g.LeaderID, &g.Level, &g.Experience, &g.CreatedAt, &g.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return g, err
}

func (r *PostgresRepository) List(ctx context.Context) ([]guilds.Guild, error) {
	query := `SELECT id, name, description, leader_id, level, experience, created_at, updated_at FROM guilds`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []guilds.Guild
	for rows.Next() {
		var g guilds.Guild
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.LeaderID, &g.Level, &g.Experience, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, g)
	}
	return result, nil
}

func (r *PostgresRepository) AddMember(ctx context.Context, m *guilds.GuildMember) error {
	query := `INSERT INTO guild_members (guild_id, fighter_id, role, joined_at) VALUES ($1, $2, $3, $4)`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, m.GuildID, m.FighterID, m.Role, now)
	return err
}

func (r *PostgresRepository) RemoveMember(ctx context.Context, guildID, fighterID string) error {
	query := `DELETE FROM guild_members WHERE guild_id = $1 AND fighter_id = $2`
	_, err := r.db.ExecContext(ctx, query, guildID, fighterID)
	return err
}

func (r *PostgresRepository) GetMembers(ctx context.Context, guildID string) ([]guilds.GuildMember, error) {
	query := `SELECT guild_id, fighter_id, role, joined_at FROM guild_members WHERE guild_id = $1`
	rows, err := r.db.QueryContext(ctx, query, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []guilds.GuildMember
	for rows.Next() {
		var m guilds.GuildMember
		if err := rows.Scan(&m.GuildID, &m.FighterID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (r *PostgresRepository) GetFighterGuild(ctx context.Context, fighterID string) (*guilds.Guild, error) {
	query := `
		SELECT g.id, g.name, g.description, g.leader_id, g.level, g.experience, g.created_at, g.updated_at 
		FROM guilds g
		JOIN guild_members gm ON g.id = gm.guild_id
		WHERE gm.fighter_id = $1
	`
	g := &guilds.Guild{}
	err := r.db.QueryRowContext(ctx, query, fighterID).Scan(&g.ID, &g.Name, &g.Description, &g.LeaderID, &g.Level, &g.Experience, &g.CreatedAt, &g.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return g, err
}

func (r *PostgresRepository) CreateRequest(ctx context.Context, req *guilds.GuildRequest) error {
	if req.ID == "" {
		req.ID = uuid.NewString()
	}
	query := `INSERT INTO guild_requests (id, guild_id, fighter_id, status, created_at) VALUES ($1, $2, $3, $4, $5)`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, req.ID, req.GuildID, req.FighterID, req.Status, now)
	return err
}

func (r *PostgresRepository) ListRequests(ctx context.Context, guildID string) ([]guilds.GuildRequest, error) {
	query := `SELECT id, guild_id, fighter_id, status, created_at FROM guild_requests WHERE guild_id = $1`
	rows, err := r.db.QueryContext(ctx, query, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []guilds.GuildRequest
	for rows.Next() {
		var req guilds.GuildRequest
		if err := rows.Scan(&req.ID, &req.GuildID, &req.FighterID, &req.Status, &req.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, req)
	}
	return result, nil
}

func (r *PostgresRepository) UpdateRequest(ctx context.Context, requestID string, status string) error {
	query := `UPDATE guild_requests SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, requestID)
	return err
}

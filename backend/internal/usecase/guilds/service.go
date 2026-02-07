package guilds

import (
	"context"
	"errors"
	"empoweredpixels/internal/domain/guilds"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateGuild(ctx context.Context, fighterID, name, description string) (*guilds.Guild, error) {
	// Check if already in a guild
	existing, _ := s.repo.GetFighterGuild(ctx, fighterID)
	if existing != nil {
		return nil, errors.New("already in a guild")
	}

	guild := &guilds.Guild{
		Name:        name,
		Description: description,
		LeaderID:    fighterID,
		Level:       1,
	}

	if err := s.repo.Create(ctx, guild); err != nil {
		return nil, err
	}

	// Add leader as member
	member := &guilds.GuildMember{
		GuildID:  guild.ID,
		FighterID: fighterID,
		Role:     "leader",
	}
	s.repo.AddMember(ctx, member)

	return guild, nil
}

func (s *Service) JoinGuild(ctx context.Context, fighterID, guildID string) error {
	req := &guilds.GuildRequest{
		GuildID:   guildID,
		FighterID: fighterID,
		Status:    "pending",
	}
	return s.repo.CreateRequest(ctx, req)
}

func (s *Service) ListGuilds(ctx context.Context) ([]guilds.Guild, error) {
	return s.repo.List(ctx)
}

package guilds

import (
	"context"
	"empoweredpixels/internal/domain/guilds"
)

type Repository interface {
	Create(ctx context.Context, guild *guilds.Guild) error
	GetByID(ctx context.Context, id string) (*guilds.Guild, error)
	GetByName(ctx context.Context, name string) (*guilds.Guild, error)
	List(ctx context.Context) ([]guilds.Guild, error)
	
	AddMember(ctx context.Context, member *guilds.GuildMember) error
	RemoveMember(ctx context.Context, guildID, fighterID string) error
	GetMembers(ctx context.Context, guildID string) ([]guilds.GuildMember, error)
	GetFighterGuild(ctx context.Context, fighterID string) (*guilds.Guild, error)
	
	CreateRequest(ctx context.Context, req *guilds.GuildRequest) error
	ListRequests(ctx context.Context, guildID string) ([]guilds.GuildRequest, error)
	UpdateRequest(ctx context.Context, requestID string, status string) error
}

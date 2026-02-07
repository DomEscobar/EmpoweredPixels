package guilds

import "time"

type Guild struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LeaderID    string    `json:"leaderId"`
	Level       int       `json:"level"`
	Experience  int       `json:"experience"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GuildMember struct {
	GuildID  string    `json:"guildId"`
	FighterID string   `json:"fighterId"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joinedAt"`
}

type GuildRequest struct {
	ID        string    `json:"id"`
	GuildID   string    `json:"guildId"`
	FighterID string    `json:"fighterId"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

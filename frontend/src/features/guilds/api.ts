export interface Guild {
  id: string;
  name: string;
  description: string;
  leaderId: string;
  level: number;
  experience: number;
  createdAt: string;
  updatedAt: string;
}

export interface GuildMember {
  guildId: string;
  fighterId: string;
  role: string;
  joinedAt: string;
}

export interface GuildRequest {
  id: string;
  guildId: string;
  fighterId: string;
  status: string;
  createdAt: string;
}

export async function getGuilds(token: string): Promise<Guild[]> {
  const response = await fetch('/api/guilds', {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (!response.ok) throw new Error('Failed to fetch guilds');
  return response.json();
}

export async function createGuild(token: string, name: string, description: string): Promise<Guild> {
  const response = await fetch('/api/guilds', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ name, description })
  });
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to create guild');
  }
  return response.json();
}

export async function joinGuild(token: string, guildId: string): Promise<{ message: string }> {
  const response = await fetch(`/api/guilds/${guildId}/join`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to join guild');
  }
  return response.json();
}

import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

export interface Squad {
  id: string;
  userId: number;
  name: string;
  isActive: boolean;
  members: SquadMember[];
  createdAt: string;
  updatedAt: string;
}

export interface SquadMember {
  fighterId: string;
  slotIndex: number;
}

export interface SetSquadRequest {
  name: string;
  fighterIds: string[];
}

export async function getActiveSquad(token: string) {
  return request<Squad>(endpoints.squads.active, { token });
}

export async function setActiveSquad(
  token: string,
  name: string,
  fighterIds: string[]
) {
  return request<Squad>(endpoints.squads.setActive, {
    method: "POST",
    token,
    body: { name, fighterIds },
  });
}

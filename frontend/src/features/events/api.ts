// Events API
const API_URL = import.meta.env.VITE_API_URL || "";

export interface WeekendEvent {
  id: string;
  name: string;
  description: string;
  event_type: string;
  multiplier: number;
  start_day: number;
  end_day: number;
  start_hour: number;
  end_hour: number;
  is_active: boolean;
}

export interface ActiveEvent {
  id: string;
  event_id: string;
  event?: WeekendEvent;
  started_at: string;
  ends_at: string;
  is_active: boolean;
}

export interface EventStatus {
  has_active_event: boolean;
  active_event?: ActiveEvent;
  time_remaining?: string;
  multiplier: number;
  type?: string;
}

export const EVENT_TYPE_LABELS: Record<string, { name: string; icon: string; color: string }> = {
  double_drops: { name: "Double Drops", icon: "📦", color: "text-purple-400" },
  double_xp: { name: "Double XP", icon: "⚡", color: "text-yellow-400" },
  bonus_gold: { name: "Bonus Gold", icon: "💰", color: "text-amber-400" },
  half_price: { name: "Half Price", icon: "🏷️", color: "text-green-400" },
};

export async function getEventStatus(token: string): Promise<EventStatus> {
  const response = await fetch(`${API_URL}/api/events/status`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!response.ok) throw new Error("Failed to fetch event status");
  return response.json();
}

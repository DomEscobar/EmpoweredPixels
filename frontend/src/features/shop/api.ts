import { endpoints } from "@/shared/api/endpoints";
import { request } from "@/shared/api/http";

// Rarity constants matching backend
export const Rarity = {
  Basic: 0,
  Common: 1,
  Rare: 2,
  Fabled: 3,
  Mythic: 4,
  Legendary: 5,
} as const;

export type RarityLevel = typeof Rarity[keyof typeof Rarity];

export interface ShopItem {
  id: string;
  shop_id: string;
  name: string;
  description: string;
  item_type: 'gold_package' | 'bundle' | 'equipment' | 'consumable';
  price_amount: number;
  price_currency: 'usd' | 'gold' | 'particles';
  gold_amount?: number;
  rarity: RarityLevel;
  image_url?: string;
  metadata: Record<string, unknown>;
  is_active: boolean;
  sort_order: number;
  created: string;
}

export interface PlayerGold {
  balance: number;
  lifetime_earned: number;
  lifetime_spent: number;
}

export interface Transaction {
  id: string;
  user_id: number;
  shop_item_id?: string;
  item_type: string;
  item_name: string;
  price_amount: number;
  price_currency: string;
  gold_change: number;
  status: string;
  metadata: Record<string, unknown>;
  created: string;
}

export interface PurchaseResponse {
  success: boolean;
  transaction_id: string;
  new_balance: number;
  items_received?: string[];
  message: string;
}

// Get all shop items
export async function getShopItems(token: string): Promise<ShopItem[]> {
  return request<ShopItem[]>(`${endpoints.shop}/items`, { token });
}

// Get gold packages only
export async function getGoldPackages(token: string): Promise<ShopItem[]> {
  return request<ShopItem[]>(`${endpoints.shop}/gold`, { token });
}

// Get bundles only
export async function getBundles(token: string): Promise<ShopItem[]> {
  return request<ShopItem[]>(`${endpoints.shop}/bundles`, { token });
}

// Get a specific shop item
export async function getShopItem(token: string, itemId: string): Promise<ShopItem> {
  return request<ShopItem>(`${endpoints.shop}/item/${itemId}`, { token });
}

// Purchase an item
export async function purchaseItem(token: string, itemId: string): Promise<PurchaseResponse> {
  return request<PurchaseResponse>(`${endpoints.shop}/purchase`, {
    method: "POST",
    token,
    body: { item_id: itemId },
  });
}

// Get player's gold balance
export async function getPlayerGold(token: string): Promise<PlayerGold> {
  return request<PlayerGold>(endpoints.playerGold, { token });
}

// Get player's transaction history
export async function getTransactions(token: string): Promise<Transaction[]> {
  return request<Transaction[]>(endpoints.playerTransactions, { token });
}

// Helper: Get rarity display name
export function getRarityName(rarity: RarityLevel): string {
  switch (rarity) {
    case Rarity.Basic: return "Basic";
    case Rarity.Common: return "Common";
    case Rarity.Rare: return "Rare";
    case Rarity.Fabled: return "Fabled";
    case Rarity.Mythic: return "Mythic";
    case Rarity.Legendary: return "Legendary";
    default: return "Unknown";
  }
}

// Helper: Get rarity color classes (Tailwind)
export function getRarityColor(rarity: RarityLevel): string {
  switch (rarity) {
    case Rarity.Basic: return "text-slate-400 border-slate-600 bg-slate-800/50";
    case Rarity.Common: return "text-green-400 border-green-600 bg-green-900/30";
    case Rarity.Rare: return "text-blue-400 border-blue-600 bg-blue-900/30";
    case Rarity.Fabled: return "text-purple-400 border-purple-600 bg-purple-900/30";
    case Rarity.Mythic: return "text-orange-400 border-orange-600 bg-orange-900/30";
    case Rarity.Legendary: return "text-yellow-400 border-yellow-600 bg-yellow-900/30";
    default: return "text-slate-400 border-slate-600 bg-slate-800/50";
  }
}

// Helper: Get rarity glow effect
export function getRarityGlow(rarity: RarityLevel): string {
  switch (rarity) {
    case Rarity.Basic: return "";
    case Rarity.Common: return "shadow-green-500/20";
    case Rarity.Rare: return "shadow-blue-500/30";
    case Rarity.Fabled: return "shadow-purple-500/40";
    case Rarity.Mythic: return "shadow-orange-500/50";
    case Rarity.Legendary: return "shadow-yellow-500/60";
    default: return "";
  }
}

// Helper: Format price for display
export function formatPrice(amount: number, currency: string): string {
  if (currency === 'usd') {
    return `$${(amount / 100).toFixed(2)}`;
  }
  if (currency === 'gold') {
    return `${amount.toLocaleString()} Gold`;
  }
  if (currency === 'particles') {
    return `${amount.toLocaleString()} Particles`;
  }
  return `${amount}`;
}

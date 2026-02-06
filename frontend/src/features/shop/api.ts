import { request } from "@/shared/api/http";
import { useAuthStore } from "@/features/auth/store";
import type { ShopItem, PlayerGold, Transaction, PurchaseResponse } from './types'

const SHOP_BASE = "/api/shop";
const PLAYER_BASE = "/api/player";

function getOptions(options: any = {}) {
  const auth = useAuthStore();
  return {
    ...options,
    token: auth.token || undefined
  };
}

export const shopApi = {
  // Get all shop items
  async getItems(): Promise<ShopItem[]> {
    return request<ShopItem[]>(`${SHOP_BASE}/items`, getOptions());
  },

  // Get gold packages only
  async getGoldPackages(): Promise<ShopItem[]> {
    return request<ShopItem[]>(`${SHOP_BASE}/gold`, getOptions());
  },

  // Get bundles only
  async getBundles(): Promise<ShopItem[]> {
    return request<ShopItem[]>(`${SHOP_BASE}/bundles`, getOptions());
  },

  // Get single item details
  async getItem(id: number): Promise<ShopItem> {
    return request<ShopItem>(`${SHOP_BASE}/item/${id}`, getOptions());
  },

  // Purchase an item
  async purchase(itemId: number): Promise<PurchaseResponse> {
    return request<PurchaseResponse>(`${SHOP_BASE}/purchase`, getOptions({
      method: "POST",
      body: { item_id: itemId }
    }));
  },

  // Get player gold balance
  async getGoldBalance(): Promise<PlayerGold> {
    return request<PlayerGold>(`${PLAYER_BASE}/gold`, getOptions());
  },

  // Get transaction history
  async getTransactions(limit = 50): Promise<Transaction[]> {
    return request<Transaction[]>(`${PLAYER_BASE}/transactions?limit=${limit}`, getOptions());
  }
}

// Rarity color mapping (matching 9-tier system)
export const rarityColors: Record<number, string> = {
  0: '#9ca3af', // Basic - Gray
  1: '#22c55e', // Common - Green
  2: '#3b82f6', // Rare - Blue
  3: '#a855f7', // Fabled - Purple
  4: '#f59e0b', // Mythic - Orange
  5: '#ef4444'  // Legendary - Red
}

export const rarityNames: Record<number, string> = {
  0: 'Basic',
  1: 'Common',
  2: 'Rare',
  3: 'Fabled',
  4: 'Mythic',
  5: 'Legendary'
}

// Format price for display
export function formatPrice(amount: number, currency: string): string {
  if (currency === 'usd') {
    return `$${(amount / 100).toFixed(2)}`
  }
  if (currency === 'gold') {
    return `${amount.toLocaleString()} Gold`
  }
  return `${amount} ${currency}`
}

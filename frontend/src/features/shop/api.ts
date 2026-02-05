import { api } from '@/shared/api'
import type { ShopItem, PlayerGold, Transaction, PurchaseResponse } from './types'

export const shopApi = {
  // Get all shop items
  async getItems(): Promise<ShopItem[]> {
    const response = await api.get('/shop/items')
    return response.data
  },

  // Get gold packages only
  async getGoldPackages(): Promise<ShopItem[]> {
    const response = await api.get('/shop/gold')
    return response.data
  },

  // Get bundles only
  async getBundles(): Promise<ShopItem[]> {
    const response = await api.get('/shop/bundles')
    return response.data
  },

  // Get single item details
  async getItem(id: number): Promise<ShopItem> {
    const response = await api.get(`/shop/item/${id}`)
    return response.data
  },

  // Purchase an item
  async purchase(itemId: number): Promise<PurchaseResponse> {
    const response = await api.post('/shop/purchase', { item_id: itemId })
    return response.data
  },

  // Get player gold balance
  async getGoldBalance(): Promise<PlayerGold> {
    const response = await api.get('/player/gold')
    return response.data
  },

  // Get transaction history
  async getTransactions(limit = 50): Promise<Transaction[]> {
    const response = await api.get('/player/transactions', { params: { limit } })
    return response.data
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

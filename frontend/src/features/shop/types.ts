export interface ShopItem {
  id: number
  shop_id: number
  name: string
  description: string
  item_type: 'gold_package' | 'bundle' | 'equipment' | 'consumable'
  price_amount: number
  price_currency: 'usd' | 'gold' | 'particles'
  gold_amount?: number
  rarity: number
  image_url?: string
  metadata?: Record<string, any>
  is_active: boolean
  sort_order: number
  created: string
}

export interface PlayerGold {
  user_id: number
  balance: number
  lifetime_earned: number
  lifetime_spent: number
  updated: string
}

export interface Transaction {
  id: number
  user_id: number
  shop_item_id?: number
  item_type: string
  item_name: string
  price_amount: number
  price_currency: string
  gold_change: number
  status: 'pending' | 'completed' | 'failed' | 'refunded'
  metadata?: Record<string, any>
  created: string
}

export interface PurchaseResponse {
  success: boolean
  transaction_id: number
  new_balance: number
  items_received?: string[]
  message: string
}

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { shopApi, rarityColors, rarityNames, formatPrice } from './api'
import type { ShopItem, PlayerGold, Transaction, PurchaseResponse } from './types'

export const useShopStore = defineStore('shop', () => {
  // State
  const items = ref<ShopItem[]>([])
  const goldPackages = ref<ShopItem[]>([])
  const bundles = ref<ShopItem[]>([])
  const goldBalance = ref<PlayerGold | null>(null)
  const transactions = ref<Transaction[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const purchaseInProgress = ref(false)

  // Getters
  const goldItems = computed(() => 
    items.value.filter(item => item.item_type === 'gold_package')
  )
  
  const bundleItems = computed(() => 
    items.value.filter(item => item.item_type === 'bundle')
  )

  const hasInsufficientGold = computed(() => (item: ShopItem) => {
    if (item.price_currency !== 'gold') return false
    return (goldBalance.value?.balance || 0) < item.price_amount
  })

  // Actions
  async function fetchItems() {
    loading.value = true
    error.value = null
    try {
      items.value = await shopApi.getItems()
    } catch (err: any) {
      error.value = err.message || 'Failed to load shop items'
    } finally {
      loading.value = false
    }
  }

  async function fetchGoldPackages() {
    loading.value = true
    error.value = null
    try {
      goldPackages.value = await shopApi.getGoldPackages()
    } catch (err: any) {
      error.value = err.message || 'Failed to load gold packages'
    } finally {
      loading.value = false
    }
  }

  async function fetchBundles() {
    loading.value = true
    error.value = null
    try {
      bundles.value = await shopApi.getBundles()
    } catch (err: any) {
      error.value = err.message || 'Failed to load bundles'
    } finally {
      loading.value = false
    }
  }

  async function fetchGoldBalance() {
    try {
      goldBalance.value = await shopApi.getGoldBalance()
    } catch (err: any) {
      console.error('Failed to fetch gold balance:', err)
    }
  }

  async function fetchTransactions(limit = 50) {
    try {
      transactions.value = await shopApi.getTransactions(limit)
    } catch (err: any) {
      console.error('Failed to fetch transactions:', err)
    }
  }

  async function purchaseItem(itemId: number): Promise<PurchaseResponse> {
    purchaseInProgress.value = true
    error.value = null
    try {
      const result = await shopApi.purchase(itemId)
      if (result.success) {
        // Refresh balance after successful purchase
        await fetchGoldBalance()
        await fetchTransactions()
      }
      return result
    } catch (err: any) {
      error.value = err.message || 'Purchase failed'
      return {
        success: false,
        transaction_id: 0,
        new_balance: goldBalance.value?.balance || 0,
        message: error.value
      }
    } finally {
      purchaseInProgress.value = false
    }
  }

  function getRarityColor(rarity: number): string {
    return rarityColors[rarity] || '#9ca3af'
  }

  function getRarityName(rarity: number): string {
    return rarityNames[rarity] || 'Unknown'
  }

  function formatItemPrice(item: ShopItem): string {
    return formatPrice(item.price_amount, item.price_currency)
  }

  return {
    // State
    items,
    goldPackages,
    bundles,
    goldBalance,
    transactions,
    loading,
    error,
    purchaseInProgress,
    // Getters
    goldItems,
    bundleItems,
    hasInsufficientGold,
    // Actions
    fetchItems,
    fetchGoldPackages,
    fetchBundles,
    fetchGoldBalance,
    fetchTransactions,
    purchaseItem,
    getRarityColor,
    getRarityName,
    formatItemPrice
  }
})

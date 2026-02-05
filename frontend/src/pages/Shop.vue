<template>
  <div class="shop-page">
    <header class="shop-header">
      <h1>🏪 Shop</h1>
      <GoldDisplay />
    </header>

    <nav class="shop-tabs">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        :class="['tab', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        {{ tab.label }}
      </button>
    </nav>

    <main class="shop-content">
      <!-- Gold Packages Tab -->
      <section v-if="activeTab === 'gold'" class="tab-content">
        <h2>Gold Emporium</h2>
        <p class="tab-description">Purchase gold for upgrades and items</p>
        
        <div v-if="shopStore.loading" class="loading">
          Loading...
        </div>
        
        <div v-else-if="shopStore.error" class="error">
          {{ shopStore.error }}
        </div>
        
        <div v-else class="gold-grid">
          <GoldPackageCard 
            v-for="item in shopStore.goldPackages" 
            :key="item.id"
            :item="item"
            @purchase="openPurchaseModal"
          />
        </div>
      </section>

      <!-- Bundles Tab -->
      <section v-if="activeTab === 'bundles'" class="tab-content">
        <h2>Equipment Bundles</h2>
        <p class="tab-description">Curated bundles with guaranteed rarity items</p>
        
        <div v-if="shopStore.loading" class="loading">
          Loading...
        </div>
        
        <div v-else-if="shopStore.error" class="error">
          {{ shopStore.error }}
        </div>
        
        <div v-else class="bundles-grid">
          <BundleCard 
            v-for="item in shopStore.bundles" 
            :key="item.id"
            :item="item"
            @purchase="openPurchaseModal"
          />
        </div>
      </section>

      <!-- History Tab -->
      <section v-if="activeTab === 'history'" class="tab-content">
        <h2>Purchase History</h2>
        <p class="tab-description">Your recent transactions</p>
        
        <div v-if="shopStore.transactions.length === 0" class="empty">
          No transactions yet
        </div>
        
        <div v-else class="transactions-list">
          <div 
            v-for="tx in shopStore.transactions" 
            :key="tx.id"
            :class="['transaction', tx.status]"
          >
            <div class="tx-info">
              <span class="tx-name">{{ tx.item_name }}</span>
              <span class="tx-date">{{ formatDate(tx.created) }}</span>
            </div>
            <div class="tx-amount" :class="{ negative: tx.gold_change < 0 }">
              {{ tx.gold_change > 0 ? '+' : '' }}{{ tx.gold_change }} Gold
            </div>
          </div>
        </div>
      </section>
    </main>

    <PurchaseModal
      :is-open="isModalOpen"
      :item="selectedItem"
      @close="closeModal"
      @success="onPurchaseSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useShopStore } from '../features/shop/store'
import GoldDisplay from '../features/shop/components/GoldDisplay.vue'
import GoldPackageCard from '../features/shop/components/GoldPackageCard.vue'
import BundleCard from '../features/shop/components/BundleCard.vue'
import PurchaseModal from '../features/shop/components/PurchaseModal.vue'
import type { ShopItem, PurchaseResponse } from '../features/shop/types'

const shopStore = useShopStore()
const activeTab = ref('gold')
const isModalOpen = ref(false)
const selectedItem = ref<ShopItem | null>(null)

const tabs = [
  { id: 'gold', label: 'Gold' },
  { id: 'bundles', label: 'Bundles' },
  { id: 'history', label: 'History' }
]

onMounted(() => {
  loadTabData()
})

function loadTabData() {
  shopStore.fetchGoldBalance()
  
  if (activeTab.value === 'gold') {
    shopStore.fetchGoldPackages()
  } else if (activeTab.value === 'bundles') {
    shopStore.fetchBundles()
  } else if (activeTab.value === 'history') {
    shopStore.fetchTransactions()
  }
}

function openPurchaseModal(itemId: number) {
  // Find item in either list
  const item = [...shopStore.goldPackages, ...shopStore.bundles]
    .find(i => i.id === itemId)
  
  if (item) {
    selectedItem.value = item
    isModalOpen.value = true
  }
}

function closeModal() {
  isModalOpen.value = false
  selectedItem.value = null
}

function onPurchaseSuccess(result: PurchaseResponse) {
  // Could show toast notification here
  console.log('Purchase successful:', result)
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString()
}
</script>

<style scoped>
.shop-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem;
}

.shop-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.shop-header h1 {
  margin: 0;
  font-size: 1.875rem;
  color: white;
}

.shop-tabs {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding-bottom: 0.5rem;
}

.tab {
  padding: 0.625rem 1.25rem;
  background: none;
  border: none;
  color: #9ca3af;
  font-weight: 500;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s;
}

.tab:hover {
  background: rgba(255, 255, 255, 0.05);
  color: white;
}

.tab.active {
  background: #3b82f6;
  color: white;
}

.tab-content h2 {
  margin: 0 0 0.5rem 0;
  color: white;
  font-size: 1.5rem;
}

.tab-description {
  margin: 0 0 1.5rem 0;
  color: #9ca3af;
}

.gold-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 1.5rem;
}

.bundles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.loading, .error, .empty {
  text-align: center;
  padding: 3rem;
  color: #9ca3af;
}

.error {
  color: #ef4444;
}

.transactions-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.transaction {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  border-left: 3px solid #4b5563;
}

.transaction.completed {
  border-left-color: #22c55e;
}

.transaction.failed {
  border-left-color: #ef4444;
}

.tx-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.tx-name {
  color: white;
  font-weight: 500;
}

.tx-date {
  color: #6b7280;
  font-size: 0.75rem;
}

.tx-amount {
  color: #22c55e;
  font-weight: 600;
}

.tx-amount.negative {
  color: #ef4444;
}
</style>

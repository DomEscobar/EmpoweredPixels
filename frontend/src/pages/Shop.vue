<template>
  <div class="shop-page" data-testid="shop-page">
    <header class="shop-header" data-testid="shop-header">
      <h1>🏪 Shop</h1>
      <GoldDisplay />
    </header>

    <nav class="shop-tabs" data-testid="shop-tabs">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        :class="['tab', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
        :data-testid="`tab-${tab.id}`"
      >
        {{ tab.label }}
      </button>
    </nav>

    <main class="shop-content">
      <!-- Gold Packages Tab -->
      <section v-if="activeTab === 'gold'" class="tab-content" data-testid="tab-gold">
        <h2>Gold Emporium</h2>
        <p class="tab-description">Purchase gold for upgrades and items</p>
        
        <div v-if="shopStore.loading" class="gold-grid">
          <div v-for="i in 4" :key="i" class="skeleton-card">
            <Skeleton height="160px" border-radius="16px" />
            <div class="skeleton-footer">
              <Skeleton width="40%" height="24px" />
              <Skeleton width="30%" height="32px" />
            </div>
          </div>
        </div>
        
        <div v-else-if="shopStore.error" class="error" data-testid="error-message">
          {{ shopStore.error }}
        </div>

        <div v-else-if="shopStore.goldPackages.length === 0" class="empty-state-wrapper" data-testid="empty-gold">
          <EmptyState 
            icon="🪙" 
            title="No Gold Packages" 
            message="The gold emporium is currently empty. Check back later!" 
          />
        </div>
        
        <div v-else class="gold-grid" data-testid="gold-grid">
          <GoldPackageCard 
            v-for="item in shopStore.goldPackages" 
            :key="item.id"
            :item="item"
            @purchase="openPurchaseModal"
            :data-testid="`gold-package-${item.id}`"
          />
        </div>
      </section>

      <!-- Bundles Tab -->
      <section v-if="activeTab === 'bundles'" class="tab-content" data-testid="tab-bundles">
        <h2>Equipment Bundles</h2>
        <p class="tab-description">Curated bundles with guaranteed rarity items</p>
        
        <div v-if="shopStore.loading" class="bundles-grid">
          <div v-for="i in 3" :key="i" class="skeleton-card">
            <Skeleton height="200px" border-radius="12px" />
            <div class="skeleton-footer">
              <Skeleton width="40%" height="24px" />
              <Skeleton width="30%" height="32px" />
            </div>
          </div>
        </div>
        
        <div v-else-if="shopStore.error" class="error" data-testid="error-message">
          {{ shopStore.error }}
        </div>

        <div v-else-if="shopStore.bundles.length === 0" class="empty-state-wrapper" data-testid="empty-bundles">
          <EmptyState 
            icon="📦" 
            title="No Bundles Available" 
            message="We're currently out of equipment bundles. Our blacksmith is working hard!" 
          />
        </div>
        
        <div v-else class="bundles-grid" data-testid="bundles-grid">
          <BundleCard 
            v-for="item in shopStore.bundles" 
            :key="item.id"
            :item="item"
            @purchase="openPurchaseModal"
            :data-testid="`bundle-${item.id}`"
          />
        </div>
      </section>

      <!-- History Tab -->
      <section v-if="activeTab === 'history'" class="tab-content" data-testid="tab-history">
        <h2>Purchase History</h2>
        <p class="tab-description">Your recent transactions</p>
        
        <div v-if="shopStore.loading" class="transactions-list">
          <Skeleton v-for="i in 5" :key="i" height="64px" border-radius="8px" />
        </div>

        <div v-else-if="shopStore.transactions.length === 0" class="empty-state-wrapper" data-testid="empty-history">
          <EmptyState 
            icon="📜" 
            title="No History" 
            message="You haven't made any purchases yet." 
          />
        </div>
        
        <div v-else class="transactions-list" data-testid="transactions-list">
          <div 
            v-for="tx in shopStore.transactions" 
            :key="tx.id"
            :class="['transaction', tx.status]"
            :data-testid="`tx-${tx.id}`"
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

    <ToastManager ref="toastManager" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useShopStore } from '../features/shop/store'
import GoldDisplay from '../features/shop/components/GoldDisplay.vue'
import GoldPackageCard from '../features/shop/components/GoldPackageCard.vue'
import BundleCard from '../features/shop/components/BundleCard.vue'
import PurchaseModal from '../features/shop/components/PurchaseModal.vue'
import Skeleton from '../shared/ui/Skeleton.vue'
import ToastManager from '../shared/ui/ToastManager.vue'
import EmptyState from '../shared/ui/EmptyState.vue'
import type { ShopItem, PurchaseResponse } from '../features/shop/types'

const shopStore = useShopStore()
const activeTab = ref('gold')
const isModalOpen = ref(false)
const selectedItem = ref<ShopItem | null>(null)
const toastManager = ref<any>(null)

const tabs = [
  { id: 'gold', label: 'Gold' },
  { id: 'bundles', label: 'Bundles' },
  { id: 'history', label: 'History' }
]

onMounted(() => {
  loadTabData()
})

// Watch tab changes to reload data
watch(activeTab, () => {
  loadTabData()
})

function loadTabData() {
  shopStore.fetchGoldBalance()
  
  if (activeTab.value === 'gold') {
    if (shopStore.goldPackages.length === 0) shopStore.fetchGoldPackages()
  } else if (activeTab.value === 'bundles') {
    if (shopStore.bundles.length === 0) shopStore.fetchBundles()
  } else if (activeTab.value === 'history') {
    shopStore.fetchTransactions()
  }
}

function openPurchaseModal(itemId: number) {
  // Find item in either list
  const item = [...shopStore.goldPackages, ...shopStore.bundles]
    .find(i => i.id === itemId)
  
  if (item) {
    // Client-side check for gold
    if (item.price_currency === 'gold' && shopStore.hasInsufficientGold(item)) {
      toastManager.value?.addToast(`Insufficient Gold to purchase ${item.name}`, 'error')
      return
    }

    selectedItem.value = item
    isModalOpen.value = true
  }
}

function closeModal() {
  isModalOpen.value = false
  selectedItem.value = null
}

function onPurchaseSuccess(result: PurchaseResponse) {
  toastManager.value?.addToast(result.message || 'Purchase successful!', 'success')
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

.empty-state-wrapper {
  margin: 2rem 0;
}

.skeleton-card {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.skeleton-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

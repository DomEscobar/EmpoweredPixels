<template>
  <div class="gold-package" :class="`tier-${tier}`" :data-testid="`gold-package-card-${item.id}`">
    <div class="package-header">
      <span v-if="isPopular" class="popular-badge">Most Popular</span>
      <span v-if="isBestValue" class="value-badge">Best Value</span>
      <h3 class="package-name">{{ item.name }}</h3>
    </div>
    
    <div class="gold-amount">
      <span class="gold-icon">🪙</span>
      <span class="amount">{{ item.gold_amount?.toLocaleString() }}</span>
      <span class="gold-label">Gold</span>
    </div>

    <div v-if="bonusPercent > 0" class="bonus-badge">
      +{{ bonusPercent }}% Bonus
    </div>

    <div class="package-footer">
      <span class="price">{{ formattedPrice }}</span>
      <button 
        class="buy-button"
        :disabled="shopStore.purchaseInProgress"
        @click="$emit('purchase', item.id)"
        :data-testid="`buy-gold-${item.id}`"
      >
        {{ shopStore.purchaseInProgress ? 'Processing...' : 'Buy' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ShopItem } from '../types'
import { useShopStore } from '../store'

const props = defineProps<{
  item: ShopItem
}>()

defineEmits<{
  purchase: [id: number]
}>()

const shopStore = useShopStore()

const tier = computed(() => {
  // Map gold amounts to tiers
  const amount = props.item.gold_amount || 0
  if (amount >= 6500) return 'legendary'
  if (amount >= 1200) return 'epic'
  if (amount >= 550) return 'rare'
  return 'common'
})

const isPopular = computed(() => props.item.name.includes('Satchel'))
const isBestValue = computed(() => props.item.name.includes('Hoard'))

const bonusPercent = computed(() => {
  // Calculate bonus from metadata or hardcoded
  const metadata = props.item.metadata || {}
  if (metadata.bonus_percent) return metadata.bonus_percent
  
  // Hardcoded based on seed data
  if (props.item.name.includes('Hoard')) return 30
  if (props.item.name.includes('Chest')) return 20
  if (props.item.name.includes('Satchel')) return 10
  return 0
})

const formattedPrice = computed(() => shopStore.formatItemPrice(props.item))
</script>

<style scoped>
.gold-package {
  border-radius: 16px;
  padding: 1.5rem;
  text-align: center;
  background: linear-gradient(135deg, #374151 0%, #1f2937 100%);
  border: 2px solid #4b5563;
  transition: transform 0.2s, box-shadow 0.2s;
  position: relative;
}

.gold-package:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
}

.tier-rare {
  border-color: #3b82f6;
  background: linear-gradient(135deg, #1e3a8a 0%, #1e40af 100%);
}

.tier-epic {
  border-color: #a855f7;
  background: linear-gradient(135deg, #581c87 0%, #7c3aed 100%);
}

.tier-legendary {
  border-color: #f59e0b;
  background: linear-gradient(135deg, #92400e 0%, #d97706 100%);
}

.popular-badge, .value-badge {
  position: absolute;
  top: -10px;
  left: 50%;
  transform: translateX(-50%);
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.625rem;
  font-weight: 700;
  text-transform: uppercase;
  white-space: nowrap;
}

.popular-badge {
  background: #3b82f6;
  color: white;
}

.value-badge {
  background: #f59e0b;
  color: white;
}

.package-header {
  margin-bottom: 1rem;
}

.package-name {
  margin: 0.5rem 0 0 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: white;
}

.gold-amount {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  margin: 1.5rem 0;
}

.gold-icon {
  font-size: 3rem;
}

.amount {
  font-size: 2.5rem;
  font-weight: 800;
  color: #fbbf24;
  line-height: 1;
}

.gold-label {
  font-size: 0.875rem;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

.bonus-badge {
  display: inline-block;
  padding: 0.375rem 0.75rem;
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.package-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.price {
  font-size: 1.25rem;
  font-weight: 700;
  color: white;
}

.buy-button {
  padding: 0.625rem 1.5rem;
  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.buy-button:hover:not(:disabled) {
  background: linear-gradient(135deg, #16a34a 0%, #15803d 100%);
  transform: scale(1.05);
}

.buy-button:disabled {
  background: #4b5563;
  cursor: not-allowed;
  opacity: 0.7;
}
</style>

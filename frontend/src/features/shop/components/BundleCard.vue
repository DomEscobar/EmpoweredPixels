<template>
  <div class="bundle-card" :style="{ borderColor: rarityColor }">
    <div class="bundle-header" :style="{ backgroundColor: rarityColor + '20' }">
      <span class="rarity-badge" :style="{ backgroundColor: rarityColor }">
        {{ rarityName }}
      </span>
      <h3 class="bundle-name">{{ item.name }}</h3>
    </div>
    
    <div class="bundle-content">
      <p class="bundle-description">{{ item.description }}</p>
      
      <div v-if="item.gold_amount" class="bonus-gold">
        <span class="bonus-icon">+</span>
        <span class="bonus-amount">{{ item.gold_amount.toLocaleString() }} Gold</span>
      </div>

      <div v-if="item.metadata?.drop_boosts" class="bonus-item">
        <span class="bonus-icon">🎁</span>
        <span>{{ item.metadata.drop_boosts }} Drop Boosts</span>
      </div>

      <div v-if="item.metadata?.equipment_count" class="bonus-item">
        <span class="bonus-icon">⚔️</span>
        <span>{{ item.metadata.equipment_count }} Equipment Items</span>
      </div>
    </div>

    <div class="bundle-footer">
      <span class="price">{{ formattedPrice }}</span>
      <button 
        class="buy-button"
        :disabled="isDisabled"
        @click="$emit('purchase', item.id)"
      >
        {{ buttonText }}
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

const rarityColor = computed(() => shopStore.getRarityColor(props.item.rarity))
const rarityName = computed(() => shopStore.getRarityName(props.item.rarity))
const formattedPrice = computed(() => shopStore.formatItemPrice(props.item))

const isDisabled = computed(() => {
  if (shopStore.purchaseInProgress) return true
  if (props.item.price_currency === 'gold') {
    return shopStore.hasInsufficientGold(props.item)
  }
  return false
})

const buttonText = computed(() => {
  if (shopStore.purchaseInProgress) return 'Processing...'
  if (props.item.price_currency === 'gold' && shopStore.hasInsufficientGold(props.item)) {
    return 'Insufficient Gold'
  }
  return 'Buy Now'
})
</script>

<style scoped>
.bundle-card {
  border: 2px solid;
  border-radius: 12px;
  overflow: hidden;
  background: #1f2937;
  transition: transform 0.2s, box-shadow 0.2s;
}

.bundle-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
}

.bundle-header {
  padding: 1rem;
  position: relative;
}

.rarity-badge {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.625rem;
  font-weight: 700;
  text-transform: uppercase;
  color: white;
}

.bundle-name {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: white;
  padding-right: 4rem;
}

.bundle-content {
  padding: 1rem;
}

.bundle-description {
  margin: 0 0 1rem 0;
  color: #9ca3af;
  font-size: 0.875rem;
  line-height: 1.5;
}

.bonus-gold, .bonus-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  margin-bottom: 0.5rem;
  color: #fbbf24;
  font-weight: 600;
}

.bonus-icon {
  font-size: 1.25rem;
}

.bundle-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.price {
  font-size: 1.25rem;
  font-weight: 700;
  color: white;
}

.buy-button {
  padding: 0.5rem 1.5rem;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.buy-button:hover:not(:disabled) {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
  transform: scale(1.05);
}

.buy-button:disabled {
  background: #4b5563;
  cursor: not-allowed;
  opacity: 0.7;
}
</style>

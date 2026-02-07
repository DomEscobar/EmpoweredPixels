<template>
  <div class="gold-display" :class="{ 'is-loading': loading }" data-testid="gold-display">
    <div class="gold-icon">
      <span class="icon">🪙</span>
    </div>
    <div class="gold-info">
      <span class="gold-amount">{{ formattedBalance }}</span>
      <span class="gold-label">Gold</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useShopStore } from '../store'

const shopStore = useShopStore()

const formattedBalance = computed(() => {
  const balance = shopStore.goldBalance?.balance || 0
  return balance.toLocaleString()
})

const loading = computed(() => shopStore.loading)
</script>

<style scoped>
.gold-display {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  border-radius: 9999px;
  color: white;
  font-weight: 600;
  transition: opacity 0.2s;
}

.gold-display.is-loading {
  opacity: 0.7;
}

.gold-icon {
  font-size: 1.25rem;
}

.gold-info {
  display: flex;
  flex-direction: column;
  line-height: 1;
}

.gold-amount {
  font-size: 1rem;
}

.gold-label {
  font-size: 0.625rem;
  opacity: 0.9;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
</style>

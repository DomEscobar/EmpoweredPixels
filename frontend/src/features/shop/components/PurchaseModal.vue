<template>
  <div v-if="isOpen" class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Confirm Purchase</h2>
        <button class="close-button" @click="close">&times;</button>
      </div>

      <div v-if="item" class="modal-body">
        <div class="item-preview" :style="{ borderColor: rarityColor }">
          <h3>{{ item.name }}</h3>
          <p>{{ item.description }}</p>
          <div class="price-tag">{{ formattedPrice }}</div>
        </div>

        <div class="balance-info">
          <div class="balance-row">
            <span>Current Balance:</span>
            <span class="balance">{{ currentBalance.toLocaleString() }} Gold</span>
          </div>
          <div v-if="isGoldPurchase" class="balance-row">
            <span>After Purchase:</span>
            <span :class="['balance', isAffordable ? 'positive' : 'negative']">
              {{ (currentBalance - item.price_amount).toLocaleString() }} Gold
            </span>
          </div>
        </div>

        <div v-if="result" :class="['result-message', result.success ? 'success' : 'error']">
          {{ result.message }}
        </div>

        <div class="modal-actions">
          <button class="cancel-button" @click="close">Cancel</button>
          <button 
            class="confirm-button"
            :disabled="!canPurchase || shopStore.purchaseInProgress"
            @click="confirmPurchase"
          >
            <span v-if="shopStore.purchaseInProgress">Processing...</span>
            <span v-else>Confirm Purchase</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { ShopItem, PurchaseResponse } from '../types'
import { useShopStore } from '../store'

const props = defineProps<{
  isOpen: boolean
  item: ShopItem | null
}>()

const emit = defineEmits<{
  close: []
  success: [result: PurchaseResponse]
}>()

const shopStore = useShopStore()
const result = ref<PurchaseResponse | null>(null)

const rarityColor = computed(() => {
  if (!props.item) return '#9ca3af'
  return shopStore.getRarityColor(props.item.rarity)
})

const formattedPrice = computed(() => {
  if (!props.item) return ''
  return shopStore.formatItemPrice(props.item)
})

const currentBalance = computed(() => 
  shopStore.goldBalance?.balance || 0
)

const isGoldPurchase = computed(() => 
  props.item?.price_currency === 'gold'
)

const isAffordable = computed(() => {
  if (!isGoldPurchase.value || !props.item) return true
  return currentBalance.value >= props.item.price_amount
})

const canPurchase = computed(() => {
  if (shopStore.purchaseInProgress) return false
  if (isGoldPurchase.value) return isAffordable.value
  return true // USD purchases always possible (mock)
})

async function confirmPurchase() {
  if (!props.item || !canPurchase.value) return
  
  result.value = await shopStore.purchaseItem(props.item.id)
  
  if (result.value.success) {
    setTimeout(() => {
      emit('success', result.value!)
      close()
    }, 1500)
  }
}

function close() {
  result.value = null
  emit('close')
}

// Reset result when modal opens/closes
watch(() => props.isOpen, (isOpen) => {
  if (!isOpen) result.value = null
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 1rem;
}

.modal-content {
  background: #1f2937;
  border-radius: 16px;
  width: 100%;
  max-width: 420px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h2 {
  margin: 0;
  font-size: 1.25rem;
  color: white;
}

.close-button {
  background: none;
  border: none;
  color: #9ca3af;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.2s;
}

.close-button:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.modal-body {
  padding: 1.5rem;
}

.item-preview {
  border: 2px solid;
  border-radius: 12px;
  padding: 1.25rem;
  text-align: center;
  margin-bottom: 1.5rem;
}

.item-preview h3 {
  margin: 0 0 0.5rem 0;
  color: white;
  font-size: 1.125rem;
}

.item-preview p {
  margin: 0 0 1rem 0;
  color: #9ca3af;
  font-size: 0.875rem;
}

.price-tag {
  font-size: 1.5rem;
  font-weight: 700;
  color: #fbbf24;
}

.balance-info {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
}

.balance-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  color: #9ca3af;
}

.balance-row:not(:last-child) {
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.balance {
  font-weight: 600;
  color: white;
}

.balance.positive {
  color: #22c55e;
}

.balance.negative {
  color: #ef4444;
}

.result-message {
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  text-align: center;
  font-weight: 500;
}

.result-message.success {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.result-message.error {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.modal-actions {
  display: flex;
  gap: 1rem;
}

.cancel-button, .confirm-button {
  flex: 1;
  padding: 0.75rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.cancel-button {
  background: rgba(255, 255, 255, 0.1);
  color: #9ca3af;
  border: none;
}

.cancel-button:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.confirm-button {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  border: none;
}

.confirm-button:hover:not(:disabled) {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
}

.confirm-button:disabled {
  background: #4b5563;
  cursor: not-allowed;
  opacity: 0.7;
}
</style>

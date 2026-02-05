<template>
  <Teleport to="body">
    <Transition name="modal">
      <div 
        v-if="visible" 
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="$emit('close')"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" />

        <!-- Modal -->
        <div class="modal-content relative w-full max-w-md bg-slate-900 border-4 border-slate-700 rounded-lg shadow-2xl overflow-hidden">
          <!-- Header -->
          <div class="p-4 border-b border-slate-700 flex items-center justify-between">
            <h2 class="text-xl font-bold">Confirm Purchase</h2>
            <button 
              class="text-slate-400 hover:text-white transition-colors"
              @click="$emit('close')"
            >
              ✕
            </button>
          </div>

          <!-- Content -->
          <div class="p-6 space-y-4">
            <!-- Item Preview -->
            <div class="flex items-center gap-4 p-4 bg-slate-800/50 rounded-lg border border-slate-700">
              <div class="w-16 h-16 rounded-lg flex items-center justify-center" :class="iconBgClass">
                <img :src="itemIcon" alt="" class="w-12 h-12 pixelated" />
              </div>
              <div class="flex-1">
                <h3 class="font-bold text-lg">{{ item?.name }}</h3>
                <p class="text-sm text-slate-400">{{ item?.description }}</p>
              </div>
            </div>

            <!-- Price Summary -->
            <div class="space-y-2">
              <div class="flex items-center justify-between text-sm">
                <span class="text-slate-400">Price</span>
                <span class="font-bold">{{ formattedPrice }}</span>
              </div>
              <div v-if="isGoldPurchase" class="flex items-center justify-between text-sm">
                <span class="text-slate-400">Current Balance</span>
                <span class="font-bold text-amber-400">{{ currentBalance.toLocaleString() }} Gold</span>
              </div>
              <div v-if="isGoldPurchase" class="flex items-center justify-between text-sm border-t border-slate-700 pt-2">
                <span class="text-slate-400">After Purchase</span>
                <span class="font-bold" :class="canAfford ? 'text-green-400' : 'text-red-400'">
                  {{ afterPurchaseBalance.toLocaleString() }} Gold
                </span>
              </div>
              <div v-if="item?.gold_amount" class="flex items-center justify-between text-sm border-t border-slate-700 pt-2">
                <span class="text-slate-400">You'll Receive</span>
                <div class="flex items-center gap-1">
                  <img :src="GOLD_ICON" alt="" class="w-4 h-4 pixelated" />
                  <span class="font-bold text-amber-400">{{ item.gold_amount.toLocaleString() }} Gold</span>
                </div>
              </div>
            </div>

            <!-- Warning for insufficient funds -->
            <div 
              v-if="isGoldPurchase && !canAfford" 
              class="p-3 bg-red-900/30 border border-red-600/50 rounded-lg text-sm text-red-400"
            >
              ⚠️ You don't have enough gold for this purchase.
            </div>
          </div>

          <!-- Actions -->
          <div class="p-4 border-t border-slate-700 flex gap-3">
            <button 
              class="flex-1 py-3 rounded-lg font-bold uppercase tracking-wide bg-slate-700 hover:bg-slate-600 transition-colors"
              @click="$emit('close')"
            >
              Cancel
            </button>
            <button 
              class="flex-1 py-3 rounded-lg font-bold uppercase tracking-wide transition-all"
              :class="confirmButtonClass"
              :disabled="isGoldPurchase && !canAfford || loading"
              @click="handleConfirm"
            >
              <span v-if="loading" class="animate-spin">⏳</span>
              <span v-else>{{ confirmText }}</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { ShopItem, formatPrice } from '../api';
import { useShopStore } from '../store';

const props = defineProps<{
  visible: boolean;
  item: ShopItem | null;
  loading?: boolean;
}>();

const emit = defineEmits<{
  close: [];
  confirm: [item: ShopItem];
}>();

const GOLD_ICON = 'https://vibemedia.space/icon_gold_coin_nav_8f7e6d.png?prompt=golden%20coin%20with%20shine%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON';
const BUNDLE_ICON = 'https://vibemedia.space/icon_bundle_common_1a2b3c.png?prompt=wooden%20treasure%20chest%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON';

const shop = useShopStore();

const isGoldPurchase = computed(() => props.item?.price_currency === 'gold');
const currentBalance = computed(() => shop.goldBalance);
const afterPurchaseBalance = computed(() => {
  if (!props.item || !isGoldPurchase.value) return currentBalance.value;
  return currentBalance.value - props.item.price_amount;
});
const canAfford = computed(() => props.item ? shop.canAfford(props.item) : false);

const formattedPrice = computed(() => {
  if (!props.item) return '';
  return formatPrice(props.item.price_amount, props.item.price_currency);
});

const itemIcon = computed(() => {
  if (props.item?.item_type === 'gold_package') return GOLD_ICON;
  return BUNDLE_ICON;
});

const iconBgClass = computed(() => {
  if (props.item?.item_type === 'gold_package') return 'bg-amber-900/50';
  return 'bg-indigo-900/50';
});

const confirmButtonClass = computed(() => {
  if (props.loading) return 'bg-slate-600 cursor-wait';
  if (isGoldPurchase.value && !canAfford.value) {
    return 'bg-slate-700 text-slate-500 cursor-not-allowed';
  }
  if (props.item?.item_type === 'gold_package') {
    return 'bg-indigo-600 hover:bg-indigo-500 text-white';
  }
  return 'bg-amber-600 hover:bg-amber-500 text-white';
});

const confirmText = computed(() => {
  if (props.item?.price_currency === 'usd') {
    return 'Purchase';
  }
  return 'Buy Now';
});

function handleConfirm() {
  if (props.item && (!isGoldPurchase.value || canAfford.value)) {
    emit('confirm', props.item);
  }
}
</script>

<style scoped>
.pixelated {
  image-rendering: pixelated;
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-content,
.modal-leave-to .modal-content {
  transform: scale(0.95) translateY(20px);
}
</style>

<template>
  <div 
    class="gold-display flex items-center gap-2 px-3 py-1.5 rounded"
    :class="displayClass"
  >
    <img 
      :src="GOLD_ICON" 
      alt="Gold" 
      class="w-5 h-5 pixelated"
    />
    <span class="font-bold text-sm">
      {{ formattedBalance }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useShopStore } from '../store';

const props = defineProps<{
  compact?: boolean;
}>();

const GOLD_ICON = 'https://vibemedia.space/icon_gold_coin_nav_8f7e6d.png?prompt=golden%20coin%20with%20shine%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON';

const shop = useShopStore();

const formattedBalance = computed(() => {
  const balance = shop.goldBalance;
  if (balance >= 10000) {
    return `${(balance / 1000).toFixed(1)}K`;
  }
  return balance.toLocaleString();
});

const displayClass = computed(() => {
  if (props.compact) {
    return 'bg-amber-900/30 border border-amber-600/50';
  }
  return 'bg-gradient-to-r from-amber-900/40 to-yellow-900/40 border-2 border-amber-600/70 shadow-lg shadow-amber-500/10';
});
</script>

<style scoped>
.pixelated {
  image-rendering: pixelated;
}

.gold-display {
  text-shadow: 1px 1px 0 rgba(0, 0, 0, 0.5);
}
</style>

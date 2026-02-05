<template>
  <div 
    class="gold-package-card relative overflow-hidden rounded-lg border-3 transition-all duration-200 hover:scale-[1.02] cursor-pointer"
    :class="cardClasses"
    @click="$emit('select', item)"
  >
    <!-- Popular Badge -->
    <div 
      v-if="isPopular" 
      class="absolute top-2 right-2 bg-amber-500 text-amber-950 text-xs font-bold px-2 py-1 rounded"
    >
      POPULAR
    </div>

    <!-- Best Value Badge -->
    <div 
      v-if="isBestValue" 
      class="absolute top-2 right-2 bg-green-500 text-green-950 text-xs font-bold px-2 py-1 rounded"
    >
      BEST VALUE
    </div>

    <!-- Content -->
    <div class="p-4 space-y-3">
      <!-- Icon & Name -->
      <div class="flex items-center gap-3">
        <div class="gold-icon-container w-14 h-14 rounded-lg flex items-center justify-center" :class="iconBgClass">
          <img :src="goldIcon" alt="" class="w-10 h-10 pixelated animate-bounce-slow" />
        </div>
        <div>
          <h3 class="font-bold text-lg">{{ item.name }}</h3>
          <p class="text-xs opacity-70">{{ item.description }}</p>
        </div>
      </div>

      <!-- Gold Amount -->
      <div class="gold-amount text-center py-3 rounded-lg" :class="amountBgClass">
        <div class="flex items-center justify-center gap-2">
          <img :src="goldCoinIcon" alt="" class="w-6 h-6 pixelated" />
          <span class="text-2xl font-bold text-amber-400">
            {{ formattedGoldAmount }}
          </span>
        </div>
        <p v-if="bonusPercent" class="text-xs text-green-400 mt-1">
          +{{ bonusPercent }}% Bonus!
        </p>
      </div>

      <!-- Price -->
      <button
        class="w-full py-3 rounded-lg font-bold text-lg uppercase tracking-wide transition-all bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 text-white shadow-lg hover:shadow-indigo-500/30"
        @click.stop="$emit('purchase', item)"
      >
        {{ formattedPrice }}
      </button>
    </div>

    <!-- Coin particles for larger packages -->
    <div v-if="item.gold_amount && item.gold_amount >= 1000" class="absolute inset-0 pointer-events-none overflow-hidden">
      <div 
        v-for="i in 5" 
        :key="i"
        class="particle absolute w-2 h-2 bg-amber-400 rounded-full opacity-50"
        :style="particleStyle(i)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { ShopItem, formatPrice } from '../api';

const props = defineProps<{
  item: ShopItem;
}>();

defineEmits<{
  select: [item: ShopItem];
  purchase: [item: ShopItem];
}>();

const ICONS = {
  GOLD_SMALL: 'https://vibemedia.space/icon_gold_pile_small_3a4b5c.png?prompt=small%20pile%20of%20gold%20coins%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  GOLD_MEDIUM: 'https://vibemedia.space/icon_gold_bag_6d7e8f.png?prompt=leather%20bag%20of%20gold%20coins%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  GOLD_LARGE: 'https://vibemedia.space/icon_gold_chest_9g0h1i.png?prompt=open%20treasure%20chest%20full%20of%20gold%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  GOLD_VAULT: 'https://vibemedia.space/icon_gold_vault_2j3k4l.png?prompt=overflowing%20vault%20of%20gold%20treasure%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  GOLD_COIN: 'https://vibemedia.space/icon_gold_coin_nav_8f7e6d.png?prompt=golden%20coin%20with%20shine%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
};

const goldIcon = computed(() => {
  const amount = props.item.gold_amount ?? 0;
  if (amount >= 5000) return ICONS.GOLD_VAULT;
  if (amount >= 1000) return ICONS.GOLD_LARGE;
  if (amount >= 500) return ICONS.GOLD_MEDIUM;
  return ICONS.GOLD_SMALL;
});

const goldCoinIcon = ICONS.GOLD_COIN;

const formattedGoldAmount = computed(() => {
  const amount = props.item.gold_amount ?? 0;
  return amount.toLocaleString();
});

const formattedPrice = computed(() => formatPrice(props.item.price_amount, props.item.price_currency));

const bonusPercent = computed(() => {
  const amount = props.item.gold_amount ?? 0;
  const basePrice = props.item.price_amount; // in cents
  
  // Calculate base rate from smallest package (100 gold = $0.99)
  const baseRate = 100 / 99; // gold per cent
  const actualRate = amount / basePrice;
  
  if (actualRate > baseRate) {
    const bonus = Math.round(((actualRate / baseRate) - 1) * 100);
    return bonus > 0 ? bonus : null;
  }
  return null;
});

const isPopular = computed(() => {
  const amount = props.item.gold_amount ?? 0;
  return amount === 550;
});

const isBestValue = computed(() => {
  const amount = props.item.gold_amount ?? 0;
  return amount >= 6000;
});

const cardClasses = computed(() => {
  if (isBestValue.value) {
    return 'border-green-500 bg-gradient-to-br from-slate-900 to-green-950/50';
  }
  if (isPopular.value) {
    return 'border-amber-500 bg-gradient-to-br from-slate-900 to-amber-950/50';
  }
  return 'border-slate-600 bg-gradient-to-br from-slate-900 to-slate-800';
});

const iconBgClass = computed(() => {
  if (isBestValue.value) return 'bg-green-900/50';
  if (isPopular.value) return 'bg-amber-900/50';
  return 'bg-slate-800/50';
});

const amountBgClass = computed(() => {
  return 'bg-slate-950/50';
});

const particleStyle = (index: number) => {
  const delay = index * 0.5;
  const left = 10 + (index * 18);
  return {
    left: `${left}%`,
    animation: `float ${2 + index * 0.3}s ease-in-out ${delay}s infinite`,
  };
};
</script>

<style scoped>
.pixelated {
  image-rendering: pixelated;
}

.border-3 {
  border-width: 3px;
}

.animate-bounce-slow {
  animation: bounce-slow 2s ease-in-out infinite;
}

@keyframes bounce-slow {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

@keyframes float {
  0%, 100% { 
    transform: translateY(100%) scale(0);
    opacity: 0;
  }
  50% { 
    transform: translateY(-50vh) scale(1);
    opacity: 0.6;
  }
}

.particle {
  bottom: 0;
}
</style>

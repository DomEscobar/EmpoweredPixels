<template>
  <div 
    class="bundle-card relative overflow-hidden rounded-lg transition-all duration-200 hover:scale-[1.02] cursor-pointer"
    :class="[rarityClasses, glowClasses]"
    @click="$emit('select', item)"
  >
    <!-- Rarity Banner -->
    <div 
      class="absolute top-0 left-0 right-0 h-1"
      :class="rarityBarClass"
    />

    <!-- Content -->
    <div class="p-4 space-y-3">
      <!-- Header -->
      <div class="flex items-start justify-between">
        <div>
          <h3 class="font-bold text-lg">{{ item.name }}</h3>
          <span class="text-xs uppercase tracking-wider opacity-70">
            {{ rarityName }}
          </span>
        </div>
        <div class="bundle-icon w-12 h-12 rounded flex items-center justify-center" :class="iconBgClass">
          <img :src="bundleIcon" alt="" class="w-8 h-8 pixelated" />
        </div>
      </div>

      <!-- Description -->
      <p class="text-sm opacity-80 line-clamp-2">
        {{ item.description }}
      </p>

      <!-- Metadata (equipment count, etc.) -->
      <div v-if="equipmentCount" class="flex items-center gap-2 text-xs opacity-70">
        <span>📦 {{ equipmentCount }} item{{ equipmentCount > 1 ? 's' : '' }}</span>
      </div>

      <!-- Price -->
      <div class="flex items-center justify-between pt-2 border-t border-current/20">
        <div class="flex items-center gap-2">
          <img :src="currencyIcon" alt="" class="w-5 h-5 pixelated" />
          <span class="font-bold text-lg">{{ formattedPrice }}</span>
        </div>
        <button
          class="purchase-btn px-4 py-2 rounded font-bold text-sm uppercase tracking-wide transition-all"
          :class="buttonClasses"
          :disabled="!canAfford"
          @click.stop="$emit('purchase', item)"
        >
          {{ canAfford ? 'Buy' : 'Not Enough' }}
        </button>
      </div>
    </div>

    <!-- Shine effect for high rarity -->
    <div 
      v-if="item.rarity >= 3" 
      class="absolute inset-0 pointer-events-none shine-effect"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { ShopItem, Rarity, getRarityName, getRarityColor, getRarityGlow, formatPrice } from '../api';
import { useShopStore } from '../store';

const props = defineProps<{
  item: ShopItem;
}>();

defineEmits<{
  select: [item: ShopItem];
  purchase: [item: ShopItem];
}>();

const ICONS = {
  GOLD: 'https://vibemedia.space/icon_gold_coin_nav_8f7e6d.png?prompt=golden%20coin%20with%20shine%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  BUNDLE_COMMON: 'https://vibemedia.space/icon_bundle_common_1a2b3c.png?prompt=wooden%20treasure%20chest%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  BUNDLE_RARE: 'https://vibemedia.space/icon_bundle_rare_4d5e6f.png?prompt=silver%20treasure%20chest%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  BUNDLE_FABLED: 'https://vibemedia.space/icon_bundle_fabled_7g8h9i.png?prompt=purple%20magic%20treasure%20chest%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  BUNDLE_MYTHIC: 'https://vibemedia.space/icon_bundle_mythic_0j1k2l.png?prompt=golden%20ornate%20treasure%20chest%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
};

const shop = useShopStore();

const rarityName = computed(() => getRarityName(props.item.rarity));
const rarityClasses = computed(() => getRarityColor(props.item.rarity));
const glowClasses = computed(() => {
  const glow = getRarityGlow(props.item.rarity);
  return glow ? `shadow-xl ${glow}` : '';
});

const rarityBarClass = computed(() => {
  switch (props.item.rarity) {
    case Rarity.Common: return 'bg-green-500';
    case Rarity.Rare: return 'bg-blue-500';
    case Rarity.Fabled: return 'bg-purple-500';
    case Rarity.Mythic: return 'bg-orange-500';
    case Rarity.Legendary: return 'bg-yellow-500';
    default: return 'bg-slate-500';
  }
});

const iconBgClass = computed(() => {
  switch (props.item.rarity) {
    case Rarity.Common: return 'bg-green-900/50';
    case Rarity.Rare: return 'bg-blue-900/50';
    case Rarity.Fabled: return 'bg-purple-900/50';
    case Rarity.Mythic: return 'bg-orange-900/50';
    case Rarity.Legendary: return 'bg-yellow-900/50';
    default: return 'bg-slate-900/50';
  }
});

const bundleIcon = computed(() => {
  switch (props.item.rarity) {
    case Rarity.Rare: return ICONS.BUNDLE_RARE;
    case Rarity.Fabled: return ICONS.BUNDLE_FABLED;
    case Rarity.Mythic: return ICONS.BUNDLE_MYTHIC;
    case Rarity.Legendary: return ICONS.BUNDLE_MYTHIC;
    default: return ICONS.BUNDLE_COMMON;
  }
});

const currencyIcon = computed(() => ICONS.GOLD);

const formattedPrice = computed(() => {
  if (props.item.price_currency === 'gold') {
    return props.item.price_amount.toLocaleString();
  }
  return formatPrice(props.item.price_amount, props.item.price_currency);
});

const canAfford = computed(() => shop.canAfford(props.item));

const buttonClasses = computed(() => {
  if (!canAfford.value) {
    return 'bg-slate-700 text-slate-500 cursor-not-allowed';
  }
  return 'bg-amber-600 hover:bg-amber-500 text-white shadow-md hover:shadow-lg';
});

const equipmentCount = computed(() => {
  return props.item.metadata?.equipment_count as number | undefined;
});
</script>

<style scoped>
.pixelated {
  image-rendering: pixelated;
}

.bundle-card {
  border-width: 3px;
}

.shine-effect {
  background: linear-gradient(
    135deg,
    transparent 40%,
    rgba(255, 255, 255, 0.1) 50%,
    transparent 60%
  );
  animation: shine 3s infinite;
}

@keyframes shine {
  0% { transform: translateX(-100%) translateY(-100%); }
  100% { transform: translateX(100%) translateY(100%); }
}

.purchase-btn:not(:disabled):active {
  transform: translateY(1px);
}
</style>

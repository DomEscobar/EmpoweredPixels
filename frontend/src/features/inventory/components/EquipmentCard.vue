<template>
  <div 
    class="relative group pixel-box-iron transition-all hover:-translate-y-1"
    :class="[rarityBorderColor, { 'is-favorite': item.isFavorite }]"
  >
    <!-- CRT Overlay for Item Slot -->
    <div class="pointer-events-none absolute inset-0 z-0 opacity-[0.05] bg-[repeating-linear-gradient(0deg,transparent,transparent_2px,rgba(255,255,255,0.1)_2px,rgba(255,255,255,0.1)_4px)]"></div>

    <!-- Favorite Indicator -->
    <div v-if="item.isFavorite" class="absolute -top-1 -right-1 z-20 text-yellow-400 drop-shadow-md">
      ★
    </div>

    <!-- Icon/Level/Enhancement -->
    <div class="flex flex-col items-center justify-center p-3 relative z-10 h-full">
      <div 
        class="h-16 w-16 mb-2 flex items-center justify-center relative p-1"
        :class="rarityIconBg"
      >
        <img :src="itemIcon" class="w-full h-full pixelated" :alt="item.type" />
        
        <div 
          v-if="item.enhancement > 0"
          class="absolute -top-2 -right-2 bg-indigo-600 text-white text-[10px] font-bold px-1.5 py-0.5 border border-indigo-400 shadow-sm"
        >
          +{{ item.enhancement }}
        </div>
      </div>
      
      <div class="text-[10px] text-slate-400 font-bold uppercase tracking-tighter">
        LVL {{ item.level }} {{ rarityName }}
      </div>
    </div>

    <!-- Ethereal Iron Tooltip -->
    <div class="item-tooltip pixel-box-iron">
      <div class="tooltip-header" :class="'rarity-' + item.rarity">
        <div class="flex justify-between items-center mb-1">
          <span class="text-xs font-bold">{{ rarityName }} {{ item.type }}</span>
          <span class="text-[10px] text-white/60">LVL {{ item.level }}</span>
        </div>
        <div v-if="item.enhancement > 0" class="enhance-text">+{{ item.enhancement }} Enhanced</div>
      </div>
      
      <div class="tooltip-stats">
        <div class="stat-row">
          <span class="label">Power Rating</span>
          <span class="value text-amber-400">{{ powerRating }}</span>
        </div>
        <!-- Derived stats based on item type and rarity for the mockup -->
        <div class="stat-row">
          <span class="label">Durability</span>
          <span class="value">100/100</span>
        </div>
      </div>

      <div class="tooltip-actions">
        <button @click.stop="$emit('enhance', item)" class="tip-btn enhance">ENHANCE</button>
        <button @click.stop="$emit('salvage', item)" class="tip-btn salvage">SALVAGE</button>
        <button @click.stop="$emit('toggle-favorite', item)" class="tip-btn favorite">
          {{ item.isFavorite ? 'UNFAVORITE' : 'FAVORITE' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Equipment } from '../api';

const props = defineProps<{
  item: Equipment
}>();

defineEmits<{
  (e: 'enhance', item: Equipment): void;
  (e: 'salvage', item: Equipment): void;
  (e: 'toggle-favorite', item: Equipment): void;
}>();

const itemIcon = computed(() => {
  const type = props.item.type.toLowerCase();
  if (type.includes('wpn')) return 'https://vibemedia.space/wpn_steel_sword_123.png?prompt=pixel%20steel%20sword%20sprite&style=pixel_game_asset&key=NOGON';
  if (type.includes('ring')) return 'https://vibemedia.space/ring_gold_456.png?prompt=pixel%20gold%20ring%20sprite&style=pixel_game_asset&key=NOGON';
  return 'https://vibemedia.space/item_generic_789.png?prompt=pixel%20mystery%20item%20sprite&style=pixel_game_asset&key=NOGON';
});

const powerRating = computed(() => {
  return (props.item.level * 10) + (props.item.rarity * 25) + (props.item.enhancement * 15);
});

const rarityConfig = {
  0: { name: 'BASIC', border: 'border-slate-800', icon: 'bg-slate-900 border-slate-700' },
  1: { name: 'COMMON', border: 'border-slate-700', icon: 'bg-slate-800 border-slate-700' },
  2: { name: 'RARE', border: 'border-blue-900', icon: 'bg-blue-900/40 border-blue-700' },
  3: { name: 'FABLED', border: 'border-purple-900', icon: 'bg-purple-900/40 border-purple-700' },
  4: { name: 'MYTHIC', border: 'border-amber-900', icon: 'bg-amber-900/40 border-amber-700' },
  5: { name: 'LEGENDARY', border: 'border-red-900', icon: 'bg-red-900/40 border-red-700' },
};

const config = computed(() => rarityConfig[props.item.rarity as keyof typeof rarityConfig] || rarityConfig[0]);

const rarityName = computed(() => config.value.name);
const rarityBorderColor = computed(() => config.value.border);
const rarityIconBg = computed(() => config.value.icon);
</script>

<style scoped>
.pixel-box-iron {
  background: rgba(15, 23, 42, 0.95);
  border: 3px solid #1e293b;
  box-shadow: 
    inset 0 0 0 1px #334155,
    0 4px 6px rgba(0, 0, 0, 0.3);
  image-rendering: pixelated;
  height: 100px;
  width: 100px;
  display: flex;
  margin: 0 auto;
}

.pixelated {
  image-rendering: pixelated;
}

/* Tooltip container */
.item-tooltip {
  position: absolute;
  top: -10px;
  left: 50%;
  transform: translate(-50%, -100%);
  width: 180px;
  z-index: 50;
  opacity: 0;
  pointer-events: none;
  transition: all 0.2s;
  height: auto;
  display: block;
}

.group:hover .item-tooltip {
  opacity: 1;
  pointer-events: auto;
  transform: translate(-50%, -110%);
}

.tooltip-header {
  padding: 0.5rem;
  background: #1e293b;
  border-bottom: 2px solid #334155;
}

.rarity-1 { border-top: 3px solid #94a3b8; }
.rarity-2 { border-top: 3px solid #3b82f6; }
.rarity-3 { border-top: 3px solid #a855f7; }
.rarity-4 { border-top: 3px solid #f59e0b; }
.rarity-5 { border-top: 3px solid #ef4444; }

.enhance-text {
  font-size: 8px;
  color: #818cf8;
  text-transform: uppercase;
  font-weight: bold;
}

.tooltip-stats {
  padding: 0.5rem;
  background: rgba(15, 23, 42, 0.9);
}

.stat-row {
  display: flex;
  justify-content: space-between;
  font-size: 10px;
  margin-bottom: 0.25rem;
}

.stat-row .label { color: #64748b; }
.stat-row .value { font-weight: bold; color: #e2e8f0; }

.tooltip-actions {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2px;
  padding: 0.25rem;
  background: #0f172a;
}

.tip-btn {
  font-size: 8px;
  padding: 0.25rem;
  font-weight: bold;
  text-transform: uppercase;
  border: 1px solid #334155;
  cursor: pointer;
  transition: all 0.1s;
}

.tip-btn:hover { background: #1e293b; border-color: #4a9eff; }
.tip-btn.enhance { color: #818cf8; }
.tip-btn.salvage { color: #ef4444; }
.tip-btn.favorite { color: #fbbf24; }
</style>

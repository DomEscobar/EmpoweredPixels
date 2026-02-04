<template>
  <div 
    class="relative group rounded-lg border bg-slate-800/50 p-4 transition-all hover:bg-slate-800"
    :class="[rarityBorderColor, { 'ring-2 ring-indigo-500': item.isFavorite }]"
  >
    <!-- Rarity/Level Header -->
    <div class="flex items-start justify-between mb-3">
      <div>
        <span 
          class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium border"
          :class="rarityBadgeColor"
        >
          {{ rarityName }}
        </span>
        <div class="mt-1 text-xs text-slate-400">Lvl {{ item.level }}</div>
      </div>
      <div class="flex gap-2">
        <button 
          @click.stop="$emit('toggle-favorite', item)"
          class="text-slate-500 hover:text-yellow-400 transition-colors"
          :class="{ 'text-yellow-400': item.isFavorite }"
          title="Toggle Favorite"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Icon/Enhancement -->
    <div class="flex items-center justify-center py-4 relative">
      <div 
        class="h-16 w-16 rounded-lg flex items-center justify-center text-3xl shadow-lg border"
        :class="rarityIconBg"
      >
        ⚔️
      </div>
      <div 
        v-if="item.enhancement > 0"
        class="absolute -right-2 top-2 bg-indigo-600 text-white text-xs font-bold px-2 py-0.5 rounded-full border border-indigo-400 shadow-sm"
      >
        +{{ item.enhancement }}
      </div>
    </div>

    <!-- Actions -->
    <div class="mt-4 grid grid-cols-2 gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
      <button 
        @click.stop="$emit('enhance', item)"
        class="flex items-center justify-center px-3 py-1.5 text-xs font-medium text-white bg-indigo-600 rounded hover:bg-indigo-500 transition-colors"
      >
        Enhance
      </button>
      <button 
        @click.stop="$emit('salvage', item)"
        class="flex items-center justify-center px-3 py-1.5 text-xs font-medium text-red-300 bg-red-900/30 rounded hover:bg-red-900/50 transition-colors border border-red-900/50"
      >
        Salvage
      </button>
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

const rarityConfig = {
  0: { name: 'Basic', border: 'border-slate-700', badge: 'bg-slate-700 text-slate-300 border-slate-600', icon: 'bg-slate-800 border-slate-700' },
  1: { name: 'Common', border: 'border-slate-600', badge: 'bg-slate-700 text-white border-slate-500', icon: 'bg-slate-700 border-slate-600' },
  2: { name: 'Rare', border: 'border-blue-500/30', badge: 'bg-blue-900/30 text-blue-300 border-blue-500/30', icon: 'bg-blue-900/20 border-blue-500/30' },
  3: { name: 'Fabled', border: 'border-purple-500/30', badge: 'bg-purple-900/30 text-purple-300 border-purple-500/30', icon: 'bg-purple-900/20 border-purple-500/30' },
  4: { name: 'Mythic', border: 'border-amber-500/30', badge: 'bg-amber-900/30 text-amber-300 border-amber-500/30', icon: 'bg-amber-900/20 border-amber-500/30' },
  5: { name: 'Legendary', border: 'border-red-500/30', badge: 'bg-red-900/30 text-red-300 border-red-500/30', icon: 'bg-red-900/20 border-red-500/30' },
};

const config = computed(() => rarityConfig[props.item.rarity as keyof typeof rarityConfig] || rarityConfig[0]);

const rarityName = computed(() => config.value.name);
const rarityBorderColor = computed(() => config.value.border);
const rarityBadgeColor = computed(() => config.value.badge);
const rarityIconBg = computed(() => config.value.icon);
</script>

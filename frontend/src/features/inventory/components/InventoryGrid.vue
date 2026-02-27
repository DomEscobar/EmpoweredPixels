<template>
  <div class="inventory-container" data-testid="inventory-grid">
    <!-- Filtering Bar -->
    <div class="pixel-box-iron bg-slate-900/90 p-4 mb-6 flex flex-wrap items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <h2 class="text-amber-400 font-bold tracking-wider">
          FILTER:
        </h2>
        <div class="flex gap-2">
          <button 
            v-for="filter in filters" 
            :key="filter.id"
            class="filter-btn"
            :class="{ 'active': activeFilter === filter.id }"
            :data-testid="'filter-' + filter.id"
            @click="activeFilter = filter.id"
          >
            {{ filter.label }}
          </button>
        </div>
      </div>
      
      <div class="text-slate-500 text-xs">
        ITEMS: <span class="text-amber-300">{{ filteredEquipment.length }}</span> / {{ equipment.length }}
      </div>
    </div>

    <!-- Inventory Grid -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center py-20">
      <div class="animate-spin text-4xl mb-4">
        ⚙️
      </div>
      <p class="text-amber-400 animate-pulse">
        SCANNING VAULT...
      </p>
    </div>

    <div v-else-if="filteredEquipment.length === 0" class="pixel-box-iron bg-slate-900/60 py-20 text-center">
      <p class="text-slate-500">
        NO ITEMS FOUND IN THIS SECTOR
      </p>
    </div>

    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
      <EquipmentCard 
        v-for="item in filteredEquipment" 
        :key="item.id" 
        :item="item"
        data-testid="inventory-item"
        @enhance="$emit('enhance', item)"
        @salvage="$emit('salvage', item)"
        @toggle-favorite="$emit('toggle-favorite', item)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Equipment } from '../api';
import EquipmentCard from './EquipmentCard.vue';

const props = defineProps<{
  equipment: Equipment[];
  isLoading: boolean;
  search?: string;
  sortBy?: string;
}>();

defineEmits<{
  (e: 'enhance', item: Equipment): void;
  (e: 'salvage', item: Equipment): void;
  (e: 'toggle-favorite', item: Equipment): void;
}>();

const activeFilter = ref('all');

const filters = [
  { id: 'all', label: 'ALL' },
  { id: 'weapon', label: 'WEAPONS' },
  { id: 'armor', label: 'ARMOR' }, // Placeholder for future armor types
  { id: 'consumable', label: 'CONSUMABLES' },
  { id: 'ring', label: 'RINGS' }
];

// Rarity helper
const rarityLevels = [
  { id: 0, name: 'Broken' },
  { id: 1, name: 'Common' },
  { id: 2, name: 'Uncommon' },
  { id: 3, name: 'Rare' },
  { id: 4, name: 'Epic' },
  { id: 5, name: 'Legendary' },
  { id: 6, name: 'Mythic' },
  { id: 7, name: 'Divine' }
];

const rarityName = (level: number): string => {
  const found = rarityLevels.find(r => r.id === level);
  return found?.name || `Rarity ${level}`;
};

const filteredEquipment = computed(() => {
  let items = [...props.equipment];
  
  // 1. Search Filter (keyword)
  if (props.search && props.search.trim()) {
    const query = props.search.toLowerCase().trim();
    items = items.filter(item => {
      // Search in item ID, type, or rarity
      return (
        item.type.toLowerCase().includes(query) ||
        item.id.toLowerCase().includes(query) ||
        rarityName(item.rarity).toLowerCase().includes(query)
      );
    });
  }
  
  // 2. Sort Order
  if (props.sortBy) {
    switch (props.sortBy) {
      case 'level-desc':
        items.sort((a, b) => b.level - a.level);
        break;
      case 'rarity-desc':
        items.sort((a, b) => b.rarity - a.rarity);
        break;
      case 'power-desc':
        // Approximate power by rarity + level + enhancement
        const power = (item: Equipment) => (item.rarity * 10) + item.level + (item.enhancement * 2);
        items.sort((a, b) => power(b) - power(a));
        break;
      case 'recent':
      default:
        // Sort by ID (assuming newer items have higher IDs or are added at end)
        items.sort((a, b) => b.id.localeCompare(a.id));
        break;
    }
  }
  
  // 3. Category Filter
  if (activeFilter.value === 'all') return items;
  
  return items.filter(item => {
    // Prefer backend category if available
    if (item.category && item.category !== 'unknown') {
        if (activeFilter.value === 'weapon') return item.category === 'weapon';
        if (activeFilter.value === 'armor') return item.category === 'armor';
        if (activeFilter.value === 'ring') return item.category === 'trinket';
    }

    const type = item.type.toLowerCase();
    // Fallback patterns for items without explicit category
    if (activeFilter.value === 'weapon') return type.startsWith('wpn') || type.includes('sword') || type.includes('axe');
    if (activeFilter.value === 'armor') return type.startsWith('arm') || type.includes('helmet') || type.includes('chest');
    if (activeFilter.value === 'ring') return type.startsWith('ring') || type.includes('necklace');
    if (activeFilter.value === 'consumable') return type.startsWith('pot') || type.startsWith('scroll');
    return false;
  });
});
</script>

<style scoped>
.inventory-container {
  width: 100%;
}

.pixel-box-iron {
  border: 3px solid #1e293b;
  box-shadow: 
    inset 0 0 0 1px #334155,
    4px 4px 0 #0f172a;
  image-rendering: pixelated;
}

.filter-btn {
  background: #1e293b;
  border: 2px solid #334155;
  color: #64748b;
  padding: 0.4rem 1rem;
  font-size: 0.75rem;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s;
  text-transform: uppercase;
}

.filter-btn:hover {
  border-color: #4a9eff;
  color: #e2e8f0;
}

.filter-btn.active {
  background: #4a9eff;
  border-color: #e2e8f0;
  color: white;
  box-shadow: 0 0 10px rgba(74, 158, 255, 0.4);
}
</style>

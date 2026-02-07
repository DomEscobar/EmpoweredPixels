<template>
  <div class="inventory-container">
    <!-- Filtering Bar -->
    <div class="pixel-box-iron bg-slate-900/90 p-4 mb-6 flex flex-wrap items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <h2 class="text-amber-400 font-bold tracking-wider">FILTER:</h2>
        <div class="flex gap-2">
          <button 
            v-for="filter in filters" 
            :key="filter.id"
            @click="activeFilter = filter.id"
            class="filter-btn"
            :class="{ 'active': activeFilter === filter.id }"
            :data-testid="'filter-' + filter.id"
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
      <div class="animate-spin text-4xl mb-4">⚙️</div>
      <p class="text-amber-400 animate-pulse">SCANNING VAULT...</p>
    </div>

    <div v-else-if="filteredEquipment.length === 0" class="pixel-box-iron bg-slate-900/60 py-20 text-center">
      <p class="text-slate-500">NO ITEMS FOUND IN THIS SECTOR</p>
    </div>

    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
      <EquipmentCard 
        v-for="item in filteredEquipment" 
        :key="item.id" 
        :item="item"
        @enhance="$emit('enhance', item)"
        @salvage="$emit('salvage', item)"
        @toggle-favorite="$emit('toggle-favorite', item)"
        data-testid="inventory-item"
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

const filteredEquipment = computed(() => {
  if (activeFilter.value === 'all') return props.equipment;
  
  // Note: Backend 'type' currently maps to ItemID. 
  // We'll map some known prefixes or assume type strings for now.
  return props.equipment.filter(item => {
    const type = item.type.toLowerCase();
    if (activeFilter.value === 'weapon') return type.includes('wpn');
    if (activeFilter.value === 'ring') return type.includes('ring');
    if (activeFilter.value === 'consumable') return type.includes('pot') || type.includes('scroll');
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

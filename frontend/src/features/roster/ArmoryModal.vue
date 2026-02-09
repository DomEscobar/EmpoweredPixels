<template>
  <div class="fixed inset-0 z-110 flex items-center justify-center p-4">
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-black/90 pointer-events-auto" @click="$emit('close')"></div>
    
    <!-- Modal -->
    <div class="pixel-box bg-slate-900 w-full max-w-4xl max-h-[90vh] flex flex-col relative z-20 pointer-events-auto" data-testid="armory-modal">
      <!-- Header -->
      <div class="p-6 border-b-4 border-slate-800 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <div class="text-3xl">⚔️</div>
          <div>
            <h2 class="text-2xl font-bold text-amber-400 text-shadow-retro">ARMORY</h2>
            <p class="text-xs text-slate-500 uppercase tracking-[0.2em]">Select equipment for {{ fighter.name }}</p>
          </div>
        </div>
        <button @click="$emit('close')" class="rpg-btn-small">✕</button>
      </div>

      <!-- Filters & Stats -->
      <div class="p-4 bg-slate-950/50 border-b border-slate-800 flex gap-4 overflow-x-auto">
        <button 
          v-for="type in filterTypes" 
          :key="type"
          @click="currentFilter = type"
          :class="[
            'px-4 py-1.5 text-[10px] font-bold uppercase tracking-widest border transition-all whitespace-nowrap',
            currentFilter === type 
              ? 'bg-amber-600 border-amber-400 text-white' 
              : 'bg-slate-800 border-slate-700 text-slate-400 hover:border-slate-500'
          ]"
        >
          {{ type }}
        </button>
      </div>

      <!-- Inventory Grid (Scrollable) -->
      <div class="flex-1 overflow-y-auto p-6">
        <div v-if="inventory.isLoading" class="flex flex-col items-center justify-center py-20">
          <div class="text-4xl animate-spin mb-4">⚙️</div>
          <p class="text-amber-500 font-mono italic">Searching storage...</p>
        </div>

        <div v-else-if="filteredItems.length === 0" class="flex flex-col items-center justify-center py-20 text-slate-500 border-2 border-dashed border-slate-800 rounded">
          <div class="text-4xl mb-4 opacity-20">📦</div>
          <p class="font-bold">Thy inventory is void of such items.</p>
        </div>

        <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
          <div
            v-for="item in filteredItems"
            :key="item.id"
            @click="handleEquip(item.id)"
            class="group cursor-pointer"
            :data-testid="`inventory-item-${item.id}`"
          >
            <div 
              class="relative aspect-square pixel-box-sm bg-slate-950 transition-all hover:scale-105 active:scale-95 flex items-center justify-center overflow-hidden"
              :class="item.fighterId ? 'opacity-40 grayscale' : 'hover:border-amber-500'"
            >
              <!-- Rarity Glow -->
              <div class="absolute inset-0 opacity-10" :class="rarityColors[getRarityName(item.rarity)]"></div>
              
              <!-- Icon -->
              <div class="text-3xl relative z-10">{{ getWeaponIcon(item.type) }}</div>
              
              <!-- Level & Enhancement -->
              <div class="absolute top-1 left-1 bg-black/80 px-1 text-[8px] font-bold text-white border border-slate-800">
                L{{ item.level }}
              </div>
              <div v-if="item.enhancement > 0" class="absolute bottom-1 right-1 text-emerald-400 text-[10px] font-black drop-shadow-lg">
                +{{ item.enhancement }}
              </div>
              
              <!-- Binding Indicator -->
              <div v-if="item.fighterId" class="absolute inset-0 flex items-center justify-center">
                 <div class="bg-black/80 p-1 text-[8px] uppercase font-bold text-red-500 border border-red-900 border-t-2 border-b-2 rotate-[-15deg]">
                   BOUND
                 </div>
              </div>
            </div>
            
            <div class="mt-2 text-[10px] text-center">
              <div class="truncate font-bold text-slate-300">{{ item.type }}</div>
              <div :class="['font-bold uppercase tracking-tighter', rarityTextColors[getRarityName(item.rarity)]]">
                {{ getRarityName(item.rarity) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer Info -->
      <div class="p-4 bg-slate-900 border-t-4 border-slate-800 text-[10px] text-slate-500 flex justify-between">
        <p>※ Items bound to other warriors cannot be re-equipped here.</p>
        <p class="uppercase font-bold tracking-widest text-slate-400">Inventory: {{ inventory.equipment.length }} Items</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useInventoryStore } from '@/features/inventory/store';
import { useRosterStore } from '@/features/roster/store';
import type { Fighter } from '@/features/roster/api';

const props = defineProps<{
  fighter: Fighter;
}>();

const emit = defineEmits(['close']);

const inventory = useInventoryStore();
const roster = useRosterStore();

const currentFilter = ref('ALL');
const filterTypes = ['ALL', 'WEAPON', 'ARMOR', 'TRINKET'];

onMounted(() => {
  inventory.fetchInventory(1, 100);
});

const filteredItems = computed(() => {
  let items = inventory.equipment;
  
  if (currentFilter.value !== 'ALL') {
    items = items.filter(i => {
        const t = i.type.toLowerCase();
        if (currentFilter.value === 'WEAPON') return t.includes('sword') || t.includes('axe') || t.includes('bow') || t.includes('staff') || t.includes('blade');
        return true; 
    });
  }
  
  return items;
});

const handleEquip = async (itemId: string) => {
  const item = inventory.equipment.find(i => i.id === itemId);
  if (item?.fighterId) return; // Already bound

  try {
    await roster.equipItemToFighter(props.fighter.id, itemId);
    emit('close');
  } catch (e) {
    console.error("Equip failed", e);
  }
};

const getRarityName = (rarity: number) => {
  if (rarity >= 4) return 'Legendary';
  if (rarity === 3) return 'Epic';
  if (rarity === 2) return 'Rare';
  return 'Common';
};

const getWeaponIcon = (type: string) => {
  const t = type.toLowerCase();
  if (t.includes('axe')) return '🪓';
  if (t.includes('sword')) return '⚔️';
  if (t.includes('staff')) return '🪄';
  if (t.includes('bow')) return '🏹';
  if (t.includes('mask')) return '🎭';
  if (t.includes('ring')) return '💍';
  return '📦';
};

const rarityColors: Record<string, string> = {
    'Common': 'bg-slate-500',
    'Rare': 'bg-emerald-500',
    'Epic': 'bg-purple-600',
    'Legendary': 'bg-orange-500'
};

const rarityTextColors: Record<string, string> = {
    'Common': 'text-slate-500',
    'Rare': 'text-emerald-400',
    'Epic': 'text-purple-400',
    'Legendary': 'text-orange-400'
};
</script>

<style scoped>
.pixel-box {
  border: 4px solid #1e293b;
  box-shadow: 
    inset 0 0 0 2px #334155,
    4px 4px 0 #0f172a;
}

.pixel-box-sm {
  border: 2px solid #334155;
  box-shadow: 2px 2px 0 #0f172a;
}

.text-shadow-retro {
  text-shadow: 2px 2px 0 #000;
}

.rpg-btn-small {
  padding: 0.375rem 0.75rem;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  color: #e2e8f0;
  background: linear-gradient(to bottom, #475569, #334155);
  border: 2px solid #1e293b;
  box-shadow: 0 2px 0 #0f172a;
  cursor: pointer;
}
</style>

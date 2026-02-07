<template>
  <div 
    class="min-h-screen p-4 md:p-8 font-mono"
    :style="{ 
      backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')`,
      backgroundSize: '128px 128px',
      imageRendering: 'pixelated'
    }"
  >
    <!-- CRT Scanline Overlay -->
    <div class="pointer-events-none fixed inset-0 z-50 opacity-[0.03] bg-[repeating-linear-gradient(0deg,transparent,transparent_2px,rgba(0,0,0,0.3)_2px,rgba(0,0,0,0.3)_4px)]"></div>
    
    <!-- Vignette -->
    <div class="pointer-events-none fixed inset-0 z-40 bg-[radial-gradient(ellipse_at_center,transparent_0%,rgba(0,0,0,0.4)_100%)]"></div>

    <div class="relative z-10 max-w-7xl mx-auto space-y-6">
      
      <!-- Header Banner -->
      <header class="pixel-box bg-slate-900/95 p-6">
        <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div class="flex items-center gap-4">
            <img 
              :src="PIXEL_ASSETS.ICON_CHEST" 
              alt="" 
              class="w-12 h-12 pixelated"
            />
            <div>
              <h1 class="text-2xl md:text-3xl font-bold text-amber-400 text-shadow-retro tracking-wide">
                TREASURY VAULT
              </h1>
              <p class="text-slate-400 text-sm mt-1">
                Manage thy equipment and resources
              </p>
            </div>
          </div>

          <!-- Resource Counters -->
          <div class="flex flex-wrap gap-3">
            <!-- Particles -->
            <div class="pixel-box-sm bg-slate-800/90 px-4 py-2 flex items-center gap-3">
              <img :src="PIXEL_ASSETS.ICON_CRYSTAL" alt="" class="w-6 h-6 pixelated" />
              <div>
                <div class="text-[10px] uppercase text-slate-500 tracking-wider">Particles</div>
                <div class="text-lg font-bold text-emerald-400">{{ inventoryStore.particles.toLocaleString() }}</div>
              </div>
            </div>

            <!-- Tokens -->
            <div class="pixel-box-sm bg-slate-800/90 px-4 py-2 flex items-center gap-3">
              <img :src="PIXEL_ASSETS.ICON_COIN" alt="" class="w-6 h-6 pixelated" />
              <div>
                <div class="text-[10px] uppercase text-slate-500 tracking-wider">Tokens</div>
                <div class="flex gap-2 text-sm font-bold">
                  <span class="text-slate-300" title="Common">{{ inventoryStore.commonTokens }}</span>
                  <span class="text-blue-400" title="Rare">{{ inventoryStore.rareTokens }}</span>
                  <span class="text-purple-400" title="Fabled">{{ inventoryStore.fabledTokens }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- Controls Bar -->
      <div class="pixel-box bg-slate-900/90 p-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <img :src="PIXEL_ASSETS.ICON_SWORD" alt="" class="w-5 h-5 pixelated" />
          <h2 class="text-lg font-bold text-amber-300">EQUIPMENT</h2>
          <span class="pixel-badge bg-slate-700 text-slate-300 px-2 py-0.5 text-xs">
            {{ inventoryStore.totalEquipment }}
          </span>
        </div>
        <div class="flex gap-2">
          <button class="rpg-btn-small opacity-50 cursor-not-allowed" disabled>
            <img :src="PIXEL_ASSETS.ICON_SCROLL" alt="" class="w-4 h-4 pixelated mr-1" />
            FILTER
          </button>
          <button class="rpg-btn-small" @click="fetchData">
            <span class="mr-1">↻</span>
            REFRESH
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="inventoryStore.isLoading" class="pixel-box bg-slate-900/90 py-16 text-center">
        <div class="inline-block animate-spin text-4xl mb-4">⚙️</div>
        <p class="text-amber-400 text-shadow-retro animate-pulse">Accessing vault...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="inventoryStore.equipment.length === 0" class="pixel-box bg-slate-900/90 py-16 text-center">
        <img :src="PIXEL_ASSETS.ICON_CHEST_EMPTY" alt="" class="w-16 h-16 pixelated mx-auto mb-4 opacity-50" />
        <h3 class="text-lg font-bold text-slate-400">Vault is Empty</h3>
        <p class="text-slate-500 text-sm mt-2">Win battles to earn equipment crates</p>
      </div>

      <!-- Equipment Grid -->
      <InventoryGrid 
        :equipment="inventoryStore.equipment"
        :isLoading="inventoryStore.isLoading"
        @enhance="openEnhanceModal"
        @salvage="openSalvageModal"
        @toggle-favorite="inventoryStore.toggleFavorite($event.id, !$event.isFavorite)"
      />

      <!-- Enhance Modal -->
      <Teleport to="body">
        <div v-if="showEnhanceModal" class="fixed inset-0 z-100 flex items-center justify-center p-4">
          <!-- Backdrop -->
          <div class="absolute inset-0 bg-black/80" @click="closeEnhanceModal"></div>
          
          <!-- Modal -->
          <div class="pixel-box bg-slate-900 w-full max-w-md relative z-10 p-6 space-y-4">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-xl font-bold text-amber-400 text-shadow-retro flex items-center gap-2">
                <img :src="PIXEL_ASSETS.ICON_ANVIL" alt="" class="w-6 h-6 pixelated" />
                ENHANCE GEAR
              </h3>
              <button class="text-slate-500 hover:text-white text-xl" @click="closeEnhanceModal">✕</button>
            </div>

            <div v-if="selectedItem" class="space-y-4">
              <!-- Current vs Next -->
              <div class="pixel-box-sm bg-slate-800/80 p-4 flex items-center justify-between">
                <div class="text-center">
                  <div class="text-xs text-slate-500 uppercase">Current</div>
                  <div class="text-lg font-bold text-white">
                    Lvl {{ selectedItem.level }}
                    <span class="text-indigo-400 text-sm">(+{{ selectedItem.enhancement }})</span>
                  </div>
                </div>
                <div class="text-2xl text-amber-500 animate-pulse">→</div>
                <div class="text-center">
                  <div class="text-xs text-slate-500 uppercase">Next</div>
                  <div class="text-lg font-bold text-emerald-400">+{{ selectedItem.enhancement + 1 }}</div>
                </div>
              </div>

              <!-- Cost -->
              <div v-if="enhanceCost" class="space-y-2">
                <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider">Required</h4>
                <div class="grid grid-cols-2 gap-2">
                  <div class="pixel-box-sm bg-slate-800/60 p-3 flex items-center justify-between">
                    <span class="text-slate-400 text-sm">Particles</span>
                    <span class="font-bold" :class="canAffordParticles ? 'text-emerald-400' : 'text-red-400'">
                      {{ enhanceCost.particles }}
                    </span>
                  </div>
                  <div class="pixel-box-sm bg-slate-800/60 p-3 flex items-center justify-between">
                    <span class="text-slate-400 text-sm">{{ enhanceCost.tokenType }}</span>
                    <span class="font-bold" :class="canAffordTokens ? 'text-emerald-400' : 'text-red-400'">
                      {{ enhanceCost.tokens }}
                    </span>
                  </div>
                </div>
              </div>

              <!-- Error -->
              <div v-if="enhanceError" class="pixel-box-sm bg-red-900/30 border-red-500/50 p-3 text-red-300 text-sm">
                {{ enhanceError }}
              </div>

              <!-- Actions -->
              <div class="flex gap-3 pt-2">
                <button class="rpg-btn-small flex-1 bg-slate-700 hover:bg-slate-600" @click="closeEnhanceModal">
                  CANCEL
                </button>
                <button 
                  class="rpg-btn flex-1"
                  :class="{ 'opacity-50 cursor-not-allowed': !canAfford || isEnhancing }"
                  :disabled="!canAfford || isEnhancing"
                  @click="confirmEnhance"
                >
                  {{ isEnhancing ? 'FORGING...' : 'ENHANCE' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Salvage Modal -->
      <Teleport to="body">
        <div v-if="showSalvageModal" class="fixed inset-0 z-100 flex items-center justify-center p-4">
          <!-- Backdrop -->
          <div class="absolute inset-0 bg-black/80" @click="closeSalvageModal"></div>
          
          <!-- Modal -->
          <div class="pixel-box bg-slate-900 w-full max-w-md relative z-10 p-6 space-y-4">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-xl font-bold text-red-400 text-shadow-retro flex items-center gap-2">
                <img :src="PIXEL_ASSETS.ICON_SKULL" alt="" class="w-6 h-6 pixelated" />
                SALVAGE GEAR
              </h3>
              <button class="text-slate-500 hover:text-white text-xl" @click="closeSalvageModal">✕</button>
            </div>

            <div v-if="selectedItem" class="space-y-4">
              <p class="text-slate-300">
                Art thou certain? This cannot be undone.
              </p>
              
              <div class="pixel-box-sm bg-red-900/20 border-red-500/30 p-4 flex items-center gap-3">
                <span class="text-2xl">⚠️</span>
                <span class="text-red-200 text-sm">Item will be destroyed and converted to Particles.</span>
              </div>

              <!-- Actions -->
              <div class="flex gap-3 pt-2">
                <button class="rpg-btn-small flex-1 bg-slate-700 hover:bg-slate-600" @click="closeSalvageModal">
                  CANCEL
                </button>
                <button 
                  class="rpg-btn flex-1 bg-red-700 hover:bg-red-600 border-red-800"
                  :class="{ 'opacity-50 cursor-not-allowed': isSalvaging }"
                  :disabled="isSalvaging"
                  @click="confirmSalvage"
                >
                  {{ isSalvaging ? 'SALVAGING...' : 'DESTROY' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Teleport>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useInventoryStore } from '@/features/inventory/store';
import { getEnhanceCost, type Equipment } from '@/features/inventory/api';
import { useAuthStore } from '@/features/auth/store';
import InventoryGrid from '@/features/inventory/components/InventoryGrid.vue';

const PIXEL_ASSETS = {
  BG_DUNGEON: 'https://vibemedia.space/bg_dungeon_vault_8x9y0z_v1.png?prompt=dark%20dungeon%20stone%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
  ICON_CHEST: 'https://vibemedia.space/icon_chest_inv_1a2b3c_v1.png?prompt=golden%20treasure%20chest%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_CHEST_EMPTY: 'https://vibemedia.space/icon_chest_empty_4d5e6f_v1.png?prompt=empty%20wooden%20chest%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_CRYSTAL: 'https://vibemedia.space/icon_crystal_7g8h9i_v1.png?prompt=glowing%20green%20energy%20crystal%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_COIN: 'https://vibemedia.space/icon_coin_stack_0j1k2l_v1.png?prompt=stack%20of%20gold%20coins%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SWORD: 'https://vibemedia.space/icon_sword_inv_3m4n5o_v1.png?prompt=steel%20sword%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_SCROLL: 'https://vibemedia.space/icon_scroll_inv_6p7q8r_v1.png?prompt=ancient%20scroll%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_ANVIL: 'https://vibemedia.space/icon_anvil_9s0t1u_v1.png?prompt=blacksmith%20anvil%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_SKULL: 'https://vibemedia.space/icon_skull_inv_2v3w4x_v1.png?prompt=skull%20warning%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON'
};

const inventoryStore = useInventoryStore();
const authStore = useAuthStore();

const showEnhanceModal = ref(false);
const showSalvageModal = ref(false);
const selectedItem = ref<Equipment | null>(null);
const enhanceCost = ref<{ particles: number; tokens: number; tokenType: string } | null>(null);
const enhanceError = ref<string | null>(null);
const isEnhancing = ref(false);
const isSalvaging = ref(false);

const canAffordParticles = computed(() => {
   if (!enhanceCost.value) return false;
   return inventoryStore.particles >= enhanceCost.value.particles;
});

const canAffordTokens = computed(() => {
   if (!enhanceCost.value) return false;
   const type = enhanceCost.value.tokenType.toLowerCase();
   if (type === 'common') return inventoryStore.commonTokens >= enhanceCost.value.tokens;
   if (type === 'rare') return inventoryStore.rareTokens >= enhanceCost.value.tokens;
   if (type === 'fabled') return inventoryStore.fabledTokens >= enhanceCost.value.tokens;
   if (type === 'mythic') return inventoryStore.mythicTokens >= enhanceCost.value.tokens;
   return false;
});

const canAfford = computed(() => canAffordParticles.value && canAffordTokens.value);

const fetchData = async () => {
  await Promise.all([
    inventoryStore.fetchBalances(),
    inventoryStore.fetchInventory()
  ]);
};

const openEnhanceModal = async (item: Equipment) => {
  selectedItem.value = item;
  enhanceCost.value = null;
  enhanceError.value = null;
  showEnhanceModal.value = true;
  
  try {
     if (authStore.token) {
        enhanceCost.value = await getEnhanceCost(authStore.token, item.id);
     }
  } catch (e) {
     enhanceError.value = "Failed to calculate enhancement cost.";
  }
};

const closeEnhanceModal = () => {
  showEnhanceModal.value = false;
  selectedItem.value = null;
};

const confirmEnhance = async () => {
   if (!selectedItem.value) return;
   isEnhancing.value = true;
   try {
      await inventoryStore.enhanceItem(selectedItem.value.id);
      closeEnhanceModal();
   } catch (e) {
      enhanceError.value = "Enhancement failed.";
   } finally {
      isEnhancing.value = false;
   }
};

const openSalvageModal = (item: Equipment) => {
  selectedItem.value = item;
  showSalvageModal.value = true;
};

const closeSalvageModal = () => {
  showSalvageModal.value = false;
  selectedItem.value = null;
};

const confirmSalvage = async () => {
   if (!selectedItem.value) return;
   isSalvaging.value = true;
   try {
      await inventoryStore.salvageItem(selectedItem.value.id);
      closeSalvageModal();
   } catch (e) {
      console.error(e);
   } finally {
      isSalvaging.value = false;
   }
};

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.pixelated {
  image-rendering: pixelated;
}

.text-shadow-retro {
  text-shadow: 2px 2px 0 #000, 4px 4px 0 rgba(0, 0, 0, 0.3);
}

.pixel-box {
  border: 4px solid #1e293b;
  box-shadow: 
    inset 0 0 0 2px #334155,
    4px 4px 0 #0f172a;
  image-rendering: pixelated;
}

.pixel-box-sm {
  border: 2px solid #334155;
  box-shadow: 2px 2px 0 #0f172a;
}

.pixel-badge {
  border: 2px solid #475569;
}

.rpg-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.625rem 1.25rem;
  font-weight: 700;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #fef3c7;
  background: linear-gradient(to bottom, #d97706, #b45309);
  border: 3px solid #92400e;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 4px 0 #78350f,
    0 6px 4px rgba(0, 0, 0, 0.4);
  cursor: pointer;
  transition: all 0.1s;
}

.rpg-btn:hover:not(:disabled) {
  background: linear-gradient(to bottom, #f59e0b, #d97706);
  transform: translateY(-1px);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.3),
    0 5px 0 #78350f,
    0 7px 6px rgba(0, 0, 0, 0.4);
}

.rpg-btn:active:not(:disabled) {
  transform: translateY(2px);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 2px 0 #78350f,
    0 3px 2px rgba(0, 0, 0, 0.3);
}

.rpg-btn-small {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.375rem 0.75rem;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #e2e8f0;
  background: linear-gradient(to bottom, #475569, #334155);
  border: 2px solid #1e293b;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 2px 0 #0f172a;
  cursor: pointer;
  transition: all 0.1s;
}

.rpg-btn-small:hover:not(:disabled) {
  background: linear-gradient(to bottom, #64748b, #475569);
}

.rpg-btn-small:active:not(:disabled) {
  transform: translateY(1px);
  box-shadow: 0 1px 0 #0f172a;
}
</style>

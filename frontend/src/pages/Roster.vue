<template>
  <div 
    class="min-h-screen p-4 md:p-8 font-mono"
    :style="{ 
      backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')`,
      backgroundSize: '128px 128px',
      imageRendering: 'pixelated'
    }"
    data-testid="roster-page"
  >
    <!-- CRT Scanline Overlay -->
    <div class="pointer-events-none fixed inset-0 z-50 opacity-[0.03] bg-[repeating-linear-gradient(0deg,transparent,transparent_2px,rgba(0,0,0,0.3)_2px,rgba(0,0,0,0.3)_4px)]"></div>
    
    <!-- Vignette -->
    <div class="pointer-events-none fixed inset-0 z-40 bg-[radial-gradient(ellipse_at_center,transparent_0%,rgba(0,0,0,0.4)_100%)]"></div>

    <div class="relative z-10 max-w-7xl mx-auto space-y-6">
      <!-- Header Banner -->
      <header class="pixel-box bg-slate-900/95 p-6" data-testid="roster-header">
        <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div class="flex items-center gap-4">
            <img 
              :src="PIXEL_ASSETS.ICON_BARRACKS" 
              alt="" 
              class="w-12 h-12 pixelated"
            />
            <div>
              <h1 class="text-2xl md:text-3xl font-bold text-amber-400 text-shadow-retro tracking-wide">
                WAR ROOM
              </h1>
              <p class="text-slate-400 text-sm mt-1">
                Train, equip, and prepare thy champions
              </p>
            </div>
          </div>

          <div class="flex items-center gap-4">
            <!-- Fighter Count -->
            <div class="pixel-box-sm bg-slate-800/80 px-4 py-2 text-center">
              <div class="text-xs text-slate-500 uppercase tracking-wider">
                Warriors
              </div>
              <div class="text-2xl font-bold text-white">
                {{ roster.fighters.length }}
              </div>
            </div>

            <!-- Recruit Button -->
            <button class="rpg-btn flex items-center gap-2" data-testid="recruit-button" @click="openCreateWizard">
              <img :src="PIXEL_ASSETS.ICON_PLUS" alt="" class="w-4 h-4 pixelated" />
              RECRUIT
            </button>
          </div>
        </div>
      </header>

      <!-- Error Alert -->
      <div v-if="roster.error" class="pixel-box bg-red-900/90 border-red-500 p-4 flex items-center gap-3" data-testid="roster-error">
        <span class="text-red-300 text-xl">⚠️</span>
        <p class="text-red-100">
          {{ roster.error }}
        </p>
        <button class="ml-auto text-red-300 hover:text-red-100" @click="roster.error = null">
          ✕
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="roster.isLoading && !roster.fighters.length" class="pixel-box bg-slate-900/90 py-20 text-center">
        <div class="inline-block animate-spin text-4xl mb-4">
          ⚔️
        </div>
        <p class="text-amber-400 text-shadow-retro animate-pulse">
          Summoning warriors...
        </p>
      </div>

      <!-- Empty State -->
      <div v-else-if="!roster.fighters.length" class="pixel-box bg-slate-900/90 py-16 text-center" data-testid="roster-empty">
        <img :src="PIXEL_ASSETS.ICON_HELMET_EMPTY" alt="" class="w-20 h-20 pixelated mx-auto mb-4 opacity-50" />
        <h2 class="text-xl font-bold text-slate-400 mb-2">
          No Warriors Yet
        </h2>
        <p class="text-slate-500 text-sm mb-6">
          Thy arena awaits its first champion
        </p>
        <button class="rpg-btn" data-testid="recruit-first-button" @click="openCreateWizard">
          RECRUIT FIRST WARRIOR
        </button>
      </div>

      <!-- Fighter Grid -->
      <div v-else class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3" data-testid="fighter-grid">
        <article
          v-for="fighter in roster.fighters"
          :key="fighter.id"
          class="pixel-box bg-slate-900/90 overflow-hidden group hover:border-indigo-500/50 transition-colors"
          :data-testid="`fighter-card-${fighter.id}`"
        >
          <!-- Attunement Glow Header -->
          <div 
            class="h-2 w-full bg-slate-700"
          ></div>

          <div class="p-4 space-y-4">
            <!-- Top Row: Avatar + Name + Level -->
            <div class="flex items-start gap-4">
              <!-- Fighter Avatar -->
              <div class="relative">
                <div
                  class="pixel-box-sm flex h-16 w-16 items-center justify-center text-3xl overflow-hidden bg-slate-800"
                >
                  <VoxelFighter :seed="fighter.id" :animate="true" />
                </div>
                <!-- Level Badge -->
                <div class="absolute -bottom-1 -right-1 bg-amber-600 text-white text-xs font-bold px-1.5 py-0.5 border-2 border-slate-900">
                  {{ fighter.level }}
                </div>
              </div>

              <!-- Name + Status -->
              <div class="flex-1 min-w-0">
                <h3 class="truncate text-lg font-bold text-white">
                  {{ fighter.name }}
                </h3>
                <div class="mt-1 flex items-center gap-2">
                  <span class="pixel-badge bg-emerald-900/50 border-emerald-500/50 text-emerald-400 text-xs px-2 py-0.5 flex items-center gap-1">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
                    Ready
                  </span>
                </div>
                
                <!-- XP Bar -->
                <div class="mt-2">
                  <div class="mb-0.5 flex justify-between text-[10px]">
                    <span class="text-slate-500 uppercase">EXP</span>
                    <span class="text-slate-400">{{ fighter.currentExp !== undefined ? fighter.currentExp : (fighter.xp || 0) }}/{{ fighter.levelExp || fighter.xpToNextLevel }}</span>
                  </div>
                  <div class="h-1.5 pixel-box-sm bg-slate-800/80 overflow-hidden">
                    <div
                      class="h-full bg-gradient-to-r from-indigo-600 to-purple-500 transition-all"
                      :style="{ width: `${getExpPercent(fighter)}%` }"
                    ></div>
                  </div>
                </div>

                <!-- Match Stats -->
                <div class="mt-2 flex items-center gap-3">
                  <div class="flex items-center gap-1 text-[10px]">
                    <span class="text-emerald-400 font-bold">{{ fighter.matchesWon || 0 }}</span>
                    <span class="text-slate-500">W</span>
                  </div>
                  <div class="flex items-center gap-1 text-[10px]">
                    <span class="text-red-400 font-bold">{{ fighter.matchesLost || 0 }}</span>
                    <span class="text-slate-500">L</span>
                  </div>
                  <div class="flex items-center gap-1 text-[10px] ml-auto">
                    <span class="text-amber-400 font-bold">{{ fighter.totalMatches || 0 }}</span>
                    <span class="text-slate-500">Total</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Core Stats Grid -->
            <div class="grid grid-cols-4 gap-1">
              <div
                v-for="stat in getCoreStats(fighter)"
                :key="stat.key"
                class="pixel-box-sm bg-slate-800/60 p-1.5 text-center"
              >
                <div class="text-[10px] text-slate-500 uppercase">
                  {{ stat.label }}
                </div>
                <div :class="['text-sm font-bold', stat.color]">
                  {{ stat.value }}
                </div>
              </div>
            </div>

            <!-- Stat Bars -->
            <div class="space-y-1.5">
              <div v-for="bar in getStatBars(fighter)" :key="bar.key">
                <div class="mb-0.5 flex justify-between text-[10px]">
                  <span class="text-slate-500 uppercase">{{ bar.label }}</span>
                  <span class="text-slate-400">{{ bar.value }}</span>
                </div>
                <div class="h-1 pixel-box-sm bg-slate-800/80 overflow-hidden">
                  <div
                    :class="['h-full transition-all', bar.colorClass]"
                    :style="{ width: `${bar.percent}%` }"
                  ></div>
                </div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-2 pt-2 border-t-2 border-slate-800">
              <button
                class="rpg-btn flex-1 text-sm py-2"
                :data-testid="`manage-fighter-${fighter.id}`"
                @click="openFighterPanel(fighter)"
              >
                MANAGE
              </button>
              <button
                class="rpg-btn-small bg-red-900/50 hover:bg-red-800/50 border-red-700 text-red-300 px-3"
                :data-testid="`dismiss-fighter-${fighter.id}`"
                @click="confirmDismiss(fighter)"
              >
                <img :src="PIXEL_ASSETS.ICON_SKULL" alt="" class="w-4 h-4 pixelated" />
              </button>
            </div>
          </div>
        </article>
      </div>

      <!-- Fighter Detail Panel (Slide-out) -->
      <Teleport to="body">
        <Transition
          enter-active-class="transition duration-300"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="transition duration-200"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="selectedFighter" class="fixed inset-0 z-100">
            <!-- Backdrop -->
            <div class="absolute inset-0 bg-black/80" @click="selectedFighter = null"></div>
            
            <!-- Panel -->
            <Transition
              enter-active-class="transition duration-300 ease-out"
              enter-from-class="translate-x-full"
              enter-to-class="translate-x-0"
              leave-active-class="transition duration-200 ease-in"
              leave-from-class="translate-x-0"
              leave-to-class="translate-x-full"
            >
              <aside
                v-if="selectedFighter"
                class="absolute right-0 top-0 h-full w-full max-w-xl overflow-y-auto pixel-box bg-slate-900 border-l-4 border-slate-700"
                :style="{ 
                  backgroundImage: `url('${PIXEL_ASSETS.BG_DUNGEON}')`,
                  backgroundSize: '64px 64px',
                  imageRendering: 'pixelated'
                }"
              >
                <!-- Panel Header -->
                <div class="sticky top-0 z-10 pixel-box bg-slate-900/95 border-t-0 border-x-0 p-4">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-3">
                      <div
                      class="pixel-box-sm flex h-12 w-12 items-center justify-center overflow-hidden bg-slate-800"
                    >
                      <VoxelFighter :seed="selectedFighter.id" :animate="true" />
                    </div>
                      <div>
                        <h2 class="text-xl font-bold text-amber-400 text-shadow-retro">
                          {{ selectedFighter.name }}
                        </h2>
                        <p class="text-sm text-slate-400">
                          Level {{ selectedFighter.level }} Warrior
                        </p>
                      </div>
                    </div>
                    <button
                      class="rpg-btn-small text-lg px-2 py-1"
                      @click="selectedFighter = null"
                    >
                      ✕
                    </button>
                  </div>
                </div>

                <!-- Panel Content -->
                <div class="p-4 space-y-6">
                  <FighterStats
                    :fighter="selectedFighter"
                    :equipment="roster.equipment[selectedFighter.id] || []"
                    :data-testid="`fighter-stats-${selectedFighter.id}`"
                    @open-armory="showArmory = true"
                  />
                </div>
              </aside>
            </Transition>
          </div>
        </Transition>
      </Teleport>

      <!-- Armory Modal -->
      <Teleport to="body">
        <Transition
          enter-active-class="transition duration-200"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="transition duration-150"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <ArmoryModal
            v-if="showArmory && selectedFighter"
            :fighter="selectedFighter"
            @close="showArmory = false"
          />
        </Transition>
      </Teleport>

      <!-- Create Fighter Wizard -->
      <Teleport to="body">
        <Transition
          enter-active-class="transition duration-200"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="transition duration-150"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="showCreate" class="fixed inset-0 z-100 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/80" @click="closeCreateWizard"></div>
            
            <Transition
              enter-active-class="transition duration-300 ease-out"
              enter-from-class="opacity-0 scale-95"
              enter-to-class="opacity-100 scale-100"
              leave-active-class="transition duration-200 ease-in"
              leave-from-class="opacity-100 scale-100"
              leave-to-class="opacity-0 scale-95"
            >
              <div v-if="showCreate" class="pixel-box bg-slate-900 w-full max-w-lg relative z-10" data-testid="create-wizard">
                <!-- Wizard Header -->
                <div class="p-4 border-b-4 border-slate-800 flex items-center gap-3">
                  <img :src="PIXEL_ASSETS.ICON_SCROLL" alt="" class="w-6 h-6 pixelated" />
                  <div>
                    <h2 class="text-xl font-bold text-amber-400 text-shadow-retro">
                      RECRUIT WARRIOR
                    </h2>
                    <p class="text-xs text-slate-500">
                      Create a new champion
                    </p>
                  </div>
                </div>

                <!-- Wizard Content -->
                <form class="p-6 space-y-6" @submit.prevent="handleCreate">
                  <!-- Fighter Preview -->
                  <div class="flex justify-center">
                    <div class="relative">
                      <div
                        class="pixel-box flex h-24 w-24 items-center justify-center transition-all overflow-hidden bg-slate-800"
                      >
                        <VoxelFighter :seed="previewSeed" :animate="true" />
                      </div>
                      <div class="absolute -bottom-3 left-1/2 -translate-x-1/2 whitespace-nowrap pixel-box-sm bg-slate-800 px-3 py-1 text-sm font-bold text-white">
                        {{ newName || 'Unnamed' }}
                      </div>
                    </div>
                  </div>

                  <!-- Name Input -->
                  <div class="pt-4">
                    <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Fighter Name</label>
                    <input
                      v-model="newName"
                      type="text"
                      placeholder="Enter a legendary name..."
                      required
                      minlength="2"
                      maxlength="24"
                      class="w-full pixel-box-sm bg-slate-800 px-4 py-3 text-lg text-white placeholder-slate-600 focus:border-amber-500 focus:outline-none"
                      data-testid="new-fighter-name-input"
                    />
                    <p class="mt-1 text-[10px] text-slate-600">
                      {{ newName.length }}/24 characters
                    </p>
                  </div>

                  <!-- Info Box -->
                  <div class="pixel-box-sm bg-slate-800/50 p-3 flex gap-3">
                    <span class="text-indigo-400">ℹ️</span>
                    <p class="text-xs text-slate-400">
                      Thy fighter shall begin at Level 1. Win battles to earn XP!
                    </p>
                  </div>

                  <!-- Actions -->
                  <div class="flex gap-3 pt-2">
                    <button
                      type="button"
                      class="rpg-btn-small flex-1"
                      data-testid="cancel-recruit"
                      @click="closeCreateWizard"
                    >
                      CANCEL
                    </button>
                    <button
                      type="submit"
                      :disabled="roster.isLoading || !newName.trim()"
                      :class="['rpg-btn flex-1', { 'opacity-50 cursor-not-allowed': roster.isLoading || !newName.trim() }]"
                      data-testid="confirm-recruit"
                    >
                      <span v-if="roster.isLoading" class="flex items-center justify-center gap-2">
                        <span class="animate-spin">⚙️</span>
                        RECRUITING...
                      </span>
                      <span v-else>RECRUIT</span>
                    </button>
                  </div>
                </form>
              </div>
            </Transition>
          </div>
        </Transition>
      </Teleport>

      <!-- Dismiss Confirmation -->
      <Teleport to="body">
        <Transition
          enter-active-class="transition duration-200"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="transition duration-150"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="dismissTarget" class="fixed inset-0 z-100 flex items-center justify-center p-4">
            <div class="absolute inset-0 bg-black/80" @click="dismissTarget = null"></div>
            
            <div class="pixel-box bg-slate-900 w-full max-w-sm relative z-10 p-6" data-testid="dismiss-confirmation">
              <div class="flex items-center gap-3 mb-4">
                <img :src="PIXEL_ASSETS.ICON_SKULL" alt="" class="w-8 h-8 pixelated" />
                <h3 class="text-lg font-bold text-red-400 text-shadow-retro">
                  DISMISS WARRIOR
                </h3>
              </div>
              
              <p class="text-slate-300 mb-2">
                Art thou certain about dismissing <span class="text-white font-bold">{{ dismissTarget.name }}</span>?
              </p>
              <div class="pixel-box-sm bg-red-900/20 border-red-500/30 p-3 mb-6 flex gap-2">
                <span>⚠️</span>
                <p class="text-xs text-red-200">
                  This cannot be undone. All progress shall be lost.
                </p>
              </div>
              
              <div class="flex gap-3">
                <button
                  class="rpg-btn-small flex-1"
                  data-testid="cancel-dismiss"
                  @click="dismissTarget = null"
                >
                  CANCEL
                </button>
                <button
                  class="rpg-btn flex-1 bg-red-700 hover:bg-red-600 border-red-800"
                  data-testid="confirm-dismiss"
                  @click="handleDismiss"
                >
                  DISMISS
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </Teleport>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRosterStore } from '@/features/roster/store';
import { useAuthStore } from '@/features/auth/store';
import type { Fighter } from '@/features/roster/api';
import FighterStats from '@/features/roster/FighterStats.vue';
import ArmoryModal from '@/features/roster/ArmoryModal.vue';
import VoxelFighter from '@/shared/ui/VoxelFighter.vue';

const auth = useAuthStore();

const PIXEL_ASSETS = {
  BG_DUNGEON: 'https://vibemedia.space/bg_dungeon_roster_7x8y9z_v1.png?prompt=dark%20dungeon%20stone%20floor%20tile%20texture%20seamless&style=pixel_game_asset&key=NOGON',
  ICON_BARRACKS: 'https://vibemedia.space/icon_barracks_1a2b3c_v1.png?prompt=medieval%20barracks%20building%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_PLUS: 'https://vibemedia.space/icon_plus_gold_4d5e6f_v1.png?prompt=golden%20plus%20sign%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_HELMET_EMPTY: 'https://vibemedia.space/icon_helmet_empty_7g8h9i_v1.png?prompt=empty%20knight%20helmet%20stand%20pixel%20art&style=pixel_game_asset&key=NOGON',
  ICON_SKULL: 'https://vibemedia.space/icon_skull_red_0j1k2l_v1.png?prompt=red%20skull%20warning%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
  ICON_SCROLL: 'https://vibemedia.space/icon_scroll_recruit_3m4n5o_v1.png?prompt=recruitment%20scroll%20pixel%20art%20icon&style=pixel_game_asset&key=NOGON',
};

const roster = useRosterStore();
const showCreate = ref(false);
const showArmory = ref(false);
const newName = ref('');
const selectedFighter = ref<Fighter | null>(null);
const dismissTarget = ref<Fighter | null>(null);

// Preview seed for new fighter
const previewSeed = computed(() => newName.value || 'preview');

const getExpPercent = (fighter: Fighter) => {
  const next = fighter.levelExp || fighter.xpToNextLevel;
  if (!next || next === 0) return 0;
  const current = fighter.currentExp !== undefined ? fighter.currentExp : (fighter.xp || 0);
  return Math.min(100, (current / next) * 100);
};

const getCoreStats = (fighter: Fighter) => [
  { key: 'power', label: 'PWR', value: fighter.power, color: 'text-red-400' },
  { key: 'vitality', label: 'VIT', value: fighter.vitality, color: 'text-emerald-400' },
  { key: 'armor', label: 'ARM', value: fighter.armor, color: 'text-amber-400' },
  { key: 'speed', label: 'SPD', value: fighter.speed, color: 'text-sky-400' },
];

const getStatBars = (fighter: Fighter) => {
  const maxStat = 100;
  return [
    { key: 'accuracy', label: 'ACC', value: fighter.accuracy, percent: (fighter.accuracy / maxStat) * 100, colorClass: 'bg-indigo-500' },
    { key: 'agility', label: 'AGI', value: fighter.agility, percent: (fighter.agility / maxStat) * 100, colorClass: 'bg-purple-500' },
    { key: 'precision', label: 'PRE', value: fighter.precision, percent: (fighter.precision / maxStat) * 100, colorClass: 'bg-pink-500' },
  ];
};

const openCreateWizard = () => {
  newName.value = '';
  showCreate.value = true;
};

const closeCreateWizard = () => {
  showCreate.value = false;
};

const handleCreate = async () => {
  if (!newName.value.trim()) return;
  roster.error = null; // Clear previous errors
  try {
    await roster.addFighter(newName.value);
    closeCreateWizard();
  } catch (e) {
    // Error is displayed via roster.error binding
  }
};

const openFighterPanel = async (fighter: Fighter) => {
  selectedFighter.value = fighter;
  await roster.fetchFighterEquipment(fighter.id);
};

const confirmDismiss = (fighter: Fighter) => {
  dismissTarget.value = fighter;
};

const handleDismiss = async () => {
  if (!dismissTarget.value) return;
  await roster.removeFighter(dismissTarget.value.id);
  dismissTarget.value = null;
};

onMounted(() => {
  roster.fetchFighters();
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
  border: 2px solid currentColor;
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

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
              <div class="text-xs text-slate-500 uppercase tracking-wider">Warriors</div>
              <div class="text-2xl font-bold text-white">{{ roster.fighters.length }}</div>
            </div>

            <!-- Recruit Button -->
            <button @click="openCreateWizard" class="rpg-btn flex items-center gap-2">
              <img :src="PIXEL_ASSETS.ICON_PLUS" alt="" class="w-4 h-4 pixelated" />
              RECRUIT
            </button>
          </div>
        </div>
      </header>

      <!-- Loading State -->
      <div v-if="roster.isLoading && !roster.fighters.length" class="pixel-box bg-slate-900/90 py-20 text-center">
        <div class="inline-block animate-spin text-4xl mb-4">⚔️</div>
        <p class="text-amber-400 text-shadow-retro animate-pulse">Summoning warriors...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="!roster.fighters.length" class="pixel-box bg-slate-900/90 py-16 text-center">
        <img :src="PIXEL_ASSETS.ICON_HELMET_EMPTY" alt="" class="w-20 h-20 pixelated mx-auto mb-4 opacity-50" />
        <h2 class="text-xl font-bold text-slate-400 mb-2">No Warriors Yet</h2>
        <p class="text-slate-500 text-sm mb-6">
          Thy arena awaits its first champion
        </p>
        <button @click="openCreateWizard" class="rpg-btn">
          RECRUIT FIRST WARRIOR
        </button>
      </div>

      <!-- Fighter Grid -->
      <div v-else class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        <article
          v-for="fighter in roster.fighters"
          :key="fighter.id"
          class="pixel-box bg-slate-900/90 overflow-hidden group hover:border-indigo-500/50 transition-colors"
        >
          <!-- Attunement Glow Header -->
          <div 
            class="h-2 w-full"
            :class="fighter.attunementId ? getAttunementBarColor(fighter.attunementId) : 'bg-slate-700'"
          ></div>

          <div class="p-4 space-y-4">
            <!-- Top Row: Avatar + Name + Level -->
            <div class="flex items-start gap-4">
              <!-- Fighter Avatar -->
              <div class="relative">
                <div
                  :class="[
                    'pixel-box-sm flex h-16 w-16 items-center justify-center text-3xl overflow-hidden',
                    fighter.attunementId ? getAttunementBg(fighter.attunementId) : 'bg-slate-800'
                  ]"
                >
                  <VoxelFighter :seed="fighter.id" :attunement="fighter.attunementId" :animate="true" />
                </div>
                <!-- Level Badge -->
                <div class="absolute -bottom-1 -right-1 bg-amber-600 text-white text-xs font-bold px-1.5 py-0.5 border-2 border-slate-900">
                  {{ fighter.level }}
                </div>
              </div>

              <!-- Name + Status -->
              <div class="flex-1 min-w-0">
                <h3 class="truncate text-lg font-bold text-white">{{ fighter.name }}</h3>
                <div class="mt-1 flex items-center gap-2">
                  <span class="pixel-badge bg-emerald-900/50 border-emerald-500/50 text-emerald-400 text-xs px-2 py-0.5 flex items-center gap-1">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
                    Ready
                  </span>
                  <span v-if="fighter.attunementId" class="text-xs text-slate-500">
                    {{ fighter.attunementId }}
                  </span>
                </div>
                
                <!-- XP Bar -->
                <div class="mt-2">
                  <div class="mb-0.5 flex justify-between text-[10px]">
                    <span class="text-slate-500 uppercase">EXP</span>
                    <span class="text-slate-400">{{ fighter.currentExp }}/{{ fighter.levelExp }}</span>
                  </div>
                  <div class="h-1.5 pixel-box-sm bg-slate-800/80 overflow-hidden">
                    <div
                      class="h-full bg-linear-to-r from-indigo-600 to-purple-500 transition-all"
                      :style="{ width: `${getExpPercent(fighter)}%` }"
                    ></div>
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
                <div class="text-[10px] text-slate-500 uppercase">{{ stat.label }}</div>
                <div :class="['text-sm font-bold', stat.color]">{{ stat.value }}</div>
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

            <!-- Attunement Quick Select -->
            <div>
              <div class="mb-1.5 text-[10px] font-bold uppercase tracking-wider text-slate-500">Attunement</div>
              <div class="flex gap-1">
                <button
                  v-for="att in attunements"
                  :key="att.id"
                  @click.stop="quickSetAttunement(fighter.id, att.id)"
                  :class="[
                    'pixel-box-sm flex h-8 w-8 items-center justify-center text-base transition-all',
                    fighter.attunementId === att.id
                      ? `${att.bgActive} border-2 ${att.borderActive}`
                      : 'bg-slate-800/50 opacity-40 hover:opacity-100 hover:bg-slate-700'
                  ]"
                  :title="att.name"
                >
                  {{ att.icon }}
                </button>
                <button
                  v-if="fighter.attunementId"
                  @click.stop="quickSetAttunement(fighter.id, null)"
                  class="ml-auto rpg-btn-small text-[10px] px-2"
                >
                  CLEAR
                </button>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-2 pt-2 border-t-2 border-slate-800">
              <button
                @click="openFighterPanel(fighter)"
                class="rpg-btn flex-1 text-sm py-2"
              >
                MANAGE
              </button>
              <button
                @click="confirmDismiss(fighter)"
                class="rpg-btn-small bg-red-900/50 hover:bg-red-800/50 border-red-700 text-red-300 px-3"
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
                        :class="[
                          'pixel-box-sm flex h-12 w-12 items-center justify-center overflow-hidden',
                          selectedFighter.attunementId ? getAttunementBg(selectedFighter.attunementId) : 'bg-slate-800'
                        ]"
                      >
                        <VoxelFighter :seed="selectedFighter.id" :attunement="selectedFighter.attunementId" :animate="true" />
                      </div>
                      <div>
                        <h2 class="text-xl font-bold text-amber-400 text-shadow-retro">{{ selectedFighter.name }}</h2>
                        <p class="text-sm text-slate-400">Level {{ selectedFighter.level }} Warrior</p>
                      </div>
                    </div>
                    <button
                      @click="selectedFighter = null"
                      class="rpg-btn-small text-lg px-2 py-1"
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
                  />
                </div>
              </aside>
            </Transition>
          </div>
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
              <div v-if="showCreate" class="pixel-box bg-slate-900 w-full max-w-lg relative z-10">
                <!-- Wizard Header -->
                <div class="p-4 border-b-4 border-slate-800 flex items-center gap-3">
                  <img :src="PIXEL_ASSETS.ICON_SCROLL" alt="" class="w-6 h-6 pixelated" />
                  <div>
                    <h2 class="text-xl font-bold text-amber-400 text-shadow-retro">RECRUIT WARRIOR</h2>
                    <p class="text-xs text-slate-500">Create a new champion</p>
                  </div>
                </div>

                <!-- Wizard Content -->
                <form @submit.prevent="handleCreate" class="p-6 space-y-6">
                  <!-- Fighter Preview -->
                  <div class="flex justify-center">
                    <div class="relative">
                      <div
                        :class="[
                          'pixel-box flex h-24 w-24 items-center justify-center transition-all overflow-hidden',
                          createAttunement ? getAttunementBg(createAttunement) : 'bg-slate-800'
                        ]"
                      >
                        <VoxelFighter :seed="previewSeed" :attunement="createAttunement" :animate="true" />
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
                    />
                    <p class="mt-1 text-[10px] text-slate-600">{{ newName.length }}/24 characters</p>
                  </div>

                  <!-- Attunement Selection -->
                  <div>
                    <label class="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Starting Attunement (Optional)</label>
                    <div class="grid grid-cols-5 gap-2">
                      <button
                        v-for="att in attunements"
                        :key="att.id"
                        type="button"
                        @click="createAttunement = createAttunement === att.id ? null : att.id"
                        :class="[
                          'pixel-box-sm flex flex-col items-center gap-1 p-3 transition-all',
                          createAttunement === att.id
                            ? `${att.bgActive} ${att.borderActive}`
                            : 'bg-slate-800/60 hover:bg-slate-700'
                        ]"
                      >
                        <span class="text-xl">{{ att.icon }}</span>
                        <span class="text-[10px] font-bold" :class="createAttunement === att.id ? 'text-white' : 'text-slate-500'">
                          {{ att.name }}
                        </span>
                      </button>
                    </div>
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
                      @click="closeCreateWizard"
                      class="rpg-btn-small flex-1"
                    >
                      CANCEL
                    </button>
                    <button
                      type="submit"
                      :disabled="roster.isLoading || !newName.trim()"
                      :class="['rpg-btn flex-1', { 'opacity-50 cursor-not-allowed': roster.isLoading || !newName.trim() }]"
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
            
            <div class="pixel-box bg-slate-900 w-full max-w-sm relative z-10 p-6">
              <div class="flex items-center gap-3 mb-4">
                <img :src="PIXEL_ASSETS.ICON_SKULL" alt="" class="w-8 h-8 pixelated" />
                <h3 class="text-lg font-bold text-red-400 text-shadow-retro">DISMISS WARRIOR</h3>
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
                  @click="dismissTarget = null"
                  class="rpg-btn-small flex-1"
                >
                  CANCEL
                </button>
                <button
                  @click="handleDismiss"
                  class="rpg-btn flex-1 bg-red-700 hover:bg-red-600 border-red-800"
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
import type { Fighter } from '@/features/roster/api';
import FighterStats from '@/features/roster/FighterStats.vue';
import VoxelFighter from '@/shared/ui/VoxelFighter.vue';

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
const newName = ref('');
const createAttunement = ref<string | null>(null);
const selectedFighter = ref<Fighter | null>(null);
const dismissTarget = ref<Fighter | null>(null);

// Preview seed for new fighter
const previewSeed = computed(() => newName.value || 'preview');

const attunements = [
  { id: 'Fire', name: 'Fire', icon: '🔥', bgActive: 'bg-orange-900/60', borderActive: 'border-orange-500' },
  { id: 'Water', name: 'Water', icon: '💧', bgActive: 'bg-blue-900/60', borderActive: 'border-blue-500' },
  { id: 'Earth', name: 'Earth', icon: '🪨', bgActive: 'bg-amber-900/60', borderActive: 'border-amber-600' },
  { id: 'Wind', name: 'Wind', icon: '💨', bgActive: 'bg-teal-900/60', borderActive: 'border-teal-500' },
  { id: 'Lightning', name: 'Lightning', icon: '⚡', bgActive: 'bg-yellow-900/60', borderActive: 'border-yellow-500' },
];

const getAttunementBg = (id: string | null | undefined) => {
  if (!id) return 'bg-slate-800';
  return attunements.find(a => a.id === id)?.bgActive || 'bg-slate-800';
};

const getAttunementBarColor = (id: string) => {
  const colors: Record<string, string> = {
    Fire: 'bg-linear-to-r from-orange-600 to-red-600',
    Water: 'bg-linear-to-r from-blue-500 to-cyan-500',
    Earth: 'bg-linear-to-r from-amber-600 to-stone-600',
    Wind: 'bg-linear-to-r from-teal-500 to-emerald-500',
    Lightning: 'bg-linear-to-r from-yellow-400 to-amber-500',
  };
  return colors[id] || 'bg-slate-700';
};

const getExpPercent = (fighter: Fighter) => {
  if (!fighter.levelExp) return 0;
  return Math.min(100, (fighter.currentExp / fighter.levelExp) * 100);
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
  createAttunement.value = null;
  showCreate.value = true;
};

const closeCreateWizard = () => {
  showCreate.value = false;
};

const handleCreate = async () => {
  if (!newName.value.trim()) return;
  
  try {
    const fighter = await roster.addFighter(newName.value);
    if (createAttunement.value && fighter) {
      await roster.updateAttunement(fighter.id, createAttunement.value);
    }
    closeCreateWizard();
  } catch (e) {
    // Error handled in store
  }
};

const openFighterPanel = async (fighter: Fighter) => {
  selectedFighter.value = fighter;
  await roster.fetchFighterEquipment(fighter.id);
};

const quickSetAttunement = async (fighterId: string, attunementId: string | null) => {
  await roster.updateAttunement(fighterId, attunementId);
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

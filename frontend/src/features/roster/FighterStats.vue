<template>
  <div class="overhaul-theme-container space-y-8">
    <!-- EXPERIENCE PANORAMA (GW2 Influence) -->
    <section class="ep-card-iron p-6 overhaul-experience-panel relative overflow-hidden">
      <!-- Artistic Background Splatter -->
      <div class="absolute -right-10 -top-10 artistic-splatter ep-splatter-ink bg-purple-500/20 w-48 h-48 pointer-events-none"></div>

      <div class="relative z-10 flex flex-col sm:flex-row items-center gap-6">
        <div class="relative">
          <div class="flex h-16 w-16 items-center justify-center rounded-lg border-2 border-amber-500 bg-slate-900 text-2xl font-black text-amber-500 shadow-[0_0_15px_rgba(245,158,11,0.3)]">
            {{ fighter.level }}
          </div>
          <div class="absolute -bottom-2 -right-2 h-6 w-6 bg-amber-500 rounded-full flex items-center justify-center text-[10px] text-black font-bold">
            LV
          </div>
        </div>
        
        <div class="flex-1 w-full">
          <div class="flex justify-between items-end mb-2">
            <h3 class="ep-header-gold text-2xl">
              Commander Level
            </h3>
            <span class="text-xs text-slate-500">{{ getExpDisplay }} XP</span>
          </div>
          <div class="h-4 overflow-hidden rounded-sm bg-black border border-slate-800">
            <div
              class="h-full bg-gradient-to-r from-amber-600 via-yellow-500 to-amber-600 transition-all duration-1000 ease-out shadow-[0_0_10px_rgba(245,158,11,0.5)]"
              :style="{ width: `${getExpPercent}%` }"
            ></div>
          </div>
        </div>
      </div>
    </section>

    <!-- CORE ATTRIBUTES GRID (WoW influence - high readability) -->
    <section class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div v-for="(group, name) in statGroups" :key="name" class="ep-card-iron overhaul-stat-group">
        <div class="px-4 py-3 border-b border-slate-800 bg-black/40">
          <h4 class="text-xs font-bold uppercase tracking-[0.2em] text-slate-400">
            {{ name }}
          </h4>
        </div>
        <div class="p-4 space-y-3">
          <div v-for="stat in group" :key="stat.key" class="flex items-center justify-between group">
            <span class="text-sm text-slate-400 group-hover:text-amber-200 transition-colors">{{ stat.label }}</span>
            <div class="flex items-center gap-2">
              <span class="text-base font-mono font-bold text-white">{{ stat.value ?? 0 }}{{ stat.suffix || '' }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- LEAGUE ENROLLMENT SECTION (New Integration) -->
    <section class="ep-card-iron overhaul-stat-group overflow-hidden">
      <div class="px-4 py-3 border-b border-slate-800 bg-black/40 flex justify-between items-center">
        <h4 class="text-xs font-bold uppercase tracking-[0.2em] text-slate-400">
          Active League Enrollments
        </h4>
        <router-link to="/leagues" class="text-[10px] text-amber-500 hover:text-amber-400 font-bold uppercase">
          View All Leagues
        </router-link>
      </div>
      <div class="p-4">
        <div v-if="enrolledLeagues.length > 0" class="flex flex-wrap gap-3">
          <div
            v-for="league in enrolledLeagues" :key="league.id" 
            class="px-3 py-2 bg-slate-900 border border-amber-900/50 rounded flex items-center gap-2 group cursor-pointer hover:border-amber-500 transition-colors"
            @click="router.push(`/leagues?id=${league.id}`)"
          >
            <span class="text-lg">🏆</span>
            <div>
              <p class="text-xs font-bold text-amber-100 group-hover:text-amber-400">
                {{ league.name }}
              </p>
              <p class="text-[10px] text-slate-500">
                {{ league.options?.tier || 'General' }}
              </p>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-4">
          <p class="text-xs text-slate-500 italic mb-2">
            This fighter is not enrolled in any active leagues.
          </p>
          <button class="text-[10px] bg-slate-800 hover:bg-slate-700 text-amber-500 px-3 py-1 rounded border border-slate-700 uppercase font-black" @click="router.push('/leagues')">
            Find Leagues
          </button>
        </div>
      </div>
    </section>

    <!-- EQUIPMENT SLOTS (D4 influence - moody frames) -->
    <section class="ep-card-iron overhaul-weapon-panel p-6">
      <h4 class="ep-header-gold mb-6">
        Equipped Armament
      </h4>
        
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <!-- Weapon Slot -->
        <div class="flex items-center gap-4">
          <div class="relative group">
            <div class="w-24 h-24 bg-slate-950 border-2 border-slate-800 rounded flex items-center justify-center relative overflow-hidden">
              <div v-if="equippedWeapon" class="absolute inset-0 bg-gradient-to-t from-red-900/20 to-transparent"></div>
              <div v-if="equippedWeapon" class="absolute inset-0 flex items-center justify-center z-10">
                  <img
                    :src="getEquipmentImageUrl(equippedWeapon.type, equippedWeapon.id).url"
                    class="w-16 h-16 pixelated"
                    :alt="equippedWeapon.type"
                    :data-testid="`weapon-item-${equippedWeapon.id}`"
                    @error="handleImageError"
                  />
              </div>
              <span v-else class="text-slate-800 text-3xl">⚔️</span>
              <div v-if="equippedWeapon" class="absolute top-0 left-0 w-6 h-6 -translate-x-3 -translate-y-3 rotate-45" :class="rarityColors[getRarityName(equippedWeapon.rarity)]"></div>
            </div>
          </div>
          <div v-if="equippedWeapon" class="flex-1 min-w-0">
            <p class="text-[10px] text-slate-500 uppercase font-bold">Weapon</p>
            <h5 class="text-sm font-bold text-white truncate">{{ equippedWeapon.type }}</h5>
            <button class="text-[10px] text-red-400 hover:text-red-300 font-bold uppercase mt-1" @click="handleUnequip(equippedWeapon.id)">Unequip</button>
          </div>
          <div v-else class="flex-1">
            <p class="text-[10px] text-slate-600 italic">No weapon equipped</p>
            <button class="text-[10px] text-amber-500 font-bold uppercase mt-1" @click="$emit('openArmory')">Equip</button>
          </div>
        </div>

        <!-- Armor Slot -->
        <div class="flex items-center gap-4">
          <div class="relative group">
            <div class="w-24 h-24 bg-slate-950 border-2 border-slate-800 rounded flex items-center justify-center relative overflow-hidden">
              <div v-if="equippedArmor" class="absolute inset-0 bg-gradient-to-t from-blue-900/20 to-transparent"></div>
              <div v-if="equippedArmor" class="absolute inset-0 flex items-center justify-center z-10">
                  <img
                    :src="getEquipmentImageUrl(equippedArmor.type, equippedArmor.id).url"
                    class="w-16 h-16 pixelated"
                    :alt="equippedArmor.type"
                    :data-testid="`armor-item-${equippedArmor.id}`"
                    @error="handleImageError"
                  />
              </div>
              <span v-else class="text-slate-800 text-3xl">🛡️</span>
              <div v-if="equippedArmor" class="absolute top-0 left-0 w-6 h-6 -translate-x-3 -translate-y-3 rotate-45" :class="rarityColors[getRarityName(equippedArmor.rarity)]"></div>
            </div>
          </div>
          <div v-if="equippedArmor" class="flex-1 min-w-0">
            <p class="text-[10px] text-slate-500 uppercase font-bold">Armor</p>
            <h5 class="text-sm font-bold text-white truncate">{{ equippedArmor.type }}</h5>
            <button class="text-[10px] text-red-400 hover:text-red-300 font-bold uppercase mt-1" @click="handleUnequip(equippedArmor.id)">Unequip</button>
          </div>
          <div v-else class="flex-1">
            <p class="text-[10px] text-slate-600 italic">No armor equipped</p>
            <button class="text-[10px] text-amber-500 font-bold uppercase mt-1" @click="$emit('openArmory')">Equip</button>
          </div>
        </div>
      </div>

      <!-- Action Footer (Legacy Support) -->
      <div v-if="equippedWeapon || equippedArmor" class="mt-6 pt-4 border-t border-slate-800 flex justify-end">
          <button 
            class="px-4 py-1.5 bg-amber-600 hover:bg-amber-500 text-[10px] font-bold uppercase tracking-widest text-black transition-all border border-amber-400"
            data-testid="upgrade-item-button"
            @click="handleUpgrade"
          >
            Manage All in Inventory
          </button>
      </div>
    </section>
    
    <ToastManager ref="toastManager" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRosterStore } from '@/features/roster/store';
import { useLeaguesStore } from '@/features/leagues/store';
 import { getEquipmentImageUrl } from '@/shared/utils/equipmentImages';
import ToastManager from '@/shared/ui/ToastManager.vue';
import { useRouter } from 'vue-router';
import type { Fighter, Equipment } from '@/features/roster/api';

const props = defineProps<{
  fighter: Fighter;
  equipment: Equipment[];
}>();

const rosterStore = useRosterStore();
const leaguesStore = useLeaguesStore();
const router = useRouter();
const toastManager = ref<any>(null);

const emit = defineEmits(['openArmory']);

onMounted(() => {
  leaguesStore.fetchLeagues();
});

const enrolledLeagues = computed(() => {
  const allLeagues = leaguesStore.leagues;
  return allLeagues.filter(league => {
    const subs = leaguesStore.subscriptions[league.id] || [];
    return subs.some(s => s.fighterId === props.fighter.id);
  });
});

const handleUnequip = async (itemId: string) => {
    try {
        await rosterStore.unequipItemFromFighter(props.fighter.id, itemId);
    } catch (e) {
        toastManager.value?.addToast('Failed to unequip item', 'error');
    }
};

const handleUpgrade = () => {
    // Navigate to Inventory/Enhancement page or open modal
    router.push('/inventory');
};

const handleImageError = (event: Event) => {
    const img = event.target as HTMLImageElement;
    img.src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="%23333" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="12" cy="12" r="3"/><path d="M12 8v8M8 12h8"/></svg>';
};

const getExpPercent = computed(() => {
  const next = props.fighter.levelExp || props.fighter.xpToNextLevel;
  if (!next) return 0;
  const current = props.fighter.currentExp !== undefined ? props.fighter.currentExp : (props.fighter.xp || 0);
  return Math.round((current / next) * 100);
});

const getExpDisplay = computed(() => {
  const current = props.fighter.currentExp !== undefined ? props.fighter.currentExp : (props.fighter.xp || 0);
  const next = props.fighter.levelExp || props.fighter.xpToNextLevel;
  return `${current} / ${next}`;
});

// Primary item detection (weapon or armor)
const equippedWeapon = computed(() => {
  return props.equipment.find(item => 
    (item.category === 'weapon' || 
     item.type.toLowerCase().includes('weapon') || 
     item.type.toLowerCase().includes('axe') || 
     item.type.toLowerCase().includes('sword') ||
     item.type.toLowerCase().includes('staff') ||
     item.type.toLowerCase().includes('bow')) &&
    item.fighterId === props.fighter.id
  );
});

const equippedArmor = computed(() => {
  return props.equipment.find(item => 
    (item.category === 'armor' || 
     item.type.toLowerCase().includes('helmet') || 
     item.type.toLowerCase().includes('chest') || 
     item.type.toLowerCase().includes('gloves') || 
     item.type.toLowerCase().includes('boots') || 
     item.type.toLowerCase().includes('armor') || 
     item.type.toLowerCase().includes('shield')) &&
    item.fighterId === props.fighter.id
  );
});

// For backward compatibility and stats display
const equippedItem = computed(() => equippedWeapon.value || equippedArmor.value);

const getItemCategory = (type: string) => {
  const t = type.toLowerCase();
  if (t.includes('weapon') || t.includes('axe') || t.includes('sword') || t.includes('staff') || t.includes('bow')) return 'Weapon';
  if (t.includes('helmet') || t.includes('chest') || t.includes('gloves') || t.includes('boots') || t.includes('armor') || t.includes('shield')) return 'Armor';
  return 'Equipment';
};

const rarityColors: Record<string, string> = {
    'Common': 'bg-slate-500',
    'Rare': 'bg-emerald-500',
    'Legendary': 'bg-orange-500',
    'Epic': 'bg-purple-600'
};

const getRarityName = (rarity: number) => {
  if (rarity >= 4) return 'Legendary';
  if (rarity === 3) return 'Epic';
  if (rarity === 2) return 'Rare';
  return 'Common';
};

const itemStats = computed(() => {
    const stats = [];
    if (equippedWeapon.value) {
        stats.push({ label: 'Weapon', value: equippedWeapon.value.type });
        stats.push({ label: 'W-Power', value: (equippedWeapon.value.level * 10) + (equippedWeapon.value.rarity * 20) });
    }
    if (equippedArmor.value) {
        stats.push({ label: 'Armor', value: equippedArmor.value.type });
        stats.push({ label: 'A-Defense', value: (equippedArmor.value.level * 10) + (equippedArmor.value.rarity * 20) });
    }
    return stats;
});

const statGroups = computed(() => ({
  'Offense': [
    { key: 'power', label: 'Combat Power', value: props.fighter.power, suffix: '' },
    { key: 'precision', label: 'Precision', value: props.fighter.precision, suffix: '' },
    { key: 'ferocity', label: 'Ferocity', value: props.fighter.ferocity, suffix: '' },
  ],
  'Defense': [
    { key: 'vitality', label: 'Vitality Pool', value: props.fighter.vitality, suffix: '' },
    { key: 'armor', label: 'Iron Plating', value: props.fighter.armor, suffix: '' },
    { key: 'parry', label: 'Blade Parry', value: props.fighter.parryChance, suffix: '%' },
  ],
  'Precision': [
    { key: 'speed', label: 'Movement', value: props.fighter.speed, suffix: '' },
    { key: 'accuracy', label: 'Strike Rating', value: props.fighter.accuracy, suffix: '' },
    { key: 'vision', label: 'Battle Awareness', value: props.fighter.vision, suffix: '' },
  ]
}));

</script>

<style scoped>
.pixelated {
    image-rendering: pixelated;
}
</style>

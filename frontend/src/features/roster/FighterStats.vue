<template>
  <div class="space-y-8">
    <!-- Experience Progress -->
    <section class="rounded-xl border border-slate-800 bg-slate-800/30 p-4">
      <div class="mb-3 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-gradient-to-br from-amber-500 to-orange-600 text-lg font-black text-white shadow-lg">
            {{ fighter.level }}
          </div>
          <div>
            <div class="text-sm font-medium text-white">Level {{ fighter.level }}</div>
            <div class="text-xs text-slate-500">{{ fighter.currentExp }} / {{ fighter.levelExp }} XP</div>
          </div>
        </div>
        <div class="text-right">
          <div class="text-2xl font-bold text-white">{{ getExpPercent }}%</div>
          <div class="text-xs text-slate-500">to next level</div>
        </div>
      </div>
      <div class="h-3 overflow-hidden rounded-full bg-slate-800">
        <div
          class="h-full rounded-full bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 transition-all duration-700"
          :style="{ width: `${getExpPercent}%` }"
        ></div>
      </div>
    </section>

    <!-- Stats Grouped -->
    <section>
      <h4 class="mb-4 flex items-center gap-2 text-sm font-bold uppercase tracking-widest text-slate-400">
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        Combat Attributes
      </h4>
      
      <!-- Weapons Section -->
      <div class="mb-6">
        <div class="mb-2 text-xs font-medium uppercase tracking-wider text-purple-400">Primary Weapon</div>
        <div v-if="equippedWeapon" class="rounded-xl border border-purple-500/30 bg-purple-950/30 p-4">
            <div class="flex items-center justify-between">
                <div class="flex items-center gap-4">
                    <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-purple-800 to-purple-900 text-2xl shadow-lg">
                        {{ getWeaponIcon(equippedWeapon.type) }}
                    </div>
                    <div>
                        <div class="flex items-center gap-2">
                            <span class="font-bold text-purple-400">{{ equippedWeapon.name }}</span>
                            <span v-if="equippedWeapon.enhancement" class="rounded bg-purple-500/20 px-1.5 py-0.5 text-xs font-bold text-purple-400">
                                +{{ equippedWeapon.enhancement }}
                            </span>
                        </div>
                        <div class="text-xs text-slate-500">
                            {{ equippedWeapon.type }} • {{ equippedWeapon.rarity }}
                        </div>
                    </div>
                </div>
                <button 
                  @click="handleUnequipWeapon"
                  class="rounded-lg bg-slate-800 px-3 py-1.5 text-xs font-bold text-slate-400 hover:bg-red-900/40 hover:text-red-400"
                >
                    Unequip
                </button>
            </div>
            <div class="mt-3 grid grid-cols-3 gap-2 border-t border-purple-500/10 pt-3">
                <div class="text-center">
                    <div class="text-xs text-slate-500">Damage</div>
                    <div class="text-sm font-bold text-white">{{ equippedWeapon.damage }}</div>
                </div>
                <div class="text-center">
                    <div class="text-xs text-slate-500">Speed</div>
                    <div class="text-sm font-bold text-white">{{ equippedWeapon.attackSpeed }}</div>
                </div>
                <div class="text-center">
                    <div class="text-xs text-slate-500">Crit</div>
                    <div class="text-sm font-bold text-white">{{ equippedWeapon.critChance }}%</div>
                </div>
            </div>
        </div>
        <div v-else class="space-y-4">
            <div class="rounded-xl border-2 border-dashed border-slate-800 p-4 text-center">
                <p class="text-sm text-slate-500">No primary weapon equipped</p>
            </div>
            <div v-if="availableWeapons.length" class="space-y-2">
                 <div class="text-xs font-medium text-slate-500">Available in Armory:</div>
                 <div class="grid grid-cols-1 gap-2">
                    <div 
                      v-for="weapon in availableWeapons" 
                      :key="weapon.id"
                      class="flex items-center justify-between rounded-lg border border-slate-800 bg-slate-900/50 p-2 pl-3"
                    >
                        <div class="flex items-center gap-3">
                            <span class="text-lg">{{ getWeaponIcon(weapon.type) }}</span>
                            <div>
                                <span class="text-sm font-bold text-white">{{ weapon.name }}</span>
                                <span class="ml-2 text-[10px] uppercase text-slate-500">{{ weapon.rarity }}</span>
                            </div>
                        </div>
                        <button 
                          @click="handleEquipWeapon(weapon.id)"
                          class="rounded bg-indigo-600 px-3 py-1 text-xs font-bold text-white hover:bg-indigo-500"
                        >
                            Equip
                        </button>
                    </div>
                 </div>
            </div>
        </div>
      </div>

      <!-- Offense Stats -->
      <div class="mb-6">
        <div class="mb-2 text-xs font-medium uppercase tracking-wider text-red-400">Offense</div>
        <div class="grid grid-cols-2 gap-3">
          <div v-for="stat in offenseStats" :key="stat.key" class="rounded-lg border border-slate-800 bg-slate-900/50 p-3">
            <div class="mb-2 flex items-center justify-between">
              <span class="text-sm text-slate-400">{{ stat.label }}</span>
              <span class="text-lg font-bold text-white">{{ stat.value }}</span>
            </div>
            <div class="h-1.5 overflow-hidden rounded-full bg-slate-800">
              <div
                class="h-full rounded-full bg-gradient-to-r from-red-500 to-orange-500 transition-all"
                :style="{ width: `${Math.min(100, stat.value)}%` }"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Defense Stats -->
      <div class="mb-6">
        <div class="mb-2 text-xs font-medium uppercase tracking-wider text-emerald-400">Defense</div>
        <div class="grid grid-cols-2 gap-3">
          <div v-for="stat in defenseStats" :key="stat.key" class="rounded-lg border border-slate-800 bg-slate-900/50 p-3">
            <div class="mb-2 flex items-center justify-between">
              <span class="text-sm text-slate-400">{{ stat.label }}</span>
              <span class="text-lg font-bold text-white">{{ stat.value }}</span>
            </div>
            <div class="h-1.5 overflow-hidden rounded-full bg-slate-800">
              <div
                class="h-full rounded-full bg-gradient-to-r from-emerald-500 to-teal-500 transition-all"
                :style="{ width: `${Math.min(100, stat.value)}%` }"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Utility Stats -->
      <div>
        <div class="mb-2 text-xs font-medium uppercase tracking-wider text-sky-400">Utility</div>
        <div class="grid grid-cols-2 gap-3">
          <div v-for="stat in utilityStats" :key="stat.key" class="rounded-lg border border-slate-800 bg-slate-900/50 p-3">
            <div class="mb-2 flex items-center justify-between">
              <span class="text-sm text-slate-400">{{ stat.label }}</span>
              <span class="text-lg font-bold text-white">{{ stat.value }}{{ stat.suffix || '' }}</span>
            </div>
            <div class="h-1.5 overflow-hidden rounded-full bg-slate-800">
              <div
                class="h-full rounded-full bg-gradient-to-r from-sky-500 to-indigo-500 transition-all"
                :style="{ width: `${Math.min(100, stat.rawValue)}%` }"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Equipment Section -->
    <section>
      <h4 class="mb-4 flex items-center gap-2 text-sm font-bold uppercase tracking-widest text-slate-400">
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
        </svg>
        Equipment
      </h4>
      
      <div v-if="equipment.length" class="space-y-2">
        <div
          v-for="item in equipment"
          :key="item.id"
          :class="[
            'flex items-center justify-between rounded-xl p-4 transition-all',
            getRarityBorder(item.rarity)
          ]"
        >
          <div class="flex items-center gap-4">
            <div
              :class="[
                'flex h-12 w-12 items-center justify-center rounded-xl text-2xl shadow-lg',
                getRarityClass(item.rarity)
              ]"
            >
              {{ getItemIcon(item.type) }}
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span :class="['font-bold', getRarityTextClass(item.rarity)]">
                  {{ item.type }}
                </span>
                <span v-if="item.enhancement" class="rounded bg-emerald-500/20 px-1.5 py-0.5 text-xs font-bold text-emerald-400">
                  +{{ item.enhancement }}
                </span>
              </div>
              <div class="flex items-center gap-2 text-xs text-slate-500">
                <span>Level {{ item.level }}</span>
                <span class="h-1 w-1 rounded-full bg-slate-600"></span>
                <span :class="getRarityTextClass(item.rarity)">{{ getRarityName(item.rarity) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="rounded-xl border-2 border-dashed border-slate-800 p-8 text-center">
        <div class="mb-2 text-3xl">📦</div>
        <p class="text-sm text-slate-500">No equipment assigned</p>
        <p class="mt-1 text-xs text-slate-600">Visit the Inventory to equip gear</p>
      </div>
    </section>

    <!-- Attunement Section -->
    <section>
      <h4 class="mb-4 flex items-center gap-2 text-sm font-bold uppercase tracking-widest text-slate-400">
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        Elemental Attunement
      </h4>
      
      <div class="grid grid-cols-5 gap-2">
        <button
          v-for="type in attunementTypes"
          :key="type.id"
          @click="updateAttunement(type.id)"
          :class="[
            'flex flex-col items-center gap-2 rounded-xl p-4 transition-all',
            fighter.attunementId === type.id
              ? `${type.bgActive} ring-2 ${type.ring} shadow-lg ${type.shadow}`
              : 'border border-slate-800 bg-slate-900/50 hover:bg-slate-800'
          ]"
        >
          <span class="text-3xl">{{ type.icon }}</span>
          <span
            :class="[
              'text-xs font-bold uppercase tracking-wider',
              fighter.attunementId === type.id ? 'text-white' : 'text-slate-500'
            ]"
          >
            {{ type.name }}
          </span>
        </button>
      </div>
      
      <p v-if="fighter.attunementId" class="mt-3 text-sm text-slate-500">
        {{ getAttunementDescription(fighter.attunementId) }}
      </p>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useRosterStore } from '@/features/roster/store';
import { useWeaponsStore } from '@/features/weapons/store';
import type { Fighter, Equipment } from '@/features/roster/api';

const props = defineProps<{
  fighter: Fighter;
  equipment: Equipment[];
}>();

const roster = useRosterStore();
const weaponsStore = useWeaponsStore();

onMounted(() => {
    weaponsStore.fetchWeapons();
});

const equippedWeapon = computed(() => {
    return weaponsStore.weapons.find(w => w.fighterId === props.fighter.id);
});

const availableWeapons = computed(() => {
    return weaponsStore.weapons.filter(w => !w.isEquipped);
});

const handleEquipWeapon = async (weaponId: string) => {
    try {
        await weaponsStore.equipWeapon(weaponId, props.fighter.id);
    } catch (e) {
        console.error("Failed to equip weapon", e);
    }
}

const handleUnequipWeapon = async () => {
    if (equippedWeapon.value) {
        try {
            await weaponsStore.unequipWeapon(equippedWeapon.value.id);
        } catch (e) {
            console.error("Failed to unequip weapon", e);
        }
    }
}

const getWeaponIcon = (type: string) => {
  const lower = type.toLowerCase();
  if (lower.includes('sword')) return '⚔️';
  if (lower.includes('axe')) return '🪓';
  if (lower.includes('staff')) return '🪄';
  if (lower.includes('dagger')) return '🗡️';
  if (lower.includes('bow')) return '🏹';
  return '⚔️';
};

const attunementTypes = [
  { id: 'Fire', name: 'Fire', icon: '🔥', bgActive: 'bg-orange-900/50', ring: 'ring-orange-500', shadow: 'shadow-orange-500/30' },
  { id: 'Water', name: 'Water', icon: '💧', bgActive: 'bg-blue-900/50', ring: 'ring-blue-500', shadow: 'shadow-blue-500/30' },
  { id: 'Earth', name: 'Earth', icon: '🪨', bgActive: 'bg-amber-900/50', ring: 'ring-amber-600', shadow: 'shadow-amber-500/30' },
  { id: 'Wind', name: 'Wind', icon: '💨', bgActive: 'bg-teal-900/50', ring: 'ring-teal-500', shadow: 'shadow-teal-500/30' },
  { id: 'Lightning', name: 'Lightning', icon: '⚡', bgActive: 'bg-yellow-900/50', ring: 'ring-yellow-500', shadow: 'shadow-yellow-500/30' },
];

const getAttunementDescription = (id: string) => {
  const descriptions: Record<string, string> = {
    Fire: 'Increases damage output and burn effects in combat.',
    Water: 'Enhances healing abilities and resistance to fire.',
    Earth: 'Boosts armor and stability, reducing knockback.',
    Wind: 'Improves speed and evasion, enabling faster attacks.',
    Lightning: 'Adds chain damage and stun potential to attacks.',
  };
  return descriptions[id] || '';
};

const updateAttunement = async (id: string) => {
  const newId = props.fighter.attunementId === id ? null : id;
  await roster.updateAttunement(props.fighter.id, newId);
};

const getExpPercent = computed(() => {
  if (!props.fighter.xpToNextLevel) return 0;
  return Math.round((props.fighter.xp / props.fighter.xpToNextLevel) * 100);
});

const offenseStats = computed(() => [
  { key: 'power', label: 'Power', value: props.fighter.power },
  { key: 'precision', label: 'Precision', value: props.fighter.precision },
  { key: 'ferocity', label: 'Ferocity', value: props.fighter.ferocity },
  { key: 'conditionPower', label: 'Condition Power', value: props.fighter.conditionPower },
]);

const defenseStats = computed(() => [
  { key: 'vitality', label: 'Vitality', value: props.fighter.vitality },
  { key: 'armor', label: 'Armor', value: props.fighter.armor },
  { key: 'parryChance', label: 'Parry Chance', value: props.fighter.parryChance, rawValue: props.fighter.parryChance, suffix: '%' },
  { key: 'healingPower', label: 'Healing Power', value: props.fighter.healingPower },
]);

const utilityStats = computed(() => [
  { key: 'speed', label: 'Speed', value: props.fighter.speed, rawValue: props.fighter.speed, suffix: '' },
  { key: 'agility', label: 'Agility', value: props.fighter.agility, rawValue: props.fighter.agility, suffix: '' },
  { key: 'accuracy', label: 'Accuracy', value: props.fighter.accuracy, rawValue: props.fighter.accuracy, suffix: '' },
  { key: 'vision', label: 'Vision', value: props.fighter.vision, rawValue: props.fighter.vision, suffix: '' },
]);

const getItemIcon = (type: string) => {
  const lower = type.toLowerCase();
  if (lower.includes('sword') || lower.includes('weapon') || lower.includes('axe') || lower.includes('dagger')) return '⚔️';
  if (lower.includes('staff') || lower.includes('wand')) return '🪄';
  if (lower.includes('bow')) return '🏹';
  if (lower.includes('chest') || lower.includes('armor') || lower.includes('plate')) return '🛡️';
  if (lower.includes('head') || lower.includes('helmet') || lower.includes('helm')) return '🪖';
  if (lower.includes('legs') || lower.includes('pants') || lower.includes('greaves')) return '👖';
  if (lower.includes('shoes') || lower.includes('boots') || lower.includes('feet')) return '🥾';
  if (lower.includes('gloves') || lower.includes('hands') || lower.includes('gauntlets')) return '🧤';
  if (lower.includes('ring')) return '💍';
  if (lower.includes('amulet') || lower.includes('necklace')) return '📿';
  return '📦';
};

const getRarityClass = (rarity: number) => {
  switch (rarity) {
    case 0: return 'bg-slate-700 text-slate-300';
    case 1: return 'bg-gradient-to-br from-emerald-800 to-emerald-900 text-emerald-300';
    case 2: return 'bg-gradient-to-br from-blue-800 to-blue-900 text-blue-300';
    case 3: return 'bg-gradient-to-br from-purple-800 to-purple-900 text-purple-300';
    case 4: return 'bg-gradient-to-br from-amber-700 to-orange-800 text-amber-200';
    default: return 'bg-slate-700';
  }
};

const getRarityBorder = (rarity: number) => {
  switch (rarity) {
    case 0: return 'border border-slate-800 bg-slate-900/50';
    case 1: return 'border border-emerald-500/30 bg-emerald-950/30';
    case 2: return 'border border-blue-500/30 bg-blue-950/30';
    case 3: return 'border border-purple-500/30 bg-purple-950/30';
    case 4: return 'border border-amber-500/30 bg-amber-950/30';
    default: return 'border border-slate-800 bg-slate-900/50';
  }
};

const getRarityTextClass = (rarity: number) => {
  switch (rarity) {
    case 0: return 'text-slate-400';
    case 1: return 'text-emerald-400';
    case 2: return 'text-blue-400';
    case 3: return 'text-purple-400';
    case 4: return 'text-amber-400';
    default: return 'text-slate-400';
  }
};

const getRarityName = (rarity: number) => {
  switch (rarity) {
    case 0: return 'Common';
    case 1: return 'Rare';
    case 2: return 'Fabled';
    case 3: return 'Mythic';
    case 4: return 'Legendary';
    default: return 'Unknown';
  }
};
</script>

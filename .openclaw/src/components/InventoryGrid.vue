<template>
  <div class="inventory-container">
    <div class="forge-header">
      <h2>Darkened Forge Inventory</h2>
    </div>
    <div class="inventory-grid">
      <div 
        v-for="index in 32" 
        :key="index" 
        class="grid-slot"
        @mouseenter="showTooltip($event, getItemAt(index - 1))"
        @mouseleave="hideTooltip"
      >
        <div v-if="getItemAt(index - 1)" class="item-icon-wrapper">
          <img :src="getItemAt(index - 1).icon" :alt="getItemAt(index - 1).name" class="item-icon" />
        </div>
      </div>
    </div>

    <!-- Tooltip -->
    <div 
      v-if="hoveredItem" 
      class="item-tooltip" 
      :style="{ top: tooltipY + 'px', left: tooltipX + 'px' }"
    >
      <div class="tooltip-header" :class="hoveredItem.rarity">
        {{ hoveredItem.name }}
      </div>
      <div class="tooltip-body">
        <p class="item-type">{{ hoveredItem.type }}</p>
        <div class="item-stats">
          <div v-for="(val, stat) in hoveredItem.stats" :key="stat" class="stat-line">
            <span class="stat-name">{{ stat }}:</span>
            <span class="stat-value">{{ val }}</span>
          </div>
        </div>
        <p class="item-description">{{ hoveredItem.description }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const items = ref([
  {
    name: "Infernal Cleaver",
    type: "Two-Handed Axe",
    rarity: "legendary",
    icon: "https://vibemedia.space/wpn_infernal_cleaver_8a2b3c4d.png?prompt=legendary%20infernal%20cleaver%20axe%20with%20glowing%20magma%20edges%20and%20demon%20bone%20handle&style=pixel_game_asset&key=NOGON",
    stats: { "Damage": "145-180", "Strength": "+25", "Fire Damage": "+40" },
    description: "Forged in the deepest pits of the abyss, this blade hungers for souls.",
    slot: 2
  },
  {
    name: "Shadow-Stiched Robes",
    type: "Chest Armor",
    rarity: "rare",
    icon: "https://vibemedia.space/icon_shadow_robe_1a9b2c3d.png?prompt=dark%20shadow%20priest%20robes%20with%20purple%20glow%20and%20mystical%20patterns&style=pixel_game_asset&key=NOGON",
    stats: { "Armor": "320", "Intelligence": "+18", "Mana Regen": "+5%" },
    description: "Threads of pure darkness woven by a mad sorcerer.",
    slot: 5
  },
  {
    name: "Crystalline Bulwark",
    type: "Shield",
    rarity: "legendary",
    icon: "https://vibemedia.space/wpn_crystal_shield_7x8y9z0w.png?prompt=legendary%20crystalline%20shield%20made%20of%20glowing%20blue%20shards&style=pixel_game_asset&key=NOGON",
    stats: { "Block Chance": "25%", "Armor": "450", "All Resistances": "+15" },
    description: "Impenetrable. Reflects the light of a thousand suns.",
    slot: 12
  },
  {
    name: "Void-Touched Dagger",
    type: "One-Handed Dagger",
    rarity: "uncommon",
    icon: "https://vibemedia.space/wpn_void_dagger_q1w2e3r4.png?prompt=void%20touched%20dagger%20with%20pulsing%20purple%20energy%20and%20ebony%20blade&style=pixel_game_asset&key=NOGON",
    stats: { "Damage": "45-60", "Attack Speed": "+15%", "Poison Damage": "+12" },
    description: "A single scratch from this blade is often enough to end a life.",
    slot: 18
  },
  {
    name: "Ironclad Helm",
    type: "Headpiece",
    rarity: "common",
    icon: "https://vibemedia.space/icon_iron_helm_5t6y7u8i.png?prompt=heavy%20iron%20knight%20helmet%20with%20battle%20scars%20and%20plume&style=pixel_game_asset&key=NOGON",
    stats: { "Armor": "120", "Vitality": "+10" },
    description: "Standard issue for the Vanguard, seasoned by countless sieges.",
    slot: 26
  }
]);

const hoveredItem = ref(null);
const tooltipX = ref(0);
const tooltipY = ref(0);

const getItemAt = (index) => {
  return items.value.find(item => item.slot === index);
};

const showTooltip = (event, item) => {
  if (!item) return;
  hoveredItem.value = item;
  tooltipX.value = event.clientX + 15;
  tooltipY.value = event.clientY + 15;
};

const hideTooltip = () => {
  hoveredItem.value = null;
};
</script>

<style scoped>
.inventory-container {
  background-color: #1a1a1a;
  padding: 20px;
  border: 4px solid #3d3d3d;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.8);
  width: fit-content;
  user-select: none;
  font-family: 'Cinzel', serif; /* Assuming a gothic-style font is available or fall back */
}

.forge-header h2 {
  color: #c0c0c0;
  text-align: center;
  margin-bottom: 20px;
  text-transform: uppercase;
  letter-spacing: 2px;
  border-bottom: 2px solid #4a3c2c;
  padding-bottom: 10px;
}

.inventory-grid {
  display: grid;
  grid-template-columns: repeat(8, 64px);
  grid-template-rows: repeat(4, 64px);
  gap: 4px;
  background-color: #0d0d0d;
  padding: 8px;
  border: 2px solid #4a3c2c;
}

.grid-slot {
  width: 64px;
  height: 64px;
  background-color: #151515;
  border: 1px solid #2a2a2a;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: inset 0 0 10px rgba(0, 0, 0, 0.5);
}

.grid-slot:hover {
  border-color: #666;
  background-color: #202020;
}

.item-icon-wrapper {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.item-icon {
  max-width: 100%;
  max-height: 100%;
  image-rendering: pixelated;
}

/* Tooltip Styles */
.item-tooltip {
  position: fixed;
  z-index: 100;
  width: 240px;
  background-color: rgba(10, 10, 10, 0.95);
  border: 2px solid #5a4a3a;
  box-shadow: 10px 10px 20px rgba(0, 0, 0, 0.9);
  pointer-events: none;
  padding: 0;
}

.tooltip-header {
  padding: 8px 12px;
  font-weight: bold;
  font-size: 1.1rem;
  border-bottom: 1px solid #444;
}

.tooltip-header.common { color: #ffffff; }
.tooltip-header.uncommon { color: #1eff00; }
.tooltip-header.rare { color: #0070dd; }
.tooltip-header.legendary { color: #ff8000; text-shadow: 0 0 5px #ff8000; }

.tooltip-body {
  padding: 12px;
}

.item-type {
  color: #9d9d9d;
  font-style: italic;
  font-size: 0.9rem;
  margin-bottom: 8px;
}

.item-stats {
  margin-bottom: 12px;
}

.stat-line {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
  font-size: 0.95rem;
}

.stat-name {
  color: #adb5bd;
}

.stat-value {
  color: #eee;
}

.item-description {
  color: #ca9e00;
  font-size: 0.85rem;
  line-height: 1.4;
  margin-top: 10px;
  border-top: 1px solid #333;
  padding-top: 8px;
}
</style>

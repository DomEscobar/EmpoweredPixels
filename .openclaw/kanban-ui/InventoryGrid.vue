<template>
  <div class="inventory-container">
    <div class="inventory-frame">
      <div class="inventory-header">
        <span class="inventory-title">INVENTORY</span>
      </div>
      <div class="inventory-grid">
        <div 
          v-for="(slot, index) in slots" 
          :key="index" 
          class="inventory-slot"
          @mouseenter="hoverItem = slot.item"
          @mouseleave="hoverItem = null"
        >
          <div v-if="slot.item" class="item-icon-wrapper">
            <img :src="slot.item.icon" :alt="slot.item.name" class="item-icon" />
          </div>
        </div>
      </div>
    </div>

    <!-- Tooltip -->
    <div v-if="hoverItem" class="item-tooltip" :style="tooltipStyle">
      <div class="tooltip-header" :class="hoverItem.rarity">
        {{ hoverItem.name }}
      </div>
      <div class="tooltip-type">{{ hoverItem.type }}</div>
      <div class="tooltip-stats">
        <div v-for="(stat, key) in hoverItem.stats" :key="key" class="tooltip-stat">
          <span class="stat-name">{{ key }}:</span>
          <span class="stat-value">{{ stat }}</span>
        </div>
      </div>
      <div class="tooltip-description">{{ hoverItem.description }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'InventoryGrid',
  data() {
    const items = [
      {
        name: 'Iron Sword of Might',
        type: 'One-Handed Sword',
        rarity: 'magic',
        icon: 'https://vibemedia.space/wpn_iron_sword_v1.png?prompt=iron%20sword%20with%20leather%20grip%20and%20sharp%20blade&style=pixel_game_asset&key=NOGON',
        stats: { Damage: '12-18', Speed: 'Fast' },
        description: 'A sturdy iron sword, favored by guards and militia alike.'
      },
      {
        name: 'Aegis of Protection',
        type: 'Shield',
        rarity: 'rare',
        icon: 'https://vibemedia.space/shd_aegis_v1.png?prompt=ornate%20iron%20shield%20with%20wolf%20emblem&style=pixel_game_asset&key=NOGON',
        stats: { Defense: '+45', Block: '15%' },
        description: 'Rumored to have blocked a dragon\'s breath.'
      },
      {
        name: 'Greater Healing Potion',
        type: 'Consumable',
        rarity: 'common',
        icon: 'https://vibemedia.space/pot_healing_v1.png?prompt=glass%20bottle%20with%20swirling%20red%20liquid&style=pixel_game_asset&key=NOGON',
        stats: { Heal: '200 HP' },
        description: 'Tastes like strawberries and old parchment.'
      },
      {
        name: 'Mana Crystal',
        type: 'Consumable',
        rarity: 'magic',
        icon: 'https://vibemedia.space/itm_mana_v1.png?prompt=glowing%20blue%20crystal%20shard&style=pixel_game_asset&key=NOGON',
        stats: { Mana: '+100' },
        description: 'Vibrates faintly in your hand.'
      },
      {
        name: 'Shadow Dagger',
        type: 'One-Handed Dagger',
        rarity: 'legendary',
        icon: 'https://vibemedia.space/wpn_dagger_v1.png?prompt=obsidian%20dagger%20with%20purple%20glow&style=pixel_game_asset&key=NOGON',
        stats: { Damage: '8-12', Crit: '+20%' },
        description: 'Forged in the void between worlds.'
      }
    ];

    // Create 8x4 grid (32 slots)
    const slots = Array(32).fill(null).map((_, i) => ({
      item: i < items.length ? items[i] : null
    }));

    return {
      slots,
      hoverItem: null,
      mouseX: 0,
      mouseY: 0
    };
  },
  computed: {
    tooltipStyle() {
      return {
        left: `${this.mouseX + 15}px`,
        top: `${this.mouseY + 15}px`
      };
    }
  },
  mounted() {
    window.addEventListener('mousemove', this.updateMousePosition);
  },
  beforeUnmount() {
    window.removeEventListener('mousemove', this.updateMousePosition);
  },
  methods: {
    updateMousePosition(e) {
      this.mouseX = e.clientX;
      this.mouseY = e.clientY;
    }
  }
};
</script>

<style scoped>
.inventory-container {
  display: flex;
  justify-content: center;
  padding: 20px;
  background: #000;
  min-height: 400px;
}

.inventory-frame {
  background: #1a1a1a;
  border: 4px solid #4a4a4a;
  box-shadow: inset 0 0 15px #000;
  padding: 8px;
  border-radius: 4px;
  width: fit-content;
}

.inventory-header {
  text-align: center;
  border-bottom: 2px solid #333;
  margin-bottom: 8px;
  padding: 4px;
}

.inventory-title {
  color: #8a0707;
  font-family: 'Georgia', serif;
  font-weight: bold;
  letter-spacing: 2px;
  text-shadow: 1px 1px 2px #000;
}

.inventory-grid {
  display: grid;
  grid-template-columns: repeat(8, 48px);
  grid-template-rows: repeat(4, 48px);
  gap: 2px;
  background: #2b2b2b;
  border: 2px solid #111;
}

.inventory-slot {
  width: 48px;
  height: 48px;
  background: #0f0f0f;
  border: 1px solid #222;
  box-shadow: inset 0 0 5px #000;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background 0.2s;
}

.inventory-slot:hover {
  background: #1a1a1a;
  border-color: #555;
}

.item-icon-wrapper {
  width: 40px;
  height: 40px;
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
  z-index: 1000;
  background: rgba(0, 0, 0, 0.9);
  border: 2px solid #5a5a5a;
  padding: 10px;
  width: 220px;
  pointer-events: none;
  font-family: sans-serif;
  color: #ccc;
  box-shadow: 0 4px 15px rgba(0,0,0,0.8);
}

.tooltip-header {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 4px;
  border-bottom: 1px solid #333;
  padding-bottom: 4px;
}

.common { color: #fff; }
.magic { color: #4d4dff; }
.rare { color: #ffff00; }
.legendary { color: #ff8000; }

.tooltip-type {
  font-size: 12px;
  color: #888;
  font-style: italic;
  margin-bottom: 8px;
}

.tooltip-stats {
  margin-bottom: 8px;
}

.tooltip-stat {
  font-size: 13px;
  display: flex;
  justify-content: space-between;
}

.stat-name { color: #aaa; }
.stat-value { color: #4e4; }

.tooltip-description {
  font-size: 11px;
  color: #666;
  line-height: 1.2;
}
</style>

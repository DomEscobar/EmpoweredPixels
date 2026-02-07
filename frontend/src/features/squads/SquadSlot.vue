<template>
  <div class="squad-slot">
    <div class="slot-frame">
      <!-- Empty State -->
      <div v-if="!fighter" class="slot-empty">
        <span class="slot-number">{{ slotIndex + 1 }}</span>
        <div class="slot-icon">⚔️</div>
        <p class="slot-placeholder">Empty</p>
      </div>

      <!-- Fighter State -->
      <div v-else class="slot-content">
        <img
          :src="fighter.avatarUrl || fighterIcon"
          :alt="fighter.name"
          class="fighter-avatar"
        />
        <div class="fighter-info">
          <h4 class="fighter-name">{{ fighter.name }}</h4>
          <p class="fighter-level">Lv.{{ fighter.level }}</p>
          <div class="fighter-stats">
            <span class="stat">
              <span class="stat-icon">⚡</span>
              {{ fighter.power }}
            </span>
            <span class="stat">
              <span class="stat-icon">🛡️</span>
              {{ fighter.armor }}
            </span>
          </div>
        </div>
      </div>

      <!-- Active Badge -->
      <div v-if="isActive && slotIndex === 0" class="active-badge">
        <span>★ ACTIVE</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  fighter: any; // Fighter data from roster
  slotIndex: number;
  isActive?: boolean;
}>();

const fighterIcon = computed(() => {
  return `https://vibemedia.space/char_knight_walk_012jasd.png?prompt=fantasy%20knight%20fighter%20avatar&style=pixel_game_asset&key=NOGON`;
});
</script>

<style scoped>
.squad-slot {
  perspective: 1000px;
}

.slot-frame {
  position: relative;
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.95), rgba(30, 41, 59, 0.95));
  border: 3px solid rgba(251, 191, 36, 0.3);
  border-radius: 16px;
  overflow: hidden;
  box-shadow:
    0 0 20px rgba(251, 191, 36, 0.1),
    inset 0 0 40px rgba(0, 0, 0, 0.5);
  transition: all 0.3s ease;
}

.slot-frame:hover {
  border-color: rgba(251, 191, 36, 0.6);
  box-shadow:
    0 0 30px rgba(251, 191, 36, 0.2),
    inset 0 0 40px rgba(0, 0, 0, 0.5);
  transform: translateY(-2px);
}

.slot-empty {
  height: 100%;
  min-height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #64748b;
  position: relative;
}

.slot-number {
  position: absolute;
  top: 8px;
  left: 8px;
  font-size: 10px;
  font-weight: bold;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

.slot-icon {
  font-size: 48px;
  opacity: 0.3;
  margin-bottom: 8px;
}

.slot-placeholder {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.slot-content {
  height: 100%;
  min-height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.fighter-avatar {
  width: 64px;
  height: 64px;
  border-radius: 8px;
  border: 2px solid rgba(251, 191, 36, 0.3);
  margin-bottom: 8px;
  image-rendering: pixelated;
}

.fighter-info {
  text-align: center;
}

.fighter-name {
  color: #fbbf24;
  font-weight: bold;
  font-size: 14px;
  margin-bottom: 2px;
  text-shadow: 0 0 10px rgba(251, 191, 36, 0.3);
}

.fighter-level {
  color: #94a3b8;
  font-size: 12px;
  margin-bottom: 8px;
}

.fighter-stats {
  display: flex;
  gap: 12px;
}

.stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  font-size: 10px;
  font-weight: bold;
  color: #e2e8f0;
}

.stat-icon {
  font-size: 14px;
  margin-bottom: 2px;
}

.active-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: black;
  font-size: 10px;
  font-weight: bold;
  padding: 4px 8px;
  border-radius: 4px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  box-shadow: 0 0 15px rgba(245, 158, 11, 0.5);
}

.active-badge span {
  display: block;
  text-align: center;
}
</style>

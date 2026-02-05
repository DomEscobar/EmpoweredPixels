<template>
  <div class="attunement-page">
    <header class="page-header">
      <h1>✨ Attunements</h1>
      <p class="subtitle">Master the 6 elements to gain powerful bonuses</p>
    </header>

    <div v-if="attunementStore.loading" class="loading">
      Loading attunements...
    </div>

    <div v-else-if="attunementStore.error" class="error">
      {{ attunementStore.error }}
    </div>

    <div v-else class="attunements-grid">
      <div 
        v-for="att in attunementStore.attunements" 
        :key="att.element"
        class="attunement-card"
        :style="{ borderColor: getConfig(att.element).color }"
      >
        <div class="card-header" :style="{ backgroundColor: getConfig(att.element).color + '20' }">
          <span class="element-icon">{{ getConfig(att.element).icon }}</span>
          <h3>{{ getConfig(att.element).name }}</h3>
          <span class="level-badge">Level {{ att.level }}</span>
        </div>

        <p class="description">{{ getConfig(att.element).description }}</p>

        <div class="progress-section">
          <div class="progress-bar">
            <div 
              class="progress-fill" 
              :style="{ 
                width: getProgress(att) + '%',
                backgroundColor: getConfig(att.element).color
              }"
            />
          </div>
          <span class="xp-text">{{ formatXP(att) }} XP</span>
        </div>

        <div class="bonuses">
          <div class="bonus-row">
            <span>Power</span>
            <span :style="{ color: getConfig(att.element).color }">+{{ getBonus(att.element, 'power') }}%</span>
          </div>
          <div class="bonus-row">
            <span>Defense</span>
            <span :style="{ color: getConfig(att.element).color }">+{{ getBonus(att.element, 'defense') }}%</span>
          </div>
          <div class="bonus-row">
            <span>Speed</span>
            <span :style="{ color: getConfig(att.element).color }">+{{ getBonus(att.element, 'speed') }}%</span>
          </div>
          <div class="bonus-row">
            <span>Precision</span>
            <span :style="{ color: getConfig(att.element).color }">+{{ getBonus(att.element, 'precision') }}%</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="attunementStore.bonuses" class="total-bonuses">
      <h2>Total Bonuses</h2>
      <div class="bonus-grid">
        <div class="total-item">
          <span class="label">Total Power</span>
          <span class="value">+{{ attunementStore.bonuses.total_power.toFixed(1) }}%</span>
        </div>
        <div class="total-item">
          <span class="label">Total Defense</span>
          <span class="value">+{{ attunementStore.bonuses.total_defense.toFixed(1) }}%</span>
        </div>
        <div class="total-item">
          <span class="label">Total Speed</span>
          <span class="value">+{{ attunementStore.bonuses.total_speed.toFixed(1) }}%</span>
        </div>
        <div class="total-item">
          <span class="label">Total Precision</span>
          <span class="value">+{{ attunementStore.bonuses.total_precision.toFixed(1) }}%</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useAttunementStore } from '../features/attunement/store';
import type { Attunement } from '../features/attunement/types';

const attunementStore = useAttunementStore();

function getConfig(element: string) {
  return attunementStore.getElementConfig(element);
}

function getProgress(att: Attunement): number {
  // Get XP required for next level (simplified calculation)
  const xpRequired = getXPRequired(att.level);
  return Math.min(100, Math.round((att.current_xp / xpRequired) * 100));
}

function getXPRequired(level: number): number {
  const requirements = [0, 100, 250, 450, 700, 1000, 1350, 1750, 2200, 2700, 3250, 3850, 4500, 5200, 5950, 6750, 7600, 8500, 9450, 10450, 11500, 12600, 13750, 14950, 16200];
  return requirements[level] || 16200;
}

function formatXP(att: Attunement): string {
  const xpRequired = getXPRequired(att.level);
  return `${att.current_xp.toLocaleString()} / ${xpRequired.toLocaleString()}`;
}

function getBonus(element: string, type: 'power' | 'defense' | 'speed' | 'precision'): string {
  const bonuses = attunementStore.bonuses?.by_element[element];
  if (!bonuses) return '0.0';
  return bonuses[type].toFixed(1);
}

onMounted(() => {
  attunementStore.fetchAttunements();
  attunementStore.fetchBonuses();
});
</script>

<style scoped>
.attunement-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem;
}

.page-header {
  text-align: center;
  margin-bottom: 2rem;
}

.page-header h1 {
  margin: 0;
  font-size: 2rem;
  color: white;
}

.subtitle {
  color: #9ca3af;
  margin-top: 0.5rem;
}

.attunements-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.attunement-card {
  border: 2px solid;
  border-radius: 12px;
  overflow: hidden;
  background: #1f2937;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
}

.element-icon {
  font-size: 1.5rem;
}

.card-header h3 {
  margin: 0;
  flex: 1;
  color: white;
}

.level-badge {
  padding: 0.25rem 0.75rem;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 700;
  color: white;
}

.description {
  padding: 0 1rem;
  margin: 0 0 1rem 0;
  color: #9ca3af;
  font-size: 0.875rem;
}

.progress-section {
  padding: 0 1rem;
  margin-bottom: 1rem;
}

.progress-bar {
  height: 8px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.xp-text {
  font-size: 0.75rem;
  color: #9ca3af;
}

.bonuses {
  padding: 0 1rem 1rem;
}

.bonus-row {
  display: flex;
  justify-content: space-between;
  padding: 0.375rem 0;
  font-size: 0.875rem;
}

.bonus-row span:first-child {
  color: #9ca3af;
}

.bonus-row span:last-child {
  font-weight: 600;
}

.total-bonuses {
  background: #1f2937;
  border-radius: 12px;
  padding: 1.5rem;
}

.total-bonuses h2 {
  margin: 0 0 1rem 0;
  color: white;
  text-align: center;
}

.bonus-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

.total-item {
  text-align: center;
  padding: 1rem;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
}

.total-item .label {
  display: block;
  font-size: 0.75rem;
  color: #9ca3af;
  margin-bottom: 0.25rem;
}

.total-item .value {
  font-size: 1.25rem;
  font-weight: 700;
  color: #22c55e;
}

.loading, .error {
  text-align: center;
  padding: 3rem;
  color: #9ca3af;
}

.error {
  color: #ef4444;
}
</style>

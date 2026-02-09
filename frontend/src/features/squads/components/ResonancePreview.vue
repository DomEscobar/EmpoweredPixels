<template>
  <div class="resonance-preview">
    <div class="resonance-header">
      <h3>Squad Resonance Analysis</h3>
    </div>

    <!-- Attunement Circles -->
    <div class="attunement-circles">
      <div
        v-for="(element, idx) in displayElements"
        :key="idx"
        class="attunement-circle"
        :style="{ borderColor: getElementColor(element) }"
      >
        <span class="element-icon">{{ getElementIcon(element) }}</span>
        <span class="element-name">{{ formatElement(element) }}</span>
      </div>
    </div>

    <!-- Harmony Lines -->
    <svg v-if="resonanceState" class="harmony-lines" viewBox="0 0 300 100">
      <!-- Harmonic connections (green) -->
      <g stroke="#22c55e" stroke-width="2" stroke-dasharray="0" opacity="0.6">
        <line
          v-for="(pair, idx) in resonanceState.pattern?.HarmonicPairs || []"
          :key="`harmonic-${idx}`"
          :x1="getElementX(pair.Element1)"
          y1="50"
          :x2="getElementX(pair.Element2)"
          y2="50"
        />
      </g>

      <!-- Dissonant connections (red dashed) -->
      <g stroke="#ef4444" stroke-width="2" stroke-dasharray="5,5" opacity="0.6">
        <line
          v-for="(pair, idx) in resonanceState.pattern?.DissonantPairs || []"
          :key="`dissonant-${idx}`"
          :x1="getElementX(pair.Element1)"
          y1="50"
          :x2="getElementX(pair.Element2)"
          y2="50"
        />
      </g>
    </svg>

    <!-- Harmony Score Gauge -->
    <div class="harmony-gauge">
      <div class="gauge-label">Harmony Score</div>
      <div class="gauge-bar-container">
        <div
          class="gauge-bar-fill"
          :style="{
            width: `${resonanceState?.harmonyScore || 0}%`,
            backgroundColor: getGaugeColor(resonanceState?.harmonyScore || 0),
          }"
        ></div>
      </div>
      <div class="gauge-value">{{ resonanceState?.harmonyScore || 0 }}/100</div>
    </div>

    <!-- Tier Name and Description -->
    <div v-if="resonanceState" class="tier-info">
      <div class="tier-name" :style="{ color: getTierColor(resonanceState.tierName) }">
        {{ resonanceState.tierName }}
      </div>
      <div class="tier-description">
        {{ getTierDescription(resonanceState.tierName) }}
      </div>
    </div>

    <!-- Bonus Table -->
    <div v-if="resonanceState" class="bonus-table">
      <div class="bonus-row">
        <span class="bonus-label">Damage Bonus</span>
        <span class="bonus-value" :style="{ color: getBonusColor(resonanceState.bonusDamage) }">
          {{ formatBonus(resonanceState.bonusDamage) }}
        </span>
      </div>
      <div class="bonus-row">
        <span class="bonus-label">Defense Bonus</span>
        <span class="bonus-value" :style="{ color: getBonusColor(resonanceState.bonusDefense) }">
          {{ formatBonus(resonanceState.bonusDefense) }}
        </span>
      </div>
    </div>

    <!-- Dissonance Warning -->
    <div v-if="hasDissonance" class="dissonance-warning">
      ⚠️ Dissonance Detected
      <div class="warning-text">
        {{ getDissonanceText() }}
      </div>
    </div>

    <!-- Aura Preview -->
    <div v-if="resonanceState" class="aura-preview">
      <div class="aura-circle" :style="{ backgroundColor: resonanceState.auraColor }"></div>
      <span class="aura-label">Aura Effect</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface ResonanceState {
  squadID: string
  harmonyScore: number
  tierName: string
  harmonicElements: string[]
  dissonantElements: string[]
  bonuses: {
    damage: number
    defense: number
  }
  bonusDamage: number
  bonusDefense: number
  auraColor: string
  pattern?: {
    PrimaryElements: string[]
    HarmonicPairs: Array<{ Element1: string; Element2: string }>
    DissonantPairs: Array<{ Element1: string; Element2: string }>
  }
}

interface Props {
  resonanceState?: ResonanceState
}

defineProps<Props>()

const displayElements = computed(() => {
  if (!props.resonanceState) return []
  return props.resonanceState.pattern?.PrimaryElements || []
})

const hasDissonance = computed(() => {
  return (
    props.resonanceState &&
    props.resonanceState.pattern?.DissonantPairs &&
    props.resonanceState.pattern.DissonantPairs.length > 0
  )
})

const props = defineProps<Props>()

function getElementColor(element: string): string {
  const colors: Record<string, string> = {
    Fire: '#ff6b6b',
    Water: '#4ecdc4',
    Earth: '#8b7355',
    Air: '#a8d8ff',
    Light: '#ffd700',
    Dark: '#2d3436',
  }
  return colors[element] || '#888'
}

function getElementIcon(element: string): string {
  const icons: Record<string, string> = {
    Fire: '🔥',
    Water: '💧',
    Earth: '🌍',
    Air: '🌬️',
    Light: '☀️',
    Dark: '🌙',
  }
  return icons[element] || '❓'
}

function formatElement(element: string): string {
  return element.charAt(0).toUpperCase() + element.slice(1).toLowerCase()
}

function getElementX(element: string): number {
  const elements = ['Fire', 'Water', 'Earth', 'Air', 'Light', 'Dark']
  const index = elements.indexOf(element)
  return 50 + index * 40 // Spread elements across x-axis
}

function getGaugeColor(score: number): string {
  if (score >= 76) return '#22c55e' // Green for Resonant
  if (score >= 51) return '#3b82f6' // Blue for Harmonized
  if (score >= 26) return '#f59e0b' // Orange for Aligned
  return '#ef4444' // Red for Discordant
}

function getTierColor(tierName: string): string {
  const colors: Record<string, string> = {
    Resonant: '#22c55e',
    Harmonized: '#3b82f6',
    Aligned: '#f59e0b',
    Discordant: '#ef4444',
  }
  return colors[tierName] || '#888'
}

function getTierDescription(tierName: string): string {
  const descriptions: Record<string, string> = {
    Resonant: 'Perfect elemental harmony! Maximum bonuses active.',
    Harmonized: 'Strong elemental synergy with good bonuses.',
    Aligned: 'Some elemental cooperation detected.',
    Discordant: 'Elemental conflict reducing effectiveness.',
  }
  return descriptions[tierName] || ''
}

function formatBonus(multiplier: number): string {
  if (multiplier === 1) return '+0%'
  if (multiplier < 1) return `${((multiplier - 1) * 100).toFixed(0)}%`
  return `+${((multiplier - 1) * 100).toFixed(0)}%`
}

function getBonusColor(multiplier: number): string {
  if (multiplier > 1) return '#22c55e' // Green for bonus
  if (multiplier < 1) return '#ef4444' // Red for penalty
  return '#888' // Gray for neutral
}

function getDissonanceText(): string {
  if (!props.resonanceState) return ''
  const dissonant = props.resonanceState.dissonantElements || []
  return `${dissonant.join(' and ')} are in conflict.`
}
</script>

<style scoped>
.resonance-preview {
  background: linear-gradient(135deg, #1f2937 0%, #111827 100%);
  border: 1px solid #374151;
  border-radius: 12px;
  padding: 24px;
  color: #e5e7eb;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.resonance-header {
  margin-bottom: 24px;
}

.resonance-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #fff;
}

.attunement-circles {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  justify-content: center;
  flex-wrap: wrap;
}

.attunement-circle {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border: 3px solid;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.05);
  transition: all 0.3s ease;
}

.attunement-circle:hover {
  transform: scale(1.1);
  box-shadow: 0 0 20px currentColor;
}

.element-icon {
  font-size: 2rem;
  margin-bottom: 4px;
}

.element-name {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.harmony-lines {
  width: 100%;
  height: 100px;
  margin-bottom: 24px;
  opacity: 0.8;
}

.harmony-gauge {
  margin-bottom: 24px;
}

.gauge-label {
  font-size: 0.875rem;
  font-weight: 600;
  margin-bottom: 8px;
  color: #d1d5db;
}

.gauge-bar-container {
  background: #374151;
  border-radius: 8px;
  height: 24px;
  overflow: hidden;
  margin-bottom: 8px;
}

.gauge-bar-fill {
  height: 100%;
  transition: width 0.5s ease, background-color 0.3s ease;
  border-radius: 8px;
}

.gauge-value {
  text-align: right;
  font-size: 0.875rem;
  color: #9ca3af;
}

.tier-info {
  text-align: center;
  margin-bottom: 24px;
}

.tier-name {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 8px;
  transition: color 0.3s ease;
}

.tier-description {
  font-size: 0.875rem;
  color: #9ca3af;
  line-height: 1.5;
}

.bonus-table {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 24px;
}

.bonus-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  font-size: 0.875rem;
}

.bonus-row:not(:last-child) {
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.bonus-label {
  color: #d1d5db;
}

.bonus-value {
  font-weight: 700;
  transition: color 0.3s ease;
}

.dissonance-warning {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 24px;
  color: #fca5a5;
  font-weight: 600;
  font-size: 0.875rem;
}

.warning-text {
  font-size: 0.8rem;
  color: #f87171;
  margin-top: 4px;
  font-weight: normal;
}

.aura-preview {
  display: flex;
  align-items: center;
  gap: 12px;
  justify-content: center;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}

.aura-circle {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  animation: aura-pulse 2s ease-in-out infinite;
  box-shadow: 0 0 20px currentColor;
}

.aura-label {
  font-size: 0.875rem;
  color: #d1d5db;
  font-weight: 600;
}

@keyframes aura-pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}

@media (max-width: 640px) {
  .resonance-preview {
    padding: 16px;
  }

  .attunement-circles {
    gap: 12px;
  }

  .attunement-circle {
    width: 70px;
    height: 70px;
  }

  .element-icon {
    font-size: 1.5rem;
  }
}
</style>

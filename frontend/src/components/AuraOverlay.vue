<template>
  <div
    v-if="resonance"
    class="aura-overlay"
    :style="{
      borderColor: resonance.auraColor,
      boxShadow: `0 0 30px ${resonance.auraColor}80, 0 0 60px ${resonance.auraColor}40`,
      animation: getAnimationClass(),
    }"
  >
    <!-- Inner glow -->
    <div
      class="aura-glow-inner"
      :style="{
        backgroundColor: `${resonance.auraColor}40`,
        animation: getInnerGlowAnimation(),
      }"
    ></div>

    <!-- Tier indicator badge -->
    <div v-if="showBadge" class="aura-tier-badge" :style="{ backgroundColor: resonance.auraColor }">
      <span class="badge-text">{{ resonance.tierName }}</span>
    </div>

    <!-- Bonus indicators -->
    <div v-if="showBonusIndicators" class="aura-bonuses">
      <span v-if="resonance.bonuses?.damage > 1" class="bonus-icon damage-bonus">
        ⚔️ +{{ formatPercent(resonance.bonuses.damage) }}
      </span>
      <span v-if="resonance.bonuses?.defense > 1" class="bonus-icon defense-bonus">
        🛡️ +{{ formatPercent(resonance.bonuses.defense) }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type ResonanceState = {
  squadID: string
  harmonyScore: number
  tierName: string
  harmonicElements: string[]
  dissonantElements: string[]
  bonuses: {
    damage: number
    defense: number
  }
  auraColor: string
}

type Props = {
  resonance?: ResonanceState
  showBadge?: boolean
  showBonusIndicators?: boolean
  animated?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showBadge: true,
  showBonusIndicators: false,
  animated: true,
})

const getAnimationClass = () => {
  if (!props.animated || !props.resonance) return 'none'

  switch (props.resonance.tierName) {
    case 'Resonant':
      return 'aura-pulse-intense 1.5s ease-in-out infinite'
    case 'Harmonized':
      return 'aura-pulse 2s ease-in-out infinite'
    case 'Aligned':
      return 'aura-pulse-subtle 3s ease-in-out infinite'
    default:
      return 'none'
  }
}

const getInnerGlowAnimation = () => {
  if (!props.animated || !props.resonance) return 'none'

  return props.resonance.harmonyScore >= 76
    ? 'inner-glow-intense 1.5s ease-in-out infinite'
    : props.resonance.harmonyScore >= 51
    ? 'inner-glow 2s ease-in-out infinite'
    : 'none'
}

const formatPercent = (multiplier: number) => {
  return Math.round((multiplier - 1) * 100)
}
</script>

<style scoped>
.aura-overlay {
  position: absolute;
  inset: -15px;
  border: 3px solid;
  border-radius: 8px;
  pointer-events: none;
  transition: all 0.3s ease;
  will-change: transform, box-shadow;
}

.aura-glow-inner {
  position: absolute;
  inset: 3px;
  border-radius: 5px;
  opacity: 0.2;
}

.aura-tier-badge {
  position: absolute;
  top: -15px;
  left: 50%;
  transform: translateX(-50%);
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 700;
  color: #000;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  white-space: nowrap;
  z-index: 10;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.badge-text {
  display: block;
}

.aura-bonuses {
  position: absolute;
  bottom: -20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
  font-size: 0.7rem;
  font-weight: 600;
  white-space: nowrap;
  z-index: 10;
}

.bonus-icon {
  background: rgba(0, 0, 0, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 2px 6px;
  border-radius: 4px;
  color: #fff;
}

.damage-bonus {
  border-color: #ff6b6b;
  color: #ff6b6b;
}

.defense-bonus {
  border-color: #4ecdc4;
  color: #4ecdc4;
}

/* Animations */
@keyframes aura-pulse {
  0%,
  100% {
    opacity: 0.7;
    transform: scale(1);
  }
  50% {
    opacity: 1;
    transform: scale(1.05);
  }
}

@keyframes aura-pulse-intense {
  0%,
  100% {
    opacity: 0.8;
    transform: scale(1);
  }
  50% {
    opacity: 1;
    transform: scale(1.1);
  }
}

@keyframes aura-pulse-subtle {
  0%,
  100% {
    opacity: 0.5;
  }
  50% {
    opacity: 0.8;
  }
}

@keyframes inner-glow {
  0%,
  100% {
    opacity: 0.2;
  }
  50% {
    opacity: 0.4;
  }
}

@keyframes inner-glow-intense {
  0%,
  100% {
    opacity: 0.3;
  }
  50% {
    opacity: 0.6;
  }
}

@media (prefers-reduced-motion: reduce) {
  .aura-overlay,
  .aura-glow-inner {
    animation: none !important;
    opacity: 0.6 !important;
  }
}
</style>

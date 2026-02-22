import { ref, computed } from 'vue'

export interface ResonanceState {
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

export interface FighterResonance {
  fighterId: string
  resonance: ResonanceState
}

/**
 * Composable for managing resonance state in match viewer
 * Handles WebSocket resonance messages and applies aura effects
 */
export function useMatchResonance() {
  const resonanceMap = ref<Map<string, ResonanceState>>(new Map())
  const resonanceInitialized = ref(false)

  /**
   * Process resonance state message from WebSocket
   */
   const onResonanceMessage = (payload: { type: string; resonances: Map<string, ResonanceState> | Record<string, ResonanceState> }) => {
    if (payload.type === 'match.resonance_state' && payload.resonances) {
      resonanceMap.value.clear()
      
      // Handle both Map and object formats
      const resonances = payload.resonances instanceof Map 
        ? payload.resonances 
        : new Map(Object.entries(payload.resonances))
      
      for (const [fighterId, resonance] of resonances) {
        resonanceMap.value.set(fighterId, resonance as ResonanceState)
      }
      
      resonanceInitialized.value = true
    }
  }

  /**
   * Get resonance state for a specific fighter
   */
  const getResonance = (fighterId: string): ResonanceState | null => {
    return resonanceMap.value.get(fighterId) || null
  }

  /**
   * Get all resonances
   */
  const getAllResonances = computed(() => {
    return Array.from(resonanceMap.value.values())
  })

  /**
   * Apply CSS filter for aura effect
   */
  const getAuraStyle = (fighterId: string) => {
    const resonance = getResonance(fighterId)
    if (!resonance) return {}

    const color = resonance.auraColor
    // Convert hex to RGB for drop-shadow filter
    const rgb = hexToRgb(color)
    if (!rgb) return {}

    return {
      filter: `drop-shadow(0 0 20px ${color}) drop-shadow(0 0 10px ${color})`,
      animation: resonance.harmonyScore >= 76 
        ? 'resonance-pulse-intense 1.5s ease-in-out infinite'
        : resonance.harmonyScore >= 51
        ? 'resonance-pulse 2s ease-in-out infinite'
        : 'none'
    }
  }

  /**
   * Convert hex color to RGB
   */
  const hexToRgb = (hex: string): string | null => {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
    if (!result) return null
    return `${parseInt(result[1], 16)}, ${parseInt(result[2], 16)}, ${parseInt(result[3], 16)}`
  }

  /**
   * Get aura color with opacity for background
   */
  const getAuraBackgroundColor = (fighterId: string): string => {
    const resonance = getResonance(fighterId)
    if (!resonance) return 'transparent'
    return addAlphaToHex(resonance.auraColor, 0.15)
  }

  /**
   * Add alpha channel to hex color
   */
  const addAlphaToHex = (hex: string, alpha: number): string => {
    const rgb = hexToRgb(hex)
    if (!rgb) return 'transparent'
    return `rgba(${rgb}, ${alpha})`
  }

  /**
   * Get tier-based color intensity
   */
  const getAuraIntensity = (fighterId: string): number => {
    const resonance = getResonance(fighterId)
    if (!resonance) return 0

    switch (resonance.tierName) {
      case 'Resonant':
        return 1
      case 'Harmonized':
        return 0.7
      case 'Aligned':
        return 0.4
      case 'Discordant':
        return 0
      default:
        return 0
    }
  }

  return {
    resonanceMap,
    resonanceInitialized,
    onResonanceMessage,
    getResonance,
    getAllResonances,
    getAuraStyle,
    getAuraBackgroundColor,
    getAuraIntensity,
    hexToRgb,
  }
}

/**
 * Global styles for aura animations (add to your main CSS)
 */
export const resonanceStyles = `
  @keyframes resonance-pulse {
    0%, 100% {
      filter: drop-shadow(0 0 10px currentColor) drop-shadow(0 0 5px currentColor);
      opacity: 1;
    }
    50% {
      filter: drop-shadow(0 0 20px currentColor) drop-shadow(0 0 10px currentColor);
      opacity: 0.85;
    }
  }

  @keyframes resonance-pulse-intense {
    0%, 100% {
      filter: drop-shadow(0 0 20px currentColor) drop-shadow(0 0 15px currentColor);
      opacity: 1;
    }
    50% {
      filter: drop-shadow(0 0 30px currentColor) drop-shadow(0 0 20px currentColor);
      opacity: 0.9;
    }
  }

  .resonance-aura-container {
    position: relative;
    transition: all 0.3s ease;
  }

  .resonance-aura-glow {
    position: absolute;
    inset: -10px;
    border-radius: 50%;
    opacity: 0.3;
    animation: resonance-pulse 2s ease-in-out infinite;
    pointer-events: none;
  }

  .resonance-glow-intense {
    animation: resonance-pulse-intense 1.5s ease-in-out infinite;
  }
`

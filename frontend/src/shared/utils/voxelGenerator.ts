// Deterministic random number generator based on string seed
const seededRandom = (seed: string) => {
  let h = 0xdeadbeef;
  for (let i = 0; i < seed.length; i++) {
    h = Math.imul(h ^ seed.charCodeAt(i), 2654435761);
  }
  return () => {
    h = Math.imul(h ^ (h >>> 16), 2246822507);
    h = Math.imul(h ^ (h >>> 13), 3266489909);
    return (h >>> 0) / 4294967296;
  };
};

export type FighterAppearance = {
  skin: string;
  armorPrimary: string;
  armorSecondary: string;
  heightScale: number;
  buildScale: number;
  headType: 'standard' | 'helmet';
};

export const generateFighterAppearance = (id: string, attunement?: string | null): FighterAppearance => {
  const rng = seededRandom(id);
  
  // 1. Skin Tones (Humanoid ranges + fantasy)
  const skinTones = ['#fca5a5', '#fcd34d', '#ad5e25', '#5f370e', '#e5e7eb', '#d1d5db', '#9ca3af'];
  const skin = skinTones[Math.floor(rng() * skinTones.length)];

  // 2. Base Armor Colors (Default or overridden by Attunement)
  let armorPrimary = '#3b82f6'; // Default Blue
  let armorSecondary = '#1e3a8a';
  
  // Base pools for non-attuned
  const primaryColors = ['#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6', '#6366f1'];
  const secondaryColors = ['#1e3a8a', '#991b1b', '#065f46', '#b45309', '#5b21b6', '#4338ca'];

  if (attunement === 'Fire') {
     armorPrimary = '#ef4444'; armorSecondary = '#991b1b';
  } else if (attunement === 'Earth') {
     armorPrimary = '#d97706'; armorSecondary = '#78350f';
  } else if (attunement === 'Water') {
     armorPrimary = '#3b82f6'; armorSecondary = '#1e3a8a';
  } else if (attunement === 'Wind') {
     armorPrimary = '#10b981'; armorSecondary = '#065f46';
  } else if (attunement === 'Lightning') {
     armorPrimary = '#eab308'; armorSecondary = '#854d0e';
  } else {
     // Randomize if no attunement
     const idx = Math.floor(rng() * primaryColors.length);
     armorPrimary = primaryColors[idx];
     armorSecondary = secondaryColors[idx];
  }

  // 3. Physical Variations
  const heightScale = 0.9 + rng() * 0.2; // 0.9x to 1.1x height
  const buildScale = 0.85 + rng() * 0.3; // Width

  return { 
    skin, 
    armorPrimary, 
    armorSecondary, 
    heightScale, 
    buildScale,
    headType: rng() > 0.5 ? 'helmet' : 'standard' 
  };
};

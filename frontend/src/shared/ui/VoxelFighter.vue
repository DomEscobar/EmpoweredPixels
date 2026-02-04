<template>
  <canvas ref="canvas" class="w-full h-full object-contain pixelated"></canvas>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { generateFighterAppearance, type FighterAppearance } from '@/shared/utils/voxelGenerator';

const props = defineProps<{
  seed: string;       // Fighter ID
  attunement?: string | null;
  animate?: boolean;
  isAttacking?: boolean;
  facing?: 'left' | 'right' | 'front';
}>();

const canvas = ref<HTMLCanvasElement | null>(null);
let animationHandle: number;

// --- 3D MATH & DRAWING HELPERS ---

const TILE_HEIGHT = 32;

const rotatePt = (x: number, y: number, angle: number) => {
   const cos = Math.cos(angle);
   const sin = Math.sin(angle);
   return {
      x: x * cos - y * sin,
      y: x * sin + y * cos
   };
};

// Simple isometric projection
const toScreen = (x: number, y: number, width: number, height: number, zoom: number) => {
  // Isometric: x goes down-right, y goes down-left
  // We center it
  const isoX = (x - y) * 32; 
  const isoY = (x + y) * 16;
  
  return {
    x: width / 2 + isoX * zoom,
    y: height / 2 + isoY * zoom
  };
};

const shadeColor = (color: string, percent: number) => {
    const f = parseInt(color.slice(1), 16),
          t = percent < 0 ? 0 : 255,
          p = percent < 0 ? percent * -1 : percent,
          R = f >> 16,
          G = f >> 8 & 0x00FF,
          B = f & 0x0000FF;
    return "#" + (0x1000000 + (Math.round((t - R) * p) + R) * 0x10000 + (Math.round((t - G) * p) + G) * 0x100 + (Math.round((t - B) * p) + B)).toString(16).slice(1);
};

const drawIsoCuboid = (
  ctx: CanvasRenderingContext2D,
  cx: number, cy: number, cz: number, 
  w: number, d: number, h: number,   
  angle: number,
  color: string,
  zoom: number,
  canvasW: number, canvasH: number
) => {
  const hw = w / 2;
  const hd = d / 2;
  
  const rawCorners = [
     { x: hw, y: hd },   
     { x: -hw, y: hd },  
     { x: -hw, y: -hd }, 
     { x: hw, y: -hd }   
  ];
  
  const topWorld = rawCorners.map(p => {
     const rot = rotatePt(p.x, p.y, angle);
     return { x: cx + rot.x, y: cy + rot.y, z: cz + h };
  });
  
  const botWorld = rawCorners.map(p => {
     const rot = rotatePt(p.x, p.y, angle);
     return { x: cx + rot.x, y: cy + rot.y, z: cz };
  });

  const project = (p: {x:number, y:number, z:number}) => {
      const s = toScreen(p.x, p.y, canvasW, canvasH, zoom);
      return { x: s.x, y: s.y - p.z * TILE_HEIGHT * zoom };
  };

  const topScreen = topWorld.map(project);
  const botScreen = botWorld.map(project);
  
  ctx.fillStyle = shadeColor(color, 0.1); 
  ctx.beginPath();
  ctx.moveTo(topScreen[0].x, topScreen[0].y);
  for (let i = 1; i < 4; i++) ctx.lineTo(topScreen[i].x, topScreen[i].y);
  ctx.closePath();
  ctx.fill();
  ctx.strokeStyle = shadeColor(color, -0.2);
  ctx.lineWidth = 1;
  ctx.stroke();

  let maxYi = 0;
  let maxY = topScreen[0].y;
  for (let i = 1; i < 4; i++) {
     if (topScreen[i].y > maxY) {
        maxY = topScreen[i].y;
        maxYi = i;
     }
  }
  
  const drawSide = (idx1: number, idx2: number, shade: number) => {
     ctx.fillStyle = shadeColor(color, shade);
     ctx.beginPath();
     ctx.moveTo(topScreen[idx1].x, topScreen[idx1].y);
     ctx.lineTo(topScreen[idx2].x, topScreen[idx2].y);
     ctx.lineTo(botScreen[idx2].x, botScreen[idx2].y);
     ctx.lineTo(botScreen[idx1].x, botScreen[idx1].y);
     ctx.closePath();
     ctx.fill();
     ctx.stroke();
  };
  
  drawSide(maxYi, (maxYi + 1) % 4, -0.3);
  drawSide(maxYi, (maxYi + 3) % 4, -0.1);
};

const drawVoxelCharacter = (
  ctx: CanvasRenderingContext2D,
  appearance: FighterAppearance,
  time: number,
  canvasW: number, 
  canvasH: number,
  isAttacking: boolean
) => {
   const zoom = 1.5; // Fixed zoom for roster view
   
   // Animation States
   const speed = 0.005;
   const breathe = Math.sin(time * speed) * 0.05;
   const limbCycle = props.animate ? Math.sin(time * speed * 2) : 0;
   
   // Angles
   // In roster, we want them facing front-ish (South-East in iso) or customizable
   // Iso 0 angle = Facing South-East
   // We rotate slightly to show volume
   const baseAngle = props.facing === 'left' ? Math.PI / 2 : (props.facing === 'right' ? 0 : Math.PI / 4);
   
   // Limbs
   const legAngle = limbCycle * 0.2;
   let lArmAngle = -limbCycle * 0.2;
   let rArmAngle = limbCycle * 0.2;
   
   if (isAttacking) {
      const attackCycle = Math.sin(time * 0.03); 
      rArmAngle = -Math.PI / 2 + attackCycle * 0.5;
   }
   
   // Dimensions from Appearance
   const bodyW = 0.25 * appearance.buildScale;
   const bodyD = 0.12 * appearance.buildScale;
   const bodyH = 0.3 * appearance.heightScale;
   const headS = 0.25 * (appearance.buildScale * 0.9); // Heads don't scale as much
   const limbW = 0.08 * appearance.buildScale;
   const limbH = 0.3 * appearance.heightScale;
   
   const pz = 0; // Ground level
   const bodyZ = pz + limbH; 
   
   const parts = [];
   
   // Body
   parts.push({
      cx: 0, cy: 0, cz: bodyZ + breathe,
      w: bodyW, d: bodyD, h: bodyH,
      offsetVec: {x:0,y:0,z:0},
      color: appearance.armorPrimary
   });
   
   // Head
   parts.push({
      cx: 0, cy: 0, cz: bodyZ + bodyH + breathe,
      w: headS, d: headS, h: headS,
      offsetVec: {x:0, y:0, z: 0},
      color: appearance.headType === 'helmet' ? appearance.armorSecondary : appearance.skin
   });
   
   // Legs
   const legOffset = 0.06 * appearance.buildScale;
   parts.push({
      cx: 0, cy: legOffset, cz: pz,
      w: limbW, d: limbW, h: limbH,
      offsetVec: { x: Math.sin(legAngle)*0.1, y: 0, z: 0 },
      color: appearance.armorSecondary
   });
   parts.push({
      cx: 0, cy: -legOffset, cz: pz,
      w: limbW, d: limbW, h: limbH,
      offsetVec: { x: Math.sin(-legAngle)*0.1, y: 0, z: 0 },
      color: appearance.armorSecondary
   });
   
   // Arms
   const armZ = bodyZ + bodyH - 0.05;
   const armOffset = 0.16 * appearance.buildScale;
   
   // Left Arm
   parts.push({
      cx: 0, cy: armOffset, cz: armZ - limbH + breathe,
      w: limbW, d: limbW, h: limbH,
      offsetVec: { x: Math.sin(lArmAngle)*0.1, y: 0, z: Math.cos(lArmAngle)*0.1 - 0.1 },
      color: appearance.armorPrimary
   });
   
   // Right Arm (Main Hand)
   const handX = Math.sin(rArmAngle)*0.2;
   const handZ = armZ - limbH + Math.cos(rArmAngle)*0.1 - 0.1 + breathe;
   
   parts.push({
      cx: 0, cy: -armOffset, cz: armZ - limbH + breathe,
      w: limbW, d: limbW, h: limbH,
      offsetVec: { x: handX, y: 0, z: handZ - (armZ - limbH + breathe) }, // relative
      color: appearance.armorPrimary
   });

   // Weapon (Sword) if attacking or just holding it
   // Always show weapon in roster
   const swordLen = 0.4;
   parts.push({
      cx: 0.1, cy: -armOffset - 0.05, cz: handZ, 
      w: swordLen, d: 0.05, h: 0.05,
      offsetVec: { x: handX + 0.2, y: 0, z: 0 },
      color: '#cbd5e1' // Blade
   });

   const transformed = parts.map(p => {
      let lx = p.cx + p.offsetVec.x;
      let ly = p.cy + p.offsetVec.y;
      let lz = p.cz + p.offsetVec.z;
      
      const r = rotatePt(lx, ly, baseAngle);
      
      const wx = r.x;
      const wy = r.y;
      
      const depth = wx + wy + lz * 0.1; 
      
      return { ...p, wx, wy, wz: lz, depth };
   });
   
   transformed.sort((a, b) => a.depth - b.depth);
   
   transformed.forEach(p => {
      drawIsoCuboid(ctx, p.wx, p.wy, p.wz, p.w, p.d, p.h, baseAngle, p.color, zoom, canvasW, canvasH);
   });
};

const render = (time: number) => {
  if (!canvas.value) return;
  const ctx = canvas.value.getContext('2d');
  if (!ctx) return;
  
  // Resize if needed
  const rect = canvas.value.getBoundingClientRect();
  if (canvas.value.width !== rect.width || canvas.value.height !== rect.height) {
     canvas.value.width = rect.width;
     canvas.value.height = rect.height;
  }
  
  // 1. Generate Look
  const appearance = generateFighterAppearance(props.seed, props.attunement);
  
  // 2. Clear & Draw
  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height);
  
  drawVoxelCharacter(
     ctx, 
     appearance,
     time,
     canvas.value.width, 
     canvas.value.height,
     props.isAttacking || false
  );
  
  if (props.animate) {
    animationHandle = requestAnimationFrame(render);
  }
};

onMounted(() => {
  animationHandle = requestAnimationFrame(render);
});

onUnmounted(() => {
  cancelAnimationFrame(animationHandle);
});

watch(() => [props.seed, props.attunement], () => {
   // Re-render immediately if props change
   render(performance.now());
});
</script>

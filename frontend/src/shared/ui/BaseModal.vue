<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4 sm:p-6" :class="containerClass">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-slate-950/80 backdrop-blur-sm" @click="$emit('close')"></div>
        
        <!-- Modal Panel -->
        <Transition
          enter-active-class="transition duration-300 ease-out"
          enter-from-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          enter-to-class="opacity-100 translate-y-0 sm:scale-100"
          leave-active-class="transition duration-200 ease-in"
          leave-from-class="opacity-100 translate-y-0 sm:scale-100"
          leave-to-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
        >
          <div v-if="show" class="relative w-full transform overflow-hidden rounded-2xl border border-slate-800 bg-slate-900 shadow-2xl transition-all flex flex-col" :class="panelClass" :style="panelStyle">
            <!-- Header -->
            <div class="flex items-center justify-between border-b border-slate-800 px-6 py-4 flex-shrink-0">
              <h3 class="text-xl font-bold text-white">
                <slot name="title"></slot>
              </h3>
              <button 
                @click="$emit('close')" 
                class="rounded-lg p-1 text-slate-400 hover:bg-slate-800 hover:text-white transition-colors"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <!-- Body -->
            <div class="flex-1 overflow-y-auto px-6 py-4">
              <slot></slot>
            </div>
            
            <!-- Footer -->
            <div v-if="$slots.footer" class="border-t border-slate-800 px-6 py-4 flex-shrink-0">
              <slot name="footer"></slot>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  show: boolean;
  variant?: 'default' | 'bottom-sheet' | 'full-screen';
  containerClass?: string;
  panelClass?: string;
}>();

const containerClasses = computed(() => {
  const base = '';
  const variantClasses = {
    'bottom-sheet': 'items-end',
    'full-screen': 'p-0',
    'default': ''
  };
  return [base, props.containerClass || ''].filter(Boolean).join(' ');
});

const panelClasses = computed(() => {
  const base = 'relative transform overflow-hidden rounded-2xl border border-slate-800 bg-slate-900 shadow-2xl transition-all flex flex-col';
  const variantClasses = {
    'bottom-sheet': 'max-h-[90vh] w-full rounded-t-2xl',
    'full-screen': 'max-w-full h-full rounded-none',
    'default': 'max-w-2xl max-h-[90vh]'
  };
  return [base, variantClasses[props.variant || 'default'], props.panelClass || ''].filter(Boolean).join(' ');
});

const panelStyle = computed(() => {
  if (props.variant === 'full-screen') {
    return { height: '100%' };
  }
  return {};
});

defineEmits<{
  (e: 'close'): void;
}>();
</script>

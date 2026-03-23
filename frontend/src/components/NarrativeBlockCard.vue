<script setup lang="ts">
import type { NarrativeBlock } from '@/types/character'

const props = defineProps<{
  block: NarrativeBlock
  locked: boolean
  title: string
  isSecret?: boolean
}>()

const emit = defineEmits<{
  refresh: []
  toggleLock: []
}>()
</script>

<template>
  <!-- Secret: dark treatment -->
  <div v-if="props.isSecret" class="bg-surface-container-low p-6 relative group">
    <div class="absolute -top-3 -right-3 flex items-center gap-1">
      <button @click="emit('refresh')" class="bg-surface-container-lowest text-outline p-1 shadow-sm hover:text-primary transition-colors">
        <span class="material-symbols-outlined text-base">refresh</span>
      </button>
    </div>
    <p class="text-xs font-bold uppercase tracking-widest text-primary-container mb-2">{{ props.title }}</p>
    <p class="font-headline italic text-lg text-primary leading-relaxed">{{ props.block.content }}</p>
  </div>

  <!-- Standard narrative block -->
  <div v-else class="space-y-2 relative group">
    <div class="flex justify-between items-center">
      <h4 class="font-headline text-xl font-bold text-on-surface">{{ props.title }}</h4>
      <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
        <button @click="emit('toggleLock')" class="text-outline hover:text-primary transition-colors">
          <span class="material-symbols-outlined text-sm">{{ props.locked ? 'lock' : 'lock_open' }}</span>
        </button>
        <button @click="emit('refresh')" class="text-outline hover:text-primary transition-colors">
          <span class="material-symbols-outlined text-sm">refresh</span>
        </button>
      </div>
    </div>
    <p class="font-body text-on-surface-variant leading-relaxed text-sm">{{ props.block.content }}</p>
  </div>
</template>

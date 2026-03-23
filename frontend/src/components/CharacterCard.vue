<script setup lang="ts">
import type { Character } from '@/types/character'

const props = defineProps<{
  character: Character
}>()

const emit = defineEmits<{
  select: [id: string]
}>()

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('es-AR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1)
}
</script>

<template>
  <article
    class="bg-surface-container rounded-lg p-5 cursor-pointer hover:bg-surface-container-high transition-colors border border-outline-variant/20 flex flex-col gap-3"
    @click="emit('select', props.character.id)"
  >
    <div class="flex justify-between items-start gap-2">
      <h3 class="font-headline text-on-surface text-lg leading-snug">{{ props.character.name }}</h3>
      <span class="text-[10px] font-label font-bold uppercase tracking-widest text-secondary whitespace-nowrap mt-1">
        Nv. {{ props.character.level }}
      </span>
    </div>

    <div class="flex flex-wrap gap-2">
      <span class="bg-primary/10 text-primary text-[11px] font-label font-semibold uppercase tracking-widest px-2 py-0.5 rounded">
        {{ capitalize(props.character.class) }}
      </span>
      <span class="bg-secondary/10 text-secondary text-[11px] font-label font-semibold uppercase tracking-widest px-2 py-0.5 rounded">
        {{ capitalize(props.character.species) }}
      </span>
    </div>

    <div class="flex gap-4 text-xs font-label text-on-surface-variant">
      <span><span class="font-bold text-on-surface">{{ props.character.derived.hp }}</span> PV</span>
      <span><span class="font-bold text-on-surface">{{ props.character.derived.ac }}</span> CA</span>
    </div>

    <p class="text-[11px] text-outline mt-auto">{{ formatDate(props.character.createdAt) }}</p>
  </article>
</template>

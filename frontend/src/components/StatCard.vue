<script setup lang="ts">
const props = defineProps<{
  stat: 'STR' | 'DEX' | 'CON' | 'INT' | 'WIS' | 'CHA'
  value: number
  modifier: number
  isPrimary?: boolean
}>()

const statLabels: Record<string, string> = {
  STR: 'FUE', DEX: 'DES', CON: 'CON', INT: 'INT', WIS: 'SAB', CHA: 'CAR',
}

function formatModifier(mod: number): string {
  return mod >= 0 ? `+${mod}` : `${mod}`
}
</script>

<template>
  <div
    :class="[
      'p-4 flex flex-col items-center transition-colors group relative',
      props.isPrimary
        ? 'bg-surface-container-low border-2 border-primary-container/20'
        : 'bg-surface-container-low hover:bg-surface-container',
    ]"
  >
    <span
      :class="[
        'text-[10px] font-bold uppercase tracking-widest',
        props.isPrimary ? 'text-primary' : 'text-secondary group-hover:text-primary',
      ]"
    >
      {{ statLabels[props.stat] }}
    </span>
    <span
      :class="[
        'text-3xl font-headline font-bold',
        props.isPrimary ? 'text-primary' : 'text-on-surface',
      ]"
    >
      {{ props.value }}
    </span>
    <div
      :class="[
        'mt-1 px-2 py-0.5 rounded-full text-[10px] font-bold font-body',
        props.isPrimary ? 'bg-primary-container text-on-primary' : 'bg-surface-variant text-on-surface',
      ]"
    >
      {{ formatModifier(props.modifier) }}
    </div>
  </div>
</template>

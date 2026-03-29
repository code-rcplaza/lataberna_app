<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useGeneratorHistoryStore } from '@/stores/useGeneratorHistoryStore'
import { useCharacterStore } from '@/stores/useCharacterStore'
import { useAuthStore } from '@/stores/useAuthStore'

const router = useRouter()
const historyStore = useGeneratorHistoryStore()
const characterStore = useCharacterStore()
const authStore = useAuthStore()

const speciesLabels: Record<string, string> = {
  human: 'Humano', elf: 'Elfo', dwarf: 'Enano', halfling: 'Mediano', gnome: 'Gnomo',
  'half-elf': 'Semielfo', 'half-orc': 'Semiorco', tiefling: 'Tiefling', dragonborn: 'Dragonborn',
}
const subSpeciesLabels: Record<string, string> = {
  'high-elf': 'Alto Elfo', 'wood-elf': 'Elfo del Bosque', 'drow': 'Drow',
  'hill-dwarf': 'Enano de las Colinas', 'mountain-dwarf': 'Enano de la Montaña',
  'lightfoot': 'Pie Ligero', 'stout': 'Robusto',
  'forest-gnome': 'Gnomo del Bosque', 'rock-gnome': 'Gnomo de Roca',
  'tiefling-infernal': 'Linaje Infernal', 'tiefling-virtue': 'Linaje Virtud',
}
const classLabels: Record<string, string> = {
  barbarian: 'Bárbaro', bard: 'Bardo', cleric: 'Clérigo', druid: 'Druida', fighter: 'Guerrero',
  monk: 'Monje', paladin: 'Paladín', ranger: 'Explorador', rogue: 'Pícaro', sorcerer: 'Hechicero',
  warlock: 'Brujo', wizard: 'Mago', artificer: 'Artificiero',
}

function speciesLabel(species: string, subSpecies?: string | null): string {
  if (subSpecies && subSpeciesLabels[subSpecies]) return subSpeciesLabels[subSpecies]
  return speciesLabels[species] ?? species
}

function loadOrNavigate(index: number) {
  const entry = historyStore.recent[index]
  if (!entry) return

  historyStore.rotateIn(index, characterStore.current, characterStore.isSaved)

  if (entry.isSaved) {
    router.push(`/biblioteca/${entry.character.id}`)
  } else {
    characterStore.setCharacter(entry.character)
  }
}
</script>

<template>
  <div v-if="authStore.isAuthenticated && historyStore.recent.length > 0" class="space-y-3 pt-4 border-t border-outline-variant/20">
    <p class="text-[10px] font-bold uppercase tracking-widest text-secondary">Últimos generados</p>
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
      <button
        v-for="(entry, i) in historyStore.recent"
        :key="entry.character.id"
        @click="loadOrNavigate(i)"
        class="group text-left bg-surface-container-low border border-outline-variant/20 px-4 py-3 hover:border-primary/40 hover:bg-surface-container transition-colors"
      >
        <div class="flex items-start justify-between gap-2">
          <p class="font-headline text-sm font-bold text-on-surface truncate leading-tight">
            {{ entry.character.name }}
          </p>
          <span
            v-if="entry.isSaved"
            class="shrink-0 text-[9px] font-bold uppercase tracking-widest text-primary bg-primary/10 px-1.5 py-0.5"
          >
            Guardado
          </span>
        </div>
        <p class="text-[11px] text-secondary mt-1 truncate">
          {{ speciesLabel(entry.character.species, entry.character.subSpecies) }}
          · {{ classLabels[entry.character.class] ?? entry.character.class }}
        </p>
        <p class="text-[10px] text-outline mt-1">
          {{ entry.character.derived.hp }} HP · CA {{ entry.character.derived.ac }}
        </p>
      </button>
    </div>
  </div>
</template>

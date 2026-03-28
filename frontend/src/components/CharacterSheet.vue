<script setup lang="ts">
import { computed } from 'vue'
import StatCard from './StatCard.vue'
import NarrativeBlockCard from './NarrativeBlockCard.vue'
import type { Character, GeneratorLocks } from '@/types/character'
import { useCharacterStore } from '@/stores/useCharacterStore'
import { useCharacterAPI } from '@/composables/useCharacterAPI'

const props = defineProps<{
  character: Character | null
  locks: GeneratorLocks
}>()

const store = useCharacterStore()
const { generate, regenerateField } = useCharacterAPI()

// Primary stat per class
const primaryStat: Record<string, string> = {
  barbarian: 'STR',
  bard:      'CHA',
  cleric:    'WIS',
  druid:     'WIS',
  fighter:   'STR',
  monk:      'DEX',
  paladin:   'CHA',
  ranger:    'DEX',
  rogue:     'DEX',
  sorcerer:  'CHA',
  warlock:   'CHA',
  wizard:    'INT',
  artificer: 'INT',
}

const characterPrimary = computed(() =>
  props.character ? (primaryStat[props.character.class] ?? '') : '',
)

// Spanish labels
const speciesLabels: Record<string, string> = {
  human: 'Humano', elf: 'Elfo', dwarf: 'Enano', halfling: 'Mediano', gnome: 'Gnomo',
  'half-elf': 'Semielfo', 'half-orc': 'Semiorco', tiefling: 'Tiefling', dragonborn: 'Dragonborn',
}
const subSpeciesLabels: Record<string, string> = {
  'high-elf':           'Alto Elfo',
  'wood-elf':           'Elfo del Bosque',
  'drow':               'Drow',
  'hill-dwarf':         'Enano de las Colinas',
  'mountain-dwarf':     'Enano de la Montaña',
  'lightfoot':          'Pie Ligero',
  'stout':              'Robusto',
  'forest-gnome':       'Gnomo del Bosque',
  'rock-gnome':         'Gnomo de Roca',
  'tiefling-infernal':  'Linaje Infernal',
  'tiefling-virtue':    'Linaje Virtud',
}
const classLabels: Record<string, string> = {
  barbarian: 'Bárbaro', bard: 'Bardo', cleric: 'Clérigo', druid: 'Druida', fighter: 'Guerrero',
  monk: 'Monje', paladin: 'Paladín', ranger: 'Montaraz', rogue: 'Pícaro', sorcerer: 'Hechicero',
  warlock: 'Brujo', wizard: 'Mago', artificer: 'Artificiero',
}

const statKeys = ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA'] as const

function proficiencyBonus(level: number): string {
  const bonus = Math.ceil(level / 4) + 1
  return `+${bonus}`
}

</script>

<template>
  <!-- Empty state -->
  <div
    v-if="!character"
    class="flex flex-col items-center justify-center min-h-[600px] text-center space-y-4 opacity-40"
  >
    <span class="material-symbols-outlined text-6xl text-outline">auto_stories</span>
    <p class="font-headline text-2xl text-on-surface-variant italic">Genera tu primer personaje</p>
    <p class="font-body text-sm text-outline max-w-xs">
      Configura los parámetros a la izquierda y haz clic en "¡Sorpréndeme!" para dar vida a un aventurero.
    </p>
    <button
      @click="generate()"
      :disabled="store.isLoading"
      class="mt-4 bg-primary text-on-primary px-6 py-3 font-label font-bold uppercase tracking-widest text-xs hover:bg-primary-container hover:text-on-primary-container transition-colors disabled:opacity-60"
    >
      {{ store.isLoading ? 'Generando…' : '¡Sorpréndeme!' }}
    </button>
  </div>

  <!-- Character sheet -->
  <div v-else class="space-y-8">

    <!-- ── Header ── -->
    <div class="flex items-start justify-between">
      <div>
        <div class="flex items-center gap-3">
          <h1 class="font-headline text-4xl font-bold text-on-surface">{{ character.name }}</h1>
          <button
            @click="regenerateField('name')"
            :disabled="store.isLoading || store.locks.species"
            class="text-outline hover:text-primary transition-colors disabled:opacity-40"
            title="Regenerar nombre"
          >
            <span class="material-symbols-outlined">refresh</span>
          </button>
        </div>
        <p class="font-body text-secondary mt-1">
          {{ speciesLabels[character.species] ?? character.species }}
          <span v-if="character.subSpecies"> · {{ subSpeciesLabels[character.subSpecies] ?? character.subSpecies }}</span>
          · {{ classLabels[character.class] ?? character.class }}
          · Nivel {{ character.level }}
        </p>
      </div>
      <div class="flex flex-col items-end gap-2">
        <span class="text-[10px] font-bold uppercase tracking-widest text-secondary bg-surface-container px-3 py-1">
          {{ character.ruleset }}
        </span>
        <span v-if="store.isSaved" class="text-[10px] font-bold uppercase tracking-widest text-on-primary bg-primary px-3 py-1">
          Guardado
        </span>
      </div>
    </div>

    <!-- ── Atributos ── -->
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <h3 class="font-headline text-lg font-bold text-on-surface">Atributos</h3>
        <div class="flex items-center gap-3">
          <span class="text-[10px] font-bold uppercase tracking-widest text-secondary bg-surface-container px-2 py-1">
            Competencia {{ proficiencyBonus(character.level) }}
          </span>
          <button
            @click="regenerateField('stats')"
            :disabled="store.isLoading || store.locks.stats"
            class="text-outline hover:text-primary transition-colors disabled:opacity-40"
            title="Retirar estadísticas"
          >
            <span class="material-symbols-outlined text-sm">casino</span>
          </button>
        </div>
      </div>
      <div class="grid grid-cols-3 sm:grid-cols-6 gap-2">
        <StatCard
          v-for="stat in statKeys"
          :key="stat"
          :stat="stat"
          :value="character.finalStats[stat]"
          :modifier="character.modifiers[stat]"
          :is-primary="characterPrimary === stat"
        />
      </div>
    </div>

    <!-- ── Combat Stats ── -->
    <div class="grid grid-cols-2 gap-4">
      <div class="bg-surface-container-low p-5 flex flex-col items-center justify-center space-y-1">
        <span class="text-[10px] font-bold uppercase tracking-widest text-secondary">Puntos de Golpe</span>
        <span class="font-headline text-5xl font-bold text-primary">{{ character.derived.hp }}</span>
        <span class="text-xs text-outline font-body">HP máximos</span>
      </div>
      <div class="bg-surface-container-low p-5 flex flex-col items-center justify-center space-y-1">
        <span class="text-[10px] font-bold uppercase tracking-widest text-secondary">Clase de Armadura</span>
        <span class="font-headline text-5xl font-bold text-primary">{{ character.derived.ac }}</span>
        <span class="text-xs text-outline font-body">CA base</span>
      </div>
    </div>

    <!-- ── Narrativa ── -->
    <div class="space-y-6">
      <h3 class="font-headline text-lg font-bold text-on-surface">Narrativa</h3>

      <NarrativeBlockCard
        :block="character.background"
        :locked="store.locks.background"
        title="Origen"
        @refresh="regenerateField('background')"
        @toggle-lock="store.toggleLock('background')"
      />

      <div class="border-t border-outline-variant/20" />

      <NarrativeBlockCard
        :block="character.motivation"
        :locked="store.locks.motivation"
        title="Motivación"
        @refresh="regenerateField('motivation')"
        @toggle-lock="store.toggleLock('motivation')"
      />

      <div class="border-t border-outline-variant/20" />

      <NarrativeBlockCard
        :block="character.secret"
        :locked="store.locks.secret"
        title="Secreto"
        :is-secret="true"
        @refresh="regenerateField('secret')"
        @toggle-lock="store.toggleLock('secret')"
      />
    </div>

  </div>
</template>

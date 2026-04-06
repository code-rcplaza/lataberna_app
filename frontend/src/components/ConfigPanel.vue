<script setup lang="ts">
import { computed } from 'vue'
import LockButton from './LockButton.vue'
import { useCharacterStore } from '@/stores/useCharacterStore'
import { useCharacterAPI } from '@/composables/useCharacterAPI'
import type { Species, SubSpecies, Class } from '@/types/character'

const store = useCharacterStore()
const { generate } = useCharacterAPI()

// ── Species ──────────────────────────────────────────────────────────────────
const speciesOptions: Array<{ value: Species | 'random'; label: string }> = [
  { value: 'random',     label: 'Aleatorio' },
  { value: 'human',      label: 'Humano' },
  { value: 'elf',        label: 'Elfo' },
  { value: 'dwarf',      label: 'Enano' },
  { value: 'halfling',   label: 'Mediano' },
  { value: 'gnome',      label: 'Gnomo' },
  { value: 'half-elf',   label: 'Semielfo' },
  { value: 'half-orc',   label: 'Semiorco' },
  { value: 'tiefling',   label: 'Tiefling' },
  { value: 'dragonborn', label: 'Dragonborn' },
]

// ── SubSpecies (filtered by species selection) ────────────────────────────────
const allSubSpecies: Array<{ value: SubSpecies | 'random'; label: string; species: Species[] }> = [
  { value: 'random',           label: 'Aleatorio',         species: ['human','elf','dwarf','halfling','gnome','half-elf','half-orc','tiefling','dragonborn'] },
  { value: 'high-elf',         label: 'Alto Elfo',            species: ['elf'] },
  { value: 'wood-elf',         label: 'Elfo del Bosque',     species: ['elf'] },
  { value: 'drow',             label: 'Drow',                species: ['elf'] },
  { value: 'hill-dwarf',       label: 'Enano de las Colinas',species: ['dwarf'] },
  { value: 'mountain-dwarf',   label: 'Enano de la Montaña', species: ['dwarf'] },
  { value: 'lightfoot',        label: 'Pie Ligero',           species: ['halfling'] },
  { value: 'stout',            label: 'Robusto',              species: ['halfling'] },
  { value: 'forest-gnome',     label: 'Gnomo del Bosque',    species: ['gnome'] },
  { value: 'rock-gnome',       label: 'Gnomo de Roca',       species: ['gnome'] },
  { value: 'tiefling-infernal',label: 'Linaje Infernal',     species: ['tiefling'] },
  { value: 'tiefling-virtue',  label: 'Linaje Virtud',       species: ['tiefling'] },
]

const subSpeciesOptions = computed(() => {
  const selected = store.input.species
  if (!selected || selected === 'random') return allSubSpecies
  return allSubSpecies.filter(s => s.value === 'random' || s.species.includes(selected as Species))
})

// ── Classes ───────────────────────────────────────────────────────────────────
const classOptions: Array<{ value: Class | 'random'; label: string }> = [
  { value: 'random',     label: 'Aleatorio' },
  { value: 'barbarian',  label: 'Bárbaro' },
  { value: 'bard',       label: 'Bardo' },
  { value: 'cleric',     label: 'Clérigo' },
  { value: 'druid',      label: 'Druida' },
  { value: 'fighter',    label: 'Guerrero' },
  { value: 'monk',       label: 'Monje' },
  { value: 'paladin',    label: 'Paladín' },
  { value: 'ranger',     label: 'Explorador' },
  { value: 'rogue',      label: 'Pícaro' },
  { value: 'sorcerer',   label: 'Hechicero' },
  { value: 'warlock',    label: 'Brujo' },
  { value: 'wizard',     label: 'Mago' },
  { value: 'artificer',  label: 'Artificiero' },
]

// ── Alignment ─────────────────────────────────────────────────────────────────
const alignmentOptions = [
  { value: 'random',          label: 'Aleatorio' },
  { value: 'legal-bueno',     label: 'Legal Bueno' },
  { value: 'legal-neutral',   label: 'Legal Neutral' },
  { value: 'legal-malvado',   label: 'Legal Malvado' },
  { value: 'neutral-bueno',   label: 'Neutral Bueno' },
  { value: 'neutral',         label: 'Neutral' },
  { value: 'neutral-malvado', label: 'Neutral Malvado' },
  { value: 'caotico-bueno',   label: 'Caótico Bueno' },
  { value: 'caotico-neutral', label: 'Caótico Neutral' },
  { value: 'caotico-malvado', label: 'Caótico Malvado' },
]

// ── Dirty check ───────────────────────────────────────────────────────────────
const isDirty = computed(() => {
  const i = store.input
  const l = store.locks
  const inputChanged =
    (i.species   && i.species   !== 'random') ||
    (i.subSpecies && i.subSpecies !== 'random') ||
    (i.class     && i.class     !== 'random') ||
    (i.gender    && i.gender    !== 'random') ||
    (i.alignment && i.alignment !== 'random') ||
    i.seed != null
  const locksChanged = Object.values(l).some(Boolean)
  return inputChanged || locksChanged
})

// ── Helper ─────────────────────────────────────────────────────────────────────
function fieldClass(locked: boolean) {
  return locked
    ? 'w-full bg-surface-dim text-on-surface-variant border border-outline-variant/30 px-3 py-2 pr-9 text-sm font-body cursor-not-allowed opacity-70 focus:outline-none appearance-none'
    : 'w-full bg-surface-container-lowest text-on-surface border border-outline-variant/30 px-3 py-2 pr-9 text-sm font-body focus:outline-none focus:border-primary transition-colors appearance-none'
}
</script>

<template>
  <div class="bg-surface-container-low p-6 space-y-6">
    <!-- Panel header -->
    <div class="border-b border-outline-variant/30 pb-4">
      <h2 class="font-headline text-2xl font-bold text-on-surface">Escritorio de Encarnación</h2>
      <p class="text-xs text-secondary mt-1 font-body">Configura los parámetros del personaje antes de generar</p>
    </div>

    <!-- Especie -->
    <div class="space-y-1">
      <div class="flex items-center justify-between">
        <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Especie</label>
        <LockButton :locked="store.locks.species" @toggle="store.toggleLock('species')" />
      </div>
      <div class="relative">
        <select
          v-model="store.input.species"
          :disabled="store.locks.species"
          :class="fieldClass(store.locks.species)"
        >
          <option v-for="opt in speciesOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <span class="material-symbols-outlined absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-secondary text-base">expand_more</span>
      </div>
    </div>

    <!-- Sub-especie / Linaje -->
    <div class="space-y-1">
      <div class="flex items-center justify-between">
        <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Sub-especie / Linaje</label>
        <LockButton :locked="store.locks.subSpecies" @toggle="store.toggleLock('subSpecies')" />
      </div>
      <div class="relative">
        <select
          v-model="store.input.subSpecies"
          :disabled="store.locks.subSpecies"
          :class="fieldClass(store.locks.subSpecies)"
        >
          <option v-for="opt in subSpeciesOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <span class="material-symbols-outlined absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-secondary text-base">expand_more</span>
      </div>
    </div>

    <!-- Vocación -->
    <div class="space-y-1">
      <div class="flex items-center justify-between">
        <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Vocación</label>
        <LockButton :locked="store.locks.class" @toggle="store.toggleLock('class')" />
      </div>
      <div class="relative">
        <select
          v-model="store.input.class"
          :disabled="store.locks.class"
          :class="fieldClass(store.locks.class)"
        >
          <option v-for="opt in classOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <span class="material-symbols-outlined absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-secondary text-base">expand_more</span>
      </div>
    </div>

    <!-- Género + Semilla (2-col) -->
    <div class="grid grid-cols-2 gap-4 items-end">
      <!-- Género -->
      <div class="space-y-1">
        <div class="flex items-center justify-between">
          <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Género</label>
          <LockButton :locked="store.locks.gender" @toggle="store.toggleLock('gender')" />
        </div>
        <div class="relative">
          <select
            v-model="store.input.gender"
            :disabled="store.locks.gender"
            :class="fieldClass(store.locks.gender)"
          >
            <option value="random">Aleatorio</option>
            <option value="male">Masculino</option>
            <option value="female">Femenino</option>
          </select>
          <span class="material-symbols-outlined absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-secondary text-base">expand_more</span>
        </div>
      </div>

      <!-- Semilla -->
      <div class="space-y-1">
        <div class="flex items-center justify-between">
          <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Semilla</label>
          <span class="w-6 h-6"></span>
        </div>
        <input
          v-model.number="store.input.seed"
          type="number"
          placeholder="Aleatoria"
          class="w-full bg-surface-container-lowest text-on-surface border border-outline-variant/30 px-3 py-2 text-sm font-body focus:outline-none focus:border-primary transition-colors"
        />
      </div>
    </div>

    <!-- Alineamiento (cosmético) -->
    <div class="space-y-1">
      <div class="flex items-center justify-between">
        <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Alineamiento</label>
        <LockButton :locked="store.locks.alignment" @toggle="store.toggleLock('alignment')" />
      </div>
      <div class="relative">
        <select
          v-model="store.input.alignment"
          :disabled="store.locks.alignment"
          :class="fieldClass(store.locks.alignment)"
        >
          <option v-for="opt in alignmentOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <span class="material-symbols-outlined absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-secondary text-base">expand_more</span>
      </div>
      <p class="text-[10px] text-outline font-body italic">Solo cosmético — no afecta la generación</p>
    </div>

    <!-- Actions -->
    <div class="flex flex-col gap-2">
      <button
        @click="generate()"
        :disabled="store.isLoading"
        class="w-full bg-primary text-on-primary font-label font-bold uppercase tracking-widest py-4 text-sm hover:bg-primary-container hover:text-on-primary-container transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
      >
        <span v-if="store.isLoading" class="material-symbols-outlined text-base animate-spin">autorenew</span>
        <span>{{ store.isLoading ? 'Generando…' : '¡Sorpréndeme!' }}</span>
      </button>
      <button
        v-if="isDirty"
        @click="store.reset()"
        :disabled="store.isLoading"
        class="w-full border border-outline-variant text-secondary font-label font-bold uppercase tracking-widest py-2 text-xs hover:border-outline hover:text-on-surface transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
      >
        <span class="material-symbols-outlined text-sm">close</span>
        Limpiar
      </button>
    </div>

    <!-- Flavor note -->
    <div class="border-l-2 border-primary-container/40 pl-4 py-2 bg-surface-container rounded-sm">
      <p class="text-xs text-on-surface-variant font-body italic leading-relaxed">
        "Cada alma nace del caos y el azar. Los dados dicen quién eres — tú decides en qué te conviertes."
      </p>
    </div>
  </div>
</template>

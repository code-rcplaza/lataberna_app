<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLibraryStore } from '@/stores/useLibraryStore'
import { useLibraryAPI } from '@/composables/useCharacterAPI'
import ConfirmModal from '@/components/ConfirmModal.vue'
import type { EditCharacterInput } from '@/types/character'

const route = useRoute()
const router = useRouter()
const libraryStore = useLibraryStore()
const { getCharacter, deleteCharacter, editCharacter } = useLibraryAPI()

const deleteError = ref<string | null>(null)
const editError = ref<string | null>(null)
const editSuccess = ref(false)
const showDeleteModal = ref(false)

// Name edit state
const editingName = ref(false)
const nameInput = ref('')
const nameValidationError = ref<string | null>(null)

// Narrative edit state
type NarrativeKey = 'background' | 'motivation' | 'secret'
const editingBlock = ref<NarrativeKey | null>(null)
const narrativeInput = ref('')
const narrativeValidationError = ref<string | null>(null)
const narrativeSuccess = ref<NarrativeKey | null>(null)

const statKeys = ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA'] as const

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
const statLabels: Record<string, string> = {
  STR: 'FUE', DEX: 'DES', CON: 'CON', INT: 'INT', WIS: 'SAB', CHA: 'CAR',
}
const narrativeTitles: Record<string, string> = {
  background: 'Origen', motivation: 'Motivación', secret: 'Secreto',
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('es-AR', {
    day: '2-digit', month: 'long', year: 'numeric',
  })
}

function modifierDisplay(mod: number): string {
  return mod >= 0 ? `+${mod}` : `${mod}`
}

onMounted(async () => {
  const id = route.params.id as string
  libraryStore.setSelected(null)
  try {
    await getCharacter(id)
  } catch {
    // error is set in libraryStore.error
  }
})

async function confirmDelete() {
  const character = libraryStore.selected
  if (!character) return
  showDeleteModal.value = false
  deleteError.value = null
  try {
    await deleteCharacter(character.id)
    router.push('/biblioteca')
  } catch (err) {
    deleteError.value = err instanceof Error ? err.message : 'Error al eliminar el personaje'
  }
}

function startEditName() {
  if (!libraryStore.selected) return
  nameInput.value = libraryStore.selected.name
  nameValidationError.value = null
  editError.value = null
  editSuccess.value = false
  editingName.value = true
}

function cancelEditName() {
  editingName.value = false
  nameInput.value = ''
  nameValidationError.value = null
}

function startEditNarrative(key: NarrativeKey) {
  if (!libraryStore.selected) return
  narrativeInput.value = libraryStore.selected[key].content
  narrativeValidationError.value = null
  editError.value = null
  narrativeSuccess.value = null
  editingBlock.value = key
}

function cancelEditNarrative() {
  editingBlock.value = null
  narrativeInput.value = ''
  narrativeValidationError.value = null
}

async function submitEditNarrative() {
  const block = editingBlock.value
  if (!block || !libraryStore.selected) return
  const trimmed = narrativeInput.value.trim()
  if (!trimmed) {
    narrativeValidationError.value = 'El contenido no puede estar vacío.'
    return
  }
  narrativeValidationError.value = null
  editError.value = null
  try {
    await editCharacter(libraryStore.selected.id, { [block]: trimmed })
    editingBlock.value = null
    narrativeSuccess.value = block
    setTimeout(() => { narrativeSuccess.value = null }, 3000)
  } catch (err) {
    editError.value = err instanceof Error ? err.message : 'Error al editar el personaje'
  }
}

async function submitEditName() {
  const trimmed = nameInput.value.trim()
  if (!trimmed) {
    nameValidationError.value = 'El nombre no puede estar vacío.'
    return
  }
  nameValidationError.value = null
  editError.value = null
  const patch: EditCharacterInput = { name: trimmed }
  try {
    await editCharacter(libraryStore.selected!.id, patch)
    editingName.value = false
    editSuccess.value = true
    setTimeout(() => { editSuccess.value = false }, 3000)
  } catch (err) {
    editError.value = err instanceof Error ? err.message : 'Error al editar el personaje'
  }
}
</script>

<template>
  <div class="max-w-[900px] mx-auto px-8 pb-12">

    <!-- Back link -->
    <div class="mb-6">
      <RouterLink
        to="/biblioteca"
        class="flex items-center gap-1 text-secondary text-xs font-label font-bold uppercase tracking-widest hover:text-primary transition-colors"
      >
        <span class="material-symbols-outlined text-sm">arrow_back</span>
        Biblioteca
      </RouterLink>
    </div>

    <!-- Loading state -->
    <div v-if="libraryStore.isLoading" class="flex flex-col items-center justify-center py-24 gap-4">
      <span class="material-symbols-outlined text-primary text-5xl animate-spin">refresh</span>
      <p class="text-outline text-sm font-label">Cargando personaje…</p>
    </div>

    <!-- Error / not found state -->
    <div
      v-else-if="libraryStore.error || !libraryStore.selected"
      class="flex flex-col items-center justify-center py-24 gap-4"
    >
      <span class="material-symbols-outlined text-outline text-6xl">person_off</span>
      <p class="font-headline text-on-surface text-xl">Personaje no encontrado</p>
      <p class="text-outline text-sm">
        {{ libraryStore.error ?? 'Este personaje no existe o no te pertenece.' }}
      </p>
      <RouterLink
        to="/biblioteca"
        class="mt-2 text-xs font-bold uppercase tracking-widest text-primary hover:text-primary-container transition-colors"
      >
        Volver a la Biblioteca
      </RouterLink>
    </div>

    <!-- Character detail -->
    <div v-else class="space-y-8">
      <!-- Header -->
      <div class="flex items-start justify-between gap-4">
        <div class="flex-1">
          <!-- Name + edit -->
          <div v-if="!editingName" class="flex items-center gap-3">
            <h1 class="font-headline text-4xl font-bold text-on-surface">{{ libraryStore.selected.name }}</h1>
            <button
              @click="startEditName"
              class="text-outline hover:text-primary transition-colors"
              title="Editar nombre"
            >
              <span class="material-symbols-outlined text-base">edit</span>
            </button>
          </div>
          <div v-else class="flex items-center gap-2">
            <input
              v-model="nameInput"
              type="text"
              class="font-headline text-2xl bg-surface-container border border-outline-variant text-on-surface px-3 py-1 focus:outline-none focus:border-primary"
              @keyup.enter="submitEditName"
              @keyup.escape="cancelEditName"
            />
            <button
              @click="submitEditName"
              :disabled="libraryStore.isLoading"
              class="text-primary hover:text-primary-container transition-colors disabled:opacity-50"
              title="Guardar"
            >
              <span class="material-symbols-outlined">check</span>
            </button>
            <button
              @click="cancelEditName"
              class="text-outline hover:text-on-surface transition-colors"
              title="Cancelar"
            >
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>
          <p v-if="nameValidationError" class="text-error text-xs mt-1">{{ nameValidationError }}</p>
          <p v-if="editSuccess" class="text-xs font-label mt-1 text-primary">¡Nombre actualizado!</p>

          <p class="font-body text-secondary mt-1">
            {{ speciesLabels[libraryStore.selected.species] ?? libraryStore.selected.species }}
            <span v-if="libraryStore.selected.subSpecies">
              · {{ subSpeciesLabels[libraryStore.selected.subSpecies] ?? libraryStore.selected.subSpecies }}
            </span>
            · {{ classLabels[libraryStore.selected.class] ?? libraryStore.selected.class }}
            · Nivel {{ libraryStore.selected.level }}
          </p>
          <p class="text-xs text-outline mt-1">Guardado el {{ formatDate(libraryStore.selected.createdAt) }}</p>
        </div>
        <span class="text-[10px] font-bold uppercase tracking-widest text-secondary bg-surface-container px-3 py-1 mt-1">
          {{ libraryStore.selected.ruleset }}
        </span>
      </div>

      <!-- Edit error -->
      <div v-if="editError" class="bg-error-container text-on-error-container text-sm px-4 py-3 rounded">
        {{ editError }}
      </div>

      <!-- Combat stats -->
      <div class="grid grid-cols-2 gap-4">
        <div class="bg-surface-container-low p-5 flex flex-col items-center justify-center space-y-1">
          <span class="text-[10px] font-bold uppercase tracking-widest text-secondary">Puntos de Golpe</span>
          <span class="font-headline text-5xl font-bold text-primary">{{ libraryStore.selected.derived.hp }}</span>
          <span class="text-xs text-outline font-body">HP máximos</span>
        </div>
        <div class="bg-surface-container-low p-5 flex flex-col items-center justify-center space-y-1">
          <span class="text-[10px] font-bold uppercase tracking-widest text-secondary">Clase de Armadura</span>
          <span class="font-headline text-5xl font-bold text-primary">{{ libraryStore.selected.derived.ac }}</span>
          <span class="text-xs text-outline font-body">CA base</span>
        </div>
      </div>

      <!-- Trasfondo 5.5e -->
      <div v-if="libraryStore.selected.backgroundType" class="bg-surface-container-low p-5 space-y-3">
        <h3 class="font-headline text-lg font-bold text-on-surface">Trasfondo</h3>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div class="flex flex-col gap-1">
            <span class="text-[10px] font-bold uppercase tracking-widest text-secondary">Trasfondo</span>
            <span class="font-body text-on-surface text-sm">{{ libraryStore.selected.backgroundType }}</span>
          </div>
          <div class="flex flex-col gap-1">
            <span class="text-[10px] font-bold uppercase tracking-widest text-secondary">Dote de origen</span>
            <span class="font-body text-on-surface text-sm">{{ libraryStore.selected.originFeat ?? '—' }}</span>
          </div>
          <div class="flex flex-col gap-1">
            <span class="text-[10px] font-bold uppercase tracking-widest text-secondary">Distribución de atributos</span>
            <span class="font-body text-on-surface text-sm">
              {{ libraryStore.selected.asiDistribution === 'standard' ? '+2 / +1' : libraryStore.selected.asiDistribution === 'spread' ? '+1 / +1 / +1' : (libraryStore.selected.asiDistribution ?? '—') }}
            </span>
          </div>
        </div>
      </div>

      <!-- Attributes -->
      <div class="space-y-3">
        <h3 class="font-headline text-lg font-bold text-on-surface">Atributos</h3>
        <div class="grid grid-cols-3 sm:grid-cols-6 gap-2">
          <div
            v-for="stat in statKeys"
            :key="stat"
            class="bg-surface-container-low p-3 flex flex-col items-center gap-1"
          >
            <span class="text-[10px] font-bold uppercase tracking-widest text-secondary">{{ statLabels[stat] }}</span>
            <span class="font-headline text-2xl font-bold text-on-surface">{{ libraryStore.selected.finalStats[stat] }}</span>
            <span class="text-xs font-label text-outline">{{ modifierDisplay(libraryStore.selected.modifiers[stat]) }}</span>
          </div>
        </div>
      </div>

      <!-- Narrative -->
      <div class="space-y-6">
        <h3 class="font-headline text-lg font-bold text-on-surface">Narrativa</h3>
        <div
          v-for="blockKey in (['background', 'motivation', 'secret'] as const)"
          :key="blockKey"
        >
          <div class="border-t border-outline-variant/20 pt-4 first:border-t-0 first:pt-0">
            <!-- Label row -->
            <div class="flex items-center gap-2 mb-2">
              <p class="text-[10px] font-bold uppercase tracking-widest text-secondary">
                {{ narrativeTitles[blockKey] }}
              </p>
              <button
                v-if="editingBlock !== blockKey"
                @click="startEditNarrative(blockKey)"
                class="text-secondary hover:text-primary transition-colors"
                title="Editar"
              >
                <span class="material-symbols-outlined text-base">edit</span>
              </button>
              <span
                v-if="narrativeSuccess === blockKey"
                class="text-[10px] font-label font-semibold text-primary uppercase tracking-widest"
              >
                ¡Guardado!
              </span>
            </div>

            <!-- Read mode -->
            <p
              v-if="editingBlock !== blockKey"
              class="font-body text-on-surface text-sm leading-relaxed"
            >
              {{ libraryStore.selected[blockKey].content }}
            </p>

            <!-- Edit mode -->
            <div v-else class="flex flex-col gap-2">
              <textarea
                v-model="narrativeInput"
                rows="4"
                class="w-full bg-surface-container border border-outline-variant text-on-surface text-sm px-3 py-2 focus:outline-none focus:border-primary resize-none leading-relaxed"
                @keyup.escape="cancelEditNarrative"
              />
              <p v-if="narrativeValidationError" class="text-error text-xs font-label">
                {{ narrativeValidationError }}
              </p>
              <div class="flex items-center gap-2">
                <button
                  @click="submitEditNarrative"
                  :disabled="libraryStore.isLoading"
                  class="flex items-center gap-1 px-3 py-1 bg-primary text-on-primary font-label font-bold uppercase tracking-widest text-xs hover:bg-primary-container hover:text-on-primary-container transition-colors disabled:opacity-50"
                >
                  <span class="material-symbols-outlined text-sm">check</span>
                  Guardar
                </button>
                <button
                  @click="cancelEditNarrative"
                  class="flex items-center gap-1 px-3 py-1 border border-outline-variant text-secondary font-label font-bold uppercase tracking-widest text-xs hover:border-outline hover:text-on-surface transition-colors"
                >
                  <span class="material-symbols-outlined text-sm">close</span>
                  Cancelar
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Delete error -->
      <div v-if="deleteError" class="bg-error-container text-on-error-container text-sm px-4 py-3 rounded">
        {{ deleteError }}
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-4 pt-4 border-t border-outline-variant/20">
        <RouterLink
          to="/biblioteca"
          class="flex items-center gap-2 px-4 py-2 border border-outline-variant text-secondary font-label font-bold uppercase tracking-widest text-xs hover:border-primary hover:text-primary transition-colors"
        >
          <span class="material-symbols-outlined text-sm">arrow_back</span>
          Volver
        </RouterLink>
        <button
          @click="showDeleteModal = true"
          :disabled="libraryStore.isLoading"
          class="flex items-center gap-2 px-4 py-2 border border-error text-error font-label font-bold uppercase tracking-widest text-xs hover:bg-error hover:text-on-error transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
        >
          <span class="material-symbols-outlined text-sm">delete</span>
          Eliminar
        </button>
      </div>
    </div>
  </div>

  <ConfirmModal
    :open="showDeleteModal"
    title="Eliminar personaje"
    :message="`¿Eliminar a ${libraryStore.selected?.name}? Esta acción es irreversible.`"
    @confirm="confirmDelete"
    @cancel="showDeleteModal = false"
  />
</template>

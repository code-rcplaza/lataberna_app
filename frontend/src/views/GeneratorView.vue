<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import ConfigPanel from '@/components/ConfigPanel.vue'
import CharacterSheet from '@/components/CharacterSheet.vue'
import RecentCharacters from '@/components/RecentCharacters.vue'
import { useCharacterStore } from '@/stores/useCharacterStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { useLibraryStore } from '@/stores/useLibraryStore'
import { useLibraryAPI } from '@/composables/useCharacterAPI'

const router = useRouter()
const characterStore = useCharacterStore()
const authStore = useAuthStore()
const libraryStore = useLibraryStore()
const { saveCharacter } = useLibraryAPI()

const saveError = ref<string | null>(null)
const seedCopied = ref(false)

async function copySeed() {
  const seed = characterStore.current?.seed
  if (seed === undefined || seed === null) return
  await navigator.clipboard.writeText(String(seed))
  seedCopied.value = true
  setTimeout(() => { seedCopied.value = false }, 2000)
}

async function handleSave() {
  saveError.value = null
  try {
    await saveCharacter()
    router.push('/biblioteca')
  } catch (err) {
    saveError.value = err instanceof Error ? err.message : 'Error al guardar el personaje'
  }
}

const isSaveDisabled = () =>
  !characterStore.current ||
  characterStore.isSaved ||
  libraryStore.isLoading

function exportJSON() {
  const character = characterStore.current
  if (!character) return
  const blob = new Blob([JSON.stringify(character, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${character.name.toLowerCase().replace(/\s+/g, '-')}.json`
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<template>
  <div class="flex flex-col xl:flex-row gap-12 max-w-[1400px] mx-auto px-8 pb-12">
    <ConfigPanel class="xl:w-1/3" />
    <div class="xl:w-2/3 flex flex-col gap-4">
      <CharacterSheet :character="characterStore.current" :locks="characterStore.locks" />

      <!-- Seed display — solo para usuarios autenticados -->
      <div v-if="authStore.isAuthenticated && characterStore.current" class="flex items-center gap-2 text-outline text-xs font-label">
        <span class="uppercase tracking-widest font-semibold">Semilla</span>
        <span class="font-mono text-on-surface-variant">{{ characterStore.current.seed }}</span>
        <button
          @click="copySeed"
          class="flex items-center gap-1 px-2 py-0.5 border border-outline-variant/40 hover:border-outline transition-colors text-[11px] font-label font-semibold uppercase tracking-widest"
          :class="seedCopied ? 'text-primary' : 'text-outline'"
          type="button"
        >
          <span class="material-symbols-outlined text-sm leading-none">{{ seedCopied ? 'check' : 'content_copy' }}</span>
          {{ seedCopied ? '¡Copiado!' : 'Copiar' }}
        </button>
      </div>

      <!-- Autenticado: acciones -->
      <div v-if="authStore.isAuthenticated && characterStore.current" class="flex flex-col gap-2 pt-2">
        <div class="flex items-center gap-4">
          <button
            @click="exportJSON"
            class="flex items-center gap-2 px-4 py-2 border border-outline-variant text-secondary font-label font-bold uppercase tracking-widest text-xs hover:border-primary hover:text-primary transition-colors"
          >
            <span class="material-symbols-outlined text-sm">download</span>
            Exportar JSON
          </button>
          <button
            @click="handleSave"
            :disabled="isSaveDisabled()"
            class="flex items-center gap-2 px-5 py-2 bg-primary text-on-primary font-label font-bold uppercase tracking-widest text-xs hover:bg-primary-container hover:text-on-primary-container transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
          >
            <span class="material-symbols-outlined text-sm">bookmark_add</span>
            {{ characterStore.isSaved ? '¡Guardado!' : libraryStore.isLoading ? 'Guardando…' : 'Guardar en Biblioteca' }}
          </button>
        </div>
        <p v-if="saveError" class="text-error text-xs font-label">{{ saveError }}</p>
      </div>

      <!-- No autenticado + personaje generado: CTA login -->
      <div v-else-if="!authStore.isAuthenticated && characterStore.current" class="pt-2">
        <RouterLink
          to="/auth"
          class="inline-flex items-center gap-2 px-5 py-2 border border-outline-variant text-on-surface-variant font-label font-semibold uppercase tracking-widest text-xs hover:border-outline hover:text-on-surface transition-colors"
        >
          <span class="material-symbols-outlined text-sm">lock</span>
          Iniciá sesión para guardar tu personaje
        </RouterLink>
      </div>

      <!-- Historial de la sesión -->
      <RecentCharacters />
    </div>
  </div>
</template>

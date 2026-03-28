<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLibraryStore } from '@/stores/useLibraryStore'
import { useLibraryAPI } from '@/composables/useCharacterAPI'
import CharacterCard from '@/components/CharacterCard.vue'

const router = useRouter()
const libraryStore = useLibraryStore()
const { listCharacters } = useLibraryAPI()

onMounted(async () => {
  try {
    await listCharacters()
  } catch {
    // error is already set in libraryStore.error by the composable
  }
})

function goToDetail(id: string) {
  router.push(`/biblioteca/${id}`)
}
</script>

<template>
  <div class="max-w-[1400px] mx-auto px-8 pb-12">
    <h1 class="font-headline text-on-surface text-3xl mb-8">Biblioteca</h1>

    <!-- Loading skeleton -->
    <div v-if="libraryStore.isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <div
        v-for="i in 8"
        :key="i"
        class="bg-surface-container rounded-lg p-5 h-40 animate-pulse border border-outline-variant/20"
      >
        <div class="bg-outline-variant/30 rounded h-5 w-2/3 mb-3"></div>
        <div class="flex gap-2 mb-3">
          <div class="bg-outline-variant/20 rounded h-4 w-16"></div>
          <div class="bg-outline-variant/20 rounded h-4 w-16"></div>
        </div>
        <div class="bg-outline-variant/20 rounded h-3 w-1/3"></div>
      </div>
    </div>

    <!-- Error state -->
    <div v-else-if="libraryStore.error" class="flex flex-col items-center justify-center py-24 gap-4">
      <span class="material-symbols-outlined text-error text-5xl">error_outline</span>
      <p class="text-on-surface font-label text-sm text-center max-w-xs">{{ libraryStore.error }}</p>
      <button
        @click="listCharacters()"
        class="text-xs font-bold uppercase tracking-widest text-primary hover:text-primary-container transition-colors"
      >
        Reintentar
      </button>
    </div>

    <!-- Empty state -->
    <div v-else-if="libraryStore.characters.length === 0" class="flex flex-col items-center justify-center py-24 gap-4">
      <span class="material-symbols-outlined text-outline text-6xl">library_books</span>
      <p class="text-on-surface font-headline text-xl">No tenés personajes guardados aún</p>
      <p class="text-outline text-sm">Generá uno y guardalo para verlo acá.</p>
      <RouterLink
        to="/forja"
        class="mt-2 bg-primary text-on-primary px-6 py-2 font-label font-bold uppercase tracking-widest text-xs hover:bg-primary-container hover:text-on-primary-container transition-colors"
      >
        Generá tu primer personaje
      </RouterLink>
    </div>

    <!-- Character grid -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <CharacterCard
        v-for="character in libraryStore.characters"
        :key="character.id"
        :character="character"
        @select="goToDetail"
      />
    </div>
  </div>
</template>

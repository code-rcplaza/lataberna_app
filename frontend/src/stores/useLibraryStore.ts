import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Character } from '@/types/character'

export const useLibraryStore = defineStore('library', () => {
  const characters = ref<Character[]>([])
  const selected   = ref<Character | null>(null)
  const isLoading  = ref(false)
  const error      = ref<string | null>(null)

  function setCharacters(list: Character[]) {
    characters.value = list
  }

  function addCharacter(c: Character) {
    characters.value.unshift(c)
  }

  function removeCharacter(id: string) {
    characters.value = characters.value.filter(c => c.id !== id)
  }

  function setSelected(c: Character | null) {
    selected.value = c
  }

  function updateSelected(patch: Partial<Character>) {
    if (!selected.value) return
    selected.value = { ...selected.value, ...patch }
    const idx = characters.value.findIndex(c => c.id === selected.value!.id)
    if (idx !== -1) characters.value[idx] = selected.value
  }

  function reset() {
    characters.value = []
    selected.value = null
    isLoading.value = false
    error.value = null
  }

  return {
    characters,
    selected,
    isLoading,
    error,
    setCharacters,
    addCharacter,
    removeCharacter,
    setSelected,
    updateSelected,
    reset,
  }
})

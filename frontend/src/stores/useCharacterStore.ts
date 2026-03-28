import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Character, GeneratorInput, GeneratorLocks } from '@/types/character'

export const useCharacterStore = defineStore('character', () => {
  const current = ref<Character | null>(null)
  const isLoading = ref(false)
  const isSaved = ref(false)

  const input = ref<GeneratorInput>({
    species: 'random',
    subSpecies: 'random',
    class: 'random',
    gender: 'random',
    alignment: 'random',
  })

  const locks = ref<GeneratorLocks>({
    species: false,
    subSpecies: false,
    class: false,
    gender: false,
    alignment: false,
    stats: false,
    background: false,
    motivation: false,
    secret: false,
  })

  function toggleLock(field: keyof GeneratorLocks) {
    locks.value[field] = !locks.value[field]
  }

  function setCharacter(c: Character) {
    current.value = c
    isSaved.value = false
  }

  // Updates current with the persisted library character (which has the library ID).
  // Called after a successful saveCharacter — preserves the library ID for history navigation.
  function setSaved(c: Character) {
    current.value = c
    isSaved.value = true
  }

  return { current, isLoading, isSaved, input, locks, toggleLock, setCharacter, setSaved }
})

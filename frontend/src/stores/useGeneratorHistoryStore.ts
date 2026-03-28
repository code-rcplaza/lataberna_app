import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import type { Character } from '@/types/character'

export interface HistoryEntry {
  character: Character
  isSaved: boolean
}

const STORAGE_KEY = 'forge_generator_recent'
const MAX_RECENT = 3

function loadFromStorage(): HistoryEntry[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? (JSON.parse(raw) as HistoryEntry[]) : []
  } catch {
    return []
  }
}

export const useGeneratorHistoryStore = defineStore('generatorHistory', () => {
  const recent = ref<HistoryEntry[]>(loadFromStorage())

  watch(recent, (val) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(val))
  }, { deep: true })

  // Push a character to the front of recent history, trimming to MAX_RECENT.
  function push(character: Character, isSaved: boolean) {
    recent.value = [{ character, isSaved }, ...recent.value].slice(0, MAX_RECENT)
  }

  // Remove the entry at index, and optionally prepend current to the front.
  // Used when the user loads a recent character: it leaves its slot and current takes its place.
  function rotateIn(index: number, current: Character | null, currentIsSaved: boolean) {
    const without = recent.value.filter((_, i) => i !== index)
    if (current) {
      recent.value = [{ character: current, isSaved: currentIsSaved }, ...without].slice(0, MAX_RECENT)
    } else {
      recent.value = without
    }
  }

  return { recent, push, rotateIn }
})

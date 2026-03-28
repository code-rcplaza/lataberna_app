import { defineStore } from 'pinia'
import { ref } from 'vue'
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

function persist(entries: HistoryEntry[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(entries))
}

export const useGeneratorHistoryStore = defineStore('generatorHistory', () => {
  const recent = ref<HistoryEntry[]>(loadFromStorage())

  function push(character: Character, isSaved: boolean) {
    recent.value = [{ character, isSaved }, ...recent.value].slice(0, MAX_RECENT)
    persist(recent.value)
  }

  function rotateIn(index: number, current: Character | null, currentIsSaved: boolean) {
    const without = recent.value.filter((_, i) => i !== index)
    recent.value = current
      ? [{ character: current, isSaved: currentIsSaved }, ...without].slice(0, MAX_RECENT)
      : without
    persist(recent.value)
  }

  return { recent, push, rotateIn }
})

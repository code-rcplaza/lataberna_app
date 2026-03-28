import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGeneratorHistoryStore } from '../useGeneratorHistoryStore'
import type { Character } from '@/types/character'

function makeCharacter(id: string, name = 'Personaje'): Character {
  return {
    id,
    name,
    species: 'human',
    subSpecies: null,
    class: 'fighter',
    level: 1,
    ruleset: '5e',
    seed: 1,
    baseStats:  { STR: 10, DEX: 10, CON: 10, INT: 10, WIS: 10, CHA: 10 },
    finalStats: { STR: 10, DEX: 10, CON: 10, INT: 10, WIS: 10, CHA: 10 },
    modifiers:  { STR: 0,  DEX: 0,  CON: 0,  INT: 0,  WIS: 0,  CHA: 0  },
    derived: { hp: 8, ac: 10 },
    background: { category: 'background', content: '', tags: [] },
    motivation: { category: 'motivation', content: '', tags: [] },
    secret:     { category: 'secret',     content: '', tags: [] },
    locks: { name: false, stats: false, background: false, motivation: false, secret: false },
  }
}

describe('useGeneratorHistoryStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  describe('push', () => {
    it('adds a character to the front of recent', () => {
      const store = useGeneratorHistoryStore()
      store.push(makeCharacter('c1'), false)
      expect(store.recent).toHaveLength(1)
      expect(store.recent[0].character.id).toBe('c1')
    })

    it('prepends — most recent is always at index 0', () => {
      const store = useGeneratorHistoryStore()
      store.push(makeCharacter('c1'), false)
      store.push(makeCharacter('c2'), false)
      expect(store.recent[0].character.id).toBe('c2')
      expect(store.recent[1].character.id).toBe('c1')
    })

    it('caps at 3 entries', () => {
      const store = useGeneratorHistoryStore()
      store.push(makeCharacter('c1'), false)
      store.push(makeCharacter('c2'), false)
      store.push(makeCharacter('c3'), false)
      store.push(makeCharacter('c4'), false)
      expect(store.recent).toHaveLength(3)
      expect(store.recent[0].character.id).toBe('c4')
      expect(store.recent[2].character.id).toBe('c2')
    })

    it('stores the isSaved flag correctly', () => {
      const store = useGeneratorHistoryStore()
      store.push(makeCharacter('c1'), true)
      store.push(makeCharacter('c2'), false)
      expect(store.recent[0].isSaved).toBe(false)
      expect(store.recent[1].isSaved).toBe(true)
    })
  })

  describe('rotateIn', () => {
    it('removes the entry at the given index', () => {
      const store = useGeneratorHistoryStore()
      store.push(makeCharacter('c1'), false)
      store.push(makeCharacter('c2'), false)
      store.push(makeCharacter('c3'), false)
      // recent = [c3, c2, c1]
      store.rotateIn(1, null, false) // remove c2
      expect(store.recent.map(e => e.character.id)).toEqual(['c3', 'c1'])
    })

    it('prepends current when provided', () => {
      const store = useGeneratorHistoryStore()
      store.push(makeCharacter('c1'), false)
      store.push(makeCharacter('c2'), false)
      // recent = [c2, c1]
      store.rotateIn(0, makeCharacter('active'), true)
      expect(store.recent[0].character.id).toBe('active')
      expect(store.recent[0].isSaved).toBe(true)
      expect(store.recent[1].character.id).toBe('c1')
    })

    it('does not add current when null', () => {
      const store = useGeneratorHistoryStore()
      store.push(makeCharacter('c1'), false)
      store.push(makeCharacter('c2'), false)
      // recent = [c2, c1]
      store.rotateIn(0, null, false) // remove c2, no current to add
      expect(store.recent).toHaveLength(1)
      expect(store.recent[0].character.id).toBe('c1')
    })

    it('respects the 3-entry cap after rotation', () => {
      const store = useGeneratorHistoryStore()
      store.push(makeCharacter('c1'), false)
      store.push(makeCharacter('c2'), false)
      store.push(makeCharacter('c3'), false)
      // recent = [c3, c2, c1] — already at cap
      // rotate index 2 (c1): remove c1, add active → [active, c3, c2]
      store.rotateIn(2, makeCharacter('active'), false)
      expect(store.recent).toHaveLength(3)
      expect(store.recent.map(e => e.character.id)).toEqual(['active', 'c3', 'c2'])
    })
  })

  describe('localStorage persistence', () => {
    it('persists recent on push', async () => {
      const store = useGeneratorHistoryStore()
      store.push(makeCharacter('c1'), false)
      // flush watchers
      await new Promise(r => setTimeout(r, 0))
      const stored = JSON.parse(localStorage.getItem('forge_generator_recent') ?? '[]')
      expect(stored).toHaveLength(1)
      expect(stored[0].character.id).toBe('c1')
    })

    it('loads from localStorage on store creation', async () => {
      // Pre-populate storage
      const entry = { character: makeCharacter('persisted'), isSaved: false }
      localStorage.setItem('forge_generator_recent', JSON.stringify([entry]))

      // New pinia instance simulates a fresh page load
      setActivePinia(createPinia())
      const store = useGeneratorHistoryStore()
      expect(store.recent).toHaveLength(1)
      expect(store.recent[0].character.id).toBe('persisted')
    })
  })
})

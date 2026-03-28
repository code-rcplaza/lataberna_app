import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCharacterStore } from '../useCharacterStore'
import type { Character } from '@/types/character'

function makeCharacter(overrides: Partial<Character> = {}): Character {
  return {
    id: 'char-001',
    name: 'Thorin',
    species: 'dwarf',
    subSpecies: 'mountain-dwarf',
    class: 'fighter',
    level: 1,
    ruleset: '5e',
    abilityBonusSource: 'species',
    seed: 12345,
    createdAt: '2026-01-01T00:00:00Z',
    baseStats:  { STR: 15, DEX: 10, CON: 14, INT: 8, WIS: 12, CHA: 9 },
    finalStats: { STR: 15, DEX: 10, CON: 16, INT: 8, WIS: 12, CHA: 9 },
    modifiers:  { STR: 2,  DEX: 0,  CON: 3,  INT: -1, WIS: 1, CHA: -1 },
    derived: { hp: 11, ac: 16 },
    background: { category: 'background', content: 'Criado en las montañas.', tags: [] },
    motivation: { category: 'motivation', content: 'Proteger a su clan.', tags: [] },
    secret:     { category: 'secret',     content: 'Teme a la oscuridad.', tags: [] },
    locks: { name: false, stats: false, background: false, motivation: false, secret: false },
    ...overrides,
  }
}

describe('useCharacterStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with no current character', () => {
    const store = useCharacterStore()
    expect(store.current).toBeNull()
    expect(store.isSaved).toBe(false)
  })

  describe('setCharacter', () => {
    it('sets the current character', () => {
      const store = useCharacterStore()
      const char = makeCharacter()
      store.setCharacter(char)
      expect(store.current).toEqual(char)
    })

    it('resets isSaved to false', () => {
      const store = useCharacterStore()
      const char = makeCharacter()
      store.setSaved(char)
      expect(store.isSaved).toBe(true)

      store.setCharacter(makeCharacter({ id: 'char-002' }))
      expect(store.isSaved).toBe(false)
    })
  })

  describe('setSaved', () => {
    it('updates current with the saved character', () => {
      const store = useCharacterStore()
      const draft = makeCharacter({ id: 'draft-001' })
      store.setCharacter(draft)

      const saved = makeCharacter({ id: 'library-001' })
      store.setSaved(saved)

      expect(store.current?.id).toBe('library-001')
    })

    it('sets isSaved to true', () => {
      const store = useCharacterStore()
      store.setCharacter(makeCharacter())
      store.setSaved(makeCharacter({ id: 'library-001' }))
      expect(store.isSaved).toBe(true)
    })
  })

  describe('toggleLock', () => {
    it('flips the lock for a given field', () => {
      const store = useCharacterStore()
      expect(store.locks.stats).toBe(false)
      store.toggleLock('stats')
      expect(store.locks.stats).toBe(true)
      store.toggleLock('stats')
      expect(store.locks.stats).toBe(false)
    })

    it('only affects the targeted field', () => {
      const store = useCharacterStore()
      store.toggleLock('background')
      expect(store.locks.background).toBe(true)
      expect(store.locks.motivation).toBe(false)
      expect(store.locks.secret).toBe(false)
    })
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCharacterStore } from '@/stores/useCharacterStore'
import { useGeneratorHistoryStore } from '@/stores/useGeneratorHistoryStore'
import { useCharacterAPI } from '@/composables/useCharacterAPI'
import type { Character } from '@/types/character'

// Mock gql so no real HTTP calls are made
vi.mock('@/composables/useGraphQL', () => ({
  gql: vi.fn(),
}))

import { gql } from '@/composables/useGraphQL'
const mockGql = vi.mocked(gql)

function makeCharacter(seed: number): Character {
  return {
    id: 'char-001',
    name: 'Thorin',
    species: 'dwarf',
    subSpecies: 'mountain-dwarf',
    class: 'fighter',
    level: 1,
    ruleset: '5e',
    abilityBonusSource: 'species',
    seed,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
    baseStats:  { STR: 15, DEX: 10, CON: 14, INT: 8, WIS: 12, CHA: 9 },
    finalStats: { STR: 15, DEX: 10, CON: 16, INT: 8, WIS: 12, CHA: 9 },
    modifiers:  { STR: 2,  DEX: 0,  CON: 3,  INT: -1, WIS: 1, CHA: -1 },
    derived: { hp: 11, ac: 16 },
    background: { category: 'background', content: 'Un pasado oscuro.', tags: [] },
    motivation: { category: 'motivation', content: 'Proteger a su clan.', tags: [] },
    secret:     { category: 'secret',     content: 'Teme a la oscuridad.', tags: [] },
    locks: { name: false, stats: false, background: false, motivation: false, secret: false },
  }
}

describe('useCharacterAPI — generate()', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGql.mockReset()
  })

  it('does NOT send seed to the API when store.input.seed is not set', async () => {
    const store = useCharacterStore()
    store.input.seed = undefined

    mockGql.mockResolvedValueOnce({ generateCharacter: makeCharacter(99999) })

    const { generate } = useCharacterAPI()
    await generate()

    const variables = mockGql.mock.calls[0][1] as { input: Record<string, unknown> }
    expect(variables.input).not.toHaveProperty('seed')
  })

  it('sends the seed to the API when store.input.seed is set', async () => {
    const store = useCharacterStore()
    store.input.seed = 12345

    mockGql.mockResolvedValueOnce({ generateCharacter: makeCharacter(12345) })

    const { generate } = useCharacterAPI()
    await generate()

    const variables = mockGql.mock.calls[0][1] as { input: Record<string, unknown> }
    expect(variables.input.seed).toBe(12345)
  })

  it('sends the seed as a number, not a string', async () => {
    const store = useCharacterStore()
    store.input.seed = 99887766

    mockGql.mockResolvedValueOnce({ generateCharacter: makeCharacter(99887766) })

    const { generate } = useCharacterAPI()
    await generate()

    const variables = mockGql.mock.calls[0][1] as { input: Record<string, unknown> }
    expect(typeof variables.input.seed).toBe('number')
  })

  it('does NOT update store.input.seed after generation — user seed stays as-is', async () => {
    const store = useCharacterStore()
    store.input.seed = 12345

    // Backend returns a character with a DIFFERENT seed (e.g. it normalized it)
    mockGql.mockResolvedValueOnce({ generateCharacter: makeCharacter(99999) })

    const { generate } = useCharacterAPI()
    await generate()

    // store.input.seed must not be overwritten by the character's seed
    expect(store.input.seed).toBe(12345)
  })
})

describe('useCharacterAPI — generate() history deduplication', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGql.mockReset()
    localStorage.clear()
  })

  it('does NOT push to history when re-generating with the same seed', async () => {
    const store = useCharacterStore()
    const history = useGeneratorHistoryStore()

    // Current character has seed 12345
    store.setCharacter(makeCharacter(12345))
    store.input.seed = 12345

    // Backend returns the same character again (deterministic seed)
    mockGql.mockResolvedValueOnce({ generateCharacter: makeCharacter(12345) })

    const { generate } = useCharacterAPI()
    await generate()

    expect(history.recent).toHaveLength(0)
  })

  it('DOES push to history when generating with a different seed', async () => {
    const store = useCharacterStore()
    const history = useGeneratorHistoryStore()

    // Current character has seed 12345, user changed to seed 99999
    store.setCharacter(makeCharacter(12345))
    store.input.seed = 99999

    mockGql.mockResolvedValueOnce({ generateCharacter: makeCharacter(99999) })

    const { generate } = useCharacterAPI()
    await generate()

    expect(history.recent).toHaveLength(1)
    expect(history.recent[0].character.seed).toBe(12345)
  })

  it('DOES push to history when generating without a seed (random)', async () => {
    const store = useCharacterStore()
    const history = useGeneratorHistoryStore()

    store.setCharacter(makeCharacter(12345))
    store.input.seed = undefined

    mockGql.mockResolvedValueOnce({ generateCharacter: makeCharacter(77777) })

    const { generate } = useCharacterAPI()
    await generate()

    expect(history.recent).toHaveLength(1)
    expect(history.recent[0].character.seed).toBe(12345)
  })
})

describe('useCharacterAPI — regenerateField()', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockGql.mockReset()
  })

  it('does NOT send top-level $seed — field re-rolls use current.seed for determinism', async () => {
    const store = useCharacterStore()
    store.setCharacter(makeCharacter(12345))

    mockGql.mockResolvedValueOnce({ regenerateDraft: makeCharacter(12345) })

    const { regenerateField } = useCharacterAPI()
    await regenerateField('stats')

    const variables = mockGql.mock.calls[0][1] as Record<string, unknown>
    // $seed top-level is not sent — determinism comes from current.seed inside the current object
    expect(variables.seed).toBeUndefined()
  })

  it('passes current.seed as part of the current character state', async () => {
    const store = useCharacterStore()
    store.setCharacter(makeCharacter(12345))

    mockGql.mockResolvedValueOnce({ regenerateDraft: makeCharacter(12345) })

    const { regenerateField } = useCharacterAPI()
    await regenerateField('stats')

    const variables = mockGql.mock.calls[0][1] as { current: Record<string, unknown> }
    expect(variables.current.seed).toBe(12345)
  })
})

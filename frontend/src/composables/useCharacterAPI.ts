import { useCharacterStore } from '@/stores/useCharacterStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { useLibraryStore } from '@/stores/useLibraryStore'
import { useGeneratorHistoryStore } from '@/stores/useGeneratorHistoryStore'
import { gql } from './useGraphQL'
import type { Character, EditCharacterInput } from '@/types/character'

const GENERATE_CHARACTER = `
  mutation GenerateCharacter($input: GenerateCharacterInput) {
    generateCharacter(input: $input) {
      id name species subSpecies class level ruleset seed
      backgroundType asiDistribution originFeat
      baseStats { STR DEX CON INT WIS CHA }
      finalStats { STR DEX CON INT WIS CHA }
      modifiers  { STR DEX CON INT WIS CHA }
      derived    { hp ac }
      background { category content tags }
      motivation { category content tags }
      secret     { category content tags }
      locks { name stats background motivation secret }
    }
  }
`

const REGENERATE_DRAFT = `
  mutation RegenerateDraft($current: CurrentCharacterInput!, $locks: CharacterLocksInput!, $seed: Int) {
    regenerateDraft(current: $current, locks: $locks, seed: $seed) {
      id name species subSpecies class level ruleset seed
      backgroundType asiDistribution originFeat
      baseStats { STR DEX CON INT WIS CHA }
      finalStats { STR DEX CON INT WIS CHA }
      modifiers  { STR DEX CON INT WIS CHA }
      derived    { hp ac }
      background { category content tags }
      motivation { category content tags }
      secret     { category content tags }
      locks { name stats background motivation secret }
    }
  }
`

// Summary fields — used in BibliotecaView card grid (avoids over-fetching stats/narrative)
const CHARACTER_LIST_FIELDS = `
  id name species subSpecies class level ruleset createdAt
  derived { hp ac }
`

// Full fields — used in CharacterDetailView
const CHARACTER_DETAIL_FIELDS = `
  id name species subSpecies class level ruleset seed createdAt updatedAt
  baseStats { STR DEX CON INT WIS CHA }
  finalStats { STR DEX CON INT WIS CHA }
  modifiers  { STR DEX CON INT WIS CHA }
  derived    { hp ac }
  background { category content tags }
  motivation { category content tags }
  secret     { category content tags }
  locks { name stats background motivation secret }
  backgroundType asiDistribution originFeat
`

const LIST_CHARACTERS = `
  query ListCharacters {
    characters {
      ${CHARACTER_LIST_FIELDS}
    }
  }
`

const GET_CHARACTER = `
  query GetCharacter($id: ID!) {
    character(id: $id) {
      ${CHARACTER_DETAIL_FIELDS}
    }
  }
`

const SAVE_CHARACTER = `
  mutation SaveCharacter($input: GenerateCharacterInput!, $seed: Int!) {
    saveCharacter(input: $input, seed: $seed) {
      ${CHARACTER_DETAIL_FIELDS}
    }
  }
`

const DELETE_CHARACTER = `
  mutation DeleteCharacter($id: ID!) {
    deleteCharacter(id: $id)
  }
`

const EDIT_CHARACTER = `
  mutation EditCharacter($id: ID!, $input: EditCharacterInput!) {
    editCharacter(id: $id, input: $input) {
      ${CHARACTER_DETAIL_FIELDS}
    }
  }
`

export function useCharacterAPI() {
  const store = useCharacterStore()
  const auth = useAuthStore()
  const historyStore = useGeneratorHistoryStore()

  async function generate() {
    store.isLoading = true
    try {
      const input: Record<string, unknown> = {}
      if (store.input.class    !== 'random') input.class    = store.input.class
      if (store.input.species  !== 'random') input.species  = store.input.species
      if (store.input.subSpecies !== 'random') input.subSpecies = store.input.subSpecies
      if (store.input.gender   !== 'random') input.gender   = store.input.gender
      if (store.input.seed)                  input.seed     = store.input.seed

      const data = await gql<{ generateCharacter: Character }>(
        GENERATE_CHARACTER,
        { input },
        auth.sessionId,
      )
      // Push current to history before replacing it with the new character.
      // Skip if we're re-generating with the same seed — same seed = same character,
      // adding it again would create duplicates in the history.
      const isRepeatSeed = store.input.seed !== undefined && store.current?.seed === store.input.seed
      if (store.current && !isRepeatSeed) {
        historyStore.push(store.current, store.isSaved)
      }
      store.setCharacter(data.generateCharacter)
    } finally {
      store.isLoading = false
    }
  }

  // field: which field to refresh — everything else is implicitly locked.
  async function regenerateField(field: 'name' | 'stats' | 'background' | 'motivation' | 'secret') {
    const c = store.current
    if (!c) return

    store.isLoading = true
    try {
      // All locks true except the target field
      const locks = {
        name:       field !== 'name',
        stats:      field !== 'stats',
        background: field !== 'background',
        motivation: field !== 'motivation',
        secret:     field !== 'secret',
      }

      const current = {
        name:       c.name,
        class:      c.class,
        species:    c.species,
        subSpecies: c.subSpecies ?? null,
        seed:       c.seed ?? null,
        finalStats: c.finalStats,
        derived:    c.derived,
        background: { category: c.background.category, content: c.background.content, tags: c.background.tags },
        motivation: { category: c.motivation.category, content: c.motivation.content, tags: c.motivation.tags },
        secret:     { category: c.secret.category,     content: c.secret.content,     tags: c.secret.tags },
      }

      const data = await gql<{ regenerateDraft: Character }>(
        REGENERATE_DRAFT,
        { current, locks },
        auth.sessionId,
      )
      store.setCharacter(data.regenerateDraft)
    } finally {
      store.isLoading = false
    }
  }

  return { generate, regenerateField }
}

export function useLibraryAPI() {
  const auth = useAuthStore()
  const libraryStore = useLibraryStore()
  const characterStore = useCharacterStore()

  async function listCharacters(): Promise<void> {
    libraryStore.isLoading = true
    libraryStore.error = null
    try {
      const data = await gql<{ characters: Character[] }>(
        LIST_CHARACTERS,
        {},
        auth.sessionId,
      )
      libraryStore.setCharacters(data.characters)
    } catch (err) {
      libraryStore.error = err instanceof Error ? err.message : 'Error al cargar los personajes'
      throw err
    } finally {
      libraryStore.isLoading = false
    }
  }

  async function getCharacter(id: string): Promise<void> {
    libraryStore.isLoading = true
    libraryStore.error = null
    try {
      const data = await gql<{ character: Character | null }>(
        GET_CHARACTER,
        { id },
        auth.sessionId,
      )
      libraryStore.setSelected(data.character)
    } catch (err) {
      libraryStore.error = err instanceof Error ? err.message : 'Error al cargar el personaje'
      throw err
    } finally {
      libraryStore.isLoading = false
    }
  }

  // NOTE: seed is non-nullable in the schema. This function MUST NOT be called
  // when characterStore.current is null. The caller (GeneratorView) gates the
  // "Guardar" button on current !== null to guarantee safety.
  async function saveCharacter(): Promise<void> {
    const current = characterStore.current
    if (!current || current.seed === undefined || current.seed === null) return

    libraryStore.isLoading = true
    libraryStore.error = null
    try {
      const input: Record<string, unknown> = {}
      if (characterStore.input.class     !== 'random') input.class     = characterStore.input.class
      if (characterStore.input.species   !== 'random') input.species   = characterStore.input.species
      if (characterStore.input.subSpecies !== 'random') input.subSpecies = characterStore.input.subSpecies
      if (characterStore.input.gender    !== 'random') input.gender    = characterStore.input.gender

      const data = await gql<{ saveCharacter: Character }>(
        SAVE_CHARACTER,
        { input, seed: current.seed },
        auth.sessionId,
      )
      libraryStore.addCharacter(data.saveCharacter)
      characterStore.setSaved(data.saveCharacter)
    } catch (err) {
      libraryStore.error = err instanceof Error ? err.message : 'Error al guardar el personaje'
      throw err
    } finally {
      libraryStore.isLoading = false
    }
  }

  async function deleteCharacter(id: string): Promise<void> {
    libraryStore.isLoading = true
    libraryStore.error = null
    try {
      await gql<{ deleteCharacter: boolean }>(
        DELETE_CHARACTER,
        { id },
        auth.sessionId,
      )
      libraryStore.removeCharacter(id)
    } catch (err) {
      libraryStore.error = err instanceof Error ? err.message : 'Error al eliminar el personaje'
      throw err
    } finally {
      libraryStore.isLoading = false
    }
  }

  async function editCharacter(id: string, patch: EditCharacterInput): Promise<void> {
    libraryStore.isLoading = true
    libraryStore.error = null
    try {
      const data = await gql<{ editCharacter: Character }>(
        EDIT_CHARACTER,
        { id, input: patch },
        auth.sessionId,
      )
      libraryStore.updateSelected(data.editCharacter)
    } catch (err) {
      libraryStore.error = err instanceof Error ? err.message : 'Error al editar el personaje'
      throw err
    } finally {
      libraryStore.isLoading = false
    }
  }

  return { listCharacters, getCharacter, saveCharacter, deleteCharacter, editCharacter }
}

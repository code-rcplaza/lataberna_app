export interface User {
  id: string
  email: string
}

export type Species =
  | 'human' | 'elf' | 'dwarf' | 'halfling' | 'gnome'
  | 'half-elf' | 'half-orc' | 'tiefling' | 'dragonborn'

export type SubSpecies =
  | 'high-elf' | 'wood-elf' | 'drow'
  | 'hill-dwarf' | 'mountain-dwarf'
  | 'lightfoot' | 'stout'
  | 'forest-gnome' | 'rock-gnome'
  | 'tiefling-infernal' | 'tiefling-virtue'

export type Class =
  | 'barbarian' | 'bard' | 'cleric' | 'druid' | 'fighter'
  | 'monk' | 'paladin' | 'ranger' | 'rogue' | 'sorcerer'
  | 'warlock' | 'wizard' | 'artificer'

export type Ruleset = '5e' | '5.5e'
export type AbilityBonusSource = 'species' | 'background' | 'none'
export type NarrativeCategory = 'background' | 'motivation' | 'secret'

export interface Stats {
  STR: number
  DEX: number
  CON: number
  INT: number
  WIS: number
  CHA: number
}

export interface Modifiers {
  STR: number
  DEX: number
  CON: number
  INT: number
  WIS: number
  CHA: number
}

export interface DerivedStats {
  hp: number
  ac: number
}

export interface NarrativeBlock {
  category: NarrativeCategory
  content: string
  tags: string[]
}

export interface CharacterLocks {
  name: boolean
  stats: boolean
  background: boolean
  motivation: boolean
  secret: boolean
}

export interface Character {
  id: string
  name: string
  species: Species
  subSpecies?: SubSpecies
  class: Class
  level: number
  ruleset: Ruleset
  abilityBonusSource: AbilityBonusSource
  baseStats: Stats
  finalStats: Stats
  modifiers: Modifiers
  derived: DerivedStats
  background: NarrativeBlock
  motivation: NarrativeBlock
  secret: NarrativeBlock
  locks: CharacterLocks
  seed?: number
  backgroundType?: string
  backgroundDescription?: string
  asiDistribution?: string
  originFeat?: string
  createdAt: string
  updatedAt: string
}

export type Alignment =
  | 'legal-bueno' | 'legal-neutral' | 'legal-malvado'
  | 'neutral-bueno' | 'neutral' | 'neutral-malvado'
  | 'caotico-bueno' | 'caotico-neutral' | 'caotico-malvado'

export type Gender = 'male' | 'female' | 'random'

export interface GeneratorInput {
  name?: string
  species?: Species | 'random'
  subSpecies?: SubSpecies | 'random'
  class?: Class | 'random'
  gender?: Gender
  alignment?: Alignment | 'random'
  seed?: number
}

export interface EditCharacterInput {
  name?: string
  background?: string
  motivation?: string
  secret?: string
}

export interface GeneratorLocks {
  name: boolean
  species: boolean
  subSpecies: boolean
  class: boolean
  gender: boolean
  alignment: boolean
  stats: boolean
  background: boolean
  motivation: boolean
  secret: boolean
}

# domain-model.md — rpg_engine

Entidades, interfaces y tipos del dominio.
Este archivo cambia solo cuando cambia la estructura del modelo, no los datos.

---

## Entidad principal: Character

```typescript
type Ruleset = '5e' | '5.5e'
type AbilityBonusSource = 'species' | 'background' | 'none'
type ArmorCategory = 'none' | 'light' | 'medium' | 'heavy' | 'shield'

interface Character {
  id: string

  // Identidad
  name: string
  species: Species
  subSpecies?: SubSpecies
  class: Class
  level: number

  // Configuración de reglas
  ruleset: Ruleset
  abilityBonusSource: AbilityBonusSource

  // Núcleo mecánico
  baseStats: Stats       // antes de bonos
  finalStats: Stats      // después de bonos
  modifiers: Modifiers   // calculados sobre finalStats
  derived: DerivedStats

  // Narrativa
  background: NarrativeBlock
  motivation: NarrativeBlock
  secret: NarrativeBlock

  // Estado de regeneración
  locks: CharacterLocks

  // Metadata
  seed?: number
  createdAt: Date
  updatedAt: Date
}
```

---

## Stats y derivados

```typescript
interface Stats {
  STR: number
  DEX: number
  CON: number
  INT: number
  WIS: number
  CHA: number
}

interface Modifiers {
  STR: number
  DEX: number
  CON: number
  INT: number
  WIS: number
  CHA: number
}

interface DerivedStats {
  hp: number
  ac: number
}
```

---

## Narrativa

```typescript
interface NarrativeBlock {
  category: 'background' | 'motivation' | 'secret'
  content: string
  tags: string[]  // clases y species compatibles; 'any' = universal
}
```

---

## Armadura

```typescript
interface ArmorType {
  name: string
  category: ArmorCategory
  baseAC: number
  maxDex?: number         // undefined → sin límite, 2 → medium, 0 → heavy
  strRequirement?: number
}
```

---

## Bonos (desacoplados del origen)

```typescript
interface AbilityBonus {
  stat: keyof Stats
  value: number
  source: 'species' | 'background'
}
```

---

## Locks (regeneración parcial)

```typescript
interface CharacterLocks {
  name: boolean
  stats: boolean
  background: boolean
  motivation: boolean
  secret: boolean
}
```

Un parámetro provisto en el request actúa como lock implícito.
El orquestador preserva los campos bloqueados al regenerar.

---

## Resolución de bonos

```typescript
function resolveAbilityBonuses(input, config): AbilityBonus[] {
  switch (config.abilityBonusSource) {
    case 'species':    return getSpeciesBonuses(input.species, input.subSpecies)
    case 'background': return getBackgroundBonuses(input.background)
    default:           return []
  }
}
```

MVP opera con: `ruleset: '5e'`, `abilityBonusSource: 'species'`

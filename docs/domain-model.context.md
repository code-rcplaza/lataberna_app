# domain-model.md — rpg_engine

Entidades, interfaces y tipos del dominio.
Este archivo cambia solo cuando cambia la estructura del modelo, no los datos.

---

## Entidad principal: Character

```typescript
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
  ruleset: string             // siempre "5.5e" — fijo tras el pivote; guardado para fidelidad en save/load
  abilityBonusSource: string  // siempre "background" en 5.5e — las species no otorgan ASIs
  asiDistribution: string     // "standard" (+2/+1) | "spread" (+1/+1/+1)

  // Núcleo mecánico
  baseStats: Stats       // antes de bonos
  finalStats: Stats      // después de bonos de background
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

## Background (mecánico — 5.5e)

```typescript
interface Background {
  name: string
  asiPool: [string, string, string]  // ej: ["WIS", "INT", "CON"] — definido por el background
  originFeat: string                 // feat fijo otorgado por este background (no elegido)
  tags: string[]                     // tags de coherencia por clase/species; 'any' = universal
}
```

`asiPool` define los 3 stats elegibles para el bonus. La distribución (`standard` o `spread`)
se elige o genera aleatoriamente y se almacena en `Character.asiDistribution`.

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
  source: 'background'  // en 5.5e el único origen de ASIs es el background
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

En 5.5e los ASIs provienen exclusivamente del background. La función recibe el background
resuelto y la distribución elegida:

```typescript
function resolveAbilityBonuses(background: Background, distribution: 'standard' | 'spread'): AbilityBonus[] {
  if (distribution === 'spread') {
    // +1 a los 3 stats del asiPool
    return background.asiPool.map(stat => ({ stat, value: 1, source: 'background' }))
  }
  // distribution === 'standard': +2 al primero del pool, +1 al segundo
  return [
    { stat: background.asiPool[0], value: 2, source: 'background' },
    { stat: background.asiPool[1], value: 1, source: 'background' },
  ]
}
```

El MVP opera con: `ruleset: '5.5e'`, `abilityBonusSource: 'background'`.

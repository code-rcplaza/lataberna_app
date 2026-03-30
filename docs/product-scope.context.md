# rpg_engine — Scope del MVP (v3, definitivo)

## Nombre tentativo
**La Taberna RPG** — Generador de contenido RPG en español para D&D 5.5e (2024 PHB)

---

## Visión
Herramienta en español que permite a jugadores y Dungeon Masters generar
personajes y NPCs completos, coherentes y listos para usar en mesa,
combinando mecánicas y narrativa en un solo flujo.

---

## Problema
Las herramientas existentes:
- Están en inglés (o son malas traducciones)
- Generan contenido desconectado entre sí
- No están diseñadas para el flujo real de una sesión

Esto obliga a los DMs y jugadores a improvisar bajo presión, traducir
contenido al vuelo y perder tiempo en preparación.

---

## Usuario objetivo
**Primario:** DMs y jugadores de D&D 5.5e (2024 PHB) hispanohablantes, nivel principiante a intermedio.
**Flujo unificado:** el mismo generador sirve para crear un personaje propio (PC)
o un NPC para la sesión. La diferencia es intencional, no técnica.

---

## Propuesta de valor
> Genera un personaje o NPC completo, coherente y en español en segundos.
> Narrativa + mecánicas integradas. Listo para jugar.

**Diferenciadores:**
- Español nativo (no traducido)
- Coherencia entre nombre, species, clase, stats e historia
- Regeneración parcial (cambia solo lo que no te gusta)
- Biblioteca personal de personajes guardados

---

## Principio de generación

> **Todo parámetro es opcional. Lo que se omite se genera aleatoriamente. Lo que se provee queda locked.**

Esto permite dos flujos extremos con el mismo generador:

| Flujo | Input | Resultado |
|---|---|---|
| Cero clicks | Ningún parámetro | Personaje 100% aleatorio |
| Totalmente guiado | Todos los parámetros | Personaje con restricciones exactas |

Parámetros opcionales del generador:

| Parámetro | Si se omite |
|---|---|
| `class` | Se elige aleatoriamente |
| `species` | Se elige aleatoriamente |
| `subSpecies` | Se elige aleatoriamente (válida para la species resuelta) |
| `gender` | Se elige aleatoriamente |
| `seed` | Random — resultado no reproducible |

Un DM puede pedir "un guerrero enano" y obtener nombre, stats y narrativa coherentes.
Un DM también puede pedir "sorpréndeme" y obtener un personaje completamente inesperado.

---



### Módulo 1 — Generador de Nombres
- 9 species core con subspecies:
  - Human
  - Elf → High Elf, Wood Elf, Dark Elf (Drow)
  - Dwarf → Hill Dwarf, Mountain Dwarf
  - Halfling → Lightfoot, Stout
  - Gnome → Forest Gnome, Rock Gnome
  - Half-Elf
  - Half-Orc
  - Tiefling (infernal vs virtue names)
  - Dragonborn (clan first)
- Respeta las convenciones de nomenclatura de cada species/subSpecies
- Mínimo 50 entradas por género/species en la seed — criterio de aceptación no negociable

### Módulo 2 — Generador de Stat Block

**Clases del MVP (13):**

| Clase | Hit Die | Stats principales |
|---|---|---|
| Barbarian | d12 | STR / CON |
| Bard | d8 | CHA / DEX |
| Cleric | d8 | WIS / CON |
| Druid | d8 | WIS / CON |
| Fighter | d10 | STR o DEX / CON |
| Monk | d8 | DEX / WIS |
| Paladin | d10 | STR / CHA |
| Ranger | d10 | DEX / WIS |
| Rogue | d8 | DEX / INT |
| Sorcerer | d6 | CHA |
| Warlock | d8 | CHA |
| Wizard | d6 | INT |
| Artificer | d8 | INT |

> Artificer es de Tasha's Cauldron of Everything, incluido por decisión explícita de producto.

**Competencia de armadura y AC por clase (MVP):**

| Clase | Competencia máxima | Armadura MVP (Nivel 1) | Lógica de AC |
|---|---|---|---|
| Wizard / Sorcerer | `none` | Ropas (10) | `10 + DEX` |
| Rogue / Bard | `light` | Cuero (11) | `11 + DEX` |
| Monk | `none` | Unarmored (10) | `10 + DEX + WIS` |
| Warlock | `light` | Cuero (11) | `11 + DEX` |
| Ranger / Druid | `medium` | Camisote (13) | `13 + min(DEX, 2)` |
| Cleric / Artificer | `medium` | Camisote (13) | `13 + min(DEX, 2)` |
| Barbarian | `medium` | Unarmored (10) | `10 + DEX + CON` |
| Fighter / Paladin | `heavy` | Cota de malla (16) | `16` (fijo) |

> Monk y Barbarian usan Unarmored Defense — fórmula propia, sin armadura asignada.

**Modelo de armadura:**

```typescript
type ArmorCategory = 'none' | 'light' | 'medium' | 'heavy' | 'shield'

interface ArmorType {
  name: string
  category: ArmorCategory
  baseAC: number
  maxDex?: number  // undefined → sin límite (light/none), 2 → medium, 0 → heavy
  strRequirement?: number
}
```

La lógica de AC vive en una función central — no en cada clase:

```typescript
function calculateAC(armor: ArmorType, modifiers: Modifiers): number {
  switch (armor.category) {
    case 'heavy':  return armor.baseAC
    case 'medium': return armor.baseAC + Math.min(modifiers.DEX, 2)
    default:       return armor.baseAC + modifiers.DEX
  }
}
```

Casos especiales (Unarmored Defense):
```typescript
// Barbarian
AC = 10 + modifiers.DEX + modifiers.CON

// Monk
AC = 10 + modifiers.DEX + modifiers.WIS
```

**Método de generación de stats:**
- Baseline por clase (array [STR, DEX, CON, INT, WIS, CHA]) → `baseStats`
- Variación controlada ±1/±2 sobre el baseline
- Los ASIs del background (ASI pool) se aplican vía `resolveAbilityBonuses` → `finalStats`
- Las species NO otorgan ASIs en 5.5e — solo rasgos raciales
- Coherencia sobre aleatoriedad pura

**Baselines por clase:**

| Clase | STR | DEX | CON | INT | WIS | CHA | Lógica |
|---|---|---|---|---|---|---|---|
| Barbarian | 15 | 13 | 14 | 8 | 10 | 12 | Fuerza y Aguante |
| Bard | 8 | 14 | 12 | 10 | 13 | 15 | Carisma máximo, buena agilidad |
| Cleric | 14 | 8 | 13 | 10 | 15 | 12 | Sabiduría para conjuros, Fuerza para armadura |
| Druid | 10 | 12 | 14 | 13 | 15 | 8 | Sabiduría pura y buena Constitución |
| Fighter | 15 | 13 | 14 | 10 | 12 | 8 | El estándar físico (STR/CON) |
| Monk | 10 | 15 | 13 | 8 | 14 | 12 | Destreza y Sabiduría (AC sin armadura) |
| Paladin | 15 | 8 | 13 | 10 | 12 | 14 | Fuerza y Carisma (Aura/Conjuros) |
| Ranger | 12 | 15 | 13 | 10 | 14 | 8 | Destreza y supervivencia (WIS) |
| Rogue | 8 | 15 | 12 | 14 | 13 | 10 | Destreza máxima e Inteligencia |
| Sorcerer | 8 | 13 | 14 | 12 | 10 | 15 | Carisma innato y vida (CON) |
| Warlock | 10 | 14 | 13 | 12 | 8 | 15 | Carisma y Destreza para defensa ligera |
| Wizard | 8 | 13 | 14 | 15 | 12 | 10 | Inteligencia máxima |
| Artificer | 10 | 12 | 14 | 15 | 13 | 8 | Inteligencia y herramientas (CON/WIS) |

Ejemplo con variación:
```
Fighter baseline: [15, 13, 14, 10, 12, 8]
Fighter generado: [15, 14, 13, 10, 12, 8]  // DEX+1, CON-1
```

**Fórmulas (nivel 1, fijas para el MVP):**
```
// Se aplican siempre sobre finalStats, nunca sobre baseStats
modifier  = ⌊(finalStat - 10) / 2⌋
HP        = hit_die_clase + modifiers.CON
AC        = calculateAC(classDefaultArmor, modifiers)
```

**Seed:**
- Parámetro opcional en el request
- Si viene → resultado reproducible
- Si no viene → random
- Habilita testing y debugging sin infraestructura extra

**MVP genera solo personajes de nivel 1.**
Diseñado para escalar sin reescribir:
- `level` como campo explícito en `Character`
- `baseStats` almacenados antes de bonos, `finalStats` después
- Funciones puras: `calculateHP(level, class, modifiers.CON)`,
  `calculateProficiency(level)`, `calculateAC(armor, modifiers)`

Features de clase, spell slots, ASI y progresión completa:
extensibles en el diseño, **no implementados en el MVP**.

### Módulo 3 — Generador de Ganchos de Historia
Genera tres `NarrativeBlock` en español, uno por categoría:
- **background** — origen, vida anterior al aventurerismo
- **motivation** — qué lo mueve, qué busca
- **secret** — qué oculta, su sombra

Cada bloque incluye `content` (texto generado) y `tags` que determinan
con qué clases y species es compatible ese template.

**La narrativa recibe `class` y `species` como contexto** — los templates
se filtran por compatibilidad antes de seleccionar. Esto garantiza coherencia
narrativa incluso en generaciones 100% aleatorias.

Ejemplo de template con tags:
```typescript
{
  category: 'background',
  content: 'Pasó sus años de juventud en una biblioteca estudiando textos arcanos',
  tags: ['wizard', 'sorcerer', 'artificer']  // incompatible con barbarian
}

{
  category: 'background',
  content: 'Creció en las fronteras, aprendiendo a sobrevivir en lo salvaje',
  tags: ['barbarian', 'ranger', 'druid']
}

{
  category: 'background',
  content: 'Fue rechazado por su comunidad y aprendió a valerse solo',
  tags: ['tiefling', 'half-orc', 'any']  // 'any' = compatible con todo
}
```

El tag `any` marca templates universales — válidos para cualquier clase o species.
Cuantos más tags específicos tenga un template, más coherente y distintivo
será el personaje resultante.

### Módulo 4 — Creador de Personaje (orquestador)
Ejecuta el pipeline completo de generación en orden:

```
1. Resolver inputs     → usar parámetros provistos o elegir aleatoriamente
2. Generar baseStats   → según clase resuelta
3. Aplicar variación   → variación controlada ±1/±2
4. Resolver bonos      → background ASI pool (+2/+1 standard O +1/+1/+1 spread)
5. Aplicar bonos       → finalStats (desde el background ASI pool, no species)
6. Calcular modifiers  → sobre finalStats
7. Resolver armadura   → classDefaultArmor según competencia de clase
8. Calcular derived    → HP = hit_die + modifiers.CON
                      → AC = calculateAC(armor, modifiers)
9. Generar narrativa   → background, motivation, secret (NarrativeBlock)
```

- Los inputs no provistos se resuelven aleatoriamente en el paso 1
- La species determina el nombre (Módulo 1) y los bonos vía pipeline
- La clase determina el baseline de stats, el hit die y la armadura por defecto
- Output: entidad `Character` completa, lista para usar en mesa

**Estado de bloqueo para regeneración parcial:**

```typescript
interface CharacterLocks {
  name: boolean
  stats: boolean
  background: boolean
  motivation: boolean
  secret: boolean
}
```

El orquestador respeta los locks al regenerar — los campos bloqueados
se preservan, los desbloqueados se regeneran. Esto habilita la regeneración
parcial del Módulo 5 sin lógica adicional.

### Módulo 5 — Gestión de personajes
- Guardar personaje/NPC generado (asociado al usuario)
- Listar personajes guardados
- Abrir y releer un personaje guardado
- Regeneración parcial: nombre, stats o historia de forma independiente
- Edición básica de cualquier campo generado

### Módulo 6 — Autenticación (Magic Link)
- El usuario se registra/inicia sesión solo con su email
- El backend envía un link de un solo uso con expiración
- Sin contraseñas
- Base limpia para agregar auth fuerte y suscripciones en el futuro

---

## Modelo de dominio

### Entidad principal: Character

```typescript
interface Character {
  id: string

  // Identidad
  name: string
  species: Species
  subSpecies?: SubSpecies
  class: Class
  level: number

  // Configuración de reglas
  ruleset: string             // siempre "5.5e" — fijo tras el pivote
  abilityBonusSource: string  // siempre "background" en 5.5e
  asiDistribution: string     // "standard" (+2/+1) | "spread" (+1/+1/+1)

  // Núcleo mecánico
  baseStats: Stats       // antes de bonos
  finalStats: Stats      // después de bonos de background
  modifiers: Modifiers
  derived: DerivedStats  // HP, AC

  // Narrativa
  background: NarrativeBlock
  motivation: NarrativeBlock
  secret: NarrativeBlock

  // Metadata
  seed?: number
  createdAt: Date
  updatedAt: Date
}
```

### Stats y derivados

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

interface NarrativeBlock {
  category: 'background' | 'motivation' | 'secret'
  content: string
  tags: string[]  // clases y species compatibles; 'any' = universal
}

type ArmorCategory = 'none' | 'light' | 'medium' | 'heavy' | 'shield'

interface ArmorType {
  name: string
  category: ArmorCategory
  baseAC: number
  maxDex?: number        // undefined → sin límite, 2 → medium, 0 → heavy
  strRequirement?: number
}

interface CharacterLocks {
  name: boolean
  stats: boolean
  background: boolean
  motivation: boolean
  secret: boolean
}
```

### Bonos desacoplados del origen

```typescript
interface AbilityBonus {
  stat: keyof Stats
  value: number
  source: 'background'  // en 5.5e el único origen de ASIs es el background
}
```

Los bonos no viven dentro de la raza ni de la clase — son un resultado
de resolución del background. Esto mantiene el pipeline limpio y extensible.

### Background (mecánico — 5.5e)

```typescript
interface Background {
  name: string
  asiPool: [string, string, string]  // 3 stats elegibles para el bonus
  originFeat: string                 // feat fijo otorgado por este background
  tags: string[]                     // coherencia por clase/species; 'any' = universal
}
```

### Pipeline de generación

```
1. Resolver inputs     → usar parámetros provistos o elegir aleatoriamente
2. Generar baseStats   → según clase resuelta
3. Aplicar variación   → variación controlada ±1/±2
4. Resolver bonos      → background ASI pool (+2/+1 standard O +1/+1/+1 spread)
5. Aplicar bonos       → finalStats (desde el background ASI pool, no species)
6. Calcular modifiers  → ⌊(stat - 10) / 2⌋ sobre finalStats
7. Resolver armadura   → classDefaultArmor según competencia de clase
8. Calcular derived    → HP = hit_die + modifiers.CON
                      → AC = calculateAC(armor, modifiers)
9. Generar narrativa   → background, motivation, secret (NarrativeBlock)
                         filtrados por class y species resueltos en paso 1
```

### Resolución de bonos

```typescript
function resolveAbilityBonuses(background: Background, distribution: 'standard' | 'spread'): AbilityBonus[] {
  if (distribution === 'spread') {
    return background.asiPool.map(stat => ({ stat, value: 1, source: 'background' }))
  }
  return [
    { stat: background.asiPool[0], value: 2, source: 'background' },
    { stat: background.asiPool[1], value: 1, source: 'background' },
  ]
}
```

El MVP opera con `ruleset: '5.5e'` y `abilityBonusSource: 'background'`.

### Anti-patrones a evitar

```typescript
// ❌ Bonos hardcodeados dentro de la raza
Elf = { bonus: { DEX: +2 } }

// ❌ Mezclar base con final
stats = generateStatsWithBonuses()

// ❌ Calcular modifiers antes de aplicar bonos
modifiers = calculateModifiers(baseStats)  // debe ser sobre finalStats
```

---

## Fuera del scope del MVP
- Auth con contraseña o OAuth (Google, Discord, etc.)
- Suscripciones o planes de pago
- Subclases (Paladin Oath, Wizard School, etc.)
- Equipamiento e inventario
- Hechizos y spell slots
- Rasgos y habilidades de clase detallados
- Multiclase
- Generador de quests o ciudades
- Sistema de combate
- Integraciones externas (VTT, PDF, etc.)

---

## Roadmap post-MVP
1. Progresión de nivel (level up de personajes existentes, niveles 2–20)
2. Auth fuerte (OAuth, 2FA, recuperación de cuenta)
3. Modelo de suscripción (límites por plan, features premium)
4. Generador de quests / ganchos de campaña
5. Equipamiento e inventario coherente por clase
6. Exportación (PDF, integración con VTT)

---

## Riesgos identificados
| Riesgo | Impacto | Mitigación |
|---|---|---|
| Contenido repetitivo | Alto | Combinar templates + variaciones en datos |
| Mala calidad narrativa en español | Alto | Iterar fuerte en seed data |
| Scope creep (ciudades, items, etc.) | Medio | Este documento como ancla |
| Latencia si se usa IA | Medio | Fallback con templates determinísticos |
| Sobreingeniería en el motor | Medio | Empezar simple, crecer con tests |

---

## Supuestos de producto
- Los DMs y jugadores valoran velocidad sobre profundidad extrema
- El español nativo es una ventaja competitiva real
- Un personaje usable tiene más valor que 10 ideas vagas
- Los usuarios quieren poder editar, no solo aceptar lo generado
- El flujo principal es: generar → leer → ajustar → guardar → usar

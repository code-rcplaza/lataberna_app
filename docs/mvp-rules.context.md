# mvp-rules.md — rpg_engine

Datos operativos del MVP: clases, species, backgrounds, baselines, armaduras y fórmulas.
Reglas basadas en D&D 5.5e (2024 Player's Handbook).
Este archivo cambia cuando se agregan clases, species, backgrounds o se ajustan valores.

---

## Clases (13)

| Clase | Hit Die | Competencia armadura | Armadura MVP | Lógica AC |
|---|---|---|---|---|
| Barbarian | d12 | medium | Unarmored | `10 + DEX + CON` |
| Bard | d8 | light | Cuero (11) | `11 + DEX` |
| Cleric | d8 | medium | Camisote (13) | `13 + min(DEX, 2)` |
| Druid | d8 | medium | Camisote (13) | `13 + min(DEX, 2)` |
| Fighter | d10 | heavy | Cota de malla (16) | `16` (fijo) |
| Monk | d8 | none | Unarmored | `10 + DEX + WIS` |
| Paladin | d10 | heavy | Cota de malla (16) | `16` (fijo) |
| Ranger | d10 | medium | Camisote (13) | `13 + min(DEX, 2)` |
| Rogue | d8 | light | Cuero (11) | `11 + DEX` |
| Sorcerer | d6 | none | Ropas (10) | `10 + DEX` |
| Warlock | d8 | light | Cuero (11) | `11 + DEX` |
| Wizard | d6 | none | Ropas (10) | `10 + DEX` |
| Artificer | d8 | medium | Camisote (13) | `13 + min(DEX, 2)` |

> Artificer es de Tasha's Cauldron of Everything — incluido por decisión explícita de producto.
> Monk y Barbarian usan Unarmored Defense — sin ArmorType asignado, fórmula propia.

---

## Baselines de stats por clase

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

Variación aplicada sobre el baseline: ±1/±2 controlado.

---

## Species (9 con subSpecies)

| Species | SubSpecies |
|---|---|
| Human | — |
| Elf | High Elf, Wood Elf, Dark Elf (Drow) |
| Dwarf | Hill Dwarf, Mountain Dwarf |
| Halfling | Lightfoot, Stout |
| Gnome | Forest Gnome, Rock Gnome |
| Half-Elf | — |
| Half-Orc | — |
| Tiefling | infernal names, virtue names |
| Dragonborn | — (clan name primero) |

Mínimo 50 entradas de nombres por género/species en la seed.

> **5.5e:** Las species NO otorgan bonos de estadísticas (ASIs).
> Solo proveen rasgos y habilidades raciales. Los ASIs provienen exclusivamente del background.

---

## Backgrounds (ASI pool — 5.5e)

Cada background define un `asiPool` de 3 stats y un `OriginFeat` fijo.
El generador aplica una distribución "standard" (+2/+1 a dos stats del pool)
o "spread" (+1/+1/+1 a los tres).

Ejemplos de backgrounds MVP:

| Background | ASI Pool | Origin Feat | Tags |
|---|---|---|---|
| Acolyte | WIS, INT, CON | Magic Initiate (Cleric) | cleric, paladin, druid |
| Criminal | DEX, INT, CHA | Alert | rogue, warlock, any |
| Folk Hero | STR, CON, WIS | Tough | fighter, barbarian, ranger |
| Noble | STR, INT, CHA | Skilled | paladin, bard, warlock |
| Sage | INT, WIS, CHA | Magic Initiate (Wizard) | wizard, sorcerer, artificer |
| Soldier | STR, DEX, CON | Savage Attacker | fighter, barbarian, paladin |

> Esta tabla es orientativa para el MVP. La seed data de backgrounds debe incluir suficientes
> entradas con tags variados para garantizar coherencia en cualquier combinación clase/species.

---

## Fórmulas (nivel 1)

```
// Siempre sobre finalStats
modifier = ⌊(finalStat - 10) / 2⌋

// HP
HP = hit_die_clase + modifiers.CON

// AC — función central
calculateAC(armor, modifiers):
  heavy  → armor.baseAC
  medium → armor.baseAC + min(modifiers.DEX, 2)
  other  → armor.baseAC + modifiers.DEX

// Unarmored Defense
Barbarian → 10 + modifiers.DEX + modifiers.CON
Monk      → 10 + modifiers.DEX + modifiers.WIS
```

---

## Seed

- Parámetro opcional en el request
- Con seed → resultado reproducible
- Sin seed → random
- Mismo seed + mismos parámetros = mismo resultado siempre

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**FORGE RPG / La Taberna RPG** — D&D 5.5e (2024 Player's Handbook) character generator that produces complete, coherent, narrative-rich characters in Spanish. Core value: generate a full PC or NPC in seconds, with optional partial regeneration and a personal library.

**Status:** Specification-complete. Source code not yet initialized. All architecture decisions are resolved in `docs/`.

---

## Stack

| Layer | Technology |
|---|---|
| Backend | Go — clean architecture |
| API | GraphQL via `gqlgen` |
| Frontend | Vue or Angular (TBD) |
| Database | SQLite + Atlas HCL migrations (PostgreSQL post-MVP) |
| Auth | Magic link (email + single-use token, no passwords) |
| Mobile | Flutter (post-MVP) |

---

## Commands (once initialized)

```bash
# Backend
go run ./cmd/server       # dev server
go test ./...             # all tests
go test ./internal/...    # domain + usecase tests only
go generate ./...         # gqlgen code generation from schema

# Database
atlas migrate apply       # run pending migrations
atlas migrate diff        # generate migration from schema diff

# Frontend (pnpm preferred)
pnpm dev                  # dev server
pnpm test                 # unit tests
pnpm build                # production build
```

---

## Architecture

Clean architecture with strict layer isolation. Dependencies point inward only.

```
cmd/                  → entrypoints (server, migrate, seed)
internal/
  domain/             → entities, value objects, interfaces (no external deps)
  usecase/            → orchestrators, business logic (depends on domain only)
  infrastructure/     → DB, email, GraphQL resolvers (implements domain interfaces)
docs/                 → specifications (CONTEXT.md is the source of truth)
```

### Key patterns
- **Builder** — character generation with optional parameters
- **Factory** — `ArmorType` and `NarrativeBlock` creation
- **Prototype** — partial regeneration: clone → re-execute unlocked fields
- **Repository** — persistence abstraction; domain never touches DB directly

---

## Character Generation Pipeline (9 steps — immutable order)

```
1. resolve inputs       → use provided params or choose randomly
2. generate baseStats   → per-class baseline
3. apply variation      → ±1/±2 controlled variation
4. resolve bonuses      → background ASI pool (+2/+1 standard OR +1/+1/+1 spread)
5. apply bonuses        → produce finalStats (from background ASI pool, not species)
6. calculate modifiers  → ⌊(finalStat − 10) / 2⌋ on finalStats (NOT baseStats)
7. resolve armor        → class default ArmorType
8. calculate derived    → hp = hit_die + modifiers.CON
                          ac = calculateAC(armor, modifiers)
9. generate narrative   → background, motivation, secret (filtered by class + species)
```

Modifiers are **always** computed from `finalStats`, never `baseStats`. This is a hard invariant.

---

## Domain Model (core types)

```go
type Character struct {
    ID               string
    Name             string
    Species          Species
    SubSpecies       *SubSpecies
    Class            Class
    Level            int
    Ruleset          string           // always "5.5e" — fixed after pivot; kept for save/load fidelity
    AbilityBonusSource string         // always "background" in 5.5e — species no longer provide ASIs
    ASIDistribution  string           // "standard" (+2/+1) | "spread" (+1/+1/+1)
    BaseStats        Stats
    FinalStats       Stats
    Modifiers        Modifiers
    Derived          DerivedStats     // { HP, AC }
    Background       NarrativeBlock
    Motivation       NarrativeBlock
    Secret           NarrativeBlock
    Locks            CharacterLocks
    Seed             *int64
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

type Stats     struct { STR, DEX, CON, INT, WIS, CHA int }
type Modifiers struct { STR, DEX, CON, INT, WIS, CHA int }
type DerivedStats struct { HP, AC int }

type NarrativeBlock struct {
    Category string   // "background" | "motivation" | "secret"
    Content  string
    Tags     []string // class/species compatibility tags; "any" = universal — same filtering applies to Background selection
}

type Background struct {
    Name        string
    ASIPool     [3]string  // e.g., ["WIS", "INT", "CON"] — defined per background; player picks 2 stats to boost
    OriginFeat  string     // fixed feat granted by this background (not chosen by player — simplifies generation)
    Tags        []string   // class/species coherence tags; "any" = universal
}

type ArmorType struct {
    Name           string
    Category       string  // "none" | "light" | "medium" | "heavy" | "shield"
    BaseAC         int
    MaxDex         *int    // nil = unlimited, 2 = medium, 0 = heavy
    StrRequirement *int
}

type CharacterLocks struct {
    Name       bool
    Stats      bool
    Background bool
    Motivation bool
    Secret     bool
}
```

---

## Architectural Decisions (non-negotiable without explicit discussion)

1. **All parameters optional** — omitted fields generate randomly; never require input.
2. **Ruleset fixed to "5.5e"** — pivoted from 5e; stronger market position for Spanish-language content targeting the 2024 PHB audience. Field kept on Character for save/load fidelity.
3. **Backgrounds provide ASIs, not species** — per 5.5e rules; background defines an ASIPool of 3 stats, player (or generator) picks a distribution. Species provide racial traits/abilities only.
4. **ASIDistribution on Character** — "standard" (+2/+1) or "spread" (+1/+1/+1); stored on the entity (not generator config) so saved characters reload with identical stat allocation.
5. **OriginFeat fixed per background** — the feat granted is determined by the background definition, not chosen by the player; this simplifies generation and keeps the engine deterministic.
6. **Separate baseStats / finalStats** — enables future level scaling and multiclass without rewrites.
7. **Modifiers on finalStats** — hard invariant in the pipeline.
8. **Narrative filtered by class + species** — coherence is guaranteed, not optional. Same Tags filtering applies to Background selection.
9. **ArmorType as a resource** — not embedded in class logic; enables future inventory.
10. **Magic link auth** — no passwords; tokens: 32+ byte entropy, 15min TTL, one-time use, hashed at rest.
11. **SQLite in MVP** — migrate to Postgres post-validation, Atlas handles migration diffing.
12. **Level 1 only in MVP** — engine designed to scale to 1–20, not implemented yet.
13. **Seed for reproducibility** — same seed + same params = same result; enables deterministic testing.

---

## Coding Conventions

### Early returns (mandatory)
Validate preconditions first; no nested else blocks.

```go
// ✅
func (b *Builder) Build() (*Character, error) {
    if b.class == "" {
        return nil, errors.New("class required")
    }
    // happy path
}

// ❌
func (b *Builder) Build() (*Character, error) {
    if b.class != "" {
        // nested ...
    }
}
```

### Testing
- **TDD is mandatory** for domain logic and usecases — write tests before implementation.
- Table-driven tests with named subtests in Go.
- Minimum 80% coverage on `internal/domain/` and `internal/usecase/`.
- Infrastructure (DB, email) and GraphQL handlers: tests after implementation.
- E2E for full generate → show → save flows.

---

## Content Data

- **Classes (13):** Barbarian, Bard, Cleric, Druid, Fighter, Monk, Paladin, Ranger, Rogue, Sorcerer, Warlock, Wizard, Artificer
- **Species (9 with sub-species):** Human, Elf (High/Wood/Drow), Dwarf (Hill/Mountain), Halfling (Lightfoot/Stout), Gnome (Forest/Rock), Half-Elf, Half-Orc, Tiefling, Dragonborn
- **Species no longer provide ability score bonuses** — per 5.5e rules, ASIs come exclusively from backgrounds. Species data contains racial traits and abilities only.
- Minimum 50 name entries per gender per species in seed data.
- New classes/species must be **data additions only** — no logic changes required.

---

## Specs Reference

All requirements and operational rules live in `docs/`:

| File | Content |
|---|---|
| `CONTEXT.md` | Stable project context, pipeline definition, all architecture decisions |
| `product-scope.context.md` | MVP scope, 6 modules, problem statement |
| `functional_requirements.md` | 24 RFs (RF-00-001 through RF-05-005) |
| `mvp-rules.context.md` | Class baselines, species bonuses, armor table, AC formulas |
| `domain-model.context.md` | Full type definitions with examples |
| `non_functional_requirements.md` | Performance (<500ms full, <200ms partial), security, observability |

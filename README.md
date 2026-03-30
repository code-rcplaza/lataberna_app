> # La Taberna RPG

Generador de personajes D&D 5e que produce personajes completos y narrativamente ricos en español. Backend en Go con arquitectura limpia, API GraphQL, base de datos SQLite y frontend Vue 3.

---

## Stack

| Capa            | Tecnología                                                        |
| --------------- | ----------------------------------------------------------------- |
| Backend         | Go 1.25 — arquitectura limpia (domain / usecase / infrastructure) |
| API             | GraphQL via `gqlgen`                                              |
| Base de datos   | SQLite3 + migraciones Atlas HCL                                   |
| Email           | Resend (magic link auth)                                          |
| Frontend        | Vue 3 + TypeScript + Pinia + Vue Router                           |
| CSS             | Tailwind CSS                                                      |
| Build frontend  | Vite                                                              |
| Deploy backend  | Railway (Docker)                                                  |
| Deploy frontend | Vercel (SPA rewrite)                                              |

---

## Estructura del proyecto

```
FORGE_RPG/
├── backend/
│   ├── cmd/server/           # Entrypoint — main.go
│   ├── internal/
│   │   ├── domain/           # Entidades, value objects, interfaces (sin deps externas)
│   │   │   └── ports/        # Interfaces de repositorios y servicios
│   │   ├── usecase/          # Lógica de negocio
│   │   │   ├── character/    # Pipeline de generación (9 pasos)
│   │   │   ├── statblock/    # Cálculo de stats, HP y AC
│   │   │   ├── namegen/      # Generación de nombres por especie
│   │   │   ├── narrativegen/ # Generación de narrativa ponderada
│   │   │   └── auth/         # Magic link authentication
│   │   └── infrastructure/
│   │       ├── db/           # Repositorios SQLite + seed data
│   │       ├── email/        # Integración Resend
│   │       └── graphql/      # Middleware de autenticación
│   ├── graph/                # Schema GraphQL + resolvers (generados con gqlgen)
│   ├── migrations/           # Schema Atlas HCL
│   ├── Dockerfile
│   └── go.mod
│
├── frontend/
│   ├── src/
│   │   ├── components/       # CharacterSheet, ConfigPanel, LockButton, StatCard...
│   │   ├── views/            # AuthView, CharacterDetailView, BibliotecaView
│   │   ├── stores/           # Pinia: auth, character, library
│   │   ├── composables/      # useGraphQL, useAuthAPI, useCharacterAPI
│   │   └── types/            # Tipos TypeScript
│   ├── vercel.json           # SPA rewrite para Vue Router
│   └── vite.config.ts
│
├── docs/                     # Especificaciones del dominio
│   ├── CONTEXT.md            # Contexto estable + decisiones arquitectónicas
│   ├── functional_requirements.md
│   ├── mvp-rules.context.md  # Baselines por clase, bonuses por especie, fórmulas AC
│   └── domain-model.context.md
│
└── Makefile
```

---

## Comandos

### Backend

```bash
# Desarrollo
go run ./cmd/server              # Servidor en :8080

# Tests
go test ./...                    # Todos los paquetes
go test ./internal/...           # Solo domain + usecase + infrastructure

# Generación de código GraphQL
go generate ./...                # Regenera desde schema.graphqls

# Migraciones
atlas schema apply               # Aplicar migraciones pendientes
atlas schema diff                # Generar diff desde esquema actual
```

### Frontend

```bash
pnpm dev                         # Servidor de desarrollo en :5173
pnpm build                       # Build de producción → dist/
```

### Makefile (desde la raíz)

```bash
make dev-backend                 # go run ./cmd/server
make dev-frontend                # pnpm dev
make test                        # go test ./...
make build                       # Backend + frontend
make generate                    # gqlgen code generation
```

---

## Variables de entorno

### Backend (`.env`)

| Variable         | Descripción                                              | Default                              |
| ---------------- | -------------------------------------------------------- | ------------------------------------ |
| `DB_PATH`        | Ruta a la base de datos SQLite                           | `forge.db`                           |
| `PORT`           | Puerto del servidor                                      | `8080`                               |
| `CORS_ORIGIN`    | Origen permitido por CORS                                | `http://localhost:5173`              |
| `LINK_BASE`      | URL base para verificación del magic link                | `http://localhost:8080/auth/verify`  |
| `RESEND_API_KEY` | API key de Resend (opcional — sin ella imprime a stdout) | —                                    |
| `RESEND_FROM`    | Header From de los emails                                | `La Taberna <noreply@lataberna.app>` |

### Frontend (`.env`)

| Variable           | Descripción                  |
| ------------------ | ---------------------------- |
| `VITE_GRAPHQL_URL` | Endpoint GraphQL del backend |

---

## Arquitectura

### Capas (las dependencias apuntan hacia adentro)

```
graph/ (resolvers)
    ↓
usecase/ (lógica de negocio)
    ↓
domain/ (entidades + ports)
    ↑
infrastructure/ (implementaciones: DB, email)
```

### Modelo de dominio (tipos principales)

```go
type Character struct {
    ID                 string
    Name               string
    Species            Species
    SubSpecies         *SubSpecies
    Class              Class
    Level              int
    Ruleset            Ruleset           // "5e" | "5.5e"
    AbilityBonusSource string
    BaseStats          Stats
    FinalStats         Stats
    Modifiers          Modifiers
    Derived            DerivedStats      // { HP, AC }
    Background         NarrativeBlock
    Motivation         NarrativeBlock
    Secret             NarrativeBlock
    Locks              CharacterLocks
    Seed               *int64
}

type Stats        struct { STR, DEX, CON, INT, WIS, CHA int }
type Modifiers    struct { STR, DEX, CON, INT, WIS, CHA int }
type DerivedStats struct { HP, AC int }
```

---

## Pipeline de generación de personajes (9 pasos — orden inmutable)

```
1. resolve inputs       → usa parámetros provistos o elige aleatoriamente
2. generate baseStats   → baseline por clase
3. apply variation      → ±1/±2 variación controlada
4. resolve bonuses      → bonuses por especie/subespecie
5. apply bonuses        → produce finalStats
6. calculate modifiers  → ⌊(finalStat − 10) / 2⌋ sobre finalStats (invariante)
7. resolve armor        → ArmorType default por clase
8. calculate derived    → HP = hit_die + MOD.CON
                          AC = calculateAC(armor, modifiers)
9. generate narrative   → trasfondo, motivación, secreto (filtrados por clase + especie)
```

> **Invariante**: Los modificadores se calculan SIEMPRE sobre `finalStats`, nunca sobre `baseStats`.

---

## Generación de nombres

Sistema de composición multi-componente por especie. Cada especie tiene reglas de ensamblado propias:

| Especie                   | Formato                                      |
| ------------------------- | -------------------------------------------- |
| Humano                    | Nombre + Apellido                            |
| Enano (Hill/Mountain)     | Nombre + Nombre de clan                      |
| Elfo (High/Wood/Drow)     | Nombre + Apellido familiar                   |
| Mediano (Lightfoot/Stout) | Nombre + Apellido                            |
| Gnomo (Forest/Rock)       | Nombre + Clan + "Apodo"                      |
| Semielfo                  | Nombre + Apellido (humano o élfico, al azar) |
| Semiorco                  | Nombre [+ Apellido con 30% de probabilidad]  |
| Dragonborn                | Clan + Nombre (el clan va primero)           |
| Tiefling Infernal         | Nombre infernal (componente único)           |
| Tiefling Virtud           | Palabra de virtud en español                 |

Todo el contenido está en español. Los datos se almacenan en SQLite (`name_entries`) con columnas `species_key`, `gender`, `name_type` y `name`.

---

## Generación de narrativa

Las entradas narrativas (trasfondo, motivación, secreto) se filtran y ponderan por clase y especie:

- **Primario** (`primary`): mayor probabilidad de selección
- **Secundario** (`secondary`): probabilidad media
- **Excluido** (`excluded`): nunca se selecciona para esa combinación
- **Universal** (`any`): disponible para todos

La compatibilidad se almacena en `narrative_compatibility` con dimensiones `class` y `species`.

---

## Autenticación (Magic Link)

Sin contraseñas. El flujo completo:

1. Usuario ingresa su email → `requestMagicLink(email)`
2. Backend genera token de 32 bytes de entropía, lo hashea y almacena
3. Se envía email con link de verificación (TTL: 15 minutos, uso único)
4. Usuario hace click → `verifyMagicLink(token)` → backend crea sesión
5. Sesión se usa para autenticar requests posteriores

---

## API GraphQL

### Queries

```graphql
me: User                              # Usuario autenticado
characters: [Character!]!             # Personajes del usuario
character(id: ID!): Character         # Personaje por ID
```

### Mutations

```graphql
# Auth
requestMagicLink(email: String!): Boolean!
verifyMagicLink(token: String!): AuthPayload!
logout: Boolean!

# Generación (sin autenticación requerida)
generateCharacter(input: GenerateCharacterInput): Character!
regenerateDraft(current: CurrentCharacterInput!, locks: CharacterLocksInput!, seed: Int): Character!

# Gestión (requiere autenticación)
regenerateCharacter(id: ID!, locks: CharacterLocksInput!, seed: Int): Character!
saveCharacter(input: GenerateCharacterInput!, seed: Int!): Character!
editCharacter(id: ID!, input: EditCharacterInput!): Character!
deleteCharacter(id: ID!): Boolean!
```

### Regeneración parcial (sistema de locks)

```graphql
type CharacterLocks {
  name: Boolean!
  stats: Boolean!
  background: Boolean!
  motivation: Boolean!
  secret: Boolean!
}
```

Los campos bloqueados se preservan. Los desbloqueados se regeneran. Mismo `seed` + mismo input = mismo resultado.

---

## Base de datos

SQLite3 con WAL journal mode y foreign keys habilitadas. Migraciones versionadas en `migrate()` al startup.

| Tabla                     | Propósito                                           |
| ------------------------- | --------------------------------------------------- |
| `users`                   | Cuentas de usuario (email)                          |
| `sessions`                | Sesiones activas                                    |
| `magic_link_tokens`       | Tokens de auth (hasheados, TTL 15min)               |
| `characters`              | Biblioteca de personajes (stats como JSON)          |
| `narrative_entries`       | Bloques narrativos (trasfondo, motivación, secreto) |
| `narrative_compatibility` | Tags de compatibilidad clase/especie                |
| `name_entries`            | Nombres por especie, género y tipo                  |
| `seed_version`            | Control de versión del seed data                    |

---

## Contenido disponible

- **13 clases**: Bárbaro, Bardo, Clérigo, Druida, Guerrero, Monje, Paladín, Explorador, Pícaro, Hechicero, Brujo, Mago, Artífice
- **9 especies**: Humano, Elfo (Alto/Bosque/Drow), Enano (Colina/Montaña), Mediano (Lightfoot/Stout), Gnomo (Bosque/Roca), Semielfo, Semiorco, Tiefling, Dragonborn
- **Mínimo 50 nombres** por género por especie
- **Narrativa completa** en español con compatibilidad por clase y especie

---

## Convenciones de código

### Early returns (obligatorio)

```go
// ✅ Correcto
func (b *Builder) Build() (*Character, error) {
    if b.class == "" {
        return nil, errors.New("class required")
    }
    // happy path sin else anidados
}
```

### Testing

- **TDD obligatorio** para domain y usecase — tests antes de implementación
- Tests con tabla y subtests nombrados en Go
- Cobertura mínima 80% en `internal/domain/` y `internal/usecase/`
- Infrastructure y GraphQL: tests después de implementación

### Commits

Formato Conventional Commits. Sin `Co-Authored-By`.

---

## Decisiones arquitectónicas no negociables

1. **Todos los parámetros opcionales** — sin input obligatorio, todo se genera aleatoriamente si se omite
2. **Bonuses desacoplados de la especie** — resueltos en el pipeline, no hardcodeados por especie
3. **baseStats / finalStats separados** — habilita escala de nivel y multiclase sin reescribir
4. **Modificadores sobre finalStats** — invariante absoluta del pipeline
5. **Narrativa filtrada por clase + especie** — coherencia garantizada, no opcional
6. **ArmorType como recurso** — no embebido en la clase, habilita inventario futuro
7. **Magic link auth** — sin contraseñas; tokens 32+ bytes de entropía, TTL 15min, uso único, hasheados en reposo
8. **SQLite en MVP** — migrar a Postgres post-validación; Atlas maneja el diff
9. **Solo nivel 1 en MVP** — motor diseñado para 1-20, no implementado aún
10. **Seed para reproducibilidad** — mismo seed + mismos params = mismo resultado

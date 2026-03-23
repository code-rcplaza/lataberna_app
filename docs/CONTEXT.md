# CONTEXT.md — rpg_engine

Contexto estable del proyecto. Leer completo antes de generar código,
sugerencias o decisiones técnicas. Para detalles operativos ver `docs/`.

---

## Qué es este proyecto

Generador de personajes y NPCs para D&D 5e en español.
Nombre tentativo: **La Taberna RPG**.
Audiencia: jugadores y Dungeon Masters hispanohablantes, nivel principiante a intermedio.

---

## Stack tecnológico

| Capa | Tecnología |
|---|---|
| Backend | Go (clean architecture) |
| Frontend web | Vue o Angular (por decidir) |
| Frontend móvil | Flutter |
| Base de datos | SQLite + Atlas (migraciones en HCL) |
| API | GraphQL (gqlgen) |
| Auth | Magic link (email + token de un solo uso) |

---

## Principio fundamental del generador

> Todo parámetro es opcional. Lo que se omite se genera aleatoriamente.
> Lo que se provee queda locked.

El generador funciona con cero inputs (personaje 100% aleatorio)
o con cualquier combinación de parámetros restringidos.

---

## Pipeline de generación (orden estricto)

```
1. Resolver inputs     → usar parámetros provistos o elegir aleatoriamente
2. Generar baseStats   → según clase resuelta (baseline por clase)
3. Aplicar variación   → variación controlada ±1/±2 sobre el baseline
4. Resolver bonos      → resolveAbilityBonuses(species, subSpecies)
5. Aplicar bonos       → finalStats
6. Calcular modifiers  → ⌊(finalStat - 10) / 2⌋ sobre finalStats
7. Resolver armadura   → classDefaultArmor según competencia de clase
8. Calcular derived    → HP = hit_die + modifiers.CON
                      → AC = calculateAC(armor, modifiers)
9. Generar narrativa   → filtrada por class y species resueltos en paso 1
```

Reglas invariables del pipeline:
- Los modifiers se calculan **siempre** sobre `finalStats`, nunca sobre `baseStats`
- La narrativa **siempre** recibe class y species como contexto de filtrado
- Los bonos **nunca** viven dentro de la species — son un resultado de resolución

---

## Decisiones explícitas (no revertir sin discusión)

| Decisión | Razón |
|---|---|
| Todo parámetro es opcional | Minimizar fricción, maximizar aleatoriedad útil |
| Bonos desacoplados de species | Soportar 5.5e sin reescribir el motor |
| `baseStats` separado de `finalStats` | Escalabilidad de nivel y multiclase |
| Modifiers sobre `finalStats` | Orden correcto del pipeline |
| Narrativa filtrada por class/species | Coherencia en generaciones aleatorias |
| `ArmorType` como recurso, no lógica de clase | DRY, preparado para inventario |
| Magic link sobre sesión anónima | Persistencia real sin perder datos al borrar cookies |
| SQLite en MVP | Simplicidad — migración a Postgres post-validación |
| MVP solo nivel 1 | Foco — el motor escala por diseño, no por implementación |

---

## Prácticas de desarrollo

### Early returns
Toda función valida primero y retorna temprano en caso de error.
No se permiten bloques `else` anidados cuando un `return` es suficiente.

```go
// ✅
func (b *CharacterBuilder) Build() (*Character, error) {
    if b.class == "" {
        return nil, errors.New("class is required")
    }
    if b.species == "" {
        return nil, errors.New("species is required")
    }
    // happy path
}

// ❌
func (b *CharacterBuilder) Build() (*Character, error) {
    if b.class != "" {
        if b.species != "" {
            // happy path
        }
    }
}
```

---

### TDD
El test se escribe **antes** de la implementación — sin excepciones en lógica de negocio.
Se usa table-driven testing en Go para cubrir variantes del pipeline.

```go
func TestCalculateHP(t *testing.T) {
    tests := []struct {
        name   string
        class  Class
        conMod int
        want   int
    }{
        {"fighter con mod 2",   ClassFighter,   2, 12}, // d10 + 2
        {"wizard con mod 0",    ClassWizard,    0, 6},  // d6 + 0
        {"barbarian con mod 3", ClassBarbarian, 3, 15}, // d12 + 3
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CalculateHP(tt.class, tt.conMod)
            if got != tt.want {
                t.Errorf("got %d, want %d", got, tt.want)
            }
        })
    }
}
```

**Estrategia de testing por capa:**

| Capa | Estrategia |
|---|---|
| Lógica de dominio (pipeline, fórmulas) | TDD estricto |
| Usecases / orquestador | TDD estricto |
| Lógica de UI (formateo, validaciones) | TDD estricto |
| Comportamiento de componentes (locks, estado) | TDD estricto |
| Infraestructura (DB, email) | Tests después |
| GraphQL handlers | Tests de integración |
| Flujos de UI (generar → mostrar → guardar) | Tests de integración / E2E |
| CSS, layout, posición visual | Sin tests — revisión manual |
| Narrativa / seed data | Revisión humana |

El valor del test de integración en frontend está en el **comportamiento observable**, no en la implementación visual:

```typescript
// ✅ Valioso — verifica comportamiento
test('al generar sin parámetros, todos los campos del personaje se muestran', async () => {
  await userEvent.click(screen.getByText('Generar'))
  expect(screen.getByTestId('character-name')).not.toBeEmpty()
  expect(screen.getByTestId('character-class')).not.toBeEmpty()
  expect(screen.getByTestId('stat-STR')).not.toBeEmpty()
})

// ❌ No valioso — verifica implementación visual
test('el botón tiene color primario y margen de 16px', () => { ... })
```

---

### Patrones creacionales adoptados

#### Builder — generación de Character
El pipeline de 9 pasos con parámetros opcionales se implementa como Builder.
Parámetro omitido → generado aleatoriamente en `Build()`.

```go
character, err := NewCharacterBuilder().
    WithClass(ClassFighter).
    WithSpecies(SpeciesElf).
    WithSeed(12345).
    Build()
```

#### Factory — ArmorType y NarrativeBlock
Evita switches dispersos. La factory recibe la clase y devuelve el objeto correcto.

```go
armor  := ArmorFactory.ForClass(ClassBarbarian)
blocks := NarrativeFactory.ForClass(ClassWizard, SpeciesElf)
```

#### Prototype — regeneración parcial
La regeneración con locks clona el personaje existente y re-ejecuta
el pipeline solo en los campos desbloqueados.

```go
newChar := existingChar.Clone()
newChar.Regenerate(locks)
```

---



```typescript
// ❌ Bonos dentro de la species
Elf = { bonus: { DEX: +2 } }

// ❌ Mezclar base con final
stats = generateStatsWithBonuses()

// ❌ Modifiers antes de aplicar bonos
modifiers = calculateModifiers(baseStats)

// ❌ Lógica de AC dentro de cada clase
Fighter.calculateAC = () => 16

// ❌ Narrativa sin contexto de clase/species
generateBackground()
```

---

## Fuera del scope del MVP (no sugerir implementar)

- Auth con contraseña, OAuth, 2FA
- Suscripciones o planes de pago
- Subclases (Paladin Oath, Wizard School, etc.)
- Equipamiento e inventario
- Hechizos y spell slots
- Rasgos y habilidades de clase
- Multiclase
- Generador de quests o ciudades
- Sistema de combate
- Exportación PDF / integración VTT
- Progresión de nivel (nivel 2–20)

---

## Documentación relacionada

| Archivo | Contenido |
|---|---|
| `docs/domain-model.md` | Entidades, interfaces y tipos del dominio |
| `docs/mvp-rules.md` | Clases, species, baselines, armaduras, fórmulas |
| `docs/product-scope.md` | Alcance funcional, módulos, roadmap |

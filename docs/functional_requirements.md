# rpg_engine — Requerimientos Funcionales (MVP)

---

## Módulo 0 — Gestión de identidad y acceso

### RF-00-001 — Solicitud de acceso por email
El sistema debe permitir que un usuario solicite acceso ingresando su correo electrónico.

**Criterios de aceptación:**
- El usuario puede ingresar un email válido.
- El sistema valida el formato del email.
- El sistema genera una solicitud de acceso.
- El sistema informa que se envió un enlace de acceso si el email es aceptado.

---

### RF-00-002 — Inicio de sesión mediante magic link
El sistema debe permitir que un usuario inicie sesión mediante un enlace de un solo uso enviado a su correo.

**Criterios de aceptación:**
- El enlace contiene un token único.
- El token expira luego de un tiempo definido.
- El token solo puede usarse una vez.
- Si el token es válido, el usuario inicia sesión.
- Si el token es inválido o expiró, el sistema rechaza el acceso con mensaje claro.

---

### RF-00-003 — Persistencia de sesión
El sistema debe mantener la sesión del usuario autenticado para permitir acceso continuo a su biblioteca de personajes.

**Criterios de aceptación:**
- El usuario no debe autenticarse nuevamente en cada solicitud.
- El sistema puede identificar al usuario autenticado en operaciones protegidas.
- El usuario puede cerrar sesión explícitamente.

---

## Módulo 1 — Generador de Nombres

### RF-01-001 — Todo parámetro es opcional
El generador de nombres funciona con cero parámetros.
Lo que se omite se genera aleatoriamente. Lo que se provee queda locked.

**Criterios de aceptación:**
- El sistema acepta `species`, `subSpecies` y `gender` como parámetros opcionales.
- Si `species` no se provee, el sistema elige una aleatoriamente entre las 9 disponibles.
- Si `subSpecies` no se provee, el sistema elige una válida para la `species` resuelta.
- Si `gender` no se provee, el sistema elige aleatoriamente entre masculino, femenino y neutro.
- Un parámetro provisto produce siempre un nombre coherente con ese valor.

---

### RF-01-002 — Seed para reproducibilidad
Si se provee un `seed`, el nombre generado es reproducible.

**Criterios de aceptación:**
- El sistema acepta `seed` como parámetro opcional.
- El mismo `seed` + mismos parámetros produce siempre el mismo nombre.
- Sin `seed`, el resultado es aleatorio en cada llamada.

---

### RF-01-003 — Convenciones de nomenclatura por species/subSpecies
El nombre generado respeta las convenciones canónicas de D&D 5.5e de cada species y subSpecies.

**Criterios de aceptación:**
- Tiefling genera nombres infernales o de virtud según subSpecies.
- Dragonborn genera el nombre de clan primero.
- Gnome incluye apodo opcional.
- Half-Elf y Half-Orc pueden mezclar convenciones de sus species padre.
- El resto de species respeta sus convenciones canónicas de D&D 5.5e.

---

### RF-01-004 — Variedad suficiente de nombres
La seed data de nombres debe garantizar variedad mínima para evitar repetición perceptible.

**Criterios de aceptación:**
- Mínimo 50 entradas por género/species en la seed.
- La probabilidad de obtener el mismo nombre en dos generaciones consecutivas
  sin seed es estadísticamente baja.

---

### RF-01-005 — Regeneración independiente de nombre
El nombre puede regenerarse sin afectar ningún otro campo del personaje.

**Criterios de aceptación:**
- El sistema expone una operación de regeneración exclusiva para el nombre.
- Los demás campos del personaje permanecen inalterados.
- La regeneración respeta los parámetros originales (species, subSpecies, gender).
- Si el nombre estaba locked, la operación es rechazada.

---

## Módulo 2 — Generador de Stat Block

### RF-02-001 — Todo parámetro es opcional
El generador de stats funciona con cero parámetros.

**Criterios de aceptación:**
- El sistema acepta `class` como parámetro opcional.
- Si `class` no se provee, el sistema elige una aleatoriamente entre las 13 disponibles.
- Una `class` provista produce siempre stats coherentes con ese baseline.

---

### RF-02-002 — Nivel fijo con firma escalable
El sistema genera personajes en nivel 1. Las funciones de cálculo reciben `level`
como parámetro desde el inicio para soportar escalabilidad futura sin cambiar firmas.

**Criterios de aceptación:**
- `level` siempre es `1` en el MVP.
- `level` es un campo explícito en la entidad `Character`.
- `calculateHP(level, class, conMod)`, `calculateProficiency(level)` y
  `calculateAC(level, armor, modifiers)` reciben `level` como parámetro
  aunque en el MVP el valor siempre sea `1`.

---

### RF-02-003 — Generación de baseStats por clase
El sistema genera los 6 stats base aplicando variación controlada sobre el baseline de la clase.

**Criterios de aceptación:**
- Los stats generados parten del baseline definido para la clase resuelta.
- La variación aplicada es de ±1/±2 por stat, de forma controlada.
- El resultado se almacena como `baseStats` — antes de cualquier bono.
- Dos generaciones con el mismo `seed` y clase producen el mismo `baseStats`.

---

### RF-02-004 — Resolución de bonos de background (ASI pool)
El sistema resuelve los ASIs del background seleccionado y los aplica sobre `baseStats` para producir `finalStats`.

**Criterios de aceptación:**
- Los bonos no viven dentro de la species — en 5.5e los ASIs provienen exclusivamente del background.
- El background resuelto define un `asiPool` de 3 stats; el generador elige o aplica una distribución (`standard` o `spread`).
- La distribución se almacena en `Character.asiDistribution`.
- El resultado se almacena como `finalStats` — separado de `baseStats`.
- `finalStats` refleja correctamente los bonos del background ASI pool aplicados sobre `baseStats`.
- El MVP opera con `abilityBonusSource: 'background'`.

---

### RF-02-005 — Cálculo de modifiers
El sistema calcula los 6 modifiers a partir de `finalStats`.

**Criterios de aceptación:**
- Fórmula: `modifier = ⌊(finalStat - 10) / 2⌋`
- Los modifiers se calculan **siempre** sobre `finalStats`, nunca sobre `baseStats`.
- Los modifiers se almacenan en la entidad `Character`.

---

### RF-02-006 — Resolución de armadura por clase
El sistema asigna la armadura por defecto según la competencia de la clase resuelta.

**Criterios de aceptación:**
- Cada clase tiene una armadura por defecto definida en `mvp-rules.md`.
- Monk y Barbarian usan Unarmored Defense — sin ArmorType asignado.
- La armadura resuelta se usa como input del cálculo de AC.

---

### RF-02-007 — Cálculo de HP y AC
El sistema calcula HP y AC usando `level`, modifiers y armadura resuelta.

**Criterios de aceptación:**
- `HP = calculateHP(level, class, modifiers.CON)`
- `AC` se calcula según la categoría de armadura:
  - `heavy`      → `armor.baseAC`
  - `medium`     → `armor.baseAC + min(modifiers.DEX, 2)`
  - `light/none` → `armor.baseAC + modifiers.DEX`
  - Barbarian    → `10 + modifiers.DEX + modifiers.CON`
  - Monk         → `10 + modifiers.DEX + modifiers.WIS`
- HP y AC se almacenan en `derived` dentro de `Character`.
- El cálculo usa los modifiers de `finalStats`, no de `baseStats`.

---

### RF-02-008 — Regeneración independiente de stats
Los stats pueden regenerarse sin afectar nombre ni narrativa del personaje.

**Criterios de aceptación:**
- El sistema expone una operación de regeneración exclusiva para stats.
- Nombre, background, motivation y secret permanecen inalterados.
- La regeneración respeta la `class` original si estaba locked.
- La regeneración re-ejecuta los pasos 2–7 del pipeline completo.
- Si los stats estaban locked, la operación es rechazada.

---

## Módulo 3 — Generador de Ganchos de Historia

### RF-03-001 — Generación de tres NarrativeBlocks
El sistema genera un NarrativeBlock por cada categoría: background, motivation y secret.

**Criterios de aceptación:**
- Cada bloque tiene `category`, `content` y `tags`.
- El `content` está en español.
- Los tres bloques se generan en una sola operación.

---

### RF-03-002 — Seed para reproducibilidad
Si se provee un `seed`, los bloques generados son reproducibles.

**Criterios de aceptación:**
- El mismo `seed` + mismos parámetros produce siempre los mismos bloques.
- Sin `seed`, el resultado es aleatorio en cada llamada.

---

### RF-03-003 — Filtrado por class y species
La narrativa recibe `class` y `species` como contexto y filtra templates por compatibilidad.

**Criterios de aceptación:**
- Solo se seleccionan templates compatibles con la `class` y `species` resueltas.
- Templates con tag `any` son válidos para cualquier `class` y `species`.
- Nunca se selecciona un template incompatible con la `class` o `species` del personaje.
- Un Barbarian nunca recibe trasfondos de erudito.
- Un Wizard nunca recibe trasfondos de combatiente salvaje.

---

### RF-03-004 — Variedad de templates
La seed data de narrativa debe garantizar variedad suficiente por categoría.

**Criterios de aceptación:**
- Mínimo 10 templates por categoría por clase en la seed.
- Mínimo 10 templates con tag `any` por categoría.
- La probabilidad de obtener el mismo bloque en dos generaciones consecutivas
  sin seed es estadísticamente baja.

---

### RF-03-005 — Regeneración independiente de ganchos
Cada NarrativeBlock puede regenerarse de forma independiente sin afectar otros campos.

**Criterios de aceptación:**
- El sistema expone operaciones de regeneración individuales para
  background, motivation y secret.
- Los campos no regenerados permanecen inalterados.
- La regeneración respeta el contexto de `class` y `species` del personaje.
- Si el bloque estaba locked, la operación es rechazada.

---

## Módulo 4 — Creador de Personaje (orquestador)

### RF-04-001 — Todo parámetro es opcional
El sistema genera un personaje completo sin ningún input obligatorio.

**Criterios de aceptación:**
- Una llamada sin parámetros produce un `Character` completo y válido.
- Todos los campos de `Character` están poblados al finalizar.

---

### RF-04-002 — Parámetros opcionales como locks implícitos
Un parámetro provisto en el request actúa como lock implícito para ese campo.

**Criterios de aceptación:**
- `class` provista → el personaje tiene esa clase, sin excepción.
- `species` provista → el personaje tiene esa species, sin excepción.
- `subSpecies` provista → el personaje tiene esa subSpecies, sin excepción.
- `gender` provisto → el nombre es coherente con ese género, sin excepción.
- `seed` provisto → el resultado completo es reproducible.

---

### RF-04-003 — Ejecución del pipeline en orden estricto
El orquestador ejecuta los 9 pasos del pipeline en el orden definido en `CONTEXT.md`.

**Criterios de aceptación:**
- Paso 1: inputs resueltos antes de cualquier generación.
- Pasos 2–6: stats resueltos antes de calcular derived.
- Paso 7: armadura resuelta antes de calcular AC.
- Paso 8: derived calculado sobre `finalStats` y armadura resuelta.
- Paso 9: narrativa filtrada por `class` y `species` resueltos en paso 1.
- Ningún paso usa resultados de un paso posterior.

---

### RF-04-004 — Coherencia entre módulos
El personaje generado es coherente entre nombre, stats y narrativa.

**Criterios de aceptación:**
- El nombre es coherente con `species`, `subSpecies` y `gender`.
- Los stats son coherentes con `class` y el background ASI pool resuelto.
- La narrativa es compatible con `class` y `species`.
- El background seleccionado es coherente con `class` y `species` (filtrado por Tags).
- No existe combinación de inputs válidos que produzca un personaje incoherente.

---

### RF-04-005 — Regeneración parcial con locks
El orquestador regenera campos específicos respetando los locks activos.

**Criterios de aceptación:**
- Los campos en `CharacterLocks` con valor `true` no se regeneran.
- Los campos con valor `false` se regeneran ejecutando los pasos del pipeline
  correspondientes.
- Los campos regenerados son coherentes con los campos locked.
- El `seed` original se descarta en regeneraciones sin seed nuevo.

---

## Módulo 5 — Gestión de personajes

### RF-05-001 — Guardar personaje
El usuario autenticado puede guardar un personaje generado en su biblioteca.

**Criterios de aceptación:**
- El personaje se asocia al usuario autenticado.
- El personaje guardado incluye todos los campos de la entidad `Character`.
- El sistema confirma el guardado con el `id` del personaje creado.

---

### RF-05-002 — Listar personajes guardados
El usuario puede ver la lista de personajes guardados en su biblioteca.

**Criterios de aceptación:**
- El sistema retorna solo los personajes del usuario autenticado.
- Cada ítem muestra al menos: nombre, species, class y nivel.
- La lista está ordenada por `createdAt` descendente.

---

### RF-05-003 — Abrir personaje guardado
El usuario puede abrir un personaje guardado y ver su ficha completa.

**Criterios de aceptación:**
- El sistema retorna todos los campos de la entidad `Character`.
- Solo el propietario del personaje puede acceder a él.
- Si el personaje no existe o no pertenece al usuario, el sistema retorna error claro.

---

### RF-05-004 — Edición básica de campos
El usuario puede editar campos de texto del personaje guardado.

**Criterios de aceptación:**
- El usuario puede modificar: nombre y `content` de cualquier NarrativeBlock.
- Los cambios se persisten en base de datos.
- Los campos mecánicos (stats, HP, AC) no son editables manualmente en el MVP.

---

### RF-05-005 — Eliminar personaje
El usuario puede eliminar un personaje de su biblioteca.

**Criterios de aceptación:**
- Solo el propietario puede eliminar el personaje.
- La operación es irreversible.
- El sistema confirma la eliminación.

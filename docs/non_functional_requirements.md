# rpg_engine — Requerimientos No Funcionales (MVP)

---

## Rendimiento

### RNF-01-001 — Tiempo de respuesta del generador
El pipeline completo de generación debe sentirse instantáneo para el usuario.

**Criterios de aceptación:**
- La generación de un personaje completo (9 pasos) responde en menos de 500ms.
- La regeneración parcial de un campo responde en menos de 200ms.
- Estos tiempos aplican bajo carga normal de un MVP (no escenario de alta concurrencia).

---

### RNF-01-002 — Tiempo de respuesta de operaciones de datos
Las operaciones de lectura y escritura en base de datos deben ser imperceptibles.

**Criterios de aceptación:**
- Guardar un personaje responde en menos de 300ms.
- Listar personajes guardados responde en menos de 300ms.
- Abrir un personaje guardado responde en menos de 200ms.

---

## Seguridad

### RNF-02-001 — Protección de endpoints autenticados
Los endpoints que operan sobre datos de usuario requieren sesión válida.

**Criterios de aceptación:**
- Toda operación del Módulo 5 requiere sesión autenticada.
- Una request sin sesión válida recibe un error 401.
- Una request con sesión válida pero accediendo a datos de otro usuario recibe 403.

---

### RNF-02-002 — Manejo seguro de tokens de magic link
Los tokens de acceso deben ser seguros, únicos y de vida limitada.

**Criterios de aceptación:**
- Los tokens son generados con entropía suficiente (mínimo 32 bytes aleatorios).
- Los tokens expiran en un tiempo definido (máximo 15 minutos).
- Un token usado no puede reutilizarse — se invalida inmediatamente tras el primer uso.
- Los tokens se almacenan hasheados en base de datos, no en texto plano.

---

### RNF-02-003 — Aislamiento de datos entre usuarios
Un usuario no puede acceder, modificar ni eliminar datos de otro usuario.

**Criterios de aceptación:**
- Toda query a personajes incluye el `user_id` del usuario autenticado como filtro.
- No existe endpoint que devuelva personajes de todos los usuarios.
- Las operaciones de edición y eliminación validan ownership antes de ejecutar.

---

## Disponibilidad y confiabilidad

### RNF-03-001 — Comportamiento ante fallos del generador
Un fallo en cualquier paso del pipeline no debe corromper datos existentes.

**Criterios de aceptación:**
- Si el pipeline falla en cualquier paso, no se persiste ningún dato parcial.
- El sistema retorna un error claro indicando que la generación falló.
- El usuario puede reintentar la operación sin efectos secundarios.

---

### RNF-03-002 — Integridad de datos guardados
Un personaje guardado debe ser siempre recuperable en el estado exacto en que fue guardado.

**Criterios de aceptación:**
- Los datos guardados son idénticos a los datos leídos para el mismo personaje.
- No existe operación que modifique parcialmente un personaje sin completarse.
- Las operaciones de escritura son atómicas.

---

## Escalabilidad

### RNF-04-001 — Extensibilidad del motor de generación
El motor debe soportar agregar nuevas clases, species o backgrounds sin cambios estructurales.

**Criterios de aceptación:**
- Agregar una nueva clase requiere solo agregar datos (baseline, hit die, armadura)
  sin modificar la lógica del pipeline.
- Agregar una nueva species requiere solo agregar datos (nombres, rasgos raciales)
  sin modificar la lógica del generador — las species no aportan ASIs en 5.5e.
- Agregar un nuevo background requiere solo agregar datos (asiPool, originFeat, tags)
  sin modificar el pipeline.

---

### RNF-04-002 — Escalabilidad de nivel sin cambio de firmas
Las funciones de cálculo deben soportar nivel 2–20 sin reescribirse.

**Criterios de aceptación:**
- `calculateHP(level, class, conMod)`, `calculateProficiency(level)` y
  `calculateAC(level, armor, modifiers)` reciben `level` como parámetro
  desde el MVP.
- Activar niveles 2–20 no requiere cambiar las firmas de estas funciones.

---

## Mantenibilidad

### RNF-05-001 — Cobertura de tests mínima en lógica de dominio
La lógica de dominio y los usecases deben tener cobertura de tests suficiente
para detectar regresiones.

**Criterios de aceptación:**
- Cobertura mínima de 80% en lógica de dominio (pipeline, fórmulas, resolución de bonos).
- Cobertura mínima de 80% en usecases (orquestador, regeneración parcial).
- Todo nuevo RF implementado incluye sus tests antes del merge.

---

### RNF-05-002 — Separación de capas (clean architecture)
El código debe respetar la separación entre dominio, usecases e infraestructura.

**Criterios de aceptación:**
- El dominio no importa nada de infraestructura.
- Los usecases dependen de interfaces, no de implementaciones concretas.
- La infraestructura (SQLite, email) es reemplazable sin tocar dominio ni usecases.
- Compile-time interface check en Go para repositorios (`var _ InterfaceName = (*Impl)(nil)`).

---

### RNF-05-003 — Early returns en toda la base de código
El código no usa bloques else anidados cuando un return es suficiente.

**Criterios de aceptación:**
- Las funciones validan precondiciones al inicio y retornan temprano en caso de error.
- El happy path siempre está al final de la función, sin anidamiento.
- El code review rechaza PRs que introduzcan else anidados evitables.

---

## Usabilidad

### RNF-06-001 — Feedback inmediato ante errores
El sistema comunica errores de forma clara y accionable.

**Criterios de aceptación:**
- Los errores de validación indican exactamente qué campo falló y por qué.
- Los errores de auth indican si el token es inválido o si expiró.
- Los errores del generador no exponen detalles internos al cliente.
- Ningún error retorna un stack trace al usuario final.

---

### RNF-06-002 — Generación perceptiblemente instantánea
El usuario no debe percibir latencia en la generación.

**Criterios de aceptación:**
- Si la generación supera 500ms, el cliente muestra un indicador de carga.
- El sistema nunca deja al usuario sin feedback por más de 1 segundo.

---

## Observabilidad

### RNF-07-001 — Logs de generación
El sistema registra cada generación con suficiente información para reproducirla.

**Criterios de aceptación:**
- Cada generación registra: timestamp, parámetros de entrada, seed usado y duración.
- Los logs no contienen datos sensibles del usuario (email, tokens).
- Los logs son estructurados (JSON) para facilitar búsqueda y filtrado.

---

### RNF-07-002 — Trazabilidad de errores
Los errores del sistema son trazables desde el log hasta el punto exacto de fallo.

**Criterios de aceptación:**
- Cada error registra: timestamp, módulo, operación y mensaje descriptivo.
- Los errores de infraestructura (DB, email) se distinguen de los errores de dominio.
- Un error en producción puede reproducirse localmente con la información del log.

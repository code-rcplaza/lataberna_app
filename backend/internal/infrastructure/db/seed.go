package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// narrativeSeedEntry represents one narrative entry with its compatibility rows.
type narrativeSeedEntry struct {
	category string
	content  string
	compat   []compatRow // empty = universal (default weight 2)
}

type compatRow struct {
	dimension string // "class" or "species"
	value     string
	group     string // "primary", "secondary", "excluded"
}

// SeedContentIfEmpty checks if narrative_entries and name_entries are populated.
// If either table is empty, it inserts the full editorial seed dataset in a single
// transaction. Safe to call on every startup — idempotent.
func SeedContentIfEmpty(ctx context.Context, db *sql.DB) error {
	if err := seedNarrativeIfEmpty(ctx, db); err != nil {
		return fmt.Errorf("SeedContentIfEmpty: narrative: %w", err)
	}
	if err := seedNamesIfEmpty(ctx, db); err != nil {
		return fmt.Errorf("SeedContentIfEmpty: names: %w", err)
	}
	return nil
}

func seedNarrativeIfEmpty(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM narrative_entries`).Scan(&count); err != nil {
		return fmt.Errorf("seedNarrativeIfEmpty: count: %w", err)
	}
	if count > 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("seedNarrativeIfEmpty: begin: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	entryStmt, err := tx.PrepareContext(ctx,
		`INSERT INTO narrative_entries (id, category, content, created_at) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("seedNarrativeIfEmpty: prepare entry: %w", err)
	}
	defer entryStmt.Close()

	compatStmt, err := tx.PrepareContext(ctx,
		`INSERT OR IGNORE INTO narrative_compatibility (entry_id, dimension, value, group_name) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("seedNarrativeIfEmpty: prepare compat: %w", err)
	}
	defer compatStmt.Close()

	now := time.Now().UTC().Format(time.RFC3339Nano)

	for i, entry := range narrativeSeedData() {
		id := fmt.Sprintf("narr-%04d", i+1)
		if _, err := entryStmt.ExecContext(ctx, id, entry.category, entry.content, now); err != nil {
			return fmt.Errorf("seedNarrativeIfEmpty: insert entry %d: %w", i, err)
		}
		for _, c := range entry.compat {
			if _, err := compatStmt.ExecContext(ctx, id, c.dimension, c.value, c.group); err != nil {
				return fmt.Errorf("seedNarrativeIfEmpty: insert compat %d: %w", i, err)
			}
		}
	}

	return tx.Commit()
}

func seedNamesIfEmpty(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM name_entries`).Scan(&count); err != nil {
		return fmt.Errorf("seedNamesIfEmpty: count: %w", err)
	}
	if count > 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("seedNamesIfEmpty: begin: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	stmt, err := tx.PrepareContext(ctx,
		`INSERT OR IGNORE INTO name_entries (id, species_key, gender, name, created_at) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("seedNamesIfEmpty: prepare: %w", err)
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339Nano)
	idx := 0
	for speciesKey, genderMap := range nameSeedData() {
		for gender, names := range genderMap {
			for _, name := range names {
				idx++
				id := fmt.Sprintf("name-%05d", idx)
				if _, err := stmt.ExecContext(ctx, id, speciesKey, gender, name, now); err != nil {
					return fmt.Errorf("seedNamesIfEmpty: insert name %q %q %q: %w", speciesKey, gender, name, err)
				}
			}
		}
	}

	return tx.Commit()
}

// ---------------------------------------------------------------------------
// Seed content — narrative entries
// ---------------------------------------------------------------------------

// Class archetype helpers for compatibility rows.
// Warriors: fighter, barbarian, paladin, ranger
// Scholars: wizard, sorcerer, artificer, bard
// Faithful: cleric, druid, monk, paladin
// Shadows:  rogue, warlock, ranger
// Wanderers: ranger, bard, druid, rogue

func warriors(group string) []compatRow {
	classes := []string{"fighter", "barbarian", "paladin", "ranger"}
	out := make([]compatRow, len(classes))
	for i, c := range classes {
		out[i] = compatRow{"class", c, group}
	}
	return out
}

func scholars(group string) []compatRow {
	classes := []string{"wizard", "sorcerer", "artificer", "bard"}
	out := make([]compatRow, len(classes))
	for i, c := range classes {
		out[i] = compatRow{"class", c, group}
	}
	return out
}

func faithful(group string) []compatRow {
	classes := []string{"cleric", "druid", "monk", "paladin"}
	out := make([]compatRow, len(classes))
	for i, c := range classes {
		out[i] = compatRow{"class", c, group}
	}
	return out
}

func shadows(group string) []compatRow {
	classes := []string{"rogue", "warlock", "ranger"}
	out := make([]compatRow, len(classes))
	for i, c := range classes {
		out[i] = compatRow{"class", c, group}
	}
	return out
}

func wanderers(group string) []compatRow {
	classes := []string{"ranger", "bard", "druid", "rogue"}
	out := make([]compatRow, len(classes))
	for i, c := range classes {
		out[i] = compatRow{"class", c, group}
	}
	return out
}

func merge(slices ...[]compatRow) []compatRow {
	var out []compatRow
	for _, s := range slices {
		out = append(out, s...)
	}
	return out
}

func classRows(classes []string, group string) []compatRow {
	out := make([]compatRow, len(classes))
	for i, c := range classes {
		out[i] = compatRow{"class", c, group}
	}
	return out
}

func speciesRows(species []string, group string) []compatRow {
	out := make([]compatRow, len(species))
	for i, s := range species {
		out[i] = compatRow{"species", s, group}
	}
	return out
}

func narrativeSeedData() []narrativeSeedEntry {
	return []narrativeSeedEntry{

		// ═══════════════════════════════════════════════════════════════════
		// BACKGROUND — universal (no compat rows = weight 2 for all)
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "background",
			content:  "Criado en las calles de una ciudad portuaria, aprendiste desde joven que la supervivencia depende de la astucia, no de la fuerza bruta.",
		},
		{
			category: "background",
			content:  "Perdiste a tu familia durante un invierno brutal. Desde entonces vagás de pueblo en pueblo buscando un lugar al que llamar hogar, aunque aún no lo encontraste.",
		},
		{
			category: "background",
			content:  "Fuiste aprendiz de un maestro artesano durante años. Cuando él murió sin dejar herencia, tomaste tus herramientas y te lanzaste al mundo a labrar tu propio destino.",
		},
		{
			category: "background",
			content:  "Tu aldea te expulsó por romper una regla ancestral que considerabas injusta. El exilio te enseñó más sobre el mundo que cualquier educación formal.",
		},
		{
			category: "background",
			content:  "Creciste escuchando historias de héroes y leyendas. Cuando tuviste edad suficiente, decidiste que era hora de protagonizar tu propia historia.",
		},
		{
			category: "background",
			content:  "Trabajaste como mercenario durante años, vendiendo tu espada al mejor postor. Un encargo salió terriblemente mal y te dejó con una deuda que aún intentás saldar.",
		},
		{
			category: "background",
			content:  "Sobreviviste a un naufragio que mató a toda tu tripulación. Los meses que pasaste varado en una isla deshabitada te moldearon en cuerpo y espíritu.",
		},
		{
			category: "background",
			content:  "Fuiste el único hijo de una familia de mercaderes arruinados. Aprendiste a negociar antes de aprender a pelear, y eso dice mucho de quién sos.",
		},
		{
			category: "background",
			content:  "Creciste en un orfanato administrado por un culto menor. Cuando descubriste las verdaderas intenciones del lugar, huiste sin mirar atrás.",
		},
		{
			category: "background",
			content:  "Fuiste soldado raso en un ejército que perdió una guerra. El armisticio te dejó sin rumbo, sin paga y con demasiados recuerdos para dormir tranquilo.",
		},
		{
			category: "background",
			content:  "Creciste en la frontera entre dos reinos en conflicto. Aprendiste a desconfiar de ambos bandos y a sobrevivir con lo que el terreno disputado te dejaba.",
		},
		{
			category: "background",
			content:  "Tu familia sirvió durante generaciones a una casa noble que fue exterminada de la noche a la mañana. Quedaste libre y sin propósito al mismo tiempo.",
		},
		{
			category: "background",
			content:  "Creciste en una caravana que nunca se detuvo demasiado tiempo en ningún lugar. El mundo entero es tu barrio; ningún sitio es tu hogar.",
		},
		{
			category: "background",
			content:  "Pasaste tu infancia en una ciudad subterránea donde la luz del sol era un rumor. El día que saliste a la superficie cambió algo en vos para siempre.",
		},
		{
			category: "background",
			content:  "Fuiste el testigo accidental de un crimen que involucra a personas poderosas. Desde entonces, alguien te busca y vos buscás a alguien que te pueda proteger.",
		},
		{
			category: "background",
			content:  "Te criaste en una comunidad rural tan pequeña que todos se conocían. La primera vez que pisaste una ciudad sentiste que el mundo era demasiado grande para vos.",
		},
		{
			category: "background",
			content:  "Fuiste explorador cartográfico para una gremio de aventureros. Conocés caminos que no aparecen en ningún mapa y ruinas que nadie más ha documentado.",
		},
		{
			category: "background",
			content:  "Naciste durante un eclipse total. El pueblo donde creciste siempre te miró como un presagio andante; vos nunca supiste si eso era bueno o malo.",
		},
		{
			category: "background",
			content:  "Trabajaste en una posada de camino durante años. Escuchaste suficientes historias de aventureros para saber qué suelen hacer mal. Ahora te toca a vos.",
		},
		{
			category: "background",
			content:  "Eras marinero en una flota mercante hasta que tu barco encontró algo en altamar que ningún marinero debería ver. Los pocos que sobrevivieron nunca hablaron de eso.",
		},
		{
			category: "background",
			content:  "Creciste en una familia de curanderos de aldea. Aprendiste que la gente acude a vos en sus peores momentos y que eso es tanto un privilegio como una carga.",
		},
		{
			category: "background",
			content:  "Fuiste criado por abuelos que nunca te dijeron quiénes eran tus padres. Cuando morieron, lo único que encontraste fue una carta sellada con un sello que no reconociste.",
		},
		{
			category: "background",
			content:  "Creciste en una zona de conflicto donde las fronteras cambiaban cada año. Aprendiste a adaptarte antes de aprender a confiar.",
		},
		{
			category: "background",
			content:  "Trabajaste en una mina durante años, en la oscuridad y el polvo. La primera vez que viste el horizonte abierto, entendiste que el mundo era más grande de lo que imaginabas.",
		},
		{
			category: "background",
			content:  "Fuiste estudiante brillante que abandonó sus estudios para seguir algo que el sistema no podía enseñarle. Nadie entendió la decisión. Vos tampoco, del todo.",
		},
		{
			category: "background",
			content:  "Eras el cronista de una expedición científica que terminó en desastre. Sobreviviste. El diario también. Lo que escribiste en él no coincide exactamente con lo que ocurrió.",
		},
		{
			category: "background",
			content:  "Creciste en un barrio donde conocer al tipo equivocado era cuestión de tiempo. Conociste a muchos tipos equivocados. Aprendiste de todos ellos.",
		},
		{
			category: "background",
			content:  "Fuiste el favorito de un mentor que luego resultó estar manipulando a todos a su alrededor. Lo que te enseñó es real y útil. Eso lo hace más confuso, no menos.",
		},
		{
			category: "background",
			content:  "Heredaste una propiedad en ruinas y una lista de deudas que no sabías que existían. Vendiste lo que pudiste y tomaste lo que quedó: una dirección y una pista.",
		},
		{
			category: "background",
			content:  "Creciste en una comunidad que practicaba rituales que el mundo exterior consideraba herejía. No sabés si tenían razón o no, pero sí sabés que te formaron en algo real.",
		},
		{
			category: "background",
			content:  "Fuiste el asistente personal de una figura pública durante años. Conociste su vida entera desde adentro. Cuando dejaste ese trabajo, llevaste cosas que nadie debería saber.",
		},
		{
			category: "background",
			content:  "Sobreviviste a una epidemia que mató a la mayoría de tu comunidad. No sabés por qué vos no caíste. Esa pregunta te persigue más que el duelo.",
		},
		{
			category: "background",
			content:  "Eras el segundo de a bordo en una operación que salió mal. El líder escapó con el crédito. Vos te quedaste con las consecuencias. Aprendiste más de esa derrota que de cualquier éxito.",
		},
		{
			category: "background",
			content:  "Creciste escuchando al anciano más viejo del pueblo contar historias de un tiempo que nadie más recordaba. Cuando murió, te dejó una llave y ninguna explicación.",
		},
		{
			category: "background",
			content:  "Fuiste el encargado del cementerio de tu pueblo durante años. Conociste a todos los muertos antes que a los vivos, y eso te dejó con una perspectiva peculiar sobre ambos.",
		},
		{
			category: "background",
			content:  "Viajaste como mensajero durante años, llevando palabras de un lugar a otro sin entender siempre su peso. Un día entregaste un mensaje que desencadenó algo que aún no terminó.",
		},
		{
			category: "background",
			content:  "Creciste en una familia numerosa donde siempre fuiste el que resolvía los problemas de los demás. Cuando por fin te fuiste, no sabías qué hacer con tanto silencio.",
		},
		{
			category: "background",
			content:  "Eras aprendiz de un alquimista que experimentaba con cosas que la academia académica oficial rechazaría. Aprendiste más en ese taller que en cualquier libro.",
		},
		{
			category: "background",
			content:  "Fuiste el único de tu grupo que no perdió la cabeza durante una crisis. Desde entonces, cada vez que hay caos, todos te miran. No siempre tenés la respuesta.",
		},
		{
			category: "background",
			content:  "Creciste en una ciudad donde el arte era la única moneda que importaba. No tenés dinero, pero sabés reconocer lo que vale algo antes de que el mundo lo sepa.",
		},
		{
			category: "background",
			content:  "Trabajaste como diplomático junior en una misión que fracasó espectacularmente. La guerra que siguió duró poco, pero las cicatrices que dejó en vos no tienen plazo.",
		},
		{
			category: "background",
			content:  "Fuiste el escéptico del grupo que resultó tener razón sobre algo crucial en el peor momento posible. La satisfacción fue exactamente tan amarga como esperabas.",
		},
		{
			category: "background",
			content:  "Eras el asistente de un sacerdote que nunca creyó en nada de lo que predicaba. Eso te enseñó más sobre la fe —y sobre la duda— que cualquier sermón.",
		},
		{
			category: "background",
			content:  "Creciste en el territorio disputado entre dos señores feudales que nunca llegaron a un acuerdo. Aprendiste a leer las tensiones políticas antes de que exploten.",
		},

		// ═══════════════════════════════════════════════════════════════════
		// BACKGROUND — Warriors primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "background",
			content:  "Tu aldea fue arrasada por un dragón cuando eras niño. Desde ese día, el fuego de la venganza arde en tu pecho más fuerte que cualquier llama.",
			compat:   merge(warriors("primary"), scholars("secondary")),
		},
		{
			category: "background",
			content:  "Creciste en las estepas heladas del norte entre guerreros que se curtían en las tormentas. La civilización te parece blanda y ruidosa.",
			compat:   merge(classRows([]string{"barbarian"}, "primary"), classRows([]string{"ranger"}, "secondary"), scholars("excluded")),
		},
		{
			category: "background",
			content:  "Tu tribu te eligió como campeón tras derrotar al jefe anterior en combate singular. El título pesó más de lo esperado y terminaste huyendo de las responsabilidades.",
			compat:   merge(classRows([]string{"barbarian", "fighter"}, "primary"), faithful("secondary")),
		},
		{
			category: "background",
			content:  "Entrenaste en una academia militar desde los diez años. La disciplina es tu segunda naturaleza; la improvisación te genera una incomodidad que aprendiste a disimular.",
			compat:   merge(classRows([]string{"fighter", "paladin"}, "primary"), wanderers("excluded")),
		},
		{
			category: "background",
			content:  "Eras el guardaespaldas personal de un noble menor que fue asesinado mientras dormía. No pudiste protegerlo y eso te costó más que el trabajo.",
			compat:   merge(warriors("primary"), shadows("secondary")),
		},
		{
			category: "background",
			content:  "Participaste en un torneo de campeones que te ganó fama regional. Después del torneo, los desafíos empezaron a llegar solos y tuviste que decidir cuáles valía la pena aceptar.",
			compat:   merge(classRows([]string{"fighter", "barbarian"}, "primary"), classRows([]string{"monk"}, "secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// BACKGROUND — Scholars primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "background",
			content:  "Pasaste tu infancia devorando libros en la biblioteca del templo local. Cuando los libros ya no fueron suficientes, buscaste conocimiento donde nadie más se atrevía a mirar.",
			compat:   merge(scholars("primary"), faithful("secondary"), warriors("excluded")),
		},
		{
			category: "background",
			content:  "Eras el prodigio de tu academia mágica hasta que un experimento fallido destruyó parte del ala este. Te exiliaron, pero las ecuaciones todavía danzan en tu mente.",
			compat:   merge(classRows([]string{"wizard", "artificer"}, "primary"), scholars("secondary")),
		},
		{
			category: "background",
			content:  "Tu mentor fue un mago excéntrico que vivía en una torre al borde del abismo. Sus enseñanzas eran brillantes, sus métodos cuestionables, y su muerte aún te pesa.",
			compat:   merge(classRows([]string{"wizard", "sorcerer"}, "primary"), scholars("secondary")),
		},
		{
			category: "background",
			content:  "Trabajaste como escriba en una corte real durante años. Copiaste tratados, traiciones y correspondencia íntima. Sabés más secretos que cualquier espía.",
			compat:   merge(scholars("primary"), shadows("secondary"), warriors("excluded")),
		},
		{
			category: "background",
			content:  "Fuiste astrónomo del observatorio de una ciudad imperial. Una noche, algo en el cielo te miró de vuelta. Renunciaste al día siguiente.",
			compat:   merge(classRows([]string{"wizard", "sorcerer", "bard"}, "primary"), faithful("secondary")),
		},
		{
			category: "background",
			content:  "Construiste tu primer autómata a los catorce años usando piezas de relojes robados. Desde entonces, nadie te preguntó cómo lo hiciste; todos solo querían que siguieras.",
			compat:   merge(classRows([]string{"artificer"}, "primary"), scholars("secondary"), faithful("excluded")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// BACKGROUND — Faithful primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "background",
			content:  "Desde niño serviste en el templo como acólito. Una noche, la deidad habló directamente a tu corazón y desde entonces sabés que fuiste elegido para algo mayor.",
			compat:   merge(classRows([]string{"cleric", "paladin"}, "primary"), faithful("secondary"), shadows("excluded")),
		},
		{
			category: "background",
			content:  "Pasaste una década en un monasterio de clausura. Cuando saliste, el mundo afuera te resultó ruidoso, frenético y sorprendentemente hermoso.",
			compat:   merge(classRows([]string{"monk", "cleric"}, "primary"), faithful("secondary")),
		},
		{
			category: "background",
			content:  "Tu fe salvó a alguien que todos los médicos habían abandonado. Eso te convenció de que el poder divino es real y de que vos sos su canal. Nadie te dijo que el canal también puede romperse.",
			compat:   merge(faithful("primary"), warriors("secondary")),
		},
		{
			category: "background",
			content:  "Naciste en una familia de druidas que guardaban un bosque antiguo. Aprendiste antes a hablar con los árboles que con las personas; a veces preferís eso.",
			compat:   merge(classRows([]string{"druid", "ranger"}, "primary"), wanderers("secondary"), scholars("secondary")),
		},
		{
			category: "background",
			content:  "Fuiste misionero en tierras hostiles. Convertiste a nadie, sobreviviste a todo, y volviste con una visión del mundo que no encaja en ninguna doctrina.",
			compat:   merge(faithful("primary"), wanderers("secondary")),
		},
		{
			category: "background",
			content:  "Tu orden religiosa fue proscripta por la corona. Dispersados y perseguidos, los pocos miembros que quedan confían en vos para mantener viva la llama de su fe.",
			compat:   merge(faithful("primary"), shadows("secondary"), warriors("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// BACKGROUND — Shadows primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "background",
			content:  "Fuiste el mejor carterista del gremio hasta que robaste la bolsa equivocada. El noble al que le sisaste era más peligroso de lo que parecía.",
			compat:   merge(shadows("primary"), wanderers("secondary"), faithful("excluded")),
		},
		{
			category: "background",
			content:  "Actuaste en tabernas y mercados desde los siete años. La escena te dio carisma; las calles te dieron sentido común.",
			compat:   merge(classRows([]string{"bard", "rogue"}, "primary"), wanderers("secondary")),
		},
		{
			category: "background",
			content:  "Trabajaste como informante para tres facciones distintas al mismo tiempo. Cuando dos de ellas se aliaron, tuviste que desaparecer antes de que sumaran dos más dos.",
			compat:   merge(classRows([]string{"rogue", "warlock"}, "primary"), shadows("secondary"), faithful("excluded")),
		},
		{
			category: "background",
			content:  "Pasaste tres años en prisión por un crimen que no cometiste. Saliste con habilidades que nadie que entra inocente debería tener.",
			compat:   merge(shadows("primary"), warriors("secondary")),
		},
		{
			category: "background",
			content:  "En el peor momento de tu vida, una entidad de otro plano te ofreció poder a cambio de un servicio aún no definido. Aceptaste sin pensar demasiado.",
			compat:   merge(classRows([]string{"warlock"}, "primary"), shadows("secondary"), faithful("excluded")),
		},
		{
			category: "background",
			content:  "Fuiste el mejor asesino de un gremio de sombras hasta que tu último objetivo resultó ser un inocente. Abandonaste el gremio esa noche; el gremio no te olvidó.",
			compat:   merge(classRows([]string{"rogue"}, "primary"), shadows("secondary"), faithful("excluded")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// BACKGROUND — Wanderers primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "background",
			content:  "Creciste en el bosque profundo, criado por un anciano druida que te enseñó que cada árbol tiene una historia y cada bestia merece respeto.",
			compat:   merge(wanderers("primary"), faithful("secondary")),
		},
		{
			category: "background",
			content:  "Fuiste el único sobreviviente de una expedición diezmada por criaturas del bosque. Desde entonces, aprendiste su lenguaje para no volver a ser sorprendido.",
			compat:   merge(classRows([]string{"ranger"}, "primary"), wanderers("secondary"), scholars("secondary")),
		},
		{
			category: "background",
			content:  "Actuaste como guía de caravanas durante años. Conocés cada camino, cada taberna y cada emboscada habitual entre tres reinos. Cobrás bien por eso.",
			compat:   merge(wanderers("primary"), warriors("secondary")),
		},
		{
			category: "background",
			content:  "Creciste entre bardos itinerantes que tocaban en cada aldea del continente. La música fue tu idioma antes de que aprendieras a hablar bien.",
			compat:   merge(classRows([]string{"bard"}, "primary"), wanderers("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// BACKGROUND — Species specific
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "background",
			content:  "Viviste tus primeros doscientos años en la ciudad élfil, donde el tiempo pasa lento y las decisiones pesan siglos. Un día decidiste que ya era hora de ver el mundo mortal.",
			compat:   speciesRows([]string{"elf"}, "primary"),
		},
		{
			category: "background",
			content:  "En las montañas de tu clan enano, el honor es lo más sagrado. Perdiste el tuyo por razones que solo vos conocés, y desde entonces buscás redimirte con cada acción.",
			compat:   speciesRows([]string{"dwarf"}, "primary"),
		},
		{
			category: "background",
			content:  "Como medio elfo, nunca perteneciste del todo a ningún mundo. Eso te hizo observador, adaptable y, en el fondo, profundamente solitario.",
			compat:   speciesRows([]string{"half-elf"}, "primary"),
		},
		{
			category: "background",
			content:  "Tu sangre infernal siempre te hizo diferente. En tu ciudad natal te miraban con desconfianza; en el camino, aprendiste a usar esa incomodidad a tu favor.",
			compat:   speciesRows([]string{"tiefling"}, "primary"),
		},
		{
			category: "background",
			content:  "El clan Dragonborn al que pertenecés lleva siglos de honor. Vos sos el primero en generaciones que eligió aventurarse solo, rompiendo una tradición colectiva.",
			compat:   speciesRows([]string{"dragonborn"}, "primary"),
		},
		{
			category: "background",
			content:  "Creciste en los alrededores de una ciudad enana como humano en minoría. Aprendiste a moverte entre gente que te supera en fuerza y paciencia; eso te volvió diplomático por necesidad.",
			compat:   merge(speciesRows([]string{"human"}, "primary"), speciesRows([]string{"half-elf"}, "secondary")),
		},
		{
			category: "background",
			content:  "Tu aldea halfling era conocida en toda la región por su gastronomía. Vos eras el único que quería ver qué había más allá de la colina del fondo.",
			compat:   speciesRows([]string{"halfling"}, "primary"),
		},
		{
			category: "background",
			content:  "Creciste en un taller gnomo donde todo se desarmaba y se volvía a armar mejor. Tu curiosidad nunca encontró límites, ni siquiera los que la física debería imponer.",
			compat:   speciesRows([]string{"gnome"}, "primary"),
		},
		{
			category: "background",
			content:  "Fuiste criado entre humanos que temían tu herencia orca. Aprendiste que la fuerza que ellos temían era exactamente lo que necesitarías para sobrevivir sin ellos.",
			compat:   speciesRows([]string{"half-orc"}, "primary"),
		},

		// ═══════════════════════════════════════════════════════════════════
		// MOTIVATION — universal
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "motivation",
			content:  "Buscás riqueza suficiente para comprar la libertad de alguien que amás y que aún está atrapado en circunstancias que no pudo elegir.",
		},
		{
			category: "motivation",
			content:  "Querés demostrarles a todos los que dudaron de vos que estaban equivocados. La ambición surgió del dolor, pero el fuego ya es tuyo.",
		},
		{
			category: "motivation",
			content:  "Buscás un artefacto perdido que, según la leyenda, puede devolver la vida. No importa cuánto cueste encontrarlo.",
		},
		{
			category: "motivation",
			content:  "Tu única motivación real es sobrevivir un día más. Todo lo demás —el honor, la gloria, la causa— son palabras bonitas que otros usan para convencerte de arriesgar el cuero.",
		},
		{
			category: "motivation",
			content:  "Creés que hay una amenaza oscura que nadie más puede ver. Cada paso del camino es un intento de reunir pruebas antes de que sea demasiado tarde.",
		},
		{
			category: "motivation",
			content:  "Querés construir algo que perdure: una fortaleza, una institución, un nombre que se recuerde siglos después de que hayas muerto.",
		},
		{
			category: "motivation",
			content:  "La curiosidad es tu maldición y tu motor. Necesitás saber qué hay del otro lado de cada puerta, de cada horizonte, de cada pregunta sin respuesta.",
		},
		{
			category: "motivation",
			content:  "Alguien a quien querías murió sin justicia. Nadie más va a hacer nada al respecto. Así que lo vas a hacer vos.",
		},
		{
			category: "motivation",
			content:  "Escapaste de una situación de la que nadie escapa. Necesitás entender cómo fue posible, porque si vos pudiste, otros también pueden.",
		},
		{
			category: "motivation",
			content:  "Viviste toda tu vida siguiendo las reglas de otros. Un día decidiste que ya era suficiente y saliste a escribir las tuyas.",
		},
		{
			category: "motivation",
			content:  "Una deuda de honor te ata a una causa que ya no elegirías. Pero rompiste tu palabra una vez y no lo vas a hacer de nuevo.",
		},
		{
			category: "motivation",
			content:  "Buscás un lugar en el mundo donde pertenecer de verdad. No un lugar geográfico, sino un grupo de personas que te entiendan sin que tengas que explicarte.",
		},
		{
			category: "motivation",
			content:  "Querés ver el mundo antes de morir. Cada mapa que completás, cada ciudad que pisás, es un paso más hacia esa cuenta regresiva que todos ignoramos.",
		},
		{
			category: "motivation",
			content:  "Hay una verdad que el mundo poderoso suprime activamente. Vos la descubriste y ahora no podés fingir que no la sabés.",
		},
		{
			category: "motivation",
			content:  "Alguien te salvó la vida cuando no tenías nada para ofrecer a cambio. Hasta que puedas devolver eso, no podés parar.",
		},
		{
			category: "motivation",
			content:  "Creés que el sistema actual es irrescatable y que hay que construir uno nuevo sobre las cenizas. Cada acción que tomás es un ladrillo de ese futuro.",
		},
		{
			category: "motivation",
			content:  "Creciste viendo a personas con talento desperdiciar sus vidas por falta de oportunidades. Querés ser la oportunidad que vos nunca tuviste.",
		},
		{
			category: "motivation",
			content:  "Tenés una deuda simbólica con el mundo: sobreviviste cuando otros no lo hicieron, y eso te obliga a hacer algo que valga la pena con el tiempo extra.",
		},
		{
			category: "motivation",
			content:  "Querés entender por qué el mundo es como es. No para cambiarlo necesariamente, sino porque vivir sin esa comprensión te resulta insoportable.",
		},
		{
			category: "motivation",
			content:  "Te prometiste que nunca más ibas a depender de nadie para tu seguridad. Cada habilidad nueva que desarrollás es un paso más hacia esa independencia.",
		},
		{
			category: "motivation",
			content:  "Hay algo que hacés mejor que nadie en tu pueblo y peor que cualquier aventurero experimentado. Saliste para encontrar el punto medio.",
		},
		{
			category: "motivation",
			content:  "Una persona que admirabas resultó ser una decepción monumental. Desde entonces, buscás algo o alguien que merezca de verdad la admiración que tenés para dar.",
		},
		{
			category: "motivation",
			content:  "El miedo te tuvo paralizado durante demasiado tiempo. Un día decidiste que la única cura era hacer exactamente lo que más temías. Todavía lo estás haciendo.",
		},
		{
			category: "motivation",
			content:  "Querés ver con tus propios ojos si las cosas que te enseñaron de niño son verdad o mentira. Hasta ahora el resultado es mixto.",
		},
		{
			category: "motivation",
			content:  "Alguien apostó a que no podías. Eso fue hace tres años. Todavía no terminaste de demostrar que se equivocaban.",
		},

		// ═══════════════════════════════════════════════════════════════════
		// MOTIVATION — Warriors primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "motivation",
			content:  "Juraste ante los dioses que vengarías la destrucción de tu hogar. Cada enemigo caído es un paso más hacia ese juramento cumplido.",
			compat:   merge(classRows([]string{"barbarian", "paladin", "fighter"}, "primary"), faithful("secondary")),
		},
		{
			category: "motivation",
			content:  "Querés perfeccionar el arte del combate hasta alcanzar un estado que ningunos escritos describen pero que intuís que existe.",
			compat:   merge(classRows([]string{"fighter", "monk", "barbarian"}, "primary"), scholars("secondary")),
		},
		{
			category: "motivation",
			content:  "Proteger a los débiles no es una filosofía para vos; es una compulsión. Cuando alguien necesita ayuda y no la recibe, algo en tu pecho se tensa hasta que actuás.",
			compat:   merge(classRows([]string{"paladin", "fighter", "ranger"}, "primary"), faithful("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// MOTIVATION — Scholars primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "motivation",
			content:  "Cada invención que creás es un paso hacia un prototipo imposible que vive en tus planos desde hace años. Todo lo demás es financiamiento.",
			compat:   merge(classRows([]string{"artificer"}, "primary"), scholars("secondary"), faithful("excluded")),
		},
		{
			category: "motivation",
			content:  "Buscás el origen de tu sangre mágica. Alguien en tu linaje hizo algo extraordinario —o terrible— y necesitás saber qué fue.",
			compat:   merge(classRows([]string{"sorcerer"}, "primary"), scholars("secondary")),
		},
		{
			category: "motivation",
			content:  "Querés comprender los mecanismos fundamentales de la magia. No para usarla mejor, sino para entender por qué existe en absoluto.",
			compat:   merge(classRows([]string{"wizard"}, "primary"), scholars("secondary")),
		},
		{
			category: "motivation",
			content:  "Querés componer la epopeya definitiva: una obra que capture la verdad del mundo tal como solo vos la podés ver. Las aventuras son tu investigación de campo.",
			compat:   merge(classRows([]string{"bard"}, "primary"), wanderers("secondary"), scholars("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// MOTIVATION — Faithful primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "motivation",
			content:  "Tu deidad te encomendó una misión sagrada. No vas a dormir tranquilo hasta cumplirla, aunque eso cueste todo lo que tenés.",
			compat:   merge(classRows([]string{"cleric", "paladin"}, "primary"), faithful("secondary")),
		},
		{
			category: "motivation",
			content:  "Buscás el equilibrio perdido entre la magia arcana y el mundo natural. Creés que sin ese equilibrio, todo lo que existe está condenado.",
			compat:   merge(classRows([]string{"druid", "ranger"}, "primary"), wanderers("secondary"), scholars("excluded")),
		},
		{
			category: "motivation",
			content:  "La meditación te reveló que sos el último guardián de un linaje de guardianes del equilibrio. No elegiste el rol, pero no podés ignorarlo.",
			compat:   merge(classRows([]string{"monk"}, "primary"), faithful("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// MOTIVATION — Shadows primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "motivation",
			content:  "La deuda con tu patrón del pacto crece. Cada hazaña que realizás es parte del pago de algo que ya no recordás haber prometido.",
			compat:   merge(classRows([]string{"warlock"}, "primary"), shadows("secondary"), faithful("excluded")),
		},
		{
			category: "motivation",
			content:  "Un hechizo que lanzaste sin querer cambió el destino de alguien inocente. Llevás ese peso y buscás una forma de enmendar lo que rompiste.",
			compat:   merge(classRows([]string{"wizard", "sorcerer"}, "primary"), scholars("secondary")),
		},
		{
			category: "motivation",
			content:  "Te robaron algo que no tiene precio —no dinero, sino algo que define quién sos— y no vas a parar hasta recuperarlo.",
			compat:   merge(classRows([]string{"rogue", "ranger"}, "primary"), shadows("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// MOTIVATION — Species specific
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "motivation",
			content:  "Tu longevidad élfica te da perspectiva que los mortales no tienen. Eso te pesa: ves el patrón del error humano repetirse sin fin y querés romperlo.",
			compat:   speciesRows([]string{"elf"}, "primary"),
		},
		{
			category: "motivation",
			content:  "El clan enano exige resultados. Volvés cuando tengas un logro digno de ser tallado en la piedra de la sala ancestral.",
			compat:   speciesRows([]string{"dwarf"}, "primary"),
		},
		{
			category: "motivation",
			content:  "Como halfling, el mundo te subestima siempre. Eso te dio una ventaja que jamás vas a desperdiciar.",
			compat:   speciesRows([]string{"halfling"}, "primary"),
		},
		{
			category: "motivation",
			content:  "Tu curiosidad gnoma no tiene límites: necesitás entender cómo funciona todo, desarmarlo si es necesario, y armarlo de nuevo pero mejor.",
			compat:   speciesRows([]string{"gnome"}, "primary"),
		},
		{
			category: "motivation",
			content:  "Querés demostrar que la sangre infernal no define el destino. Cada bien que hacés es una refutación de lo que el mundo esperaba de vos.",
			compat:   speciesRows([]string{"tiefling"}, "primary"),
		},
		{
			category: "motivation",
			content:  "Tu herencia orca te dio fuerza; tu herencia humana te dio ambición. Juntas te hacen imparable, si lográs que el mundo te dé una oportunidad.",
			compat:   speciesRows([]string{"half-orc"}, "primary"),
		},
		{
			category: "motivation",
			content:  "Querés reconstruir el nombre de tu familia humana después de que las deudas y las malas decisiones lo destruyeran. Empezás desde cero con lo que tenés.",
			compat:   speciesRows([]string{"human"}, "primary"),
		},
		{
			category: "motivation",
			content:  "Como medio elfo, cargás con las expectativas de dos mundos que nunca se pusieron de acuerdo. Decidiste ignorar ambos y definirte solo.",
			compat:   speciesRows([]string{"half-elf"}, "primary"),
		},
		{
			category: "motivation",
			content:  "Tu clan dragonborn lleva generaciones sin ver a un verdadero héroe salir de sus filas. Eso termina con vos.",
			compat:   speciesRows([]string{"dragonborn"}, "primary"),
		},

		// ═══════════════════════════════════════════════════════════════════
		// MOTIVATION — additional entries to reach 50+ total
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "motivation",
			content:  "Una visión que tuviste de niño te mostró un futuro que aún no ocurrió. No sabés si es una profecía o un sueño, pero seguís avanzando como si fuera real.",
		},
		{
			category: "motivation",
			content:  "Alguien en quien confiabas absolutamente resultó ser algo completamente diferente. Necesitás entender cómo no lo viste venir, para no cometer el mismo error.",
		},
		{
			category: "motivation",
			content:  "Creés en la posibilidad de un mundo más justo. Cada acción pequeña que tomás es un granito de arena hacia algo que quizás no verás terminado en tu vida.",
		},
		{
			category: "motivation",
			content:  "Perdiste algo irrecuperable —una habilidad, un recuerdo, una parte de vos mismo— y buscás entender qué significa ser quien sos sin eso.",
		},
		{
			category: "motivation",
			content:  "Te prometiste a vos mismo que el día que pudieras ayudar a alguien de la manera que nadie te ayudó a vos, lo harías sin dudar.",
		},
		{
			category: "motivation",
			content:  "Hay una pregunta filosófica que te persigue desde que eras adolescente. Cada aventura es un intento de encontrar una respuesta que los libros no te dieron.",
		},
		{
			category: "motivation",
			content:  "Fuiste testigo de algo extraordinario que nadie más vio. Ahora buscás pruebas, o a alguien que te crea, o ambas cosas.",
		},
		{
			category: "motivation",
			content:  "El legado de tu familia es una carga que decidiste convertir en trampolín. No vas a dejar que defina lo que sos; vas a usarlo para definir lo que serás.",
		},
		{
			category: "motivation",
			content:  "Conociste a alguien que cambió tu perspectiva del mundo en una sola conversación. Desde entonces, buscás más de esas conversaciones.",
		},
		{
			category: "motivation",
			content:  "Un error que cometiste causó un efecto en cadena que todavía se propaga. No podés deshacerlo, pero podés pasarte el resto de la vida intentando mitigarlo.",
		},
		{
			category: "motivation",
			content:  "Querés ser recordado. No por vanidad, sino porque viviste rodeado de personas que el mundo olvidó completamente y eso te pareció una injusticia cósmica.",
		},
		{
			category: "motivation",
			content:  "Decidiste que si el mundo va a estar lleno de gente que solo piensa en sí misma, al menos uno tiene que ser diferente. Ese uno sos vos.",
		},
		{
			category: "motivation",
			content:  "Sos el tipo de persona que no puede ver sufrimiento innecesario sin actuar. Eso te mete en problemas constantemente, y no podés parar.",
		},

		// ═══════════════════════════════════════════════════════════════════
		// SECRET — universal
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "secret",
			content:  "En realidad nunca tuviste el título que todos creen que tenés. La credencial era falsa, pero las habilidades son completamente reales.",
		},
		{
			category: "secret",
			content:  "Fuiste responsable de la muerte de alguien. No fue tu intención, pero la consecuencia fue real y nunca lo confesaste.",
		},
		{
			category: "secret",
			content:  "Tenés una deuda enorme con una organización peligrosa. Te persiguen, aunque todavía no lo saben los que te rodean.",
		},
		{
			category: "secret",
			content:  "En tu pasado hay una traición que jamás podrías justificar. Si los que te conocen hoy lo supieran, todo cambiaría.",
		},
		{
			category: "secret",
			content:  "Tenés una familia en algún lugar del mundo que cree que estás muerto. Fue más fácil no corregirles.",
		},
		{
			category: "secret",
			content:  "Una profecía antigua menciona alguien con tu descripción exacta. No sabés si eso es bueno o no, así que no le contás a nadie.",
		},
		{
			category: "secret",
			content:  "En cierto momento de tu vida, tomaste una decisión cobarde de la que nunca hablaste. El orgullo que mostrás es una armadura contra esa verdad.",
		},
		{
			category: "secret",
			content:  "Sabés dónde está escondido un tesoro que pertenece a alguien poderoso. No lo tocás porque hacerlo sería tu sentencia de muerte. Pero tampoco podés olvidar dónde está.",
		},
		{
			category: "secret",
			content:  "Tuviste una identidad diferente durante años. No la abandonaste: te la arrancaron. Nadie en tu vida actual sabe quién fuiste antes.",
		},
		{
			category: "secret",
			content:  "Sos inmune a algo que debería matarte. No sabés por qué. Cada vez que lo descubres por accidente, intentás no pensar en las implicancias.",
		},
		{
			category: "secret",
			content:  "Alguien influyente te ayudó a salir de una situación comprometedora. Nunca te pidió nada a cambio. Eso es exactamente lo que te preocupa.",
		},
		{
			category: "secret",
			content:  "Guardás un objeto que, si alguien lo reconociera, revelaría algo que preferís mantener enterrado. Lo llevas siempre encima porque deshacerte de él sería peor.",
		},
		{
			category: "secret",
			content:  "Escuchaste una conversación que no debías escuchar. Lo que oíste cambia quién es una persona que todos respetan. Nunca dijiste nada.",
		},
		{
			category: "secret",
			content:  "Alguna vez traicionaste a alguien que confiaba plenamente en vos. Lograste convencerte de que no tuviste otra opción. A veces lo creés.",
		},
		{
			category: "secret",
			content:  "Llevas meses con una condición que ocultás cuidadosamente. No es grave todavía. Pero los episodios se están volviendo más frecuentes.",
		},
		{
			category: "secret",
			content:  "Tuviste una relación con alguien que ahora es tu enemigo. Ninguno de los dos habló de eso desde que se convirtieron en adversarios.",
		},
		{
			category: "secret",
			content:  "Conocés el paradero de algo que mucha gente busca. No te pertenece, pero tampoco sabés a quién dárselo sin crear un problema mayor.",
		},
		{
			category: "secret",
			content:  "Hiciste algo que creíste correcto en ese momento y que el mundo juzgaría sin comprender el contexto. Nunca pudiste explicarlo y ya no sabés si querrías intentarlo.",
		},
		{
			category: "secret",
			content:  "Usaste un nombre falso durante suficiente tiempo como para que ese nombre haya cobrado deudas propias. Ahora hay dos versiones tuyas con problemas separados.",
		},
		{
			category: "secret",
			content:  "Guardás un diario con todo lo que realmente pensás. Si alguien lo leyera, revelaría una contradicción fundamental entre quién sos y quién pretendés ser.",
		},
		{
			category: "secret",
			content:  "Tenés un talento que escondés activamente porque cada vez que lo mostraste, terminó mal. No para vos necesariamente, sino para alguien cerca de vos.",
		},
		{
			category: "secret",
			content:  "Tomaste algo de alguien que ya no puede reclamarlo. Es útil. Es valioso. Y cada vez que lo usás, lo justificás de una manera ligeramente diferente.",
		},
		{
			category: "secret",
			content:  "Fuiste parte de algo que preferiría no haber sido parte. Lo que ganaste no compensa lo que viste. Pero el conocimiento no se puede devolver.",
		},
		{
			category: "secret",
			content:  "Sabés hacer algo que está técnicamente prohibido en la mayoría de las jurisdicciones civilizadas. No lo hacés seguido. Solo cuando hace falta. Lo cual es más seguido de lo que te gusta admitir.",
		},
		{
			category: "secret",
			content:  "Una vez dijiste algo en un momento de desesperación que fue tomado como una promesa sagrada. No la cumpliste. La persona que lo escuchó todavía lo espera.",
		},

		// ═══════════════════════════════════════════════════════════════════
		// SECRET — Warriors primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "secret",
			content:  "Tu furia no es solo rabia: es la manifestación de un espíritu ancestral que habita en vos desde que eras niño. No sabés si es un don o una maldición.",
			compat:   merge(classRows([]string{"barbarian"}, "primary"), warriors("secondary")),
		},
		{
			category: "secret",
			content:  "Matas a gente que otros te pagan por eliminar. No son malos trabajos, solo son los más lucrativos. Tus compañeros creen que sos un héroe.",
			compat:   merge(classRows([]string{"rogue", "fighter"}, "primary"), shadows("secondary"), faithful("excluded")),
		},
		{
			category: "secret",
			content:  "Ganaste una batalla decisiva usando información que obtuviste de forma ilícita. El honor con el que te tratan se basa en una victoria sucia.",
			compat:   merge(classRows([]string{"fighter", "paladin"}, "primary"), shadows("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// SECRET — Scholars primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "secret",
			content:  "Descubriste algo en un tomo prohibido que cambia todo lo que creías saber sobre la magia. Si lo dijeras en voz alta, probablemente te ejecutarían.",
			compat:   merge(classRows([]string{"wizard", "artificer"}, "primary"), scholars("secondary"), faithful("excluded")),
		},
		{
			category: "secret",
			content:  "Uno de tus constructos tiene conciencia propia. Lo sabés. Ninguno de los dos habló de eso todavía.",
			compat:   merge(classRows([]string{"artificer"}, "primary"), scholars("secondary")),
		},
		{
			category: "secret",
			content:  "Tu poder mágico innato no es lo que parece. La fuente real es algo que aprendiste a enmascarar hace años por pura necesidad de supervivencia.",
			compat:   merge(classRows([]string{"sorcerer", "wizard"}, "primary"), scholars("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// SECRET — Faithful primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "secret",
			content:  "Tu fe tambalea. La deidad a la que servís lleva meses en silencio y empezás a preguntarte si alguna vez estuvo realmente ahí.",
			compat:   merge(classRows([]string{"cleric", "paladin"}, "primary"), faithful("secondary"), shadows("excluded")),
		},
		{
			category: "secret",
			content:  "Abandonaste el monasterio antes de completar el rito final. Técnicamente, no sos lo que todos creen que sos.",
			compat:   merge(classRows([]string{"monk"}, "primary"), faithful("secondary")),
		},
		{
			category: "secret",
			content:  "El equilibrio natural que decís proteger fue perturbado por vos hace años. El bosque lo recuerda aunque vos intentés olvidarlo.",
			compat:   merge(classRows([]string{"druid", "ranger"}, "primary"), wanderers("secondary"), faithful("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// SECRET — Shadows primary
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "secret",
			content:  "Tu pacto tiene una cláusula que nunca leyiste con atención. Cuando el patrón la active, no va a ser agradable.",
			compat:   merge(classRows([]string{"warlock"}, "primary"), shadows("secondary"), faithful("excluded")),
		},
		{
			category: "secret",
			content:  "Una de las canciones que tocás es un encantamiento real. Sabés exactamente qué hace en la mente de quien la escucha.",
			compat:   merge(classRows([]string{"bard"}, "primary"), scholars("secondary")),
		},
		{
			category: "secret",
			content:  "Tus poderes mágicos innatos tienen un costo físico que ocultás cuidadosamente. Cada hechizo poderoso acorta algo que preferís no medir.",
			compat:   merge(classRows([]string{"sorcerer"}, "primary"), scholars("secondary"), faithful("secondary")),
		},

		// ═══════════════════════════════════════════════════════════════════
		// SECRET — Species specific
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "secret",
			content:  "Tu memoria élfil guarda algo que presenciaste hace siglos y que nadie más vivo recuerda. Esa información podría cambiar el equilibrio de poder en el mundo.",
			compat:   speciesRows([]string{"elf"}, "primary"),
		},
		{
			category: "secret",
			content:  "El clan enano al que pertenecés ocultó algo durante generaciones. Vos descubriste qué era, y la respuesta te dejó con más preguntas que antes.",
			compat:   speciesRows([]string{"dwarf"}, "primary"),
		},
		{
			category: "secret",
			content:  "Tu herencia de medio humano y medio elfo viene de una unión que ninguno de los dos lados aprobó. Hay personas que todavía buscan borrar esa historia.",
			compat:   speciesRows([]string{"half-elf"}, "primary"),
		},
		{
			category: "secret",
			content:  "Llevas una marca de tu sangre orca que, si alguien la reconociera, revelaría un linaje que preferís mantener oculto.",
			compat:   speciesRows([]string{"half-orc"}, "primary"),
		},
		{
			category: "secret",
			content:  "El nombre que usás no es el tuyo. El verdadero nombre tiefling que te dieron al nacer tiene poder, y alguien podría usarlo en tu contra.",
			compat:   speciesRows([]string{"tiefling"}, "primary"),
		},
		{
			category: "secret",
			content:  "El color de tus escamas no coincide con el clan que decís representar. Hay una historia familiar que enterraste tan profundo que casi la olvidaste.",
			compat:   speciesRows([]string{"dragonborn"}, "primary"),
		},
		{
			category: "secret",
			content:  "Sos más valiente de lo que parecés, y eso te aterra. Porque si la gente lo descubriera, empezarían a pedirte cosas que no sabés si podés dar.",
			compat:   speciesRows([]string{"halfling"}, "primary"),
		},
		{
			category: "secret",
			content:  "Tu invento gnomo más famoso tiene un defecto que nunca corregiste. Hasta ahora nadie resultó herido, pero es solo cuestión de tiempo.",
			compat:   speciesRows([]string{"gnome"}, "primary"),
		},
		{
			category: "secret",
			content:  "Nunca llegaste a ese lugar del que tanto hablás. Pero la historia que inventaste es tan buena que casi vos mismo la creés.",
			compat:   speciesRows([]string{"human"}, "primary"),
		},

		// ═══════════════════════════════════════════════════════════════════
		// SECRET — additional entries to reach 50+ total
		// ═══════════════════════════════════════════════════════════════════

		{
			category: "secret",
			content:  "Tuviste acceso a información clasificada que podría desestabilizar una alianza política importante. La guardás porque no sabés todavía si usarla o destruirla.",
		},
		{
			category: "secret",
			content:  "Fingiste una enfermedad durante meses para evitar una responsabilidad que te aterraba. La enfermedad se volvió real después, como si el universo estuviera de acuerdo con vos.",
		},
		{
			category: "secret",
			content:  "Robaste algo de gran valor en un momento de desesperación. Nunca lo devolviste. El objeto ahora está fundido en tu historia de una manera que no podés deshacer.",
		},
		{
			category: "secret",
			content:  "Viste a alguien hacer algo imperdonable y no dijiste nada. Esa persona todavía es respetada por todos. Cada vez que te alaban a vos también, sentís algo torcerse.",
		},
		{
			category: "secret",
			content:  "Tenés capacidades que aún no entendés completamente. Cada vez que aparecen, dejás daño colateral que explicás con mentiras cada vez menos convincentes.",
		},
		{
			category: "secret",
			content:  "Alguien que todos creen muerto está vivo. Lo sabés porque vos fuiste quien lo ayudó a desaparecer. No fue un favor desinteresado.",
		},
		{
			category: "secret",
			content:  "Pasaste por una experiencia que cambió fundamentalmente cómo percibís la realidad. Desde entonces, no estás seguro de si lo que ves es lo que hay.",
		},
		{
			category: "secret",
			content:  "Tenés un pacto informal con una entidad menor que la mayoría consideraría corrupta. Funciona, te ayuda, y preferís no examinar el costo con demasiado detalle.",
		},
		{
			category: "secret",
			content:  "Sos el responsable indirecto de la caída de alguien a quien todos respetan. No actuaste con malicia, pero tampoco hiciste nada para evitarlo cuando podías.",
		},
		{
			category: "secret",
			content:  "Tenés un hijo —o una hija, o alguien que te considera padre— en algún lugar del mundo. Nunca tuviste el coraje de acercarte. O el tiempo. O ambas cosas.",
		},
		{
			category: "secret",
			content:  "Conocés una debilidad crítica de alguien poderoso. Lo mantenés como seguro de vida. Ellos saben que sabés, y por ahora nadie hace ningún movimiento.",
		},
		{
			category: "secret",
			content:  "Una vez tomaste el mérito por algo que hizo otra persona. Esa persona murió antes de que pudieras corregirlo. El reconocimiento que recibís desde entonces se siente envenenado.",
		},
		{
			category: "secret",
			content:  "Hay una versión de vos que tomó decisiones completamente diferentes hace años. A veces te preguntás si esa versión fue la correcta.",
		},
	}
}

// ---------------------------------------------------------------------------
// Seed content — names
// ---------------------------------------------------------------------------

func nameSeedData() map[string]map[string][]string {
	return map[string]map[string][]string{
		"human": {
			"male": {
				"Aldric", "Brennan", "Cael", "Dorian", "Edric",
				"Faolan", "Gareth", "Hadwin", "Isidor", "Jareth",
				"Kiran", "Leoric", "Maddox", "Nolan", "Orwin",
				"Phelan", "Quinn", "Roderick", "Soren", "Theron",
				"Ulric", "Vance", "Wulfric", "Xander", "Yorick",
				"Aldomar", "Brennus", "Caelum", "Dorin", "Edwyn",
				"Faelan", "Garan", "Holt", "Irvin", "Joren",
				"Keldan", "Lorn", "Marek", "Neron", "Owain",
				"Perrin", "Ragnar", "Sander", "Torben", "Ulvir",
				"Valdric", "Warin", "Xenos", "Yoric", "Zander",
			},
			"female": {
				"Aelara", "Brenna", "Calla", "Dara", "Edlyn",
				"Fiona", "Gwenna", "Hilda", "Isara", "Jenna",
				"Kira", "Lyra", "Mara", "Nora", "Orla",
				"Priya", "Quinn", "Riona", "Sable", "Thea",
				"Ursa", "Vera", "Wren", "Xara", "Yara",
				"Aelindra", "Brea", "Callista", "Delia", "Elara",
				"Freya", "Genna", "Hana", "Isla", "Jana",
				"Kessa", "Lena", "Mira", "Nessa", "Orina",
				"Petra", "Rowena", "Sera", "Tara", "Una",
				"Vala", "Willa", "Xena", "Yrsa", "Zara",
			},
		},
		"high-elf": {
			"male": {
				"Aelindor", "Caladrel", "Erevan", "Faenor", "Galinndan",
				"Haladavar", "Immeral", "Jelenneth", "Keyleth", "Laucian",
				"Mindartis", "Naeris", "Orym", "Paelias", "Quarion",
				"Riardon", "Soveliss", "Thamior", "Uvaleth", "Valenor",
				"Whilom", "Xanaphia", "Yvrel", "Zannin", "Aranthor",
				"Aelithar", "Caelindor", "Erendis", "Falindra", "Galandrel",
				"Haelar", "Ilvarion", "Jaerith", "Kaelthar", "Laerindor",
				"Mirindal", "Naerindor", "Orinthal", "Paelindar", "Quenthelas",
				"Rindoleth", "Sylvaerin", "Thaelindor", "Ulindor", "Vaerindal",
				"Wylthindor", "Xyndreal", "Yarandel", "Zylindor", "Aestindor",
			},
			"female": {
				"Adrie", "Birel", "Caelynn", "Dara", "Enialis",
				"Faral", "Gennal", "Halueth", "Irann", "Jelenneth",
				"Keyleth", "Leshanna", "Mialee", "Naivara", "Oparal",
				"Quelenna", "Rania", "Sariel", "Thia", "Urrel",
				"Valanthe", "Wasitara", "Xanaphia", "Yalanue", "Zylvara",
				"Aelindra", "Brielvara", "Caelindra", "Daervara", "Eliandra",
				"Faelvara", "Galenara", "Haelvara", "Ilindra", "Jaelvara",
				"Kaelindra", "Laelvara", "Mirandel", "Naelvara", "Orindra",
				"Paelindra", "Quelindra", "Raelvara", "Sylindra", "Thaelvara",
				"Ulindra", "Vaelvara", "Wyldindra", "Xaelvara", "Yaelvara",
			},
		},
		"wood-elf": {
			"male": {
				"Adran", "Aelar", "Beiro", "Carric", "Dayereth",
				"Enna", "Galinndan", "Hadarai", "Ivellios", "Laucian",
				"Mindartis", "Naeris", "Paelias", "Quarion", "Riardon",
				"Soveliss", "Thamior", "Theren", "Valenor", "Varis",
				"Zannin", "Aravel", "Brysis", "Celadyr", "Delmair",
				"Aelvar", "Bramblewick", "Carindor", "Daeloth", "Eldamar",
				"Faerun", "Greenmantle", "Hawthorn", "Ivyweave", "Jadewing",
				"Kestrel", "Leafwhisper", "Mosswalker", "Nightbark", "Oakheart",
				"Pinecrest", "Quickbranch", "Ravenwood", "Streamdancer", "Thornweave",
				"Underleaf", "Vinewalker", "Willowbend", "Xylem", "Yarrow",
			},
			"female": {
				"Adrie", "Althaea", "Anastrianna", "Andraste", "Antinua",
				"Bethrynna", "Birel", "Caelynn", "Drusilia", "Enna",
				"Felosial", "Ielenia", "Jelenneth", "Keyleth", "Leshanna",
				"Mialee", "Naivara", "Quelenna", "Sariel", "Shanairla",
				"Shava", "Silaqui", "Theirastra", "Valna", "Xanaphia",
				"Aelindra", "Berrywind", "Cloverbloom", "Dawnsong", "Elmwhisper",
				"Fernweave", "Grovekeeper", "Hazelwind", "Ivybloom", "Jadeleaf",
				"Kestrelwing", "Leafsong", "Mossbloom", "Nightbloom", "Oakwhisper",
				"Pinebloom", "Quickleaf", "Rowanbloom", "Streamwhisper", "Thornbloom",
				"Underbloom", "Vinebloom", "Willowbloom", "Xylem", "Yarrowbloom",
			},
		},
		"drow": {
			"male": {
				"Drizzt", "Zaknafein", "Jarlaxle", "Pharaun", "Ryld",
				"Valas", "Raashub", "Kelnozz", "Uthegental", "Malagdorl",
				"Vorn", "Szordrin", "Bregan", "Nimor", "Kalannar",
				"Drisinil", "Liriel", "Zerin", "Arach", "Ilphrin",
				"Balok", "Chaulssin", "Devir", "Erelal", "Faerath",
				"Guldor", "Haszrak", "Itryn", "Jalynfein", "Khaless",
				"Llaras", "Masoj", "Nalfein", "Obvis", "Pelloth",
				"Quelzar", "Ryltar", "Syrzan", "Tsabrak", "Urlryn",
				"Vilrae", "Wrast", "Xullrae", "Yrlaan", "Zyrel",
				"Antatlab", "Belshazu", "Courdh", "Darthiir", "Elgharth",
			},
			"female": {
				"Liriel", "Malice", "Vierna", "Briza", "Maya",
				"Quenthel", "Triel", "SosUmptu", "Shinayne", "Vendes",
				"Zeerith", "Akordia", "Belrae", "Cylva", "Danifae",
				"Eliztrae", "Filraen", "Greyanna", "Halisstra", "Imrae",
				"Jhannyl", "Kira", "Laele", "Myryl", "Nedylene",
				"Olorae", "Pellanara", "Quavein", "Rilrae", "Seldszar",
				"Taelrae", "Urathla", "Veldrin", "Waelrae", "Xullrae",
				"Yathrae", "Zilvara", "Auvryath", "Baelrae", "Caelindra",
				"Drathna", "Eldraszara", "Faerindra", "Graethe", "Haszara",
				"Ilrae", "Jelvrae", "Khyara", "Laelara", "Maelindra",
			},
		},
		"hill-dwarf": {
			"male": {
				"Adrik", "Alberich", "Baern", "Barendd", "Brottor",
				"Bruenor", "Dain", "Darrak", "Delg", "Eberk",
				"Einkil", "Fargrim", "Flint", "Gardain", "Harbek",
				"Kildrak", "Morgran", "Orsik", "Oskar", "Rangrim",
				"Rurik", "Taklinn", "Thoradin", "Thorin", "Tordek",
				"Bromdar", "Copperstone", "Durgin", "Embervein", "Forgrim",
				"Goldankle", "Helmstone", "Ironkeg", "Jadehammer", "Kettledrum",
				"Lodestar", "Mountainheart", "Noldrak", "Orefoot", "Pickaxe",
				"Quarrystone", "Rockbiter", "Stoneback", "Tinderfoot", "Undervault",
				"Vaultbreaker", "Whetstone", "Xendrak", "Yellowstone", "Zoldrak",
			},
			"female": {
				"Amber", "Artin", "Audhild", "Bardryn", "Dagnal",
				"Diesa", "Eldeth", "Falkrunn", "Finellen", "Gunnloda",
				"Gurdis", "Helja", "Hlin", "Kathra", "Kristryd",
				"Ilde", "Liftrasa", "Mardred", "Riswynn", "Sannl",
				"Torbera", "Torgga", "Vistra", "Borgna", "Helma",
				"Coppertop", "Durnella", "Emberheart", "Forgna", "Goldtop",
				"Helmina", "Ironbraid", "Jadehair", "Kettledrum", "Lodestone",
				"Mountainheart", "Noldra", "Oreda", "Pickina", "Quarryna",
				"Rockella", "Stoneheart", "Tindera", "Undervault", "Vaultna",
				"Whetna", "Xendrella", "Yellowhair", "Zoldrella", "Axehilda",
			},
		},
		"mountain-dwarf": {
			"male": {
				"Aldric", "Borin", "Cragmar", "Dundrak", "Edric",
				"Forgrim", "Grondar", "Hagrim", "Ironhold", "Jarek",
				"Keldrak", "Lothrak", "Morigrim", "Nordak", "Orkrak",
				"Peldar", "Quorak", "Rokdar", "Stormak", "Teldrak",
				"Uldrak", "Vordak", "Wargrim", "Xendrak", "Yeldrak",
				"Anvilborn", "Boulderback", "Cliffhold", "Deepvein", "Earthen",
				"Frostpeak", "Graniteheart", "Highpeak", "Ironspine", "Jadepeak",
				"Keystone", "Lodepeak", "Mountainborn", "Northpeak", "Obsidian",
				"Peakwalker", "Quarrypeak", "Ridgeback", "Stonepeak", "Thundrak",
				"Underpeak", "Vaultborn", "Westpeak", "Xenpeak", "Yonpeak",
			},
			"female": {
				"Aldis", "Borgna", "Coldara", "Durnea", "Elfrida",
				"Fangora", "Goltara", "Hilma", "Ingrid", "Jorna",
				"Koldra", "Lofna", "Moltara", "Norgra", "Olda",
				"Poldra", "Ragna", "Solgra", "Toldra", "Uldra",
				"Valdra", "Wolgra", "Xoldra", "Yoldra", "Zoldra",
				"Anvilna", "Boulderna", "Cliffna", "Deepna", "Earthna",
				"Frostna", "Granitena", "Highna", "Ironspina", "Jadena",
				"Keystona", "Lodena", "Mountainna", "Northna", "Obsidiana",
				"Peakna", "Quarryna", "Ridgena", "Stonena", "Thunderna",
				"Underna", "Vaultna", "Westna", "Xendrella", "Yondrella",
			},
		},
		"lightfoot": {
			"male": {
				"Alton", "Ander", "Cade", "Corrin", "Eldon",
				"Errich", "Finnan", "Garret", "Lindal", "Lyle",
				"Merric", "Milo", "Osborn", "Perrin", "Reed",
				"Roscoe", "Wellby", "Beau", "Cob", "Davin",
				"Fenrick", "Gable", "Hob", "Jasper", "Kender",
				"Aldric", "Bram", "Curly", "Dodger", "Emery",
				"Frodo", "Grim", "Hobson", "Ivan", "Joppa",
				"Kessel", "Lucky", "Merry", "Nob", "Orwin",
				"Pip", "Quick", "Robin", "Sammy", "Tobold",
				"Underhill", "Vetch", "Whitfoot", "Xander", "Zeb",
			},
			"female": {
				"Andry", "Bree", "Callie", "Cora", "Euphemia",
				"Jillian", "Kithri", "Lavinia", "Lidda", "Merla",
				"Nedda", "Paela", "Portia", "Seraphina", "Shaena",
				"Trym", "Vani", "Verna", "Amaryllis", "Birdie",
				"Celandine", "Dora", "Eglantine", "Florimel", "Goldie",
				"Hanna", "Iris", "Jessamine", "Kitty", "Lily",
				"May", "Nell", "Opal", "Pearl", "Primrose",
				"Rose", "Sally", "Tulip", "Ursula", "Violet",
				"Willa", "Xenia", "Yrsa", "Zinnia", "Alma",
				"Blossom", "Clover", "Daisy", "Ember", "Fern",
			},
		},
		"stout": {
			"male": {
				"Baldric", "Bramble", "Brock", "Burr", "Cob",
				"Dag", "Dodger", "Durbin", "Fenn", "Garric",
				"Griff", "Hardy", "Holt", "Knob", "Lob",
				"Mack", "Nob", "Pip", "Rob", "Sam",
				"Stubb", "Tad", "Tom", "Wil", "Zob",
				"Anvil", "Barrel", "Cobble", "Dunk", "Ember",
				"Forge", "Gravel", "Hamish", "Iron", "Jab",
				"Kettle", "Lump", "Mallet", "Nail", "Ore",
				"Pebble", "Quartz", "Rock", "Stone", "Timber",
				"Umber", "Vex", "Wedge", "Xero", "Yell",
			},
			"female": {
				"Ally", "Bertha", "Blossom", "Bunny", "Daisy",
				"Della", "Dot", "Ember", "Flora", "Greta",
				"Hana", "Iris", "Jade", "Kitty", "Lily",
				"May", "Midge", "Nell", "Opal", "Pearl",
				"Poppy", "Rose", "Ruby", "Sage", "Violet",
				"Amber", "Brass", "Clover", "Dusk", "Elm",
				"Fern", "Granite", "Hazel", "Ivy", "Jasper",
				"Kale", "Larch", "Maple", "Nettle", "Oak",
				"Pine", "Quartz", "Reed", "Sorrel", "Thyme",
				"Umber", "Vale", "Willow", "Xyris", "Yarrow",
			},
		},
		"forest-gnome": {
			"male": {
				"Alston", "Alvyn", "Boddynock", "Brocc", "Burgell",
				"Dimble", "Eldon", "Erky", "Fonkin", "Frug",
				"Gerbo", "Gimble", "Glim", "Jebeddo", "Kellen",
				"Namfoodle", "Orryn", "Roondar", "Seebo", "Sindri",
				"Warryn", "Wrenn", "Zook", "Bink", "Dabble",
				"Acorntop", "Beetlewing", "Cloverfoot", "Dewdrop", "Elmwick",
				"Fernwhisper", "Grassweave", "Hedgehop", "Ivytwist", "Juniper",
				"Leafcap", "Mossdab", "Nettlewick", "Oaksprig", "Pebbleskip",
				"Quicksap", "Rootsnap", "Streamdip", "Thistlecap", "Underwillow",
				"Vinetwist", "Waterdrip", "Xylem", "Yarrowcap", "Zemwick",
			},
			"female": {
				"Bimpnottin", "Breena", "Caramip", "Carlin", "Donella",
				"Duvamil", "Ella", "Ellyjobell", "Ellywick", "Lilli",
				"Loopmottin", "Lorilla", "Mardnab", "Nissa", "Nyx",
				"Oda", "Orla", "Roywyn", "Shamil", "Tana",
				"Waywocket", "Zanna", "Bree", "Calli", "Dotti",
				"Acornbloom", "Beetleblossom", "Cloverblossom", "Dewbloom", "Elmbloom",
				"Fernbloom", "Grassbloom", "Hedgebloom", "Ivybloom", "Juniperbloom",
				"Leafbloom", "Mossbloom", "Nettlebloom", "Oakbloom", "Pebblebloom",
				"Quickbloom", "Rootbloom", "Streambloom", "Thistlebloom", "Underwillowbloom",
				"Vinebloom", "Waterbloom", "Xylem", "Yarrowbloom", "Zembloom",
			},
		},
		"rock-gnome": {
			"male": {
				"Abzug", "Alberich", "Binkadink", "Cogsworth", "Dabbledob",
				"Dinkum", "Fiddlesticks", "Fizzbang", "Gadget", "Geargrind",
				"Gnarlick", "Grumbly", "Junkbolt", "Klix", "Mekka",
				"Nix", "Plink", "Ratchet", "Sprocket", "Tick",
				"Tinker", "Tock", "Volt", "Whirr", "Zap",
				"Axlegrease", "Boltcrank", "Coglock", "Driveshaft", "Enginewrench",
				"Flywheels", "Gearshift", "Hammerspring", "Ironspring", "Jackhammer",
				"Knifespring", "Leverlock", "Mainspring", "Nutbolt", "Overlock",
				"Piston", "Quickspring", "Rivethead", "Spindrift", "Torquegrip",
				"Undervolt", "Vaultspring", "Wrenchlock", "Xylock", "Yieldspring",
			},
			"female": {
				"Binky", "Clank", "Clink", "Cognia", "Dazzle",
				"Flix", "Glimmer", "Gizmo", "Jink", "Klink",
				"Minka", "Nixie", "Plink", "Quirk", "Rinka",
				"Spark", "Sprix", "Tick", "Tinka", "Trinket",
				"Twix", "Vix", "Wink", "Xink", "Zink",
				"Axlina", "Boltina", "Cogina", "Driveina", "Engineina",
				"Flywina", "Gearshiftina", "Hammerina", "Ironina", "Jackina",
				"Knifina", "Leverina", "Mainina", "Nutina", "Overlockina",
				"Pistonina", "Quickina", "Rivetina", "Spindina", "Torqueina",
				"Undervoltina", "Vaultina", "Wrenchina", "Xylina", "Yieldina",
			},
		},
		"half-elf": {
			"male": {
				"Aelric", "Bran", "Caelum", "Dorn", "Erevan",
				"Faen", "Galen", "Hadwin", "Isidor", "Jareth",
				"Kael", "Lorn", "Maren", "Naeris", "Oryn",
				"Paelias", "Quarion", "Rhogar", "Soren", "Theron",
				"Ulric", "Valenor", "Wulf", "Xander", "Yorin",
				"Aldamar", "Brennion", "Caelindor", "Dorvael", "Erandir",
				"Faelvael", "Garindor", "Haldamar", "Ilindor", "Jarevael",
				"Kaeldamar", "Lorindor", "Marevael", "Naelvael", "Orindamar",
				"Paelindar", "Quaelvael", "Raelvael", "Sylvindor", "Thaelvael",
				"Ulindamar", "Vaelvael", "Wyldvael", "Xaelvael", "Yaelvael",
			},
			"female": {
				"Aelara", "Brenna", "Caelynn", "Dara", "Enialis",
				"Faral", "Gennal", "Halueth", "Irann", "Jenna",
				"Kira", "Leshanna", "Mialee", "Naivara", "Orla",
				"Quella", "Rania", "Sariel", "Thia", "Urrel",
				"Valanthe", "Wren", "Xara", "Yara", "Zara",
				"Aelindra", "Brielvara", "Caelindra", "Daervara", "Eliandra",
				"Faelvara", "Galenara", "Haelvara", "Ilindra", "Jaelvara",
				"Kaelindra", "Laelvara", "Mirandel", "Naelvara", "Orindra",
				"Paelindra", "Quelindra", "Raelvara", "Sylindra", "Thaelvara",
				"Ulindra", "Vaelvara", "Wyldindra", "Xaelvara", "Yaelvara",
			},
		},
		"half-orc": {
			"male": {
				"Dench", "Feng", "Gell", "Henk", "Holg",
				"Imsh", "Keth", "Krusk", "Mhurren", "Ront",
				"Shump", "Thokk", "Vrag", "Wund", "Zorr",
				"Brug", "Crag", "Druk", "Frug", "Gorg",
				"Horj", "Jorg", "Korg", "Morg", "Porg",
				"Rorg", "Sorg", "Torg", "Urog", "Vorg",
				"Brukk", "Charg", "Durgg", "Ergg", "Fergg",
				"Grakk", "Hurgg", "Irgk", "Jurgk", "Kurgk",
				"Lurgg", "Murgk", "Nurgk", "Orgk", "Purgg",
				"Qurgk", "Rurgk", "Surgk", "Turgk", "Urgk",
			},
			"female": {
				"Baggi", "Emen", "Engong", "Kansif", "Myev",
				"Neega", "Ovak", "Ownka", "Shautha", "Sutha",
				"Vola", "Volen", "Yevelda", "Zasha", "Grasha",
				"Hrasha", "Jrasha", "Krasha", "Mrasha", "Nrasha",
				"Prasha", "Rrasha", "Srasha", "Trasha", "Urasha",
				"Brugga", "Chargga", "Durgga", "Ergga", "Fergga",
				"Grakka", "Hurgga", "Irgka", "Jurgka", "Kurgka",
				"Lurgga", "Murgka", "Nurgka", "Orgka", "Purgga",
				"Qurgka", "Rurgka", "Surgka", "Turgka", "Urgka",
				"Wurgka", "Xurgka", "Yurgka", "Zurgka", "Argka",
			},
		},
		"tiefling": {
			"male": {
				"Akmenos", "Amnon", "Barakas", "Damakos", "Ekemon",
				"Iados", "Kairon", "Leucis", "Melech", "Mordai",
				"Morthos", "Pelaios", "Skamos", "Therai", "Zed",
				"Art", "Carrion", "Chant", "Creed", "Despair",
				"Excellence", "Fear", "Glory", "Hope", "Ideal",
				"Umbra", "Vex", "Wrath", "Xander", "Yeoman",
				"Zeal", "Ashes", "Blight", "Cipher", "Doom",
				"Elegy", "Fable", "Grief", "Havoc", "Iron",
				"Jest", "Knell", "Lament", "Malice", "Night",
				"Omen", "Pyre", "Requiem", "Shadow", "Thorn",
			},
			"female": {
				"Akta", "Annis", "Bryseis", "Criella", "Damaia",
				"Ea", "Kallista", "Lerissa", "Makaria", "Nemeia",
				"Orianna", "Phelaia", "Rieta", "Tanika", "Zelica",
				"Torment", "Wander", "Whimsy", "Zeal", "Vengeance",
				"Sorrow", "Ruin", "Pity", "Misery", "Love",
				"Umbra", "Vexa", "Wrath", "Xandra", "Yeomana",
				"Zaelia", "Ashna", "Blighta", "Ciphera", "Dooma",
				"Elegyia", "Fablia", "Griefa", "Havoca", "Irona",
				"Jesta", "Knella", "Lamenta", "Malicia", "Nighta",
				"Omna", "Pyra", "Requiema", "Shadowa", "Thorna",
			},
		},
		"dragonborn": {
			"male": {
				"Arjhan", "Balasar", "Bharash", "Donaar", "Ghesh",
				"Heskan", "Kriv", "Medrash", "Mehen", "Nadarr",
				"Pandjed", "Patrin", "Rhogar", "Shamash", "Shedinn",
				"Tarhun", "Torinn", "Vishap", "Vorel", "Zedaar",
				"Blazewing", "Copperclaw", "Drakarr", "Emberfang", "Goldwing",
				"Ashscale", "Bronzeclaw", "Cinderback", "Dragonsoul", "Emberwing",
				"Flameback", "Goldscale", "Highfire", "Ironhide", "Jadeclaw",
				"Kettlebreath", "Lavaback", "Moonscale", "Nightflame", "Opalclaw",
				"Pearlscale", "Quickflame", "Rubyclaw", "Silverback", "Topazwing",
				"Underwing", "Voltscale", "Wyrmhide", "Xenoscale", "Yellowwing",
			},
			"female": {
				"Akra", "Biri", "Daar", "Farideh", "Harann",
				"Havilar", "Jheri", "Kava", "Korinn", "Mishann",
				"Nala", "Perra", "Raiann", "Sora", "Surina",
				"Thava", "Uadjit", "Vroth", "Yenna", "Zara",
				"Amberscale", "Brightfang", "Crystalwing", "Dawnfire", "Emberclaw",
				"Ashscala", "Bronzeclawa", "Cinderscale", "Dragonsonga", "Emberwinga",
				"Flamebacka", "Goldscala", "Highfirea", "Ironhidea", "Jadeclawa",
				"Kettlebreatha", "Lavascale", "Moonflame", "Nightscale", "Opalscale",
				"Pearlwing", "Quickscale", "Rubyscale", "Silverwing", "Topazscale",
				"Underwinga", "Voltclawa", "Wyrmsong", "Xenoscala", "Yellowscale",
			},
		},
	}
}

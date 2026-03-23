package narrativegen

import "forge-rpg/internal/domain"

// narrativeTemplates holds all narrative template blocks indexed by category.
// Each block has tags controlling compatibility:
//   - "any"          → universal, valid for any class and species
//   - class name     → only valid for that specific class (e.g. "wizard")
//   - species name   → only valid for that specific species (e.g. "elf")
//
// RF-03-004: Mínimo 10 templates por categoría (cumplidos con margen).
var narrativeTemplates = map[domain.NarrativeCategory][]domain.NarrativeBlock{

	// ─── BACKGROUND ──────────────────────────────────────────────────────────

	domain.NarrativeBackground: {
		// Universal (any)
		{
			Category: domain.NarrativeBackground,
			Content:  "Criado en las calles de una ciudad portuaria, aprendiste desde joven que la supervivencia depende de la astucia, no de la fuerza bruta.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Perdiste a tu familia durante un invierno brutal. Desde entonces vagás de pueblo en pueblo buscando un lugar al que llamar hogar, aunque aún no lo encontraste.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Fuiste aprendiz de un maestro artesano durante años. Cuando él murió sin dejar herencia, tomaste tus herramientas y te lanzaste al mundo a labrar tu propio destino.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Tu aldea te expulsó por romper una regla ancestral que considerabas injusta. El exilio te enseñó más sobre el mundo que cualquier educación formal.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Creciste escuchando historias de héroes y leyendas. Cuando tuviste edad suficiente, decidiste que era hora de protagonizar tu propia historia.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Trabajaste como mercenario durante años, vendiendo tu espada al mejor postor. Un encargo salió terriblemente mal y te dejó con una deuda que aún intentás saldar.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Sobreviviste a un naufragio que mató a toda tu tripulación. Los meses que pasaste varado en una isla deshabitada te moldearon en cuerpo y espíritu.",
			Tags:     []string{"any"},
		},
		// Barbarian
		{
			Category: domain.NarrativeBackground,
			Content:  "Tu aldea fue arrasada por un dragón cuando eras niño. Desde ese día, el fuego de la venganza arde en tu pecho más fuerte que cualquier llama.",
			Tags:     []string{"barbarian", "fighter", "paladin"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Creciste en las estepas heladas del norte entre guerreros que se curtían en las tormentas. La civilización te parece blanda y ruidosa.",
			Tags:     []string{"barbarian"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Tu tribu te eligió como campeón tras derrotar al jefe anterior en combate singular. El título pesó más de lo esperado y terminaste huyendo de las responsabilidades.",
			Tags:     []string{"barbarian"},
		},
		// Wizard / Artificer (eruditos)
		{
			Category: domain.NarrativeBackground,
			Content:  "Pasaste tu infancia devorando libros en la biblioteca del templo local. Cuando los libros ya no fueron suficientes, buscaste conocimiento donde nadie más se atrevía a mirar.",
			Tags:     []string{"wizard", "artificer", "cleric", "druid"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Eras el prodigio de tu academia mágica hasta que un experimento fallido destruyó parte del ala este. Te exiliaron, pero las ecuaciones todavía danzan en tu mente.",
			Tags:     []string{"wizard", "artificer"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Tu mentor fue un mago excéntrico que vivía en una torre al borde del abismo. Sus enseñanzas eran brillantes, sus métodos cuestionables, y su muerte aún te pesa.",
			Tags:     []string{"wizard"},
		},
		// Rogue / Bard
		{
			Category: domain.NarrativeBackground,
			Content:  "Fuiste el mejor carterista del gremio hasta que robaste la bolsa equivocada. El noble al que le sisaste era más peligroso de lo que parecía.",
			Tags:     []string{"rogue", "bard"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Actuaste en tabernas y mercados desde los siete años. La escena te dio carisma; las calles te dieron sentido común.",
			Tags:     []string{"bard", "rogue"},
		},
		// Cleric / Paladin
		{
			Category: domain.NarrativeBackground,
			Content:  "Desde niño serviste en el templo como acólito. Una noche, la deidad habló directamente a tu corazón y desde entonces sabés que fuiste elegido para algo mayor.",
			Tags:     []string{"cleric", "paladin"},
		},
		// Druid / Ranger
		{
			Category: domain.NarrativeBackground,
			Content:  "Creciste en el bosque profundo, criado por un anciano druida que te enseñó que cada árbol tiene una historia y cada bestia merece respeto.",
			Tags:     []string{"druid", "ranger"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Fuiste el único sobreviviente de una expedición diezmada por criaturas del bosque. Desde entonces, aprendiste su lenguaje para no volver a ser sorprendido.",
			Tags:     []string{"ranger"},
		},
		// Monk
		{
			Category: domain.NarrativeBackground,
			Content:  "Te entregaron al monasterio de niño. Décadas de disciplina te forjaron cuerpo y mente, pero una pregunta sin respuesta te hizo abandonar los muros sagrados.",
			Tags:     []string{"monk"},
		},
		// Sorcerer / Warlock
		{
			Category: domain.NarrativeBackground,
			Content:  "Desde tu primer recuerdo, las chispas mágicas escapaban de tus manos cuando te enojabas. Tu familia temía tus poderes; vos aprendiste a usarlos.",
			Tags:     []string{"sorcerer"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "En el peor momento de tu vida, una entidad de otro plano te ofreció poder a cambio de un servicio aún no definido. Aceptaste sin pensar demasiado.",
			Tags:     []string{"warlock"},
		},
		// Species-specific
		{
			Category: domain.NarrativeBackground,
			Content:  "Viviste tus primeros doscientos años en la ciudad élfil, donde el tiempo pasa lento y las decisiones pesan siglos. Un día decidiste que ya era hora de ver el mundo mortal.",
			Tags:     []string{"elf"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "En las montañas de tu clan enano, el honor es lo más sagrado. Perdiste el tuyo por razones que solo vos conocés, y desde entonces buscás redimirte con cada acción.",
			Tags:     []string{"dwarf"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Como medio elfo, nunca perteneciste del todo a ningún mundo. Eso te hizo observador, adaptable y, en el fondo, profundamente solitario.",
			Tags:     []string{"half-elf"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "Tu sangre infernal siempre te hizo diferente. En tu ciudad natal te miraban con desconfianza; en el camino, aprendiste a usar esa incomodidad a tu favor.",
			Tags:     []string{"tiefling"},
		},
		{
			Category: domain.NarrativeBackground,
			Content:  "El clan Dragonborn al que pertenecés lleva siglos de honor. Vos sos el primero en generaciones que eligió aventurarse solo, rompiendo una tradición colectiva.",
			Tags:     []string{"dragonborn"},
		},
	},

	// ─── MOTIVATION ──────────────────────────────────────────────────────────

	domain.NarrativeMotivation: {
		// Universal (any)
		{
			Category: domain.NarrativeMotivation,
			Content:  "Buscás riqueza suficiente para comprar la libertad de alguien que amás y que aún está atrapado en circunstancias que no pudo elegir.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Querés demostrarles a todos los que dudaron de vos que estaban equivocados. La ambición surgió del dolor, pero el fuego ya es tuyo.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Buscás un artefacto perdido que, según la leyenda, puede devolver la vida. No importa cuánto cueste encontrarlo.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Tu única motivación real es sobrevivir un día más. Todo lo demás —el honor, la gloria, la causa— son palabras bonitas que otros usan para convencerte de arriesgar el cuero.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Creés que hay una amenaza oscura que nadie más puede ver. Cada paso del camino es un intento de reunir pruebas antes de que sea demasiado tarde.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Querés construir algo que perdure: una fortaleza, una institución, un nombre que se recuerde siglos después de que hayas muerto.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "La curiosidad es tu maldición y tu motor. Necesitás saber qué hay del otro lado de cada puerta, de cada horizonte, de cada pregunta sin respuesta.",
			Tags:     []string{"any"},
		},
		// Class-specific
		{
			Category: domain.NarrativeMotivation,
			Content:  "Juraste ante los dioses que vengarías la destrucción de tu hogar. Cada enemigo caído es un paso más hacia ese juramento cumplido.",
			Tags:     []string{"barbarian", "paladin", "fighter"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Querés perfeccionar el arte del combate hasta alcanzar un estado que ningunos escritos describen pero que intuís que existe.",
			Tags:     []string{"fighter", "monk", "barbarian"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Tu deidad te encomendó una misión sagrada. No vas a dormir tranquilo hasta cumplirla, aunque eso cueste todo lo que tenés.",
			Tags:     []string{"cleric", "paladin"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Buscás el equilibrio perdido entre la magia arcana y el mundo natural. Creés que sin ese equilibrio, todo lo que existe está condenado.",
			Tags:     []string{"druid", "ranger"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Querés componer la epopeya definitiva: una obra que capture la verdad del mundo tal como solo vos la podés ver. Las aventuras son tu investigación de campo.",
			Tags:     []string{"bard"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "La deuda con tu patrón del pacto crece. Cada hazaña que realizás es parte del pago de algo que ya no recordás haber prometido.",
			Tags:     []string{"warlock"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Buscás el origen de tu sangre mágica. Alguien en tu linaje hizo algo extraordinario —o terrible— y necesitás saber qué fue.",
			Tags:     []string{"sorcerer"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Cada invención que creás es un paso hacia un prototipo imposible que vive en tus planos desde hace años. Todo lo demás es financiamiento.",
			Tags:     []string{"artificer"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "La meditación te reveló que sos el último guardián de un linaje de guardianes del equilibrio. No elegiste el rol, pero no podés ignorarlo.",
			Tags:     []string{"monk"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Un hechizo que lanzaste sin querer cambió el destino de alguien inocente. Llevás ese peso y buscás una forma de enmendar lo que rompiste.",
			Tags:     []string{"wizard", "sorcerer", "druid"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Te robaron algo que no tiene precio —no dinero, sino algo que define quién sos— y no vas a parar hasta recuperarlo.",
			Tags:     []string{"rogue", "ranger"},
		},
		// Species-specific
		{
			Category: domain.NarrativeMotivation,
			Content:  "Tu longevidad élfil te da perspectiva que los mortales no tienen. Eso te pesa: ves el patrón del error humano repetirse sin fin y querés romperlo.",
			Tags:     []string{"elf"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "El clan enano exige resultados. Volvés cuando tengas un logro digno de ser tallado en la piedra de la sala ancestral.",
			Tags:     []string{"dwarf"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Como halfling, el mundo te subestima siempre. Eso te dio una ventaja que jamás vas a desperdiciar.",
			Tags:     []string{"halfling"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Tu curiosidad gnoma no tiene límites: necesitás entender cómo funciona todo, desarmarlo si es necesario, y armarlo de nuevo pero mejor.",
			Tags:     []string{"gnome"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Querés demostrar que la sangre infernal no define el destino. Cada bien que hacés es una refutación de lo que el mundo esperaba de vos.",
			Tags:     []string{"tiefling"},
		},
		{
			Category: domain.NarrativeMotivation,
			Content:  "Tu herencia orca te dio fuerza; tu herencia humana te dio ambición. Juntas te hacen imparable, si lograís que el mundo te dé una oportunidad.",
			Tags:     []string{"half-orc"},
		},
	},

	// ─── SECRET ──────────────────────────────────────────────────────────────

	domain.NarrativeSecret: {
		// Universal (any)
		{
			Category: domain.NarrativeSecret,
			Content:  "En realidad nunca tuviste el título que todos creen que tenés. La credencial era falsa, pero las habilidades son completamente reales.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Fuiste responsable de la muerte de alguien. No fue tu intención, pero la consecuencia fue real y nunca lo confesaste.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Tenés una deuda enorme con una organización peligrosa. Te persiguen, aunque todavía no lo saben los que te rodean.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "En tu pasado hay una traición que jamás podrías justificar. Si los que te conocen hoy lo supieran, todo cambiaría.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Tenés una familia en algún lugar del mundo que cree que estás muerto. Fue más fácil no corregirles.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Una profecía antigua menciona alguien con tu descripción exacta. No sabés si eso es bueno o no, así que no le contás a nadie.",
			Tags:     []string{"any"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "En cierto momento de tu vida, tomaste una decisión cobarde de la que nunca hablaste. El orgullo que mostrás es una armadura contra esa verdad.",
			Tags:     []string{"any"},
		},
		// Class-specific
		{
			Category: domain.NarrativeSecret,
			Content:  "Tu furia no es solo rabia: es la manifestación de un espíritu ancestral que habita en vos desde que eras niño. No sabés si es un don o una maldición.",
			Tags:     []string{"barbarian"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Matas a gente que otros te pagan por eliminar. No son malos trabajos, solo son los más lucrativos. Tus compañeros creen que sos un héroe.",
			Tags:     []string{"rogue", "fighter"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Tu fe tambalea. La deidad a la que servís lleva meses en silencio y empezás a preguntarte si alguna vez estuvo realmente ahí.",
			Tags:     []string{"cleric", "paladin"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Descubriste algo en un tomo prohibido que cambia todo lo que creías saber sobre la magia. Si lo dijeras en voz alta, probablemente te ejecutarían.",
			Tags:     []string{"wizard", "artificer"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Tu pacto tiene una cláusula que nunca leyiste con atención. Cuando el patrón la active, no va a ser agradable.",
			Tags:     []string{"warlock"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Una de las canciones que tocás es un encantamiento real. Sabés exactamente qué hace en la mente de quien la escucha.",
			Tags:     []string{"bard"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "El equilibrio natural que decís proteger fue perturbado por vos hace años. El bosque lo recuerda aunque vos intentés olvidarlo.",
			Tags:     []string{"druid", "ranger"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Tus poderes mágicos innatos tienen un costo físico que ocultás cuidadosamente. Cada hechizo poderoso acorta algo que preferís no medir.",
			Tags:     []string{"sorcerer"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Abandonaste el monasterio antes de completar el rito final. Técnicamente, no sos lo que todos creen que sos.",
			Tags:     []string{"monk"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Uno de tus constructos tiene conciencia propia. Lo sabés. Ninguno de los dos habló de eso todavía.",
			Tags:     []string{"artificer"},
		},
		// Species-specific
		{
			Category: domain.NarrativeSecret,
			Content:  "Tu memoria élfil guarda algo que presenciaste hace siglos y que nadie más vivo recuerda. Esa información podría cambiar el equilibrio de poder en el mundo.",
			Tags:     []string{"elf"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "El clan enano al que pertenecés ocultó algo durante generaciones. Vos descubriste qué era, y la respuesta te dejó con más preguntas que antes.",
			Tags:     []string{"dwarf"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Tu herencia de medio humano y medio elfo viene de una unión que ninguno de los dos lados aprobó. Hay personas que todavía buscan borrar esa historia.",
			Tags:     []string{"half-elf"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Llevas una marca de tu sangre orca que, si alguien la reconociera, revelaría un linaje que preferís mantener oculto.",
			Tags:     []string{"half-orc"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "El nombre que usás no es el tuyo. El verdadero nombre tiefling que te dieron al nacer tiene poder, y alguien podría usarlo en tu contra.",
			Tags:     []string{"tiefling"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "El color de tus escamas no coincide con el clan que decís representar. Hay una historia familiar que enterraste tan profundo que casi la olvidaste.",
			Tags:     []string{"dragonborn"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Sos más valiente de lo que parecés, y eso te aterra. Porque si la gente lo descubriera, empezarían a pedirte cosas que no sabés si podés dar.",
			Tags:     []string{"halfling"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Tu invento gnomo más famoso tiene un defecto que nunca corregiste. Hasta ahora nadie resultó herido, pero es solo cuestión de tiempo.",
			Tags:     []string{"gnome"},
		},
		{
			Category: domain.NarrativeSecret,
			Content:  "Nunca llegaste a ese lugar del que tanto hablás. Pero la historia que inventaste es tan buena que casi vos mismo la creés.",
			Tags:     []string{"human"},
		},
	},
}

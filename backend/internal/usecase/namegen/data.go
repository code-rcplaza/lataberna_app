package namegen

type namePool struct {
	male   []string
	female []string
}

var nameData = map[string]namePool{
	// ── Humans ───────────────────────────────────────────────
	"human": {
		male: []string{
			"Aldric", "Brennan", "Cael", "Dorian", "Edric",
			"Faolan", "Gareth", "Hadwin", "Isidor", "Jareth",
			"Kiran", "Leoric", "Maddox", "Nolan", "Orwin",
			"Phelan", "Quinn", "Roderick", "Soren", "Theron",
			"Ulric", "Vance", "Wulfric", "Xander", "Yorick",
		},
		female: []string{
			"Aelara", "Brenna", "Calla", "Dara", "Edlyn",
			"Fiona", "Gwenna", "Hilda", "Isara", "Jenna",
			"Kira", "Lyra", "Mara", "Nora", "Orla",
			"Priya", "Quinn", "Riona", "Sable", "Thea",
			"Ursa", "Vera", "Wren", "Xara", "Yara",
		},
	},

	// ── Elves ─────────────────────────────────────────────────
	"high-elf": {
		male: []string{
			"Aelindor", "Caladrel", "Erevan", "Faenor", "Galinndan",
			"Haladavar", "Immeral", "Jelenneth", "Keyleth", "Laucian",
			"Mindartis", "Naeris", "Orym", "Paelias", "Quarion",
			"Riardon", "Soveliss", "Thamior", "Uvaleth", "Valenor",
			"Whilom", "Xanaphia", "Yvrel", "Zannin", "Aranthor",
		},
		female: []string{
			"Adrie", "Birel", "Caelynn", "Dara", "Enialis",
			"Faral", "Gennal", "Halueth", "Irann", "Jelenneth",
			"Keyleth", "Leshanna", "Mialee", "Naivara", "Oparal",
			"Quelenna", "Rania", "Sariel", "Thia", "Urrel",
			"Valanthe", "Wasitara", "Xanaphia", "Yalanue", "Zylvara",
		},
	},
	"wood-elf": {
		male: []string{
			"Adran", "Aelar", "Beiro", "Carric", "Dayereth",
			"Enna", "Galinndan", "Hadarai", "Ivellios", "Laucian",
			"Mindartis", "Naeris", "Paelias", "Quarion", "Riardon",
			"Soveliss", "Thamior", "Theren", "Valenor", "Varis",
			"Zannin", "Aravel", "Brysis", "Celadyr", "Delmair",
		},
		female: []string{
			"Adrie", "Althaea", "Anastrianna", "Andraste", "Antinua",
			"Bethrynna", "Birel", "Caelynn", "Drusilia", "Enna",
			"Felosial", "Ielenia", "Jelenneth", "Keyleth", "Leshanna",
			"Mialee", "Naivara", "Quelenna", "Sariel", "Shanairla",
			"Shava", "Silaqui", "Theirastra", "Valna", "Xanaphia",
		},
	},
	"drow": {
		male: []string{
			"Drizzt", "Zaknafein", "Jarlaxle", "Pharaun", "Ryld",
			"Valas", "Raashub", "Kelnozz", "Uthegental", "Malagdorl",
			"Vorn", "Szordrin", "Bregan", "Nimor", "Kalannar",
			"Drisinil", "Liriel", "Zerin", "Arach", "Ilphrin",
			"Balok", "Chaulssin", "Devir", "Erelal", "Faerath",
		},
		female: []string{
			"Liriel", "Malice", "Vierna", "Briza", "Maya",
			"Quenthel", "Triel", "SosUmptu", "Shinayne", "Vendes",
			"Zeerith", "Akordia", "Belrae", "Cylva", "Danifae",
			"Eliztrae", "Filraen", "Greyanna", "Halisstra", "Imrae",
			"Jhannyl", "Kira", "Laele", "Myryl", "Nedylene",
		},
	},

	// ── Dwarves ───────────────────────────────────────────────
	"hill-dwarf": {
		male: []string{
			"Adrik", "Alberich", "Baern", "Barendd", "Brottor",
			"Bruenor", "Dain", "Darrak", "Delg", "Eberk",
			"Einkil", "Fargrim", "Flint", "Gardain", "Harbek",
			"Kildrak", "Morgran", "Orsik", "Oskar", "Rangrim",
			"Rurik", "Taklinn", "Thoradin", "Thorin", "Tordek",
		},
		female: []string{
			"Amber", "Artin", "Audhild", "Bardryn", "Dagnal",
			"Diesa", "Eldeth", "Falkrunn", "Finellen", "Gunnloda",
			"Gurdis", "Helja", "Hlin", "Kathra", "Kristryd",
			"Ilde", "Liftrasa", "Mardred", "Riswynn", "Sannl",
			"Torbera", "Torgga", "Vistra", "Borgna", "Helma",
		},
	},
	"mountain-dwarf": {
		male: []string{
			"Aldric", "Borin", "Cragmar", "Dundrak", "Edric",
			"Forgrim", "Grondar", "Hagrim", "Ironhold", "Jarek",
			"Keldrak", "Lothrak", "Morigrim", "Nordak", "Orkrak",
			"Peldar", "Quorak", "Rokdar", "Stormak", "Teldrak",
			"Uldrak", "Vordak", "Wargrim", "Xendrak", "Yeldrak",
		},
		female: []string{
			"Aldis", "Borgna", "Coldara", "Durnea", "Elfrida",
			"Fangora", "Goltara", "Hilma", "Ingrid", "Jorna",
			"Koldra", "Lofna", "Moltara", "Norgra", "Olda",
			"Poldra", "Ragna", "Solgra", "Toldra", "Uldra",
			"Valdra", "Wolgra", "Xoldra", "Yoldra", "Zoldra",
		},
	},

	// ── Halflings ─────────────────────────────────────────────
	"lightfoot": {
		male: []string{
			"Alton", "Ander", "Cade", "Corrin", "Eldon",
			"Errich", "Finnan", "Garret", "Lindal", "Lyle",
			"Merric", "Milo", "Osborn", "Perrin", "Reed",
			"Roscoe", "Wellby", "Beau", "Cob", "Davin",
			"Fenrick", "Gable", "Hob", "Jasper", "Kender",
		},
		female: []string{
			"Andry", "Bree", "Callie", "Cora", "Euphemia",
			"Jillian", "Kithri", "Lavinia", "Lidda", "Merla",
			"Nedda", "Paela", "Portia", "Seraphina", "Shaena",
			"Trym", "Vani", "Verna", "Amaryllis", "Birdie",
			"Celandine", "Dora", "Eglantine", "Florimel", "Goldie",
		},
	},
	"stout": {
		male: []string{
			"Baldric", "Bramble", "Brock", "Burr", "Cob",
			"Dag", "Dodger", "Durbin", "Fenn", "Garric",
			"Griff", "Hardy", "Holt", "Knob", "Lob",
			"Mack", "Nob", "Pip", "Rob", "Sam",
			"Stubb", "Tad", "Tom", "Wil", "Zob",
		},
		female: []string{
			"Ally", "Bertha", "Blossom", "Bunny", "Daisy",
			"Della", "Dot", "Ember", "Flora", "Greta",
			"Hana", "Iris", "Jade", "Kitty", "Lily",
			"May", "Midge", "Nell", "Opal", "Pearl",
			"Poppy", "Rose", "Ruby", "Sage", "Violet",
		},
	},

	// ── Gnomes ───────────────────────────────────────────────
	"forest-gnome": {
		male: []string{
			"Alston", "Alvyn", "Boddynock", "Brocc", "Burgell",
			"Dimble", "Eldon", "Erky", "Fonkin", "Frug",
			"Gerbo", "Gimble", "Glim", "Jebeddo", "Kellen",
			"Namfoodle", "Orryn", "Roondar", "Seebo", "Sindri",
			"Warryn", "Wrenn", "Zook", "Bink", "Dabble",
		},
		female: []string{
			"Bimpnottin", "Breena", "Caramip", "Carlin", "Donella",
			"Duvamil", "Ella", "Ellyjobell", "Ellywick", "Lilli",
			"Loopmottin", "Lorilla", "Mardnab", "Nissa", "Nyx",
			"Oda", "Orla", "Roywyn", "Shamil", "Tana",
			"Waywocket", "Zanna", "Bree", "Calli", "Dotti",
		},
	},
	"rock-gnome": {
		male: []string{
			"Abzug", "Alberich", "Binkadink", "Cogsworth", "Dabbledob",
			"Dinkum", "Fiddlesticks", "Fizzbang", "Gadget", "Geargrind",
			"Gnarlick", "Grumbly", "Junkbolt", "Klix", "Mekka",
			"Nix", "Plink", "Ratchet", "Sprocket", "Tick",
			"Tinker", "Tock", "Volt", "Whirr", "Zap",
		},
		female: []string{
			"Binky", "Clank", "Clink", "Cognia", "Dazzle",
			"Flix", "Glimmer", "Gizmo", "Jink", "Klink",
			"Minka", "Nixie", "Plink", "Quirk", "Rinka",
			"Spark", "Sprix", "Tick", "Tinka", "Trinket",
			"Twix", "Vix", "Wink", "Xink", "Zink",
		},
	},

	// ── Half-Elf ──────────────────────────────────────────────
	"half-elf": {
		male: []string{
			"Aelric", "Bran", "Caelum", "Dorn", "Erevan",
			"Faen", "Galen", "Hadwin", "Isidor", "Jareth",
			"Kael", "Lorn", "Maren", "Naeris", "Oryn",
			"Paelias", "Quarion", "Rhogar", "Soren", "Theron",
			"Ulric", "Valenor", "Wulf", "Xander", "Yorin",
		},
		female: []string{
			"Aelara", "Brenna", "Caelynn", "Dara", "Enialis",
			"Faral", "Gennal", "Halueth", "Irann", "Jenna",
			"Kira", "Leshanna", "Mialee", "Naivara", "Orla",
			"Quella", "Rania", "Sariel", "Thia", "Urrel",
			"Valanthe", "Wren", "Xara", "Yara", "Zara",
		},
	},

	// ── Half-Orc ──────────────────────────────────────────────
	"half-orc": {
		male: []string{
			"Dench", "Feng", "Gell", "Henk", "Holg",
			"Imsh", "Keth", "Krusk", "Mhurren", "Ront",
			"Shump", "Thokk", "Vrag", "Wund", "Zorr",
			"Brug", "Crag", "Druk", "Frug", "Grak",
			"Hrak", "Jrak", "Krak", "Mrak", "Prak",
		},
		female: []string{
			"Baggi", "Emen", "Engong", "Kansif", "Myev",
			"Neega", "Ovak", "Ownka", "Shautha", "Sutha",
			"Vola", "Volen", "Yevelda", "Zasha", "Grak",
			"Hrak", "Jrak", "Krak", "Mrak", "Nrak",
			"Prak", "Rrak", "Srak", "Trak", "Urak",
		},
	},

	// ── Tiefling ──────────────────────────────────────────────
	"tiefling": {
		male: []string{
			"Akmenos", "Amnon", "Barakas", "Damakos", "Ekemon",
			"Iados", "Kairon", "Leucis", "Melech", "Mordai",
			"Morthos", "Pelaios", "Skamos", "Therai", "Zed",
			"Art", "Carrion", "Chant", "Creed", "Despair",
			"Excellence", "Fear", "Glory", "Hope", "Ideal",
		},
		female: []string{
			"Akta", "Annis", "Bryseis", "Criella", "Damaia",
			"Ea", "Kallista", "Lerissa", "Makaria", "Nemeia",
			"Orianna", "Phelaia", "Rieta", "Tanika", "Zelica",
			"Torment", "Wander", "Whimsy", "Zeal", "Vengeance",
			"Sorrow", "Ruin", "Pity", "Misery", "Love",
		},
	},

	// ── Dragonborn ────────────────────────────────────────────
	"dragonborn": {
		male: []string{
			"Arjhan", "Balasar", "Bharash", "Donaar", "Ghesh",
			"Heskan", "Kriv", "Medrash", "Mehen", "Nadarr",
			"Pandjed", "Patrin", "Rhogar", "Shamash", "Shedinn",
			"Tarhun", "Torinn", "Vishap", "Vorel", "Zedaar",
			"Blazewing", "Copperclaw", "Drakarr", "Emberfang", "Goldwing",
		},
		female: []string{
			"Akra", "Biri", "Daar", "Farideh", "Harann",
			"Havilar", "Jheri", "Kava", "Korinn", "Mishann",
			"Nala", "Perra", "Raiann", "Sora", "Surina",
			"Thava", "Uadjit", "Vroth", "Yenna", "Zara",
			"Amberscale", "Brightfang", "Crystalwing", "Dawnfire", "Emberclaw",
		},
	},
}

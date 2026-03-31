package statblock

import (
	"fmt"

	"forge-rpg/internal/domain"
)

// BackgroundEntry represents a D&D 5.5e background with its mechanical properties.
type BackgroundEntry struct {
	Name        string    // display name in Spanish (e.g., "Acólito")
	ASIPool     [3]string // the three stats the background can boost, e.g., ["WIS", "INT", "CHA"]
	OriginFeat  string    // fixed feat granted at level 1
	Description string    // short description of the origin feat benefit
	Tags        []string  // class/species coherence tags; "any" = universal
}

// classData holds the static class configuration for stat generation.
type classData struct {
	hitDie        int
	primaryStat   string // always receives 15 — the class's main stat
	secondaryStat string // always receives 14 — the class's second priority
	armorKey      string
}

// maxDex2 is a pointer to 2, used for medium armor MaxDex cap.
var maxDex2 = func() *int { v := 2; return &v }()

// classTable holds the class priorities and armor for all 13 classes.
// primaryStat receives 15, secondaryStat receives 14; remaining stats get
// [13,12,10,8] shuffled randomly — always totalling 27 in point buy cost.
var classTable = map[domain.Class]classData{
	domain.ClassBarbarian: {hitDie: 12, primaryStat: "STR", secondaryStat: "CON", armorKey: "unarmored-barbarian"},
	domain.ClassBard:      {hitDie: 8,  primaryStat: "CHA", secondaryStat: "DEX", armorKey: "leather"},
	domain.ClassCleric:    {hitDie: 8,  primaryStat: "WIS", secondaryStat: "STR", armorKey: "chain-shirt"},
	domain.ClassDruid:     {hitDie: 8,  primaryStat: "WIS", secondaryStat: "CON", armorKey: "chain-shirt"},
	domain.ClassFighter:   {hitDie: 10, primaryStat: "STR", secondaryStat: "CON", armorKey: "chain-mail"},
	domain.ClassMonk:      {hitDie: 8,  primaryStat: "DEX", secondaryStat: "WIS", armorKey: "unarmored-monk"},
	domain.ClassPaladin:   {hitDie: 10, primaryStat: "STR", secondaryStat: "CHA", armorKey: "chain-mail"},
	domain.ClassRanger:    {hitDie: 10, primaryStat: "DEX", secondaryStat: "WIS", armorKey: "chain-shirt"},
	domain.ClassRogue:     {hitDie: 8,  primaryStat: "DEX", secondaryStat: "INT", armorKey: "leather"},
	domain.ClassSorcerer:  {hitDie: 6,  primaryStat: "CHA", secondaryStat: "CON", armorKey: "clothes"},
	domain.ClassWarlock:   {hitDie: 8,  primaryStat: "CHA", secondaryStat: "DEX", armorKey: "leather"},
	domain.ClassWizard:    {hitDie: 6,  primaryStat: "INT", secondaryStat: "CON", armorKey: "clothes"},
	domain.ClassArtificer: {hitDie: 8,  primaryStat: "INT", secondaryStat: "CON", armorKey: "chain-shirt"},
}

// armorTable holds the ArmorType instances keyed by armor key.
// Values from mvp-rules.context.md.
var armorTable = map[string]domain.ArmorType{
	// Unarmored Defense — Barbarian: 10 + DEX + CON (handled in calculateAC)
	"unarmored-barbarian": {
		Name:     "Unarmored Defense (Barbarian)",
		Category: domain.ArmorCategory("unarmored-barbarian"),
		BaseAC:   10,
		MaxDex:   nil,
	},
	// Unarmored Defense — Monk: 10 + DEX + WIS (handled in calculateAC)
	"unarmored-monk": {
		Name:     "Unarmored Defense (Monk)",
		Category: domain.ArmorCategory("unarmored-monk"),
		BaseAC:   10,
		MaxDex:   nil,
	},
	// Leather armor: light, baseAC=11, full DEX
	"leather": {
		Name:     "Leather",
		Category: domain.ArmorLight,
		BaseAC:   11,
		MaxDex:   nil,
	},
	// Chain shirt: medium, baseAC=13, maxDex=2
	"chain-shirt": {
		Name:     "Chain Shirt",
		Category: domain.ArmorMedium,
		BaseAC:   13,
		MaxDex:   maxDex2,
	},
	// Chain mail: heavy, baseAC=16, no DEX
	"chain-mail": {
		Name:     "Chain Mail",
		Category: domain.ArmorHeavy,
		BaseAC:   16,
		MaxDex:   nil,
	},
	// Clothes (robes): none, baseAC=10, full DEX
	"clothes": {
		Name:     "Clothes",
		Category: domain.ArmorNone,
		BaseAC:   10,
		MaxDex:   nil,
	},
}

// backgroundTable holds the D&D 5.5e (2024 PHB) backgrounds used for ASI and feat resolution.
// Order is stable — index 0..15 maps to the 16 canonical backgrounds.
// Names and origin feats are in Spanish, matching backgrounds.md.
var backgroundTable = []BackgroundEntry{
	{Name: "Acólito",    ASIPool: [3]string{"WIS", "INT", "CHA"}, OriginFeat: "Iniciado en la magia (Clérigo)", Description: "Obtienes 2 trucos y un conjuro de nivel 1 de la lista de Clérigo.",                                                            Tags: []string{"cleric", "paladin", "druid"}},
	{Name: "Artesano",   ASIPool: [3]string{"STR", "DEX", "INT"}, OriginFeat: "Fabricante",                     Description: "Competencia con 3 herramientas, 20% de descuento en compras y fabricas objetos rápido.",                                    Tags: []string{"any"}},
	{Name: "Charlatán",  ASIPool: [3]string{"DEX", "CON", "CHA"}, OriginFeat: "Afortunado",                     Description: "Tienes puntos de suerte (igual a tu BC) para ganar ventaja o dar desventaja a atacantes.",                                  Tags: []string{"rogue", "bard", "warlock"}},
	{Name: "Criminal",   ASIPool: [3]string{"DEX", "CON", "INT"}, OriginFeat: "Alerta",                         Description: "Sumas tu BC a la iniciativa y puedes intercambiar tu iniciativa con un aliado.",                                            Tags: []string{"rogue", "ranger", "warlock"}},
	{Name: "Animador",   ASIPool: [3]string{"STR", "DEX", "CHA"}, OriginFeat: "Músico",                         Description: "Competencia con 3 instrumentos y das Inspiración Heroica a aliados tras un descanso.",                                      Tags: []string{"bard", "rogue", "sorcerer"}},
	{Name: "Campesino",  ASIPool: [3]string{"STR", "CON", "WIS"}, OriginFeat: "Duro",                           Description: "Tus puntos de golpe máximos aumentan en 2 por cada nivel que tengas.",                                                      Tags: []string{"any"}},
	{Name: "Guardia",    ASIPool: [3]string{"STR", "INT", "CHA"}, OriginFeat: "Alerta",                         Description: "Sumas tu BC a la iniciativa y puedes intercambiar tu iniciativa con un aliado.",                                            Tags: []string{"fighter", "paladin", "ranger"}},
	{Name: "Guía",       ASIPool: [3]string{"DEX", "CON", "WIS"}, OriginFeat: "Iniciado en la magia (Druida)",  Description: "Obtienes 2 trucos y un conjuro de nivel 1 de la lista de Druida.",                                                          Tags: []string{"ranger", "druid", "monk"}},
	{Name: "Ermitaño",   ASIPool: [3]string{"CON", "WIS", "CHA"}, OriginFeat: "Sanador",                        Description: "Puedes usar útiles de sanador para que un aliado gaste un dado de golpe y se cure; repites \"1\" en curaciones.",            Tags: []string{"druid", "monk", "cleric"}},
	{Name: "Comerciante", ASIPool: [3]string{"CON", "INT", "CHA"}, OriginFeat: "Afortunado",                   Description: "Tienes puntos de suerte (igual a tu BC) para ganar ventaja o dar desventaja a atacantes.",                                   Tags: []string{"any"}},
	{Name: "Noble",      ASIPool: [3]string{"STR", "INT", "CHA"}, OriginFeat: "Músico",                         Description: "Competencia con 3 instrumentos y das Inspiración Heroica a aliados tras un descanso.",                                      Tags: []string{"paladin", "fighter", "bard", "warlock"}},
	{Name: "Erudito",    ASIPool: [3]string{"CON", "INT", "WIS"}, OriginFeat: "Iniciado en la magia (Mago)",    Description: "Obtienes 2 trucos y un conjuro de nivel 1 de la lista de Mago.",                                                            Tags: []string{"wizard", "sorcerer", "artificer"}},
	{Name: "Marinero",   ASIPool: [3]string{"STR", "DEX", "WIS"}, OriginFeat: "Tabernero Peleón",               Description: "Tus ataques desarmados son 1d4, repites \"1\" en su daño y puedes empujar 1.5m al golpear.",                                Tags: []string{"any"}},
	{Name: "Escriba",    ASIPool: [3]string{"DEX", "INT", "WIS"}, OriginFeat: "Iniciado en la magia (Mago)",    Description: "Obtienes 2 trucos y un conjuro de nivel 1 de la lista de Mago.",                                                            Tags: []string{"wizard", "artificer", "cleric"}},
	{Name: "Soldado",    ASIPool: [3]string{"STR", "DEX", "CON"}, OriginFeat: "Atacante Salvaje",               Description: "Una vez por turno, al acertar con un arma, tiras dos veces el daño y eliges el mayor.",                                    Tags: []string{"fighter", "paladin", "ranger", "barbarian"}},
	{Name: "Vagabundo",  ASIPool: [3]string{"DEX", "WIS", "CHA"}, OriginFeat: "Iniciado en la magia (Druida)",  Description: "Obtienes 2 trucos y un conjuro de nivel 1 de la lista de Druida.",                                                          Tags: []string{"rogue", "ranger", "bard", "monk"}},
}

// BackgroundsForClass returns all backgrounds whose tags include the given class or "any".
func BackgroundsForClass(class string) []BackgroundEntry {
	var result []BackgroundEntry
	for _, b := range backgroundTable {
		for _, tag := range b.Tags {
			if tag == "any" || tag == class {
				result = append(result, b)
				break
			}
		}
	}
	return result
}

// AllBackgrounds returns all background entries.
func AllBackgrounds() []BackgroundEntry {
	return backgroundTable
}

// OriginFeatFor returns the origin feat for the given background name.
// Returns an empty string if the background is not found.
func OriginFeatFor(backgroundName string) string {
	for _, b := range backgroundTable {
		if b.Name == backgroundName {
			return b.OriginFeat
		}
	}
	return ""
}

// BackgroundDescriptionFor returns the origin feat description for the given background name.
// Returns an empty string if the background is not found.
func BackgroundDescriptionFor(name string) string {
	for _, b := range backgroundTable {
		if b.Name == name {
			return b.Description
		}
	}
	return ""
}

// allClasses returns a stable slice of all valid classes for random selection.
func allClasses() []domain.Class {
	return []domain.Class{
		domain.ClassBarbarian,
		domain.ClassBard,
		domain.ClassCleric,
		domain.ClassDruid,
		domain.ClassFighter,
		domain.ClassMonk,
		domain.ClassPaladin,
		domain.ClassRanger,
		domain.ClassRogue,
		domain.ClassSorcerer,
		domain.ClassWarlock,
		domain.ClassWizard,
		domain.ClassArtificer,
	}
}

// allSpecies returns a stable slice of all valid species for random selection.
func allSpecies() []domain.Species {
	return []domain.Species{
		domain.SpeciesHuman,
		domain.SpeciesElf,
		domain.SpeciesDwarf,
		domain.SpeciesHalfling,
		domain.SpeciesGnome,
		domain.SpeciesHalfElf,
		domain.SpeciesHalfOrc,
		domain.SpeciesTiefling,
		domain.SpeciesDragonborn,
	}
}

// subSpeciesFor returns the valid subspecies for a given species.
// Species without subspecies return an empty slice.
func subSpeciesFor(s domain.Species) []domain.SubSpecies {
	switch s {
	case domain.SpeciesElf:
		return []domain.SubSpecies{
			domain.SubSpeciesHighElf,
			domain.SubSpeciesWoodElf,
			domain.SubSpeciesDrow,
		}
	case domain.SpeciesDwarf:
		return []domain.SubSpecies{
			domain.SubSpeciesHillDwarf,
			domain.SubSpeciesMountainDwarf,
		}
	case domain.SpeciesHalfling:
		return []domain.SubSpecies{
			domain.SubSpeciesLightfoot,
			domain.SubSpeciesStout,
		}
	case domain.SpeciesGnome:
		return []domain.SubSpecies{
			domain.SubSpeciesForestGnome,
			domain.SubSpeciesRockGnome,
		}
	default:
		return nil
	}
}

// errUnknownClass returns a formatted error for an unknown class.
func errUnknownClass(c domain.Class) error {
	return fmt.Errorf("statblock: unknown class %q", c)
}

// errUnknownArmor returns a formatted error for an unknown armor key.
func errUnknownArmor(key string) error {
	return fmt.Errorf("statblock: unknown armor key %q", key)
}

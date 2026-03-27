package statblock

import (
	"fmt"

	"forge-rpg/internal/domain"
)

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

// speciesBonuses holds the D&D 5e PHB ability score bonuses per species/subspecies.
// Keyed by species bonus key (see speciesBonusKey).
//
// Species bonuses (5e PHB):
//   Human:          +1 to all 6 stats
//   High Elf:       +2 DEX, +1 INT
//   Wood Elf:       +2 DEX, +1 WIS
//   Drow:           +2 DEX, +1 CHA
//   Hill Dwarf:     +2 CON, +1 WIS
//   Mountain Dwarf: +2 CON, +2 STR
//   Lightfoot:      +2 DEX, +1 CHA
//   Stout:          +2 DEX, +1 CON
//   Forest Gnome:   +2 INT, +1 DEX
//   Rock Gnome:     +2 INT, +1 CON
//   Half-Elf:       +2 CHA, +1 STR, +1 CON (choose 2 — defaulting to STR, CON per spec)
//   Half-Orc:       +2 STR, +1 CON
//   Tiefling:       +2 CHA, +1 INT
//   Dragonborn:     +2 STR, +1 CHA
var speciesBonuses = map[string][]domain.AbilityBonus{
	"human": {
		{Stat: "STR", Value: 1, Source: "species"},
		{Stat: "DEX", Value: 1, Source: "species"},
		{Stat: "CON", Value: 1, Source: "species"},
		{Stat: "INT", Value: 1, Source: "species"},
		{Stat: "WIS", Value: 1, Source: "species"},
		{Stat: "CHA", Value: 1, Source: "species"},
	},
	"high-elf": {
		{Stat: "DEX", Value: 2, Source: "species"},
		{Stat: "INT", Value: 1, Source: "species"},
	},
	"wood-elf": {
		{Stat: "DEX", Value: 2, Source: "species"},
		{Stat: "WIS", Value: 1, Source: "species"},
	},
	"drow": {
		{Stat: "DEX", Value: 2, Source: "species"},
		{Stat: "CHA", Value: 1, Source: "species"},
	},
	"hill-dwarf": {
		{Stat: "CON", Value: 2, Source: "species"},
		{Stat: "WIS", Value: 1, Source: "species"},
	},
	"mountain-dwarf": {
		{Stat: "STR", Value: 2, Source: "species"},
		{Stat: "CON", Value: 2, Source: "species"},
	},
	"lightfoot": {
		{Stat: "DEX", Value: 2, Source: "species"},
		{Stat: "CHA", Value: 1, Source: "species"},
	},
	"stout": {
		{Stat: "DEX", Value: 2, Source: "species"},
		{Stat: "CON", Value: 1, Source: "species"},
	},
	"forest-gnome": {
		{Stat: "INT", Value: 2, Source: "species"},
		{Stat: "DEX", Value: 1, Source: "species"},
	},
	"rock-gnome": {
		{Stat: "INT", Value: 2, Source: "species"},
		{Stat: "CON", Value: 1, Source: "species"},
	},
	// Half-Elf: +2 CHA, then player chooses 2 stats for +1 each.
	// For generation purposes: defaulting to STR and CON per spec instructions.
	"half-elf": {
		{Stat: "CHA", Value: 2, Source: "species"},
		{Stat: "STR", Value: 1, Source: "species"},
		{Stat: "CON", Value: 1, Source: "species"},
	},
	"half-orc": {
		{Stat: "STR", Value: 2, Source: "species"},
		{Stat: "CON", Value: 1, Source: "species"},
	},
	"tiefling": {
		{Stat: "CHA", Value: 2, Source: "species"},
		{Stat: "INT", Value: 1, Source: "species"},
	},
	"dragonborn": {
		{Stat: "STR", Value: 2, Source: "species"},
		{Stat: "CHA", Value: 1, Source: "species"},
	},
}

// speciesBonusKey returns the map key to use for species bonus lookup.
// For species with subspecies, the subspecies key is used directly.
func speciesBonusKey(s domain.Species, sub *domain.SubSpecies) string {
	if sub != nil {
		return string(*sub)
	}
	return string(s)
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

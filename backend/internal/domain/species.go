package domain

// Species represents a playable species in D&D 5e.
type Species string

const (
	SpeciesHuman      Species = "human"
	SpeciesElf        Species = "elf"
	SpeciesDwarf      Species = "dwarf"
	SpeciesHalfling   Species = "halfling"
	SpeciesGnome      Species = "gnome"
	SpeciesHalfElf    Species = "half-elf"
	SpeciesHalfOrc    Species = "half-orc"
	SpeciesTiefling   Species = "tiefling"
	SpeciesDragonborn Species = "dragonborn"
)

// SubSpecies represents a subspecies variant.
type SubSpecies string

const (
	// Elf subspecies
	SubSpeciesHighElf SubSpecies = "high-elf"
	SubSpeciesWoodElf SubSpecies = "wood-elf"
	SubSpeciesDrow    SubSpecies = "drow"

	// Dwarf subspecies
	SubSpeciesHillDwarf     SubSpecies = "hill-dwarf"
	SubSpeciesMountainDwarf SubSpecies = "mountain-dwarf"

	// Halfling subspecies
	SubSpeciesLightfoot SubSpecies = "lightfoot"
	SubSpeciesStout     SubSpecies = "stout"

	// Gnome subspecies
	SubSpeciesForestGnome SubSpecies = "forest-gnome"
	SubSpeciesRockGnome   SubSpecies = "rock-gnome"

	// Tiefling subspecies (naming convention)
	SubSpeciesInfernalTiefling SubSpecies = "tiefling-infernal"
	SubSpeciesVirtueTiefling   SubSpecies = "tiefling-virtue"
)

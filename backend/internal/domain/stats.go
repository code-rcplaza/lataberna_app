package domain

// Stats represents the six D&D ability scores.
type Stats struct {
	STR int
	DEX int
	CON int
	INT int
	WIS int
	CHA int
}

// Modifiers holds the calculated modifier for each ability score.
// Modifiers MUST always be calculated from FinalStats, never BaseStats.
type Modifiers struct {
	STR int
	DEX int
	CON int
	INT int
	WIS int
	CHA int
}

// DerivedStats holds HP and AC, both derived from FinalStats and class/armor.
type DerivedStats struct {
	HP int
	AC int
}

// AbilityBonus represents a bonus to a specific ability score.
type AbilityBonus struct {
	Stat   string // "STR" | "DEX" | "CON" | "INT" | "WIS" | "CHA"
	Value  int
	Source string // "species" | "background"
}

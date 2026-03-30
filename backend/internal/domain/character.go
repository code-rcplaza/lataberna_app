package domain

import "time"

// CharacterLocks tracks which fields are locked during partial regeneration.
// A locked field is preserved; an unlocked field is regenerated.
type CharacterLocks struct {
	Name       bool
	Stats      bool
	Background bool
	Motivation bool
	Secret     bool
}

// Character is the central aggregate of the domain.
// All generation parameters are optional — omitted fields generate randomly.
type Character struct {
	ID     string
	UserID string // owner — set by the persistence layer, never by generators

	// Identity
	Name       string
	Species    Species
	SubSpecies *SubSpecies
	Class      Class
	Level      int

	// Rule configuration
	Ruleset            Ruleset
	AbilityBonusSource AbilityBonusSource
	// BackgroundType is the D&D 5.5e mechanical background entity (e.g., "acolyte", "soldier").
	// Distinct from the narrative Background NarrativeBlock.
	BackgroundType string
	// ASIDistribution defines how the background ASI pool is allocated.
	// "standard" = +2/+1 to two abilities; "spread" = +1/+1/+1 to three abilities.
	ASIDistribution string

	// Mechanics — BaseStats before bonuses, FinalStats after bonuses
	BaseStats  Stats
	FinalStats Stats

	// Modifiers MUST be calculated from FinalStats, never BaseStats
	Modifiers Modifiers

	// Derived from FinalStats + armor
	Derived DerivedStats

	// Narrative blocks
	Background NarrativeBlock
	Motivation NarrativeBlock
	Secret     NarrativeBlock

	// Regeneration state
	Locks CharacterLocks

	// Optional seed for reproducibility — same seed + same params = same result
	Seed *int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

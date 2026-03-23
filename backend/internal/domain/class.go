package domain

// Class represents a D&D 5e character class.
type Class string

const (
	ClassBarbarian Class = "barbarian"
	ClassBard      Class = "bard"
	ClassCleric    Class = "cleric"
	ClassDruid     Class = "druid"
	ClassFighter   Class = "fighter"
	ClassMonk      Class = "monk"
	ClassPaladin   Class = "paladin"
	ClassRanger    Class = "ranger"
	ClassRogue     Class = "rogue"
	ClassSorcerer  Class = "sorcerer"
	ClassWarlock   Class = "warlock"
	ClassWizard    Class = "wizard"
	ClassArtificer Class = "artificer"
)

// Ruleset identifies the D&D ruleset in use.
type Ruleset string

const (
	Ruleset5e  Ruleset = "5e"
	Ruleset55e Ruleset = "5.5e"
)

// AbilityBonusSource controls how ability bonuses are resolved.
type AbilityBonusSource string

const (
	AbilityBonusFromSpecies    AbilityBonusSource = "species"
	AbilityBonusFromBackground AbilityBonusSource = "background"
	AbilityBonusNone           AbilityBonusSource = "none"
)

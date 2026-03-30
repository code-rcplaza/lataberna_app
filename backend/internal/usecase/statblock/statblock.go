package statblock

import (
	"math/rand"
	"time"

	"forge-rpg/internal/domain"
)

// Input for stat block generation. All fields optional — omitted = random.
type Input struct {
	Class      *domain.Class
	Species    *domain.Species
	SubSpecies *domain.SubSpecies
	Level      int   // defaults to 1 if zero
	Seed       *int64
}

// Output of a stat block generation.
type Output struct {
	Class      domain.Class
	Species    domain.Species
	SubSpecies *domain.SubSpecies
	Level      int
	BaseStats  domain.Stats
	FinalStats domain.Stats
	Modifiers  domain.Modifiers
	Derived    domain.DerivedStats
	Armor      domain.ArmorType
	Seed       int64
}

// Generate produces a character stat block based on the input parameters.
// Same Seed + same Input always returns the same Output.
func Generate(in Input) (Output, error) {
	seed := resolveSeed(in.Seed)
	rng := rand.New(rand.NewSource(seed))

	class := resolveClass(in.Class, rng)
	species, subSpecies := resolveSpecies(in.Species, in.SubSpecies, rng)
	level := resolveLevel(in.Level)

	baseStats, err := generateBaseStats(class, rng)
	if err != nil {
		return Output{}, err
	}

	// TODO(5.5e): replace with background ASI resolution
	var bonuses []domain.AbilityBonus
	finalStats := applyBonuses(baseStats, bonuses)
	modifiers := calculateModifiers(finalStats)

	armor, err := resolveArmor(class)
	if err != nil {
		return Output{}, err
	}

	derived := domain.DerivedStats{
		HP: calculateHP(class, modifiers, level),
		AC: calculateAC(armor, modifiers),
	}

	return Output{
		Class:      class,
		Species:    species,
		SubSpecies: subSpecies,
		Level:      level,
		BaseStats:  baseStats,
		FinalStats: finalStats,
		Modifiers:  modifiers,
		Derived:    derived,
		Armor:      armor,
		Seed:       seed,
	}, nil
}

func resolveSeed(s *int64) int64 {
	if s != nil {
		return *s
	}
	return time.Now().UnixNano()
}

func resolveLevel(l int) int {
	if l <= 0 {
		return 1
	}
	return l
}

func resolveClass(c *domain.Class, rng *rand.Rand) domain.Class {
	if c != nil {
		return *c
	}
	all := allClasses()
	return all[rng.Intn(len(all))]
}

func resolveSpecies(s *domain.Species, sub *domain.SubSpecies, rng *rand.Rand) (domain.Species, *domain.SubSpecies) {
	if s != nil {
		// Species is fixed — resolve subspecies
		resolved := resolveSubSpecies(*s, sub, rng)
		return *s, resolved
	}
	// Pick random species from all species
	all := allSpecies()
	picked := all[rng.Intn(len(all))]
	resolved := resolveSubSpecies(picked, nil, rng)
	return picked, resolved
}

func resolveSubSpecies(s domain.Species, sub *domain.SubSpecies, rng *rand.Rand) *domain.SubSpecies {
	if sub != nil {
		return sub
	}
	subs := subSpeciesFor(s)
	if len(subs) == 0 {
		return nil
	}
	picked := subs[rng.Intn(len(subs))]
	return &picked
}

func generateBaseStats(class domain.Class, rng *rand.Rand) (domain.Stats, error) {
	data, ok := classTable[class]
	if !ok {
		return domain.Stats{}, errUnknownClass(class)
	}
	return buildStatsFromPriority(data.primaryStat, data.secondaryStat, rng), nil
}

func applyBonuses(base domain.Stats, bonuses []domain.AbilityBonus) domain.Stats {
	result := base
	for _, b := range bonuses {
		switch b.Stat {
		case "STR":
			result.STR += b.Value
		case "DEX":
			result.DEX += b.Value
		case "CON":
			result.CON += b.Value
		case "INT":
			result.INT += b.Value
		case "WIS":
			result.WIS += b.Value
		case "CHA":
			result.CHA += b.Value
		}
	}
	return result
}

func resolveArmor(class domain.Class) (domain.ArmorType, error) {
	data, ok := classTable[class]
	if !ok {
		return domain.ArmorType{}, errUnknownClass(class)
	}
	armor, ok := armorTable[data.armorKey]
	if !ok {
		return domain.ArmorType{}, errUnknownArmor(data.armorKey)
	}
	return armor, nil
}

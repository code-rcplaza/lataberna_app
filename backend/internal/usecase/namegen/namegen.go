package namegen

import (
	"fmt"
	"math/rand"
	"time"

	"forge-rpg/internal/domain"
)

// Gender of the generated character name.
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

// Input parameters for name generation. All fields are optional.
// Omitted fields are resolved randomly.
type Input struct {
	Species    domain.Species
	SubSpecies *domain.SubSpecies
	Gender     *Gender
	Seed       *int64
}

// Output of a name generation.
type Output struct {
	Name string
	Seed int64 // the seed used — allows reproduction of the same result
}

// Generate produces a character name based on the input parameters.
// Same Seed + same Input always returns the same Output.Name.
func Generate(in Input) (Output, error) {
	seed := resolveSeed(in.Seed)
	rng := rand.New(rand.NewSource(seed))

	gender := resolveGender(in.Gender, rng)

	names, err := namesFor(in.Species, in.SubSpecies, gender, rng)
	if err != nil {
		return Output{}, err
	}

	if len(names) == 0 {
		return Output{}, fmt.Errorf("namegen: no names found for species %q gender %q", in.Species, gender)
	}

	name := names[rng.Intn(len(names))]
	return Output{Name: name, Seed: seed}, nil
}

func resolveSeed(s *int64) int64 {
	if s != nil {
		return *s
	}
	return time.Now().UnixNano()
}

func resolveGender(g *Gender, rng *rand.Rand) Gender {
	if g != nil {
		return *g
	}
	if rng.Intn(2) == 0 {
		return GenderMale
	}
	return GenderFemale
}

func namesFor(species domain.Species, sub *domain.SubSpecies, gender Gender, rng *rand.Rand) ([]string, error) {
	key := speciesKey(species, sub, rng)
	pool, ok := nameData[key]
	if !ok {
		return nil, fmt.Errorf("namegen: unknown species key %q", key)
	}
	switch gender {
	case GenderMale:
		return pool.male, nil
	case GenderFemale:
		return pool.female, nil
	default:
		return nil, fmt.Errorf("namegen: unknown gender %q", gender)
	}
}

// speciesKey resolves which data key to use.
// For species with subspecies, picks randomly if none is provided.
func speciesKey(species domain.Species, sub *domain.SubSpecies, rng *rand.Rand) string {
	if sub != nil {
		return string(*sub)
	}
	switch species {
	case domain.SpeciesHuman, domain.SpeciesHalfElf, domain.SpeciesHalfOrc,
		domain.SpeciesTiefling, domain.SpeciesDragonborn:
		return string(species)
	case domain.SpeciesElf:
		subs := []string{string(domain.SubSpeciesHighElf), string(domain.SubSpeciesWoodElf), string(domain.SubSpeciesDrow)}
		return subs[rng.Intn(len(subs))]
	case domain.SpeciesDwarf:
		subs := []string{string(domain.SubSpeciesHillDwarf), string(domain.SubSpeciesMountainDwarf)}
		return subs[rng.Intn(len(subs))]
	case domain.SpeciesHalfling:
		subs := []string{string(domain.SubSpeciesLightfoot), string(domain.SubSpeciesStout)}
		return subs[rng.Intn(len(subs))]
	case domain.SpeciesGnome:
		subs := []string{string(domain.SubSpeciesForestGnome), string(domain.SubSpeciesRockGnome)}
		return subs[rng.Intn(len(subs))]
	default:
		return string(species)
	}
}

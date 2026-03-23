package narrativegen

import (
	"fmt"
	"math/rand"
	"time"

	"forge-rpg/internal/domain"
)

// Input for narrative generation. All fields optional — omitted = random.
type Input struct {
	Class      *domain.Class
	Species    *domain.Species
	Categories []domain.NarrativeCategory // nil = generate all three
	Seed       *int64
}

// Output of a narrative generation — always three blocks.
type Output struct {
	Background domain.NarrativeBlock
	Motivation domain.NarrativeBlock
	Secret     domain.NarrativeBlock
	Seed       int64
}

// Generate produces three NarrativeBlocks filtered by class and species.
// Same Seed + same Input always returns the same Output.
func Generate(in Input) (Output, error) {
	seed := resolveSeed(in.Seed)
	rng := rand.New(rand.NewSource(seed))

	class := resolveClass(in.Class, rng)
	species := resolveSpecies(in.Species, rng)

	background, err := generateBlock(domain.NarrativeBackground, class, species, rng)
	if err != nil {
		return Output{}, err
	}

	motivation, err := generateBlock(domain.NarrativeMotivation, class, species, rng)
	if err != nil {
		return Output{}, err
	}

	secret, err := generateBlock(domain.NarrativeSecret, class, species, rng)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Background: background,
		Motivation: motivation,
		Secret:     secret,
		Seed:       seed,
	}, nil
}

func resolveSeed(s *int64) int64 {
	if s != nil {
		return *s
	}
	return time.Now().UnixNano()
}

func resolveClass(c *domain.Class, rng *rand.Rand) domain.Class {
	if c != nil {
		return *c
	}
	all := allClasses()
	return all[rng.Intn(len(all))]
}

func resolveSpecies(s *domain.Species, rng *rand.Rand) domain.Species {
	if s != nil {
		return *s
	}
	all := allSpecies()
	return all[rng.Intn(len(all))]
}

func allClasses() []domain.Class {
	return []domain.Class{
		domain.ClassBarbarian, domain.ClassBard, domain.ClassCleric, domain.ClassDruid,
		domain.ClassFighter, domain.ClassMonk, domain.ClassPaladin, domain.ClassRanger,
		domain.ClassRogue, domain.ClassSorcerer, domain.ClassWarlock, domain.ClassWizard,
		domain.ClassArtificer,
	}
}

func allSpecies() []domain.Species {
	return []domain.Species{
		domain.SpeciesHuman, domain.SpeciesElf, domain.SpeciesDwarf, domain.SpeciesHalfling,
		domain.SpeciesGnome, domain.SpeciesHalfElf, domain.SpeciesHalfOrc,
		domain.SpeciesTiefling, domain.SpeciesDragonborn,
	}
}

func generateBlock(
	category domain.NarrativeCategory,
	class domain.Class,
	species domain.Species,
	rng *rand.Rand,
) (domain.NarrativeBlock, error) {
	pool := templatesFor(category)

	compatible := filterCompatible(pool, class, species)
	if len(compatible) == 0 {
		return domain.NarrativeBlock{}, fmt.Errorf(
			"narrativegen: no compatible templates for category %q class %q species %q",
			category, class, species,
		)
	}

	return compatible[rng.Intn(len(compatible))], nil
}

func templatesFor(category domain.NarrativeCategory) []domain.NarrativeBlock {
	return narrativeTemplates[category]
}

func filterCompatible(
	templates []domain.NarrativeBlock,
	class domain.Class,
	species domain.Species,
) []domain.NarrativeBlock {
	var out []domain.NarrativeBlock
	for _, t := range templates {
		if isCompatible(t, class, species) {
			out = append(out, t)
		}
	}
	return out
}

func isCompatible(t domain.NarrativeBlock, class domain.Class, species domain.Species) bool {
	for _, tag := range t.Tags {
		if tag == "any" || tag == string(class) || tag == string(species) {
			return true
		}
	}
	return false
}

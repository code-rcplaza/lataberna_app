package namegen

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
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

// Service generates character names using a NameRepository for content.
// Inject it via New — the repo is queried on every Generate call.
type Service struct {
	repo ports.NameRepository
}

// New constructs a Service with the given repository.
func New(repo ports.NameRepository) *Service {
	return &Service{repo: repo}
}

// Generate produces a character name based on the input parameters.
// Same Seed + same Input always returns the same Output.Name.
func (s *Service) Generate(ctx context.Context, in Input) (Output, error) {
	seed := resolveSeed(in.Seed)
	rng := rand.New(rand.NewSource(seed))

	gender := resolveGender(in.Gender, rng)

	names, err := s.namesFor(ctx, in.Species, in.SubSpecies, gender, rng)
	if err != nil {
		return Output{}, err
	}

	if len(names) == 0 {
		return Output{}, fmt.Errorf("namegen: no names found for species %q gender %q", in.Species, gender)
	}

	name := names[rng.Intn(len(names))]
	return Output{Name: name, Seed: seed}, nil
}

func (s *Service) namesFor(
	ctx     context.Context,
	species domain.Species,
	sub     *domain.SubSpecies,
	gender  Gender,
	rng     *rand.Rand,
) ([]string, error) {
	key := speciesKey(species, sub, rng)
	names, err := s.repo.FindBySpeciesGender(ctx, key, string(gender))
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("namegen: unknown species key %q", key)
	}
	return names, nil
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

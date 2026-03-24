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
	key := speciesKey(in.Species, in.SubSpecies, rng)

	name, err := s.compose(ctx, key, string(gender), rng)
	if err != nil {
		return Output{}, err
	}
	return Output{Name: name, Seed: seed}, nil
}

// compose assembles a species-appropriate name from component pools.
// The speciesKey determines which pools and assembly rule to use.
func (s *Service) compose(ctx context.Context, key, gender string, rng *rand.Rand) (string, error) {
	switch key {
	case "human":
		first, err := s.pick(ctx, key, gender, "first_name", rng)
		if err != nil {
			return "", err
		}
		surname, err := s.pick(ctx, key, "any", "surname", rng)
		if err != nil {
			return "", err
		}
		return first + " " + surname, nil

	case "hill-dwarf", "mountain-dwarf":
		first, err := s.pick(ctx, key, gender, "first_name", rng)
		if err != nil {
			return "", err
		}
		clan, err := s.pick(ctx, key, "any", "clan_name", rng)
		if err != nil {
			return "", err
		}
		return first + " " + clan, nil

	case "high-elf", "wood-elf", "drow":
		first, err := s.pick(ctx, key, gender, "first_name", rng)
		if err != nil {
			return "", err
		}
		family, err := s.pick(ctx, key, "any", "family_name", rng)
		if err != nil {
			return "", err
		}
		return first + " " + family, nil

	case "lightfoot", "stout":
		first, err := s.pick(ctx, key, gender, "first_name", rng)
		if err != nil {
			return "", err
		}
		surname, err := s.pick(ctx, key, "any", "surname", rng)
		if err != nil {
			return "", err
		}
		return first + " " + surname, nil

	case "dragonborn":
		// Clan name precedes first name — hard invariant.
		clan, err := s.pick(ctx, key, "any", "clan_name", rng)
		if err != nil {
			return "", err
		}
		first, err := s.pick(ctx, key, gender, "first_name", rng)
		if err != nil {
			return "", err
		}
		return clan + " " + first, nil

	case "forest-gnome", "rock-gnome":
		first, err := s.pick(ctx, key, gender, "first_name", rng)
		if err != nil {
			return "", err
		}
		clan, err := s.pick(ctx, key, "any", "clan_name", rng)
		if err != nil {
			return "", err
		}
		nick, err := s.pick(ctx, key, "any", "nickname", rng)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(`%s %s "%s"`, first, clan, nick), nil

	case "half-elf":
		first, err := s.pick(ctx, key, gender, "first_name", rng)
		if err != nil {
			return "", err
		}
		// Randomly adopts human (surname) or elven (family_name) convention.
		nameType := "surname"
		if rng.Intn(2) != 0 {
			nameType = "family_name"
		}
		second, err := s.pick(ctx, key, "any", nameType, rng)
		if err != nil {
			return "", err
		}
		return first + " " + second, nil

	case "half-orc":
		first, err := s.pick(ctx, key, gender, "first_name", rng)
		if err != nil {
			return "", err
		}
		// ~30% chance of also having a surname.
		if rng.Intn(10) < 3 {
			surname, err := s.pick(ctx, key, "any", "surname", rng)
			if err != nil {
				return "", err
			}
			return first + " " + surname, nil
		}
		return first, nil

	case string(domain.SubSpeciesInfernalTiefling):
		return s.pick(ctx, key, "any", "infernal_name", rng)

	case string(domain.SubSpeciesVirtueTiefling):
		return s.pick(ctx, key, "any", "virtue_word", rng)

	default:
		return "", fmt.Errorf("namegen.compose: unknown species key %q", key)
	}
}

// pick fetches a pool via FindByType and returns one random entry.
func (s *Service) pick(ctx context.Context, speciesKey, gender, nameType string, rng *rand.Rand) (string, error) {
	pool, err := s.repo.FindByType(ctx, speciesKey, gender, nameType)
	if err != nil {
		return "", fmt.Errorf("namegen.compose: %w", err)
	}
	return pickOne(pool, rng)
}

// pickOne returns a uniformly random element from pool.
// Returns an error wrapping ports.ErrEmptyNamePool if pool is empty.
func pickOne(pool []string, rng *rand.Rand) (string, error) {
	if len(pool) == 0 {
		return "", fmt.Errorf("pickOne: %w", ports.ErrEmptyNamePool)
	}
	return pool[rng.Intn(len(pool))], nil
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
	case domain.SpeciesHuman, domain.SpeciesHalfElf, domain.SpeciesHalfOrc, domain.SpeciesDragonborn:
		return string(species)
	case domain.SpeciesTiefling:
		// Randomly selects infernal or virtue naming tradition.
		if rng.Intn(2) == 0 {
			return string(domain.SubSpeciesInfernalTiefling)
		}
		return string(domain.SubSpeciesVirtueTiefling)
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

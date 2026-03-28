package character

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/usecase/namegen"
	"forge-rpg/internal/usecase/narrativegen"
	"forge-rpg/internal/usecase/statblock"
)

// CreateInput — all fields optional. Omitted fields generate randomly.
// Provided fields become implicit locks (they are preserved as-is).
type CreateInput struct {
	Name       *string
	Class      *domain.Class
	Species    *domain.Species
	SubSpecies *domain.SubSpecies
	Gender     *namegen.Gender // used for name generation only
	Seed       *int64
}

// RegenerateInput — regenerates unlocked fields of an existing character.
// Locked fields are preserved exactly as they are.
type RegenerateInput struct {
	Character *domain.Character
	Locks     domain.CharacterLocks
	Seed      *int64
}

// Creator orchestrates the full character generation pipeline.
// It depends on narrativegen.Service and namegen.Service for content generation.
type Creator struct {
	narrativeSvc *narrativegen.Service
	nameSvc      *namegen.Service
}

// NewCreator constructs a Creator with the given services.
func NewCreator(narrativeSvc *narrativegen.Service, nameSvc *namegen.Service) *Creator {
	return &Creator{
		narrativeSvc: narrativeSvc,
		nameSvc:      nameSvc,
	}
}

// Create runs the full 9-step pipeline and returns a fully populated Character.
// Same Seed + same CreateInput always returns the same Character.
func (c *Creator) Create(ctx context.Context, in CreateInput) (*domain.Character, error) {
	seed := resolveSeed(in.Seed)
	rng := rand.New(rand.NewSource(seed))

	// Derive independent sub-seeds from the main seed.
	// This guarantees same main seed → same character always.
	nameSeed := rng.Int63()
	statSeed := rng.Int63()
	narrativeSeed := rng.Int63()

	// Step 1: resolve class and species (needed across all modules).
	class, species, subSpecies := resolveIdentity(in, rng)

	// Steps 2–8: generate stat block.
	statOut, err := statblock.Generate(statblock.Input{
		Class:      &class,
		Species:    &species,
		SubSpecies: subSpecies,
		Level:      1,
		Seed:       &statSeed,
	})
	if err != nil {
		return nil, fmt.Errorf("character.Create: stat block: %w", err)
	}

	// Step 1 (name): generate name using class-resolved species.
	name, err := c.resolveName(ctx, in.Name, species, subSpecies, in.Gender, nameSeed)
	if err != nil {
		return nil, fmt.Errorf("character.Create: name: %w", err)
	}

	// Step 9: generate narrative.
	narrOut, err := c.narrativeSvc.Generate(ctx, narrativegen.Input{
		Class:   &class,
		Species: &species,
		Seed:    &narrativeSeed,
	})
	if err != nil {
		return nil, fmt.Errorf("character.Create: narrative: %w", err)
	}

	now := time.Now()
	return &domain.Character{
		ID:                 newID(),
		Name:               name,
		Species:            species,
		SubSpecies:         subSpecies,
		Class:              class,
		Level:              1,
		Ruleset:            domain.Ruleset5e,
		AbilityBonusSource: domain.AbilityBonusFromSpecies,
		BaseStats:          statOut.BaseStats,
		FinalStats:         statOut.FinalStats,
		Modifiers:          statOut.Modifiers,
		Derived:            statOut.Derived,
		Background:         narrOut.Background,
		Motivation:         narrOut.Motivation,
		Secret:             narrOut.Secret,
		Locks:              domain.CharacterLocks{}, // nothing locked on creation
		Seed:               &seed,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

// Regenerate re-executes only the unlocked fields of an existing character.
// Locked fields are cloned from the input character unchanged.
func (c *Creator) Regenerate(ctx context.Context, in RegenerateInput) (*domain.Character, error) {
	if in.Character == nil {
		return nil, fmt.Errorf("character.Regenerate: character is required")
	}

	seed := resolveSeed(in.Seed)

	// Clone the existing character — start from its current state.
	updated := *in.Character
	updated.Locks = in.Locks
	updated.UpdatedAt = time.Now()

	rng := rand.New(rand.NewSource(seed))
	nameSeed := rng.Int63()
	statSeed := rng.Int63()
	narrativeSeed := rng.Int63()

	if !in.Locks.Name {
		name, err := c.resolveName(ctx, nil, updated.Species, updated.SubSpecies, nil, nameSeed)
		if err != nil {
			return nil, fmt.Errorf("character.Regenerate: name: %w", err)
		}
		updated.Name = name
	}

	if !in.Locks.Stats {
		statOut, err := statblock.Generate(statblock.Input{
			Class:      &updated.Class,
			Species:    &updated.Species,
			SubSpecies: updated.SubSpecies,
			Level:      updated.Level,
			Seed:       &statSeed,
		})
		if err != nil {
			return nil, fmt.Errorf("character.Regenerate: stats: %w", err)
		}
		updated.BaseStats = statOut.BaseStats
		updated.FinalStats = statOut.FinalStats
		updated.Modifiers = statOut.Modifiers
		updated.Derived = statOut.Derived
	}

	if !in.Locks.Background || !in.Locks.Motivation || !in.Locks.Secret {
		narrOut, err := c.narrativeSvc.Generate(ctx, narrativegen.Input{
			Class:   &updated.Class,
			Species: &updated.Species,
			Seed:    &narrativeSeed,
		})
		if err != nil {
			return nil, fmt.Errorf("character.Regenerate: narrative: %w", err)
		}
		if !in.Locks.Background {
			updated.Background = narrOut.Background
		}
		if !in.Locks.Motivation {
			updated.Motivation = narrOut.Motivation
		}
		if !in.Locks.Secret {
			updated.Secret = narrOut.Secret
		}
	}

	updated.Seed = &seed
	return &updated, nil
}

// resolveSeed returns the provided seed or a random one based on current time.
// Auto-generated seeds are clamped to int32 range to prevent precision loss when
// serialized through GraphQL Int (32-bit) and parsed by JavaScript (float64).
// UnixNano values (~1.7e18) exceed Number.MAX_SAFE_INTEGER, making round-trips lossy.
func resolveSeed(s *int64) int64 {
	if s != nil {
		return *s
	}
	return time.Now().UnixNano() % (1<<31 - 1)
}

// newID generates a simple unique ID for the character.
// A proper UUID library is infrastructure concern, not domain.
func newID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// resolveIdentity determines class, species, and subSpecies.
// If provided in input, those values are used directly.
// Otherwise they are chosen randomly from all valid options.
func resolveIdentity(in CreateInput, rng *rand.Rand) (domain.Class, domain.Species, *domain.SubSpecies) {
	class := resolveClass(in.Class, rng)
	species, subSpecies := resolveSpecies(in.Species, in.SubSpecies, rng)
	return class, species, subSpecies
}

// resolveClass returns the provided class or a random one.
func resolveClass(c *domain.Class, rng *rand.Rand) domain.Class {
	if c != nil {
		return *c
	}
	all := allClasses()
	return all[rng.Intn(len(all))]
}

// resolveSpecies returns the provided species (and subSpecies) or random ones.
func resolveSpecies(s *domain.Species, sub *domain.SubSpecies, rng *rand.Rand) (domain.Species, *domain.SubSpecies) {
	if s != nil {
		resolved := resolveSubSpecies(*s, sub, rng)
		return *s, resolved
	}
	all := allSpecies()
	picked := all[rng.Intn(len(all))]
	resolved := resolveSubSpecies(picked, nil, rng)
	return picked, resolved
}

// resolveSubSpecies returns the provided subSpecies or picks a valid one at random.
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

// resolveName returns the provided name directly or calls nameSvc.Generate.
func (c *Creator) resolveName(
	ctx        context.Context,
	name       *string,
	species    domain.Species,
	subSpecies *domain.SubSpecies,
	gender     *namegen.Gender,
	seed       int64,
) (string, error) {
	if name != nil {
		return *name, nil
	}
	out, err := c.nameSvc.Generate(ctx, namegen.Input{
		Species:    species,
		SubSpecies: subSpecies,
		Gender:     gender,
		Seed:       &seed,
	})
	if err != nil {
		return "", err
	}
	return out.Name, nil
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

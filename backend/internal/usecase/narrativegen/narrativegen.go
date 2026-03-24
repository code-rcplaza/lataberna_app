package narrativegen

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
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

// Service generates narrative blocks using a NarrativeRepository for content.
// Inject it via New — the repo is queried on every Generate call.
type Service struct {
	repo ports.NarrativeRepository
}

// New constructs a Service with the given repository.
func New(repo ports.NarrativeRepository) *Service {
	return &Service{repo: repo}
}

// Generate produces three NarrativeBlocks weighted by class and species compatibility.
// Same Seed + same Input always returns the same Output.
func (s *Service) Generate(ctx context.Context, in Input) (Output, error) {
	seed := resolveSeed(in.Seed)
	rng := rand.New(rand.NewSource(seed))

	class := resolveClass(in.Class, rng)
	species := resolveSpecies(in.Species, rng)

	background, err := s.generateBlock(ctx, domain.NarrativeBackground, class, species, rng)
	if err != nil {
		return Output{}, err
	}

	motivation, err := s.generateBlock(ctx, domain.NarrativeMotivation, class, species, rng)
	if err != nil {
		return Output{}, err
	}

	secret, err := s.generateBlock(ctx, domain.NarrativeSecret, class, species, rng)
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

func (s *Service) generateBlock(
	ctx      context.Context,
	category domain.NarrativeCategory,
	class    domain.Class,
	species  domain.Species,
	rng      *rand.Rand,
) (domain.NarrativeBlock, error) {
	pool, err := s.repo.FindByCategory(ctx, category, class, species)
	if err != nil {
		return domain.NarrativeBlock{}, err
	}
	return weightedPick(pool, rng)
}

// weightedPick selects one entry from pool using the seeded rng.
// Returns an error if the pool is empty or all weights are zero.
// Pool must be ordered deterministically (ORDER BY id in SQL query) for seed reproducibility.
func weightedPick(pool []ports.WeightedNarrativeEntry, rng *rand.Rand) (domain.NarrativeBlock, error) {
	// Build cumulative weight slice — O(n)
	total := 0
	cumulative := make([]int, len(pool))
	for i, e := range pool {
		total += e.Weight
		cumulative[i] = total
	}
	if total == 0 {
		return domain.NarrativeBlock{}, errors.New("narrativegen: no compatible entries (all excluded or empty pool)")
	}
	// Draw a number in [0, total)
	draw := rng.Intn(total)
	// Linear scan for first cumulative[i] > draw
	for i, c := range cumulative {
		if draw < c {
			return pool[i].Block, nil
		}
	}
	// Unreachable — return last entry as safety fallback
	return pool[len(pool)-1].Block, nil
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

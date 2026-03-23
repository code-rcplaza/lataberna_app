package graph

import (
	"context"
	"fmt"
	"time"

	"forge-rpg/graph/model"
	"forge-rpg/internal/domain"
	graphqlmw "forge-rpg/internal/infrastructure/graphql"
	"forge-rpg/internal/usecase/character"
	"forge-rpg/internal/usecase/namegen"
)

// ---------------------------------------------------------------------------
// Context helpers
// ---------------------------------------------------------------------------

// userFromCtx extracts the authenticated user from the request context.
// Returns an "unauthorized" error if the user is absent.
func userFromCtx(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value(graphqlmw.UserContextKey).(*domain.User)
	if !ok || user == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	return user, nil
}

// ---------------------------------------------------------------------------
// Domain ↔ model conversion helpers
// ---------------------------------------------------------------------------

func toModelStats(s domain.Stats) *model.Stats {
	return &model.Stats{
		Str: s.STR,
		Dex: s.DEX,
		Con: s.CON,
		Int: s.INT,
		Wis: s.WIS,
		Cha: s.CHA,
	}
}

func toModelModifiers(m domain.Modifiers) *model.Stats {
	return &model.Stats{
		Str: m.STR,
		Dex: m.DEX,
		Con: m.CON,
		Int: m.INT,
		Wis: m.WIS,
		Cha: m.CHA,
	}
}

func toModelCharacter(c *domain.Character) *model.Character {
	mc := &model.Character{
		ID:      c.ID,
		Name:    c.Name,
		Species: string(c.Species),
		Class:   string(c.Class),
		Level:   c.Level,
		Ruleset: string(c.Ruleset),
		BaseStats:  toModelStats(c.BaseStats),
		FinalStats: toModelStats(c.FinalStats),
		Modifiers:  toModelModifiers(c.Modifiers),
		Derived: &model.DerivedStats{
			Hp: c.Derived.HP,
			Ac: c.Derived.AC,
		},
		Background: &model.NarrativeBlock{
			Category: string(c.Background.Category),
			Content:  c.Background.Content,
			Tags:     tagsOrEmpty(c.Background.Tags),
		},
		Motivation: &model.NarrativeBlock{
			Category: string(c.Motivation.Category),
			Content:  c.Motivation.Content,
			Tags:     tagsOrEmpty(c.Motivation.Tags),
		},
		Secret: &model.NarrativeBlock{
			Category: string(c.Secret.Category),
			Content:  c.Secret.Content,
			Tags:     tagsOrEmpty(c.Secret.Tags),
		},
		Locks: &model.CharacterLocks{
			Name:       c.Locks.Name,
			Stats:      c.Locks.Stats,
			Background: c.Locks.Background,
			Motivation: c.Locks.Motivation,
			Secret:     c.Locks.Secret,
		},
		CreatedAt: c.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.UTC().Format(time.RFC3339),
	}

	if c.SubSpecies != nil {
		s := string(*c.SubSpecies)
		mc.SubSpecies = &s
	}

	if c.Seed != nil {
		seedInt := int(*c.Seed)
		mc.Seed = &seedInt
	}

	return mc
}

func toModelCharacters(chars []*domain.Character) []*model.Character {
	out := make([]*model.Character, len(chars))
	for i, c := range chars {
		out[i] = toModelCharacter(c)
	}
	return out
}

func tagsOrEmpty(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	return tags
}

// inputToCreateInput converts a *model.GenerateCharacterInput to character.CreateInput.
// The seed field in GraphQL is int; domain uses int64.
func inputToCreateInput(input *model.GenerateCharacterInput, seedOverride *int64) character.CreateInput {
	ci := character.CreateInput{}

	if input == nil {
		ci.Seed = seedOverride
		return ci
	}

	if input.Name != nil {
		ci.Name = input.Name
	}
	if input.Class != nil {
		c := domain.Class(*input.Class)
		ci.Class = &c
	}
	if input.Species != nil {
		s := domain.Species(*input.Species)
		ci.Species = &s
	}
	if input.SubSpecies != nil {
		ss := domain.SubSpecies(*input.SubSpecies)
		ci.SubSpecies = &ss
	}
	if input.Gender != nil {
		g := namegen.Gender(*input.Gender)
		ci.Gender = &g
	}
	if seedOverride != nil {
		ci.Seed = seedOverride
	} else if input.Seed != nil {
		s := int64(*input.Seed)
		ci.Seed = &s
	}

	return ci
}

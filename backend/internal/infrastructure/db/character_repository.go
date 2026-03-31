package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
)

// Compile-time interface check.
var _ ports.CharacterRepository = (*CharacterRepository)(nil)

// CharacterRepository implements ports.CharacterRepository using SQLite.
// Nested structs (Stats, Modifiers, DerivedStats, NarrativeBlock, CharacterLocks)
// are stored as JSON columns for MVP simplicity.
type CharacterRepository struct {
	db *sql.DB
}

// NewCharacterRepository creates a new CharacterRepository.
func NewCharacterRepository(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

// Save inserts a new character record.
func (r *CharacterRepository) Save(ctx context.Context, c *domain.Character) error {
	row, err := marshalCharacter(c)
	if err != nil {
		return fmt.Errorf("character_repository.Save: marshal: %w", err)
	}

	const q = `
INSERT INTO characters (
	id, user_id, name, species, sub_species, class, level,
	ruleset, ability_bonus_source, background_type, asi_distribution,
	base_stats, final_stats, modifiers, derived,
	background, motivation, secret, locks,
	seed, created_at, updated_at
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err = r.db.ExecContext(ctx, q,
		row.id, row.userID, row.name, row.species, row.subSpecies,
		row.class, row.level, row.ruleset, row.abilityBonusSource,
		row.backgroundType, row.asiDistribution,
		row.baseStats, row.finalStats, row.modifiers, row.derived,
		row.background, row.motivation, row.secret, row.locks,
		row.seed, row.createdAt, row.updatedAt,
	)
	if err != nil {
		return fmt.Errorf("character_repository.Save: exec: %w", err)
	}
	return nil
}

// FindByID returns the character with the given ID, or nil if not found.
func (r *CharacterRepository) FindByID(ctx context.Context, id string) (*domain.Character, error) {
	const q = `
SELECT id, user_id, name, species, sub_species, class, level,
       ruleset, ability_bonus_source, background_type, asi_distribution,
       base_stats, final_stats, modifiers, derived,
       background, motivation, secret, locks,
       seed, created_at, updated_at
FROM characters WHERE id = ? LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, id)
	c, err := scanCharacter(row)
	if err != nil {
		return nil, fmt.Errorf("character_repository.FindByID: %w", err)
	}
	return c, nil
}

// FindByUserID returns all characters belonging to userID.
// Returns an empty (non-nil) slice when no results are found.
func (r *CharacterRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Character, error) {
	const q = `
SELECT id, user_id, name, species, sub_species, class, level,
       ruleset, ability_bonus_source, background_type, asi_distribution,
       base_stats, final_stats, modifiers, derived,
       background, motivation, secret, locks,
       seed, created_at, updated_at
FROM characters WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("character_repository.FindByUserID: query: %w", err)
	}
	defer rows.Close()

	out := make([]*domain.Character, 0)
	for rows.Next() {
		c, err := scanCharacter(rows)
		if err != nil {
			return nil, fmt.Errorf("character_repository.FindByUserID: scan: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("character_repository.FindByUserID: rows: %w", err)
	}
	return out, nil
}

// Update replaces all mutable fields of an existing character.
func (r *CharacterRepository) Update(ctx context.Context, c *domain.Character) error {
	row, err := marshalCharacter(c)
	if err != nil {
		return fmt.Errorf("character_repository.Update: marshal: %w", err)
	}

	const q = `
UPDATE characters SET
	user_id = ?, name = ?, species = ?, sub_species = ?, class = ?, level = ?,
	ruleset = ?, ability_bonus_source = ?, background_type = ?, asi_distribution = ?,
	base_stats = ?, final_stats = ?, modifiers = ?, derived = ?,
	background = ?, motivation = ?, secret = ?, locks = ?,
	seed = ?, updated_at = ?
WHERE id = ?`

	res, err := r.db.ExecContext(ctx, q,
		row.userID, row.name, row.species, row.subSpecies, row.class, row.level,
		row.ruleset, row.abilityBonusSource, row.backgroundType, row.asiDistribution,
		row.baseStats, row.finalStats, row.modifiers, row.derived,
		row.background, row.motivation, row.secret, row.locks,
		row.seed, row.updatedAt,
		row.id,
	)
	if err != nil {
		return fmt.Errorf("character_repository.Update: exec: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("character_repository.Update: rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("character_repository.Update: character %q not found", c.ID)
	}
	return nil
}

// Delete removes the character with the given ID.
func (r *CharacterRepository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM characters WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("character_repository.Delete: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// dbCharacterRow holds the flat, serialised values for a characters row.
type dbCharacterRow struct {
	id                 string
	userID             string
	name               string
	species            string
	subSpecies         *string
	class              string
	level              int
	ruleset            string
	abilityBonusSource string
	backgroundType     string
	asiDistribution    string
	baseStats          string
	finalStats         string
	modifiers          string
	derived            string
	background         string
	motivation         string
	secret             string
	locks              string
	seed               *int64
	createdAt          string
	updatedAt          string
}

func marshalCharacter(c *domain.Character) (*dbCharacterRow, error) {
	marshal := func(v any) (string, error) {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	baseStats, err := marshal(c.BaseStats)
	if err != nil {
		return nil, err
	}
	finalStats, err := marshal(c.FinalStats)
	if err != nil {
		return nil, err
	}
	modifiers, err := marshal(c.Modifiers)
	if err != nil {
		return nil, err
	}
	derived, err := marshal(c.Derived)
	if err != nil {
		return nil, err
	}
	background, err := marshal(c.Background)
	if err != nil {
		return nil, err
	}
	motivation, err := marshal(c.Motivation)
	if err != nil {
		return nil, err
	}
	secret, err := marshal(c.Secret)
	if err != nil {
		return nil, err
	}
	locks, err := marshal(c.Locks)
	if err != nil {
		return nil, err
	}

	var subSpecies *string
	if c.SubSpecies != nil {
		s := string(*c.SubSpecies)
		subSpecies = &s
	}

	return &dbCharacterRow{
		id:                 c.ID,
		userID:             c.UserID,
		name:               c.Name,
		species:            string(c.Species),
		subSpecies:         subSpecies,
		class:              string(c.Class),
		level:              c.Level,
		ruleset:            string(c.Ruleset),
		abilityBonusSource: string(c.AbilityBonusSource),
		backgroundType:     c.BackgroundType,
		asiDistribution:    c.ASIDistribution,
		baseStats:          baseStats,
		finalStats:         finalStats,
		modifiers:          modifiers,
		derived:            derived,
		background:         background,
		motivation:         motivation,
		secret:             secret,
		locks:              locks,
		seed:               c.Seed,
		createdAt:          c.CreatedAt.UTC().Format(time.RFC3339Nano),
		updatedAt:          c.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}, nil
}

// scanner abstracts *sql.Row and *sql.Rows so scanCharacter works for both.
type scanner interface {
	Scan(dest ...any) error
}

func scanCharacter(s scanner) (*domain.Character, error) {
	var row dbCharacterRow
	var createdAt, updatedAt string

	err := s.Scan(
		&row.id, &row.userID, &row.name, &row.species, &row.subSpecies,
		&row.class, &row.level, &row.ruleset, &row.abilityBonusSource,
		&row.backgroundType, &row.asiDistribution,
		&row.baseStats, &row.finalStats, &row.modifiers, &row.derived,
		&row.background, &row.motivation, &row.secret, &row.locks,
		&row.seed, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	unmarshal := func(s string, v any) error {
		return json.Unmarshal([]byte(s), v)
	}

	var c domain.Character
	c.ID = row.id
	c.UserID = row.userID
	c.Name = row.name
	c.Species = domain.Species(row.species)
	c.Class = domain.Class(row.class)
	c.Level = row.level
	c.Ruleset = domain.Ruleset(row.ruleset)
	c.AbilityBonusSource = domain.AbilityBonusSource(row.abilityBonusSource)
	c.BackgroundType = row.backgroundType
	c.ASIDistribution = row.asiDistribution
	c.Seed = row.seed

	if row.subSpecies != nil {
		ss := domain.SubSpecies(*row.subSpecies)
		c.SubSpecies = &ss
	}

	if err := unmarshal(row.baseStats, &c.BaseStats); err != nil {
		return nil, fmt.Errorf("unmarshal base_stats: %w", err)
	}
	if err := unmarshal(row.finalStats, &c.FinalStats); err != nil {
		return nil, fmt.Errorf("unmarshal final_stats: %w", err)
	}
	if err := unmarshal(row.modifiers, &c.Modifiers); err != nil {
		return nil, fmt.Errorf("unmarshal modifiers: %w", err)
	}
	if err := unmarshal(row.derived, &c.Derived); err != nil {
		return nil, fmt.Errorf("unmarshal derived: %w", err)
	}
	if err := unmarshal(row.background, &c.Background); err != nil {
		return nil, fmt.Errorf("unmarshal background: %w", err)
	}
	if err := unmarshal(row.motivation, &c.Motivation); err != nil {
		return nil, fmt.Errorf("unmarshal motivation: %w", err)
	}
	if err := unmarshal(row.secret, &c.Secret); err != nil {
		return nil, fmt.Errorf("unmarshal secret: %w", err)
	}
	if err := unmarshal(row.locks, &c.Locks); err != nil {
		return nil, fmt.Errorf("unmarshal locks: %w", err)
	}

	c.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse created_at: %w", err)
	}
	c.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse updated_at: %w", err)
	}

	return &c, nil
}

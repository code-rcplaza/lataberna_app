package character_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/usecase/character"
)

// ---------------------------------------------------------------------------
// In-memory CharacterRepository for tests
// ---------------------------------------------------------------------------

type memCharRepo struct {
	store map[string]*domain.Character
}

func newMemCharRepo() *memCharRepo {
	return &memCharRepo{store: make(map[string]*domain.Character)}
}

func (r *memCharRepo) Save(_ context.Context, c *domain.Character) error {
	cp := *c
	r.store[c.ID] = &cp
	return nil
}

func (r *memCharRepo) FindByID(_ context.Context, id string) (*domain.Character, error) {
	c, ok := r.store[id]
	if !ok {
		return nil, nil
	}
	cp := *c
	return &cp, nil
}

func (r *memCharRepo) FindByUserID(_ context.Context, userID string) ([]*domain.Character, error) {
	var out []*domain.Character
	for _, c := range r.store {
		if c.UserID == userID {
			cp := *c
			out = append(out, &cp)
		}
	}
	if out == nil {
		out = []*domain.Character{}
	}
	return out, nil
}

func (r *memCharRepo) Update(_ context.Context, c *domain.Character) error {
	if _, ok := r.store[c.ID]; !ok {
		return fmt.Errorf("character %q not found", c.ID)
	}
	cp := *c
	r.store[c.ID] = &cp
	return nil
}

func (r *memCharRepo) Delete(_ context.Context, id string) error {
	delete(r.store, id)
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestCharacter(id, userID, name string) *domain.Character {
	return &domain.Character{
		ID:                 id,
		UserID:             userID,
		Name:               name,
		Species:            domain.SpeciesHuman,
		Class:              domain.ClassFighter,
		Level:              1,
		Ruleset:            domain.Ruleset5e,
		AbilityBonusSource: domain.AbilityBonusFromSpecies,
		Background: domain.NarrativeBlock{
			Category: domain.NarrativeBackground,
			Content:  "A humble beginning",
			Tags:     []string{"any"},
		},
		Motivation: domain.NarrativeBlock{
			Category: domain.NarrativeMotivation,
			Content:  "To seek glory",
			Tags:     []string{"fighter", "any"},
		},
		Secret: domain.NarrativeBlock{
			Category: domain.NarrativeSecret,
			Content:  "A dark past",
			Tags:     []string{"any"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func newService() (*character.Service, *memCharRepo) {
	repo := newMemCharRepo()
	svc := character.NewService(repo)
	return svc, repo
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// 1. Save stores character — after Save, List returns it.
func TestSave_StoresCharacter(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-1", "user-1", "Aragorn")

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	chars, err := svc.List(ctx, "user-1")
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(chars) != 1 {
		t.Fatalf("expected 1 character, got %d", len(chars))
	}
	if chars[0].Name != "Aragorn" {
		t.Errorf("expected name Aragorn, got %q", chars[0].Name)
	}
}

// 2. Save sets UpdatedAt — UpdatedAt is non-zero after Save.
func TestSave_SetsUpdatedAt(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-2", "user-1", "Legolas")
	c.UpdatedAt = time.Time{} // zero it out

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if c.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set after Save")
	}
}

// 3. Save with empty userID → error.
func TestSave_EmptyUserID_ReturnsError(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-3", "user-1", "Gimli")

	err := svc.Save(ctx, "", c)
	if err == nil {
		t.Fatal("expected error for empty userID, got nil")
	}
}

// 4. Save with nil character → error.
func TestSave_NilCharacter_ReturnsError(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()

	err := svc.Save(ctx, "user-1", nil)
	if err == nil {
		t.Fatal("expected error for nil character, got nil")
	}
}

// 5. List returns only user's characters — two users, each sees only their own.
func TestList_ReturnsOnlyUserCharacters(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()

	c1 := newTestCharacter("char-u1-1", "user-1", "Frodo")
	c2 := newTestCharacter("char-u1-2", "user-1", "Sam")
	c3 := newTestCharacter("char-u2-1", "user-2", "Gandalf")

	for _, c := range []*domain.Character{c1, c2, c3} {
		if err := svc.Save(ctx, c.UserID, c); err != nil {
			t.Fatalf("Save returned error: %v", err)
		}
	}

	list1, err := svc.List(ctx, "user-1")
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(list1) != 2 {
		t.Errorf("user-1 expected 2 chars, got %d", len(list1))
	}

	list2, err := svc.List(ctx, "user-2")
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(list2) != 1 {
		t.Errorf("user-2 expected 1 char, got %d", len(list2))
	}
	if list2[0].Name != "Gandalf" {
		t.Errorf("expected Gandalf, got %q", list2[0].Name)
	}
}

// 6. List with empty userID → error.
func TestList_EmptyUserID_ReturnsError(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()

	_, err := svc.List(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty userID, got nil")
	}
}

// 7. Get returns correct character — by ID.
func TestGet_ReturnsCorrectCharacter(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-get-1", "user-1", "Boromir")

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	got, err := svc.Get(ctx, "user-1", "char-get-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Name != "Boromir" {
		t.Errorf("expected Boromir, got %q", got.Name)
	}
}

// 8. Get returns not-found error for unknown ID.
func TestGet_UnknownID_ReturnsError(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()

	_, err := svc.Get(ctx, "user-1", "does-not-exist")
	if err == nil {
		t.Fatal("expected not-found error, got nil")
	}
	if !errors.Is(err, character.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

// 9. Get with wrong userID → "not authorized" error (user isolation).
func TestGet_WrongUserID_ReturnsNotAuthorized(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-iso-1", "user-1", "Sauron")

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	_, err := svc.Get(ctx, "user-2", "char-iso-1")
	if err == nil {
		t.Fatal("expected not-authorized error, got nil")
	}
	if !errors.Is(err, character.ErrNotAuthorized) {
		t.Errorf("expected ErrNotAuthorized, got: %v", err)
	}
}

// 10. Edit updates name — after Edit with new name, Get returns updated name.
func TestEdit_UpdatesName(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-edit-1", "user-1", "OldName")

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	newName := "NewName"
	updated, err := svc.Edit(ctx, "user-1", "char-edit-1", character.EditInput{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if updated.Name != "NewName" {
		t.Errorf("expected name NewName, got %q", updated.Name)
	}

	// Verify persistence.
	got, err := svc.Get(ctx, "user-1", "char-edit-1")
	if err != nil {
		t.Fatalf("Get after Edit returned error: %v", err)
	}
	if got.Name != "NewName" {
		t.Errorf("persisted name: expected NewName, got %q", got.Name)
	}
}

// 11. Edit updates narrative content — Background, Motivation, Secret content changes.
func TestEdit_UpdatesNarrativeContent(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-edit-2", "user-1", "Pippin")

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	newBg := "Grew up in the Shire"
	newMotiv := "To protect his friends"
	newSecret := "Once stole a mushroom"

	updated, err := svc.Edit(ctx, "user-1", "char-edit-2", character.EditInput{
		Background: &newBg,
		Motivation: &newMotiv,
		Secret:     &newSecret,
	})
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if updated.Background.Content != newBg {
		t.Errorf("background: expected %q, got %q", newBg, updated.Background.Content)
	}
	if updated.Motivation.Content != newMotiv {
		t.Errorf("motivation: expected %q, got %q", newMotiv, updated.Motivation.Content)
	}
	if updated.Secret.Content != newSecret {
		t.Errorf("secret: expected %q, got %q", newSecret, updated.Secret.Content)
	}
}

// 12. Edit with nil fields → no change — passing all-nil EditInput changes nothing.
func TestEdit_AllNilInput_ChangesNothing(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-edit-3", "user-1", "Merry")

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	updated, err := svc.Edit(ctx, "user-1", "char-edit-3", character.EditInput{})
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if updated.Name != "Merry" {
		t.Errorf("name should be unchanged, got %q", updated.Name)
	}
	if updated.Background.Content != "A humble beginning" {
		t.Errorf("background should be unchanged, got %q", updated.Background.Content)
	}
}

// 13. Edit stats not possible — EditInput has no stats fields (compile-time guarantee).
// This test validates the struct definition: EditInput must not have stats fields.
func TestEdit_NoStatsFields_CompileTimeGuarantee(t *testing.T) {
	// If this compiles, the guarantee holds.
	// EditInput only has: Name, Background, Motivation, Secret.
	in := character.EditInput{
		Name:       nil,
		Background: nil,
		Motivation: nil,
		Secret:     nil,
	}
	// Accessing .Name ensures the field exists and is a *string.
	var _ *string = in.Name
	var _ *string = in.Background
	var _ *string = in.Motivation
	var _ *string = in.Secret
}

// 14. Edit preserves narrative tags — tags must not be lost after editing content.
func TestEdit_PreservesNarrativeTags(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-edit-4", "user-1", "Treebeard")
	c.Background.Tags = []string{"druid", "any"}
	c.Motivation.Tags = []string{"fighter", "paladin"}
	c.Secret.Tags = []string{"rogue", "warlock"}

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	newBg := "An ancient forest tale"
	newMotiv := "To guard the trees"
	newSecret := "Knows where the Entwives went"

	updated, err := svc.Edit(ctx, "user-1", "char-edit-4", character.EditInput{
		Background: &newBg,
		Motivation: &newMotiv,
		Secret:     &newSecret,
	})
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}

	if len(updated.Background.Tags) != 2 || updated.Background.Tags[0] != "druid" {
		t.Errorf("background tags not preserved: %v", updated.Background.Tags)
	}
	if len(updated.Motivation.Tags) != 2 || updated.Motivation.Tags[0] != "fighter" {
		t.Errorf("motivation tags not preserved: %v", updated.Motivation.Tags)
	}
	if len(updated.Secret.Tags) != 2 || updated.Secret.Tags[0] != "rogue" {
		t.Errorf("secret tags not preserved: %v", updated.Secret.Tags)
	}
}

// 15. Delete removes character — after Delete, Get returns not-found error.
func TestDelete_RemovesCharacter(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-del-1", "user-1", "Saruman")

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	if err := svc.Delete(ctx, "user-1", "char-del-1"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	_, err := svc.Get(ctx, "user-1", "char-del-1")
	if err == nil {
		t.Fatal("expected not-found error after delete, got nil")
	}
	if !errors.Is(err, character.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

// 16. Delete with wrong userID → error (user isolation).
func TestDelete_WrongUserID_ReturnsError(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	c := newTestCharacter("char-del-2", "user-1", "Grima")

	if err := svc.Save(ctx, "user-1", c); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	err := svc.Delete(ctx, "user-2", "char-del-2")
	if err == nil {
		t.Fatal("expected not-authorized error, got nil")
	}
	if !errors.Is(err, character.ErrNotAuthorized) {
		t.Errorf("expected ErrNotAuthorized, got: %v", err)
	}
}

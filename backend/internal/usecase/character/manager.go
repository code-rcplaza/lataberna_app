package character

import (
	"context"
	"errors"
	"fmt"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
)

// Sentinel errors for user-facing conditions.
var (
	ErrNotFound      = errors.New("character not found")
	ErrNotAuthorized = errors.New("not authorized")
)

// EditInput — only name and narrative content can be edited directly.
// Stats are immutable via this service; use character.Regenerate() for that.
type EditInput struct {
	Name       *string // nil = no change
	Background *string // nil = no change (content only, tags preserved)
	Motivation *string // nil = no change
	Secret     *string // nil = no change
}

// Service handles character persistence operations.
// It enforces user ownership on every operation.
type Service struct {
	repo ports.CharacterRepository
}

// NewService constructs a Service with the given repository.
func NewService(repo ports.CharacterRepository) *Service {
	return &Service{repo: repo}
}

// Save persists a character and associates it with userID.
// It sets a generated ID if the character has none, and stamps UpdatedAt.
func (s *Service) Save(ctx context.Context, userID string, c *domain.Character) error {
	if userID == "" {
		return fmt.Errorf("character.Service.Save: userID is required")
	}
	if c == nil {
		return fmt.Errorf("character.Service.Save: character is required")
	}
	if c.ID == "" {
		c.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	c.UserID = userID
	c.UpdatedAt = time.Now()
	return s.repo.Save(ctx, c)
}

// List returns all characters belonging to userID.
func (s *Service) List(ctx context.Context, userID string) ([]*domain.Character, error) {
	if userID == "" {
		return nil, fmt.Errorf("character.Service.List: userID is required")
	}
	return s.repo.FindByUserID(ctx, userID)
}

// Get returns the character with characterID, verifying it belongs to userID.
func (s *Service) Get(ctx context.Context, userID, characterID string) (*domain.Character, error) {
	if userID == "" {
		return nil, fmt.Errorf("character.Service.Get: userID is required")
	}
	if characterID == "" {
		return nil, fmt.Errorf("character.Service.Get: characterID is required")
	}
	c, err := s.repo.FindByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("character.Service.Get: %w", err)
	}
	if c == nil {
		return nil, ErrNotFound
	}
	if c.UserID != userID {
		return nil, ErrNotAuthorized
	}
	return c, nil
}

// Edit applies the non-nil fields of in to the character, then persists it.
// Only name and narrative content can be changed — tags and stats are preserved.
func (s *Service) Edit(ctx context.Context, userID, characterID string, in EditInput) (*domain.Character, error) {
	c, err := s.Get(ctx, userID, characterID)
	if err != nil {
		return nil, fmt.Errorf("character.Service.Edit: %w", err)
	}

	if in.Name != nil {
		c.Name = *in.Name
	}
	if in.Background != nil {
		c.Background.Content = *in.Background
	}
	if in.Motivation != nil {
		c.Motivation.Content = *in.Motivation
	}
	if in.Secret != nil {
		c.Secret.Content = *in.Secret
	}

	c.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("character.Service.Edit: %w", err)
	}
	return c, nil
}

// Delete removes the character with characterID, verifying ownership first.
func (s *Service) Delete(ctx context.Context, userID, characterID string) error {
	if _, err := s.Get(ctx, userID, characterID); err != nil {
		return fmt.Errorf("character.Service.Delete: %w", err)
	}
	return s.repo.Delete(ctx, characterID)
}

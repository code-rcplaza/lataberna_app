package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
)

// Service handles magic link authentication.
type Service struct {
	users    ports.UserRepository
	sessions ports.SessionRepository
	tokens   ports.TokenRepository
	mailer   ports.Mailer
	linkBase string // e.g. "https://app.example.com/auth/verify"
}

func NewService(
	users ports.UserRepository,
	sessions ports.SessionRepository,
	tokens ports.TokenRepository,
	mailer ports.Mailer,
	linkBase string,
) *Service {
	return &Service{
		users:    users,
		sessions: sessions,
		tokens:   tokens,
		mailer:   mailer,
		linkBase: linkBase,
	}
}

// generateID generates a cryptographically random base64url-encoded ID
// from 32 random bytes (256 bits of entropy).
func generateID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// hashToken returns the SHA-256 hex digest of the raw token.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// RequestMagicLink generates a magic link token and sends it to the given email.
// The raw token is NEVER stored — only its SHA-256 hash is persisted.
func (s *Service) RequestMagicLink(ctx context.Context, email string) error {
	if email == "" {
		return errors.New("email must not be empty")
	}

	rawToken, err := generateID()
	if err != nil {
		return err
	}

	hashed := hashToken(rawToken)

	tokenID, err := generateID()
	if err != nil {
		return err
	}

	now := time.Now()
	t := &domain.MagicLinkToken{
		ID:          tokenID,
		HashedToken: hashed,
		Email:       email,
		ExpiresAt:   now.Add(15 * time.Minute),
		CreatedAt:   now,
	}

	if err := s.tokens.Save(ctx, t); err != nil {
		return err
	}

	magicLink := s.linkBase + "?token=" + rawToken
	return s.mailer.SendMagicLink(ctx, email, magicLink)
}

// VerifyMagicLink validates the raw token and, if valid, creates an authenticated session.
// The token is marked as used immediately to prevent replay attacks.
func (s *Service) VerifyMagicLink(ctx context.Context, rawToken string) (*domain.Session, *domain.User, error) {
	hashed := hashToken(rawToken)

	t, err := s.tokens.FindByHashedToken(ctx, hashed)
	if err != nil {
		return nil, nil, err
	}
	if t == nil {
		return nil, nil, errors.New("invalid token")
	}
	if t.IsExpired() {
		return nil, nil, errors.New("token expired")
	}
	if t.IsUsed() {
		return nil, nil, errors.New("token already used")
	}

	usedAt := time.Now()
	if err := s.tokens.MarkUsed(ctx, t.ID, usedAt); err != nil {
		return nil, nil, err
	}

	user, err := s.findOrCreateUser(ctx, t.Email)
	if err != nil {
		return nil, nil, err
	}

	sessionID, err := generateID()
	if err != nil {
		return nil, nil, err
	}

	now := time.Now()
	session := &domain.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: now.Add(30 * 24 * time.Hour),
		CreatedAt: now,
	}

	if err := s.sessions.Save(ctx, session); err != nil {
		return nil, nil, err
	}

	return session, user, nil
}

// ValidateSession verifies that the given session ID is valid and not expired,
// then returns the associated user.
func (s *Service) ValidateSession(ctx context.Context, sessionID string) (*domain.User, error) {
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}
	if session.IsExpired() {
		return nil, errors.New("session expired")
	}

	user, err := s.users.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Logout invalidates a session. After this call, the session ID is no longer valid.
func (s *Service) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("auth: sessionID is required")
	}
	return s.sessions.Delete(ctx, sessionID)
}

// findOrCreateUser retrieves an existing user by email or creates a new one.
func (s *Service) findOrCreateUser(ctx context.Context, email string) (*domain.User, error) {
	existing, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	userID, err := generateID()
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		ID:        userID,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := s.users.Save(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

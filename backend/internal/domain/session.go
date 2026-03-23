package domain

import "time"

// Session represents an authenticated user session.
// Sessions are created after a successful magic link verification.
type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

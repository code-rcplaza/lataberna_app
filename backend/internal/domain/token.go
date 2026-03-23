package domain

import "time"

// MagicLinkToken represents a single-use authentication token.
// The raw token is NEVER stored — only its SHA-256 hash.
type MagicLinkToken struct {
	ID          string
	HashedToken string     // SHA-256 hash of the raw token
	Email       string
	ExpiresAt   time.Time
	UsedAt      *time.Time // nil = not yet used
	CreatedAt   time.Time
}

func (t *MagicLinkToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t *MagicLinkToken) IsUsed() bool {
	return t.UsedAt != nil
}

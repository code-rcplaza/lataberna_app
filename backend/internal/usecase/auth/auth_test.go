package auth_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/usecase/auth"
)

// ---------------------------------------------------------------------------
// In-memory test doubles
// ---------------------------------------------------------------------------

type memUserRepo struct {
	mu    sync.Mutex
	users map[string]*domain.User // keyed by ID
}

func newMemUserRepo() *memUserRepo {
	return &memUserRepo{users: make(map[string]*domain.User)}
}

func (r *memUserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil // not found — not an error
}

func (r *memUserRepo) FindByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (r *memUserRepo) Save(_ context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

// ---------------------------------------------------------------------------

type memSessionRepo struct {
	mu       sync.Mutex
	sessions map[string]*domain.Session
}

func newMemSessionRepo() *memSessionRepo {
	return &memSessionRepo{sessions: make(map[string]*domain.Session)}
}

func (r *memSessionRepo) Save(_ context.Context, s *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s.ID] = s
	return nil
}

func (r *memSessionRepo) FindByID(_ context.Context, id string) (*domain.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.sessions[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *memSessionRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, id)
	return nil
}

// ---------------------------------------------------------------------------

type memTokenRepo struct {
	mu     sync.Mutex
	tokens map[string]*domain.MagicLinkToken // keyed by HashedToken
}

func newMemTokenRepo() *memTokenRepo {
	return &memTokenRepo{tokens: make(map[string]*domain.MagicLinkToken)}
}

func (r *memTokenRepo) Save(_ context.Context, t *domain.MagicLinkToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[t.HashedToken] = t
	return nil
}

func (r *memTokenRepo) FindByHashedToken(_ context.Context, hashed string) (*domain.MagicLinkToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tokens[hashed]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *memTokenRepo) MarkUsed(_ context.Context, id string, usedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.tokens {
		if t.ID == id {
			t.UsedAt = &usedAt
			return nil
		}
	}
	return errors.New("token not found")
}

// allTokens returns a snapshot of all stored tokens (for assertions).
func (r *memTokenRepo) allTokens() []*domain.MagicLinkToken {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*domain.MagicLinkToken, 0, len(r.tokens))
	for _, t := range r.tokens {
		out = append(out, t)
	}
	return out
}

// ---------------------------------------------------------------------------

type mockMailer struct {
	mu       sync.Mutex
	sent     []mailCall
	sendFunc func(ctx context.Context, email, link string) error // optional override
}

type mailCall struct {
	email string
	link  string
}

func (m *mockMailer) SendMagicLink(ctx context.Context, email, link string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.sendFunc != nil {
		return m.sendFunc(ctx, email, link)
	}
	m.sent = append(m.sent, mailCall{email: email, link: link})
	return nil
}

func (m *mockMailer) calls() []mailCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]mailCall, len(m.sent))
	copy(out, m.sent)
	return out
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const linkBase = "https://app.example.com/auth/verify"

func newService(
	users *memUserRepo,
	sessions *memSessionRepo,
	tokens *memTokenRepo,
	mailer *mockMailer,
) *auth.Service {
	return auth.NewService(users, sessions, tokens, mailer, linkBase)
}

// requestAndExtractToken calls RequestMagicLink and extracts the raw token
// from the mailer's captured link, so tests can call VerifyMagicLink directly.
func requestAndExtractToken(t *testing.T, svc *auth.Service, mailer *mockMailer, email string) string {
	t.Helper()
	if err := svc.RequestMagicLink(context.Background(), email); err != nil {
		t.Fatalf("RequestMagicLink failed: %v", err)
	}
	calls := mailer.calls()
	if len(calls) == 0 {
		t.Fatal("mailer was not called")
	}
	link := calls[len(calls)-1].link
	// link is linkBase + "?token=" + rawToken
	parts := strings.SplitN(link, "?token=", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected link format: %s", link)
	}
	return parts[1]
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// 1. RequestMagicLink saves a hashed (not raw) token.
func TestRequestMagicLink_SavesHashedToken(t *testing.T) {
	tokens := newMemTokenRepo()
	mailer := &mockMailer{}
	svc := newService(newMemUserRepo(), newMemSessionRepo(), tokens, mailer)

	rawToken := requestAndExtractToken(t, svc, mailer, "player@example.com")

	stored := tokens.allTokens()
	if len(stored) != 1 {
		t.Fatalf("expected 1 token stored, got %d", len(stored))
	}

	if stored[0].HashedToken == rawToken {
		t.Error("stored HashedToken must not equal the raw token — it must be hashed")
	}
	if stored[0].HashedToken == "" {
		t.Error("stored HashedToken must not be empty")
	}
}

// 2. RequestMagicLink calls mailer with correct email.
func TestRequestMagicLink_CallsMailer(t *testing.T) {
	mailer := &mockMailer{}
	svc := newService(newMemUserRepo(), newMemSessionRepo(), newMemTokenRepo(), mailer)

	if err := svc.RequestMagicLink(context.Background(), "hero@forge.rpg"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	calls := mailer.calls()
	if len(calls) != 1 {
		t.Fatalf("expected 1 mail call, got %d", len(calls))
	}
	if calls[0].email != "hero@forge.rpg" {
		t.Errorf("expected email %q, got %q", "hero@forge.rpg", calls[0].email)
	}
	if !strings.HasPrefix(calls[0].link, linkBase) {
		t.Errorf("magic link %q does not start with linkBase %q", calls[0].link, linkBase)
	}
}

// 3. RequestMagicLink with empty email returns error.
func TestRequestMagicLink_EmptyEmail_Error(t *testing.T) {
	svc := newService(newMemUserRepo(), newMemSessionRepo(), newMemTokenRepo(), &mockMailer{})

	err := svc.RequestMagicLink(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty email, got nil")
	}
}

// 4. VerifyMagicLink with valid token creates session and user.
func TestVerifyMagicLink_ValidToken_CreatesSessionAndUser(t *testing.T) {
	users := newMemUserRepo()
	sessions := newMemSessionRepo()
	mailer := &mockMailer{}
	svc := newService(users, sessions, newMemTokenRepo(), mailer)

	rawToken := requestAndExtractToken(t, svc, mailer, "adventurer@forge.rpg")

	session, user, err := svc.VerifyMagicLink(context.Background(), rawToken)
	if err != nil {
		t.Fatalf("VerifyMagicLink failed: %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if session.ID == "" {
		t.Error("session ID must not be empty")
	}
	if user.Email != "adventurer@forge.rpg" {
		t.Errorf("user email = %q, want %q", user.Email, "adventurer@forge.rpg")
	}
	if session.UserID != user.ID {
		t.Errorf("session.UserID %q != user.ID %q", session.UserID, user.ID)
	}
}

// 5. VerifyMagicLink with invalid (unknown) token returns error.
func TestVerifyMagicLink_InvalidToken_Error(t *testing.T) {
	svc := newService(newMemUserRepo(), newMemSessionRepo(), newMemTokenRepo(), &mockMailer{})

	_, _, err := svc.VerifyMagicLink(context.Background(), "this-token-was-never-issued")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

// 6. VerifyMagicLink with expired token returns error.
func TestVerifyMagicLink_ExpiredToken_Error(t *testing.T) {
	tokens := newMemTokenRepo()
	mailer := &mockMailer{}
	svc := newService(newMemUserRepo(), newMemSessionRepo(), tokens, mailer)

	rawToken := requestAndExtractToken(t, svc, mailer, "late@forge.rpg")

	// Manually expire the stored token.
	stored := tokens.allTokens()
	if len(stored) != 1 {
		t.Fatalf("expected 1 token, got %d", len(stored))
	}
	past := time.Now().Add(-1 * time.Hour)
	stored[0].ExpiresAt = past

	_, _, err := svc.VerifyMagicLink(context.Background(), rawToken)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

// 7. VerifyMagicLink with used token (used twice) returns error on second use.
func TestVerifyMagicLink_UsedToken_Error(t *testing.T) {
	mailer := &mockMailer{}
	svc := newService(newMemUserRepo(), newMemSessionRepo(), newMemTokenRepo(), mailer)

	rawToken := requestAndExtractToken(t, svc, mailer, "reuse@forge.rpg")

	// First use — must succeed.
	if _, _, err := svc.VerifyMagicLink(context.Background(), rawToken); err != nil {
		t.Fatalf("first VerifyMagicLink failed: %v", err)
	}

	// Second use — must fail.
	_, _, err := svc.VerifyMagicLink(context.Background(), rawToken)
	if err == nil {
		t.Fatal("expected error on second use of token, got nil")
	}
}

// 8. VerifyMagicLink marks the token as used (UsedAt is non-nil).
func TestVerifyMagicLink_MarksTokenUsed(t *testing.T) {
	tokens := newMemTokenRepo()
	mailer := &mockMailer{}
	svc := newService(newMemUserRepo(), newMemSessionRepo(), tokens, mailer)

	rawToken := requestAndExtractToken(t, svc, mailer, "mark@forge.rpg")

	if _, _, err := svc.VerifyMagicLink(context.Background(), rawToken); err != nil {
		t.Fatalf("VerifyMagicLink failed: %v", err)
	}

	stored := tokens.allTokens()
	if len(stored) != 1 {
		t.Fatalf("expected 1 token, got %d", len(stored))
	}
	if stored[0].UsedAt == nil {
		t.Error("expected UsedAt to be non-nil after verification")
	}
}

// 9. VerifyMagicLink reuses existing user when same email is verified again.
func TestVerifyMagicLink_ExistingUser_Reused(t *testing.T) {
	users := newMemUserRepo()
	sessions := newMemSessionRepo()

	mailer1 := &mockMailer{}
	svc1 := newService(users, sessions, newMemTokenRepo(), mailer1)
	raw1 := requestAndExtractToken(t, svc1, mailer1, "returning@forge.rpg")
	_, user1, err := svc1.VerifyMagicLink(context.Background(), raw1)
	if err != nil {
		t.Fatalf("first verify failed: %v", err)
	}

	mailer2 := &mockMailer{}
	svc2 := newService(users, sessions, newMemTokenRepo(), mailer2)
	raw2 := requestAndExtractToken(t, svc2, mailer2, "returning@forge.rpg")
	_, user2, err := svc2.VerifyMagicLink(context.Background(), raw2)
	if err != nil {
		t.Fatalf("second verify failed: %v", err)
	}

	if user1.ID != user2.ID {
		t.Errorf("expected same user ID for same email: got %q and %q", user1.ID, user2.ID)
	}
}

// 10. ValidateSession with valid session returns user.
func TestValidateSession_ValidSession_ReturnsUser(t *testing.T) {
	users := newMemUserRepo()
	sessions := newMemSessionRepo()
	mailer := &mockMailer{}
	svc := newService(users, sessions, newMemTokenRepo(), mailer)

	rawToken := requestAndExtractToken(t, svc, mailer, "valid@forge.rpg")
	session, _, err := svc.VerifyMagicLink(context.Background(), rawToken)
	if err != nil {
		t.Fatalf("VerifyMagicLink failed: %v", err)
	}

	user, err := svc.ValidateSession(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("ValidateSession failed: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Email != "valid@forge.rpg" {
		t.Errorf("user email = %q, want %q", user.Email, "valid@forge.rpg")
	}
}

// 11. ValidateSession with expired session returns error.
func TestValidateSession_ExpiredSession_Error(t *testing.T) {
	sessions := newMemSessionRepo()
	mailer := &mockMailer{}
	svc := newService(newMemUserRepo(), sessions, newMemTokenRepo(), mailer)

	rawToken := requestAndExtractToken(t, svc, mailer, "expired@forge.rpg")
	session, _, err := svc.VerifyMagicLink(context.Background(), rawToken)
	if err != nil {
		t.Fatalf("VerifyMagicLink failed: %v", err)
	}

	// Manually expire the session.
	sessions.mu.Lock()
	sessions.sessions[session.ID].ExpiresAt = time.Now().Add(-1 * time.Hour)
	sessions.mu.Unlock()

	_, err = svc.ValidateSession(context.Background(), session.ID)
	if err == nil {
		t.Fatal("expected error for expired session, got nil")
	}
}

// 12. ValidateSession with unknown sessionID returns error.
func TestValidateSession_UnknownSession_Error(t *testing.T) {
	svc := newService(newMemUserRepo(), newMemSessionRepo(), newMemTokenRepo(), &mockMailer{})

	_, err := svc.ValidateSession(context.Background(), "session-that-does-not-exist")
	if err == nil {
		t.Fatal("expected error for unknown session, got nil")
	}
}

// 13. Raw token is never stored — HashedToken differs from the raw token used in Verify.
func TestRawToken_NeverStored(t *testing.T) {
	tokens := newMemTokenRepo()
	mailer := &mockMailer{}
	svc := newService(newMemUserRepo(), newMemSessionRepo(), tokens, mailer)

	rawToken := requestAndExtractToken(t, svc, mailer, "security@forge.rpg")

	stored := tokens.allTokens()
	if len(stored) != 1 {
		t.Fatalf("expected 1 token, got %d", len(stored))
	}

	if stored[0].HashedToken == rawToken {
		t.Error("SECURITY VIOLATION: raw token is stored directly — must store SHA-256 hash instead")
	}
}

// 14. Logout with valid sessionID deletes the session (ValidateSession returns error afterwards).
func TestLogout_ValidSession_DeletesSession(t *testing.T) {
	users := newMemUserRepo()
	sessions := newMemSessionRepo()
	mailer := &mockMailer{}
	svc := newService(users, sessions, newMemTokenRepo(), mailer)

	rawToken := requestAndExtractToken(t, svc, mailer, "logout@forge.rpg")
	session, _, err := svc.VerifyMagicLink(context.Background(), rawToken)
	if err != nil {
		t.Fatalf("VerifyMagicLink failed: %v", err)
	}

	if err := svc.Logout(context.Background(), session.ID); err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	_, err = svc.ValidateSession(context.Background(), session.ID)
	if err == nil {
		t.Fatal("expected ValidateSession to return error after Logout, got nil")
	}
}

// 15. Logout with empty sessionID returns error.
func TestLogout_EmptySessionID_Error(t *testing.T) {
	svc := newService(newMemUserRepo(), newMemSessionRepo(), newMemTokenRepo(), &mockMailer{})

	err := svc.Logout(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty sessionID, got nil")
	}
}

// 16. Token has 15-minute expiry (within [14min, 16min] from now).
func TestRequestMagicLink_TokenExpiry_15Minutes(t *testing.T) {
	tokens := newMemTokenRepo()
	svc := newService(newMemUserRepo(), newMemSessionRepo(), tokens, &mockMailer{})

	before := time.Now()
	if err := svc.RequestMagicLink(context.Background(), "timer@forge.rpg"); err != nil {
		t.Fatalf("RequestMagicLink failed: %v", err)
	}
	after := time.Now()

	stored := tokens.allTokens()
	if len(stored) != 1 {
		t.Fatalf("expected 1 token, got %d", len(stored))
	}

	minExpiry := before.Add(14 * time.Minute)
	maxExpiry := after.Add(16 * time.Minute)

	if stored[0].ExpiresAt.Before(minExpiry) {
		t.Errorf("token expiry %v is too early (< 14 min from now)", stored[0].ExpiresAt)
	}
	if stored[0].ExpiresAt.After(maxExpiry) {
		t.Errorf("token expiry %v is too late (> 16 min from now)", stored[0].ExpiresAt)
	}
}

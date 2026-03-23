package graphql

import (
	"context"
	"net/http"
	"strings"

	"forge-rpg/internal/usecase/auth"
)

// contextKey is a private type to avoid key collisions in context values.
type contextKey string

// UserContextKey is the context key under which the authenticated *domain.User is stored.
const UserContextKey contextKey = "user"

// SessionContextKey is the context key under which the session ID string is stored.
// Resolvers that need to delete the session (e.g. Logout) read from this key.
const SessionContextKey contextKey = "sessionID"

// AuthMiddleware reads the Authorization: Bearer <sessionId> header,
// validates the session via authSvc.ValidateSession, and injects both the
// *domain.User and the raw session ID into the request context.
// Unauthenticated requests pass through — individual resolvers decide
// whether authentication is required.
func AuthMiddleware(authSvc *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				next.ServeHTTP(w, r)
				return
			}

			sessionID := parts[1]
			user, err := authSvc.ValidateSession(r.Context(), sessionID)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			ctx = context.WithValue(ctx, SessionContextKey, sessionID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

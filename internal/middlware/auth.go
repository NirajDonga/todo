package middlware

import (
	"context"
	"net/http"
	"strings"

	"github.com/NirajDonga/todo/internal/auth"
	"github.com/google/uuid"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// NewAuthMiddleware returns an HTTP middleware that validates JWTs using the provided AuthService.
func NewAuthMiddleware(svc auth.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			claims, err := svc.ValidateToken(r.Context(), token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// ensure user id is a valid UUID (optional sanity check)
			if claims.UserID == "" {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
			if _, err := uuid.Parse(claims.UserID); err != nil {
				http.Error(w, "invalid user id in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext extracts the user ID set by the auth middleware.
func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

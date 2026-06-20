package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/muhammadjoni/mfwebapp/pkg/jwt"
	"github.com/muhammadjoni/mfwebapp/pkg/response"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
)

type AuthMiddleware struct {
	jwtMgr *jwt.Manager
}

func NewAuthMiddleware(jwtMgr *jwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{jwtMgr: jwtMgr}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearer(r)
		if token == "" {
			response.Unauthorized(w)
			return
		}

		claims, err := m.jwtMgr.ParseAccess(token)
		if err != nil {
			response.Unauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextKeyRole).(string)
			if !ok {
				response.Unauthorized(w)
				return
			}
			for _, allowed := range roles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.Forbidden(w)
		})
	}
}

func extractBearer(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return parts[1]
	}
	return ""
}

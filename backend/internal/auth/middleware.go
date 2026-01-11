// auth/middleware.go
package auth

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func (h *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := ExtractTokenFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := h.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

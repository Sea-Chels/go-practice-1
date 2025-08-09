package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/Sea-Chels/go-practice-1/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			utils.ErrorResponse(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context while preserving existing context values
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		
		// Call the next handler with the updated request
		next(w, r.WithContext(ctx))
	}
}

func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	return claims, ok
}
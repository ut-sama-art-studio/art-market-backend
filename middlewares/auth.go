package middlewares

import (
	"context"
	"net/http"

	"github.com/ut-sama-art-studio/art-market-backend/utils/jwt"
)

type contextKey string

// Extracts userID and stores it in context if auth token is provided in request header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		bearerPrefix := "Bearer "
		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			http.Error(w, "Invalid token format", http.StatusForbidden)
			return
		}

		tokenString := authHeader[len(bearerPrefix):]
		claims, err := jwt.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid authentication token", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKey("userID"), claims["user_id"])

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// Returns user's ID from context
func ContextUserID(ctx context.Context) string {
	userID, _ := ctx.Value(contextKey("userID")).(string)
	return userID
}

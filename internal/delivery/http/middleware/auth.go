package middleware

import (
	"net/http"
	"strings"

	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
	infraauth "github.com/hajimohammadinet/dabir/internal/infrastructure/auth"
)

func AuthMiddleware(jwtService *infraauth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authorization header is required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, http.StatusUnauthorized, "INVALID_AUTH_HEADER", "authorization header must be Bearer token")
				return
			}

			claims, err := jwtService.ParseAccessToken(parts[1])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "INVALID_TOKEN", "invalid or expired token")
				return
			}

			authUser := AuthUser{
				ID:       claims.UserID,
				Username: claims.Username,
				Role:     claims.Role,
			}

			next.ServeHTTP(w, r.WithContext(WithAuthUser(r.Context(), authUser)))
		})
	}
}

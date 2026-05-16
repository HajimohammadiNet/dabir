package middleware

import (
	"net/http"

	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

func RequireRoles(allowedRoles ...user.Role) func(http.Handler) http.Handler {
	allowed := make(map[user.Role]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authUser, ok := GetAuthUser(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
				return
			}

			if _, ok := allowed[authUser.Role]; !ok {
				response.Error(w, http.StatusForbidden, "FORBIDDEN", "you do not have permission to perform this action")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

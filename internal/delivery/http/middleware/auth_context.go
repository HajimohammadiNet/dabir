package middleware

import (
	"context"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

type contextKey string

const authUserKey contextKey = "auth_user"

type AuthUser struct {
	ID       string
	Username string
	Role     user.Role
}

func WithAuthUser(ctx context.Context, authUser AuthUser) context.Context {
	return context.WithValue(ctx, authUserKey, authUser)
}

func GetAuthUser(ctx context.Context) (AuthUser, bool) {
	authUser, ok := ctx.Value(authUserKey).(AuthUser)
	return authUser, ok
}

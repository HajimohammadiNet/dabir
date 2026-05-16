package user

import "context"

type Repository interface {
	CountSuperUsers(ctx context.Context) (int, error)
	Create(ctx context.Context, u *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

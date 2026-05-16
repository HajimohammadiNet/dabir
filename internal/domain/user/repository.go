package user

import "context"

type Repository interface {
	CountSuperUsers(ctx context.Context) (int, error)

	Create(ctx context.Context, u *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)

	List(ctx context.Context, filter ListFilter) ([]User, int, error)
	Update(ctx context.Context, u *User) error
	SetActive(ctx context.Context, id string, isActive bool) error
}

type ListFilter struct {
	Page     int
	PageSize int
	Search   string
	Role     *Role
	IsActive *bool
}

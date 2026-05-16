package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

var ErrUserNotFound = errors.New("user not found")

type GetUserUseCase struct {
	userRepo user.Repository
}

func NewGetUserUseCase(userRepo user.Repository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*UserDTO, error) {
	u, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if u == nil {
		return nil, ErrUserNotFound
	}

	dto := ToUserDTO(*u)

	return &dto, nil
}

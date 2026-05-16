package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

type SetUserActiveUseCase struct {
	userRepo user.Repository
}

func NewSetUserActiveUseCase(userRepo user.Repository) *SetUserActiveUseCase {
	return &SetUserActiveUseCase{
		userRepo: userRepo,
	}
}

func (uc *SetUserActiveUseCase) Execute(ctx context.Context, id string, isActive bool) error {
	if id == "" {
		return errors.New("user id is required")
	}

	u, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if u == nil {
		return ErrUserNotFound
	}

	if err := uc.userRepo.SetActive(ctx, id, isActive); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

var ErrWeakPassword = errors.New("password must be at least 8 characters")

type ResetPasswordUseCase struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
}

type ResetPasswordInput struct {
	UserID      string `json:"-"`
	NewPassword string `json:"new_password"`
}

func NewResetPasswordUseCase(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (uc *ResetPasswordUseCase) Execute(ctx context.Context, input ResetPasswordInput) error {
	if input.UserID == "" {
		return errors.New("user id is required")
	}

	if len(input.NewPassword) < 8 {
		return ErrWeakPassword
	}

	u, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if u == nil {
		return ErrUserNotFound
	}

	passwordHash, err := uc.passwordHasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := uc.userRepo.UpdatePassword(ctx, input.UserID, passwordHash); err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}

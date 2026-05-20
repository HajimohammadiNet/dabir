package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

var (
	ErrInvalidCurrentPassword = errors.New("invalid current password")
	ErrWeakPassword           = errors.New("new password must be at least 8 characters")
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type ChangePasswordUseCase struct {
	userRepo         user.Repository
	passwordComparer PasswordComparer
	passwordHasher   PasswordHasher
}

type ChangePasswordInput struct {
	UserID          string `json:"-"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func NewChangePasswordUseCase(
	userRepo user.Repository,
	passwordComparer PasswordComparer,
	passwordHasher PasswordHasher,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo:         userRepo,
		passwordComparer: passwordComparer,
		passwordHasher:   passwordHasher,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) error {
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

	if err := uc.passwordComparer.Compare(u.PasswordHash, input.CurrentPassword); err != nil {
		return ErrInvalidCurrentPassword
	}

	newHash, err := uc.passwordHasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	if err := uc.userRepo.UpdatePassword(ctx, input.UserID, newHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

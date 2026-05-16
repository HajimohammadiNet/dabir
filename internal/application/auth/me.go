package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

var ErrUserNotFound = errors.New("user not found")

type MeUseCase struct {
	userRepo user.Repository
}

type MeOutput struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Role     user.Role `json:"role"`
	IsActive bool      `json:"is_active"`
}

func NewMeUseCase(userRepo user.Repository) *MeUseCase {
	return &MeUseCase{
		userRepo: userRepo,
	}
}

func (uc *MeUseCase) Execute(ctx context.Context, userID string) (*MeOutput, error) {
	u, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if u == nil {
		return nil, ErrUserNotFound
	}

	return &MeOutput{
		ID:       u.ID,
		Username: u.Username,
		FullName: u.FullName,
		Role:     u.Role,
		IsActive: u.IsActive,
	}, nil
}

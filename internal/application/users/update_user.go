package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

type UpdateUserUseCase struct {
	userRepo user.Repository
}

type UpdateUserInput struct {
	ID       string    `json:"-"`
	FullName string    `json:"full_name"`
	Role     user.Role `json:"role"`
	IsActive *bool     `json:"is_active,omitempty"`
}

func NewUpdateUserUseCase(userRepo user.Repository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*UserDTO, error) {
	input.FullName = strings.TrimSpace(input.FullName)

	if input.ID == "" {
		return nil, errors.New("user id is required")
	}

	if input.FullName == "" {
		return nil, errors.New("full name is required")
	}

	if !isValidRole(input.Role) {
		return nil, ErrInvalidRole
	}

	u, err := uc.userRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if u == nil {
		return nil, ErrUserNotFound
	}

	u.FullName = input.FullName
	u.Role = input.Role

	if input.IsActive != nil {
		u.IsActive = *input.IsActive
	}

	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	dto := ToUserDTO(*u)

	return &dto, nil
}

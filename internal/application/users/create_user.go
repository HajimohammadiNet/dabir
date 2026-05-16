package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidRole           = errors.New("invalid role")
	ErrInvalidInput          = errors.New("invalid input")
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type CreateUserUseCase struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
}

type CreateUserInput struct {
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Password string    `json:"password"`
	Role     user.Role `json:"role"`
}

func NewCreateUserUseCase(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*UserDTO, error) {
	input.Username = strings.TrimSpace(input.Username)
	input.FullName = strings.TrimSpace(input.FullName)

	if err := validateCreateUserInput(input); err != nil {
		return nil, err
	}

	existingUser, err := uc.userRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}

	if existingUser != nil {
		return nil, ErrUsernameAlreadyExists
	}

	passwordHash, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := &user.User{
		Username:     input.Username,
		FullName:     input.FullName,
		PasswordHash: passwordHash,
		Role:         input.Role,
		IsActive:     true,
	}

	if err := uc.userRepo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	dto := ToUserDTO(*u)

	return &dto, nil
}

func validateCreateUserInput(input CreateUserInput) error {
	if input.Username == "" {
		return errors.New("username is required")
	}

	if len(input.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if input.FullName == "" {
		return errors.New("full name is required")
	}

	if len(input.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if !isValidRole(input.Role) {
		return ErrInvalidRole
	}

	return nil
}

func isValidRole(role user.Role) bool {
	switch role {
	case user.RoleSuperUser, user.RoleEditor, user.RoleReadonly:
		return true
	default:
		return false
	}
}

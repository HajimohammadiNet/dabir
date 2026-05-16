package setup

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

var ErrAlreadyInitialized = errors.New("application is already initialized")

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type InitializeUseCase struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
}

type InitializeInput struct {
	OrganizationName string                 `json:"organization_name"`
	SuperUser        InitializeSuperUser    `json:"superuser"`
	LetterConfig     InitializeLetterConfig `json:"letter_config"`
}

type InitializeSuperUser struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type InitializeLetterConfig struct {
	NumberPrefix  string `json:"number_prefix"`
	NumberPadding int    `json:"number_padding"`
}

type InitializeOutput struct {
	Initialized bool   `json:"initialized"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
}

func NewInitializeUseCase(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
) *InitializeUseCase {
	return &InitializeUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (uc *InitializeUseCase) Execute(ctx context.Context, input InitializeInput) (*InitializeOutput, error) {
	count, err := uc.userRepo.CountSuperUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check setup status: %w", err)
	}

	if count > 0 {
		return nil, ErrAlreadyInitialized
	}

	if err := validateInitializeInput(input); err != nil {
		return nil, err
	}

	passwordHash, err := uc.passwordHasher.Hash(input.SuperUser.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash superuser password: %w", err)
	}

	superUser := &user.User{
		Username:     strings.TrimSpace(input.SuperUser.Username),
		FullName:     strings.TrimSpace(input.SuperUser.FullName),
		PasswordHash: passwordHash,
		Role:         user.RoleSuperUser,
		IsActive:     true,
	}

	if err := uc.userRepo.Create(ctx, superUser); err != nil {
		return nil, fmt.Errorf("failed to create superuser: %w", err)
	}

	return &InitializeOutput{
		Initialized: true,
		UserID:      superUser.ID,
		Username:    superUser.Username,
	}, nil
}

func validateInitializeInput(input InitializeInput) error {
	if strings.TrimSpace(input.SuperUser.Username) == "" {
		return errors.New("username is required")
	}

	if strings.TrimSpace(input.SuperUser.FullName) == "" {
		return errors.New("full name is required")
	}

	if len(input.SuperUser.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

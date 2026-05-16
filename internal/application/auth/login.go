package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInactiveUser       = errors.New("user is inactive")
)

type PasswordComparer interface {
	Compare(hash string, password string) error
}

type TokenGenerator interface {
	GenerateAccessToken(u *user.User) (string, int64, error)
}

type LoginUseCase struct {
	userRepo         user.Repository
	passwordComparer PasswordComparer
	tokenGenerator   TokenGenerator
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginOutput struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        LoginUserDTO `json:"user"`
}

type LoginUserDTO struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Role     user.Role `json:"role"`
}

func NewLoginUseCase(
	userRepo user.Repository,
	passwordComparer PasswordComparer,
	tokenGenerator TokenGenerator,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:         userRepo,
		passwordComparer: passwordComparer,
		tokenGenerator:   tokenGenerator,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	username := strings.TrimSpace(input.Username)

	if username == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	u, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if u == nil {
		return nil, ErrInvalidCredentials
	}

	if !u.IsActive {
		return nil, ErrInactiveUser
	}

	if err := uc.passwordComparer.Compare(u.PasswordHash, input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, expiresIn, err := uc.tokenGenerator.GenerateAccessToken(u)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User: LoginUserDTO{
			ID:       u.ID,
			Username: u.Username,
			FullName: u.FullName,
			Role:     u.Role,
		},
	}, nil
}

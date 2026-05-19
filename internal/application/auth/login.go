package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	auditapp "github.com/hajimohammadinet/dabir/internal/application/audit"
	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
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
	auditLogger      *auditapp.Logger
}

type LoginInput struct {
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	IPAddress *string `json:"-"`
	UserAgent *string `json:"-"`
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
	auditLogger *auditapp.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:         userRepo,
		passwordComparer: passwordComparer,
		tokenGenerator:   tokenGenerator,
		auditLogger:      auditLogger,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	username := strings.TrimSpace(input.Username)

	if username == "" || input.Password == "" {
		uc.logLoginFailed(ctx, input, "empty_username_or_password")
		return nil, ErrInvalidCredentials
	}

	u, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if u == nil {
		uc.logLoginFailed(ctx, input, "user_not_found")
		return nil, ErrInvalidCredentials
	}

	if !u.IsActive {
		uc.logLoginFailed(ctx, input, "inactive_user")
		return nil, ErrInactiveUser
	}

	if err := uc.passwordComparer.Compare(u.PasswordHash, input.Password); err != nil {
		uc.logLoginFailed(ctx, input, "invalid_password")
		return nil, ErrInvalidCredentials
	}

	token, expiresIn, err := uc.tokenGenerator.GenerateAccessToken(u)
	if err != nil {
		return nil, err
	}

	uc.logLoginSuccess(ctx, input, u)

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

func (uc *LoginUseCase) logLoginFailed(ctx context.Context, input LoginInput, reason string) {
	uc.auditLogger.Log(ctx, auditapp.LogInput{
		Action:     domainaudit.ActionAuthLoginFailed,
		EntityType: "auth",
		IPAddress:  input.IPAddress,
		UserAgent:  input.UserAgent,
		NewValue: map[string]interface{}{
			"username": input.Username,
			"reason":   reason,
		},
	})
}

func (uc *LoginUseCase) logLoginSuccess(ctx context.Context, input LoginInput, u *user.User) {
	userID := u.ID

	uc.auditLogger.Log(ctx, auditapp.LogInput{
		ActorUserID: &userID,
		Action:      domainaudit.ActionAuthLoginSuccess,
		EntityType:  "auth",
		EntityID:    &userID,
		IPAddress:   input.IPAddress,
		UserAgent:   input.UserAgent,
		NewValue: map[string]interface{}{
			"username": u.Username,
			"role":     u.Role,
		},
	})
}

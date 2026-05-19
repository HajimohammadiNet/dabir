package setup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	auditapp "github.com/hajimohammadinet/dabir/internal/application/audit"
	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
	"github.com/hajimohammadinet/dabir/internal/domain/settings"
	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

var ErrAlreadyInitialized = errors.New("application is already initialized")

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type InitializeUseCase struct {
	userRepo       user.Repository
	settingsRepo   settings.Repository
	passwordHasher PasswordHasher
	auditLogger    *auditapp.Logger
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

type storedLetterConfig struct {
	NumberPrefix  string `json:"number_prefix"`
	NumberPadding int    `json:"number_padding"`
}

func NewInitializeUseCase(
	userRepo user.Repository,
	settingsRepo settings.Repository,
	passwordHasher PasswordHasher,
	auditLogger *auditapp.Logger,
) *InitializeUseCase {
	return &InitializeUseCase{
		userRepo:       userRepo,
		settingsRepo:   settingsRepo,
		passwordHasher: passwordHasher,
		auditLogger:    auditLogger,
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

	input = normalizeInitializeInput(input)

	if err := validateInitializeInput(input); err != nil {
		return nil, err
	}

	passwordHash, err := uc.passwordHasher.Hash(input.SuperUser.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash superuser password: %w", err)
	}

	superUser := &user.User{
		Username:     input.SuperUser.Username,
		FullName:     input.SuperUser.FullName,
		PasswordHash: passwordHash,
		Role:         user.RoleSuperUser,
		IsActive:     true,
	}

	if err := uc.userRepo.Create(ctx, superUser); err != nil {
		return nil, fmt.Errorf("failed to create superuser: %w", err)
	}

	if err := uc.saveInitialSettings(ctx, input); err != nil {
		return nil, err
	}

	if uc.auditLogger != nil {
		actorID := superUser.ID

		uc.auditLogger.Log(ctx, auditapp.LogInput{
			ActorUserID: &actorID,
			Action:      domainaudit.ActionSetupInitialized,
			EntityType:  "setup",
			EntityID:    nil,
			NewValue: map[string]interface{}{
				"organization_name": input.OrganizationName,
				"superuser": map[string]interface{}{
					"id":       superUser.ID,
					"username": superUser.Username,
				},
				"letter_config": input.LetterConfig,
			},
		})
	}

	return &InitializeOutput{
		Initialized: true,
		UserID:      superUser.ID,
		Username:    superUser.Username,
	}, nil
}

func (uc *InitializeUseCase) saveInitialSettings(ctx context.Context, input InitializeInput) error {
	orgValue, err := json.Marshal(input.OrganizationName)
	if err != nil {
		return fmt.Errorf("failed to marshal organization name: %w", err)
	}

	if err := uc.settingsRepo.Set(ctx, settings.KeyOrganizationName, orgValue); err != nil {
		return err
	}

	letterConfig := storedLetterConfig{
		NumberPrefix:  input.LetterConfig.NumberPrefix,
		NumberPadding: input.LetterConfig.NumberPadding,
	}

	letterConfigValue, err := json.Marshal(letterConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal letter config: %w", err)
	}

	if err := uc.settingsRepo.Set(ctx, settings.KeyLetterConfig, letterConfigValue); err != nil {
		return err
	}

	return nil
}

func normalizeInitializeInput(input InitializeInput) InitializeInput {
	input.OrganizationName = strings.TrimSpace(input.OrganizationName)
	input.SuperUser.Username = strings.TrimSpace(input.SuperUser.Username)
	input.SuperUser.FullName = strings.TrimSpace(input.SuperUser.FullName)
	input.LetterConfig.NumberPrefix = strings.TrimSpace(input.LetterConfig.NumberPrefix)

	if input.OrganizationName == "" {
		input.OrganizationName = "Dabir"
	}

	if input.LetterConfig.NumberPrefix == "" {
		input.LetterConfig.NumberPrefix = "DABIR"
	}

	if input.LetterConfig.NumberPadding == 0 {
		input.LetterConfig.NumberPadding = 6
	}

	return input
}

func validateInitializeInput(input InitializeInput) error {
	if input.OrganizationName == "" {
		return errors.New("organization name is required")
	}

	if input.SuperUser.Username == "" {
		return errors.New("username is required")
	}

	if input.SuperUser.FullName == "" {
		return errors.New("full name is required")
	}

	if len(input.SuperUser.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if input.LetterConfig.NumberPadding < 1 || input.LetterConfig.NumberPadding > 12 {
		return errors.New("number padding must be between 1 and 12")
	}

	return nil
}

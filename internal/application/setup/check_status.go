package setup

import (
	"context"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

type CheckStatusUseCase struct {
	userRepo user.Repository
}

type CheckStatusOutput struct {
	Initialized bool `json:"initialized"`
	SetupNeeded bool `json:"setup_needed"`
}

func NewCheckStatusUseCase(userRepo user.Repository) *CheckStatusUseCase {
	return &CheckStatusUseCase{
		userRepo: userRepo,
	}
}

func (uc *CheckStatusUseCase) Execute(ctx context.Context) (*CheckStatusOutput, error) {
	count, err := uc.userRepo.CountSuperUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check setup status: %w", err)
	}

	initialized := count > 0

	return &CheckStatusOutput{
		Initialized: initialized,
		SetupNeeded: !initialized,
	}, nil
}

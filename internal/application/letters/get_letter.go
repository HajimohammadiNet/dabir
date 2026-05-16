package letters

import (
	"context"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
)

type GetLetterUseCase struct {
	letterRepo     letter.Repository
	configProvider *LetterConfigProvider
}

func NewGetLetterUseCase(
	letterRepo letter.Repository,
	configProvider *LetterConfigProvider,
) *GetLetterUseCase {
	return &GetLetterUseCase{
		letterRepo:     letterRepo,
		configProvider: configProvider,
	}
}

func (uc *GetLetterUseCase) Execute(ctx context.Context, id string) (*LetterDTO, error) {
	l, err := uc.letterRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find letter: %w", err)
	}

	if l == nil {
		return nil, ErrLetterNotFound
	}

	cfg := uc.configProvider.Get(ctx)
	dto := ToLetterDTO(*l, cfg)

	return &dto, nil
}

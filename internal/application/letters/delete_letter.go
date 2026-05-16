package letters

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
)

type DeleteLetterUseCase struct {
	letterRepo letter.Repository
}

type DeleteLetterInput struct {
	ID          string
	ActorUserID string
}

func NewDeleteLetterUseCase(letterRepo letter.Repository) *DeleteLetterUseCase {
	return &DeleteLetterUseCase{
		letterRepo: letterRepo,
	}
}

func (uc *DeleteLetterUseCase) Execute(ctx context.Context, input DeleteLetterInput) error {
	input.ID = strings.TrimSpace(input.ID)
	input.ActorUserID = strings.TrimSpace(input.ActorUserID)

	if input.ID == "" {
		return errors.New("letter id is required")
	}

	if input.ActorUserID == "" {
		return errors.New("actor user id is required")
	}

	l, err := uc.letterRepo.FindByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("failed to find letter: %w", err)
	}

	if l == nil || l.IsDeleted {
		return ErrLetterNotFound
	}

	if err := uc.letterRepo.SoftDelete(ctx, input.ID, input.ActorUserID); err != nil {
		return fmt.Errorf("failed to delete letter: %w", err)
	}

	return nil
}

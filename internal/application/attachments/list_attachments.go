package attachments

import (
	"context"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/attachment"
	"github.com/hajimohammadinet/dabir/internal/domain/letter"
)

type ListAttachmentsUseCase struct {
	letterRepo     letter.Repository
	attachmentRepo attachment.Repository
}

type ListAttachmentsInput struct {
	LetterID string
}

func NewListAttachmentsUseCase(
	letterRepo letter.Repository,
	attachmentRepo attachment.Repository,
) *ListAttachmentsUseCase {
	return &ListAttachmentsUseCase{
		letterRepo:     letterRepo,
		attachmentRepo: attachmentRepo,
	}
}

func (uc *ListAttachmentsUseCase) Execute(ctx context.Context, input ListAttachmentsInput) ([]AttachmentDTO, error) {
	if input.LetterID == "" {
		return nil, fmt.Errorf("letter id is required")
	}

	letterItem, err := uc.letterRepo.FindByID(ctx, input.LetterID)
	if err != nil {
		return nil, fmt.Errorf("failed to find letter: %w", err)
	}

	if letterItem == nil || letterItem.IsDeleted {
		return nil, fmt.Errorf("letter not found")
	}

	items, err := uc.attachmentRepo.ListByLetterID(ctx, input.LetterID)
	if err != nil {
		return nil, err
	}

	return ToAttachmentDTOs(items), nil
}

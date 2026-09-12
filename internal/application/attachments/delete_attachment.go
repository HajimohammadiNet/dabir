package attachments

import (
	"context"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/attachment"
	"github.com/hajimohammadinet/dabir/internal/domain/storage"
)

type DeleteAttachmentUseCase struct {
	attachmentRepo attachment.Repository
	storageClient  storage.Client
}

type DeleteAttachmentInput struct {
	LetterID     string
	AttachmentID string
	ActorUserID  string
}

func NewDeleteAttachmentUseCase(
	attachmentRepo attachment.Repository,
	storageClient storage.Client,
) *DeleteAttachmentUseCase {
	return &DeleteAttachmentUseCase{
		attachmentRepo: attachmentRepo,
		storageClient:  storageClient,
	}
}

func (uc *DeleteAttachmentUseCase) Execute(ctx context.Context, input DeleteAttachmentInput) (*AttachmentDTO, error) {
	if input.LetterID == "" {
		return nil, fmt.Errorf("letter id is required")
	}

	if input.AttachmentID == "" {
		return nil, fmt.Errorf("attachment id is required")
	}

	if input.ActorUserID == "" {
		return nil, fmt.Errorf("actor user id is required")
	}

	item, err := uc.attachmentRepo.FindByID(ctx, input.AttachmentID)
	if err != nil {
		return nil, err
	}

	if item == nil || item.IsDeleted || item.LetterID != input.LetterID {
		return nil, fmt.Errorf("attachment not found")
	}

	if err := uc.attachmentRepo.SoftDelete(ctx, input.AttachmentID, input.ActorUserID); err != nil {
		return nil, err
	}

	_ = uc.storageClient.DeleteObject(ctx, item.ObjectKey)

	dto := ToAttachmentDTO(*item)
	return &dto, nil
}

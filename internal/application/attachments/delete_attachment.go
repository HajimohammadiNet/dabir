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

func (uc *DeleteAttachmentUseCase) Execute(ctx context.Context, input DeleteAttachmentInput) error {
	if input.LetterID == "" {
		return fmt.Errorf("letter id is required")
	}

	if input.AttachmentID == "" {
		return fmt.Errorf("attachment id is required")
	}

	if input.ActorUserID == "" {
		return fmt.Errorf("actor user id is required")
	}

	item, err := uc.attachmentRepo.FindByID(ctx, input.AttachmentID)
	if err != nil {
		return err
	}

	if item == nil || item.IsDeleted || item.LetterID != input.LetterID {
		return fmt.Errorf("attachment not found")
	}

	if err := uc.attachmentRepo.SoftDelete(ctx, input.AttachmentID, input.ActorUserID); err != nil {
		return err
	}

	_ = uc.storageClient.DeleteObject(ctx, item.ObjectKey)

	return nil
}

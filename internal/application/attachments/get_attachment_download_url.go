package attachments

import (
	"context"
	"fmt"
	"time"

	"github.com/hajimohammadinet/dabir/internal/config"
	"github.com/hajimohammadinet/dabir/internal/domain/attachment"
	"github.com/hajimohammadinet/dabir/internal/domain/storage"
)

type GetAttachmentDownloadURLUseCase struct {
	attachmentRepo attachment.Repository
	storageClient  storage.Client
	storageConfig  config.StorageConfig
}

type GetAttachmentDownloadURLInput struct {
	LetterID     string
	AttachmentID string
	ActorUserID  string
}

type GetAttachmentDownloadURLOutput struct {
	URL string `json:"url"`
}

func NewGetAttachmentDownloadURLUseCase(
	attachmentRepo attachment.Repository,
	storageClient storage.Client,
	storageConfig config.StorageConfig,
) *GetAttachmentDownloadURLUseCase {
	return &GetAttachmentDownloadURLUseCase{
		attachmentRepo: attachmentRepo,
		storageClient:  storageClient,
		storageConfig:  storageConfig,
	}
}

func (uc *GetAttachmentDownloadURLUseCase) Execute(ctx context.Context, input GetAttachmentDownloadURLInput) (*GetAttachmentDownloadURLOutput, error) {
	if input.LetterID == "" {
		return nil, fmt.Errorf("letter id is required")
	}

	if input.AttachmentID == "" {
		return nil, fmt.Errorf("attachment id is required")
	}

	item, err := uc.attachmentRepo.FindByID(ctx, input.AttachmentID)
	if err != nil {
		return nil, err
	}

	if item == nil || item.IsDeleted || item.LetterID != input.LetterID {
		return nil, fmt.Errorf("attachment not found")
	}

	ttl := time.Duration(uc.storageConfig.PresignedURLTTLMinutes) * time.Minute
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	url, err := uc.storageClient.PresignedGetObjectURL(ctx, storage.PresignedGetObjectInput{
		ObjectKey: item.ObjectKey,
		TTL:       ttl,
	})
	if err != nil {
		return nil, err
	}

	return &GetAttachmentDownloadURLOutput{
		URL: url,
	}, nil
}

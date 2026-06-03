package attachments

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hajimohammadinet/dabir/internal/config"
	"github.com/hajimohammadinet/dabir/internal/domain/attachment"
	"github.com/hajimohammadinet/dabir/internal/domain/letter"
	"github.com/hajimohammadinet/dabir/internal/domain/storage"
)

type UploadAttachmentUseCase struct {
	letterRepo     letter.Repository
	attachmentRepo attachment.Repository
	storageClient  storage.Client
	storageConfig  config.StorageConfig
}

type UploadAttachmentInput struct {
	LetterID string

	FileName    string
	ContentType string
	SizeBytes   int64
	Reader      io.Reader

	ActorUserID string
}

func NewUploadAttachmentUseCase(
	letterRepo letter.Repository,
	attachmentRepo attachment.Repository,
	storageClient storage.Client,
	storageConfig config.StorageConfig,
) *UploadAttachmentUseCase {
	return &UploadAttachmentUseCase{
		letterRepo:     letterRepo,
		attachmentRepo: attachmentRepo,
		storageClient:  storageClient,
		storageConfig:  storageConfig,
	}
}

func (uc *UploadAttachmentUseCase) Execute(ctx context.Context, input UploadAttachmentInput) (*AttachmentDTO, error) {
	input.FileName = strings.TrimSpace(input.FileName)
	input.ContentType = strings.TrimSpace(input.ContentType)

	if input.LetterID == "" {
		return nil, fmt.Errorf("letter id is required")
	}

	if input.ActorUserID == "" {
		return nil, fmt.Errorf("actor user id is required")
	}

	if input.FileName == "" {
		return nil, fmt.Errorf("file name is required")
	}

	if input.Reader == nil {
		return nil, fmt.Errorf("file reader is required")
	}

	letterItem, err := uc.letterRepo.FindByID(ctx, input.LetterID)
	if err != nil {
		return nil, fmt.Errorf("failed to find letter: %w", err)
	}

	if letterItem == nil || letterItem.IsDeleted {
		return nil, fmt.Errorf("letter not found")
	}

	maxSizeBytes := uc.storageConfig.MaxUploadSizeMB * 1024 * 1024
	if err := ValidateAttachment(input.ContentType, input.SizeBytes, maxSizeBytes); err != nil {
		return nil, err
	}

	objectKey := BuildAttachmentObjectKey(input.LetterID, input.FileName)

	if err := uc.storageClient.PutObject(ctx, storage.PutObjectInput{
		ObjectKey:   objectKey,
		Reader:      input.Reader,
		SizeBytes:   input.SizeBytes,
		ContentType: input.ContentType,
	}); err != nil {
		return nil, err
	}

	item := &attachment.Attachment{
		LetterID:    input.LetterID,
		FileName:    input.FileName,
		ObjectKey:   objectKey,
		ContentType: input.ContentType,
		SizeBytes:   input.SizeBytes,
		UploadedBy:  input.ActorUserID,
		IsDeleted:   false,
	}

	if err := uc.attachmentRepo.Create(ctx, item); err != nil {
		_ = uc.storageClient.DeleteObject(ctx, objectKey)
		return nil, err
	}

	dto := ToAttachmentDTO(*item)

	return &dto, nil
}

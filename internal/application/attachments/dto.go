package attachments

import (
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/attachment"
)

type AttachmentDTO struct {
	ID string `json:"id"`

	LetterID string `json:"letter_id"`

	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`

	IsDeleted bool `json:"is_deleted"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToAttachmentDTO(item attachment.Attachment) AttachmentDTO {
	return AttachmentDTO{
		ID:          item.ID,
		LetterID:    item.LetterID,
		FileName:    item.FileName,
		ContentType: item.ContentType,
		SizeBytes:   item.SizeBytes,
		IsDeleted:   item.IsDeleted,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func ToAttachmentDTOs(items []attachment.Attachment) []AttachmentDTO {
	result := make([]AttachmentDTO, 0, len(items))

	for _, item := range items {
		result = append(result, ToAttachmentDTO(item))
	}

	return result
}

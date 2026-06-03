package attachments

import (
	"errors"
	"strings"
)

var (
	ErrUnsupportedAttachmentType = errors.New("unsupported attachment file type")
	ErrAttachmentTooLarge        = errors.New("attachment file is too large")
)

var allowedContentTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func ValidateAttachment(contentType string, sizeBytes int64, maxSizeBytes int64) error {
	contentType = strings.ToLower(strings.TrimSpace(contentType))

	if !allowedContentTypes[contentType] {
		return ErrUnsupportedAttachmentType
	}

	if maxSizeBytes > 0 && sizeBytes > maxSizeBytes {
		return ErrAttachmentTooLarge
	}

	return nil
}

package attachments

import (
	"errors"
	"testing"
)

func TestValidateAttachmentAcceptsSupportedTypes(t *testing.T) {
	t.Parallel()

	for _, contentType := range []string{"application/pdf", "image/jpeg", "image/png"} {
		if err := ValidateAttachment(contentType, 1024, 20*1024*1024); err != nil {
			t.Errorf("ValidateAttachment(%q) = %v", contentType, err)
		}
	}
}

func TestValidateAttachmentRejectsUnsupportedType(t *testing.T) {
	t.Parallel()

	err := ValidateAttachment("text/plain", 1024, 20*1024*1024)
	if !errors.Is(err, ErrUnsupportedAttachmentType) {
		t.Fatalf("ValidateAttachment() error = %v, want %v", err, ErrUnsupportedAttachmentType)
	}
}

func TestValidateAttachmentRejectsOversizedFile(t *testing.T) {
	t.Parallel()

	err := ValidateAttachment("image/png", 20*1024*1024+1, 20*1024*1024)
	if !errors.Is(err, ErrAttachmentTooLarge) {
		t.Fatalf("ValidateAttachment() error = %v, want %v", err, ErrAttachmentTooLarge)
	}
}

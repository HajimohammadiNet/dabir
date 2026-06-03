package attachments

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func BuildAttachmentObjectKey(letterID string, fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		ext = ".bin"
	}

	return "letters/" + letterID + "/" + uuid.NewString() + ext
}

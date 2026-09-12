package attachments

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/hajimohammadinet/dabir/internal/config"
	domainattachment "github.com/hajimohammadinet/dabir/internal/domain/attachment"
	"github.com/hajimohammadinet/dabir/internal/domain/letter"
	"github.com/hajimohammadinet/dabir/internal/domain/storage"
)

func TestOutgoingAttachmentUsesExistingLetterAssociationAndStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	letterRepo := &attachmentTestLetterRepository{item: &letter.Letter{
		ID: "outgoing-letter", Direction: letter.DirectionOutgoing,
	}}
	attachmentRepo := &attachmentTestRepository{}
	storageClient := &attachmentTestStorage{}
	useCase := NewUploadAttachmentUseCase(
		letterRepo,
		attachmentRepo,
		storageClient,
		config.StorageConfig{MaxUploadSizeMB: 20},
	)

	dto, err := useCase.Execute(ctx, UploadAttachmentInput{
		LetterID: "outgoing-letter", FileName: "scan.pdf", ContentType: "application/pdf",
		SizeBytes: 4, Reader: strings.NewReader("scan"), ActorUserID: "editor",
	})
	if err != nil {
		t.Fatalf("upload outgoing attachment: %v", err)
	}
	if dto.LetterID != "outgoing-letter" || attachmentRepo.item == nil {
		t.Fatalf("attachment association = %+v", dto)
	}
	if !strings.HasPrefix(attachmentRepo.item.ObjectKey, "letters/outgoing-letter/") {
		t.Fatalf("object key = %q", attachmentRepo.item.ObjectKey)
	}
	if storageClient.put.ObjectKey != attachmentRepo.item.ObjectKey || storageClient.body != "scan" {
		t.Fatalf("stored object = %+v, body %q", storageClient.put, storageClient.body)
	}

	items, err := NewListAttachmentsUseCase(letterRepo, attachmentRepo).Execute(ctx, ListAttachmentsInput{
		LetterID: "outgoing-letter",
	})
	if err != nil || len(items) != 1 || items[0].LetterID != "outgoing-letter" {
		t.Fatalf("list outgoing attachments = %+v, err %v", items, err)
	}
}

type attachmentTestLetterRepository struct {
	item *letter.Letter
}

func (r *attachmentTestLetterRepository) NextNumber(context.Context, letter.Direction) (int64, error) {
	return 1, nil
}
func (r *attachmentTestLetterRepository) NextNumberForYear(context.Context, letter.Direction, int) (int64, error) {
	return 1, nil
}
func (r *attachmentTestLetterRepository) ExistsByDisplayLetterNumber(context.Context, letter.Direction, string) (bool, error) {
	return false, nil
}
func (r *attachmentTestLetterRepository) FindLatestDisplayLetterNumberByPrefix(context.Context, letter.Direction, string) (*string, error) {
	return nil, nil
}
func (r *attachmentTestLetterRepository) Create(context.Context, *letter.Letter) error { return nil }
func (r *attachmentTestLetterRepository) FindByID(_ context.Context, id string) (*letter.Letter, error) {
	if r.item == nil || r.item.ID != id {
		return nil, nil
	}
	return r.item, nil
}
func (r *attachmentTestLetterRepository) List(context.Context, letter.ListFilter) ([]letter.Letter, int, error) {
	return nil, 0, nil
}
func (r *attachmentTestLetterRepository) Update(context.Context, *letter.Letter) error { return nil }
func (r *attachmentTestLetterRepository) SoftDelete(context.Context, string, string) error {
	return nil
}
func (r *attachmentTestLetterRepository) BulkCreate(context.Context, []letter.Letter) error {
	return nil
}
func (r *attachmentTestLetterRepository) SetSequenceValue(context.Context, int64) error { return nil }
func (r *attachmentTestLetterRepository) FindExistingNumbers(context.Context, []int64) (map[int64]bool, error) {
	return map[int64]bool{}, nil
}

type attachmentTestRepository struct {
	item *domainattachment.Attachment
}

func (r *attachmentTestRepository) Create(_ context.Context, item *domainattachment.Attachment) error {
	item.ID = "attachment-id"
	item.CreatedAt = time.Now().UTC()
	item.UpdatedAt = item.CreatedAt
	copy := *item
	r.item = &copy
	return nil
}
func (r *attachmentTestRepository) ListByLetterID(_ context.Context, letterID string) ([]domainattachment.Attachment, error) {
	if r.item == nil || r.item.LetterID != letterID || r.item.IsDeleted {
		return []domainattachment.Attachment{}, nil
	}
	return []domainattachment.Attachment{*r.item}, nil
}
func (r *attachmentTestRepository) FindByID(_ context.Context, id string) (*domainattachment.Attachment, error) {
	if r.item == nil || r.item.ID != id {
		return nil, nil
	}
	return r.item, nil
}
func (r *attachmentTestRepository) SoftDelete(_ context.Context, id string, deletedBy string) error {
	if r.item != nil && r.item.ID == id {
		r.item.IsDeleted = true
		r.item.DeletedBy = &deletedBy
	}
	return nil
}

type attachmentTestStorage struct {
	put  storage.PutObjectInput
	body string
}

func (s *attachmentTestStorage) PutObject(_ context.Context, input storage.PutObjectInput) error {
	s.put = input
	body, err := io.ReadAll(input.Reader)
	if err != nil {
		return err
	}
	s.body = string(body)
	return nil
}
func (s *attachmentTestStorage) PresignedGetObjectURL(context.Context, storage.PresignedGetObjectInput) (string, error) {
	return "https://example.test/object", nil
}
func (s *attachmentTestStorage) DeleteObject(context.Context, string) error { return nil }

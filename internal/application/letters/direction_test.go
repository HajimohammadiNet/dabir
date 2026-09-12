package letters

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
	"github.com/hajimohammadinet/dabir/internal/domain/settings"
)

func TestCreateAndListLettersSeparateDirections(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newMemoryLetterRepository()
	provider := newTestConfigProvider(t, NumberingModeFixedPrefix)
	create := NewCreateLetterUseCase(repo, provider)
	list := NewListLettersUseCase(repo, provider)

	incoming, err := create.Execute(ctx, validCreateInput("Incoming", ""))
	if err != nil {
		t.Fatalf("create incoming: %v", err)
	}
	outgoingInput := validCreateInput("Outgoing", letter.DirectionOutgoing)
	outgoing, err := create.Execute(ctx, outgoingInput)
	if err != nil {
		t.Fatalf("create outgoing: %v", err)
	}

	if incoming.Direction != letter.DirectionIncoming {
		t.Fatalf("default direction = %q, want incoming", incoming.Direction)
	}
	if outgoing.Direction != letter.DirectionOutgoing {
		t.Fatalf("outgoing direction = %q", outgoing.Direction)
	}
	if incoming.LetterNumber != 1 || outgoing.LetterNumber != 1 {
		t.Fatalf("independent numbers = incoming %d, outgoing %d; want 1 and 1", incoming.LetterNumber, outgoing.LetterNumber)
	}

	defaultList, err := list.Execute(ctx, ListLettersInput{})
	if err != nil {
		t.Fatalf("list default incoming: %v", err)
	}
	if defaultList.Total != 1 || defaultList.Items[0].Direction != letter.DirectionIncoming {
		t.Fatalf("default list did not preserve incoming-only behavior: %+v", defaultList.Items)
	}

	outgoingList, err := list.Execute(ctx, ListLettersInput{Direction: letter.DirectionOutgoing})
	if err != nil {
		t.Fatalf("list outgoing: %v", err)
	}
	if outgoingList.Total != 1 || outgoingList.Items[0].ID != outgoing.ID {
		t.Fatalf("outgoing list = %+v", outgoingList.Items)
	}
}

func TestManualDuplicateProtectionIsDirectionScoped(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newMemoryLetterRepository()
	create := NewCreateLetterUseCase(repo, newTestConfigProvider(t, NumberingModeManual))
	number := "405-ق-001"

	incoming := validCreateInput("Incoming", letter.DirectionIncoming)
	incoming.DisplayLetterNumber = &number
	if _, err := create.Execute(ctx, incoming); err != nil {
		t.Fatalf("create incoming: %v", err)
	}

	outgoing := validCreateInput("Outgoing", letter.DirectionOutgoing)
	outgoing.DisplayLetterNumber = &number
	if _, err := create.Execute(ctx, outgoing); err != nil {
		t.Fatalf("same number in outgoing direction should be allowed: %v", err)
	}

	persianNumber := "۴۰۵-ق-۰۰۱"
	outgoing.DisplayLetterNumber = &persianNumber
	if _, err := create.Execute(ctx, outgoing); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("duplicate outgoing number error = %v", err)
	}
}

func TestJalaliYearlyNumberingIsDirectionScoped(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newMemoryLetterRepository()
	create := NewCreateLetterUseCase(repo, newTestConfigProvider(t, NumberingModeJalaliYearly))

	incoming, err := create.Execute(ctx, validCreateInput("Incoming yearly", letter.DirectionIncoming))
	if err != nil {
		t.Fatalf("create incoming yearly: %v", err)
	}

	outgoingInput := validCreateInput("Outgoing yearly", letter.DirectionOutgoing)
	outgoing, err := create.Execute(ctx, outgoingInput)
	if err != nil {
		t.Fatalf("create outgoing yearly: %v", err)
	}

	if incoming.LetterSerial == nil || *incoming.LetterSerial != 1 {
		t.Fatalf("incoming yearly serial = %v, want 1", incoming.LetterSerial)
	}
	if outgoing.LetterSerial == nil || *outgoing.LetterSerial != 1 {
		t.Fatalf("outgoing yearly serial = %v, want 1", outgoing.LetterSerial)
	}
	if incoming.LetterYear == nil || outgoing.LetterYear == nil || *incoming.LetterYear != *outgoing.LetterYear {
		t.Fatalf("yearly years = incoming %v, outgoing %v", incoming.LetterYear, outgoing.LetterYear)
	}
}

func TestSmartSuggestionIsDirectionAwareAndNormalizesPersianDigits(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newMemoryLetterRepository()
	incomingNumber := "405-ق-099"
	outgoingNumber := "405-ق-004"
	repo.items["incoming"] = &letter.Letter{
		ID: "incoming", Direction: letter.DirectionIncoming, LetterNumber: 99,
		DisplayLetterNumber: &incomingNumber,
	}
	repo.items["outgoing"] = &letter.Letter{
		ID: "outgoing", Direction: letter.DirectionOutgoing, LetterNumber: 4,
		DisplayLetterNumber: &outgoingNumber,
	}

	output, err := NewSuggestLetterNumberUseCase(repo).Execute(ctx, SuggestLetterNumberInput{
		Direction: letter.DirectionOutgoing,
		Prefix:    "۴۰۵-ق-",
	})
	if err != nil {
		t.Fatalf("suggest outgoing number: %v", err)
	}
	if output.LastNumber == nil || *output.LastNumber != outgoingNumber {
		t.Fatalf("last number = %v, want %q", output.LastNumber, outgoingNumber)
	}
	if output.SuggestedNumber == nil || *output.SuggestedNumber != "405-ق-005" {
		t.Fatalf("suggested number = %v", output.SuggestedNumber)
	}
}

func TestUpdateAndSoftDeleteOutgoingLetter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newMemoryLetterRepository()
	provider := newTestConfigProvider(t, NumberingModeManual)
	create := NewCreateLetterUseCase(repo, provider)
	update := NewUpdateLetterUseCase(repo, provider)
	remove := NewDeleteLetterUseCase(repo)
	number := "405-ق-010"
	input := validCreateInput("Before", letter.DirectionOutgoing)
	input.DisplayLetterNumber = &number
	created, err := create.Execute(ctx, input)
	if err != nil {
		t.Fatalf("create outgoing: %v", err)
	}

	newNumber := "405-ق-011"
	updated, err := update.Execute(ctx, UpdateLetterInput{
		ID: created.ID, ActorUserID: "editor", DisplayLetterNumber: &newNumber,
		Title: "After", LetterDate: "1405/06/22", Sender: "Dabir", Receiver: "Client",
	})
	if err != nil {
		t.Fatalf("update outgoing: %v", err)
	}
	if updated.Title != "After" || updated.Direction != letter.DirectionOutgoing {
		t.Fatalf("updated letter = %+v", updated)
	}

	if err := remove.Execute(ctx, DeleteLetterInput{ID: created.ID, ActorUserID: "editor"}); err != nil {
		t.Fatalf("soft delete outgoing: %v", err)
	}
	stored, _ := repo.FindByID(ctx, created.ID)
	if stored == nil || !stored.IsDeleted || stored.DeletedBy == nil || *stored.DeletedBy != "editor" {
		t.Fatalf("stored deleted letter = %+v", stored)
	}
}

func validCreateInput(title string, direction letter.Direction) CreateLetterInput {
	return CreateLetterInput{
		Direction: direction, Title: title, LetterDate: "1405/06/21",
		Sender: "Sender", Receiver: "Receiver", ActorUserID: "editor", RegistrarName: "editor",
	}
}

func newTestConfigProvider(t *testing.T, mode NumberingMode) *LetterConfigProvider {
	t.Helper()
	config := settings.DefaultLetterConfig()
	config.NumberingMode = settings.NumberingMode(mode)
	raw, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	repo := &memorySettingsRepository{values: map[string][]byte{
		settings.KeyLetterConfig:         raw,
		settings.KeyOutgoingLetterConfig: raw,
	}}
	return NewLetterConfigProvider(repo)
}

type memorySettingsRepository struct {
	values map[string][]byte
}

func (r *memorySettingsRepository) Set(_ context.Context, key string, value []byte) error {
	r.values[key] = value
	return nil
}

func (r *memorySettingsRepository) Get(_ context.Context, key string) (*settings.Setting, error) {
	value, ok := r.values[key]
	if !ok {
		return nil, nil
	}
	return &settings.Setting{Key: key, Value: value}, nil
}

type memoryLetterRepository struct {
	items      map[string]*letter.Letter
	last       map[letter.Direction]int64
	yearlyLast map[letter.Direction]map[int]int64
	nextID     int
}

func newMemoryLetterRepository() *memoryLetterRepository {
	return &memoryLetterRepository{
		items: make(map[string]*letter.Letter), last: make(map[letter.Direction]int64),
		yearlyLast: make(map[letter.Direction]map[int]int64),
	}
}

func (r *memoryLetterRepository) NextNumber(_ context.Context, direction letter.Direction) (int64, error) {
	r.last[direction]++
	return r.last[direction], nil
}

func (r *memoryLetterRepository) NextNumberForYear(_ context.Context, direction letter.Direction, year int) (int64, error) {
	if r.yearlyLast[direction] == nil {
		r.yearlyLast[direction] = make(map[int]int64)
	}
	r.yearlyLast[direction][year]++
	return r.yearlyLast[direction][year], nil
}

func (r *memoryLetterRepository) ExistsByDisplayLetterNumber(_ context.Context, direction letter.Direction, number string) (bool, error) {
	for _, item := range r.items {
		if !item.IsDeleted && item.Direction == direction && item.DisplayLetterNumber != nil && normalizePersianArabicDigits(*item.DisplayLetterNumber) == normalizePersianArabicDigits(number) {
			return true, nil
		}
	}
	return false, nil
}

func (r *memoryLetterRepository) FindLatestDisplayLetterNumberByPrefix(_ context.Context, direction letter.Direction, prefix string) (*string, error) {
	var latest *letter.Letter
	for _, item := range r.items {
		if item.IsDeleted || item.Direction != direction || item.DisplayLetterNumber == nil {
			continue
		}
		if prefix != "" && !strings.HasPrefix(normalizePersianArabicDigits(*item.DisplayLetterNumber), prefix) {
			continue
		}
		if latest == nil || item.LetterNumber > latest.LetterNumber {
			latest = item
		}
	}
	if latest == nil {
		return nil, nil
	}
	value := *latest.DisplayLetterNumber
	return &value, nil
}

func (r *memoryLetterRepository) Create(_ context.Context, item *letter.Letter) error {
	r.nextID++
	item.ID = "letter-" + time.Unix(int64(r.nextID), 0).UTC().Format("150405")
	item.CreatedAt = time.Now().UTC()
	item.UpdatedAt = item.CreatedAt
	copy := *item
	r.items[item.ID] = &copy
	return nil
}

func (r *memoryLetterRepository) FindByID(_ context.Context, id string) (*letter.Letter, error) {
	item := r.items[id]
	if item == nil {
		return nil, nil
	}
	copy := *item
	return &copy, nil
}

func (r *memoryLetterRepository) List(_ context.Context, filter letter.ListFilter) ([]letter.Letter, int, error) {
	items := make([]letter.Letter, 0)
	for _, item := range r.items {
		if !filter.IncludeDeleted && item.IsDeleted {
			continue
		}
		if filter.Direction != "" && item.Direction != filter.Direction {
			continue
		}
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items, len(items), nil
}

func (r *memoryLetterRepository) Update(_ context.Context, item *letter.Letter) error {
	if r.items[item.ID] == nil {
		return errors.New("not found")
	}
	item.UpdatedAt = time.Now().UTC()
	copy := *item
	r.items[item.ID] = &copy
	return nil
}

func (r *memoryLetterRepository) SoftDelete(_ context.Context, id string, deletedBy string) error {
	item := r.items[id]
	if item == nil || item.IsDeleted {
		return errors.New("not found")
	}
	now := time.Now().UTC()
	item.IsDeleted = true
	item.DeletedBy = &deletedBy
	item.DeletedAt = &now
	return nil
}

func (r *memoryLetterRepository) BulkCreate(ctx context.Context, items []letter.Letter) error {
	for i := range items {
		if err := r.Create(ctx, &items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *memoryLetterRepository) SetSequenceValue(_ context.Context, value int64) error {
	r.last[letter.DirectionIncoming] = value
	return nil
}

func (r *memoryLetterRepository) FindExistingNumbers(_ context.Context, numbers []int64) (map[int64]bool, error) {
	wanted := make(map[int64]bool, len(numbers))
	for _, number := range numbers {
		wanted[number] = false
	}
	for _, item := range r.items {
		if item.Direction == letter.DirectionIncoming {
			if _, ok := wanted[item.LetterNumber]; ok {
				wanted[item.LetterNumber] = true
			}
		}
	}
	return wanted, nil
}

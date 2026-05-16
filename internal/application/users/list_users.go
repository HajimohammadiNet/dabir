package users

import (
	"context"
	"strings"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

type ListUsersUseCase struct {
	userRepo user.Repository
}

type ListUsersInput struct {
	Page     int
	PageSize int
	Search   string
	Role     *user.Role
	IsActive *bool
}

type ListUsersOutput struct {
	Items      []UserDTO `json:"items"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}

func NewListUsersUseCase(userRepo user.Repository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}

	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	input.Search = strings.TrimSpace(input.Search)

	usersList, total, err := uc.userRepo.List(ctx, user.ListFilter{
		Page:     input.Page,
		PageSize: input.PageSize,
		Search:   input.Search,
		Role:     input.Role,
		IsActive: input.IsActive,
	})
	if err != nil {
		return nil, err
	}

	items := make([]UserDTO, 0, len(usersList))
	for _, u := range usersList {
		items = append(items, ToUserDTO(u))
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + input.PageSize - 1) / input.PageSize
	}

	return &ListUsersOutput{
		Items:      items,
		Total:      total,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}

package users

import (
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

type UserDTO struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Role      user.Role `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserDTO(u user.User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

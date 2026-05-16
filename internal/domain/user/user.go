package user

import "time"

type Role string

const (
	RoleSuperUser Role = "superuser"
	RoleEditor    Role = "editor"
	RoleReadonly  Role = "readonly"
)

type User struct {
	ID           string
	Username     string
	FullName     string
	PasswordHash string
	Role         Role
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

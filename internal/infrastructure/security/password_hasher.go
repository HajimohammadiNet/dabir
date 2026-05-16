package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		cost: bcrypt.DefaultCost,
	}
}

func (h *PasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

func (h *PasswordHasher) Compare(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

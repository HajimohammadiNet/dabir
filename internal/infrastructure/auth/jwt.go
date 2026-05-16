package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/hajimohammadinet/dabir/internal/config"
	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

type Claims struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Role     user.Role `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTService(cfg config.AuthConfig) *JWTService {
	return &JWTService{
		secret: []byte(cfg.JWTSecret),
		ttl:    time.Duration(cfg.JWTAccessTokenTTLMinutes) * time.Minute,
	}
}

func (s *JWTService) GenerateAccessToken(u *user.User) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(s.ttl)

	claims := Claims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign jwt token: %w", err)
	}

	return signedToken, int64(s.ttl.Seconds()), nil
}

func (s *JWTService) ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
			}

			return s.secret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid jwt token")
	}

	return claims, nil
}

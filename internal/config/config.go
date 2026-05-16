package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App  AppConfig
	DB   DatabaseConfig
	Auth AuthConfig
}

type AuthConfig struct {
	JWTSecret                string `env:"JWT_SECRET,required"`
	JWTAccessTokenTTLMinutes int    `env:"JWT_ACCESS_TOKEN_TTL_MINUTES" envDefault:"60"`
}

type AppConfig struct {
	Name string `env:"APP_NAME" envDefault:"dabir"`
	Env  string `env:"APP_ENV" envDefault:"development"`
	Host string `env:"APP_HOST" envDefault:"0.0.0.0"`
	Port int    `env:"APP_PORT" envDefault:"8080"`
}

type DatabaseConfig struct {
	Host         string `env:"DB_HOST" envDefault:"localhost"`
	Port         int    `env:"DB_PORT" envDefault:"5432"`
	User         string `env:"DB_USER,required"`
	Password     string `env:"DB_PASSWORD,required"`
	Name         string `env:"DB_NAME,required"`
	SSLMode      string `env:"DB_SSLMODE" envDefault:"disable"`
	MaxOpenConns int32  `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns int32  `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

func (c AppConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

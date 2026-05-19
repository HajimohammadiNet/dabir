package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hajimohammadinet/dabir/internal/config"
	deliveryhttp "github.com/hajimohammadinet/dabir/internal/delivery/http"
	"github.com/hajimohammadinet/dabir/internal/infrastructure/postgres"
	"github.com/hajimohammadinet/dabir/internal/shared/logger"
)

type App struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *pgxpool.Pool
	Server *http.Server
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := logger.New(cfg.App.Env)

	db, err := postgres.NewPool(ctx, cfg.DB)
	if err != nil {
		return nil, err
	}

	router := deliveryhttp.NewRouter(db, cfg, log)

	server := &http.Server{
		Addr:              cfg.App.Address(),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &App{
		Config: cfg,
		Logger: log,
		DB:     db,
		Server: server,
	}, nil
}

func (a *App) Start() error {
	a.Logger.Info(
		"starting server",
		"app", a.Config.App.Name,
		"env", a.Config.App.Env,
		"address", a.Config.App.Address(),
	)

	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Logger.Info("shutting down server")

	if a.DB != nil {
		a.DB.Close()
	}

	if a.Server != nil {
		return a.Server.Shutdown(ctx)
	}

	return nil
}

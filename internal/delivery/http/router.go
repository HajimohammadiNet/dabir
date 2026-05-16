package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	authapp "github.com/hajimohammadinet/dabir/internal/application/auth"
	lettersapp "github.com/hajimohammadinet/dabir/internal/application/letters"
	settingsapp "github.com/hajimohammadinet/dabir/internal/application/settings"
	setupapp "github.com/hajimohammadinet/dabir/internal/application/setup"
	usersapp "github.com/hajimohammadinet/dabir/internal/application/users"
	"github.com/hajimohammadinet/dabir/internal/config"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/handlers"
	httpmiddleware "github.com/hajimohammadinet/dabir/internal/delivery/http/middleware"
	"github.com/hajimohammadinet/dabir/internal/domain/user"
	infraauth "github.com/hajimohammadinet/dabir/internal/infrastructure/auth"
	"github.com/hajimohammadinet/dabir/internal/infrastructure/postgres"
	"github.com/hajimohammadinet/dabir/internal/infrastructure/security"
)

func NewRouter(db *pgxpool.Pool, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	healthHandler := handlers.NewHealthHandler(db)

	userRepo := postgres.NewUserRepository(db)
	settingsRepo := postgres.NewSettingsRepository(db)
	passwordHasher := security.NewPasswordHasher()
	jwtService := infraauth.NewJWTService(cfg.Auth)

	checkStatusUseCase := setupapp.NewCheckStatusUseCase(userRepo)
	initializeUseCase := setupapp.NewInitializeUseCase(
		userRepo,
		settingsRepo,
		passwordHasher,
	)

	setupHandler := handlers.NewSetupHandler(
		checkStatusUseCase,
		initializeUseCase,
	)

	getPublicSettingsUseCase := settingsapp.NewGetPublicSettingsUseCase(settingsRepo)
	settingsHandler := handlers.NewSettingsHandler(getPublicSettingsUseCase)

	loginUseCase := authapp.NewLoginUseCase(
		userRepo,
		passwordHasher,
		jwtService,
	)
	meUseCase := authapp.NewMeUseCase(userRepo)
	authHandler := handlers.NewAuthHandler(loginUseCase, meUseCase)

	createUserUseCase := usersapp.NewCreateUserUseCase(userRepo, passwordHasher)
	listUsersUseCase := usersapp.NewListUsersUseCase(userRepo)
	getUserUseCase := usersapp.NewGetUserUseCase(userRepo)
	updateUserUseCase := usersapp.NewUpdateUserUseCase(userRepo)
	setUserActiveUseCase := usersapp.NewSetUserActiveUseCase(userRepo)

	userHandler := handlers.NewUserHandler(
		createUserUseCase,
		listUsersUseCase,
		getUserUseCase,
		updateUserUseCase,
		setUserActiveUseCase,
	)

	letterRepo := postgres.NewLetterRepository(db)
	letterConfigProvider := lettersapp.NewLetterConfigProvider(settingsRepo)

	createLetterUseCase := lettersapp.NewCreateLetterUseCase(letterRepo, letterConfigProvider)
	listLettersUseCase := lettersapp.NewListLettersUseCase(letterRepo, letterConfigProvider)
	getLetterUseCase := lettersapp.NewGetLetterUseCase(letterRepo, letterConfigProvider)
	updateLetterUseCase := lettersapp.NewUpdateLetterUseCase(letterRepo, letterConfigProvider)
	deleteLetterUseCase := lettersapp.NewDeleteLetterUseCase(letterRepo)

	letterHandler := handlers.NewLetterHandler(
		createLetterUseCase,
		listLettersUseCase,
		getLetterUseCase,
		updateLetterUseCase,
		deleteLetterUseCase,
	)

	r.Get("/healthz", healthHandler.Healthz)
	r.Get("/readyz", healthHandler.Readyz)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", healthHandler.Healthz)
		r.Get("/readyz", healthHandler.Readyz)

		r.Route("/setup", func(r chi.Router) {
			r.Get("/status", setupHandler.Status)
			r.Post("/initialize", setupHandler.Initialize)
		})

		r.Route("/settings", func(r chi.Router) {
			r.Get("/public", settingsHandler.Public)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)

			r.Group(func(r chi.Router) {
				r.Use(httpmiddleware.AuthMiddleware(jwtService))

				r.Get("/me", authHandler.Me)

				r.With(httpmiddleware.RequireRoles(user.RoleSuperUser)).
					Get("/superuser-check", func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusNoContent)
					})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(httpmiddleware.AuthMiddleware(jwtService))
			r.Use(httpmiddleware.RequireRoles(user.RoleSuperUser))

			r.Post("/", userHandler.Create)
			r.Get("/", userHandler.List)
			r.Get("/{id}", userHandler.GetByID)
			r.Patch("/{id}", userHandler.Update)
			r.Patch("/{id}/deactivate", userHandler.Deactivate)
			r.Patch("/{id}/activate", userHandler.Activate)
		})

		r.Route("/letters", func(r chi.Router) {
			r.Use(httpmiddleware.AuthMiddleware(jwtService))

			r.With(httpmiddleware.RequireRoles(user.RoleSuperUser, user.RoleEditor, user.RoleReadonly)).
				Get("/", letterHandler.List)

			r.With(httpmiddleware.RequireRoles(user.RoleSuperUser, user.RoleEditor, user.RoleReadonly)).
				Get("/{id}", letterHandler.GetByID)

			r.With(httpmiddleware.RequireRoles(user.RoleSuperUser, user.RoleEditor)).
				Post("/", letterHandler.Create)

			r.With(httpmiddleware.RequireRoles(user.RoleSuperUser, user.RoleEditor)).
				Patch("/{id}", letterHandler.Update)

			r.With(httpmiddleware.RequireRoles(user.RoleSuperUser, user.RoleEditor)).
				Delete("/{id}", letterHandler.Delete)
		})
	})

	return r
}

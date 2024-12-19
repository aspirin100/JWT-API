package main

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/goose/v3"

	"github.com/aspirin100/JWT-API/internal/api"
	"github.com/aspirin100/JWT-API/internal/logger"
	"github.com/aspirin100/JWT-API/internal/middleware"
	"github.com/aspirin100/JWT-API/internal/notifier"
	gen "github.com/aspirin100/JWT-API/internal/oas/generated"
	"github.com/aspirin100/JWT-API/internal/token"
	"github.com/aspirin100/JWT-API/migrations"
)

func main() {
	config := Config{}

	err := envconfig.Process("jwt-server", &config)
	if err != nil {
		logger.Fatal("failed to read configuration")
	}

	logger.Default().Debug("server configuration", "config", config)

	db, err := sql.Open("postgres", config.PostgresDSN)
	if err != nil {
		logger.Fatal(err.Error())
	}

	goose.SetBaseFS(migrations.Migrations)

	err = goose.Up(db, ".")
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Default().Debug("migrations up")

	service := &token.Service{
		SecretKeys:         config.SecretKeys,
		CurrentSecretKeyID: config.SecretKeyID,
		RefreshTokenTTL:    config.RefreshTokenTTL,
		AccessTokenTTL:     config.AccessTokenTTL,
		Repository: &token.PostgresRepository{
			DB: db,
		},
		Notifier: new(notifier.Mock),
	}

	handler := api.Handler{
		TokenService: service,
	}

	authServer, err := gen.NewServer(&handler, gen.WithMiddleware(middleware.Recover, middleware.DetectIP))
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Default().Info("Listening on " + config.Hostname)

	err = http.ListenAndServe(config.Hostname, authServer) //nolint:gosec
	if err != nil {
		logger.Fatal(err.Error())
	}
}

type Config struct {
	PostgresDSN     string            `envconfig:"JWT_SERVER_POSTGRES_DSN" default:"postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"` //nolint:lll
	Hostname        string            `envconfig:"JWT_SERVER_HOSTNAME" default:":8000"`
	SecretKeys      map[string]string `envconfig:"JWT_SERVER_SECRET_KEYS"`
	SecretKeyID     string            `envconfig:"JWT_SERVER_SECRET_KEY_ID"`
	RefreshTokenTTL time.Duration     `envconfig:"JWT_SERVER_REFRESH_TTL_MINUTES" default:"43200m"`
	AccessTokenTTL  time.Duration     `envconfig:"JWT_SERVER_ACCESS_TTL_MINUTES" default:"15m"`
}

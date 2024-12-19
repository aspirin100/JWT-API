package token_test

import (
	"context"
	"database/sql"
	"net/url"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/aspirin100/JWT-API/internal/token"
	"github.com/aspirin100/JWT-API/migrations"
)

func TestPostgresRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &PostgresRepositoryTestSuite{})
}

type PostgresRepositoryTestSuite struct {
	suite.Suite

	ctx    context.Context
	cancel context.CancelFunc
	cont   testcontainers.Container

	db         *sql.DB
	repository *token.PostgresRepository
}

func (s *PostgresRepositoryTestSuite) SetupSuite() {
	postgresListeningPort, err := nat.NewPort("tcp", "5432")
	s.Require().NoError(err)

	s.ctx, s.cancel = context.WithTimeout(context.Background(), time.Minute*5)

	pgContainer, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:16-alpine",
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "auth",
			},
			ExposedPorts: []string{string(postgresListeningPort)},
			WaitingFor: wait.ForAll(
				wait.ForExposedPort(),
				wait.ForLog("database system is ready to accept connections").WithStartupTimeout(time.Minute/2),
			),
			Name: "jwt_postgres_" + uuid.NewString(),
		},
		Started: true,
	})
	s.Require().NoError(err)

	s.cont = pgContainer

	pgPort, err := pgContainer.MappedPort(s.ctx, "5432/tcp")
	s.Require().NoError(err)

	pgHost, err := pgContainer.Host(s.ctx)
	s.Require().NoError(err)

	uri := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     pgHost + ":" + pgPort.Port(),
		Path:     "/auth",
		RawQuery: "sslmode=disable",
	}

	db, err := sql.Open("postgres", uri.String())
	s.Require().NoError(err)
	s.Require().NotNil(db)

	s.db = db
	s.repository = &token.PostgresRepository{
		DB: s.db,
	}

	goose.SetBaseFS(migrations.Migrations)
	s.Require().NoError(goose.Up(s.db, "."))
}

func (s *PostgresRepositoryTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
	s.NoError(s.cont.Terminate(s.ctx))
	s.cancel()
}

func (s *PostgresRepositoryTestSuite) TestBeginTx() {
	ctx, commitOrRollback, err := s.repository.BeginTx(s.ctx)
	s.NoError(err)
	s.NotSame(s.ctx, ctx)

	pairID := uuid.New()
	userID := uuid.MustParse("e05fa11d-eec3-4fba-b223-d6516800a047")

	err = s.repository.InsertRefreshToken(ctx, pairID, userID, uuid.NewString())
	s.NoError(err)

	s.NoError(commitOrRollback(&err))
}

func (s *PostgresRepositoryTestSuite) TestInsertRefreshTokenFailed() {
	ctx, commitOrRollback, err := s.repository.BeginTx(s.ctx)
	s.NoError(err)
	s.NotSame(s.ctx, ctx)

	pairID := uuid.New()
	userID := uuid.New()

	err = s.repository.InsertRefreshToken(ctx, pairID, userID, uuid.NewString())
	s.Error(err)
	s.Error(commitOrRollback(&err))
	s.ErrorIs(err, token.ErrUserNotFound)
}

func (s *PostgresRepositoryTestSuite) TestSetRefreshTokenUsed() {
	pairID := uuid.New()
	userID := uuid.MustParse("e05fa11d-eec3-4fba-b223-d6516800a047")

	err := s.repository.InsertRefreshToken(s.ctx, pairID, userID, uuid.NewString())
	s.NoError(err)

	err = s.repository.SetRefreshTokenUsed(s.ctx, pairID)
	s.NoError(err)

	err = s.repository.SetRefreshTokenUsed(s.ctx, pairID)
	s.Error(err)
}

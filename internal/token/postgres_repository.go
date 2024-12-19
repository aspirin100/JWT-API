package token

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type ctxKey struct{}

var txContextKey = ctxKey{}

type PostgresRepository struct {
	DB *sql.DB
}

func (r *PostgresRepository) BeginTx(ctx context.Context) (context.Context, CommitOrRollback, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return context.WithValue(ctx, txContextKey, tx), func(err *error) error {
		if *err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				return errors.Join(*err, errRollback)
			}

			return *err
		}

		if errCommit := tx.Commit(); errCommit != nil {
			return fmt.Errorf("failed to commit transaction: %w", errCommit)
		}

		return nil
	}, nil
}

func (r *PostgresRepository) InsertRefreshToken(
	ctx context.Context,
	pairID, userID uuid.UUID,
	refreshToken string,
) error {
	type executor interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}

	var ex executor = r.DB

	tx, ok := ctx.Value(txContextKey).(*sql.Tx)
	if ok {
		ex = tx
	}

	// на случай, если рефреш токен слишком велик для bcrypt хэша
	h := sha512.Sum512([]byte(refreshToken))

	refreshTokenHashed, err := bcrypt.GenerateFromPassword(h[:], bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("failed to create bcrypt hash: %w", err)
	}

	_, err = ex.ExecContext(ctx, insertQuery, pairID, userID, refreshTokenHashed)
	if err != nil {
		// проверка, соответствует ли полученный error ошибке, возвращенной после запроса
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrUserNotFound
		}
		// ошибка выполнения запроса
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *PostgresRepository) SetRefreshTokenUsed(ctx context.Context, pairID uuid.UUID) error {
	type executor interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}

	var ex executor = r.DB

	tx, ok := ctx.Value(txContextKey).(*sql.Tx)
	if ok {
		ex = tx
	}

	res, err := ex.ExecContext(ctx, updateQuery, pairID)
	if err != nil {
		return fmt.Errorf("failed to execute query with context: %w", err)
	}

	rowsChanged, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected func failed: %w", err)
	}

	if rowsChanged == 0 {
		return ErrRefreshTokenUsed
	}

	return nil
}

const (
	insertQuery = `insert into refresh_tokens (pair_id, user_id, refresh_token) values ($1, $2, $3)`
	updateQuery = `update refresh_tokens set used = true where pair_id = $1 and used = false`
)

package repository

import (
	"context"
	"database/sql"
	"errors"

	"notes-app/internal/model"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type PostgresRefreshTokenRepository struct {
	db *sql.DB
}

func NewPostgresRefreshTokenRepository(db *sql.DB) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

func (r *PostgresRefreshTokenRepository) Create(
	ctx context.Context,
	token model.RefreshToken,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
	)
	return err
}

func (r *PostgresRefreshTokenRepository) GetByToken(
	ctx context.Context,
	token string,
) (*model.RefreshToken, error) {
	var rt model.RefreshToken

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, token, expires_at, created_at
		 FROM refresh_tokens WHERE token = $1`,
		token,
	).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	return &rt, nil
}

func (r *PostgresRefreshTokenRepository) DeleteByToken(
	ctx context.Context,
	token string,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM refresh_tokens WHERE token = $1`,
		token,
	)
	return err
}

func (r *PostgresRefreshTokenRepository) DeleteAllByUserID(
	ctx context.Context,
	userID string,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM refresh_tokens WHERE user_id = $1`,
		userID,
	)
	return err
}

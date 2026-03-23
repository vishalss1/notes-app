package repository

import (
	"context"
	"database/sql"
	"errors"

	"notes-app/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(
	ctx context.Context,
	user model.User,
) (*model.User, error) {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, email, password)
		VALUES ($1, $2, $3)`,
		user.ID,
		user.Email,
		user.Password,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	var n model.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, password FROM users WHERE email = $1`,
		email,
	).Scan(&n.ID, &n.Email, &n.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *PostgresUserRepository) GetByID(
	ctx context.Context,
	id string,
) (*model.User, error) {
	var n model.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, password FROM users WHERE id = $1`,
		id,
	).Scan(&n.ID, &n.Email, &n.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &n, nil
}

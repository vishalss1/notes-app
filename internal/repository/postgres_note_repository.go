package repository

import (
	"context"
	"database/sql"
	"errors"

	"notes-app/internal/model"
)

var ErrNoteNotFound = errors.New("note not found")

type PostgresNoteRepository struct {
	db *sql.DB
}

func NewPostgresNoteRepository(db *sql.DB) *PostgresNoteRepository {
	return &PostgresNoteRepository{db: db}
}

func (r *PostgresNoteRepository) Create(
	ctx context.Context,
	note model.Note,
) (*model.Note, error) {

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO notes (id, user_id, title, content)
		 VALUES ($1, $2, $3, $4)`,
		note.ID,
		note.UserID,
		note.Title,
		note.Content,
	)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *PostgresNoteRepository) GetAll(
	ctx context.Context,
	userID string,
) ([]model.Note, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, title, content
		 FROM notes WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []model.Note

	for rows.Next() {
		var n model.Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	return notes, rows.Err()
}

func (r *PostgresNoteRepository) GetByID(
	ctx context.Context,
	id string,
	userID string,
) (*model.Note, error) {

	var n model.Note

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, title, content FROM notes WHERE id = $1 AND user_id = $2`,
		id,
		userID,
	).Scan(&n.ID, &n.UserID, &n.Title, &n.Content)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoteNotFound
	}
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *PostgresNoteRepository) Update(
	ctx context.Context,
	note model.Note,
	userID string,
) (*model.Note, error) {

	result, err := r.db.ExecContext(
		ctx,
		`UPDATE notes SET title = $1, content = $2 WHERE id = $3 AND user_id = $4`,
		note.Title,
		note.Content,
		note.ID,
		userID,
	)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, ErrNoteNotFound
	}

	note.UserID = userID
	return &note, nil
}

func (r *PostgresNoteRepository) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {

	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM notes WHERE id = $1 AND user_id = $2`,
		id,
		userID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoteNotFound
	}

	return nil
}

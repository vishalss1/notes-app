package repository

import (
	"context"
	"database/sql"
	"errors"

	"notes-app/internal/model"
)

type PostgresNoteRepository struct {
	db *sql.DB
}

func NewPostgresNoteRepository(db *sql.DB) *PostgresNoteRepository {
	return &PostgresNoteRepository{
		db: db,
	}
}

func (r *PostgresNoteRepository) Create(
	ctx context.Context,
	note model.Note,
) (model.Note, error) {

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO notes (title, content)
		 VALUES ($1, $2)
		 RETURNING id`,
		note.Title,
		note.Content,
	).Scan(&note.ID)

	if err != nil {
		return model.Note{}, err
	}

	return note, nil
}

func (r *PostgresNoteRepository) GetAll(
	ctx context.Context,
) ([]model.Note, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, title, content FROM notes`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []model.Note

	for rows.Next() {
		var n model.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *PostgresNoteRepository) GetByID(
	ctx context.Context,
	id int,
) (model.Note, error) {

	var n model.Note

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, title, content
		 FROM notes
		 WHERE id = $1`,
		id,
	).Scan(&n.ID, &n.Title, &n.Content)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Note{}, sql.ErrNoRows
	}

	if err != nil {
		return model.Note{}, err
	}

	return n, nil
}

func (r *PostgresNoteRepository) Update(
	ctx context.Context,
	note model.Note,
) error {

	result, err := r.db.ExecContext(
		ctx,
		`UPDATE notes
		 SET title = $1,
		     content = $2
		 WHERE id = $3`,
		note.Title,
		note.Content,
		note.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresNoteRepository) Delete(
	ctx context.Context,
	id int,
) error {

	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM notes WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

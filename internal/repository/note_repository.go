package repository

import (
	"context"
	"notes-app/internal/model"
)

type NoteRepository interface {
	GetAll(ctx context.Context) ([]model.Note, error)
	GetByID(ctx context.Context, id int) (model.Note, error)
	Create(ctx context.Context, note model.Note) (model.Note, error)
	Update(ctx context.Context, note model.Note) error
	Delete(ctx context.Context, id int) error
}

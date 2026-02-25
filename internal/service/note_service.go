package service

import (
	"context"
	"errors"
	"notes-app/internal/model"
	"notes-app/internal/repository"
)

type NoteService struct {
	repo repository.NoteRepository
}

func NewNoteService(r repository.NoteRepository) *NoteService {
	return &NoteService{
		repo: r,
	}
}

func (s *NoteService) Create(
	ctx context.Context,
	note model.Note,
) (model.Note, error) {

	if note.Title == "" {
		return model.Note{}, errors.New("title cannot be empty")
	}

	return s.repo.Create(ctx, note)
}

func (s *NoteService) GetAll(
	ctx context.Context,
) ([]model.Note, error) {

	return s.repo.GetAll(ctx)
}

func (s *NoteService) GetByID(
	ctx context.Context,
	id int,
) (model.Note, error) {

	if id <= 0 {
		return model.Note{}, errors.New("invalid id")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *NoteService) Update(
	ctx context.Context,
	note model.Note,
) error {

	if note.ID <= 0 {
		return errors.New("invalid id")
	}

	if note.Title == "" {
		return errors.New("title cannot be empty")
	}

	return s.repo.Update(ctx, note)
}

func (s *NoteService) Delete(
	ctx context.Context,
	id int,
) error {

	if id <= 0 {
		return errors.New("invalid id")
	}

	return s.repo.Delete(ctx, id)
}

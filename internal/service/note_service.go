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
	return &NoteService{repo: r}
}

func (s *NoteService) Create(
	ctx context.Context,
	note model.Note,
) (*model.Note, error) {

	if note.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	if note.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	return s.repo.Create(ctx, note)
}

func (s *NoteService) GetAll(
	ctx context.Context,
	userID string,
) ([]model.Note, error) {

	return s.repo.GetAll(ctx, userID)
}

func (s *NoteService) GetByID(
	ctx context.Context,
	id string,
	userID string,
) (*model.Note, error) {

	if id == "" {
		return nil, errors.New("invalid id")
	}

	return s.repo.GetByID(ctx, id, userID)
}

func (s *NoteService) Update(
	ctx context.Context,
	note model.Note,
	userID string,
) (*model.Note, error) {

	if note.ID == "" {
		return nil, errors.New("invalid id")
	}

	if note.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	return s.repo.Update(ctx, note, userID)
}

func (s *NoteService) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {

	if id == "" {
		return errors.New("invalid id")
	}

	return s.repo.Delete(ctx, id, userID)
}

package repository

import (
	"context"
	"notes-app/internal/model"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user model.User) (*model.User, error)
}

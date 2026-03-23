package repository

import (
	"context"
	"notes-app/internal/model"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token model.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}

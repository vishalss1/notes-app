package service

import (
	"context"
	"errors"
	"time"

	"notes-app/internal/auth"
	"notes-app/internal/model"
	"notes-app/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtSecret        string
}

func NewUserService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtSecret string,
) *UserService {
	return &UserService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
	}
}

func (s *UserService) Signup(
	ctx context.Context,
	email string,
	password string,
) (*AuthTokens, error) {

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: string(hashed),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return s.generateTokenPair(ctx, createdUser.ID)
}

func (s *UserService) Login(
	ctx context.Context,
	email string,
	password string,
) (*AuthTokens, error) {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokenPair(ctx, user.ID)
}

func (s *UserService) RefreshToken(
	ctx context.Context,
	refreshTokenStr string,
) (*AuthTokens, error) {

	if refreshTokenStr == "" {
		return nil, errors.New("refresh token is required")
	}

	storedToken, err := s.refreshTokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		_ = s.refreshTokenRepo.DeleteByToken(ctx, refreshTokenStr)
		return nil, errors.New("refresh token expired")
	}

	// Delete old refresh token (rotation)
	_ = s.refreshTokenRepo.DeleteByToken(ctx, refreshTokenStr)

	return s.generateTokenPair(ctx, storedToken.UserID)
}

func (s *UserService) Logout(
	ctx context.Context,
	refreshTokenStr string,
) error {

	if refreshTokenStr == "" {
		return errors.New("refresh token is required")
	}

	return s.refreshTokenRepo.DeleteByToken(ctx, refreshTokenStr)
}

func (s *UserService) generateTokenPair(
	ctx context.Context,
	userID string,
) (*AuthTokens, error) {

	accessToken, err := auth.GenerateAccessToken(userID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken := model.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(auth.RefreshTokenExpiry),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}

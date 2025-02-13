package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/koccyx/avito_assignment/internal/storage"
)

var (
	ErrNoEntry = errors.New("no entry found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUsername = errors.New("invalid password")
	ErrInvalidToken = errors.New("invalid token")
)

type AuthService interface {
	Auth(ctx context.Context, username, password string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
}

type Service struct {
	Auth AuthService
}

func New(repo *storage.Repository, log *slog.Logger, secret string) *Service {
	return &Service{
		Auth: NewAuthService(repo.User, log, secret),
	}
}
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/koccyx/avito_assignment/internal/lib/jwt"
	"github.com/koccyx/avito_assignment/internal/lib/sl"
	"github.com/koccyx/avito_assignment/internal/validators"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceSt struct {
	log *slog.Logger
	userRepo UserRepository
	secret string
}

func (s *AuthServiceSt) Auth(ctx context.Context, username, password string) (string, error) {
	const op = "service.auth.Auth"

	log := s.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("auth for user")

	usr, err := s.userRepo.GetByName(ctx, username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error("error while getting user", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, err)
		}

		log.Info("user not found, starting registration")

		err = validators.ValidateUsername(username)
		if err != nil {
			log.Error("error while validating username", sl.Err(err))
			return "", ErrInvalidUsername
		}

		err = validators.ValidatePassword(password)
		if err != nil {
			log.Error("error while validating password", sl.Err(err))
			return "", ErrInvalidPassword
		}

		pswdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("error while hashing", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, err)
		}

		usr, err = s.userRepo.Create(ctx, username, string(pswdHash))
		if err != nil {
			log.Error("failed while creating a user", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, err)
		}

		log.Info("user was registred")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		log.Error("invalid credentials", sl.Err(err))

		return "", ErrInvalidCredentials
	}

	token, err := jwt.NewToken(usr.Id.String(), s.secret)
	if err != nil {
		log.Error("failed during token creation", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfuly return token")

	return token, nil
}

func (s * AuthServiceSt) VerifyToken(ctx context.Context, token string) (string, error) {
	const op = "service.auth.VerifyToken"

	log := s.log.With(
		slog.String("op", op),
		slog.String("token", token),
	)

	usrId, err := jwt.ParseToken(token, s.secret)
	if err != nil {
		log.Error("failed during token validation", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidToken)
	}

	prsdUsrId, err := uuid.Parse(usrId)
	if err != nil {
		log.Error("failed during user id parsing", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	usr, err := s.userRepo.GetOne(ctx, prsdUsrId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			s.log.Error("error while getting user", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, err)
		}
		s.log.Error("error while getting user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrNoEntry)
	}

	return usr.Id.String(), nil
}

func NewAuthService(tdRepo UserRepository, logger *slog.Logger, secret string) *AuthServiceSt{
	return &AuthServiceSt{
		userRepo: tdRepo,
		log: logger,
		secret: secret,
	}
}
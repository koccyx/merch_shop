package handlers

import (
	"context"
	"errors"

	"github.com/koccyx/avito_assignment/internal/server/models"
)

var (
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrInvalidParam        = errors.New("invalid request param")
	ErrInternalServerError = errors.New("internal server error")
	ErrInternalUserId      = errors.New("invalid user id")
	ErrUserNotFound        = errors.New("invalid request body")
	ErrNotEnoughBalance    = errors.New("not enough balance for opperation")
	ErrSameUserTransfer    = errors.New("cant transfer coins to same user")
)

type UserService interface {
	Info(ctx context.Context, userId string) (*models.InfoResponse, error)
	TransferCoins(ctx context.Context, userFromId, userToId string, amount int) error
}

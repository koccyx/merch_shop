package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/server/models"
	"github.com/koccyx/avito_assignment/internal/storage"
)

var (
	ErrNoEntry = errors.New("no entry found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUsername = errors.New("invalid password")
	ErrSameUserTransfer = errors.New("cant transfer coins to same user")
	ErrInvalidToken = errors.New("invalid token")
	ErrNotEnoughBalance = errors.New("not enough balance")
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, username, password string) (*uuid.UUID, error)
	GetOne(ctx context.Context, usrId uuid.UUID) (*entities.User, error)
	GetByName(ctx context.Context, username string) (*entities.User, error)
	PutCoins(ctx context.Context, tx *sql.Tx, userId uuid.UUID, amount int) (int, error)
	GetUserItemsInfo(ctx context.Context, userId uuid.UUID) ([]entities.InventoryItem, error)
}

type UserItemRepository interface {
	Create(ctx context.Context, tx *sql.Tx, userId uuid.UUID, itemId uuid.UUID) (*uuid.UUID, error)
	GetAllInfoByUserId(ctx context.Context, usrId uuid.UUID) ([]entities.UserItem, error)
	GetOne(ctx context.Context, userItemId uuid.UUID) (*entities.UserItem, error)
}

type ItemRepository interface {
	GetByName(ctx context.Context, name string) (*entities.Item, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, fromUserId uuid.UUID, toUserId uuid.UUID, amount int) (*uuid.UUID, error)
	GetAllWithDirection(ctx context.Context, usrId uuid.UUID, direction entities.Direction) ([]entities.CoinTransactionInfo, error)
	GetOne(ctx context.Context, transactionId uuid.UUID) (*entities.Transaction, error)
}

type AuthService interface {
	Auth(ctx context.Context, username, password string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
}

type ItemService interface {
	PurchaseItem(ctx context.Context, userId, itemName string) error
}

type UserService interface {
	Info(ctx context.Context, userId string) (*models.InfoResponse, error)
	TransferCoins(ctx context.Context, userFromId, userToId string, amount int) error
}

type Service struct {
	Auth AuthService
	Item ItemService
	User UserService
}

func New(repo *storage.Repository, log *slog.Logger, secret string, db *sql.DB) *Service {
	return &Service{
		Auth: NewAuthService(repo.User, log, secret, db),
		Item: NewItemService(repo.User, repo.Item, repo.UserItem, log, db),
		User: NewUserService(repo.User, repo.Item, repo.UserItem, repo.Transaction, log, db),
	}
}
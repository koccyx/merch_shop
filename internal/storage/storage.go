package storage

import (
	"context"
	"database/sql"
	
	"github.com/google/uuid"

	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/storage/postgres"
)

type UserRepository interface {
	GetOne(ctx context.Context, usrId uuid.UUID) (*entities.User, error)
	Create(ctx context.Context, username, password string) (*entities.User, error)
	GetByName(ctx context.Context, username string) (*entities.User, error)
}

type Repository struct {
	User *postgres.UserRepository
}

func NewRepository(db *sql.DB) *Repository {
	user := postgres.NewUserRepository(db)

	return &Repository{
		User: user,
	}
}
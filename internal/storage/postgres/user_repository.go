package postgres

import (
	"context"
	"database/sql"
	"fmt"
	
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) Create(ctx context.Context, username, password string) (*entities.User, error) {	
	const op = "repo.postgres.user.Create"

	newUser := &entities.User{
		Id: uuid.New(),
		Username: username,
		Password: password,
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	
	sql, args, err := builder.Insert("users").
	Columns("id" , "username", "password").
	Values(newUser.Id, newUser.Username, newUser.Password).
	ToSql();

	if err != nil {
		return nil, fmt.Errorf("%s: building query error: %w", op, err)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := r.GetOne(ctx, newUser.Id)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return res, nil 
}

func (r *UserRepository) GetOne(ctx context.Context, usrId uuid.UUID) (*entities.User, error) {
	const op = "repo.postgres.user.GetOne" 

	if usrId == uuid.Nil {
		return nil, fmt.Errorf("%s: missing id", op)
	}
	
    builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    selectBuilder := builder.
	Select("id", "username", "password", "balance").
	From("users").
	Where(squirrel.Eq{"id": usrId})

    query, args, err := selectBuilder.ToSql()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	var user entities.User

    err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Username, &user.Password, &user.Balance)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &user, nil
}

func (r *UserRepository) GetByName(ctx context.Context, username string) (*entities.User, error) {
	const op = "repo.postgres.user.GetByName" 
	
    builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    selectBuilder := builder.
	Select("id", "username", "password", "balance").
	From("users").
	Where(squirrel.Eq{"username": username})

    query, args, err := selectBuilder.ToSql()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	var user entities.User

    err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Username, &user.Password, &user.Balance)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &user, nil
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

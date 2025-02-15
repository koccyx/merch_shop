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
	
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	
	sql, args, err := builder.Insert("users").
	Columns("username", "password").
	Values(username, password).
	ToSql();
	if err != nil {
		return nil, fmt.Errorf("%s: building query error: %w", op, err)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := r.GetByName(ctx, username)
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

func (r *UserRepository) GetUserItemsInfo(ctx context.Context, userId uuid.UUID) ([]entities.InventoryItem, error) {
	const op = "repo.postgres.user.GetUserItems" 
	
    builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    selectBuilder := builder.
	Select("i.name AS name", "COUNT(ui.item_id) AS amount").
	From("user_items ui").
	Join("items i ON ui.item_id = i.id").
	Where(squirrel.Eq{"ui.user_id": userId}).
	GroupBy("i.name")

    query, args, err := selectBuilder.ToSql()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer rows.Close()

    items := make([]entities.InventoryItem, 0) 

    for rows.Next() {
        var item entities.InventoryItem
        if err := rows.Scan(&item.Name, &item.Amount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
        }
        items = append(items, item)
    }

    return items, nil
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

func (r *UserRepository) PutCoins(ctx context.Context, userId uuid.UUID, amount int) (int64, error) {
	const op = "repo.postgres.user.PutCoins" 
	
    usr, err := r.GetOne(ctx, userId)
	if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

	if usr.Balance + amount < 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrNotEnoughBalance)
	}
	
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    putBuilder := builder.
		Update("users").
		Set("balance", usr.Balance + amount).
		Where(squirrel.Eq{"id": usr.Id})

    query, args, err := putBuilder.ToSql()
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

    return rowsAffected, nil
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

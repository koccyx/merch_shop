package postgres

import (
	"context"
	"database/sql"
	"fmt"
	
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserItemsRepository struct {
	db *sql.DB
}

func (r *UserItemsRepository) Create(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) (*entities.UserItem, error) {	
	const op = "repo.postgres.userItem.Create" 

	if itemId == uuid.Nil {
		return nil, fmt.Errorf("%s: missing item id", op)
	}
	if userId == uuid.Nil {
		return nil, fmt.Errorf("%s: missing item id", op)
	}

	id := uuid.New()

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    createBuilder := builder.
	Insert("user_items").
	Columns("id","user_id", "item_id").
	Values(id, userId, itemId)

    sql, args, err := createBuilder.ToSql()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
	
	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userItem, err := r.GetOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
    return userItem, nil
}


func (r *UserItemsRepository) GetAllInfoByUserId(ctx context.Context, usrId uuid.UUID) ([]entities.UserItem, error) {
	const op = "repo.postgres.user.GetAll" 

	if usrId == uuid.Nil {
		return nil, fmt.Errorf("%s: missing id", op)
	}
	
    builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    selectBuilder := builder.
	Select("id", "user_id", "item_id", "created_at").
	From("user_items").
	Where(squirrel.Eq{"user_id": "usrId"}).
	OrderBy("created_at DESC")
	

    query, args, err := selectBuilder.ToSql()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	userItems := make([]entities.UserItem, 0)
	
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
    }
	
	for rows.Next() {
		var userItem entities.UserItem
		err := rows.Scan(&userItem.Id, &userItem.UserId,  &userItem.ItemId, &userItem.Created_at)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		userItems = append(userItems, userItem)
	}
    return userItems, nil
}

func (r *UserItemsRepository) GetOne(ctx context.Context, userItemId uuid.UUID) (*entities.UserItem, error) {
	const op = "repo.postgres.userItem.GetOne" 
	
    builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    selectBuilder := builder.
	Select("id", "user_id", "item_id", "created_at").
	From("user_items").
	Where(squirrel.Eq{"id": userItemId}).
	OrderBy("created_at DESC")

    query, args, err := selectBuilder.ToSql()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	var userItem entities.UserItem

    err = r.db.QueryRowContext(ctx, query, args...).Scan(&userItem.Id, &userItem.UserId, &userItem.ItemId, &userItem.Created_at)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &userItem, nil
}

func NewUserItemRepository(db *sql.DB) *UserItemsRepository {
	return &UserItemsRepository{
		db: db,
	}
}

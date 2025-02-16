package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserItem struct {
	Id         uuid.UUID `db:"id"`
	UserId     uuid.UUID `db:"user_id"`
	ItemId     uuid.UUID `db:"item_id"`
	Created_at time.Time `db:"created_at"`
}

package entities

import (
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
	Balance  int       `db:"balance"`
}

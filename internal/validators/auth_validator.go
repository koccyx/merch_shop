package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/koccyx/avito_assignment/internal/http/models"
)

func ValidateUsername(username string) error {
	if len(username) < 4 {
		return fmt.Errorf("username must be more then 4")
	}

	return nil
}


func ValidatePassword(password string) error {
	if len(password) < 5 {
		return fmt.Errorf("password must be more then 5")
	}
	
	return nil
}

func ValdateAuthRequest(r models.AuthRequest) error {
	v := validator.New(validator.WithRequiredStructEnabled())

	return v.Struct(r)
}

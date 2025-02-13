package handlers

import (
	"context"
	// "encoding/json"
	// "errors"
	"fmt"
	"log/slog"
	"net/http"

	// "github.com/koccyx/avito_assignment/internal/http/models"
	// jsonwriter "github.com/koccyx/avito_assignment/internal/lib/json_writer"
	// "github.com/koccyx/avito_assignment/internal/lib/sl"
	// "github.com/koccyx/avito_assignment/internal/service"
	// "github.com/koccyx/avito_assignment/internal/validators"
)

type ItemService interface {
	Auth(ctx context.Context, username, password string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
}

func Merch(item ItemService, logger *slog.Logger) http.HandlerFunc {
	const op = "http.handlers.Merch" 
	// log := logger.With(
	// 	slog.String("op", op),
	// )
	
	fn := func(w http.ResponseWriter, r *http.Request) {
		
		res := r.Context().Value("userId")
		fmt.Println(res)
	}

	return http.HandlerFunc(fn)
}
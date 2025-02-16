package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	jsonwriter "github.com/koccyx/avito_assignment/internal/lib/json_writer"
	"github.com/koccyx/avito_assignment/internal/lib/sl"
	"github.com/koccyx/avito_assignment/internal/server/models"
	"github.com/koccyx/avito_assignment/internal/service"
	"github.com/koccyx/avito_assignment/internal/validators"
)

type AuthService interface {
	Auth(ctx context.Context, username, password string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
}

func Auth(authService AuthService, logger *slog.Logger) http.HandlerFunc {
	const op = "http.handlers.Auth"
	log := logger.With(
		slog.String("op", op),
	)

	fn := func(w http.ResponseWriter, r *http.Request) {
		var authReq models.AuthRequest

		err := json.NewDecoder(r.Body).Decode(&authReq)
		if err != nil {
			log.Error("wrong request body", sl.Err(err))
			err = jsonwriter.WriteJSONError(ErrInvalidRequestBody, w, http.StatusBadRequest)
			if err != nil {
				log.Error("json error", sl.Err(err))
			}
			return
		}

		err = validators.ValdateAuthRequest(authReq)
		if err != nil {
			log.Error("wrong request body", sl.Err(err))
			jsonwriter.WriteJSONError(ErrInvalidRequestBody, w, http.StatusBadRequest)
			return
		}

		token, err := authService.Auth(r.Context(), authReq.Username, authReq.Password)
		if err != nil {
			if errors.Is(err, service.ErrInvalidPassword) || errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrInvalidCredentials) {
				log.Error("validation", sl.Err(err))
				err := jsonwriter.WriteJSONError(err, w, http.StatusBadRequest)

				if err != nil {
					log.Error("json error", sl.Err(err))
				}

				return
			}

			log.Error("error while getting token", sl.Err(err))
			jsonwriter.WriteJSONError(ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

		jsonwriter.WriteJSON(&models.AuthResponse{
			Token: token,
		}, w)
	}

	return http.HandlerFunc(fn)
}

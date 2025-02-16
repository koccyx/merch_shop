package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	jsonwriter "github.com/koccyx/avito_assignment/internal/lib/json_writer"
	"github.com/koccyx/avito_assignment/internal/lib/sl"
	"github.com/koccyx/avito_assignment/internal/service"
)

type AuthService interface {
	VerifyToken(ctx context.Context, token string) (string, error)
}

func AuthMiddleware(auth AuthService, secret string, log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/auth"),
		)

		log.Info("auth middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			authHdr := r.Header.Get("Authorization")
			if authHdr == "" {
				log.Error("empty authorization header field")
				err := jsonwriter.WriteJSONError(fmt.Errorf("empty authorization header field"), w, http.StatusUnauthorized)

				if err != nil {
					log.Error("json error", sl.Err(err))
					return
				}

				return
			}

			splitedAuthHdr := strings.Split(authHdr, " ")

			if len(splitedAuthHdr) != 2 || splitedAuthHdr[0] != "Bearer" {
				log.Error("wrong authorization header")
				err := jsonwriter.WriteJSONError(fmt.Errorf("wrong authorization header"), w, http.StatusUnauthorized)

				if err != nil {
					log.Error("json error", sl.Err(err))
					return
				}

				return
			}

			usrId, err := auth.VerifyToken(r.Context(), splitedAuthHdr[1])
			if err != nil {
				if errors.Is(err, service.ErrNoEntry) {
					log.Error("no user with this id")
					err := jsonwriter.WriteJSONError(fmt.Errorf("no user with this id"), w, http.StatusUnauthorized)

					if err != nil {
						log.Error("json error")
						return
					}
					return
				}

				log.Error("wrong token format")
				jsonwriter.WriteJSONError(fmt.Errorf("wrong token format"), w, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", usrId)
			req := r.WithContext(ctx)

			next.ServeHTTP(w, req)
		}

		return http.HandlerFunc(fn)
	}
}

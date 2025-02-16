package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	jsonwriter "github.com/koccyx/avito_assignment/internal/lib/json_writer"
	"github.com/koccyx/avito_assignment/internal/lib/sl"
	"github.com/koccyx/avito_assignment/internal/service"
)

type ItemService interface {
	PurchaseItem(ctx context.Context, userId, itemName string) error
}

func PurchaseItem(itemService ItemService, logger *slog.Logger) http.HandlerFunc {
	const op = "http.handlers.PurchaseItem"
	log := logger.With(
		slog.String("op", op),
	)

	fn := func(w http.ResponseWriter, r *http.Request) {
		itemName := chi.URLParam(r, "item")
		if itemName == "" {
			log.Error("empty parameter", sl.Err(ErrInvalidParam))
			jsonwriter.WriteJSONError(ErrInvalidParam, w, http.StatusBadRequest)
			return
		}

		usrId := r.Context().Value("userId")
		userIdStr, ok := usrId.(string)
		if !ok {
			log.Error("no user", sl.Err(ErrInvalidParam))
			jsonwriter.WriteJSONError(ErrInternalUserId, w, http.StatusInternalServerError)
			return
		}

		err := itemService.PurchaseItem(r.Context(), userIdStr, itemName)
		if err != nil {
			if errors.Is(err, service.ErrNotEnoughBalance) {
				log.Error("not enough balance", sl.Err(err))
				jsonwriter.WriteJSONError(ErrNotEnoughBalance, w, http.StatusBadRequest)
				return
			}

			if errors.Is(err, service.ErrNoEntry) {
				log.Error("not item found", sl.Err(err))
				jsonwriter.WriteJSONError(ErrInvalidParam, w, http.StatusBadRequest)
				return
			}

			log.Error("error while purchasing item", sl.Err(err))
			jsonwriter.WriteJSONError(ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

		jsonwriter.WriteSuccess(w)
	}

	return http.HandlerFunc(fn)
}

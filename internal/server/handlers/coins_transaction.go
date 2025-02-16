package handlers

import (
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

func TransferCoins(userService UserService, logger *slog.Logger) http.HandlerFunc {
	const op = "http.handlers.TransferCoins"
	log := logger.With(
		slog.String("op", op),
	)

	fn := func(w http.ResponseWriter, r *http.Request) {
		var coinsTransactionReq models.SendCoinRequest

		err := json.NewDecoder(r.Body).Decode(&coinsTransactionReq)
		if err != nil {
			log.Error("wrong request body", sl.Err(err))
			jsonwriter.WriteJSONError(ErrInvalidRequestBody, w, http.StatusBadRequest)
			return
		}

		err = validators.ValdateCoinsTransactionRequest(coinsTransactionReq)
		if err != nil {
			log.Error("wrong request body", sl.Err(err))
			jsonwriter.WriteJSONError(ErrInvalidRequestBody, w, http.StatusBadRequest)
			return
		}

		err = validators.ValidateCoinsAmount(coinsTransactionReq.Amount)
		if err != nil {
			log.Error("wrong amount", sl.Err(err))
			jsonwriter.WriteJSONError(ErrInvalidRequestBody, w, http.StatusBadRequest)
			return
		}

		usrId := r.Context().Value("userId")
		userIdStr, ok := usrId.(string)
		if !ok {
			log.Error("empty param", sl.Err(ErrInvalidParam))
			jsonwriter.WriteJSONError(ErrInternalUserId, w, http.StatusInternalServerError)
			return
		}

		err = userService.TransferCoins(r.Context(), userIdStr, coinsTransactionReq.ToUser, coinsTransactionReq.Amount)
		if err != nil {
			if errors.Is(err, service.ErrNoEntry) {
				log.Error("validation", sl.Err(err))
				jsonwriter.WriteJSONError(ErrUserNotFound, w, http.StatusBadRequest)
				return
			}

			if errors.Is(err, service.ErrNotEnoughBalance) {
				log.Error("validation", sl.Err(err))
				jsonwriter.WriteJSONError(ErrNotEnoughBalance, w, http.StatusBadRequest)
				return
			}

			if errors.Is(err, service.ErrInvalidUsername) {
				log.Error("validation", sl.Err(err))
				jsonwriter.WriteJSONError(ErrNotEnoughBalance, w, http.StatusBadRequest)
				return
			}

			if errors.Is(err, service.ErrSameUserTransfer) {
				log.Error("cant transfer coins to same user", sl.Err(err))
				jsonwriter.WriteJSONError(ErrSameUserTransfer, w, http.StatusBadRequest)
				return
			}

			log.Error("error while transfering coins", sl.Err(err))
			jsonwriter.WriteJSONError(ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

		jsonwriter.WriteSuccess(w)
	}

	return http.HandlerFunc(fn)
}

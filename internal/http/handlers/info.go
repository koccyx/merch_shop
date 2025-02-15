package handlers

import (
	"log/slog"
	"net/http"

	jsonwriter "github.com/koccyx/avito_assignment/internal/lib/json_writer"
	"github.com/koccyx/avito_assignment/internal/lib/sl"
)


func Info(userService UserService, logger *slog.Logger) http.HandlerFunc {
	const op = "http.handlers.Info" 
	log := logger.With(
		slog.String("op", op),
	)
	
	fn := func(w http.ResponseWriter, r *http.Request) {
		usrId := r.Context().Value("userId")
		userIdStr, ok := usrId.(string)
		if !ok {
			log.Error("empty param", sl.Err(ErrInvalidParam))
			jsonwriter.WriteJSONError(ErrInternalUserId, w, http.StatusInternalServerError)
			return
		} 

		info, err := userService.Info(r.Context(), userIdStr)
		if err != nil {
			log.Error("error getting user info", sl.Err(err))
			jsonwriter.WriteJSONError(ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

		jsonwriter.WriteJSON(info, w)
	}

	return http.HandlerFunc(fn)
}